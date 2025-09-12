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

// GetProducts godoc
// @Summary      Get list of products
// @Description  Return list of products with optional search, category filter, pagination
// @Tags         Products
// @Produce      json
// @Param        search   query     string  false  "Search by product name"
// @Param        category query     string  false  "Filter by category name"
// @Param        limit    query     int     false  "Limit number of results"   default(10)
// @Param        offset   query     int     false  "Offset for pagination"     default(0)
// @Success      200      {array}   entity.Product
// @Failure      500      {object}  map[string]interface{}
// @Example 200 {json} Success Example:
// [
//   {
//     "id": 1,
//     "name": "Roasted Almond",
//     "description": "Crunchy and delicious roasted almonds",
//     "price": 150000,
//     "stock": 50,
//     "expiry_year": 2026,
//     "created_at": "2025-09-12T21:00:00Z",
//     "updated_at": "2025-09-12T21:00:00Z",
//     "images": [
//       {"id": 1, "product_id": 1, "image_url": "https://res.cloudinary.com/.../image.jpg", "is_primary": true, "created_at": "2025-09-12T21:00:00Z"}
//     ],
//     "categories": [
//       {"id": 1, "name": "Nuts", "created_at": "2025-09-10T10:00:00Z", "updated_at": "2025-09-10T10:00:00Z"}
//     ]
//   },
//   {
//     "id": 2,
//     "name": "Chia Seed",
//     "description": "Healthy chia seeds for smoothies",
//     "price": 350000,
//     "stock": 10,
//     "expiry_year": 2025,
//     "created_at": "2025-09-10T09:30:00Z",
//     "updated_at": "2025-09-10T09:30:00Z",
//     "images": [],
//     "categories": [
//       {"id": 2, "name": "Seeds", "created_at": "2025-09-08T08:00:00Z", "updated_at": "2025-09-08T08:00:00Z"}
//     ]
//   }
// ]
// @Example 500 {json} Error Example:
// {
//   "message": "Failed to fetch products, try again later",
//   "detail": "some error message"
// }
// @Router       /products [get]
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

// GetProductByID godoc
// @Summary      Get product by ID
// @Description  Return single product details by ID
// @Tags         Products
// @Produce      json
// @Param        productId   path      int  true  "Product ID"
// @Success      200         {object}  entity.Product
// @Failure      400         {object}  map[string]interface{}
// @Failure      404         {object}  map[string]interface{}
// @Example 200 {json} Success Example:
// {
//   "id": 1,
//   "name": "Roasted Almond",
//   "description": "Crunchy and delicious roasted almonds",
//   "price": 150000,
//   "stock": 50,
//   "expiry_year": 2026,
//   "created_at": "2025-09-12T21:00:00Z",
//   "updated_at": "2025-09-12T21:00:00Z",
//   "images": [
//     {"id": 1, "product_id": 1, "image_url": "https://res.cloudinary.com/.../image.jpg", "is_primary": true, "created_at": "2025-09-12T21:00:00Z"}
//   ],
//   "categories": [
//     {"id": 1, "name": "Nuts", "created_at": "2025-09-10T10:00:00Z", "updated_at": "2025-09-10T10:00:00Z"}
//   ]
// }
// @Example 400 {json} Error Example:
// {
//   "message": "Invalid product ID",
//   "detail": "parsing error detail"
// }
// @Example 404 {json} Error Example:
// {
//   "message": "Product not found",
//   "detail": "some error message"
// }
// @Router       /products/{productId} [get]
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

// CreateProduct godoc
// @Summary      Create new product
// @Description  Admin can create a new product
// @Tags         Products
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request   body      createProductReq  true  "Product data"
// @Success      200       {object}  map[string]interface{}
// @Failure      400       {object}  map[string]interface{}
// @Failure      500       {object}  map[string]interface{}
// @Example 200 {json} Success Example:
// {
//   "message": "Product successfully created",
//   "data": {
//     "id": 3,
//     "name": "Organic Honey",
//     "description": "Pure organic honey from local farms",
//     "price": 250000,
//     "stock": 20,
//     "expiry_year": 2027,
//     "created_at": "2025-09-12T21:15:00Z",
//     "updated_at": "2025-09-12T21:15:00Z",
//     "images": [],
//     "categories": [
//       {"id": 3, "name": "Sweeteners", "created_at": "2025-09-11T12:00:00Z", "updated_at": "2025-09-11T12:00:00Z"}
//     ]
//   }
// }
// @Example 400 {json} Error Example:
// {
//   "message": "Invalid input data",
//   "detail": "binding error detail"
// }
// @Example 500 {json} Error Example:
// {
//   "message": "Failed to create product, try again later",
//   "detail": "some error message"
// }
// @Router       /admin/products [post]
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

