package repository

import (
	"context"

	"ecommerce/entity"

	"gorm.io/gorm"
)

type ProductImageRepository interface {
	Create(ctx context.Context, img *entity.ProductImage) error
	GetByProductID(ctx context.Context, productID uint) ([]entity.ProductImage, error)
	GetByImageID(ctx context.Context, id uint) (*entity.ProductImage, error)
	GetByIDIncludeDeleted(ctx context.Context, id uint) (*entity.ProductImage, error)
	CountByProductID(ctx context.Context, productID uint) (int64, error)
	UnsetPrimary(ctx context.Context, productID uint) error
	Delete(ctx context.Context, id uint) error
	Recover(ctx context.Context, id uint) error
}

type productImageRepo struct {
	db *gorm.DB
}

func NewProductImageRepository(db *gorm.DB) ProductImageRepository {
	return &productImageRepo{db: db}
}

func (r *productImageRepo) Create(ctx context.Context, img *entity.ProductImage) error {
	return r.db.WithContext(ctx).Create(img).Error
}

func (r *productImageRepo) GetByProductID(ctx context.Context, productID uint) ([]entity.ProductImage, error) {
	var imgs []entity.ProductImage
	err := r.db.WithContext(ctx).Where("product_id = ?", productID).Find(&imgs).Error
	return imgs, err
}

func (r *productImageRepo) GetByImageID(ctx context.Context, id uint) (*entity.ProductImage, error) {
	var img entity.ProductImage
	err := r.db.WithContext(ctx).First(&img, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &img, nil
}

func (r *productImageRepo) GetByIDIncludeDeleted(ctx context.Context, id uint) (*entity.ProductImage, error) {
	var img entity.ProductImage
	err := r.db.WithContext(ctx).Unscoped().First(&img, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &img, nil
}

func (r *productImageRepo) CountByProductID(ctx context.Context, productID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.ProductImage{}).
		Where("product_id = ?", productID).Count(&count).Error
	return count, err
}

func (r *productImageRepo) UnsetPrimary(ctx context.Context, productID uint) error {
	return r.db.WithContext(ctx).
		Model(&entity.ProductImage{}).
		Where("product_id = ?", productID).
		Update("is_primary", false).Error
}

func (r *productImageRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.ProductImage{}, id).Error
}

func (r *productImageRepo) Recover(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).
		Model(&entity.ProductImage{}).
		Unscoped().
		Where("id = ?", id).
		Update("deleted_at", nil).Error
}
