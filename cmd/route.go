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
		name, _ := cmd.Flags().GetString("name")
		if err != nil {
			return fmt.Errorf("failed to get project '%s': %w", projectName, err)
		}
		methodStr = strings.ToUpper(methodStr)
		httpMethod, err := route.ParseHTTPMethod(methodStr)
		if err != nil {
			return fmt.Errorf("invalid HTTP method '%s': %w", methodStr, err)
		}
		if name == "" {
			name = string(httpMethod) + " " + path[1:]
		}
		r := &route.Route{
			ProjectID:   p.ID,
			Name:        name,
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
		projectName, _ := cmd.Flags().GetString("project")
		p, err := storage.GetProject(projectName)
		if err != nil {
			return fmt.Errorf("failed to get project '%s': %w", projectName, err)
		}
		routes, err := storage.ListRoutesByProject(p.ID)
		if err != nil {
			return fmt.Errorf("failed to list routes for project: %w", err)
		}
		if len(routes) == 0 {
			fmt.Printf("No routes found for project '%s'\n", projectName)
			return nil
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		if _, err := fmt.Fprintln(w, "ID\tName\tMethod\tPath"); err != nil {
			return fmt.Errorf("failed to write header: %w", err)
		}
		for _, r := range routes {
			if _, err := fmt.Fprintf(w, "%d\t%s\t%s\t%s\n", r.ID, r.Name, r.Method, r.Path); err != nil {
				return fmt.Errorf("failed to write table line: %w", err)
			}
		}
		if err := w.Flush(); err != nil {
			return fmt.Errorf("failed to write routes table: %w", err)
		}
		return nil
	},
}

var routeUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a route",
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName, _ := cmd.Flags().GetString("project")
		routeName, _ := cmd.Flags().GetString("route")
		methodStr, _ := cmd.Flags().GetString("method")
		path, _ := cmd.Flags().GetString("path")
		description, _ := cmd.Flags().GetString("description")
		newName, _ := cmd.Flags().GetString("rename")
		p, err := storage.GetProject(projectName)
		if err != nil {
			return fmt.Errorf("project not found: %w", err)
		}

		r, err := storage.GetRouteByName(p.ID, routeName)
		if err != nil {
			return fmt.Errorf("route '%s' not found in project '%s': %w", routeName, p.Name, err)
		}
		updates := &route.UpdateRouteInput{}
		if methodStr != "" {
			httpMethod, err := route.ParseHTTPMethod(methodStr)
			if err != nil {
				return fmt.Errorf("invalid HTTP method '%s': %w", methodStr, err)
			}
			updates.Method = &httpMethod
		}
		if path != "" {
			updates.Path = &path
		}
		if description != "" {
			updates.Description = &description
		}
		if newName != "" {
			updates.Name = &newName
		}
		if err := storage.UpdateRoute(r.ID, updates); err != nil {
			return fmt.Errorf("failed to update route '%s': %w", routeName, err)
		}
		fmt.Printf("Project '%s' route %s updated successfully\n", projectName, r.Name)
		return nil
	},
}

var routeDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a route",
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName, _ := cmd.Flags().GetString("project")
		p, err := storage.GetProject(projectName)
		if err != nil {
			return fmt.Errorf("failed to find project '%s': %w", projectName, err)
		}
		routeName, _ := cmd.Flags().GetString("route")
		r, err := storage.GetRouteByName(p.ID, routeName)
		if err != nil {
			return fmt.Errorf("route '%s' not found in project '%s': %w", routeName, p.Name, err)
		}
		if err := storage.DeleteRoute(r.ID); err != nil {
			return fmt.Errorf("failed to delete route from project: %w", err)
		}
		fmt.Printf("Route %s deleted from project '%s'\n", r.Name, p.Name)

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
	routeAddCmd.Flags().StringP("name", "n", "", "Route name (default = Method + Path)")
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
	routeUpdateCmd.Flags().StringP("route", "r", "", "Route Name (required)")
	if err := routeUpdateCmd.MarkFlagRequired("route"); err != nil {
		panic(err)
	}

	routeUpdateCmd.Flags().StringP("method", "m", "", "HTTP method (optional)")
	routeUpdateCmd.Flags().StringP("path", "", "", "Route path (optional)")
	routeUpdateCmd.Flags().StringP("description", "d", "", "Route description (optional)")
	routeUpdateCmd.Flags().StringP("rename", "", "", "Rename route (optional)")

	// Delete command flags
	routeDeleteCmd.Flags().StringP("project", "p", "", "Project name (required)")
	if err := routeDeleteCmd.MarkFlagRequired("project"); err != nil {
		panic(err)
	}
	routeDeleteCmd.Flags().StringP("route", "r", "", "Route Name (required)")
	if err := routeDeleteCmd.MarkFlagRequired("route"); err != nil {
		panic(err)
	}
}
