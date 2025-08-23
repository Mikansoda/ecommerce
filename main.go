package main

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/gorm"

	"ecommerce/config"
	"ecommerce/entity"
	"ecommerce/repository"
	"ecommerce/routes"
	"ecommerce/service"
)

func CreateTableIfNotExists(db *gorm.DB, model interface{}) error {
	if !db.Migrator().HasTable(model) {
		if err := db.Migrator().CreateTable(model); err != nil {
			return err
		}
	}
	return nil
}

func MigrateTables(db *gorm.DB) error {
	if err := CreateTableIfNotExists(db, &entity.Users{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(
		&entity.Address{},
	); err != nil {
		return err
	}

	if err := db.AutoMigrate(
		&entity.Product{},
		&entity.ProductImage{},
		&entity.ProductCategory{},
	); err != nil {
		return err
	}

	if err := db.AutoMigrate(
		&entity.Payment{},
	); err != nil {
		return err
	}

	if err := db.AutoMigrate(
		&entity.Order{},
		&entity.OrderItem{},
	); err != nil {
		return err
	}

	if err := db.AutoMigrate(
		&entity.ActionLog{},
		&entity.Cart{},
		&entity.CartItem{},
	); err != nil {
		return err
	}
	return nil
}

func MigrateDatabase(db *gorm.DB) {
	err := MigrateTables(db)
	if err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}
	log.Println("Database migration successful")
}

func main() {
	_ = godotenv.Load()
	config.Init()
	db := config.ConnectDatabase()

	MigrateDatabase(db)

	logRepo := repository.NewActionLogRepository(db)
	logSvc := service.NewActionLogService(logRepo)

	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepo(db)
	paymentRepo := repository.NewPaymentRepo(db)

	productSvc := service.NewProductService(productRepo)
	paymentSvc := service.NewPaymentService(paymentRepo, orderRepo, productSvc, logSvc, db)

	// Setup router
	xenditAPIKey := os.Getenv("XENDIT_API_KEY")
	r := routes.SetupRouter(db, xenditAPIKey)

	go func() {
		ticker := time.NewTicker(1 * time.Hour) // check every 1 hour
		for range ticker.C {
			log.Println("Running auto-cancel pending payments...")
			paymentSvc.AutoCancelPendingPayments()
		}
	}()

	log.Println("listening on :" + config.C.AppPort)
	if err := r.Run(":" + config.C.AppPort); err != nil {
		log.Fatal(err)
	}
}
