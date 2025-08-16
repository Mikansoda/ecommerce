package controller

import (
	"net/http"
	"strconv"
	"time"

	"marketplace/entity"
	"marketplace/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AddressController struct {
	service service.AddressService
}

func NewAddressController(s service.AddressService) *AddressController {
	return &AddressController{service: s}
}

// Struct request untuk endpoint create address dan update
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
	limit, _ := strconv.Atoi(c.Query("limit"))
	offset, _ := strconv.Atoi(c.Query("offset"))
	if limit == 0 {
		limit = 10
	}

	addresses, err := ctl.service.GetAddresses(c.Request.Context(), search, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, addresses)
}

// GET address user version
func (ctl *AddressController) GetUserAddress(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("uid"))
	addr, err := ctl.service.GetAddressByUser(c.Request.Context(), userID)
	if err != nil || addr == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "address not found"})
		return
	}
	c.JSON(http.StatusOK, addr)
}

// Create address
func (ctl *AddressController) CreateAddress(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("uid"))
	var req createAddressReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, address)
}

// Update address
func (ctl *AddressController) UpdateAddress(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("uid"))
	var req updateAddressReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existing, _ := ctl.service.GetAddressByUser(c.Request.Context(), userID)
	if existing == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "address not found"})
		return
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, existing)
}

// Delete address
func (ctl *AddressController) DeleteAddress(c *gin.Context) {
	id := c.Param("id")
	userID, _ := uuid.Parse(c.GetString("uid"))

	if err := ctl.service.DeleteAddress(c.Request.Context(), id, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// Recover
func (ctl *AddressController) RecoverAddress(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("uid"))
	id := c.Param("id")
	if err := ctl.service.RecoverAddress(c.Request.Context(), id, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "address recovered"})
}
