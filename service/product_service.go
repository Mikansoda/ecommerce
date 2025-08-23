package service

import (
	"context"
	"errors"
	"time"

	"ecommerce/entity"
	"ecommerce/repository"
)

type ProductService interface {
	CreateProduct(ctx context.Context, p *entity.Product) error
	GetProductByID(ctx context.Context, id uint) (*entity.Product, error)
	GetProductByIDIncludeDeleted(ctx context.Context, id uint) (*entity.Product, error)
	GetProducts(ctx context.Context, search, category string, limit, offset int) ([]entity.Product, error)
	UpdateProduct(ctx context.Context, p *entity.Product) error
	DeleteProduct(ctx context.Context, id uint) error
	RecoverProduct(ctx context.Context, id uint) error
}

type productService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{repo: repo}
}

func (s *productService) CreateProduct(ctx context.Context, p *entity.Product) error {
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()

	// optional transaction usage
	if err := s.repo.Create(ctx, p); err != nil {
		return errors.New("failed to create product: " + err.Error())
	}
	return nil
}

func (s *productService) GetProductByID(ctx context.Context, id uint) (*entity.Product, error) {
	product, err := s.repo.GetByProductID(ctx, id)
	if err != nil {
		return nil, errors.New("failed to fetch product: " + err.Error())
	}
	if product == nil {
		return nil, errors.New("product not found")
	}
	return product, nil
}

func (s *productService) GetProductByIDIncludeDeleted(ctx context.Context, id uint) (*entity.Product, error) {
	product, err := s.repo.GetByIDIncludeDeleted(ctx, id)
	if err != nil {
		return nil, errors.New("failed to fetch product: " + err.Error())
	}
	if product == nil {
		return nil, errors.New("product not found")
	}
	return product, nil
}

func (s *productService) GetProducts(ctx context.Context, search, category string, limit, offset int) ([]entity.Product, error) {
	products, err := s.repo.GetProducts(ctx, search, category, limit, offset)
	if err != nil {
		return nil, errors.New("failed to fetch products: " + err.Error())
	}
	return products, nil
}

func (s *productService) UpdateProduct(ctx context.Context, p *entity.Product) error {
	existing, err := s.repo.GetByIDIncludeDeleted(ctx, p.ID)
	if err != nil {
		return errors.New("failed to fetch product: " + err.Error())
	}
	if existing == nil {
		return errors.New("product not found")
	}

	// Update fields
	if p.Name != "" {
		existing.Name = p.Name
	}
	if p.Description != "" {
		existing.Description = p.Description
	}
	if p.Price != 0 {
		existing.Price = p.Price
	}
	if p.Stock != 0 {
		existing.Stock = p.Stock
	}
	if p.Categories != nil {
		existing.Categories = p.Categories
	}
	if p.ExpiryYear != nil {
		existing.ExpiryYear = p.ExpiryYear
	}
	existing.UpdatedAt = time.Now()

	// optional transaction usage
	if err := s.repo.Update(ctx, existing); err != nil {
		return errors.New("failed to update product: " + err.Error())
	}
	return nil
}

func (s *productService) DeleteProduct(ctx context.Context, id uint) error {
	existing, err := s.repo.GetByProductID(ctx, id)
	if err != nil {
		return errors.New("failed to fetch product: " + err.Error())
	}
	if existing == nil {
		return errors.New("product not found")
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return errors.New("failed to delete product: " + err.Error())
	}
	return nil
}

func (s *productService) RecoverProduct(ctx context.Context, id uint) error {
	existing, err := s.repo.GetByIDIncludeDeleted(ctx, id)
	if err != nil {
		return errors.New("failed to fetch product: " + err.Error())
	}
	if existing == nil {
		return errors.New("product not found")
	}
	if err := s.repo.Recover(ctx, id); err != nil {
		return errors.New("failed to recover product: " + err.Error())
	}
	return nil
}
