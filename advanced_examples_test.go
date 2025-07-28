package gcurl

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

// TestExample6FileUpload 测试文件上传例子
func TestExample6FileUpload(t *testing.T) {
	// Create test file
	testFile := "/tmp/gcurl_test.txt"
	err := os.WriteFile(testFile, []byte("test file content"), 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(testFile)

	scurl := fmt.Sprintf(`curl -X POST "http://httpbin.org/post" -F "file=@%s" -F "description=test file"`, testFile)

	curl, err := Parse(scurl)
	if err != nil {
		t.Fatalf("Failed to parse curl: %v", err)
	}

	// Verify it's multipart
	if curl.Body == nil || curl.Body.Type != "multipart" {
		t.Errorf("Expected multipart body, got %v", curl.Body)
	}

	resp, err := curl.Temporary().Execute()
	if err != nil {
		t.Fatalf("Failed to execute request: %v", err)
	}

	content := string(resp.Content())
	if !strings.Contains(content, "test file") {
		t.Error("File upload description not found in response")
	}
}

// TestExample7Authentication 测试认证例子
func TestExample7Authentication(t *testing.T) {
	scurl := `curl -u "user:password" "http://httpbin.org/basic-auth/user/password"`

	curl, err := Parse(scurl)
	if err != nil {
		t.Fatalf("Failed to parse curl: %v", err)
	}

	// Verify auth is configured
	if curl.Auth == nil {
		t.Error("Expected authentication to be configured")
	}

	resp, err := curl.Temporary().Execute()
	if err != nil {
		t.Fatalf("Failed to execute request: %v", err)
	}

	if resp.GetStatusCode() != 200 {
		t.Errorf("Expected status 200, got %d", resp.GetStatusCode())
	}
} // TestExample8CustomSession 测试自定义会话例子
func TestExample8CustomSession(t *testing.T) {
	curl, err := Parse(`curl "http://httpbin.org/headers"`)
	if err != nil {
		t.Fatalf("Failed to parse curl: %v", err)
	}

	// Create and customize session
	ses := curl.CreateSession()

	// Add custom headers properly
	customHeaders := make(http.Header)
	customHeaders.Set("X-Custom-Header", "MyValue")
	customHeaders.Set("User-Agent", "MyApp/1.0")
	ses.SetHeader(customHeaders)

	ses.Config().SetTimeout(5)

	tp := curl.CreateTemporary(ses)

	start := time.Now()
	resp, err := tp.Execute()
	if err != nil {
		t.Fatalf("Failed to execute request: %v", err)
	}
	duration := time.Since(start)

	// Verify request completed reasonably quickly
	if duration > 10*time.Second {
		t.Errorf("Request took too long: %v", duration)
	}

	content := string(resp.Content())
	if !strings.Contains(content, "X-Custom-Header") {
		t.Error("Custom header not found in response")
	}

	if !strings.Contains(content, "MyApp/1.0") {
		t.Error("Custom user agent not found in response")
	}
}

// TestExample9ErrorHandling 测试错误处理例子
func TestExample9ErrorHandling(t *testing.T) {
	scurl := `curl -X POST "http://httpbin.org/status/404" -d "test data"`

	curl, err := Parse(scurl)
	if err != nil {
		t.Fatalf("Failed to parse curl: %v", err)
	}

	// Validate parsed URL
	if curl.ParsedURL == nil {
		t.Fatal("ParsedURL should not be nil")
	}

	if curl.Method != "POST" {
		t.Errorf("Method = %s, want POST", curl.Method)
	}

	resp, err := curl.Temporary().Execute()
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	// Should get 404 status
	if resp.GetStatusCode() != 404 {
		t.Errorf("Status code = %d, want 404", resp.GetStatusCode())
	}

	// Check response content - 404 may have empty body, that's ok
	content := resp.ContentString()
	t.Logf("Response content: %s", content)
	t.Logf("Status: %d", resp.GetStatusCode())
} // TestSessionReuse 测试会话重用
func TestSessionReuse(t *testing.T) {
	curl1, err := Parse(`curl "http://httpbin.org/cookies/set/session/123"`)
	if err != nil {
		t.Fatal(err)
	}

	curl2, err := Parse(`curl "http://httpbin.org/cookies"`)
	if err != nil {
		t.Fatal(err)
	}

	// Create shared session
	ses := curl1.CreateSession()

	// First request sets cookie
	resp1, err := curl1.CreateTemporary(ses).Execute()
	if err != nil {
		t.Fatal(err)
	}

	// Second request should see the cookie
	resp2, err := curl2.CreateTemporary(ses).Execute()
	if err != nil {
		t.Fatal(err)
	}

	content1 := string(resp1.Content())
	content2 := string(resp2.Content())

	// Verify cookie was set and retrieved
	if !strings.Contains(content2, "session") {
		t.Error("Session cookie not found in second request")
	}

	t.Logf("First response: %s", content1)
	t.Logf("Second response: %s", content2)
}

// TestDirectExecution 测试直接执行
func TestDirectExecution(t *testing.T) {
	curl, err := Parse(`curl "http://httpbin.org/get"`)
	if err != nil {
		t.Fatal(err)
	}

	// Use temporary execution instead
	resp, err := curl.Temporary().Execute()
	if err != nil {
		t.Fatal(err)
	}

	if resp.GetStatusCode() != 200 {
		t.Errorf("Status code = %d, want 200", resp.GetStatusCode())
	}

	content := string(resp.Content())
	if !strings.Contains(content, "httpbin.org") {
		t.Error("Response should contain httpbin.org")
	}
}
