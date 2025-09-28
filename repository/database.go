package repository

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/Xillon/golang-todo-api/models"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	_ "github.com/go-sql-driver/mysql"
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
		if host == "" {
			host = "127.0.0.1"
		}
		port := os.Getenv("DB_PORT")
		if port == "" {
			port = "3306"
		}
		database := os.Getenv("DB_NAME")

		serverDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/?parseTime=true&multiStatements=true&charset=utf8mb4&collation=utf8mb4_unicode_ci", username, password, host, port)
		sqlDB, errOpen := sql.Open("mysql", serverDSN)

		if errOpen != nil {
			return nil, fmt.Errorf("failed to open mysql server connection: %w", errOpen)
		}

		defer sqlDB.Close()

		if _, errExec := sqlDB.Exec("CREATE DATABASE IF NOT EXISTS `" + database + "` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci"); errExec != nil {
			return nil, fmt.Errorf("failed to create database %s: %w", database, errExec)
		}

		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local&charset=utf8mb4&collation=utf8mb4_unicode_ci", username, password, host, port, database)
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
