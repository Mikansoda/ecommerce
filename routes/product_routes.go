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

func RegisterProductRoutes(r *gin.Engine, db *gorm.DB) {
	// Dependency injection
	// Products
	productRepo := repository.NewProductRepository(db)
	productSvc := service.NewProductService(productRepo)
	productCtl := controller.NewProductController(productSvc)

	// Categories
	categoryRepo := repository.NewCategoryRepository(db)
	categorySvc := service.NewCategoryService(categoryRepo)
	categoryCtl := controller.NewCategoryController(categorySvc)

	// Image
	imageRepo := repository.NewProductImageRepository(db)
	imageSvc := service.NewProductImageService(imageRepo)
	imageCtl := controller.NewProductImageController(imageSvc)

	// Publik
	// GET /products?search=Roasted%20Almond&category=Consummables&limit=5&offset=10
	r.GET("/products", productCtl.GetProducts)
	r.GET("/products/:productId", productCtl.GetProductByID)
	// GET /categories?limit=5&offset=10
	r.GET("/categories", categoryCtl.GetCategories)

	// Admin-protected routes
	admin := r.Group("/admin", middleware.Auth("admin"))
	{
		admin.POST("/categories", categoryCtl.CreateCategory)
		admin.PUT("/categories/:id", categoryCtl.UpdateCategory)
		admin.DELETE("/categories/:id", categoryCtl.DeleteCategory)
		admin.PATCH("/categories/:id/recover", categoryCtl.RecoverCategory)

		admin.POST("/products", productCtl.CreateProduct)
		admin.PUT("/products/:productId", productCtl.UpdateProduct)
		admin.DELETE("/products/:productId", productCtl.DeleteProduct)
		admin.PATCH("/products/:productId/recover", productCtl.RecoverProduct)
        // 5 req/menit
		admin.POST("/products/:productId/images", middleware.RateLimit(5, time.Minute), imageCtl.UploadImage)
		admin.DELETE("/images/:imageId", middleware.RateLimit(10, time.Minute), imageCtl.DeleteImage)
		admin.POST("/images/:imageId/recover", middleware.RateLimit(10, time.Minute), imageCtl.RecoverImage)
	}
}
