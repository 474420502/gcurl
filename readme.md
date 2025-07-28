# Parse curl To golang requests

* Based on [requests](https://github.com/474420502/requests) library
* Easy to transform cURL bash commands to Go code
* Inherits all cURL functionality and adds Go's flexibility for configuration, cookies, headers, and URL handling
* Supports both Bash and Windows Cmd cURL formats

## Features

- 🚀 Complete cURL command parsing
- 🔧 Support for all major cURL options (headers, cookies, data, authentication, etc.)
- 🌐 HTTP/HTTPS requests with SSL/TLS configuration  
- 📝 Form data and file uploads
- 🍪 Cookie management
- 🔄 Redirect handling
- ⏱️ Timeout and connection settings
- 🔐 Proxy support (HTTP/HTTPS/SOCKS5)
- 🏷️ Path parameter replacement

# Installation

```bash
go get github.com/474420502/gcurl
```

# Examples

## Example 1: Basic GET Request with Headers

This example demonstrates how to parse a cURL command for a GET request with custom headers, create a session, and execute the request.

```go
package main

import (
	"fmt"
	"log"
	
	"github.com/474420502/gcurl"
)

func main() {
	// Parse cURL command
	surl := `http://httpbin.org/get -H 'Connection: keep-alive' -H 'Accept-Encoding: gzip, deflate' -H 'Accept-Language: zh-CN,zh;q=0.9'`
	curl, err := gcurl.Parse(surl)
	if err != nil {
		log.Fatal(err)
	}
	
	// Create session and temporary request
	ses := curl.CreateSession()
	tp := curl.CreateTemporary(ses)
	
	// Check headers
	fmt.Println("Headers:", ses.GetHeader())
	// Output: map[Accept-Encoding:[gzip, deflate] Accept-Language:[zh-CN,zh;q=0.9] Connection:[keep-alive]]
	
	// Execute request
	resp, err := tp.Execute()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Response:", string(resp.Content()))
	// Response will contain:
	// {
	//   "headers": {
	//     "Accept-Encoding": "gzip, deflate",
	//     "Accept-Language": "zh-CN,zh;q=0.9", 
	//     "Connection": "keep-alive",
	//     "Host": "httpbin.org",
	//     "User-Agent": "Go-http-client/1.1"
	//   },
	//   "origin": "your.ip.address",
	//   "url": "http://httpbin.org/get"
	// }
}
```

## Example 2: GET Request with Cookies

This example demonstrates how to parse a cURL command with cookies and custom headers, and verify that cookies are properly handled.

```go
package main

import (
	"fmt"
	"log"
	
	"github.com/474420502/gcurl"
)

func main() {
	scurl := `curl 'http://httpbin.org/get' 
		--connect-timeout 1 
		-H 'authority: appgrowing.cn'
		-H 'accept-encoding: gzip, deflate, br' 
		-H 'accept-language: zh' 
		-H 'cookie: _ga=GA1.2.1371058419.1533104518; _gid=GA1.2.896241740.1543307916; _gat_gtag_UA_4002880_19=1' 
		-H 'if-none-match: W/"5bf7a0a9-ca6"' 
		-H 'if-modified-since: Fri, 23 Nov 2018 06:39:37 GMT'`
	
	curl, err := gcurl.Parse(scurl)
	if err != nil {
		log.Fatal(err)
	}
	
	ses := curl.CreateSession()
	wf := curl.CreateTemporary(ses)
	
	// Check cookies were parsed correctly
	cookies := ses.GetCookies(wf.ParsedURL)
	fmt.Println("Cookies:", cookies)
	// Output: [_ga=GA1.2.1371058419.1533104518 _gid=GA1.2.896241740.1543307916 _gat_gtag_UA_4002880_19=1]
	
	resp, err := wf.Execute()
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Println("Response:", string(resp.Content()))
	// Response will show that cookies were sent in the request headers
}
```

## Example 3: POST Request with JSON Data

This example shows how to handle POST requests with JSON data.

```go
package main

