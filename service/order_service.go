package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"ecommerce/entity"
	"ecommerce/repository"
	"ecommerce/helper"

	"github.com/google/uuid"
)

type OrderService interface {
	CreateOrder(ctx context.Context, userID, addressID uuid.UUID) (*entity.Order, error)
	AutoCancelOrders()
	GetOrders(ctx context.Context, limit, offset int) ([]entity.Order, error)
	GetOrdersByUser(ctx context.Context, userID uuid.UUID) ([]entity.Order, error)
	UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status string) error
	GetOrdersByStatus(ctx context.Context, status string, limit, offset int) ([]entity.Order, error)
}

type orderService struct {
	cartRepo    repository.CartRepository
	productRepo repository.ProductRepository
	orderRepo   repository.OrderRepository
}

func NewOrderService(cartRepo repository.CartRepository, productRepo repository.ProductRepository, orderRepo repository.OrderRepository) OrderService {
	return &orderService{cartRepo, productRepo, orderRepo}
}

// create order from cart, hold stock, send email
func (s *orderService) CreateOrder(ctx context.Context, userID, addressID uuid.UUID) (*entity.Order, error) {
	// start transaction
	tx, err := s.productRepo.BeginTx(ctx)
	if err != nil {
		return nil, errors.New("failed to start transaction: " + err.Error())
	}
	defer func() {
		if r := recover(); r != nil {
			s.productRepo.RollbackTx(tx)
			panic(r)
		}
	}()

	// get cart
	cart, err := s.cartRepo.GetCartWithItems(ctx, userID)
	if err != nil {
		s.productRepo.RollbackTx(tx)
		return nil, errors.New("failed to fetch cart: " + err.Error())
	}
	if len(cart.Items) == 0 {
		s.productRepo.RollbackTx(tx)
		return nil, errors.New("cart is empty")
	}

	// check stock with row lock
	var outOfStock []string
	for _, item := range cart.Items {
		p, err := s.productRepo.GetByProductIDForUpdate(ctx, item.ProductID, tx)
		if err != nil {
			s.productRepo.RollbackTx(tx)
			return nil, errors.New("failed to fetch product for stock check: " + err.Error())
		}
		if uint(item.Quantity) > p.Stock {
			outOfStock = append(outOfStock, fmt.Sprintf("%s (available: %d)", p.Name, p.Stock))
		}
	}
	if len(outOfStock) > 0 {
		s.productRepo.RollbackTx(tx)
		return nil, fmt.Errorf("following products have insufficient stock: %v", outOfStock)
	}

	// create order
	var subtotal float64
	for _, it := range cart.Items {
		subtotal += it.Product.Price * float64(it.Quantity)
	}

	shippingFee := 20000.0
	total := subtotal + shippingFee
	now := time.Now()
	expire := now.Add(24 * time.Hour)

	order := &entity.Order{
		ID:          uuid.New(),
		UserID:      userID,
		AddressID:   addressID,
		Subtotal:    subtotal,
		ShippingFee: shippingFee,
		TotalAmount: total,
		Status:      "pending",
		CreatedAt:   now,
		UpdatedAt:   now,
		ExpiredAt:   &expire,
		OrderItems:  []entity.OrderItem{},
	}

	// assign order items + hold stock
	for _, it := range cart.Items {
		p, _ := s.productRepo.GetByProductIDForUpdate(ctx, it.ProductID, tx)
		p.Stock -= uint(it.Quantity)
		if err := s.productRepo.Update(ctx, p, tx); err != nil {
			s.productRepo.RollbackTx(tx)
			return nil, errors.New("failed to hold stock: " + err.Error())
		}

		order.OrderItems = append(order.OrderItems, entity.OrderItem{
			ID:           uuid.New(),
			OrderID:      order.ID,
			ProductID:    it.ProductID,
			Quantity:     it.Quantity,
			PriceAtOrder: it.Product.Price,
		})
	}

	if err := s.orderRepo.Create(ctx, order, tx); err != nil {
		s.productRepo.RollbackTx(tx)
		return nil, errors.New("failed to create order: " + err.Error())
	}

	s.productRepo.CommitTx(tx)

	// send email
	userEmail := ""
	username := "there"
	if cart.User != nil {
		userEmail = cart.User.Email
		username = cart.User.Username
	}
	subject := "Well Sprout - Order Activity"
	body := fmt.Sprintf(
		 "Hi %s,\n\n"+
            "Order %s has been made. Order Details:\n"+
            "Subtotal: %.2f\n"+
            "Shipping Fee: %.2f\n"+
            "Total: %.2f\n\n"+
            "Please make payment within 24 hours.\n\n"+
            "Thank you,\n"+
            "Well Sprout",
        username, order.ID, subtotal, shippingFee, total,
    )
	_ = helper.SendEmail(userEmail, subject, body)

	return order, nil
}

// auto cancel pending > 24 hours
func (s *orderService) AutoCancelOrders() {
	ctx := context.Background()
	orders, _ := s.orderRepo.GetPendingOrdersOlderThan(ctx, 24*time.Hour)
	for _, order := range orders {
		tx, err := s.productRepo.BeginTx(ctx)
		if err != nil {
			fmt.Println("failed to start transaction:", err)
			continue
		}

		o, _ := s.orderRepo.GetByIDForUpdate(ctx, order.ID, tx)
		o.Status = "cancelled"
		o.UpdatedAt = time.Now()
		s.orderRepo.Update(ctx, o, tx)

		// rollback stock
		for _, item := range o.OrderItems {
			p, _ := s.productRepo.GetByProductIDForUpdate(ctx, item.ProductID, tx)
			p.Stock += uint(item.Quantity)
			s.productRepo.Update(ctx, p, tx)
		}

		s.productRepo.CommitTx(tx)

		// send cancellation email
		userEmail := ""
		username := "there"
		if o.User != nil {
			userEmail = o.User.Email
			username = o.User.Username
		}
		subject := "Well Sprout - Order Cancelled"
		body := fmt.Sprintf(
			"Hi %s,\n\n"+
            "Order %s has been cancelled because it exceeded the 24-hour payment window.\n"+
            "If you have already made payment, please contact our support:\n\n"+ 
            "WhatsApp: +62 812 90909090"+
            "Thank you,\n"+
            "Well Sprout", 
            username, order.ID,
        )
		_ = helper.SendEmail(userEmail, subject, body)
	}
}

func (s *orderService) GetOrders(ctx context.Context, limit, offset int) ([]entity.Order, error) {
	return s.orderRepo.GetOrders(ctx, limit, offset)
}

func (s *orderService) GetOrdersByUser(ctx context.Context, userID uuid.UUID) ([]entity.Order, error) {
	return s.orderRepo.GetOrdersByUserID(ctx, userID)
}

func (s *orderService) GetOrdersByStatus(ctx context.Context, status string, limit, offset int) ([]entity.Order, error) {
	return s.orderRepo.GetByStatus(ctx, status, limit, offset)
}

func (s *orderService) UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status string) error {
    tx := s.orderRepo.BeginTx()
	if tx.Error != nil {
    return tx.Error
    }

    o, err := s.orderRepo.GetByIDForUpdate(ctx, orderID, tx)
    if err != nil {
        tx.Rollback()
        return err
    }

    o.Status = status
    o.UpdatedAt = time.Now()

    if err := s.orderRepo.Update(ctx, o, tx); err != nil {
        tx.Rollback()
        return err
    }

    return tx.Commit().Error
}
