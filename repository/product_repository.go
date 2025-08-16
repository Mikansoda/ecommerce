package repository

import (
	"context"

	"marketplace/entity"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(ctx context.Context, p *entity.Product) error
	GetProducts(ctx context.Context, search string, categoryName string, limit, offset int) ([]entity.Product, error)
	GetByProductID(ctx context.Context, id uint) (*entity.Product, error)
	GetByIDIncludeDeleted(ctx context.Context, id uint) (*entity.Product, error)
	Update(ctx context.Context, p *entity.Product) error
	Delete(ctx context.Context, id uint) error
	Recover(ctx context.Context, id uint) error
}

type productRepo struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepo{db: db}
}

func (r *productRepo) Create(ctx context.Context, p *entity.Product) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *productRepo) GetProducts(ctx context.Context, search string, categoryName string, limit, offset int) ([]entity.Product, error) {
	var products []entity.Product
	query := r.db.WithContext(ctx).Model(&entity.Product{}).Preload("Images").Preload("Categories")

	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}
	if categoryName != "" {
		query = query.Joins("JOIN product_categories_map pcm ON pcm.product_id = products.id").
			Joins("JOIN product_categories pc ON pc.id = pcm.product_category_id").
			Where("pc.name = ?", categoryName)
	}

	if err := query.Limit(limit).Offset(offset).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *productRepo) GetByProductID(ctx context.Context, id uint) (*entity.Product, error) {
	var product entity.Product
	if err := r.db.WithContext(ctx).
		Preload("Images").
		Preload("Categories").
		Preload("Ratings").
		First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepo) GetByIDIncludeDeleted(ctx context.Context, id uint) (*entity.Product, error) {
	var p entity.Product
	if err := r.db.WithContext(ctx).
		Unscoped().
		Preload("Images").
		Preload("Categories").
		Preload("Ratings").
		First(&p, id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *productRepo) Update(ctx context.Context, p *entity.Product) error {
	return r.db.WithContext(ctx).Save(p).Error
}

func (r *productRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.Product{}, id).Error
}

func (r *productRepo) Recover(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).
		Model(&entity.Product{}).
		Unscoped().
		Where("id = ?", id).
		Update("deleted_at", nil).Error
}