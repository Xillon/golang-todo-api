package cmd

import (
	"fmt"

	"github.com/Xillon/golang-todo-api/handlers"
	"github.com/Xillon/golang-todo-api/persistance"
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

		fx.Provide(
			persistance.ProvideDatabase,
			handlers.ProvideTodoHandler,
		),

		fx.Invoke(func(handler *handlers.TodoHandler) {
			r := gin.Default()

			r.POST("/todos", handler.AddTodos)
			r.PATCH("/todos", handler.UpdateTodos)
			r.GET("/todos", handler.GetTodos)

			fmt.Println("API server is running on http://localhost:8080")
			if err := r.Run(":8080"); err != nil {
				fmt.Printf("Failed to run server: %v\n", err)
			}
		}),
	)

	app.Run()
}
