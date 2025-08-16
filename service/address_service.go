package service

import (
	"context"
	"errors"
	"time"

	"marketplace/entity"
	"marketplace/repository"
	"github.com/google/uuid"
)

type AddressService interface {
	CreateAddress(ctx context.Context, a *entity.Address) error
	GetAddressByUser(ctx context.Context, userID uuid.UUID) (*entity.Address, error)
	GetAddressByIDIncludeDeleted(ctx context.Context, id string) (*entity.Address, error)
	GetAddresses(ctx context.Context, search string, limit, offset int) ([]entity.Address, error) // admin
	UpdateAddress(ctx context.Context, a *entity.Address) error
	DeleteAddress(ctx context.Context, id string, userID uuid.UUID) error
	RecoverAddress(ctx context.Context, id string, userID uuid.UUID) error
}

type addressService struct {
	repo repository.AddressRepository
}

func NewAddressService(repo repository.AddressRepository) AddressService {
	return &addressService{repo: repo}
}

func (s *addressService) CreateAddress(ctx context.Context, a *entity.Address) error {
	existing, _ := s.repo.GetByUserID(ctx, a.UserID)
	if existing != nil {
		return errors.New("user already has an address")
	}
	a.ID = uuid.New()
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()
	return s.repo.Create(ctx, a)
}

// User GET address sendiri
func (s *addressService) GetAddressByUser(ctx context.Context, userID uuid.UUID) (*entity.Address, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *addressService) GetAddressByIDIncludeDeleted(ctx context.Context, id string) (*entity.Address, error) {
	return s.repo.GetByIDIncludeDeleted(ctx, id)
}

func (s *addressService) GetAddresses(ctx context.Context, search string, limit, offset int) ([]entity.Address, error) {
	return s.repo.GetAddresses(ctx, search, limit, offset)
}

func (s *addressService) UpdateAddress(ctx context.Context, a *entity.Address) error {
	existing, err := s.repo.GetByIDIncludeDeleted(ctx, a.ID.String())
	if err != nil || existing == nil {
		return errors.New("address not found")
	}

	// cek user id dulu
	if existing.UserID != a.UserID {
		return errors.New("forbidden")
	}

	// Update fiekds sesuai yg di input aja
	if a.ReceiverName != "" {
		existing.ReceiverName = a.ReceiverName
	}
	if a.PhoneNumber != "" {
		existing.PhoneNumber = a.PhoneNumber
	}
	if a.AddressLine != "" {
		existing.AddressLine = a.AddressLine
	}
	if a.City != "" {
		existing.City = a.City
	}
	if a.Province != "" {
		existing.Province = a.Province
	}
	if a.PostalCode != "" {
		existing.PostalCode = a.PostalCode
	}

	existing.UpdatedAt = time.Now()
	return s.repo.Update(ctx, existing)
}

func (s *addressService) DeleteAddress(ctx context.Context, id string, userID uuid.UUID) error {
	existing, _ := s.repo.GetByIDIncludeDeleted(ctx, id)
	if existing == nil {
		return errors.New("address not found")
	}
	if existing.UserID != userID {
		return errors.New("forbidden")
	}
	return s.repo.Delete(ctx, id)
}

func (s *addressService) RecoverAddress(ctx context.Context, id string, userID uuid.UUID) error {
	existing, _ := s.repo.GetByIDIncludeDeleted(ctx, id)
	if existing == nil {
		return errors.New("address not found")
	}
	if existing.UserID != userID {
		return errors.New("forbidden")
	}
	return s.repo.Recover(ctx, id)
}

