package controller

import (
	"net/http"

	"ecommerce/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CartController struct {
	service service.CartService
}

func NewCartController(s service.CartService) *CartController {
	return &CartController{service: s}
}

// Struct request for endpoint add item to cart
type addCartItemReq struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required"`
}

// GET cart (self-requested by user)
func (ctl *CartController) GetCart(c *gin.Context) {
	userIDStr := c.GetString("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid user ID",
			"detail":  err.Error(),
		})
		return
	}
	cart, err := ctl.service.GetCart(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch your cart, try again later",
			"detail":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, cart)
}

// Add item to cart
func (ctl *CartController) AddItem(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("userID"))
	var req addCartItemReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid input data",
			"detail":  err.Error(),
		})
		return
	}

	item, err := ctl.service.AddItem(c.Request.Context(), userID, req.ProductID, req.Quantity)
	if err != nil {
		status := http.StatusInternalServerError
		msg := "Failed to add cart item, try again later"

		if err.Error() == "product not found" {
			status = http.StatusNotFound
			msg = "Product not found"
		} else if err.Error() == "forbidden" {
			status = http.StatusForbidden
			msg = "You are not allowed to add this product"
		} else if err.Error() == "invalid quantity" {
			status = http.StatusBadRequest
			msg = "Invalid quantity"
		}

		c.JSON(status, gin.H{
			"message": msg,
			"detail":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, item)
}

// Delete item in cart (self-delete by user)
func (ctl *CartController) DeleteItem(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("userID"))
	itemIDStr := c.Param("id")
	itemID, err := uuid.Parse(itemIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid cart item id",
			"detail":  err.Error(),
		})
		return
	}

	if err := ctl.service.RemoveItem(c.Request.Context(), userID, itemID); err != nil {
		status := http.StatusInternalServerError
		msg := "Failed to delete cart item"

		if err.Error() == "cart item not found" {
			status = http.StatusNotFound
			msg = "Cart item not found"
		} else if err.Error() == "forbidden" {
			status = http.StatusForbidden
			msg = "You are not allowed to delete this cart item"
		}
		c.JSON(status, gin.H{
			"message": msg,
			"detail":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Cart item deleted",
	})
}
