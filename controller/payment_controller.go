package controller

import (
	"bytes"
	"ecommerce/service"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PaymentController struct {
	service service.PaymentService
}

func NewPaymentController(s service.PaymentService) *PaymentController {
	return &PaymentController{service: s}
}

// Request body struct
type createPaymentReq struct {
	OrderID uuid.UUID `json:"order_id" binding:"required"`
}

// GetPayments godoc
// @Summary      Get all payments
// @Description  Return list of all payments (admin only)
// @Tags         Payments
// @Security     BearerAuth
// @Produce      json
// @Success      200  {array}   entity.Payment
// @Failure      500  {object}  map[string]interface{}
// @Example 500 {json} Error Example:
// {
//   "message": "Failed to fetch payments, try again later",
//   "detail": "some error message"
// }
// @Router       /admin/payments [get]
func (ctl *PaymentController) GetPayments(c *gin.Context) {
	payments, err := ctl.service.GetAllPayments(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch payments, try again later",
			"detail":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, payments)
}

// GetUserPayments godoc
// @Summary      Get user payments
// @Description  Return list of payments for the authenticated user
// @Tags         Payments
// @Security     BearerAuth
// @Produce      json
// @Success      200  {array}   entity.Payment
// @Failure      500  {object}  map[string]interface{}
// @Example 500 {json} Error Example:
// {
//   "message": "Failed to fetch payments, try again later",
//   "detail": "some error message"
// }
// @Router       /user/payments [get]
func (ctl *PaymentController) GetUserPayments(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	payments, err := ctl.service.GetPaymentsByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch payments, try again later",
			"detail":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, payments)
}

// CreatePayment godoc
// @Summary      Create new payment
// @Description  Create a payment invoice via Xendit (user only)
// @Tags         Payments
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request   body      createPaymentReq  true  "Order ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Example 200 {json} Success Example:
// {
//   "message": "Payment created",
//   "payment_id": "uuid-here",
//   "xendit_id": "inv-12345",
//   "invoice_url": "https://checkout.xendit.co/invoices/inv-12345",
//   "status": "pending"
// }
// @Example 400 {json} Error Example:
// {
//   "message": "Invalid order ID",
//   "detail": "some error message"
// }
// @Example 404 {json} Error Example:
// {
//   "message": "order not found",
//   "detail": "record not found"
// }
// @Example 500 {json} Error Example:
// {
//   "message": "Failed to call payment gateway, try again later",
//   "detail": "some error message"
// }
// @Router       /user/payments/xendit [post]
func (ctl *PaymentController) CreatePayment(c *gin.Context) {
	var req createPaymentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid order ID",
			"detail":  err.Error(),
		})
		return
	}

	order, err := ctl.service.GetOrderByID(c.Request.Context(), req.OrderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "order not found",
			"detail":  err.Error(),
		})
		return
	}

	// Panggil Xendit API (sandbox)
	xenditKey := os.Getenv("XENDIT_API_KEY")
	if xenditKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Xendit API key not set",
		})
		return
	}

	payload := map[string]interface{}{
		"external_id":          order.ID.String(),
		"amount":               order.TotalAmount,
		"payer_email":          order.User.Email,
		"description":          fmt.Sprintf("Payment for order %s", order.ID.String()),
		"success_redirect_url": "https://example.com/success",
		"failure_redirect_url": "https://example.com/failure",
	}

	data, _ := json.Marshal(payload)
	reqXendit, _ := http.NewRequest("POST", "https://api.xendit.co/v2/invoices", bytes.NewBuffer(data))
	reqXendit.SetBasicAuth(xenditKey, "")
	reqXendit.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(reqXendit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to call payment gateway, try again later",
			"detail":  err.Error(),
		})
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Invalid response from payment gateway",
			"detail":  err.Error(),
		})
		return
	}

	invoiceID, ok := result["id"].(string)
	invoiceURL, _ := result["invoice_url"].(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Invalid invoice id from payment gateway",
		})
		return
	}

	// Save to DB
	payment, err := ctl.service.CreatePayment(c.Request.Context(), order, invoiceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to save payment",
			"detail":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Payment created",
		"payment_id":  payment.ID,
		"xendit_id":   invoiceID,
		"invoice_url": invoiceURL,
		"status":      payment.Status,
	})
}

// XenditWebhook godoc
// @Summary      Xendit webhook
// @Description  Handle Xendit payment status update (admin only)
// @Tags         Payments
// @Accept       json
// @Produce      json
// @Param        payload   body      map[string]interface{}  true  "Webhook payload"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Example 200 {json} Success Example:
// {
//   "message": "Payment successfully updated"
// }
// @Example 400 {json} Error Example:
// {
//   "message": "Invalid payload",
//   "detail": "some error message"
// }
// @Example 500 {json} Error Example:
// {
//   "message": "Failed to update payment",
//   "detail": "some error message"
// }
// @Router       /admin/payments/webhook/xendit [post]
func (ctl *PaymentController) XenditWebhook(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid payload",
			"detail":  err.Error(),
		})
		return
	}

	invoiceID, _ := payload["id"].(string)
	status, _ := payload["status"].(string)

	if err := ctl.service.UpdatePaymentStatus(c.Request.Context(), invoiceID, status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update payment",
			"detail":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment successfully updated"})
}
