package main

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"ecommerce/config"
	"ecommerce/entity"
	"ecommerce/repository"
	"ecommerce/routes"
	"ecommerce/service"
)

func SeedAdmin(db *gorm.DB) {
	log.Println("Running admin seeding...")
	email := os.Getenv("ADMIN_EMAIL")
	username := os.Getenv("ADMIN_USERNAME")
	password := os.Getenv("ADMIN_PASSWORD")

	if email == "" || username == "" || password == "" {
		log.Println("ADMIN_EMAIL, ADMIN_USERNAME, or ADMIN_PASSWORD not set. Skipping admin seed.")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Failed to hash password:", err)
		return
	}

	newAdmin := entity.Users{
		ID:           uuid.New(),
		Email:        email,
		Username:     username,
		PasswordHash: string(hashedPassword),
		Role:         entity.RoleAdmin,
		IsActive:     true,
	}

	if err := db.Create(&newAdmin).Error; err != nil {
		log.Println("Failed to create admin user:", err)
	} else {
		log.Println("Admin user created successfully:", email)
	}
}

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
	// load env
	_ = godotenv.Load()
	config.Init()
	// connect ke database, hasil connection object (*gorm.DB) namanya db
	db := config.ConnectDatabase()
    if err := CreateTableIfNotExists(db, &entity.Users{}); err != nil {
	log.Fatal("Failed to create users table:", err)
    }
	SeedAdmin(db)
    // Migrate db, inject connection sbg context db target
	MigrateDatabase(db)
    
	// inject log repo dgn db target, query ke db
	logRepo := repository.NewActionLogRepository(db)
	// service logic bisnis, minta data ke logRepo yg query ke db
	logSvc := service.NewActionLogService(logRepo)
    
	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepo(db)
	paymentRepo := repository.NewPaymentRepo(db)

	productSvc := service.NewProductService(productRepo)

	// paymentSvc inject banyak repo, order repo: buat update status order, product repo: related dgn stok produk
	// logRepo: Business logging (non method) tdk bisa ditangkap middleware, mis user X melakukan pembayaran Order Y, status berubah jadi PAID.
	paymentSvc := service.NewPaymentService(paymentRepo, orderRepo, productSvc, logSvc, db)

	// Ambil Xendit API Key dari .env.
	xenditAPIKey := os.Getenv("XENDIT_API_KEY")
	// Panggil routes.SetupRouter → di dalamnya daftar endpoint
	// Route diarahkan ke controller, controller panggil service, service panggil repository.
	r := routes.SetupRouter(db, xenditAPIKey)
    
	// Background Job → Auto cancel pending payment
	go func() {
		ticker := time.NewTicker(1 * time.Hour) // check every 1 hour
		for range ticker.C {
			log.Println("Running auto-cancel pending payments...")
			paymentSvc.AutoCancelPendingPayments()
		}
	}()
    // Start server (nyalain Gin HTTP)
	log.Println("listening on :" + config.C.AppPort)
	if err := r.Run(":" + config.C.AppPort); err != nil {
		log.Fatal(err)
	}
}
