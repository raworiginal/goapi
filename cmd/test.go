package main

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/raworiginal/goapi/internal/api"
	"github.com/raworiginal/goapi/internal/route"
	"github.com/raworiginal/goapi/internal/storage"
	"github.com/spf13/cobra"
)

type TestResult struct {
	RouteID    uint
	RouteName  string
	Method     string
	Path       string
	StatusCode int
	Duration   time.Duration
	Error      string
}

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Test API routes",
	Long:  "Execute routes and display results",
	RunE:  testRun,
}

func init() {
	rootCmd.AddCommand(testCmd)

	// Add your flags here
	testCmd.Flags().StringP("project", "p", "", "Project name (required)")
	testCmd.Flags().String("route", "", "Route name (optional, tests all if omitted)")
	testCmd.Flags().Duration("timeout", 5*time.Second, "Request timeout")

	if err := testCmd.MarkFlagRequired("project"); err != nil {
		panic(err)
	}
}

func testRun(cmd *cobra.Command, args []string) error {
	projectName, _ := cmd.Flags().GetString("project")
	routeName, _ := cmd.Flags().GetString("route")
	timeout, _ := cmd.Flags().GetDuration("timeout")

	p, err := storage.GetProject(projectName)
	if err != nil {
		return fmt.Errorf("failed to load project '%s': %w", projectName, err)
	}
	var routes []*route.Route
	if routeName != "" {
		r, err := storage.GetRouteByName(p.ID, routeName)
		if err != nil {
			return fmt.Errorf("failed to load route '%s': %w", routeName, err)
		}
		routes = []*route.Route{r}
	} else {
		routes, err = storage.ListRoutesByProject(p.ID)
		if err != nil {
			return fmt.Errorf("failed to load routes: %w", err)
		}
	}
	config := api.Config{Timeout: timeout}
	client := api.NewHTTPClient(config)

	var results []TestResult
	for _, r := range routes {
		url := p.BaseURL + r.Path

		resp, err := client.Do(string(r.Method), url, nil)

		result := TestResult{
			RouteID:   r.ID,
			RouteName: r.Name,
			Method:    string(r.Method),
			Path:      r.Path,
		}
		if err != nil {
			result.Error = err.Error()
		} else {
			result.StatusCode = resp.StatusCode
			result.Duration = resp.Duration
		}
		results = append(results, result)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	if _, err := fmt.Fprintln(w, "ID\tName\tMethod\tPath\tStatus\tDuration"); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	for _, r := range results {
		status := "Error"
		if r.Error == "" {
			status = fmt.Sprintf("%d", r.StatusCode)
		}
		if _, err := fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%v\n", r.RouteID, r.RouteName, r.Method, r.Path, status, r.Duration); err != nil {
			return fmt.Errorf("failed to write table line: %w", err)
		}
	}
	if err := w.Flush(); err != nil {
		return fmt.Errorf("failed to flush table: %w", err)
	}

	return nil
}
