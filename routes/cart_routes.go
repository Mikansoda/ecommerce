package routes

import (
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"ecommerce/controller"
	"ecommerce/middleware"
	"ecommerce/repository"
	"ecommerce/service"
)

func RegisterCartRoutes(r *gin.Engine, db *gorm.DB) {
	cartRepo := repository.NewCartRepository(db)
	productRepo := repository.NewProductRepository(db)
	cartSvc := service.NewCartService(cartRepo, productRepo)
	cartCtl := controller.NewCartController(cartSvc)

	// User-protected routes
	auth := r.Group("/user", middleware.Auth("user"))
	{
		auth.GET("/cart", cartCtl.GetCart)
		auth.POST("/cart/items", middleware.RateLimit(20, time.Minute), cartCtl.AddItem)
		auth.DELETE("/cart/items/:id", middleware.RateLimit(30, time.Minute), cartCtl.DeleteItem)
	}
}