import (
	"fmt"
	"log"
	
	"github.com/474420502/gcurl"
)

func main() {
	scurl := `curl -X POST "http://httpbin.org/post" -H "Content-Type: application/json" -d '{"name":"test","age":25}'`
	
	curl, err := gcurl.Parse(scurl)
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("Method: %s\n", curl.Method)          // POST
	fmt.Printf("Content-Type: %s\n", curl.ContentType) // application/json
	
	resp, err := curl.Temporary().Execute()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Response:", string(resp.Content()))
	// Response will contain the JSON data in the "json" field and "data" field
}
```

## Example 4: Form Data Upload

This example demonstrates how to handle multipart form data uploads.

```go
package main

import (
	"fmt"
	"log"
	
	"github.com/474420502/gcurl"
)

func main() {
	scurl := `curl -X POST "http://httpbin.org/post" -F "name=john" -F "age=30" -F "email=john@example.com"`
	
	curl, err := gcurl.Parse(scurl)
	if err != nil {
		log.Fatal(err)
	}
	
	// Form data is automatically handled as multipart
	fmt.Printf("Body type: %s\n", curl.Body.Type) // multipart
	
	resp, err := curl.Temporary().Execute()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Response:", string(resp.Content()))
	// Response will show form data in the "form" field
}
```

## Example 6: File Upload

This example shows how to upload files using form data.

```go
package main

import (
	"fmt"
	"log"
	"os"
	
	"github.com/474420502/gcurl"
)

