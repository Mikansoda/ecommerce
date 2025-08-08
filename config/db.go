package config

import (
	"fmt"
	"os" // to read env file/operating system
	"gorm.io/driver/mysql" // for GORM yo connect to MySQL db
	"gorm.io/gorm" // GORM
)

// ConnectDatabase initializes and returns a connection to database
// This function returns a *gorm.DB pointer object, which is a connection to the database that can be used for DB operations.
func ConnectDatabase() *gorm.DB {
	// DB_HOST (database server address)
	// DB_PORT (database port)
	// DB_USER (database username)
	// DB_PASSWORD (database password)
	// DB_NAME (database name)
	// These values are stored in the dbHost, dbPort, dbUser, dbPass, and dbName variables, respectively.
	// os.Getenv("ENV_NAME") = used to read environment variable values from the OS or .env file.
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Construct Data Source Name for MySQL connection
	// dsn/Data Source Name, which is a string containing complete information for connecting to a database.
	// using sprintf to format string by replacing %s with the given variable.
	//dsn result will be = "user:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local". a standard format for MySQL connections using GORM.
  	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        dbUser, dbPass, dbHost, dbPort, dbName)
    
	// Open a connection to DB using GROM + MySQL driver, return error message if failed 
	// gorm.Open = GORM function to open database connection. Ready to connect and interact with db. (ex, phone analogy: dialed and ready to "talk")
	// mysql.Open(dsn) = using MySQL driver to input address and credibilities of db. (ex, phone analogy: input phone number)
    // &gorm.Config{} = providing GORM default rules (blank, meaning use standard rules) for database connections. (ex of rules: naming convention for tables? timeout period?)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("Failed to connect to database: " + err.Error())
    }
	// Return a message and an instance to be used in main.go (db)
    fmt.Println("Database connected successfully")
    return db
}
