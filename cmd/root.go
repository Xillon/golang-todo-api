package cmd

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "golang-todo-api",
	Short: "Manage and operate the Go Todo API service",
	Long: `Use the "api" command to boot the server and
"migrate" to apply schema updates.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	if err := godotenv.Load(); err != nil {
		_ = godotenv.Load(".env.example")
	}
}
