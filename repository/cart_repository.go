package repository

import (
	"context"

	"ecommerce/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CartRepository interface {
	GetCartByUserID(ctx context.Context, userID uuid.UUID) (*entity.Cart, error)
	GetCartWithItems(ctx context.Context, userID uuid.UUID) (*entity.Cart, error)
	GetItemByID(ctx context.Context, id uuid.UUID) (*entity.CartItem, error)
	GetItemByCartAndProduct(ctx context.Context, cartID uuid.UUID, productID uint) (*entity.CartItem, error)
	CreateItem(ctx context.Context, item *entity.CartItem) error
	UpdateItem(ctx context.Context, item *entity.CartItem) error
	DeleteItem(ctx context.Context, id uuid.UUID) error
}

type cartRepo struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) CartRepository {
	return &cartRepo{db: db}
}

func (r *cartRepo) GetCartByUserID(ctx context.Context, userID uuid.UUID) (*entity.Cart, error) {
	var cart entity.Cart
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		First(&cart).Error
	if err == gorm.ErrRecordNotFound {
		cart = entity.Cart{
			ID:     uuid.New(),
			UserID: userID,
		}
		if err2 := r.db.WithContext(ctx).Create(&cart).Error; err2 != nil {
			return nil, err2
		}
		return &cart, nil
	}
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *cartRepo) GetCartWithItems(ctx context.Context, userID uuid.UUID) (*entity.Cart, error) {
	var cart entity.Cart
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Preload("User").
		Preload("Items").
		Preload("Items.Product").
		Preload("Items.Product.Categories").
		Preload("Items.Product.Images").
		First(&cart).Error
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *cartRepo) GetItemByID(ctx context.Context, id uuid.UUID) (*entity.CartItem, error) {
	var item entity.CartItem
	if err := r.db.WithContext(ctx).
		Preload("Product").
		First(&item, "id = ?", id).
		Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *cartRepo) GetItemByCartAndProduct(ctx context.Context, cartID uuid.UUID, productID uint) (*entity.CartItem, error) {
	var item entity.CartItem
	if err := r.db.WithContext(ctx).
		Where("cart_id = ? AND product_id = ?", cartID, productID).
		First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *cartRepo) CreateItem(ctx context.Context, item *entity.CartItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *cartRepo) UpdateItem(ctx context.Context, item *entity.CartItem) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *cartRepo) DeleteItem(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.CartItem{}, "id = ?", id).Error
}
