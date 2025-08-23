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

func RegisterPaymentRoutes(r *gin.Engine, db *gorm.DB) {
	// Dependency injection
	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepo(db)
	paymentRepo := repository.NewPaymentRepo(db)
	logRepo := repository.NewActionLogRepository(db)

	productSvc := service.NewProductService(productRepo)
	logSvc := service.NewActionLogService(logRepo)
	paymentSvc := service.NewPaymentService(paymentRepo, orderRepo, productSvc, logSvc, db) // <-- tambahin db
	paymentCtl := controller.NewPaymentController(paymentSvc)

	// User-protected routes
	user := r.Group("/user", middleware.Auth("user"))
	{
		user.POST("/payments/xendit", middleware.RateLimit(10, time.Minute), paymentCtl.CreatePayment)
		user.GET("/payments", middleware.RateLimit(10, time.Minute), paymentCtl.GetUserPayments)
	}

	// Admin-protected routes
	admin := r.Group("/admin", middleware.Auth("admin"))
	{
		admin.GET("/payments", paymentCtl.GetPayments)
		admin.POST("/payments/webhook/xendit", paymentCtl.XenditWebhook)
	}
}
