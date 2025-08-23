package service

import (
	"context"
	"errors"
	"time"

	"ecommerce/entity"
	"ecommerce/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CartService interface {
	AddItem(ctx context.Context, userID uuid.UUID, productID uint, quantity int) (*entity.CartItem, error)
	RemoveItem(ctx context.Context, userID uuid.UUID, itemID uuid.UUID) error
	GetCart(ctx context.Context, userID uuid.UUID) (*entity.Cart, error)
}

type cartService struct {
	cartRepo    repository.CartRepository
	productRepo repository.ProductRepository
}

func NewCartService(cartRepo repository.CartRepository, productRepo repository.ProductRepository) CartService {
	return &cartService{cartRepo: cartRepo, productRepo: productRepo}
}

func (s *cartService) AddItem(ctx context.Context, userID uuid.UUID, productID uint, quantity int) (*entity.CartItem, error) {
	if quantity <= 0 {
		return nil, errors.New("invalid quantity")
	}

	product, err := s.productRepo.GetByProductID(ctx, productID)
	if err != nil {
		return nil, errors.New("failed to fetch product: " + err.Error())
	}
	if product == nil {
		return nil, errors.New("product not found")
	}
	if product.Stock == 0 {
		return nil, errors.New("product out of stock")
	}

	// get/create cart
	cart, err := s.cartRepo.GetCartByUserID(ctx, userID)
	if err != nil {
		return nil, errors.New("failed to fetch cart: " + err.Error())
	}

	existing, err := s.cartRepo.GetItemByCartAndProduct(ctx, cart.ID, productID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("failed to fetch cart item: " + err.Error())
	}

	if existing != nil && existing.ID != uuid.Nil {
		newQty := existing.Quantity + quantity
		if uint(newQty) > product.Stock {
			return nil, errors.New("insufficient stock")
		}
		existing.Quantity = newQty
		if err := s.cartRepo.UpdateItem(ctx, existing); err != nil {
			return nil, errors.New("failed to update cart item: " + err.Error())
		}
		return existing, nil
	}

	if uint(quantity) > product.Stock {
		return nil, errors.New("insufficient stock")
	}
	item := &entity.CartItem{
		ID:        uuid.New(),
		CartID:    cart.ID,
		ProductID: productID,
		Quantity:  quantity,
	}
	if err := s.cartRepo.CreateItem(ctx, item); err != nil {
		return nil, errors.New("failed to create cart item: " + err.Error())
	}
	return item, nil
}

func (s *cartService) RemoveItem(ctx context.Context, userID uuid.UUID, itemID uuid.UUID) error {
	cart, err := s.cartRepo.GetCartByUserID(ctx, userID)
	if err != nil {
		return errors.New("failed to fetch cart: " + err.Error())
	}

	item, err := s.cartRepo.GetItemByID(ctx, itemID)
	if err != nil {
		return errors.New("failed to fetch cart item: " + err.Error())
	}
	if item == nil {
		return errors.New("cart item not found")
	}
	if item.CartID != cart.ID {
		return errors.New("cart item not found")
	}

	if err := s.cartRepo.DeleteItem(ctx, itemID); err != nil {
		return errors.New("failed to delete cart item: " + err.Error())
	}
	return nil
}

func (s *cartService) GetCart(ctx context.Context, userID uuid.UUID) (*entity.Cart, error) {
	cart, err := s.cartRepo.GetCartWithItems(ctx, userID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("failed to fetch cart: " + err.Error())
	}
	if cart == nil {
		// create an empty cart for validity when GET cart
		cart, err = s.cartRepo.GetCartByUserID(ctx, userID)
		if err != nil {
			return nil, errors.New("failed to fetch cart: " + err.Error())
		}
	}
	cart.UpdatedAt = time.Now()
	return cart, nil
}
