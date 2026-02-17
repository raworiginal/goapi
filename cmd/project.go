package main

import (
	"fmt"
	"net/url"

	"github.com/raworiginal/goapi/internal/project"
	"github.com/raworiginal/goapi/internal/storage"
	"github.com/spf13/cobra"
)

var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Manage projects",
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new project",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get flags (name, url, description)
		name, _ := cmd.Flags().GetString("name")
		baseURL, _ := cmd.Flags().GetString("url")
		description, _ := cmd.Flags().GetString("description")
		// Create Project struct
		p := &project.Project{
			Name:        name,
			BaseURL:     baseURL,
			Description: description,
		}
		// Call storage.CreateProject()

		// validate URL first
		if _, err := url.Parse(baseURL); err != nil {
			return fmt.Errorf("invalid URL: %s", baseURL)
		}
		if err := storage.CreateProject(p); err != nil {
			return fmt.Errorf("failed to create project: %w", err)
		}

		// Print success message
		fmt.Printf("Project '%s' created\n", name)
		return nil
	},
}

// TODO: Implement listCmd
// TODO: Implement deleteCmd
func init() {
	rootCmd.AddCommand(projectCmd)
	projectCmd.AddCommand(createCmd)
	// TODO: Add listCmd and deleteCmd to projectCmd

	// TODO: Add flags to createCmd (name, url, description)
	// TODO: Add flags to deleteCmd (name)
}
