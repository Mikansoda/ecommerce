package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"ecommerce/middleware"
	"ecommerce/repository"
	"ecommerce/service"
	_ "ecommerce/docs"
)

func SetupRouter(db *gorm.DB, xenditAPIKey string) *gin.Engine {
	r := gin.Default()

	// init repo & service for logger
	logRepo := repository.NewActionLogRepository(db)
	logSvc := service.NewActionLogService(logRepo)

	// use logger middleware (global)
	r.Use(gin.Recovery())
	r.Use(middleware.ActionLogger(logSvc))

	r.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	RegisterAuthRoutes(r, db)
	RegisterAddressRoutes(r, db)
	RegisterProductRoutes(r, db)
	RegisterCartRoutes(r, db)
	RegisterOrderRoutes(r, db)
	RegisterPaymentRoutes(r, db)
	RegisterActionLogRoutes(r, db)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
