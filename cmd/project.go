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
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		projects, err := storage.ListProjects()
		if err != nil {
			return fmt.Errorf("failed to list all projects: %w", err)
		}
		if len(projects) == 0 {
			fmt.Println("No projects found")
			return nil
		}

		for _, p := range projects {
			fmt.Printf("%v - %v (%v)\n", p.Name, p.BaseURL, p.DateCreated)
		}

		return nil
	},
}

// TODO: Implement deleteCmd
func init() {
	rootCmd.AddCommand(projectCmd)
	projectCmd.AddCommand(createCmd)
	projectCmd.AddCommand(listCmd)

	createCmd.Flags().StringP("name", "n", "", "Project name (required)")
	if err := createCmd.MarkFlagRequired("name"); err != nil {
		panic(err)
	}
	createCmd.Flags().String("url", "", "Base URL for project api (required)")
	if err := createCmd.MarkFlagRequired("url"); err != nil {
		panic(err)
	}
	createCmd.Flags().StringP("description", "d", "", "Project description (optional)")
	// TODO: Add flags to deleteCmd (name)
}
