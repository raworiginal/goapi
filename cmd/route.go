package main

import (
	"fmt"
	"strings"

	"github.com/raworiginal/goapi/internal/route"
	"github.com/raworiginal/goapi/internal/storage"
	"github.com/spf13/cobra"
)

var routeCmd = &cobra.Command{
	Use:   "route",
	Short: "Manage API routes",
}

// TODO: implement addCmd
var routeAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new route to a project",
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName, _ := cmd.Flags().GetString("project")
		methodStr, _ := cmd.Flags().GetString("method")
		path, _ := cmd.Flags().GetString("path")
		description, _ := cmd.Flags().GetString("description")
		p, err := storage.GetProject(projectName)
		if err != nil {
			return fmt.Errorf("failed to get project '%s': %w", projectName, err)
		}
		methodStr = strings.ToUpper(methodStr)
		var httpMethod route.HTTPMethod
		switch methodStr {
		case "GET":
			httpMethod = route.GET
		case "POST":
			httpMethod = route.POST
		case "PUT":
			httpMethod = route.PUT
		case "DELETE":
			httpMethod = route.DELETE
		case "PATCH":
			httpMethod = route.PATCH
		default:
			return fmt.Errorf("invalid HTTP method: %s. Valid methods are: GET, POST, PUT, PATCH, DELETE", methodStr)
		}
		r := &route.Route{
			ProjectID:   p.ID,
			Method:      httpMethod,
			Path:        path,
			Description: description,
		}
		if err := storage.CreateRoute(r); err != nil {
			return fmt.Errorf("failed to create route: %w", err)
		}
		fmt.Printf("Created route: %s %s\n", r.Method, r.Path)
		return nil
	},
}

// TODO: implement listCmd
var routeListCmd = &cobra.Command{
	Use:   "list",
	Short: "List routes for a project",
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: Extract project flag
		// TODO: List routes for that project
		// TODO: Display in table format
		return nil
	},
}

// TODO: implement updateCmd
var routeUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a route",
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: Extract flags (project, id, and optional: method, path, description)
		// TODO: Validate route exists
		// TODO: Create UpdateRouteInput with non-nil values
		// TODO: Call storage.UpdateRoute()
		// TODO: Print success message
		return nil
	},
}

// TODO: implement deleteCmd
var routeDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a route",
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: Extract flags (project, id)
		// TODO: Validate route exists
		// TODO: Delete route
		// TODO: Print success message
		return nil
	},
}

func init() {
	// Add commands
	rootCmd.AddCommand(routeCmd)
	routeCmd.AddCommand(routeAddCmd)
	routeCmd.AddCommand(routeListCmd)
	routeCmd.AddCommand(routeUpdateCmd)
	routeCmd.AddCommand(routeDeleteCmd)
	// registerFlags
	routeAddCmd.Flags().StringP("project", "p", "", "Project name (required)")
	if err := routeAddCmd.MarkFlagRequired("project"); err != nil {
		panic(err)
	}
	routeAddCmd.Flags().StringP("method", "m", "", "HTTP method: GET, POST, PUT PATCH, DELETE (required)")
	if err := routeAddCmd.MarkFlagRequired("method"); err != nil {
		panic(err)
	}
	routeAddCmd.Flags().StringP("path", "", "", "Route path (required)")
	if err := routeAddCmd.MarkFlagRequired("path"); err != nil {
		panic(err)
	}
	routeAddCmd.Flags().StringP("description", "d", "", "Route description (optional)")
}
