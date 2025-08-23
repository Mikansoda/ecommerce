package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"ecommerce/config"
	"ecommerce/entity"
	"ecommerce/repository"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type ProductImageService interface {
	Upload(ctx context.Context, productID uint, filePath string, isPrimary bool) (*entity.ProductImage, error)
	GetByProductID(ctx context.Context, productID uint) ([]entity.ProductImage, error)
	Delete(ctx context.Context, id uint) error
	Recover(ctx context.Context, id uint) error
}

type productImageService struct {
	repo repository.ProductImageRepository
}

func NewProductImageService(repo repository.ProductImageRepository) ProductImageService {
	return &productImageService{repo: repo}
}

func (s *productImageService) Upload(ctx context.Context, productID uint, filePath string, isPrimary bool) (*entity.ProductImage, error) {
	// max 3 images check
	count, err := s.repo.CountByProductID(ctx, productID)
	if err != nil {
		return nil, errors.New("failed to count product images: " + err.Error())
	}
	if count >= 3 {
		return nil, errors.New("maximum 3 images per product")
	}

	f, err := os.Open(filePath)
	if err != nil {
		return nil, errors.New("failed to open file: " + err.Error())
	}
	defer f.Close()

	cld := config.InitCloud()
	uploadRes, err := cld.Upload.Upload(ctx, f, uploader.UploadParams{
		Folder: "ecommerce/products",
	})
	if err != nil {
		return nil, errors.New("failed to upload image: " + err.Error())
	}

	fmt.Printf("DEBUG uploadRes: %+v\n", uploadRes)

	if uploadRes == nil || uploadRes.SecureURL == "" {
		return nil, errors.New("failed to upload image: empty secure URL from cloud")
	}

	if isPrimary {
		if err := s.repo.UnsetPrimary(ctx, productID); err != nil {
			return nil, errors.New("failed to unset previous primary image: " + err.Error())
		}
	}

	img := &entity.ProductImage{
		ProductID: productID,
		ImageURL:  uploadRes.SecureURL,
		IsPrimary: isPrimary,
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, img); err != nil {
		return nil, errors.New("failed to save product image: " + err.Error())
	}
	return img, nil
}

func (s *productImageService) GetByProductID(ctx context.Context, productID uint) ([]entity.ProductImage, error) {
	images, err := s.repo.GetByProductID(ctx, productID)
	if err != nil {
		return nil, errors.New("failed to fetch product images: " + err.Error())
	}
	return images, nil
}

func (s *productImageService) Delete(ctx context.Context, id uint) error {
	existing, err := s.repo.GetByImageID(ctx, id)
	if err != nil {
		return errors.New("failed to fetch product image: " + err.Error())
	}
	if existing == nil {
		return errors.New("product image not found")
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return errors.New("failed to delete product image: " + err.Error())
	}
	return nil
}

func (s *productImageService) Recover(ctx context.Context, id uint) error {
	existing, err := s.repo.GetByIDIncludeDeleted(ctx, id)
	if err != nil {
		return errors.New("failed to fetch product image: " + err.Error())
	}
	if existing == nil {
		return errors.New("product image not found")
	}
	if err := s.repo.Recover(ctx, id); err != nil {
		return errors.New("failed to recover product image: " + err.Error())
	}
	return nil
}
