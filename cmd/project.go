package main

import (
	"fmt"
	"net/url"
	"os"
	"text/tabwriter"

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

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		if _, err := fmt.Fprintln(w, "Name\tBase URL\tDescription\tDate Created"); err != nil {
			return fmt.Errorf("failed to write header: %w", err)
		}
		for _, p := range projects {
			if _, err := fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", p.Name, p.BaseURL, p.Description, p.DateCreated.Format("2006-01-02 15:04")); err != nil {
				return fmt.Errorf("failed to write table line: %w", err)
			}
		}
		if err := w.Flush(); err != nil {
			return fmt.Errorf("failed to write projects table: %w", err)
		}

		return nil
	},
}

// TODO: Implement deleteCmd
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a project",
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		if err := storage.DeleteProject(name); err != nil {
			return fmt.Errorf("failed to delete project: %w", err)
		}
		fmt.Printf("Deleted Project %s \n", name)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(projectCmd)
	projectCmd.AddCommand(createCmd)
	projectCmd.AddCommand(listCmd)
	projectCmd.AddCommand(deleteCmd)
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
	deleteCmd.Flags().StringP("name", "n", "", "Project name (required)")
	if err := deleteCmd.MarkFlagRequired("name"); err != nil {
		panic(err)
	}
}
