package service

import (
	"context"
	"time"

	"marketplace/entity"
	"marketplace/repository"
)

type CategoryService interface {
	CreateCategory(ctx context.Context, c *entity.ProductCategory) error
	GetCategoryByID(ctx context.Context, id uint) (*entity.ProductCategory, error)
	GetCategoryByIDIncludeDeleted(ctx context.Context, id uint) (*entity.ProductCategory, error)
	GetCategories(ctx context.Context, limit, offset int) ([]entity.ProductCategory, error)
	UpdateCategory(ctx context.Context, c *entity.ProductCategory) error
	DeleteCategory(ctx context.Context, id uint) error
	RecoverCategory(ctx context.Context, id uint) error
}

type categoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) CreateCategory(ctx context.Context, c *entity.ProductCategory) error {
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	return s.repo.Create(ctx, c)
}

func (s *categoryService) GetCategoryByID(ctx context.Context, id uint) (*entity.ProductCategory, error) {
	return s.repo.GetByCategoryID(ctx, id)
}

func (s *categoryService) GetCategoryByIDIncludeDeleted(ctx context.Context, id uint) (*entity.ProductCategory, error) {
	return s.repo.GetByIDIncludeDeleted(ctx, id)
}

func (s *categoryService) GetCategories(ctx context.Context, limit, offset int) ([]entity.ProductCategory, error) {
	return s.repo.GetCategories(ctx, limit, offset)
}

func (s *categoryService) UpdateCategory(ctx context.Context, c *entity.ProductCategory) error {
	existing, err := s.repo.GetByIDIncludeDeleted(ctx, c.ID)
	if err != nil {
		return err
	}
	
    // Update fields
	existing.Name = c.Name
	existing.Description = c.Description
	existing.UpdatedAt = time.Now()
    
	// Update categories sesuai yg di input aja
	if c.Products != nil {
		existing.Products = c.Products
	}

	return s.repo.Update(ctx, existing)
}

func (s *categoryService) DeleteCategory(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func (s *categoryService) RecoverCategory(ctx context.Context, id uint) error {
	return s.repo.Recover(ctx, id)
}