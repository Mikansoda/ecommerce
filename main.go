package main

import (
	"log"
	"marketplace/config"
	"marketplace/entity"
	"marketplace/routes"
	"github.com/joho/godotenv"
	// "github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

// Note: Understanding db
// by using gorm.Open() → will return *gorm.DB which stores the database connection → stored result in a variable, usually labeled as db → the db variable is used for queries, migrations, etc.

// Note: flow of program
// 1. program starts at func main (entry point) which only uses MigrateDatabase(). So the MigrateDatabase() function is the starting point for the database migration process.
// 2. MigrateDatabase() loads env, connect to the database via config.ConnectDatabase(), the connection result is saved in variable named "db". The function then uses the MigrateTables(db) function to start migrating the database tables.
// MigrateDatabase() is the main "controller" function that performs all migration processes without the need for external input (hence the empty param).
// 3. MigrateTables(db *gorm.DB) This function creates or updates tables in the database. First, it manually checks and creates the User and Admin tables with CreateTableIfNotExists().
// then, AutoMigrate() is called to automatically create/update other tables. If any errors occur during the process, the function immediately returns an error to stop the migration process.
// 4. CreateTableIfNotExists(db *gorm.DB, model interface{}) error This function first checks whether a table for the model (in this project, user and admin) already exists in the database. If not, a new table based on the model structure provided will be created.

// Note: error, conflict of FK in struct order and payment
// both structs depends on respective ID as FK, to resolve this, simply make these as comment temporarily:
// PaymentID    *uuid.UUID `gorm:"type:char(36);null" json:"payment_id,omitempty"`
// Payment   *Payment   `gorm:"foreignKey:PaymentID;references:ID" json:"payment,omitempty"`
// run main.go, db should be created w/o conflict, then un-comment both, finish.

// Function to create a table if one does not exist yet, using interface as it will receive many different struct model (users, admins, etc.)
func CreateTableIfNotExists(db *gorm.DB, model interface{}) error {
	// for table user and admin, if not present in db, create one.
	if !db.Migrator().HasTable(model) {
		// db = object to db connection in GORM, to control and interact w db. db has many methods
		// .migrator is a method belonged to db, which is used for migrations (create table, delete, check, so on).
		// this line is used to call the CreateTable(model) method of the db.Migrator() object. OR the command to create a table in the database based on the model (struct that is provided).
		// errors are saved in err, hence the conditional if.
		if err := db.Migrator().CreateTable(model); err != nil {
			// return error if there's an error
			return err
		}
	}
	// otherwise, return nil (no error)
	return nil
}

// Parameter of function using pointer receiver to modify actual data in db instead of creating a copy and modifying it.
// function return error that occurs so that the function caller is notifoed regarding the problem and resolve or stop the program.

// db is the name of the parameter variable (self-named). *gorm.DB is the data type, a pointer to the gorm.DB struct provided by GORM.
// in short, db is a database connection that you can use to access and manage databases via GORM.
func MigrateTables(db *gorm.DB) error {
	// Create User and Admin table first manually, using function made above, pointer to struct User made in entity.
	if err := CreateTableIfNotExists(db, &entity.UsersAndAdmins{}); err != nil {
		return err
	}

	// AutoMigrate according to foreign key dependencies
	// This line will automigrate list of parameters (arguments) to be migrated, stored inside ().
	if err := db.AutoMigrate(
		// List the struct models/arguments that will be migrated to a db.
		&entity.Address{},
		// if value of err not nil, then return error message
	); err != nil {
		return err
	}

	// The rest of the code till func MigrateDatabase has the same pattern, etc., sans the struct models name

	if err := db.AutoMigrate(
		&entity.Product{},
		&entity.ProductImage{},
		&entity.ProductCategory{},
		&entity.Rating{},
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

	// if value of err is nil, then return nil
	return nil
}

// This creates a function called MigrateDatabase, tasked to manage the database migration process from start to finish.
func MigrateDatabase(db *gorm.DB) {
	err := MigrateTables(db)
	if err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}
	log.Println("Database migration successful")
}

func main() {
	// gin.SetMode(gin.ReleaseMode)
	_ = godotenv.Load()
	config.Init()
	db := config.ConnectDatabase() // dapetin koneksi DB

	MigrateDatabase(db) // modif biar MigrateDatabase nerima parameter db

	r := routes.SetupRouter(db) // passing koneksi langsung, nggak pake config.DB global
	log.Println("listening on :" + config.C.AppPort)
	if err := r.Run(":" + config.C.AppPort); err != nil {
		log.Fatal(err)
	}
}