package gcurl

import (
	"testing"
)

func TestConnectTimeoutParsing(t *testing.T) {
	// 测试 --connect-timeout 选项的解析
	curlCmd := `curl --connect-timeout 5 http://httpbin.org/get`
	curl, err := ParseBash(curlCmd)
	if err != nil {
		t.Fatalf("Failed to parse curl command: %v", err)
	}

	// 验证 ConnectTimeout 字段是否被正确设置
	if curl.ConnectTimeout != 5 {
		t.Errorf("Expected ConnectTimeout to be 5, got %d", curl.ConnectTimeout)
	}

	// 验证其他字段也正常
	if curl.ParsedURL.String() != "http://httpbin.org/get" {
		t.Errorf("Expected URL to be http://httpbin.org/get, got %s", curl.ParsedURL.String())
	}

	if curl.Method != "GET" {
		t.Errorf("Expected Method to be GET, got %s", curl.Method)
	}
}

func TestConnectTimeoutWithOtherOptions(t *testing.T) {
	// 测试 --connect-timeout 与其他选项一起使用
	curlCmd := `curl --connect-timeout 10 --socks5 localhost:1080 -H "User-Agent: Test" http://httpbin.org/post`
	curl, err := ParseBash(curlCmd)
	if err != nil {
		t.Fatalf("Failed to parse curl command: %v", err)
	}

	// 验证各个选项都被正确解析
	if curl.ConnectTimeout != 10 {
		t.Errorf("Expected ConnectTimeout to be 10, got %d", curl.ConnectTimeout)
	}

	if curl.Proxy != "socks5://localhost:1080" {
		t.Errorf("Expected Proxy to be socks5://localhost:1080, got %s", curl.Proxy)
	}

	if curl.Header.Get("User-Agent") != "Test" {
		t.Errorf("Expected User-Agent header to be Test, got %s", curl.Header.Get("User-Agent"))
	}
}

func TestConnectTimeoutInvalidValue(t *testing.T) {
	// 测试无效的 --connect-timeout 值
	curlCmd := `curl --connect-timeout invalid http://httpbin.org/get`
	_, err := ParseBash(curlCmd)
	if err == nil {
		t.Fatal("Expected error for invalid connect-timeout value, but got none")
	}

	// 验证错误消息包含相关信息
	if !contains(err.Error(), "invalid value for --connect-timeout") {
		t.Errorf("Expected error message to contain 'invalid value for --connect-timeout', got: %s", err.Error())
	}
}

// 辅助函数检查字符串是否包含子字符串
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			indexOfSubstring(s, substr) >= 0)))
}

func indexOfSubstring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
