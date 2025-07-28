package gcurl

import (
	"testing"
)

// TestNewOptions 测试新添加的curl选项
func TestNewOptions(t *testing.T) {
	tests := []struct {
		name        string
		curlCmd     string
		expectError bool
		checkFunc   func(*CURL) bool
		description string
	}{
		{
			name:        "Max redirects",
			curlCmd:     `curl --max-redirs 5 https://httpbin.org/redirect/3`,
			expectError: false,
			checkFunc: func(c *CURL) bool {
				return c.MaxRedirs == 5
			},
			description: "测试最大重定向次数设置",
		},
		{
			name:        "CA certificate",
			curlCmd:     `curl --cacert /dev/null https://httpbin.org/get`,
			expectError: false,
			checkFunc: func(c *CURL) bool {
				return c.CACert == "/dev/null"
			},
			description: "测试自定义CA证书设置",
		},
		{
			name:        "Client certificate",
			curlCmd:     `curl --cert /dev/null https://httpbin.org/get`,
			expectError: false,
			checkFunc: func(c *CURL) bool {
				return c.ClientCert == "/dev/null"
			},
			description: "测试客户端证书设置",
		},
		{
			name:        "Client key",
			curlCmd:     `curl --key /dev/null https://httpbin.org/get`,
			expectError: false,
			checkFunc: func(c *CURL) bool {
				return c.ClientKey == "/dev/null"
			},
			description: "测试客户端私钥设置",
		},
		{
			name:        "HTTP/2",
			curlCmd:     `curl --http2 https://httpbin.org/get`,
			expectError: false,
			checkFunc: func(c *CURL) bool {
				return c.HTTP2 == true
			},
			description: "测试强制HTTP/2设置",
		},
		{
			name:        "Invalid max-redirs",
			curlCmd:     `curl --max-redirs invalid https://httpbin.org/get`,
			expectError: true,
			checkFunc:   nil,
			description: "测试无效的最大重定向次数",
		},
		{
			name:        "Negative max-redirs",
			curlCmd:     `curl --max-redirs -1 https://httpbin.org/get`,
			expectError: true,
			checkFunc:   nil,
			description: "测试负数最大重定向次数",
		},
		{
			name:        "Nonexistent CA cert",
			curlCmd:     `curl --cacert /nonexistent/path/ca.crt https://httpbin.org/get`,
			expectError: true,
			checkFunc:   nil,
			description: "测试不存在的CA证书文件",
		},
		{
			name:        "Nonexistent client cert",
			curlCmd:     `curl --cert /nonexistent/path/client.crt https://httpbin.org/get`,
			expectError: true,
			checkFunc:   nil,
			description: "测试不存在的客户端证书文件",
		},
		{
			name:        "Nonexistent client key",
			curlCmd:     `curl --key /nonexistent/path/client.key https://httpbin.org/get`,
			expectError: true,
			checkFunc:   nil,
			description: "测试不存在的客户端私钥文件",
		},
		{
			name:        "Combined SSL options",
			curlCmd:     `curl --cacert /dev/null --cert /dev/null --key /dev/null --http2 https://httpbin.org/get`,
			expectError: false,
			checkFunc: func(c *CURL) bool {
				return c.CACert == "/dev/null" && c.ClientCert == "/dev/null" && c.ClientKey == "/dev/null" && c.HTTP2 == true
			},
			description: "测试组合SSL选项",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Testing: %s - %s", tt.name, tt.description)

			cu, err := Parse(tt.curlCmd)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for %s, but got none", tt.name)
				} else {
					t.Logf("✓ Expected error: %v", err)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error for %s: %v", tt.name, err)
				return
			}

			// 验证解析结果
			if cu.ParsedURL == nil {
				t.Errorf("URL not parsed for %s", tt.name)
				return
			}

			if tt.checkFunc != nil && !tt.checkFunc(cu) {
				t.Errorf("Check function failed for %s", tt.name)
				return
			}

			t.Logf("✓ Successfully parsed: %s", cu.ParsedURL.String())
			t.Logf("✓ MaxRedirs: %d, HTTP2: %v", cu.MaxRedirs, cu.HTTP2)
			t.Logf("✓ CACert: %s, ClientCert: %s, ClientKey: %s", cu.CACert, cu.ClientCert, cu.ClientKey)
		})
	}
}

// TestNewOptionsIntegration 集成测试
func TestNewOptionsIntegration(t *testing.T) {
	complexCmd := `curl --max-redirs 3 --http2 --cacert /dev/null -H "Accept: application/json" https://httpbin.org/get`

	cu, err := Parse(complexCmd)
	if err != nil {
		t.Fatalf("Failed to parse complex command: %v", err)
	}

	// 验证所有选项都被正确设置
	if cu.MaxRedirs != 3 {
		t.Errorf("Expected MaxRedirs=3, got %d", cu.MaxRedirs)
	}

	if !cu.HTTP2 {
		t.Error("Expected HTTP2=true")
	}

	if cu.CACert != "/dev/null" {
		t.Errorf("Expected CACert=/dev/null, got %s", cu.CACert)
	}

	if cu.Header.Get("Accept") != "application/json" {
		t.Error("Expected Accept header not found")
	}

	t.Log("✓ Complex integration test passed")
}
