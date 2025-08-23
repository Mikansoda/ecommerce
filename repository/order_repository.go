package repository

import (
	"context"
	"time"

	"ecommerce/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OrderRepository interface {
	Create(ctx context.Context, order *entity.Order, tx *gorm.DB) error
	GetOrders(ctx context.Context, limit, offset int) ([]entity.Order, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Order, error)
	GetOrdersByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Order, error)
	GetByStatus(ctx context.Context, status string, limit, offset int) ([]entity.Order, error)
	GetPendingOrdersOlderThan(ctx context.Context, duration time.Duration) ([]entity.Order, error)
	Update(ctx context.Context, order *entity.Order, tx *gorm.DB) error
	GetByIDForUpdate(ctx context.Context, id uuid.UUID, tx *gorm.DB) (*entity.Order, error)
	BeginTx() *gorm.DB
}

type orderRepo struct {
	db *gorm.DB
}

func NewOrderRepo(db *gorm.DB) OrderRepository {
	return &orderRepo{db: db}
}

func (r *orderRepo) Create(ctx context.Context, order *entity.Order, tx *gorm.DB) error {
	return tx.WithContext(ctx).Create(order).Error
}

func (r *orderRepo) GetOrders(ctx context.Context, limit, offset int) ([]entity.Order, error) {
	var orders []entity.Order
	err := r.db.WithContext(ctx).
		Preload("OrderItems").
		Limit(limit).Offset(offset).
		Find(&orders).Error
	return orders, err
}

func (r *orderRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.Order, error) {
	var order entity.Order
	err := r.db.WithContext(ctx).
		Preload("OrderItems").
		Preload("User").
		First(&order, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepo) GetOrdersByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Order, error) {
	var orders []entity.Order
	err := r.db.WithContext(ctx).
		Preload("OrderItems").
		Where("user_id = ?", userID).
		Find(&orders).Error
	return orders, err
}

func (r *orderRepo) GetByStatus(ctx context.Context, status string, limit, offset int) ([]entity.Order, error) {
	var orders []entity.Order
	if err := r.db.WithContext(ctx).
		Where("status = ?", status).
		Limit(limit).Offset(offset).
		Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *orderRepo) GetPendingOrdersOlderThan(ctx context.Context, duration time.Duration) ([]entity.Order, error) {
	var orders []entity.Order
	cutoff := time.Now().Add(-duration)
	err := r.db.WithContext(ctx).
		Preload("OrderItems").
		Where("status = ? AND created_at < ?", "pending", cutoff).
		Find(&orders).Error
	return orders, err
}

func (r *orderRepo) Update(ctx context.Context, order *entity.Order, tx *gorm.DB) error {
	return tx.WithContext(ctx).Save(order).Error
}

func (r *orderRepo) GetByIDForUpdate(ctx context.Context, id uuid.UUID, tx *gorm.DB) (*entity.Order, error) {
	var order entity.Order
	err := tx.WithContext(ctx).
		Preload("OrderItems").
		Preload("User").
		Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&order, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepo) BeginTx() *gorm.DB {
    return r.db.Begin()
}