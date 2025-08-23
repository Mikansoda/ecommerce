package repository

import (
	"context"

	"ecommerce/entity"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductRepository interface {
	Create(ctx context.Context, p *entity.Product, tx ...*gorm.DB) error
	GetProducts(ctx context.Context, search string, categoryName string, limit, offset int) ([]entity.Product, error)
	GetByProductID(ctx context.Context, id uint) (*entity.Product, error)
	GetByIDIncludeDeleted(ctx context.Context, id uint) (*entity.Product, error)
	Update(ctx context.Context, p *entity.Product, tx ...*gorm.DB) error
	Delete(ctx context.Context, id uint) error
	Recover(ctx context.Context, id uint) error
	BeginTx(ctx context.Context) (*gorm.DB, error)
	CommitTx(tx *gorm.DB) error
	RollbackTx(tx *gorm.DB) error
	GetByProductIDForUpdate(ctx context.Context, id uint, tx *gorm.DB) (*entity.Product, error)
}

type productRepo struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepo{db: db}
}

func (r *productRepo) Create(ctx context.Context, p *entity.Product, tx ...*gorm.DB) error {
	db := r.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0]
	}
	return db.WithContext(ctx).Create(p).Error
}

func (r *productRepo) GetProducts(ctx context.Context, search string, categoryName string, limit, offset int) ([]entity.Product, error) {
	var products []entity.Product
	query := r.db.WithContext(ctx).Model(&entity.Product{}).Preload("Images").Preload("Categories")

	if search != "" {
		query = query.Where("products.name LIKE ?", "%"+search+"%")
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
		First(&p, id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *productRepo) Update(ctx context.Context, p *entity.Product, tx ...*gorm.DB) error {
	db := r.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0]
	}
	return db.WithContext(ctx).Save(p).Error
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

func (r *productRepo) BeginTx(ctx context.Context) (*gorm.DB, error) {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}

func (r *productRepo) CommitTx(tx *gorm.DB) error {
	if tx == nil {
		return errors.New("no active transaction")
	}
	return tx.Commit().Error
}

func (r *productRepo) RollbackTx(tx *gorm.DB) error {
	if tx == nil {
		return errors.New("no active transaction")
	}
	return tx.Rollback().Error
}

func (r *productRepo) GetByProductIDForUpdate(ctx context.Context, id uint, tx *gorm.DB) (*entity.Product, error) {
	var p entity.Product
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&p, id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}