// UpdateProduct godoc
// @Summary      Update existing product
// @Description  Admin can update product fields by ID
// @Tags         Products
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        productId   path      int               true  "Product ID"
// @Param        request     body      updateProductReq  true  "Updated product data"
// @Success      200         {object}  map[string]interface{}
// @Failure      400         {object}  map[string]interface{}
// @Failure      404         {object}  map[string]interface{}
// @Failure      500         {object}  map[string]interface{}
// @Example 200 {json} Success Example:
// {
//   "message": "Product successfully updated",
//   "data": {
//     "id": 3,
//     "name": "Organic Honey Premium",
//     "description": "Premium organic honey from local farms",
//     "price": 300000,
//     "stock": 15,
//     "expiry_year": 2027,
//     "created_at": "2025-09-12T21:15:00Z",
//     "updated_at": "2025-09-12T21:45:00Z",
//     "images": [],
//     "categories": [
//       {"id": 3, "name": "Sweeteners", "created_at": "2025-09-11T12:00:00Z", "updated_at": "2025-09-11T12:00:00Z"}
//     ]
//   }
// }
// @Example 400 {json} Error Example:
// {
//   "message": "Invalid product ID",
//   "detail": "parsing error detail"
// }
// @Example 404 {json} Error Example:
// {
//   "message": "Product not found",
//   "detail": "some error message"
// }
// @Example 500 {json} Error Example:
// {
//   "message": "Failed to update product, try again later",
//   "detail": "some error message"
// }
// @Router       /admin/products/{productId} [put]
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

// DeleteProduct godoc
// @Summary      Delete product
// @Description  Admin can soft-delete a product by ID
// @Tags         Products
// @Security     BearerAuth
// @Produce      json
// @Param        productId   path      int  true  "Product ID"
// @Success      200         {object}  map[string]interface{}
// @Failure      400         {object}  map[string]interface{}
// @Failure      404         {object}  map[string]interface{}
// @Failure      500         {object}  map[string]interface{}
// @Example 200 {json} Success Example:
// {
//   "message": "Product successfully deleted"
// }
// @Example 400 {json} Error Example:
// {
//   "message": "Invalid product ID",
//   "detail": "parsing error detail"
// }
// @Example 404 {json} Error Example:
// {
//   "message": "Product not found",
//   "detail": "some error message"
// }
// @Example 500 {json} Error Example:
// {
//   "message": "Failed to delete product, try again later",
//   "detail": "some error message"
// }
// @Router       /admin/products/{productId} [delete]
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

// RecoverProduct godoc
// @Summary      Recover deleted product
// @Description  Admin can recover a previously deleted product
// @Tags         Products
// @Security     BearerAuth
// @Produce      json
// @Param        productId   path      int  true  "Product ID"
// @Success      200         {object}  map[string]interface{}
// @Failure      400         {object}  map[string]interface{}
// @Failure      404         {object}  map[string]interface{}
// @Failure      500         {object}  map[string]interface{}
// @Example 200 {json} Success Example:
// {
//   "message": "Product successfully recovered"
// }
// @Example 400 {json} Error Example:
// {
//   "message": "Invalid product ID",
//   "detail": "parsing error detail"
// }
// @Example 404 {json} Error Example:
// {
//   "message": "Product not found",
//   "detail": "some error message"
// }
// @Example 500 {json} Error Example:
// {
//   "message": "Failed to recover product, try again later",
//   "detail": "some error message"
// }
// @Router       /admin/products/{productId}/recover [patch]
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
