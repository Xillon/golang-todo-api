package helpers

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Xillon/golang-todo-api/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func SetupSqlMock(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()

	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	gdb, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm with sqlmock: %v", err)
	}

	return gdb, mock
}

func SeedTodos(t *testing.T, db *gorm.DB, todos ...models.Todo) []models.Todo {
	t.Helper()
	for i := range todos {
		if err := db.Create(&todos[i]).Error; err != nil {
			t.Fatalf("failed to seed todo: %v", err)
		}
	}
	return todos
}
