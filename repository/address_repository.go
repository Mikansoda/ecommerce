package repository

import (
	"context"

	"marketplace/entity"
	"gorm.io/gorm"
	"github.com/google/uuid"
)

type AddressRepository interface {
	Create(ctx context.Context, a *entity.Address) error
	GetAddresses(ctx context.Context, search string, limit, offset int) ([]entity.Address, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.Address, error)
	GetByIDIncludeDeleted(ctx context.Context, id string) (*entity.Address, error)
	Update(ctx context.Context, a *entity.Address) error
	Delete(ctx context.Context, id string) error
	Recover(ctx context.Context, id string) error
}

type addressRepo struct {
	db *gorm.DB
}

func NewAddressRepository(db *gorm.DB) AddressRepository {
	return &addressRepo{db: db}
}

func (r *addressRepo) Create(ctx context.Context, a *entity.Address) error {
	return r.db.WithContext(ctx).Create(a).Error
}

func (r *addressRepo) GetAddresses(ctx context.Context, search string, limit, offset int) ([]entity.Address, error) {
	var addresses []entity.Address
	query := r.db.WithContext(ctx).Model(&entity.Address{})

	if search != "" {
		query = query.Where("receiver_name LIKE ? OR city LIKE ? OR province LIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Limit(limit).Offset(offset).Find(&addresses).Error; err != nil {
		return nil, err
	}
	return addresses, nil
}

func (r *addressRepo) GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.Address, error) {
	var address entity.Address
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&address).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &address, nil
}

func (r *addressRepo) GetByIDIncludeDeleted(ctx context.Context, id string) (*entity.Address, error) {
	var a entity.Address
	if err := r.db.WithContext(ctx).Unscoped().Where("id = ?", id).First(&a).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *addressRepo) Update(ctx context.Context, a *entity.Address) error {
	return r.db.WithContext(ctx).Save(a).Error
}

func (r *addressRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&entity.Address{}, "id = ?", id).Error
}

func (r *addressRepo) Recover(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Model(&entity.Address{}).
		Unscoped().
		Where("id = ?", id).
		Update("deleted_at", nil).Error
}
