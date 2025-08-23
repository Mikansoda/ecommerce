package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"ecommerce/entity"
	"ecommerce/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentService interface {
	CreatePayment(ctx context.Context, order *entity.Order, invoiceID string) (*entity.Payment, error)
	GetOrderByID(ctx context.Context, orderID uuid.UUID) (*entity.Order, error)
	GetPaymentsByUserID(ctx context.Context, userID string) ([]entity.Payment, error)
	GetAllPayments(ctx context.Context) ([]entity.Payment, error)
	AutoCancelPendingPayments()
	UpdatePaymentStatus(ctx context.Context, invoiceID string, status string) error
}

type paymentService struct {
	paymentRepo  repository.PaymentRepository
	orderRepo    repository.OrderRepository
	productSvc   ProductService
	actionLogSvc ActionLogService
	db           *gorm.DB
}

func NewPaymentService(paymentRepo repository.PaymentRepository, orderRepo repository.OrderRepository, productSvc ProductService, actionLogSvc ActionLogService, db *gorm.DB) PaymentService {
	return &paymentService{paymentRepo, orderRepo, productSvc, actionLogSvc, db}
}

// Create payment from Xendit invoice
func (s *paymentService) CreatePayment(ctx context.Context, order *entity.Order, invoiceID string) (*entity.Payment, error) {
	if order == nil {
		return nil, errors.New("order is nil")
	}

	payment := &entity.Payment{
		ID:          uuid.New(),
		OrderID:     order.ID,
		InvoiceID:   invoiceID,
		PaymentType: "xendit_invoice",
		Status:      "pending",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.paymentRepo.Create(ctx, payment); err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	return payment, nil
}

func (s *paymentService) GetOrderByID(ctx context.Context, orderID uuid.UUID) (*entity.Order, error) {
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, fmt.Errorf("order not found")
	}
	return order, nil
}

func (s *paymentService) GetPaymentsByUserID(ctx context.Context, userID string) ([]entity.Payment, error) {
	var payments []entity.Payment
	err := s.paymentRepo.GetByUserID(ctx, userID, &payments)
	if err != nil {
		return nil, err
	}
	return payments, nil
}

func (s *paymentService) GetAllPayments(ctx context.Context) ([]entity.Payment, error) {
	var payments []entity.Payment
	err := s.paymentRepo.GetPayments(ctx, &payments)
	if err != nil {
		return nil, err
	}
	return payments, nil
}

// Auto cancel pending payment > 24 hours
func (s *paymentService) AutoCancelPendingPayments() {
	ctx := context.Background()
	payments, _ := s.paymentRepo.GetPendingOlderThan(ctx, 1)
	for _, p := range payments {
		tx := s.db.Begin()
		if tx.Error != nil {
			continue
		}

		p.Status = "failed"
		p.UpdatedAt = time.Now()
		s.paymentRepo.Update(ctx, &p)

		// cancel order
		order, err := s.orderRepo.GetByIDForUpdate(ctx, p.OrderID, tx)
		if err != nil {
			tx.Rollback()
			continue
		}
		order.Status = "cancelled"
		order.UpdatedAt = time.Now()
		s.orderRepo.Update(ctx, order, tx)

		// rollback stock
		for _, item := range order.OrderItems {
			prod, err := s.productSvc.GetProductByID(ctx, item.ProductID)
			if err != nil {
				continue
			}
			prod.Stock += uint(item.Quantity)
			s.productSvc.UpdateProduct(ctx, prod)
		}

		tx.Commit()
	}
}

// Update payment status from Xendit webhook
func (s *paymentService) UpdatePaymentStatus(ctx context.Context, invoiceID string, status string) error {
	payment, err := s.paymentRepo.GetByInvoiceID(ctx, invoiceID)
	if err != nil {
		return fmt.Errorf("payment not found: %w", err)
	}

	payment.Status = status
	payment.UpdatedAt = time.Now()
	if err := s.paymentRepo.Update(ctx, payment); err != nil {
		return err
	}

	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	order, err := s.orderRepo.GetByIDForUpdate(ctx, payment.OrderID, tx)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("order not found: %w", err)
	}

	if status == "PAID" || status == "paid" {
		order.Status = "paid"
		order.UpdatedAt = time.Now()
		s.orderRepo.Update(ctx, order, tx)

		var actorID *uuid.UUID
		if order.UserID != uuid.Nil {
			actorID = &order.UserID
		}

		for _, item := range order.OrderItems {
			_ = s.actionLogSvc.Log(
				ctx,
				"buyer",
				actorID,
				"sold",
				"products",
				item.ProductID,
			)
		}
	}

	if status == "FAILED" || status == "failed" {
		order.Status = "cancelled"
		order.UpdatedAt = time.Now()
		s.orderRepo.Update(ctx, order, tx)

		for _, item := range order.OrderItems {
			prod, err := s.productSvc.GetProductByID(ctx, item.ProductID)
			if err != nil {
				continue
			}
			prod.Stock += uint(item.Quantity)
			s.productSvc.UpdateProduct(ctx, prod)
		}
	}

	tx.Commit()
	return nil
}
