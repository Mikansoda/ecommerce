package routes

import (
	"marketplace/controller"
	"marketplace/middleware"
	"marketplace/repository"
	"marketplace/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterAddressRoutes(r *gin.Engine, db *gorm.DB) {
	addressRepo := repository.NewAddressRepository(db)
	addressSvc := service.NewAddressService(addressRepo)
	addressCtl := controller.NewAddressController(addressSvc)

	// user routes
	addressApi := r.Group("/addresses", middleware.Auth("user", "admin"))
	{
		addressApi.POST("/", addressCtl.CreateAddress)
		addressApi.GET("/", addressCtl.GetUserAddress)  // get address user sendiri
		addressApi.PATCH("/", addressCtl.UpdateAddress)          // update address user sendiri
		addressApi.DELETE("/:id", addressCtl.DeleteAddress)         // delete address user sendiri
		addressApi.POST("/recover/:id", addressCtl.RecoverAddress)
	}

	// admin routes
	adminApi := r.Group("/admin", middleware.Auth("admin"))
	{
		adminApi.GET("/addresses", addressCtl.GetAddresses) // get semua address
	}
}
