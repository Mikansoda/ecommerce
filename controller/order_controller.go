package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"ecommerce/entity"
	"ecommerce/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OrderController struct {
	service service.OrderService
}

func NewOrderController(s service.OrderService) *OrderController {
	return &OrderController{service: s}
}

// GET orders (admin only)
func (ctl *OrderController) GetOrders(c *gin.Context) {
	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")
	limit, _ := strconv.Atoi(limitStr)
	if limit == 0 {
		limit = 10
	}
	offset, _ := strconv.Atoi(offsetStr)

	status := c.Query("status")

	var (
		orders []entity.Order
		err    error
	)
	if status != "" {
		orders, err = ctl.service.GetOrdersByStatus(c.Request.Context(), status, limit, offset)
	} else {
		orders, err = ctl.service.GetOrders(c.Request.Context(), limit, offset)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch orders, try again later",
			"detail":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, orders)
}

// GET orders (self-requested by user)
func (ctl *OrderController) GetUserOrders(c *gin.Context) {
	userIDStr := c.GetString("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid user ID",
			"detail":  err.Error(),
		})
		return
	}
	orders, err := ctl.service.GetOrdersByUser(c.Request.Context(), userID)
	if err != nil {
		if err.Error() == "orders not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "No order found, let's start ordering",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch your orders, try again later",
			"detail":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, orders)
}

// Create Order (self-created by user)
func (ctl *OrderController) CreateOrder(c *gin.Context) {
	var req struct {
		UserID    string `json:"user_id"`
		AddressID string `json:"address_id"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
			"detail":  err.Error(),
		})
		return
	}
	uid, err := uuid.Parse(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid user ID",
			"detail":  err.Error(),
		})
		return
	}
	aid, err := uuid.Parse(req.AddressID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid address ID",
			"detail":  err.Error(),
		})
		return
	}
	order, err := ctl.service.CreateOrder(c.Request.Context(), uid, aid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to create order, try again later",
			"detail":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"order": order})
}

// Update order status (admin only)
func (ctl *OrderController) UpdateOrderStatus(c *gin.Context) {
	orderIDStr := c.Param("id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid order ID",
			"detail":  err.Error(),
		})
		return
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid input data",
			"detail":  err.Error(),
		})
		return
	}
	if err := ctl.service.UpdateOrderStatus(c.Request.Context(), orderID, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to update order status, try again later",
			"detail":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("order %s status updated to %s", orderID, req.Status),
	})
}
