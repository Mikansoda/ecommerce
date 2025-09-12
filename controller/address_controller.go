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

// GetAddresses godoc
// @Summary      Get all addresses
// @Description  Return list of all user addresses (admin only)
// @Tags         Addresses
// @Security     BearerAuth
// @Produce      json
// @Param        search   query     string  false  "Search by receiver name, city, or province"
// @Param        limit    query     int     false  "Limit number of results"   default(10)
// @Param        offset   query     int     false  "Offset for pagination"     default(0)
// @Success      200      {array}   entity.Address
// @Failure      500      {object}  map[string]interface{}
// @Example 200 {json} Success Example:
// [
//   {
//     "id": "1a2b3c4d-5678-90ab-cdef-1234567890ab",
//     "receiver_name": "John Doe",
//     "phone_number": "08123456789",
//     "address_line": "Jl. Kebon Jeruk No. 12",
//     "city": "Jakarta",
//     "province": "DKI Jakarta",
//     "postal_code": "11530",
//     "created_at": "2025-09-12T21:00:00Z",
//     "updated_at": "2025-09-12T21:00:00Z"
//   },
//   {
//     "id": "2b3c4d5e-6789-01bc-def0-2345678901bc",
//     "receiver_name": "Jane Smith",
//     "phone_number": "08234567890",
//     "address_line": "Jl. Sudirman No. 45",
//     "city": "Bandung",
//     "province": "Jawa Barat",
//     "postal_code": "40123",
//     "created_at": "2025-09-10T09:30:00Z",
//     "updated_at": "2025-09-10T09:30:00Z"
//   }
// ]
// @Example 500 {json} Error Example:
// {
//   "message": "Failed to fetch addresses, try again later",
//   "detail": "some error message"
// }
// @Router       /admin/addresses [get]
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

// GetUserAddress godoc
// @Summary      Get own address
// @Description  Return the logged-in user's address
// @Tags         Addresses
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  entity.Address
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Example 200 {json} Success Example:
// {
//   "id": "1a2b3c4d-5678-90ab-cdef-1234567890ab",
//   "receiver_name": "John Doe",
//   "phone_number": "08123456789",
//   "address_line": "Jl. Kebon Jeruk No. 12",
//   "city": "Jakarta",
//   "province": "DKI Jakarta",
//   "postal_code": "11530",
//   "created_at": "2025-09-12T21:00:00Z",
//   "updated_at": "2025-09-12T21:00:00Z"
// }
// @Example 400 {json} Error Example:
// {
//   "message": "Invalid user ID",
//   "detail": "some error message"
// }
// @Example 404 {json} Error Example:
// {
//   "message": "No address found"
// }
// @Example 500 {json} Error Example:
// {
//   "message": "Failed to fetch your address, try again later",
//   "detail": "some error message"
// }
// @Router       /user/addresses [get]
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

// CreateAddress godoc
// @Summary      Create a new address
// @Description  Create address for logged-in user (user can only have one)
// @Tags         Addresses
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body      createAddressReq  true  "Address request body"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Example 200 {json} Success Example:
// {
//   "message": "Address successfully created",
//   "data": {
//     "id": "1a2b3c4d-5678-90ab-cdef-1234567890ab",
//     "receiver_name": "John Doe",
//     "phone_number": "08123456789",
//     "address_line": "Jl. Kebon Jeruk No. 12",
//     "city": "Jakarta",
//     "province": "DKI Jakarta",
//     "postal_code": "11530",
//     "created_at": "2025-09-12T21:00:00Z",
//     "updated_at": "2025-09-12T21:00:00Z"
//   }
// }
// @Example 400 {json} Error Example:
// {
//   "message": "You already have an address",
//   "detail": "user already has an address"
// }
// @Example 500 {json} Error Example:
// {
//   "message": "Failed to create address, try again later",
//   "detail": "some error message"
// }
// @Router       /user/addresses [post]
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

// UpdateAddress godoc
// @Summary      Update address
// @Description  Update logged-in user's address
// @Tags         Addresses
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body      updateAddressReq  true  "Address update request body"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Failure      404      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Example 200 {json} Success Example:
// {
//   "message": "Address successfully updated",
//   "data": {
//     "id": "1a2b3c4d-5678-90ab-cdef-1234567890ab",
//     "receiver_name": "John Doe",
//     "phone_number": "08123456789",
//     "address_line": "Jl. Kebon Jeruk No. 12",
//     "city": "Jakarta",
//     "province": "DKI Jakarta",
//     "postal_code": "11530",
//     "created_at": "2025-09-12T21:00:00Z",
//     "updated_at": "2025-09-12T21:10:00Z"
//   }
// }
// @Example 400 {json} Error Example:
// {
//   "message": "Invalid input data",
//   "detail": "some error message"
// }
// @Example 404 {json} Error Example:
// {
//   "message": "Address not found"
// }
// @Example 500 {json} Error Example:
// {
//   "message": "Failed to update address, try again later",
//   "detail": "some error message"
// }
// @Router       /user/addresses [patch]
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

// DeleteAddress godoc
// @Summary      Delete address
// @Description  Delete logged-in user's address by ID
// @Tags         Addresses
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Address ID (UUID)"
// @Success      200  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Example 200 {json} Success Example:
// {
//   "message": "Address successfully deleted"
// }
// @Example 403 {json} Error Example:
// {
//   "message": "You are not allowed to delete this address",
//   "detail": "forbidden"
// }
// @Example 404 {json} Error Example:
// {
//   "message": "Address not found"
// }
// @Example 500 {json} Error Example:
// {
//   "message": "Failed to delete address, try again later",
//   "detail": "some error message"
// }
// @Router       /user/addresses/{id} [delete]
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

// RecoverAddress godoc
// @Summary      Recover deleted address
// @Description  Recover a soft-deleted address (admin only)
// @Tags         Addresses
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Address ID (UUID)"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Example 200 {json} Success Example:
// {
//   "message": "Address successfully recovered"
// }
// @Example 404 {json} Error Example:
// {
//   "message": "Address not found"
// }
// @Example 500 {json} Error Example:
// {
//   "message": "Failed to recover address, try again later",
//   "detail": "some error message"
// }
// @Router       /admin/addresses/{id}/recover [post]
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
