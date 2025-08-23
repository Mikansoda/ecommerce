package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"ecommerce/controller"
	"ecommerce/middleware"
	"ecommerce/repository"
	"ecommerce/service"
)

func RegisterAuthRoutes(r *gin.Engine, db *gorm.DB) {
	authRepo := repository.NewUserRepository(db)
	authSvc := service.NewAuthService(authRepo)
	authCtl := controller.NewAuthController(authSvc)

	authApi := r.Group("/auth")
	{
		authApi.POST("/register", authCtl.Register)
		authApi.POST("/verify-otp", authCtl.VerifyOTP)
		authApi.POST("/login", authCtl.Login)
		authApi.POST("/refresh", authCtl.Refresh)
		authApi.POST("/logout", authCtl.Logout)

		authApi.GET("/profile", middleware.Auth("user", "admin"), func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"userID": c.GetString("userID"),
				"email":  c.GetString("email"),
				"role":   c.GetString("role"),
			})
		})

		authApi.GET("/admin/dashboard", middleware.Auth("admin"), func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Welcome to admin dashboard"})
		})
	}
}
