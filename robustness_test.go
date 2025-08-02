package gcurl

import (
	"fmt"
	"strings"
	"testing"
)

// TestCrossPlatformQuoteHandling 测试跨平台引号处理的健壮性
func TestCrossPlatformQuoteHandling(t *testing.T) {
	testCases := []struct {
		name        string
		command     string
		platform    string
		shouldPass  bool
		description string
	}{
		{
			name:        "Chrome Linux single quotes",
			command:     `curl 'https://example.com' -H 'sec-ch-ua: "Chromium";v="124", "Google Chrome";v="124"'`,
			platform:    "linux",
			shouldPass:  true,
			description: "典型的Chrome Linux导出格式",
		},
		{
			name:        "Chrome Windows cmd format",
			command:     `curl "https://example.com" -H ^"sec-ch-ua: ^\^"Chromium^\^";v=^\^"124^\^"^"`,
			platform:    "windows",
			shouldPass:  true,
			description: "Chrome Windows cmd格式转义",
		},
		{
			name:        "Chrome mixed quotes",
			command:     `curl 'https://example.com' -H "sec-ch-ua: \"Chromium\";v=\"124\""`,
			platform:    "mixed",
			shouldPass:  true,
			description: "混合单双引号",
		},
		{
			name:        "Complex nested quotes",
			command:     `curl 'https://api.example.com' -H 'X-Custom: "value with \"nested\" quotes"'`,
			platform:    "linux",
			shouldPass:  true,
			description: "嵌套引号处理",
		},
		{
			name:        "ANSI-C quoting",
			command:     `curl $'https://example.com\ntest' -d $'data\r\nwith\tspecial'`,
			platform:    "bash",
			shouldPass:  true,
			description: "ANSI-C引用格式",
		},
		{
			name:        "Edge case empty quotes",
			command:     `curl 'https://example.com' -H '' -H ""`,
			platform:    "any",
			shouldPass:  false, // Empty header values should be rejected
			description: "空引号边界情况",
		},
		{
			name:        "Real Chrome cookie header",
			command:     `curl 'https://example.com' -H 'cookie: sessionid=abc123; csrftoken=xyz789'`,
			platform:    "chrome",
			shouldPass:  true,
			description: "真实Chrome cookie头",
		},
		{
			name:        "Malformed unclosed quote",
			command:     `curl 'https://example.com -H 'X-Test: value'`,
			platform:    "any",
			shouldPass:  false,
			description: "未闭合引号错误",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("Testing %s: %s", tc.platform, tc.description)
			t.Logf("Command: %s", tc.command)

			var curl *CURL
			var err error

			// 根据平台选择解析方法
			switch tc.platform {
			case "windows":
				curl, err = ParseCmd(tc.command)
			default:
				curl, err = Parse(tc.command)
			}

			if tc.shouldPass {
				if err != nil {
					t.Errorf("Expected success but got error: %v", err)
					return
				}
				if curl == nil {
					t.Error("Expected valid CURL object but got nil")
					return
				}

				// 验证URL解析正确
				if curl.ParsedURL == nil {
					t.Error("URL should be parsed successfully")
					return
				}

				t.Logf("✓ URL: %s", curl.ParsedURL.String())
				t.Logf("✓ Headers count: %d", len(curl.Header))

				// 记录头部用于调试
				for key, values := range curl.Header {
					for _, value := range values {
						t.Logf("  Header %s: %s", key, value)
					}
				}
			} else {
				if err == nil {
					t.Error("Expected error but parsing succeeded")
				} else {
					t.Logf("✓ Expected error: %v", err)
				}
			}
		})
	}
}

