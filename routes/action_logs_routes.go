package routes

import (
	"ecommerce/controller"
	"ecommerce/middleware"
	"ecommerce/repository"
	"ecommerce/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterActionLogRoutes(r *gin.Engine, db *gorm.DB) {
	repo := repository.NewActionLogRepository(db)
	svc := service.NewActionLogService(repo)
	ctl := controller.NewActionLogController(svc)

	admin := r.Group("/admin", middleware.Auth("admin"))
	{
		admin.GET("/logs", ctl.GetLogs)
		admin.GET("/logs/:id", ctl.GetLogByID)

		admin.GET("/reports/selling", ctl.ReportSelling)
		admin.GET("/reports/stock", ctl.ReportStock)
	}
}
