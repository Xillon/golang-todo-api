package persistance

import (
	"log"
	"os"

	"github.com/Xillon/golang-todo-api/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Database *gorm.DB

func InitDatabase() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN environment variable is not set")
	}

	var err error
	Database, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connection established")

	if err := Database.AutoMigrate(&models.Todo{}); err != nil {

		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Database migrated successfully")
}
