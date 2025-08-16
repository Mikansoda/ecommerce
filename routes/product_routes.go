package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"marketplace/controller"
	"marketplace/middleware"
	"marketplace/repository"
	"marketplace/service"
)

func RegisterProductRoutes(r *gin.Engine, db *gorm.DB) {
	// Products
productRepo := repository.NewProductRepository(db)
productSvc := service.NewProductService(productRepo)
productCtl := controller.NewProductController(productSvc)

// Categories
categoryRepo := repository.NewCategoryRepository(db)
categorySvc := service.NewCategoryService(categoryRepo)
categoryCtl := controller.NewCategoryController(categorySvc)

r.GET("/products", productCtl.GetProducts)
r.GET("/products/:id", productCtl.GetProductByID)
r.GET("/categories", categoryCtl.ListCategories)

// Admin-protected routes
admin := r.Group("/admin", middleware.Auth("admin"))
{
    admin.POST("/categories", categoryCtl.CreateCategory)
    admin.PATCH("/categories/:id", categoryCtl.UpdateCategory)
    admin.DELETE("/categories/:id", categoryCtl.DeleteCategory)
    admin.POST("/categories/:id/recover", categoryCtl.RecoverCategory)

    admin.POST("/products", productCtl.CreateProduct)
    admin.PATCH("/products/:id", productCtl.UpdateProduct)
    admin.DELETE("/products/:id", productCtl.DeleteProduct)
    admin.POST("/products/:id/recover", productCtl.RecoverProduct)
    // otw route images
    }
}
