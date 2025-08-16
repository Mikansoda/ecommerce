package service

import (
	"context"
	"time"
	"errors"

	"marketplace/entity"
	"marketplace/repository"
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
	return s.repo.Create(ctx, p)
}

func (s *productService) GetProductByID(ctx context.Context, id uint) (*entity.Product, error) {
	return s.repo.GetByProductID(ctx, id)
}

func (s *productService) GetProductByIDIncludeDeleted(ctx context.Context, id uint) (*entity.Product, error) {
	return s.repo.GetByIDIncludeDeleted(ctx, id)
}

func (s *productService) GetProducts(ctx context.Context, search, category string, limit, offset int) ([]entity.Product, error) {
	return s.repo.GetProducts(ctx, search, category, limit, offset)
}

func (s *productService) UpdateProduct(ctx context.Context, p *entity.Product) error {
	existing, err := s.repo.GetByIDIncludeDeleted(ctx, p.ID)
	if err != nil {
		return err
	}

	// Update fiekds sesuai yg di input aja
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

	existing.UpdatedAt = time.Now()

	if p.Categories != nil {
		existing.Categories = p.Categories
	}

	return s.repo.Update(ctx, existing)
}

func (s *productService) DeleteProduct(ctx context.Context, id uint) error {
	existing, _ := s.repo.GetByProductID(ctx, id)
	if existing == nil {
		return errors.New("address not found")
	}
	return s.repo.Delete(ctx, id)
}

func (s *productService) RecoverProduct(ctx context.Context, id uint) error {
	existing, _ := s.repo.GetByIDIncludeDeleted(ctx, id)
	if existing == nil {
		return errors.New("address not found")
	}
	return s.repo.Recover(ctx, id)
}
