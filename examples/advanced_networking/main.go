// Package main demonstrates advanced networking features of gcurl
// including --connect-to for connection redirection and -G for GET mode with data
package main

import (
	"fmt"
	"strings"

	"github.com/474420502/gcurl"
)

func main() {
	fmt.Println("=== Advanced Networking Features Demo ===\n")

	// Demo 1: Connection Redirection with --connect-to
	connectToDemo()

	fmt.Println("\n" + strings.Repeat("=", 50) + "\n")

	// Demo 2: GET Mode with Data using -G
	getModeDemo()

	fmt.Println("\n" + strings.Repeat("=", 50) + "\n")

	// Demo 3: Complex Integration Scenarios
	integrationDemo()
}

// connectToDemo demonstrates the --connect-to functionality
func connectToDemo() {
	fmt.Println("🌐 Connection Redirection Demo (--connect-to)")
	fmt.Println("Useful for testing, debugging, and development scenarios")

	scenarios := []struct {
		name        string
		command     string
		description string
	}{
		{
			name:        "Local Development",
			command:     `curl https://api.production.com/users --connect-to api.production.com:443:localhost:3000`,
			description: "Redirect production API calls to local development server",
		},
		{
			name:        "Load Balancer Testing",
			command:     `curl https://service.com/health --connect-to service.com:443:backend-1.internal:8080`,
			description: "Test specific backend server bypassing load balancer",
		},
		{
			name:        "Staging Environment",
			command:     `curl https://app.company.com/api --connect-to app.company.com:443:staging.internal.com:443`,
			description: "Point production domain to staging environment",
		},
		{
			name:        "Proxy All Traffic",
			command:     `curl https://example.com/test --connect-to ::proxy.company.com:8080`,
			description: "Route all connections through a proxy server",
		},
	}

	for i, scenario := range scenarios {
		fmt.Printf("%d. %s\n", i+1, scenario.name)
		fmt.Printf("   Description: %s\n", scenario.description)
		fmt.Printf("   Command: %s\n", scenario.command)

		curl, err := gcurl.Parse(scenario.command)
		if err != nil {
			fmt.Printf("   ❌ Parse error: %v\n", err)
			continue
		}

		fmt.Printf("   ✅ Parsed successfully:\n")
		fmt.Printf("      Target URL: %s\n", curl.ParsedURL.String())
		fmt.Printf("      Connection redirects: %d\n", len(curl.ConnectTo))
		for j, redirect := range curl.ConnectTo {
			fmt.Printf("        [%d] %s\n", j+1, redirect)
		}

		// Show verbose output
		if len(curl.ConnectTo) > 0 {
			curl.Verbose = true
			verbose := curl.VerboseInfo()
			fmt.Printf("   📋 Verbose preview (first 3 lines):\n")
			lines := splitLines(verbose)
			for k, line := range lines[:min(3, len(lines))] {
				if line != "" {
					fmt.Printf("      %s\n", line)
				}
				if k >= 2 {
					break
				}
			}
		}
		fmt.Println()
	}
}

