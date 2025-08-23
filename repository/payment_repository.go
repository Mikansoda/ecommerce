package repository

import (
	"context"
	"ecommerce/entity"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentRepository interface {
	Create(ctx context.Context, payment *entity.Payment) error
	GetByOrderID(ctx context.Context, orderID uuid.UUID) (*entity.Payment, error)
	GetPendingOlderThan(ctx context.Context, hours int) ([]entity.Payment, error)
	GetByUserID(ctx context.Context, userID string, out *[]entity.Payment) error
	GetPayments(ctx context.Context, out *[]entity.Payment) error
	GetByInvoiceID(ctx context.Context, invoiceID string) (*entity.Payment, error)
	Update(ctx context.Context, payment *entity.Payment) error
}

type paymentRepo struct {
	db *gorm.DB
}

func NewPaymentRepo(db *gorm.DB) PaymentRepository {
	return &paymentRepo{db: db}
}

func (r *paymentRepo) Create(ctx context.Context, payment *entity.Payment) error {
	return r.db.WithContext(ctx).Create(payment).Error
}

func (r *paymentRepo) GetByOrderID(ctx context.Context, orderID uuid.UUID) (*entity.Payment, error) {
	var payment entity.Payment
	err := r.db.WithContext(ctx).Where("order_id = ?", orderID).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepo) GetPendingOlderThan(ctx context.Context, hours int) ([]entity.Payment, error) {
	var payments []entity.Payment
	cutoff := time.Now().Add(-time.Duration(hours) * time.Hour)
	err := r.db.WithContext(ctx).Where("status = ? AND created_at < ?", "pending", cutoff).Find(&payments).Error
	return payments, err
}

func (r *paymentRepo) GetByUserID(ctx context.Context, userID string, out *[]entity.Payment) error {
    return r.db.WithContext(ctx).
        Joins("JOIN orders ON orders.id = payments.order_id").
        Where("orders.user_id = ?", userID).
		Order("payments.created_at DESC").
        Find(out).Error
}

func (r *paymentRepo) GetPayments(ctx context.Context, out *[]entity.Payment) error {
	return r.db.WithContext(ctx).Find(out).Error
}

func (r *paymentRepo) GetByInvoiceID(ctx context.Context, invoiceID string) (*entity.Payment, error) {
	fmt.Println("Searching payment with invoice_id:", invoiceID)
	var payment entity.Payment
	err := r.db.WithContext(ctx).Where("invoice_id = ?", invoiceID).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepo) Update(ctx context.Context, payment *entity.Payment) error {
	return r.db.WithContext(ctx).Save(payment).Error
}
