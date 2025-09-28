package cmd

import (
	"fmt"

	"github.com/Xillon/golang-todo-api/http"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Start the To Do API server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting To Do API server...")
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

			r.POST("/todos", handler.AddTodos)
			r.PATCH("/todos", handler.UpdateTodos)
			r.GET("/todos", handler.GetTodos)
			r.DELETE("/todos/:id", handler.DeleteTodoById)

			fmt.Println("API server is running on http://localhost:8080")
			if err := r.Run(":8080"); err != nil {
				fmt.Printf("Failed to run server: %v\n", err)
			}
		}),
	)

	app.Run()
}
