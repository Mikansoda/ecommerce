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

// GetCart godoc
// @Summary      Get cart
// @Description  Return current user's cart
// @Tags         Cart
// @Security     BearerAuth
// @Produce      json
// @Success      200      {object}  entity.Cart
// @Failure      400      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Example 200 {json} Success Example:
// {
//   "id": "1a2b3c4d-5678-90ab-cdef-1234567890ab",
//   "user_id": "9b8c7d6e-5432-10fe-dcba-0987654321fe",
//   "created_at": "2025-09-12T21:00:00Z",
//   "updated_at": "2025-09-12T21:10:00Z",
//   "items": [
//     {
//       "id": "2b3c4d5e-6789-01bc-def0-2345678901bc",
//       "cart_id": "1a2b3c4d-5678-90ab-cdef-1234567890ab",
//       "product_id": 1,
//       "quantity": 2,
//       "product": {
//         "id": 1,
//         "name": "Roasted Almond",
//         "stock": 50,
//         "price": 150000
//       }
//     },
//     {
//       "id": "3c4d5e6f-7890-12cd-ef01-3456789012cd",
//       "cart_id": "1a2b3c4d-5678-90ab-cdef-1234567890ab",
//       "product_id": 2,
//       "quantity": 1,
//       "product": {
//         "id": 2,
//         "name": "Chia Seed",
//         "stock": 10,
//         "price": 350000
//       }
//     }
//   ]
// }
// @Example 400 {json} Error Example:
// {
//   "message": "Invalid user ID",
//   "detail": "parsing error detail"
// }
// @Example 500 {json} Error Example:
// {
//   "message": "Failed to fetch your cart, try again later",
//   "detail": "some error message"
// }
// @Router       /user/cart [get]
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

// AddItem godoc
// @Summary      Add item to cart
// @Description  Add a product into the current user's cart
// @Tags         Cart
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body      addCartItemReq  true  "Add Cart Item Request"
// @Success      201      {object}  entity.CartItem
// @Failure      400      {object}  map[string]interface{}
// @Failure      403      {object}  map[string]interface{}
// @Failure      404      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Example 201 {json} Success Example:
// {
//   "id": "2b3c4d5e-6789-01bc-def0-2345678901bc",
//   "cart_id": "1a2b3c4d-5678-90ab-cdef-1234567890ab",
//   "product_id": 1,
//   "quantity": 2,
//   "product": {
//     "id": 1,
//     "name": "Roasted Almond",
//     "stock": 50,
//     "price": 150000
//   }
// }
// @Example 400 {json} Error Example:
// {
//   "message": "Invalid input data",
//   "detail": "binding error detail"
// }
// @Example 404 {json} Error Example:
// {
//   "message": "Product not found",
//   "detail": "product not found"
// }
// @Example 500 {json} Error Example:
// {
//   "message": "Failed to add cart item, try again later",
//   "detail": "some error message"
// }
// @Router       /user/cart/items [post]
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

// DeleteItem godoc
// @Summary      Delete item from cart
// @Description  Remove a cart item owned by the current user
// @Tags         Cart
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Cart Item ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Example 200 {json} Success Example:
// {
//   "message": "Cart item deleted"
// }
// @Example 400 {json} Error Example:
// {
//   "message": "Invalid cart item id",
//   "detail": "uuid parsing error detail"
// }
// @Example 404 {json} Error Example:
// {
//   "message": "Cart item not found",
//   "detail": "cart item not found"
// }
// @Example 500 {json} Error Example:
// {
//   "message": "Failed to delete cart item",
//   "detail": "some error message"
// }
// @Router       /user/cart/items/{id} [delete]
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