// TestRealWorldChromeCommands 测试真实世界的Chrome导出命令
func TestRealWorldChromeCommands(t *testing.T) {
	// 这些是从实际Chrome "Copy as cURL"功能获得的命令
	realCommands := []string{
		// Linux Chrome
		`curl 'https://httpbin.org/get' -H 'accept: application/json' -H 'accept-language: en-US,en;q=0.9' -H 'sec-ch-ua: "Google Chrome";v="119", "Chromium";v="119", "Not?A_Brand";v="24"' -H 'sec-ch-ua-mobile: ?0' -H 'sec-ch-ua-platform: "Linux"'`,

		// Windows Chrome (经过cmdformat2bash转换)
		`curl "https://httpbin.org/post" -H "accept: application/json" -H "content-type: application/json" --data-raw "{\"test\":\"value\"}"`,

		// 复杂的表单数据
		`curl 'https://httpbin.org/post' -H 'content-type: multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW' --data-raw $'------WebKitFormBoundary7MA4YWxkTrZu0gW\r\nContent-Disposition: form-data; name="field1"\r\n\r\nvalue1\r\n------WebKitFormBoundary7MA4YWxkTrZu0gW--\r\n'`,
	}

	for i, cmd := range realCommands {
		t.Run(fmt.Sprintf("RealWorld_%d", i+1), func(t *testing.T) {
			t.Logf("Testing real-world command: %s", cmd[:min(100, len(cmd))]+"...")

			curl, err := Parse(cmd)
			if err != nil {
				t.Errorf("Failed to parse real-world command: %v", err)
				return
			}

			// 基本验证
			if curl.ParsedURL == nil {
				t.Error("URL should be parsed")
				return
			}

			t.Logf("✓ Successfully parsed URL: %s", curl.ParsedURL.String())
			t.Logf("✓ Method: %s", curl.Method)
			t.Logf("✓ Headers: %d", len(curl.Header))

			// 验证关键头部
			expectedHeaders := []string{"Accept", "Sec-Ch-Ua", "Content-Type"}
			for _, header := range expectedHeaders {
				if values := curl.Header.Get(header); values != "" {
					t.Logf("✓ Found %s: %s", header, values)
				}
			}
		})
	}
}

// TestQuoteEscapeEdgeCases 测试引号转义的边界情况
func TestQuoteEscapeEdgeCases(t *testing.T) {
	edgeCases := []struct {
		name    string
		command string
		expect  string // 期望的头部值
	}{
		{
			name:    "Escaped quotes in single quotes",
			command: `curl 'https://example.com' -H 'X-Test: "value with \"escaped\" quotes"'`,
			expect:  `"value with \"escaped\" quotes"`, // 在单引号内，反斜杠转义被保留
		},
		{
			name:    "Escaped quotes in double quotes",
			command: `curl 'https://example.com' -H "X-Test: \"value with \\\"escaped\\\" quotes\""`,
			expect:  `"value with \"escaped\" quotes"`, // 在双引号内，转义被处理
		},
		{
			name:    "Backslash in single quotes",
			command: `curl 'https://example.com' -H 'X-Path: C:\\Windows\\System32'`,
			expect:  `C:\\Windows\\System32`, // 在单引号内，反斜杠被保留为字面值
		},
		{
			name:    "Backslash in double quotes",
			command: `curl 'https://example.com' -H "X-Path: C:\\\\Windows\\\\System32"`,
			expect:  `C:\\Windows\\System32`, // 在双引号内，\\\\变成\\
		},
		{
			name:    "Unicode characters",
			command: `curl 'https://example.com' -H 'X-Unicode: 测试中文字符'`,
			expect:  `测试中文字符`,
		},
	}

	for _, tc := range edgeCases {
		t.Run(tc.name, func(t *testing.T) {
			curl, err := Parse(tc.command)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}

			// 查找对应的头部值
			var headerValue string
			for headerName := range curl.Header {
				if strings.Contains(strings.ToLower(headerName), "test") {
					headerValue = curl.Header.Get(headerName)
					break
				} else if strings.Contains(strings.ToLower(headerName), "path") {
					headerValue = curl.Header.Get(headerName)
					break
				} else if strings.Contains(strings.ToLower(headerName), "unicode") {
					headerValue = curl.Header.Get(headerName)
					break
				}
			}

			if headerValue != tc.expect {
				t.Errorf("Expected header value '%s', got '%s'", tc.expect, headerValue)
			} else {
				t.Logf("✓ Correct header value: %s", headerValue)
			}
		})
	}
}

// min helper function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
