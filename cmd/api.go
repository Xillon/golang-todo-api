package cmd

import (
	"fmt"
	"log"
	"os"

	docs "github.com/Xillon/golang-todo-api/docs"
	"github.com/Xillon/golang-todo-api/http"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/fx"
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Start the To Do API server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting To Do API server...")
		godotenv.Load()
		startApiServer()
	},
}

func init() {
	rootCmd.AddCommand(apiCmd)
}

func startApiServer() {
	app := fx.New(
		FxModules,
		fx.Invoke(func(handler *http.TodoHandler) {
			r := gin.Default()

			docs.SwaggerInfo.Title = "Go Todo API"
			docs.SwaggerInfo.Description = "Batch create, update, and list todos. Protected via X-API-Key header when configured."
			docs.SwaggerInfo.Version = "1.0"
			docs.SwaggerInfo.BasePath = "/"

			r.GET("/", func(c *gin.Context) { c.Status(200) })
			r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

			apiKey := os.Getenv("API_KEY")
			if apiKey == "" {
				log.Println("warning: API_KEY not set; requests will not be authenticated")
			}

			secured := r.Group("/")
			if apiKey != "" {
				secured.Use(http.APIKeyMiddleware(apiKey))
			}

			secured.POST("/todos", handler.AddTodos)
			secured.PATCH("/todos", handler.UpdateTodos)
			secured.GET("/todos", handler.GetTodos)
			secured.DELETE("/todos/:id", handler.DeleteTodoById)

			fmt.Println("API server is running on http://localhost:8080/swagger/index.html")
			if err := r.Run(":8080"); err != nil {
				fmt.Printf("Failed to run server: %v\n", err)
			}
		}),
	)

	app.Run()
}
