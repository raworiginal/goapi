package main

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/raworiginal/goapi/internal/route"
	"github.com/raworiginal/goapi/internal/storage"
	"github.com/spf13/cobra"
)

var routeCmd = &cobra.Command{
	Use:   "route",
	Short: "Manage API routes",
}

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
		httpMethod, err := route.ParseHTTPMethod(methodStr)
		if err != nil {
			return err
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

var routeListCmd = &cobra.Command{
	Use:   "list",
	Short: "List routes for a project",
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: Extract project flag
		projectName, _ := cmd.Flags().GetString("project")
		p, err := storage.GetProject(projectName)
		if err != nil {
			return fmt.Errorf("failed to get project '%s': %w", projectName, err)
		}
		// TODO: List routes for that project
		routes, err := storage.ListRoutesByProject(p.ID)
		if err != nil {
			return fmt.Errorf("failed to list routes for project: %w", err)
		}
		// TODO: Display in table format
		if len(routes) == 0 {
			fmt.Printf("No routes found for project '%s'\n", projectName)
			return nil
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		if _, err := fmt.Fprintln(w, "ID\tMethod\tPath"); err != nil {
			return fmt.Errorf("failed to write header: %w", err)
		}
		for _, r := range routes {
			if _, err := fmt.Fprintf(w, "%d\t%s\t%s\n", r.ID, r.Method, r.Path); err != nil {
				return fmt.Errorf("failed to write table line: %w", err)
			}
		}
		if err := w.Flush(); err != nil {
			return err
		}
		return nil
	},
}

// TODO: implement updateCmd
var routeUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a route",
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: Extract flags (project, id, and optional: method, path, description)
		projectName, _ := cmd.Flags().GetString("project")
		id, _ := cmd.Flags().GetUint("id")
		methodStr, _ := cmd.Flags().GetString("method")
		path, _ := cmd.Flags().GetString("path")
		description, _ := cmd.Flags().GetString("description")
		// TODO: Validate route exists
		_, err := storage.GetProject(projectName)
		if err != nil {
			return fmt.Errorf("project not found: %w", err)
		}
		// TODO: Create UpdateRouteInput with non-nil values
		updates := &route.UpdateRouteInput{}
		if methodStr != "" {
			httpMethod, err := route.ParseHTTPMethod(methodStr)
			if err != nil {
				return err
			}
			updates.Method = &httpMethod
		}
		if path != "" {
			updates.Path = &path
		}
		if description != "" {
			updates.Description = &description
		}
		// TODO: Call storage.UpdateRoute()
		if err := storage.UpdateRoute(id, updates); err != nil {
			return err
		}
		// TODO: Print success message
		fmt.Printf("Project '%s' route id %d updated successfully\n", projectName, id)
		return nil
	},
}

var routeDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a route",
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: Extract flags (project, id)
		projectName, _ := cmd.Flags().GetString("project")
		_, err := storage.GetProject(projectName)
		if err != nil {
			return fmt.Errorf("failed to find project '%s': %w", projectName, err)
		}
		id, _ := cmd.Flags().GetUint("id")
		if err := storage.DeleteRoute(id); err != nil {
			return fmt.Errorf("failed to delete route from project: %w", err)
		}
		fmt.Printf("Route id %d deleted from project '%s'\n", id, projectName)

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
	// Add command flags
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

	// List command flags
	routeListCmd.Flags().StringP("project", "p", "", "Project name (required)")
	if err := routeListCmd.MarkFlagRequired("project"); err != nil {
		panic(err)
	}

	// Update command flags
	routeUpdateCmd.Flags().StringP("project", "p", "", "Project name (required)")
	if err := routeUpdateCmd.MarkFlagRequired("project"); err != nil {
		panic(err)
	}
	routeUpdateCmd.Flags().UintP("id", "i", 0, "Route ID (required)")
	if err := routeUpdateCmd.MarkFlagRequired("id"); err != nil {
		panic(err)
	}
	routeUpdateCmd.Flags().StringP("method", "m", "", "HTTP method (optional)")
	routeUpdateCmd.Flags().StringP("path", "", "", "Route path (optional)")
	routeUpdateCmd.Flags().StringP("description", "d", "", "Route description (optional)")

	// Delete command flags
	routeDeleteCmd.Flags().StringP("project", "p", "", "Project name (required)")
	if err := routeDeleteCmd.MarkFlagRequired("project"); err != nil {
		panic(err)
	}
	routeDeleteCmd.Flags().UintP("id", "i", 0, "Route ID (required)")
	if err := routeDeleteCmd.MarkFlagRequired("id"); err != nil {
		panic(err)
	}
}
