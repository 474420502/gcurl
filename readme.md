# Parse cURL To Golang Requests

![Go Version](https://img.shields.io/badge/Go-1.20+-00ADD8?style=flat&logo=go)
![Test Coverage](https://img.shields.io/badge/Coverage-79.1%25-brightgreen)
![License](https://img.shields.io/badge/License-MIT-blue)
![Latest Release](https://img.shields.io/github/v/release/474420502/gcurl)

* Based on [requests](https://github.com/474420502/requests) library
* Easy to transform cURL bash commands to Go code
* Inherits all cURL functionality and adds Go's flexibility for configuration, cookies, headers, and URL handling
* Supports both Bash and Windows Cmd cURL formats

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
go get github.com/474420502/gcurl
```

## 🎯 Quick Start

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

Transform cURL commands with headers into Go requests:

```go
package main

import (
	"fmt"
	"log"
	
	"github.com/474420502/gcurl"
)

func main() {
	// Parse cURL command
	surl := `curl "http://httpbin.org/get" -H 'Connection: keep-alive' -H 'Accept-Encoding: gzip, deflate' -H 'Accept-Language: zh-CN,zh;q=0.9'`
	curl, err := gcurl.Parse(surl)
	if err != nil {
		log.Fatal(err)
	}
	
	// Create session and temporary request
	ses := curl.CreateSession()
	tp := curl.CreateRequest(ses)
	
	// Check headers
	fmt.Println("Headers:", ses.GetHeader())
	// Output: map[Accept-Encoding:[gzip, deflate] Accept-Language:[zh-CN,zh;q=0.9] Connection:[keep-alive]]
	
	// Execute request
	resp, err := tp.Execute()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Response:", string(resp.Content()))
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
	scurl := `curl -X POST "https://httpbin.org/post" -H "Content-Type: application/json" -d '{"name":"test","age":25}'`
	
	curl, err := gcurl.Parse(scurl)
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("Method: %s\n", curl.Method)          // POST
	fmt.Printf("Content-Type: %s\n", curl.ContentType) // application/json
	
	resp, err := curl.Request().Execute()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Response:", string(resp.Content()))
}
```

### Example 3: File Upload with Form Data

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
	testFile := "/tmp/test.txt"
	err := os.WriteFile(testFile, []byte("test file content"), 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(testFile)
	
	scurl := fmt.Sprintf(`curl -X POST "https://httpbin.org/post" -F "file=@%s" -F "description=test file"`, testFile)
	
	curl, err := gcurl.Parse(scurl)
	if err != nil {
		log.Fatal(err)
	}
	
	resp, err := curl.Request().Execute()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("File upload response:", string(resp.Content()))
}
```

### Example 4: File Output and Downloads

Save responses to files with various output options:

```go
package main

import (
	"fmt"
	"log"
	
	"github.com/474420502/gcurl"
)

func main() {
	// Save to specific file
	curl1, _ := gcurl.Parse(`curl -o output.json https://httpbin.org/json`)
	
	// Use remote filename
	curl2, _ := gcurl.Parse(`curl -O https://httpbin.org/robots.txt`)
	
	// Save to directory with automatic directory creation
	curl3, _ := gcurl.Parse(`curl -O --output-dir /tmp/downloads --create-dirs https://httpbin.org/uuid`)
	
	// Resume download from offset
	curl4, _ := gcurl.Parse(`curl -C 1024 -o partial.dat https://httpbin.org/bytes/2048`)
	
	// Auto-resume (detect existing file size)
	curl5, _ := gcurl.Parse(`curl -C - -o resume.dat https://httpbin.org/bytes/4096`)
	
	fmt.Println("File output configurations ready!")
	fmt.Printf("curl1 output: %s\n", curl1.OutputFile)
	fmt.Printf("curl2 remote name: %v\n", curl2.RemoteName)
	fmt.Printf("curl3 output dir: %s\n", curl3.OutputDir)
	fmt.Printf("curl4 continue at: %d bytes\n", curl4.ContinueAt)
	fmt.Printf("curl5 auto-resume: %d (-1 means auto)\n", curl5.ContinueAt)
}
```

### Example 5: Authentication

Handle various authentication methods:

```go
package main

import (
	"fmt"
	"log"
	
	"github.com/474420502/gcurl"
)

func main() {
	// Basic authentication
	curl1, err := gcurl.Parse(`curl -u "user:password" "https://httpbin.org/basic-auth/user/password"`)
	if err != nil {
		log.Fatal(err)
	}
	
	// Digest authentication
	curl2, err := gcurl.Parse(`curl --digest -u "user:password" "https://httpbin.org/digest-auth/auth/user/password"`)
	if err != nil {
		log.Fatal(err)
	}
	
	// Bearer token
	curl3, err := gcurl.Parse(`curl -H "Authorization: Bearer eyJhbGci..." "https://httpbin.org/bearer"`)
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("Basic auth configured: %v\n", curl1.Auth != nil)
	fmt.Printf("Digest auth configured: %v\n", curl2.DigestAuth != nil)
	fmt.Printf("Bearer token in headers: %v\n", len(curl3.Headers) > 0)
}
```

### Example 6: HTTP Version Control

Control HTTP protocol version:

```go
package main

import (
	"fmt"
	"log"
	
	"github.com/474420502/gcurl"
)

func main() {
	// Force HTTP/1.0
	curl1, _ := gcurl.Parse(`curl --http1.0 https://httpbin.org/get`)
	
	// Force HTTP/1.1
	curl2, _ := gcurl.Parse(`curl --http1.1 https://httpbin.org/get`)
	
	// Force HTTP/2
	curl3, _ := gcurl.Parse(`curl --http2 https://httpbin.org/get`)
	
	fmt.Printf("HTTP/1.0: %s\n", curl1.HTTPVersion.String())
	fmt.Printf("HTTP/1.1: %s\n", curl2.HTTPVersion.String())
	fmt.Printf("HTTP/2: %s\n", curl3.HTTPVersion.String())
}
```

### Example 7: Debug and Verbose Output

Get detailed debugging information:

```go
package main

import (
	"fmt"
	"log"
	
	"github.com/474420502/gcurl"
)

func main() {
	curl, err := gcurl.Parse(`curl -v -H "Authorization: Bearer token123" https://api.example.com/data`)
	if err != nil {
		log.Fatal(err)
	}
	
	// Debug output (like curl -v)
	fmt.Println("=== Debug Output ===")
	fmt.Println(curl.Debug())
	
	// Verbose info (simulated curl verbose output)
	fmt.Println("\n=== Verbose Output ===")
	fmt.Println(curl.VerboseInfo())
	
	// Summary
	fmt.Println("\n=== Summary ===")
	fmt.Println(curl.Summary())
}
```

## 🔧 Advanced Usage

### Session Reuse

Sessions can be reused across multiple requests to maintain cookies and connection pooling:

```go
curl1, _ := gcurl.Parse(`curl "http://httpbin.org/cookies/set/session/123"`)
curl2, _ := gcurl.Parse(`curl "http://httpbin.org/cookies"`)

// Create a shared session
ses := curl1.CreateSession()

// Both requests will share the same session (and cookies)
resp1, _ := curl1.CreateRequest(ses).Execute()
resp2, _ := curl2.CreateRequest(ses).Execute()
```

### Direct Execution

For simple one-off requests:

```go
curl, _ := gcurl.Parse(`curl "https://httpbin.org/get"`)
resp, err := curl.Request().Execute() // Direct execution with auto-generated session
```

### Custom Session Configuration

```go
package main

import (
	"time"
	"net/http"
	"github.com/474420502/gcurl"
)

func main() {
	curl, _ := gcurl.Parse(`curl "https://httpbin.org/delay/2"`)
	
	// Create and customize session
	ses := curl.CreateSession()
	
	// Add custom headers
	customHeaders := make(http.Header)
	customHeaders.Set("X-Custom-Header", "MyValue")
	ses.SetHeader(customHeaders)
	
	// Set timeout
	ses.Config().SetTimeout(5 * time.Second)
	
	// Use customized session
	resp, err := curl.CreateRequest(ses).Execute()
	// ... handle response
}
``` 
## 📊 API Reference

### Main Functions

- `gcurl.Parse(curlCommand string) (*CURL, error)` - Parse a cURL command string
- `gcurl.ParseBash(curlCommand string) (*CURL, error)` - Parse specifically as Bash format  
- `gcurl.ParseCmd(curlCommand string) (*CURL, error)` - Parse specifically as Windows Cmd format

### CURL Methods  

- `CreateSession() *requests.Session` - Create a new session with parsed settings
- `CreateRequest(ses *requests.Session) *requests.Request` - Create request with optional session  
- `Request() *requests.Request` - Create request with auto-generated session
- `Debug() string` - Get detailed debug information (like `curl -v`)
- `VerboseInfo() string` - Get verbose output simulation
- `Summary() string` - Get brief request summary
- ~~`CreateTemporary(ses *requests.Session) *requests.Temporary`~~ - **Deprecated**: Use `CreateRequest()` instead
- ~~`Temporary() *requests.Temporary`~~ - **Deprecated**: Use `Request()` instead  

### Response Methods

- `resp.GetStatusCode() int` - Get HTTP status code
- `resp.Content() []byte` - Get response body as bytes
- `resp.ContentString() string` - Get response body as string

## ✅ Supported cURL Options

| Category | cURL Option | Description | Status |
|----------|-------------|-------------|--------|
| **HTTP Methods** | `-X, --request` | HTTP method (GET, POST, etc.) | ✅ |
| **Headers** | `-H, --header` | Custom headers | ✅ |
| **Data** | `-d, --data` | POST data | ✅ |
| | `--data-urlencode` | URL encoded data | ✅ |
| | `-F, --form` | Multipart form data | ✅ |
| **Authentication** | `-u, --user` | Basic authentication | ✅ |
| | `--digest` | Digest authentication | ✅ |
| **Cookies** | `-b, --cookie` | Send cookies | ✅ |
| | `-c, --cookie-jar` | Save cookies | ✅ |
| **File Output** | `-o, --output` | Output to file | ✅ |
| | `-O, --remote-name` | Use remote filename | ✅ |
| | `--output-dir` | Output directory | ✅ |
| | `--create-dirs` | Create directories | ✅ |
| | `-C, --continue-at` | Resume download | ✅ |
| | `--remove-on-error` | Remove file on error | ✅ |
| **HTTP Protocol** | `--http1.0` | Force HTTP/1.0 | ✅ |
| | `--http1.1` | Force HTTP/1.1 | ✅ |
| | `--http2` | Force HTTP/2 | ✅ |
| **Redirects** | `-L, --location` | Follow redirects | ✅ |
| | `--max-redirs` | Maximum redirects | ✅ |
| **Timeouts** | `--connect-timeout` | Connection timeout | ✅ |
| | `--max-time` | Maximum time | ✅ |
| **Proxy** | `--proxy` | Proxy settings | ✅ |
| | `--proxy-user` | Proxy authentication | ✅ |
| **SSL/TLS** | `-k, --insecure` | Skip SSL verification | ✅ |
| | `--cacert` | CA certificate | ✅ |
| | `--cert` | Client certificate | ✅ |
| | `--key` | Client private key | ✅ |
| **User Agent** | `-A, --user-agent` | User agent string | ✅ |
| **Debugging** | `-v, --verbose` | Verbose output | ✅ |
| | `-i, --include` | Include headers | ✅ |
| | `-I, --head` | HEAD request | ✅ |
| **Compression** | `--compressed` | Accept compression | ✅ |

## 🔍 Debug and Troubleshooting

### Debug Output

Use the `Debug()` method to get detailed information about the parsed cURL command:

```go
curl, _ := gcurl.Parse(`curl -v -X POST -H "Content-Type: application/json" -d '{"test":true}' https://api.example.com`)
fmt.Println(curl.Debug())
```

**Output:**
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

### Common Issues

1. **URL not found**: Make sure the URL is properly quoted
2. **Headers not parsed**: Check header format `-H "Key: Value"`
3. **File upload fails**: Ensure file exists and use correct syntax `-F "field=@file"`
4. **Authentication not working**: Verify credentials format `-u "user:pass"`

## 🎯 Performance Tips

1. **Reuse sessions** for multiple requests to the same host
2. **Use connection pooling** via shared sessions
3. **Set appropriate timeouts** to avoid hanging requests
4. **Enable compression** with `--compressed` for large responses

## 🧪 Testing

Run the test suite:

```bash
# Run all tests
go test -v

# Run tests with coverage
go test -v -cover

# Run specific test
go test -v -run TestFileOutput
```

Current test coverage: **79.1%** with 100+ test cases covering all major functionality.

## 🤝 Contributing

We welcome contributions! Here's how you can help:

1. **Report bugs** - Open an issue with reproduction steps
2. **Suggest features** - Propose new cURL options to support
3. **Submit pull requests** - Add tests for any new functionality
4. **Improve docs** - Help make the documentation better

### Development Setup

```bash
git clone https://github.com/474420502/gcurl.git
cd gcurl
go mod tidy
go test -v
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [requests](https://github.com/474420502/requests) - The underlying HTTP library
- [cURL](https://curl.se/) - For the excellent command-line tool that inspired this project
- All contributors who have helped improve this library

## 🔗 Related Projects

- [mholt/curl-to-go](https://github.com/mholt/curl-to-go) - Convert cURL to Go code (online tool)
- [sj26/curl-to-go](https://github.com/sj26/curl-to-go) - Another cURL to Go converter
- [requests](https://github.com/474420502/requests) - The HTTP library used by gcurl

---

⭐ **Star this project** if it helps you convert cURL commands to Go code!

