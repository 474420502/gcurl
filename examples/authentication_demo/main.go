// Package main demonstrates various authentication methods in gcurl
package main

import (
	"fmt"
	"strings"

	"github.com/474420502/gcurl"
)

func main() {
	fmt.Println("=== Authentication Methods Demo ===\n")

	// Demo 1: Basic Authentication
	basicAuthDemo()

	fmt.Println("\n" + strings.Repeat("=", 50) + "\n")

	// Demo 2: Bearer Token Authentication
	bearerTokenDemo()

	fmt.Println("\n" + strings.Repeat("=", 50) + "\n")

	// Demo 3: API Key Authentication
	apiKeyDemo()

	fmt.Println("\n" + strings.Repeat("=", 50) + "\n")

	// Demo 4: Complex Authentication Scenarios
	complexAuthDemo()
}

// basicAuthDemo demonstrates HTTP Basic Authentication
func basicAuthDemo() {
	fmt.Println("🔐 Basic Authentication Demo")
	fmt.Println("HTTP Basic Auth using username and password")

	scenarios := []struct {
		name        string
		command     string
		description string
	}{
		{
			name:        "Simple Basic Auth",
			command:     `curl -u "username:password" https://httpbin.org/basic-auth/username/password`,
			description: "Basic authentication with username and password",
		},
		{
			name:        "Interactive Password",
			command:     `curl -u "admin" https://secure.example.com/api`,
			description: "Username only - password will be prompted",
		},
		{
			name:        "Basic Auth with API",
			command:     `curl -u "api_user:api_secret" -H "Accept: application/json" https://api.example.com/users`,
			description: "API access with basic authentication",
		},
		{
			name:        "Basic Auth with POST",
			command:     `curl -u "user:pass" -X POST -d "name=test" https://api.example.com/create`,
			description: "Create resource with basic authentication",
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
		fmt.Printf("      User Auth: %s\n", curl.User)

		// Check for Authorization header
		if auth := curl.Header.Get("Authorization"); auth != "" {
			fmt.Printf("      Authorization: %s\n", maskAuthHeader(auth))
		}

		fmt.Println()
	}
}

// bearerTokenDemo demonstrates Bearer Token Authentication
func bearerTokenDemo() {
	fmt.Println("🎫 Bearer Token Authentication Demo")
	fmt.Println("Modern token-based authentication (JWT, OAuth2, etc.)")

	scenarios := []struct {
		name        string
		command     string
		description string
	}{
		{
			name:        "JWT Token",
			command:     `curl -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." https://api.example.com/profile`,
			description: "JWT token for user profile access",
		},
		{
			name:        "OAuth2 Access Token",
			command:     `curl -H "Authorization: Bearer gho_1234567890abcdef" https://api.github.com/user`,
			description: "GitHub API access with OAuth2 token",
		},
		{
			name:        "API with Bearer Token",
			command:     `curl -H "Authorization: Bearer sk-1234567890abcdef" -H "Content-Type: application/json" -d '{"prompt":"Hello"}' https://api.openai.com/v1/completions`,
			description: "OpenAI API call with bearer token",
		},
		{
			name:        "Microservice Auth",
			command:     `curl -H "Authorization: Bearer service-token-xyz" -H "X-Service-Name: user-service" https://internal.api.com/users`,
			description: "Internal microservice authentication",
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
		fmt.Printf("      Headers: %d\n", len(curl.Header))

		// Show authorization header
		if auth := curl.Header.Get("Authorization"); auth != "" {
			fmt.Printf("      Authorization: %s\n", maskAuthHeader(auth))
		}

		// Show other relevant headers
		for name, values := range curl.Header {
			if strings.ToLower(name) != "authorization" && isRelevantHeader(name) {
				fmt.Printf("      %s: %s\n", name, strings.Join(values, ", "))
			}
		}

		fmt.Println()
	}
}

// apiKeyDemo demonstrates API Key Authentication
func apiKeyDemo() {
	fmt.Println("🔑 API Key Authentication Demo")
	fmt.Println("API key authentication in headers or query parameters")

	scenarios := []struct {
		name        string
		command     string
		description string
	}{
		{
			name:        "API Key in Header",
			command:     `curl -H "X-API-Key: abc123def456ghi789" https://api.weather.com/v1/current`,
			description: "Weather API with custom header API key",
		},
		{
			name:        "Multiple API Keys",
			command:     `curl -H "X-API-Key: primary-key-123" -H "X-Client-ID: client-456" https://api.example.com/data`,
			description: "Service requiring both API key and client ID",
		},
		{
			name:        "API Key with Rate Limiting",
			command:     `curl -H "X-API-Key: rate-limited-key" -H "X-Rate-Limit-Tier: premium" https://api.service.com/premium`,
			description: "Premium API access with rate limiting headers",
		},
		{
			name:        "Legacy API Key",
			command:     `curl -H "Authorization: ApiKey username:api-key-value" https://legacy.api.com/endpoint`,
			description: "Legacy API using ApiKey authorization scheme",
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

		// Show API-related headers
		apiHeaders := []string{"X-API-Key", "X-Client-ID", "X-Rate-Limit-Tier", "Authorization"}
		for _, headerName := range apiHeaders {
			if value := curl.Header.Get(headerName); value != "" {
				fmt.Printf("      %s: %s\n", headerName, maskSensitiveValue(value))
			}
		}

		fmt.Println()
	}
}

// complexAuthDemo demonstrates complex authentication scenarios
func complexAuthDemo() {
	fmt.Println("🎭 Complex Authentication Scenarios")
	fmt.Println("Real-world complex authentication patterns")

	scenarios := []struct {
		name        string
		command     string
		description string
	}{
		{
			name:        "Multi-Factor API Access",
			command:     `curl -H "Authorization: Bearer access-token" -H "X-MFA-Token: 123456" -H "X-Device-ID: device-abc" https://secure.api.com/sensitive`,
			description: "Multi-factor authentication with device tracking",
		},
		{
			name:        "Service-to-Service Auth",
			command:     `curl -H "Authorization: Bearer service-jwt" -H "X-Service-Name: payment-service" -H "X-Request-ID: req-12345" https://internal.api.com/validate`,
			description: "Internal service authentication with tracing",
		},
		{
			name:        "Federated Identity",
			command:     `curl -H "Authorization: Bearer saml-token" -H "X-Identity-Provider: company-sso" -H "X-User-Context: department=engineering" https://federated.api.com/resources`,
			description: "Federated identity with SAML tokens",
		},
		{
			name:        "API Gateway Auth",
			command:     `curl -H "X-API-Gateway-Key: gateway-key" -H "Authorization: Bearer user-token" -H "X-Forwarded-For: 192.168.1.100" https://gateway.api.com/proxy`,
			description: "API Gateway with pass-through authentication",
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
		fmt.Printf("      Total Headers: %d\n", len(curl.Header))

		// Show authentication and security headers
		securityHeaders := []string{
			"Authorization", "X-MFA-Token", "X-Device-ID", "X-Service-Name",
			"X-Request-ID", "X-Identity-Provider", "X-User-Context",
			"X-API-Gateway-Key", "X-Forwarded-For",
		}

		fmt.Printf("      Security Headers:\n")
		headerCount := 0
		for _, headerName := range securityHeaders {
			if value := curl.Header.Get(headerName); value != "" {
				fmt.Printf("        %s: %s\n", headerName, maskSensitiveValue(value))
				headerCount++
			}
		}

		if headerCount == 0 {
			fmt.Printf("        (No security headers found)\n")
		}

		fmt.Println()
	}
}

// Helper functions
func maskAuthHeader(auth string) string {
	if strings.HasPrefix(auth, "Bearer ") {
		token := auth[7:]
		if len(token) > 10 {
			return "Bearer " + token[:6] + "..." + token[len(token)-4:]
		}
		return "Bearer " + strings.Repeat("*", len(token))
	}
	if strings.HasPrefix(auth, "Basic ") {
		return "Basic " + strings.Repeat("*", len(auth)-6)
	}
	if strings.HasPrefix(auth, "ApiKey ") {
		return "ApiKey " + strings.Repeat("*", len(auth)-7)
	}
	return strings.Repeat("*", len(auth))
}

func maskSensitiveValue(value string) string {
	// For tokens and keys, show first 6 and last 4 characters
	if len(value) > 15 && (strings.Contains(strings.ToLower(value), "token") ||
		strings.Contains(strings.ToLower(value), "key") ||
		strings.Contains(value, "-") ||
		len(value) > 20) {
		return value[:6] + "..." + value[len(value)-4:]
	}
	// For shorter values, mask partially
	if len(value) > 8 {
		return value[:3] + strings.Repeat("*", len(value)-6) + value[len(value)-3:]
	}
	// For very short values, mask completely
	return strings.Repeat("*", len(value))
}

func isRelevantHeader(name string) bool {
	relevantHeaders := []string{
		"Content-Type", "Accept", "X-Service-Name", "X-Client-ID",
		"X-Rate-Limit-Tier", "X-MFA-Token", "X-Device-ID",
	}
	for _, relevant := range relevantHeaders {
		if strings.EqualFold(name, relevant) {
			return true
		}
	}
	return false
}