func main() {
	// First create a test file
	testFile := "/tmp/test.txt"
	err := os.WriteFile(testFile, []byte("test file content"), 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(testFile)
	
	scurl := fmt.Sprintf(`curl -X POST "http://httpbin.org/post" -F "file=@%s" -F "description=test file"`, testFile)
	
	curl, err := gcurl.Parse(scurl)
	if err != nil {
		log.Fatal(err)
	}
	
	resp, err := curl.Temporary().Execute()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("File upload response:", string(resp.Content()))
}
```

## Example 7: Authentication and HTTPS

This example demonstrates basic authentication and HTTPS requests.

```go
package main

import (
	"fmt"
	"log"
	
	"github.com/474420502/gcurl"
)

func main() {
	// Basic authentication example
	scurl := `curl -u "user:password" "http://httpbin.org/basic-auth/user/password"`
	
	curl, err := gcurl.Parse(scurl)
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("Auth configured: %v\n", curl.Auth != nil)
	
	resp, err := curl.Temporary().Execute()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Status: %d\n", resp.GetStatusCode())
	fmt.Println("Auth response:", string(resp.Content()))
}
```

## Example 8: Custom Session Configuration

This example shows how to modify session settings before executing requests.

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
	curl, err := gcurl.Parse(`curl "http://httpbin.org/delay/2"`)
	if err != nil {
		log.Fatal(err)
	}
	
	// Create session and customize it
	ses := curl.CreateSession()
	
	// Add custom headers to the session
	customHeaders := make(http.Header)
	customHeaders.Set("X-Custom-Header", "MyValue")
	customHeaders.Set("User-Agent", "MyApp/1.0")
	ses.SetHeader(customHeaders)
	
	// Set timeout
	ses.Config().SetTimeout(5) // 5 seconds
	
	// Create temporary request with customized session
	tp := curl.CreateTemporary(ses)
	
	start := time.Now()
	resp, err := tp.Execute()
	if err != nil {
		log.Fatal(err)
	}
	duration := time.Since(start)
	
	fmt.Printf("Request took: %v\n", duration)
	fmt.Printf("Status: %d\n", resp.GetStatusCode())
	fmt.Println("Response:", string(resp.Content()))
}
```

## Example 9: Error Handling and Validation

This example shows proper error handling and response validation.

```go
package main

import (
	"fmt"
	"log"
	"strings"
	
	"github.com/474420502/gcurl"
)

func main() {
	scurl := `curl -X POST "http://httpbin.org/status/404" -d "test data"`
	
	curl, err := gcurl.Parse(scurl)
	if err != nil {
		log.Fatalf("Failed to parse cURL command: %v", err)
	}
	
	// Validate the parsed URL
	if curl.ParsedURL == nil {
		log.Fatal("Invalid URL in cURL command")
	}
	
	fmt.Printf("Parsed URL: %s\n", curl.ParsedURL.String())
	fmt.Printf("Method: %s\n", curl.Method)
	
	resp, err := curl.Temporary().Execute()
	if err != nil {
		log.Fatalf("Request failed: %v", err)
	}
	
	// Check response status
	fmt.Printf("Status Code: %d\n", resp.GetStatusCode())
	
	// Handle different status codes
	switch {
	case resp.GetStatusCode() >= 200 && resp.GetStatusCode() < 300:
		fmt.Println("✅ Success!")
	case resp.GetStatusCode() >= 400 && resp.GetStatusCode() < 500:
		fmt.Println("❌ Client error")
	case resp.GetStatusCode() >= 500:
		fmt.Println("💥 Server error")
	}
	
	// Check if response contains expected content
	content := string(resp.Content())
	if strings.Contains(content, "httpbin") {
		fmt.Println("✅ Response from httpbin confirmed")
	}
	
	fmt.Println("Response body:", content)
}
```

## Advanced Usage

### Reusing Sessions

Sessions can be reused across multiple requests to maintain cookies and connection pooling:

```go
curl1, _ := gcurl.Parse(`curl "http://httpbin.org/cookies/set/session/123"`)
curl2, _ := gcurl.Parse(`curl "http://httpbin.org/cookies"`)

// Create a shared session
ses := curl1.CreateSession()

// Both requests will share the same session (and cookies)
resp1, _ := curl1.CreateTemporary(ses).Execute()
resp2, _ := curl2.CreateTemporary(ses).Execute()
```

### Direct Execution

For simple one-off requests, you can execute directly:

```go
curl, _ := gcurl.Parse(`curl "http://httpbin.org/get"`)
resp, err := curl.Temporary().Execute() // Direct execution with auto-generated session
```

## API Reference

### Main Functions

- `gcurl.Parse(curlCommand string) (*CURL, error)` - Parse a cURL command string
- `gcurl.ParseBash(curlCommand string) (*CURL, error)` - Parse specifically as Bash format  
- `gcurl.ParseCmd(curlCommand string) (*CURL, error)` - Parse specifically as Windows Cmd format

### CURL Methods  

- `CreateSession() *requests.Session` - Create a new session with parsed settings
- `CreateTemporary(ses *requests.Session) *requests.Temporary` - Create request with optional session
- `Temporary() *requests.Temporary` - Create request with auto-generated session  

### Response Methods

- `resp.GetStatusCode() int` - Get HTTP status code
- `resp.Content() []byte` - Get response body as bytes
- `resp.ContentString() string` - Get response body as string

### Supported cURL Options

| cURL Option | Description | Supported |
|-------------|-------------|-----------|
| `-X, --request` | HTTP method | ✅ |
| `-H, --header` | Custom headers | ✅ |
| `-d, --data` | POST data | ✅ |
| `-F, --form` | Multipart form data | ✅ |
| `-u, --user` | Authentication | ✅ |
| `-b, --cookie` | Cookies | ✅ |
| `-L, --location` | Follow redirects | ✅ |
| `-k, --insecure` | Skip SSL verification | ✅ |
| `--connect-timeout` | Connection timeout | ✅ |
| `--max-time` | Maximum time | ✅ |
| `--proxy` | Proxy settings | ✅ |
| `-A, --user-agent` | User agent | ✅ |
| `--data-urlencode` | URL encoded data | ✅ |

## Contributing

We welcome contributions! Please feel free to submit issues and pull requests.

## License

This project is licensed under the MIT License.

