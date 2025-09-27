package repository

import (
	"fmt"
	"log"
	"os"

	"github.com/Xillon/golang-todo-api/models"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Database *gorm.DB

func ProvideDatabase() (*gorm.DB, error) {

	dbType := os.Getenv("DB_TYPE")
	if dbType == "" {
		dbType = "mysql"
	}

	var dsn string
	var db *gorm.DB
	var err error

	if dbType == "mysql" {
		username := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASS")
		host := os.Getenv("DB_HOST")
		port := os.Getenv("DB_PORT")
		database := os.Getenv("DB_NAME")

		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", username, password, host, port, database)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	} else {
		fmt.Println("Using SQLite as the default database...")
		dsn = "todo.db"
		db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto-migrate the Todo model
	if err := db.AutoMigrate(&models.Todo{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Database connected and migrated successfully!")
	return db, nil
}
