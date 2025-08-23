package controller

import (
	"net/http"
	"strconv"
	"time"

	"ecommerce/entity"
	"ecommerce/service"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	service service.ProductService
}

func NewProductController(s service.ProductService) *ProductController {
	return &ProductController{service: s}
}

// Struct request for endpoint create product and update
type createProductReq struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Price       float64 `json:"price" binding:"required"`
	CategoryIDs []uint  `json:"category_ids" binding:"required"`
	Stock       uint    `json:"stock" binding:"required"`
	ExpiryYear  *int    `json:"expiry_year,omitempty"`
}

type updateProductReq struct {
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	Price       *float64 `json:"price,omitempty"`
	CategoryIDs []uint   `json:"category_ids,omitempty"`
	Stock       *uint    `json:"stock,omitempty"`
	ExpiryYear  *int     `json:"expiry_year,omitempty"`
}

// GET products
func (ctl *ProductController) GetProducts(c *gin.Context) {
	search := c.Query("search")
	category := c.Query("category")
	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")

	limit, _ := strconv.Atoi(limitStr)
	if limit == 0 {
		limit = 10
	}
	offset, _ := strconv.Atoi(offsetStr)

	products, err := ctl.service.GetProducts(c.Request.Context(), search, category, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch products, try again later",
			"detail":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, products)
}

// GET product by ID
func (ctl *ProductController) GetProductByID(c *gin.Context) {
	idStr := c.Param("productId")
	idUint64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid product ID",
			"detail":  err.Error(),
		})
		return
	}
	id := uint(idUint64)

	product, err := ctl.service.GetProductByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Product not found",
			"detail":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, product)
}

// Create product (admin only)
func (ctl *ProductController) CreateProduct(c *gin.Context) {
	var req createProductReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid input data",
			"detail":  err.Error(),
		})
		return
	}

	var categories []entity.ProductCategory
	for _, cid := range req.CategoryIDs {
		categories = append(categories, entity.ProductCategory{ID: cid})
	}
	product := &entity.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Categories:  categories,
		Stock:       req.Stock,
		ExpiryYear:  req.ExpiryYear,
	}
	if err := ctl.service.CreateProduct(c.Request.Context(), product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create product, try again later",
			"detail":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product successfully created",
		"data":    product,
	})
}

// Update product (admin only)
func (ctl *ProductController) UpdateProduct(c *gin.Context) {
	idStr := c.Param("productId")
	idUint64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid product ID",
			"detail":  err.Error(),
		})
		return
	}
	id := uint(idUint64)

	var req updateProductReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid input data",
			"detail":  err.Error(),
		})
		return
	}

	existing, err := ctl.service.GetProductByIDIncludeDeleted(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Product not found",
			"detail":  err.Error(),
		})
		return
	}

	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Description != nil {
		existing.Description = *req.Description
	}
	if req.Price != nil {
		existing.Price = *req.Price
	}
	if req.CategoryIDs != nil {
		var categories []entity.ProductCategory
		for _, cid := range req.CategoryIDs {
			categories = append(categories, entity.ProductCategory{ID: cid})
		}
		existing.Categories = categories
	}
	if req.Stock != nil {
		existing.Stock = *req.Stock
	}
	if req.ExpiryYear != nil {
	existing.ExpiryYear = req.ExpiryYear
    }
	existing.UpdatedAt = time.Now()

	if err := ctl.service.UpdateProduct(c.Request.Context(), existing); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update product, try again later",
			"detail":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product successfully updated",
		"data":    existing,
	})
}

// Delete product (admin only)
func (ctl *ProductController) DeleteProduct(c *gin.Context) {
	idStr := c.Param("productId")
	idUint64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid product ID",
			"detail":  err.Error(),
		})
		return
	}
	id := uint(idUint64)

	if err := ctl.service.DeleteProduct(c.Request.Context(), id); err != nil {
		status := http.StatusInternalServerError
		msg := "Failed to delete product, try again later"

		if err.Error() == "product not found" {
			status = http.StatusNotFound
			msg = "Product not found"
		}

		c.JSON(status, gin.H{
			"message": msg,
			"detail":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Product successfully deleted",
	})
}

// Recover product (admin only)
func (ctl *ProductController) RecoverProduct(c *gin.Context) {
	idStr := c.Param("productId")
	idUint64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid product ID",
			"detail":  err.Error(),
		})
		return
	}
	id := uint(idUint64)

	if err := ctl.service.RecoverProduct(c.Request.Context(), id); err != nil {
		status := http.StatusInternalServerError
		msg := "Failed to recover product, try again later"

		if err.Error() == "product not found" {
			status = http.StatusNotFound
			msg = "Product not found"
		}

		c.JSON(status, gin.H{
			"message": msg,
			"detail":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Product successfully recovered"})
}
