package main

import (
	"fmt"
	"os"

	"github.com/raworiginal/goapi/internal/storage"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "goapi",
	Short: "A Go-based API testing CLI/TUI",
	Long:  `goapi is a command-line and TUI tool for testing and managing API routes.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return storage.InitDB()
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
