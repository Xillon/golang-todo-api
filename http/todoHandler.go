package http

import (
	"net/http"
	"strconv"

	"github.com/Xillon/golang-todo-api/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TodoHandler struct {
	DB *gorm.DB
}

func ProvideTodoHandler(db *gorm.DB) *TodoHandler {
	return &TodoHandler{DB: db}
}

func (h *TodoHandler) AddTodos(c *gin.Context) {
	var request struct {
		Todos []models.Todo `json:"todos"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i := range request.Todos {
		if err := h.DB.Create(&request.Todos[i]).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"todos": request.Todos})
}

func (h *TodoHandler) UpdateTodos(c *gin.Context) {
	var request struct {
		Todos []models.Todo `json:"todos"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, todo := range request.Todos {
		if err := h.DB.Model(&models.Todo{}).Where("id = ?", todo.ID).Updates(todo).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"todos": request.Todos})
}

func (h *TodoHandler) GetTodos(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var todos []models.Todo
	var total int64

	h.DB.Model(&models.Todo{}).Count(&total)
	h.DB.Limit(limit).Offset(offset).Find(&todos)

	c.JSON(http.StatusOK, gin.H{
		"todos": todos,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}
