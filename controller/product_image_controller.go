package controller

import (
	"fmt"
	"net/http"
	"path/filepath"

	"ecommerce/service"

	"github.com/gin-gonic/gin"
)

type ProductImageController struct {
	service service.ProductImageService
}

func NewProductImageController(s service.ProductImageService) *ProductImageController {
	return &ProductImageController{service: s}
}

// UploadImage godoc
// @Summary      Upload product image
// @Description  Upload image for product (admin only, max 3 images per product)
// @Tags         Product Images
// @Security     BearerAuth
// @Accept       multipart/form-data
// @Produce      json
// @Param        productId   path      int     true   "Product ID"
// @Param        image       formData  file    true   "Image file to upload"
// @Param        is_primary  formData  bool    false  "Set as primary image"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Example 400 {json} Error Example:
// {
//   "message": "Invalid product ID",
//   "detail": "strconv.ParseUint: parsing \"abc\": invalid syntax"
// }
// @Failure      500  {object}  map[string]interface{}
// @Example 500 {json} Error Example:
// {
//   "message": "Failed to save image, try again later",
//   "detail": "some error message"
// }
// @Router       /admin/products/{productId}/images [post]
func (ctl *ProductImageController) UploadImage(c *gin.Context) {
	productID := c.Param("productId")
	var pid uint
	if _, err := fmt.Sscan(productID, &pid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid product ID",
			"detail":  err.Error(),
		})
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Missing image",
			"detail":  err.Error(),
		})
		return
	}

	isPrimary := c.DefaultPostForm("is_primary", "false") == "true"

	tempPath := filepath.Join("/tmp", file.Filename)
	if err := c.SaveUploadedFile(file, tempPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to save image, try again later",
			"detail":  err.Error(),
		})
		return
	}

	img, err := ctl.service.Upload(c, pid, tempPath, isPrimary)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to upload image",
			"detail":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Image successfully uploaded",
		"detail":  img,
	})
}

// DeleteImage godoc
// @Summary      Delete product image
// @Description  Soft delete product image by ID (admin only)
// @Tags         Product Images
// @Security     BearerAuth
// @Produce      json
// @Param        imageId  path      int  true  "Image ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Example 400 {json} Error Example:
// {
//   "message": "Invalid image ID",
//   "detail": "strconv.ParseUint: parsing \"xyz\": invalid syntax"
// }
// @Failure      404  {object}  map[string]interface{}
// @Example 404 {json} Error Example:
// {
//   "message": "Image not found",
//   "detail": "record not found"
// }
// @Failure      500  {object}  map[string]interface{}
// @Example 500 {json} Error Example:
// {
//   "message": "Failed to delete image, try again later",
//   "detail": "some error message"
// }
// @Router       /admin/images/{imageId} [delete]
func (ctl *ProductImageController) DeleteImage(c *gin.Context) {
	var id uint
	if _, err := fmt.Sscan(c.Param("imageId"), &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid image ID",
			"detail":  err.Error(),
		})
		return
	}

	if err := ctl.service.Delete(c, id); err != nil {
		status := http.StatusInternalServerError
		msg := "Failed to delete image, try again later"

		if err.Error() == "image not found" {
			status = http.StatusNotFound
			msg = "Image not found"
		}

		c.JSON(status, gin.H{
			"message": msg,
			"detail":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Image successfully deleted",
	})
}

// RecoverImage godoc
// @Summary      Recover deleted product image
// @Description  Restore soft-deleted image (admin only)
// @Tags         Product Images
// @Security     BearerAuth
// @Produce      json
// @Param        imageId  path      int  true  "Image ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Example 400 {json} Error Example:
// {
//   "message": "Invalid image ID",
//   "detail": "strconv.ParseUint: parsing \"abc\": invalid syntax"
// }
// @Failure      404  {object}  map[string]interface{}
// @Example 404 {json} Error Example:
// {
//   "message": "Image not found",
//   "detail": "record not found"
// }
// @Failure      500  {object}  map[string]interface{}
// @Example 500 {json} Error Example:
// {
//   "message": "Failed to recover image, try again later",
//   "detail": "some error message"
// }
// @Router       /admin/images/{imageId}/recover [post]
func (ctl *ProductImageController) RecoverImage(c *gin.Context) {
	var id uint
	if _, err := fmt.Sscan(c.Param("imageId"), &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid image ID",
			"detail":  err.Error(),
		})
		return
	}

	if err := ctl.service.Recover(c, id); err != nil {
		status := http.StatusInternalServerError
		msg := "Failed to recover image, try again later"

		if err.Error() == "image not found" {
			status = http.StatusNotFound
			msg = "Image not found"
		}

		c.JSON(status, gin.H{
			"message": msg,
			"detail":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Image successfully recovered",
	})
}
