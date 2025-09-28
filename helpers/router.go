package helpers

import (
	"testing"

	"github.com/Xillon/golang-todo-api/http"
	"github.com/Xillon/golang-todo-api/models"
	"github.com/gin-gonic/gin"
	sqlite "github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func SetupRouterWithSQLite(t *testing.T) (*gin.Engine, *gorm.DB) {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite: %v", err)
	}

	if err := db.AutoMigrate(&models.Todo{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	if err := db.Exec("DELETE FROM todos").Error; err != nil {
		t.Fatalf("failed to reset todos table: %v", err)
	}
	handler := http.ProvideTodoHandler(db)
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.POST("/todos", handler.AddTodos)
	router.PATCH("/todos", handler.UpdateTodos)
	router.GET("/todos", handler.GetTodos)
	router.DELETE("/todos/:id", handler.DeleteTodoById)

	return router, db

}
