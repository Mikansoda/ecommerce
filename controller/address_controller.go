package controller

import (
	"net/http"
	"strconv"
	"time"

	"ecommerce/entity"
	"ecommerce/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AddressController struct {
	service service.AddressService
}

func NewAddressController(s service.AddressService) *AddressController {
	return &AddressController{service: s}
}

// Struct request for endpoint create address and update
type createAddressReq struct {
	ReceiverName string `json:"receiver_name" binding:"required"`
	PhoneNumber  string `json:"phone_number" binding:"required"`
	AddressLine  string `json:"address_line" binding:"required"`
	City         string `json:"city" binding:"required"`
	Province     string `json:"province" binding:"required"`
	PostalCode   string `json:"postal_code" binding:"required"`
}

type updateAddressReq struct {
	ReceiverName *string `json:"receiver_name,omitempty"`
	PhoneNumber  *string `json:"phone_number,omitempty"`
	AddressLine  *string `json:"address_line,omitempty"`
	City         *string `json:"city,omitempty"`
	Province     *string `json:"province,omitempty"`
	PostalCode   *string `json:"postal_code,omitempty"`
}

// GET addresses (admin only)
func (ctl *AddressController) GetAddresses(c *gin.Context) {
	search := c.Query("search")
	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")

	limit, _ := strconv.Atoi(limitStr)
	if limit == 0 {
		limit = 10
	}
	offset, _ := strconv.Atoi(offsetStr)

	addresses, err := ctl.service.GetAddresses(c.Request.Context(), search, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch addresses, try again later",
			"detail":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, addresses)
}

// GET address (self-requested by user)
func (ctl *AddressController) GetUserAddress(c *gin.Context) {
	userIDStr := c.GetString("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid user ID",
			"detail":  err.Error(),
		})
		return
	}
	address, err := ctl.service.GetAddressByUser(c.Request.Context(), userID)
	if err != nil {
		if err.Error() == "address not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "No address found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch your address, try again later",
			"detail":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, address)
}

// Create address (self-created by user)
func (ctl *AddressController) CreateAddress(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("userID"))
	var req createAddressReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid input data",
			"detail":  err.Error(),
		})
		return
	}
	address := &entity.Address{
		UserID:       userID,
		ReceiverName: req.ReceiverName,
		PhoneNumber:  req.PhoneNumber,
		AddressLine:  req.AddressLine,
		City:         req.City,
		Province:     req.Province,
		PostalCode:   req.PostalCode,
	}
	if err := ctl.service.CreateAddress(c.Request.Context(), address); err != nil {
		statusCode := http.StatusInternalServerError
		msg := "Failed to create address, try again later"

		if err.Error() == "user already has an address" {
			statusCode = http.StatusBadRequest
			msg = "You already have an address"
		}

		c.JSON(statusCode, gin.H{
			"message": msg,
			"detail":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Address successfully created",
		"data":    address,
	})
}

// Update address (self-updated by user)
func (ctl *AddressController) UpdateAddress(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("userID"))
	var req updateAddressReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid input data",
			"detail":  err.Error(),
		})
		return
	}

	existing, err := ctl.service.GetAddressByUser(c.Request.Context(), userID)
	if err != nil {
		if err.Error() == "address not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Address not found",
			})
			return
		}
	}
	if req.ReceiverName != nil {
		existing.ReceiverName = *req.ReceiverName
	}
	if req.PhoneNumber != nil {
		existing.PhoneNumber = *req.PhoneNumber
	}
	if req.AddressLine != nil {
		existing.AddressLine = *req.AddressLine
	}
	if req.City != nil {
		existing.City = *req.City
	}
	if req.Province != nil {
		existing.Province = *req.Province
	}
	if req.PostalCode != nil {
		existing.PostalCode = *req.PostalCode
	}
	existing.UpdatedAt = time.Now()

	if err := ctl.service.UpdateAddress(c.Request.Context(), existing); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update address, try again later",
			"detail":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Address successfully updated",
		"data":    existing,
	})
}

// Delete address (self-delete by user)
func (ctl *AddressController) DeleteAddress(c *gin.Context) {
	id := c.Param("id")
	userID, _ := uuid.Parse(c.GetString("userID"))

	if err := ctl.service.DeleteAddress(c.Request.Context(), id, userID); err != nil {
		status := http.StatusInternalServerError
		msg := "Failed to delete address, try again later"

		if err.Error() == "address not found" {
			status = http.StatusNotFound
			msg = "Address not found"
		} else if err.Error() == "forbidden" {
			status = http.StatusForbidden
			msg = "You are not allowed to delete this address"
		}

		c.JSON(status, gin.H{
			"message": msg,
			"detail":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Address successfully deleted"})
}

// Recover (admin only)
func (ctl *AddressController) RecoverAddress(c *gin.Context) {
	id := c.Param("id")
	if err := ctl.service.RecoverAddress(c.Request.Context(), id, uuid.Nil); err != nil {
		status := http.StatusInternalServerError
		msg := "Failed to recover address, try again later"

		if err.Error() == "address not found" {
			status = http.StatusNotFound
			msg = "Address not found"
		}

		c.JSON(status, gin.H{
			"message": msg,
			"detail":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Address successfully recovered"})
}
