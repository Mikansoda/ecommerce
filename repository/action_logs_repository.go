package repository

import (
	"context"
	"ecommerce/entity"

	"gorm.io/gorm"
)

type ActionLogRepository interface {
	Create(ctx context.Context, log *entity.ActionLog) error
	Getlogs(ctx context.Context) ([]entity.ActionLog, error)
	GetByLogID(ctx context.Context, id string) (*entity.ActionLog, error)
	ReportSelling(ctx context.Context, reportType string, limit int) ([]map[string]interface{}, error)
	ReportStock(ctx context.Context, reportType string, limit int) ([]entity.Product, error)
}

type actionLogRepo struct {
	db *gorm.DB
}

func NewActionLogRepository(db *gorm.DB) ActionLogRepository {
	return &actionLogRepo{db}
}

func (r *actionLogRepo) Create(ctx context.Context, log *entity.ActionLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *actionLogRepo) Getlogs(ctx context.Context) ([]entity.ActionLog, error) {
	var logs []entity.ActionLog
	err := r.db.WithContext(ctx).Preload("User").Find(&logs).Error
	return logs, err
}

func (r *actionLogRepo) GetByLogID(ctx context.Context, id string) (*entity.ActionLog, error) {
	var log entity.ActionLog
	err := r.db.WithContext(ctx).Preload("User").Where("id = ?", id).First(&log).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

func (r *actionLogRepo) ReportSelling(ctx context.Context, reportType string, limit int) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	order := "DESC"
	if reportType == "least" {
		order = "ASC"
	}
	err := r.db.WithContext(ctx).
		Table("order_items AS oi").
		Select("oi.product_id, p.name AS product_name, SUM(oi.quantity) AS total_sold").
		Joins("JOIN orders o ON oi.order_id = o.id").
		Joins("JOIN products p ON oi.product_id = p.id").
		Where("o.status IN ?", []string{"paid", "shipped", "completed"}).
		Group("oi.product_id, p.name").
		Order("total_sold " + order).
		Limit(limit).
		Scan(&result).Error

	return result, err
}

func (r *actionLogRepo) ReportStock(ctx context.Context, reportType string, limit int) ([]entity.Product, error) {
	var products []entity.Product
	order := "ASC"
	if reportType == "high" {
		order = "DESC"
	}

	err := r.db.WithContext(ctx).
		Model(&entity.Product{}).
		Preload("Categories").
		Order("stock " + order).
		Limit(limit).
		Find(&products).Error
	return products, err
}
