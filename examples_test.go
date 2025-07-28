package gcurl

import (
	"log"
	"strings"
	"testing"
)

// TestExample1BasicGET 测试例子1：基本GET请求
func TestExample1BasicGET(t *testing.T) {
	surl := `http://httpbin.org/get -H 'Connection: keep-alive' -H 'Accept-Encoding: gzip, deflate' -H 'Accept-Language: zh-CN,zh;q=0.9'`
	curl, err := Parse(surl)
	if err != nil {
		t.Fatalf("Failed to parse curl: %v", err)
	}

	ses := curl.CreateSession()
	tp := curl.CreateTemporary(ses)

	// 验证headers是否正确设置
	headers := ses.GetHeader()
	expectedHeaders := map[string]string{
		"Connection":      "keep-alive",
		"Accept-Encoding": "gzip, deflate",
		"Accept-Language": "zh-CN,zh;q=0.9",
	}

	for key, expectedValue := range expectedHeaders {
		if actualValue := headers.Get(key); actualValue != expectedValue {
			t.Errorf("Header %s = %s, want %s", key, actualValue, expectedValue)
		}
	}

	// 执行请求
	resp, err := tp.Execute()
	if err != nil {
		t.Fatalf("Failed to execute request: %v", err)
	}

	// 验证响应
	content := string(resp.Content())
	if !strings.Contains(content, "httpbin.org") {
		t.Error("Response does not contain expected content")
	}

	log.Println("Example 1 headers:", ses.GetHeader())
	log.Println("Example 1 response:", content)
}

// TestExample2WithCookies 测试例子2：带Cookie的GET请求
func TestExample2WithCookies(t *testing.T) {
	scurl := `curl 'http://httpbin.org/get' 
	--connect-timeout 1 
	-H 'authority: appgrowing.cn'
	-H 'accept-encoding: gzip, deflate, br' 
	-H 'accept-language: zh' 
	-H 'cookie: _ga=GA1.2.1371058419.1533104518; _gid=GA1.2.896241740.1543307916; _gat_gtag_UA_4002880_19=1' 
	-H 'if-none-match: W/"5bf7a0a9-ca6"' 
	-H 'if-modified-since: Fri, 23 Nov 2018 06:39:37 GMT'`

	curl, err := Parse(scurl)
	if err != nil {
		t.Fatalf("Failed to parse curl: %v", err)
	}

	ses := curl.CreateSession()
	wf := curl.CreateTemporary(ses)

	// 验证cookies
	cookies := ses.GetCookies(wf.ParsedURL)
	if len(cookies) == 0 {
		t.Error("Expected cookies to be set")
	}

	expectedCookies := []string{"_ga", "_gid", "_gat_gtag_UA_4002880_19"}
	cookieMap := make(map[string]bool)
	for _, cookie := range cookies {
		cookieMap[cookie.Name] = true
	}

	for _, expectedCookie := range expectedCookies {
		if !cookieMap[expectedCookie] {
			t.Errorf("Expected cookie %s not found", expectedCookie)
		}
	}

	// 执行请求
	resp, err := wf.Execute()
	if err != nil {
		t.Fatalf("Failed to execute request: %v", err)
	}

	log.Println("Example 2 cookies:", ses.GetCookies(wf.ParsedURL))
	log.Println("Example 2 response:", string(resp.Content()))
}

// TestExample3PathParams 测试例子3：路径参数
func TestExample3PathParams(t *testing.T) {
	c, err := Parse(`curl -X GET "http://httpbin.org/anything/1" -H "accept: application/json"`)
	if err != nil {
		t.Fatalf("Failed to parse curl: %v", err)
	}

	tp := c.Temporary()

	// 验证PathParam功能
	pp := tp.PathParam(`anything/(\d+)`)
	if pp == nil {
		t.Skip("PathParam method not available in current requests version")
		return
	}

	pp.IntSet(100) // Set Param.
	resp, err := tp.Execute()
	if err != nil {
		t.Fatalf("Failed to execute request: %v", err)
	}

	content := string(resp.Content())
	if !strings.Contains(content, "anything/100") {
		t.Error("Path parameter replacement did not work correctly")
	}

	log.Println("Example 3 response:", content)
}

// TestExample4PostWithData 测试例子4：POST请求带数据
func TestExample4PostWithData(t *testing.T) {
	scurl := `curl -X POST "http://httpbin.org/post" -H "Content-Type: application/json" -d '{"name":"test","age":25}'`

	curl, err := Parse(scurl)
	if err != nil {
		t.Fatalf("Failed to parse curl: %v", err)
	}

	// 验证方法和内容类型
	if curl.Method != "POST" {
		t.Errorf("Method = %s, want POST", curl.Method)
	}

	if curl.ContentType != "application/json" {
		t.Errorf("ContentType = %s, want application/json", curl.ContentType)
	}

	// 执行请求
	resp, err := curl.Temporary().Execute()
	if err != nil {
		t.Fatalf("Failed to execute request: %v", err)
	}

	content := string(resp.Content())
	if !strings.Contains(content, "test") || !strings.Contains(content, "25") {
		t.Error("POST data not found in response")
	}

	log.Println("Example 4 response:", content)
}

// TestExample5FormData 测试例子5：表单数据上传
func TestExample5FormData(t *testing.T) {
	scurl := `curl -X POST "http://httpbin.org/post" -F "name=john" -F "age=30" -F "email=john@example.com"`

	curl, err := Parse(scurl)
	if err != nil {
		t.Fatalf("Failed to parse curl: %v", err)
	}

	// 验证是否为multipart
	if curl.Body == nil || curl.Body.Type != "multipart" {
		t.Errorf("Expected multipart body, got %v", curl.Body)
	}

	// 执行请求
	resp, err := curl.Temporary().Execute()
	if err != nil {
		t.Fatalf("Failed to execute request: %v", err)
	}

	content := string(resp.Content())
	expectedValues := []string{"john", "30", "john@example.com"}
	for _, value := range expectedValues {
		if !strings.Contains(content, value) {
			t.Errorf("Form value %s not found in response", value)
		}
	}

	log.Println("Example 5 response:", content)
}
