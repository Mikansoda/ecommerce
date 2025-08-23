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

func RegisterOrderRoutes(r *gin.Engine, db *gorm.DB) {
	// Dependency injection
	cartRepo := repository.NewCartRepository(db)
	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepo(db)

	orderSvc := service.NewOrderService(cartRepo, productRepo, orderRepo)
	orderCtl := controller.NewOrderController(orderSvc)

	// User-protected routes
	user := r.Group("/user", middleware.Auth("user"))
	{
		user.POST("/orders", middleware.RateLimit(10, time.Minute), orderCtl.CreateOrder)
		user.GET("/orders", middleware.RateLimit(10, time.Minute), orderCtl.GetUserOrders)
	}

	// Admin-protected routes
	admin := r.Group("/admin", middleware.Auth("admin"))
	{   // GET /admin/orders?status=pending&limit=5&offset=10
		admin.GET("/orders", orderCtl.GetOrders)
		admin.PUT("/orders/:id/status", orderCtl.UpdateOrderStatus)
	}
}
