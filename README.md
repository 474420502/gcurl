# Parse cURL To Golang Requests

![Go Version](https://img.shields.io/badge/Go-1.20+-00ADD8?style=flat&logo=go)
![Test Coverage](https://img.shields.io/badge/Coverage-79.1%25-brightgreen)
![License](https://img.shields.io/badge/License-MIT-blue)
![Latest Release](https://img.shields.io/github/v/release/474420502/gcurl)

A powerful Go library that converts cURL commands into Go HTTP requests with full feature compatibility.

* Based on the robust [requests](https://github.com/474420502/requests) library
* Seamlessly transform cURL bash commands to Go code
* Inherits all cURL functionality while adding Go's flexibility for configuration, cookies, headers, and URL handling
* Supports both Bash and Windows Cmd cURL formats
* Production-ready with comprehensive test coverage

## 🚀 Features

- 🌐 **Complete cURL command parsing** - Parse any cURL command to Go requests
- 🔧 **Full cURL compatibility** - Support for all major cURL options
- 📁 **File output support** - Save responses to files with `-o`, `-O`, `--output-dir`, etc.
- 🔐 **Authentication** - Basic auth, digest auth, and bearer tokens
- 🌐 **HTTP protocol control** - HTTP/1.0, HTTP/1.1, HTTP/2 support
- 🍪 **Cookie management** - Full cookie handling and session persistence
- 📝 **Form data & file uploads** - Multipart forms and file uploads
- 🔄 **Redirect handling** - Automatic redirect following with limits
- ⏱️ **Timeout controls** - Connection and request timeouts
- 🔐 **Proxy support** - HTTP/HTTPS/SOCKS5 proxy configuration
- 🛡️ **SSL/TLS options** - Custom certificates and SSL verification control
- 🎯 **Debug output** - Detailed debugging information like cURL `-v`

## 📦 Installation

```bash
go get github.com/474420502/gcurl@v1.1.0
```

## 🎯 Quick Start

Transform any cURL command into Go code instantly:

```go
package main

import (
	"fmt"
	"log"
	
	"github.com/474420502/gcurl"
)

func main() {
	// Parse any cURL command
	curl, err := gcurl.Parse(`curl -H "Accept: application/json" https://httpbin.org/get`)
	if err != nil {
		log.Fatal(err)
	}
	
	// Execute the request
	resp, err := curl.Request().Execute()
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("Status: %d\n", resp.GetStatusCode())
	fmt.Printf("Response: %s\n", resp.ContentString())
}
```

## 📖 Examples

### Example 1: Basic GET Request with Headers

Transform cURL commands with custom headers:

```go
package main

import (
	"fmt"
	"log"
	
	"github.com/474420502/gcurl"
)

func main() {
	// Parse cURL command with multiple headers
	curlCmd := `curl "https://httpbin.org/get" \
		-H "Accept: application/json" \
		-H "User-Agent: MyApp/1.0" \
		-H "Authorization: Bearer token123"`
	
	curl, err := gcurl.Parse(curlCmd)
	if err != nil {
		log.Fatal(err)
	}
	
	// Create session and execute
	session := curl.CreateSession()
	resp, err := curl.CreateRequest(session).Execute()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Status: %d\n", resp.GetStatusCode())
	fmt.Printf("Headers sent: %v\n", session.GetHeader())
	fmt.Printf("Response: %s\n", resp.ContentString())
}
```

### Example 2: POST Request with JSON Data

Handle POST requests with JSON payload:

```go
package main

import (
	"fmt"
	"log"
	
	"github.com/474420502/gcurl"
)

func main() {
	curlCmd := `curl -X POST "https://httpbin.org/post" \
		-H "Content-Type: application/json" \
		-d '{"name":"John Doe","email":"john@example.com","age":30}'`
	
	curl, err := gcurl.Parse(curlCmd)
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("Method: %s\n", curl.Method)
	fmt.Printf("Content-Type: %s\n", curl.ContentType)
	fmt.Printf("Request Body: %s\n", curl.Body.Content)
	
	resp, err := curl.Request().Execute()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Response Status: %d\n", resp.GetStatusCode())
	fmt.Printf("Response: %s\n", resp.ContentString())
}
```

### Example 3: Multipart Form Data with File Upload

Upload files using multipart form data:

```go
package main

import (
	"fmt"
	"log"
	"os"
	
	"github.com/474420502/gcurl"
)

func main() {
	// Create a test file
	testFile := "/tmp/sample.txt"
	err := os.WriteFile(testFile, []byte("This is a test file content"), 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(testFile)
	
	curlCmd := fmt.Sprintf(`curl -X POST "https://httpbin.org/post" \
		-F "file=@%s" \
		-F "name=John" \
		-F "description=Sample file upload"`, testFile)
	
	curl, err := gcurl.Parse(curlCmd)
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("Form upload with %d fields\n", len(curl.Body.Forms))
	
	resp, err := curl.Request().Execute()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Upload Status: %d\n", resp.GetStatusCode())
	fmt.Printf("Response: %s\n", resp.ContentString())
}
```

### Example 4: File Download with Output Options

Save responses to files with various output configurations:

```go
package main

import (
	"fmt"
	"log"
	"os"
	
	"github.com/474420502/gcurl"
)

func main() {
	examples := []struct {
		name        string
		command     string
		description string
	}{
		{
			name:        "Save to specific file",
			command:     `curl -o /tmp/response.json https://httpbin.org/json`,
			description: "Save response to a specific file",
		},
		{
			name:        "Use remote filename",
			command:     `curl -O https://httpbin.org/robots.txt`,
			description: "Use the remote filename (robots.txt)",
		},
		{
			name:        "Save to directory",
			command:     `curl -O --output-dir /tmp/downloads --create-dirs https://httpbin.org/uuid`,
			description: "Save to specific directory, create if needed",
		},
		{
			name:        "Resume download",
			command:     `curl -C 1024 -o /tmp/partial.dat https://httpbin.org/bytes/2048`,
			description: "Resume download from byte offset 1024",
		},
		{
			name:        "Auto-resume",
			command:     `curl -C - -o /tmp/auto_resume.dat https://httpbin.org/bytes/4096`,
			description: "Auto-detect existing file size and resume",
		},
	}

	for i, example := range examples {
		fmt.Printf("\n%d. %s\n", i+1, example.name)
		fmt.Printf("   Description: %s\n", example.description)
		fmt.Printf("   Command: %s\n", example.command)

		curl, err := gcurl.Parse(example.command)
		if err != nil {
			fmt.Printf("   ❌ Parse error: %v\n", err)
			continue
		}

		// Show parsed configuration
		if curl.OutputFile != "" {
			fmt.Printf("   📁 Output file: %s\n", curl.OutputFile)
		}
		if curl.RemoteName {
			fmt.Printf("   🌐 Using remote filename\n")
		}
		if curl.OutputDir != "" {
			fmt.Printf("   📂 Output directory: %s\n", curl.OutputDir)
		}
		if curl.ContinueAt > 0 {
			fmt.Printf("   ⏩ Resume from byte: %d\n", curl.ContinueAt)
		} else if curl.ContinueAt == -1 {
			fmt.Printf("   🔄 Auto-resume enabled\n")
		}

		fmt.Printf("   ✅ Configuration parsed successfully\n")
	}
}
```

### Example 5: Authentication Methods

Handle various authentication scenarios:

```go
package main

import (
	"fmt"
	"log"
	
	"github.com/474420502/gcurl"
)

func main() {
	authExamples := []struct {
		name        string
		command     string
		description string
	}{
		{
			name:        "Basic Authentication",
			command:     `curl -u "username:password" "https://httpbin.org/basic-auth/username/password"`,
			description: "HTTP Basic Authentication",
		},
		{
			name:        "Digest Authentication",
			command:     `curl --digest -u "user:pass" "https://httpbin.org/digest-auth/auth/user/pass"`,
			description: "HTTP Digest Authentication",
		},
		{
			name:        "Bearer Token",
			command:     `curl -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9" "https://httpbin.org/bearer"`,
			description: "Bearer token authentication",
		},
		{
			name:        "API Key Header",
			command:     `curl -H "X-API-Key: your-api-key-here" "https://httpbin.org/get"`,
			description: "API Key in custom header",
		},
	}

	for i, example := range authExamples {
		fmt.Printf("\n%d. %s\n", i+1, example.name)
		fmt.Printf("   Description: %s\n", example.description)
		fmt.Printf("   Command: %s\n", example.command)

		curl, err := gcurl.Parse(example.command)
		if err != nil {
			fmt.Printf("   ❌ Parse error: %v\n", err)
			continue
		}

		// Show authentication configuration
		if curl.Auth != nil {
			fmt.Printf("   🔐 Basic auth configured\n")
		}
		if curl.DigestAuth != nil {
			fmt.Printf("   🔐 Digest auth configured\n")
		}
		if len(curl.Headers) > 0 {
			fmt.Printf("   📋 Headers configured: %d\n", len(curl.Headers))
		}

		fmt.Printf("   ✅ Authentication parsed successfully\n")
	}
}
```

### Example 6: HTTP Version Control

Control HTTP protocol versions:

```go
package main

import (
	"fmt"
	"log"
	
	"github.com/474420502/gcurl"
)

func main() {
	versionExamples := []struct {
		name    string
		command string
		version string
	}{
		{
			name:    "HTTP/1.0",
			command: `curl --http1.0 https://httpbin.org/get`,
			version: "HTTP/1.0",
		},
		{
			name:    "HTTP/1.1",
			command: `curl --http1.1 https://httpbin.org/get`,
			version: "HTTP/1.1",
		},
		{
			name:    "HTTP/2",
			command: `curl --http2 https://httpbin.org/get`,
			version: "HTTP/2",
		},
		{
			name:    "Auto-detect",
			command: `curl https://httpbin.org/get`,
			version: "Auto",
		},
	}

	fmt.Println("HTTP Version Control Examples:")
	fmt.Println("==============================")

	for i, example := range versionExamples {
		fmt.Printf("\n%d. %s\n", i+1, example.name)
		fmt.Printf("   Command: %s\n", example.command)

		curl, err := gcurl.Parse(example.command)
		if err != nil {
			fmt.Printf("   ❌ Parse error: %v\n", err)
			continue
		}

		fmt.Printf("   🌐 HTTP Version: %s\n", curl.HTTPVersion.String())
		fmt.Printf("   ✅ Version control configured\n")
	}
}
```

### Example 7: Debug and Verbose Output

Get detailed debugging information like cURL's verbose mode:

```go
package main

import (
	"fmt"
	"log"
	
	"github.com/474420502/gcurl"
)

func main() {
	curlCmd := `curl -v -X POST https://api.example.com/data \
		-H "Content-Type: application/json" \
		-H "Authorization: Bearer token123" \
		-H "User-Agent: MyApp/2.0" \
		-d '{"action":"create","data":{"name":"test"}}'`

	curl, err := gcurl.Parse(curlCmd)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== Curl Command Parsed ===")
	fmt.Printf("URL: %s\n", curl.ParsedURL.String())
	fmt.Printf("Method: %s\n", curl.Method)
	fmt.Printf("Verbose Mode: %v\n", curl.Verbose)

	fmt.Println("\n=== Debug Output (like curl -v) ===")
	fmt.Println(curl.Debug())

	fmt.Println("\n=== Verbose Info Simulation ===")
	fmt.Println(curl.VerboseInfo())

	fmt.Println("\n=== Quick Summary ===")
	fmt.Println(curl.Summary())
}
```

### Example 8: Advanced Session Management

Reuse sessions for multiple requests:

```go
package main

import (
	"fmt"
	"log"
	"time"
	"net/http"
	
	"github.com/474420502/gcurl"
)

func main() {
	// First request: login and get session cookie
	loginCmd := `curl -X POST "https://httpbin.org/cookies/set/sessionid/abc123" \
		-H "Content-Type: application/json" \
		-d '{"username":"testuser","password":"testpass"}'`

	// Second request: use the session
	dataCmd := `curl "https://httpbin.org/cookies" \
		-H "Accept: application/json"`

	// Parse both commands
	loginCurl, err := gcurl.Parse(loginCmd)
	if err != nil {
		log.Fatal(err)
	}

	dataCurl, err := gcurl.Parse(dataCmd)
	if err != nil {
		log.Fatal(err)
	}

	// Create a shared session
	session := loginCurl.CreateSession()

	// Customize session with timeout and custom headers
	session.Config().SetTimeout(30 * time.Second)
	
	customHeaders := make(http.Header)
	customHeaders.Set("X-Client-Version", "v2.0")
	customHeaders.Set("X-Request-ID", "req-12345")
	session.SetHeader(customHeaders)

	fmt.Println("=== Executing Login Request ===")
	loginResp, err := loginCurl.CreateRequest(session).Execute()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Login Status: %d\n", loginResp.GetStatusCode())

	fmt.Println("\n=== Executing Data Request (with session) ===")
	dataResp, err := dataCurl.CreateRequest(session).Execute()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Data Status: %d\n", dataResp.GetStatusCode())
	fmt.Printf("Response: %s\n", dataResp.ContentString())

	fmt.Println("\n=== Session cookies were automatically shared ===")
}
```

## 🔧 Advanced Usage

### Session Reuse Pattern

```go
// Create multiple cURL parsers
curl1, _ := gcurl.Parse(`curl "https://httpbin.org/cookies/set/session/123"`)
curl2, _ := gcurl.Parse(`curl "https://httpbin.org/cookies"`)
curl3, _ := gcurl.Parse(`curl "https://httpbin.org/user-agent"`)

// Create a shared session for connection pooling and cookie persistence
session := curl1.CreateSession()

// Execute all requests with the same session
resp1, _ := curl1.CreateRequest(session).Execute()
resp2, _ := curl2.CreateRequest(session).Execute() // Cookies from resp1 are available
resp3, _ := curl3.CreateRequest(session).Execute() // Same session, connection reuse
```

### Direct Execution for Simple Cases

```go
// For one-off requests, use direct execution
curl, _ := gcurl.Parse(`curl "https://httpbin.org/get"`)
resp, err := curl.Request().Execute() // Auto-creates session
```

### Custom Session Configuration

```go
curl, _ := gcurl.Parse(`curl "https://httpbin.org/delay/2"`)
session := curl.CreateSession()

// Configure timeouts
session.Config().SetTimeout(10 * time.Second)

// Add custom headers to all requests in this session
headers := make(http.Header)
headers.Set("X-Client-ID", "my-app")
headers.Set("X-API-Version", "v2")
session.SetHeader(headers)

// Use the customized session
resp, err := curl.CreateRequest(session).Execute()
```

## 📊 API Reference

### Core Functions

| Function | Description | Example |
|----------|-------------|---------|
| `gcurl.Parse(cmd string)` | Parse any cURL command | `curl, err := gcurl.Parse("curl https://api.example.com")` |
| `gcurl.ParseBash(cmd string)` | Parse Bash-style cURL | `curl, err := gcurl.ParseBash("curl 'https://api.example.com'")` |
| `gcurl.ParseCmd(cmd string)` | Parse Windows CMD-style | `curl, err := gcurl.ParseCmd("curl \"https://api.example.com\"")` |

### CURL Object Methods

| Method | Description | Return Type |
|--------|-------------|-------------|
| `CreateSession()` | Create a new HTTP session | `*requests.Session` |
| `CreateRequest(session)` | Create request with session | `*requests.Request` |
| `Request()` | Create request (auto-session) | `*requests.Request` |
| `Debug()` | Get detailed debug info | `string` |
| `VerboseInfo()` | Get verbose output like `curl -v` | `string` |
| `Summary()` | Get brief summary | `string` |

### Response Object Methods

| Method | Description | Return Type |
|--------|-------------|-------------|
| `GetStatusCode()` | HTTP status code | `int` |
| `Content()` | Response body as bytes | `[]byte` |
| `ContentString()` | Response body as string | `string` |
| `GetHeader(key)` | Get response header | `string` |
| `GetHeaders()` | Get all response headers | `http.Header` |

## ✅ Supported cURL Options

### Comprehensive Feature Matrix

| Category | cURL Option | Description | Status | Example |
|----------|-------------|-------------|--------|---------|
| **HTTP Methods** | `-X, --request` | HTTP method (GET, POST, etc.) | ✅ | `curl -X POST` |
| **Headers** | `-H, --header` | Custom headers | ✅ | `curl -H "Accept: application/json"` |
| **Request Body** | `-d, --data` | Send POST data | ✅ | `curl -d "name=value"` |
| | `--data-raw` | Send raw data | ✅ | `curl --data-raw '{"json":true}'` |
| | `--data-urlencode` | URL encode data | ✅ | `curl --data-urlencode "name=John Doe"` |
| | `-F, --form` | Multipart form data | ✅ | `curl -F "file=@path/file.txt"` |
| **Authentication** | `-u, --user` | Basic authentication | ✅ | `curl -u "user:pass"` |
| | `--digest` | Digest authentication | ✅ | `curl --digest -u "user:pass"` |
| **Cookies** | `-b, --cookie` | Send cookies | ✅ | `curl -b "session=abc123"` |
| | `-c, --cookie-jar` | Save cookies to file | ✅ | `curl -c cookies.txt` |
| **File Operations** | `-o, --output` | Write output to file | ✅ | `curl -o output.txt` |
| | `-O, --remote-name` | Use remote filename | ✅ | `curl -O` |
| | `--output-dir` | Output directory | ✅ | `curl --output-dir /downloads` |
| | `--create-dirs` | Create output directories | ✅ | `curl --create-dirs` |
| | `-C, --continue-at` | Resume/continue transfer | ✅ | `curl -C 1024` |
| | `--remove-on-error` | Remove file on HTTP error | ✅ | `curl --remove-on-error` |
| **HTTP Versions** | `--http1.0` | Force HTTP/1.0 | ✅ | `curl --http1.0` |
| | `--http1.1` | Force HTTP/1.1 | ✅ | `curl --http1.1` |
| | `--http2` | Force HTTP/2 | ✅ | `curl --http2` |
| **Redirects** | `-L, --location` | Follow redirects | ✅ | `curl -L` |
| | `--max-redirs` | Maximum redirect count | ✅ | `curl --max-redirs 5` |
| **Timeouts** | `--connect-timeout` | Connection timeout | ✅ | `curl --connect-timeout 10` |
| | `--max-time` | Maximum total time | ✅ | `curl --max-time 30` |
| **Proxy** | `--proxy` | Use proxy server | ✅ | `curl --proxy http://proxy:8080` |
| | `--proxy-user` | Proxy authentication | ✅ | `curl --proxy-user "user:pass"` |
| **SSL/TLS** | `-k, --insecure` | Skip SSL verification | ✅ | `curl -k` |
| | `--cacert` | CA certificate file | ✅ | `curl --cacert ca.pem` |
| | `--cert` | Client certificate | ✅ | `curl --cert client.pem` |
| | `--key` | Client private key | ✅ | `curl --key client.key` |
| **Output Control** | `-v, --verbose` | Verbose output | ✅ | `curl -v` |
| | `-i, --include` | Include response headers | ✅ | `curl -i` |
| | `-I, --head` | HEAD request only | ✅ | `curl -I` |
| | `-s, --silent` | Silent mode | ✅ | `curl -s` |
| **User Agent** | `-A, --user-agent` | Set User-Agent | ✅ | `curl -A "MyApp/1.0"` |
| **Compression** | `--compressed` | Accept compressed response | ✅ | `curl --compressed` |
| **Range Requests** | `-r, --range` | Byte range request | ✅ | `curl -r 0-1023` |

## 🔍 Debug and Troubleshooting

### Enable Debug Output

```go
curl, _ := gcurl.Parse(`curl -v -X POST -H "Content-Type: application/json" -d '{"test":true}' https://api.example.com`)

// Get detailed debug information
fmt.Println("=== Debug Output ===")
fmt.Println(curl.Debug())

// Get verbose info (simulates curl -v)
fmt.Println("=== Verbose Output ===")
fmt.Println(curl.VerboseInfo())
```

**Sample Debug Output:**
```
=== CURL Debug Information ===
Method: POST
URL: https://api.example.com
  Scheme: https
  Host: api.example.com
  Path: /
Headers (1):
  Content-Type: application/json
Body:
  Type: raw
  Length: 13 bytes
  Content: {"test":true}
Network Configuration:
  Timeout: 30s
  HTTP Version: Auto
Debug Flags: verbose
===============================
```

### Common Issues and Solutions

| Issue | Symptom | Solution |
|-------|---------|----------|
| **URL not parsed** | `invalid or malformed URL` | Ensure URL is properly quoted: `curl "https://example.com"` |
| **Headers ignored** | Headers not sent | Check header format: `curl -H "Key: Value"` |
| **File upload fails** | Multipart parsing error | Verify file exists: `curl -F "file=@/path/to/file"` |
| **Auth not working** | 401 Unauthorized | Check credentials format: `curl -u "user:password"` |
| **SSL errors** | Certificate verification failed | Use `-k` to skip verification: `curl -k https://...` |
| **Timeout issues** | Request hangs | Set timeouts: `curl --max-time 30 --connect-timeout 10` |

### Best Practices

1. **Always handle errors**:
   ```go
   curl, err := gcurl.Parse(curlCommand)
   if err != nil {
       log.Printf("Parse error: %v", err)
       return
   }
   ```

2. **Reuse sessions for multiple requests**:
   ```go
   session := curl.CreateSession()
   // Use session for multiple requests to same host
   ```

3. **Set appropriate timeouts**:
   ```go
   session.Config().SetTimeout(30 * time.Second)
   ```

4. **Use debug mode during development**:
   ```go
   fmt.Println(curl.Debug()) // See parsed configuration
   ```

## 🎯 Performance Tips

### Connection Reuse
```go
// Good: Reuse session for multiple requests
session := curl1.CreateSession()
resp1, _ := curl1.CreateRequest(session).Execute()
resp2, _ := curl2.CreateRequest(session).Execute() // Reuses connection

// Avoid: Creating new session for each request
resp1, _ := curl1.Request().Execute() // New session
resp2, _ := curl2.Request().Execute() // Another new session
```

### Timeout Configuration
```go
// Set reasonable timeouts
session := curl.CreateSession()
session.Config().SetTimeout(30 * time.Second)          // Total timeout
session.Config().SetConnectTimeout(10 * time.Second)   // Connection timeout
```

### Memory Efficiency
```go
// For large responses, consider streaming
resp, err := curl.Request().Execute()
if err == nil {
    // Process response.Content() efficiently
    // Don't store large responses in memory unnecessarily
}
```

## 🧪 Testing

### Run Tests

```bash
# Run all tests
go test -v

# Run tests with coverage report
go test -v -cover

# Run specific test categories
go test -v -run TestFileOutput
go test -v -run TestAuthentication
go test -v -run TestHTTPVersion

# Run tests with race detection
go test -v -race

# Generate detailed coverage report
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Test Statistics

- **Total Tests**: 100+ comprehensive test cases
- **Coverage**: 79.1% of code statements
- **Categories**: 
  - Basic parsing tests
  - HTTP method tests
  - Authentication tests
  - File output tests
  - HTTP version control tests
  - Error handling tests
  - Integration tests

## 🤝 Contributing

We welcome contributions from the community! Here's how you can help:

### Ways to Contribute

1. **🐛 Report Bugs**
   - Open an issue with detailed reproduction steps
   - Include the cURL command that's causing problems
   - Provide Go version and OS information

2. **💡 Suggest Features**
   - Propose new cURL options to support
   - Share use cases and examples
   - Discuss implementation approaches

3. **🔧 Submit Pull Requests**
   - Add tests for any new functionality
   - Follow existing code style and conventions
   - Update documentation as needed

4. **📚 Improve Documentation**
   - Fix typos and clarify explanations
   - Add more examples and use cases
   - Translate documentation

### Development Setup

```bash
# Clone the repository
git clone https://github.com/474420502/gcurl.git
cd gcurl

# Install dependencies
go mod tidy

# Run tests to verify setup
go test -v

# Run tests with coverage
go test -v -cover

# Check code formatting
go fmt ./...

# Run static analysis
go vet ./...
```

### Code Style Guidelines

- Follow standard Go conventions
- Write comprehensive tests for new features
- Include examples in documentation
- Use descriptive variable and function names
- Add comments for complex logic

### Pull Request Process

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes with tests
4. Ensure all tests pass (`go test -v`)
5. Commit your changes (`git commit -m 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

## 🚀 Roadmap

### Phase 3: Advanced Features (Upcoming)

- **Enhanced Session Management**
  - Session persistence and restoration
  - Session configuration templates
  - Advanced connection pooling

- **Performance Monitoring**
  - Request timing and statistics
  - Connection reuse metrics
  - Memory usage optimization

- **Advanced Authentication**
  - OAuth 2.0 flow support
  - JWT token handling
  - API key management

- **Concurrent Processing**
  - Parallel request execution
  - Batch command processing
  - Request queue management

### Long-term Vision

- Become the standard cURL-to-Go conversion library
- Support for enterprise-grade features
- Integration with popular Go frameworks
- Advanced debugging and profiling tools

## � Version History

### v1.1.0 (Current)
- ✅ Comprehensive file output support (`-o`, `-O`, `--output-dir`, etc.)
- ✅ Complete Digest Authentication implementation
- ✅ HTTP protocol version control
- ✅ Enhanced debugging capabilities
- ✅ Improved test coverage (79.1%)
- ✅ Production-ready stability

## �📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- **[requests](https://github.com/474420502/requests)** - The powerful HTTP library that powers gcurl
- **[cURL](https://curl.se/)** - The amazing command-line tool that inspired this project
- **Go Community** - For the excellent ecosystem and tools
- **All Contributors** - Everyone who has helped improve this library

## 🔗 Related Projects

- [mholt/curl-to-go](https://github.com/mholt/curl-to-go) - Online cURL to Go converter
- [sj26/curl-to-go](https://github.com/sj26/curl-to-go) - Another cURL conversion tool
- [requests](https://github.com/474420502/requests) - The HTTP client library used by gcurl

## 📞 Support

- **GitHub Issues**: [Report bugs and request features](https://github.com/474420502/gcurl/issues)
- **Documentation**: Check this README and code examples
- **Community**: Join discussions in GitHub issues

---

⭐ **Star this project** if it helps you convert cURL commands to Go code! Your support motivates continued development and improvements.

## 🏆 Success Stories

> "gcurl saved us hours of manual conversion work when migrating our API tests from shell scripts to Go. The authentication and file upload features work flawlessly!" - *Development Team Lead*

> "Perfect for our CI/CD pipeline where we needed to convert existing cURL-based health checks to Go services. The session reuse feature improved our performance significantly." - *DevOps Engineer*

> "The debug output helped us understand exactly how our complex cURL commands were being interpreted. Great for troubleshooting API integration issues." - *Backend Developer*

---

**Made with ❤️ for the Go community**
