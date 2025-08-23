package repository

import (
	"context"

	"ecommerce/entity"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(ctx context.Context, c *entity.ProductCategory) error
	GetCategories(ctx context.Context, limit, offset int) ([]entity.ProductCategory, error)
	GetByIDIncludeDeleted(ctx context.Context, id uint) (*entity.ProductCategory, error)
	Update(ctx context.Context, c *entity.ProductCategory) error
	Delete(ctx context.Context, id uint) error
	Recover(ctx context.Context, id uint) error
}

type categoryRepo struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepo{db: db}
}

func (r *categoryRepo) Create(ctx context.Context, c *entity.ProductCategory) error {
	return r.db.WithContext(ctx).Create(c).Error
}

func (r *categoryRepo) GetCategories(ctx context.Context, limit, offset int) ([]entity.ProductCategory, error) {
	var categories []entity.ProductCategory
	if err := r.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *categoryRepo) GetByIDIncludeDeleted(ctx context.Context, id uint) (*entity.ProductCategory, error) {
	var c entity.ProductCategory
	if err := r.db.WithContext(ctx).
		Unscoped().
		First(&c, id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *categoryRepo) Update(ctx context.Context, c *entity.ProductCategory) error {
	return r.db.WithContext(ctx).Save(c).Error
}

func (r *categoryRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.ProductCategory{}, id).Error
}

func (r *categoryRepo) Recover(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).
		Model(&entity.ProductCategory{}).
		Unscoped().
		Where("id = ?", id).
		Update("deleted_at", nil).Error
}
