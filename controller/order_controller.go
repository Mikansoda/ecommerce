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

// GetOrders godoc
// @Summary      Get all orders
// @Description  Return list of all orders (admin only)
// @Tags         Orders
// @Security     BearerAuth
// @Produce      json
// @Param        status   query     string  false  "Filter by order status (e.g. pending, completed)"
// @Param        limit    query     int     false  "Limit number of results"   default(10)
// @Param        offset   query     int     false  "Offset for pagination"     default(0)
// @Success      200      {array}   entity.Order
// @Failure      500      {object}  map[string]interface{}
// @Example 500 {json} Error Example:
// {
//   "message": "Failed to fetch orders, try again later",
//   "detail": "some error message"
// }
// @Router       /admin/orders [get]
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

// GetUserOrders godoc
// @Summary      Get user orders
// @Description  Return list of orders belonging to the logged-in user
// @Tags         Orders
// @Security     BearerAuth
// @Produce      json
// @Success      200  {array}   entity.Order
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Example 400 {json} Error Example:
// {
//   "message": "invalid user ID",
//   "detail": "parsing error detail"
// }
// @Example 404 {json} Error Example:
// {
//   "message": "No order found, let's start ordering"
// }
// @Example 500 {json} Error Example:
// {
//   "message": "Failed to fetch your orders, try again later",
//   "detail": "some error message"
// }
// @Router       /user/orders [get]
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

// CreateOrder godoc
// @Summary      Create a new order
// @Description  Create an order from the user's cart
// @Tags         Orders
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body      object  true  "Create Order Request"  Example({"user_id": "uuid-string", "address_id": "uuid-string"})
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Example 400 {json} Error Example:
// {
//   "message": "invalid request",
//   "detail": "binding error detail"
// }
// @Example 500 {json} Error Example:
// {
//   "message": "failed to create order, try again later",
//   "detail": "some error message"
// }
// @Router       /user/orders [post]
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

// UpdateOrderStatus godoc
// @Summary      Update order status
// @Description  Update the status of an existing order (admin only)
// @Tags         Orders
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path      string  true  "Order ID (UUID)"
// @Param        request  body      object  true  "Update Order Status Request"  Example({"status": "completed"})
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Example 200 {json} Success Example:
// {
//   "message": "order 123e4567-e89b-12d3-a456-426614174000 status updated to completed"
// }
// @Example 400 {json} Error Example:
// {
//   "message": "invalid order ID",
//   "detail": "parsing error detail"
// }
// @Example 500 {json} Error Example:
// {
//   "message": "failed to update order status, try again later",
//   "detail": "some error message"
// }
// @Router       /admin/orders/{id}/status [put]
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
