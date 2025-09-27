/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// apiCmd represents the api command
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Start the To Do API server",
	Long: `This command starts the To Do API server which allows you to manage your tasks. 
	It contains endpoints for creating, reading, updating, and deleting tasks.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("starting To Do API server...")
	},
}

func init() {
	rootCmd.AddCommand(apiCmd)
}

func startApiServer() {
}
