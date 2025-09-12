package controller

import (
	"net/http"
	"strconv"
	"time"

	"ecommerce/entity"
	"ecommerce/service"

	"github.com/gin-gonic/gin"
)

type CategoryController struct {
	service service.CategoryService
}

func NewCategoryController(s service.CategoryService) *CategoryController {
	return &CategoryController{service: s}
}

// Struct request for endpoint create category dan update
type createCategoryReq struct {
	Name string `json:"name" binding:"required"`
}

type updateCategoryReq struct {
	Name *string `json:"name,omitempty"`
}

// GetCategories godoc
// @Summary      Get all categories
// @Description  Return list of categories (public)
// @Tags         Categories
// @Produce      json
// @Param        limit    query     int  false  "Limit number of results"   default(10)
// @Param        offset   query     int  false  "Offset for pagination"     default(0)
// @Success      200      {array}   entity.ProductCategory
// @Failure      500      {object}  map[string]interface{}
// @Example 200 {json} Success Example:
// [
//   {
//     "id": 1,
//     "name": "Nuts",
//     "created_at": "2025-09-12T21:00:00Z",
//     "updated_at": "2025-09-12T21:10:00Z",
//   },
//   {
//     "id": 2,
//     "name": "Seeds",
//     "created_at": "2025-09-10T09:30:00Z",
//     "updated_at": "2025-09-10T09:45:00Z",
//   }
// ]
// @Example 500 {json} Error Example:
// {
//   "message": "Failed to fetch categories, try again later",
//   "detail": "some error message"
// }
// @Router       /categories [get]
func (ctl *CategoryController) GetCategories(c *gin.Context) {
	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")

	limit, _ := strconv.Atoi(limitStr)
	if limit == 0 {
		limit = 10
	}
	offset, _ := strconv.Atoi(offsetStr)

	categories, err := ctl.service.GetCategories(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch categories, try again later",
			"detail":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, categories)
}

// CreateCategory godoc
// @Summary      Create category
// @Description  Create a new category (admin only)
// @Tags         Categories
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body      createCategoryReq  true  "Create Category Request"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Example 200 {json} Success Example:
// {
//   "message": "Category successfully created",
//   "data": {
//     "id": 3,
//     "name": "Dried Fruits",
//     "created_at": "2025-09-12T22:00:00Z",
//     "updated_at": "2025-09-12T22:00:00Z",
//   }
// }
// @Example 400 {json} Error Example:
// {
//   "message": "Invalid input data",
//   "detail": "binding error detail"
// }
// @Example 500 {json} Error Example:
// {
//   "message": "Failed to create category, try again later",
//   "detail": "some error message"
// }
// @Router       /admin/categories [post]
func (ctl *CategoryController) CreateCategory(c *gin.Context) {
	var req createCategoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid input data",
			"detail":  err.Error(),
		})
		return
	}

	category := &entity.ProductCategory{
		Name: req.Name,
	}
	if err := ctl.service.CreateCategory(c.Request.Context(), category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create category, try again later",
			"detail":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Category successfully created",
		"data":    category,
	})
}

// UpdateCategory godoc
// @Summary      Update category
// @Description  Update an existing category (admin only)
// @Tags         Categories
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path      int                true  "Category ID"
// @Param        request  body      updateCategoryReq  true  "Update Category Request"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Failure      404      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Example 400 {json} Error Example:
// {
//   "message": "Invalid category ID",
//   "detail": "parsing error detail"
// }
// @Example 404 {json} Error Example:
// {
//   "message": "Category not found",
//   "detail": "gorm.ErrRecordNotFound"
// }
// @Example 500 {json} Error Example:
// {
//   "message": "Failed to update category, try again later",
//   "detail": "some error message"
// }
// @Router       /admin/categories/{id} [patch]
func (ctl *CategoryController) UpdateCategory(c *gin.Context) {
	idStr := c.Param("id")
	idUint64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid category ID",
			"detail":  err.Error(),
		})
		return
	}
	id := uint(idUint64)

	var req updateCategoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid input data",
			"detail":  err.Error(),
		})
		return
	}

	existing, err := ctl.service.GetCategoryByIDIncludeDeleted(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Category not found",
			"detail":  err.Error(),
		})
		return
	}

	if req.Name != nil {
		existing.Name = *req.Name
	}
	existing.UpdatedAt = time.Now()

	if err := ctl.service.UpdateCategory(c.Request.Context(), existing); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update category, try again later",
			"detail":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Category successfully updated",
		"data":    existing,
	})
}

// DeleteCategory godoc
// @Summary      Delete category
// @Description  Soft delete a category (admin only)
// @Tags         Categories
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "Category ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Example 200 {json} Success Example:
// {
//   "message": "Category successfully deleted"
// }
// @Example 400 {json} Error Example:
// {
//   "message": "Invalid category ID",
//   "detail": "parsing error detail"
// }
// @Example 404 {json} Error Example:
// {
//   "message": "Category not found",
//   "detail": "category not found"
// }
// @Example 500 {json} Error Example:
// {
//   "message": "Failed to delete category, try again later",
//   "detail": "some error message"
// }
// @Router       /admin/categories/{id} [delete]
func (ctl *CategoryController) DeleteCategory(c *gin.Context) {
	idStr := c.Param("id")
	idUint64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid category ID",
			"detail":  err.Error(),
		})
		return
	}
	id := uint(idUint64)

	if err := ctl.service.DeleteCategory(c.Request.Context(), id); err != nil {
		status := http.StatusInternalServerError
		msg := "Failed to delete category, try again later"

		if err.Error() == "category not found" {
			status = http.StatusNotFound
			msg = "Category not found"
		}

		c.JSON(status, gin.H{
			"message": msg,
			"detail":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Category successfully deleted",
	})
}

// RecoverCategory godoc
// @Summary      Recover category
// @Description  Restore a soft-deleted category (admin only)
// @Tags         Categories
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "Category ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Example 200 {json} Success Example:
// {
//   "message": "Category successfully recovered"
// }
// @Example 400 {json} Error Example:
// {
//   "message": "Invalid category ID",
//   "detail": "parsing error detail"
// }
// @Example 404 {json} Error Example:
// {
//   "message": "Category not found",
//   "detail": "category not found"
// }
// @Example 500 {json} Error Example:
// {
//   "message": "Failed to recover category, try again later",
//   "detail": "some error message"
// }
// @Router       /admin/categories/{id}/recover [post]
func (ctl *CategoryController) RecoverCategory(c *gin.Context) {
	idStr := c.Param("id")
	idUint64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid category ID",
			"detail":  err.Error(),
		})
		return
	}
	id := uint(idUint64)

	if err := ctl.service.RecoverCategory(c.Request.Context(), id); err != nil {
		status := http.StatusInternalServerError
		msg := "Failed to recover category, try again later"

		if err.Error() == "category not found" {
			status = http.StatusNotFound
			msg = "Category not found"
		}

		c.JSON(status, gin.H{
			"message": msg,
			"detail":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Category successfully recovered",
	})
}
