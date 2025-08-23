package service

import (
	"context"
	"errors"
	"time"

	"ecommerce/entity"
	"ecommerce/repository"

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
	existing, err := s.repo.GetByUserID(ctx, a.UserID)
	if err != nil {
		return errors.New("failed to check existing address: " + err.Error())
	}
	if existing != nil {
		return errors.New("user already has an address")
	}
	a.ID = uuid.New()
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()
	if err := s.repo.Create(ctx, a); err != nil {
		return errors.New("failed to create address: " + err.Error())
	}
	return nil
}

// User GET own address
func (s *addressService) GetAddressByUser(ctx context.Context, userID uuid.UUID) (*entity.Address, error) {
	address, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, errors.New("failed to fetch address: " + err.Error())
	}
	if address == nil {
		return nil, errors.New("address not found")
	}
	return address, nil
}

func (s *addressService) GetAddressByIDIncludeDeleted(ctx context.Context, id string) (*entity.Address, error) {
	address, err := s.repo.GetByIDIncludeDeleted(ctx, id)
	if err != nil {
		return nil, errors.New("failed to fetch address: " + err.Error())
	}
	if address == nil {
		return nil, errors.New("address not found")
	}
	return address, nil
}

func (s *addressService) GetAddresses(ctx context.Context, search string, limit, offset int) ([]entity.Address, error) {
	addresses, err := s.repo.GetAddresses(ctx, search, limit, offset)
	if err != nil {
		return nil, errors.New("failed to fetch addresses: " + err.Error())
	}
	return addresses, nil
}

func (s *addressService) UpdateAddress(ctx context.Context, a *entity.Address) error {
	existing, err := s.repo.GetByIDIncludeDeleted(ctx, a.ID.String())
	if err != nil {
		return errors.New("failed to fetch address: " + err.Error())
	}
	if existing == nil {
		return errors.New("address not found")
	}

	if existing.UserID != a.UserID {
		return errors.New("forbidden")
	}

	// Update fields
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
	if err := s.repo.Update(ctx, existing); err != nil {
		return errors.New("failed to update address: " + err.Error())
	}
	return nil
}

func (s *addressService) DeleteAddress(ctx context.Context, id string, userID uuid.UUID) error {
	existing, err := s.repo.GetByIDIncludeDeleted(ctx, id)
	if err != nil {
		return errors.New("failed to fetch address: " + err.Error())
	}
	if existing == nil {
		return errors.New("address not found")
	}
	if existing.UserID != userID {
		return errors.New("forbidden")
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return errors.New("failed to delete address: " + err.Error())
	}
	return nil
}

func (s *addressService) RecoverAddress(ctx context.Context, id string, _ uuid.UUID) error {
	existing, err := s.repo.GetByIDIncludeDeleted(ctx, id)
	if err != nil {
		return errors.New("failed to fetch address: " + err.Error())
	}
	if existing == nil {
		return errors.New("address not found")
	}
	if err := s.repo.Recover(ctx, id); err != nil {
		return errors.New("failed to recover address: " + err.Error())
	}
	return nil
}
