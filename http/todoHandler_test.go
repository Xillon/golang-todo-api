package http_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/Xillon/golang-todo-api/helpers"
	"github.com/Xillon/golang-todo-api/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddTodos(t *testing.T) {
	router, _ := helpers.SetupRouterWithSQLite(t)

	payload := map[string][]models.Todo{
		"todos": {{
			Title:       "Buy milk",
			Description: "Semi skimmed",
		}},
	}

	jsonPayload, err := json.Marshal(payload)

	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, "/todos", bytes.NewReader(jsonPayload))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
}

func TestUpdateTodos(t *testing.T) {
	router, db := helpers.SetupRouterWithSQLite(t)
	seeded := helpers.SeedTodos(t, db, models.Todo{Title: "Seed update", Description: "Original"})
	require.Len(t, seeded, 1)

	payload := map[string][]models.Todo{
		"todos": {{
			ID:          seeded[0].ID,
			Title:       "Seed update updated",
			Description: "Now full cream",
		}},
	}

	jsonPayload, err := json.Marshal(payload)

	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}

	req, err := http.NewRequest(http.MethodPatch, "/todos", bytes.NewReader(jsonPayload))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var updated models.Todo
	require.NoError(t, db.First(&updated, seeded[0].ID).Error)
	assert.Equal(t, "Seed update updated", updated.Title)
	assert.Equal(t, "Now full cream", updated.Description)
}

func TestGetTodos(t *testing.T) {
	router, db := helpers.SetupRouterWithSQLite(t)
	helpers.SeedTodos(t, db, models.Todo{Title: "Seed list"})

	req, err := http.NewRequest(http.MethodGet, "/todos", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestDeleteTodoByID(t *testing.T) {
	router, db := helpers.SetupRouterWithSQLite(t)
	seeded := helpers.SeedTodos(t, db, models.Todo{Title: "Seed delete"})
	require.Len(t, seeded, 1)

	req, err := http.NewRequest(http.MethodDelete, "/todos/"+strconv.Itoa(int(seeded[0].ID)), nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	var count int64
	db.Model(&models.Todo{}).Count(&count)
	assert.Equal(t, int64(0), count)
}
