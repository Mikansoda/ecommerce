package routes

import (
	"ecommerce/controller"
	"ecommerce/middleware"
	"ecommerce/repository"
	"ecommerce/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterAddressRoutes(r *gin.Engine, db *gorm.DB) {
	addressRepo := repository.NewAddressRepository(db)
	addressSvc := service.NewAddressService(addressRepo)
	addressCtl := controller.NewAddressController(addressSvc)

	// user routes
	addressApi := r.Group("/user", middleware.Auth("user"))
	{
		addressApi.POST("/addresses", addressCtl.CreateAddress)
		addressApi.GET("/addresses", addressCtl.GetUserAddress) // get address by user
		addressApi.PATCH("/addresses", addressCtl.UpdateAddress)
		addressApi.DELETE("/addresses/:id", addressCtl.DeleteAddress)
	}

	// Admin-protected routes
	adminApi := r.Group("/admin", middleware.Auth("admin"))
	{   // GET /admin/addresses?search=Jakarta&limit=5&offset=10
		adminApi.GET("/addresses", addressCtl.GetAddresses) // get all addresses
		adminApi.PATCH("/addresses/:id/recover", addressCtl.RecoverAddress)
	}
}
