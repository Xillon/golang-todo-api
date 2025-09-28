package http

import (
	"fmt"
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

// AddTodos godoc
// @Summary      Add a list of todos
// @Description  Creates one or more todos
// @Tags         todos
// @Accept       json
// @Produce      json
// @Param        X-API-Key  header  string  true  "API key"
// @Param        request    body    map[string][]models.Todo  true  "Todos payload"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Router       /todos [post]
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

// UpdateTodos godoc
// @Summary      Update a list of todos
// @Description  Updates one or more todos by id
// @Tags         todos
// @Accept       json
// @Produce      json
// @Param        X-API-Key  header  string  true  "API key"
// @Param        request    body    map[string][]models.Todo  true  "Todos payload"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Router       /todos [patch]
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

// GetTodos godoc
// @Summary      List todos
// @Description  Returns a paginated list of todos
// @Tags         todos
// @Produce      json
// @Param        X-API-Key  header  string  true  "API key"
// @Param        page       query   int     false "Page number"  default(1)
// @Param        limit      query   int     false "Items per page"  default(10)
// @Success      200  {object}  map[string]interface{}
// @Router       /todos [get]
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

// DeleteTodoById godoc
// @Summary      Delete todo by ID
// @Tags         todos
// @Param        X-API-Key  header  string  true  "API key"
// @Param        id         path    int     true "Todo ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Router       /todos/{id} [delete]
func (h *TodoHandler) DeleteTodoById(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Id is required"})
		return
	}

	if err := h.DB.Delete(&models.Todo{}, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Todo with id %s deleted successfully", id)})
}