// getModeDemo demonstrates the -G/--get functionality
func getModeDemo() {
	fmt.Println("🔍 GET Mode Demo (-G/--get)")
	fmt.Println("Convert POST data to query parameters for GET requests")

	scenarios := []struct {
		name        string
		command     string
		description string
	}{
		{
			name:        "Simple Search",
			command:     `curl -G -d "q=golang" -d "limit=10" https://api.github.com/search/repositories`,
			description: "Search GitHub repositories with query parameters",
		},
		{
			name:        "Complex Filters",
			command:     `curl -G -d "filters[status]=active" -d "filters[type]=user" -d "sort=created_at" https://api.example.com/users`,
			description: "Apply multiple filters and sorting",
		},
		{
			name:        "Analytics Query",
			command:     `curl -G -d "start_date=2023-01-01" -d "end_date=2023-12-31" -d "metrics=views,clicks,conversions" https://analytics.example.com/api`,
			description: "Query analytics data with date range and metrics",
		},
		{
			name:        "Pagination",
			command:     `curl -G -d "page=2" -d "per_page=50" -d "include=profile,settings" https://api.example.com/users`,
			description: "Paginated API request with includes",
		},
	}

	for i, scenario := range scenarios {
		fmt.Printf("%d. %s\n", i+1, scenario.name)
		fmt.Printf("   Description: %s\n", scenario.description)
		fmt.Printf("   Command: %s\n", scenario.command)

		curl, err := gcurl.Parse(scenario.command)
		if err != nil {
			fmt.Printf("   ❌ Parse error: %v\n", err)
			continue
		}

		fmt.Printf("   ✅ Parsed successfully:\n")
		fmt.Printf("      Method: %s\n", curl.Method)
		fmt.Printf("      GET Mode: %t\n", curl.GetMode)
		fmt.Printf("      Base URL: %s\n", curl.ParsedURL.String())

		// In a real implementation, data would be converted to query parameters
		if curl.Body != nil {
			fmt.Printf("      Has data for query parameter conversion: Yes\n")
			fmt.Printf("      Note: Data will be appended as query parameters when executed\n")
		}
		fmt.Println()
	}
}

// integrationDemo shows complex integration scenarios
func integrationDemo() {
	fmt.Println("🎯 Integration Demo")
	fmt.Println("Combining multiple advanced features")

	scenarios := []struct {
		name        string
		command     string
		description string
	}{
		{
			name:        "Local API Testing",
			command:     `curl -G -v -d "debug=true" -d "env=development" https://api.myapp.com/status --connect-to api.myapp.com:443:localhost:8080`,
			description: "Test production API locally with debug parameters",
		},
		{
			name:        "Staging with Authentication",
			command:     `curl -G -H "Authorization: Bearer test-token" -d "include=metadata" https://secure.example.com/api/data --connect-to secure.example.com:443:staging.internal:443`,
			description: "Authenticated request to staging environment",
		},
		{
			name:        "Proxy with Complex Query",
			command:     `curl -G -d "filters[date_range][start]=2023-01-01" -d "filters[date_range][end]=2023-12-31" https://analytics.company.com/api --connect-to ::proxy.company.com:8080`,
			description: "Complex analytics query through corporate proxy",
		},
	}

	for i, scenario := range scenarios {
		fmt.Printf("%d. %s\n", i+1, scenario.name)
		fmt.Printf("   Description: %s\n", scenario.description)
		fmt.Printf("   Command: %s\n", scenario.command)

		curl, err := gcurl.Parse(scenario.command)
		if err != nil {
			fmt.Printf("   ❌ Parse error: %v\n", err)
			continue
		}

		fmt.Printf("   ✅ Parsed successfully:\n")
		fmt.Printf("      Method: %s\n", curl.Method)
		fmt.Printf("      URL: %s\n", curl.ParsedURL.String())
		fmt.Printf("      GET Mode: %t\n", curl.GetMode)
		fmt.Printf("      Verbose: %t\n", curl.Verbose)
		fmt.Printf("      Headers: %d\n", len(curl.Header))
		fmt.Printf("      Connection redirects: %d\n", len(curl.ConnectTo))

		if len(curl.ConnectTo) > 0 {
			fmt.Printf("      Redirection details:\n")
			for j, redirect := range curl.ConnectTo {
				fmt.Printf("        [%d] %s\n", j+1, redirect)
			}
		}
		fmt.Println()
	}
}

// Helper functions
func splitLines(text string) []string {
	lines := []string{}
	current := ""
	for _, char := range text {
		if char == '\n' {
			lines = append(lines, current)
			current = ""
		} else {
			current += string(char)
		}
	}
	if current != "" {
		lines = append(lines, current)
	}
	return lines
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
