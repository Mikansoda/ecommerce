package repository

import (
	"context"
	"time"

	"marketplace/entity"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, u *entity.UsersAndAdmins) error
	FindByEmail(ctx context.Context, email string) (*entity.UsersAndAdmins, error)
	FindByUsername(ctx context.Context, username string) (*entity.UsersAndAdmins, error)
	Update(ctx context.Context, u *entity.UsersAndAdmins) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, u *entity.UsersAndAdmins) error {
	return r.db.WithContext(ctx).Create(u).Error
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entity.UsersAndAdmins, error) {
	var u entity.UsersAndAdmins
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*entity.UsersAndAdmins, error) {
	var u entity.UsersAndAdmins
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) Update(ctx context.Context, u *entity.UsersAndAdmins) error {
	u.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(u).Error
}
