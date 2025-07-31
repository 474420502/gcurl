package gcurl

import (
	"testing"
)

func TestHTTPVersionControl(t *testing.T) {
	tests := []struct {
		name            string
		curlCommand     string
		expectedVersion HTTPVersion
		expectedHTTP2   bool
		wantErr         bool
	}{
		{
			name:            "HTTP/1.0 forced",
			curlCommand:     `curl --http1.0 https://httpbin.org/get`,
			expectedVersion: HTTPVersion10,
			expectedHTTP2:   false,
			wantErr:         false,
		},
		{
			name:            "HTTP/1.1 forced",
			curlCommand:     `curl --http1.1 https://httpbin.org/get`,
			expectedVersion: HTTPVersion11,
			expectedHTTP2:   false,
			wantErr:         false,
		},
		{
			name:            "HTTP/2 forced",
			curlCommand:     `curl --http2 https://httpbin.org/get`,
			expectedVersion: HTTPVersion2,
			expectedHTTP2:   true,
			wantErr:         false,
		},
		{
			name:            "Default auto version",
			curlCommand:     `curl https://httpbin.org/get`,
			expectedVersion: HTTPVersionAuto,
			expectedHTTP2:   false,
			wantErr:         false,
		},
		{
			name:            "Multiple protocol options - last wins",
			curlCommand:     `curl --http1.0 --http1.1 --http2 https://httpbin.org/get`,
			expectedVersion: HTTPVersion2,
			expectedHTTP2:   true,
			wantErr:         false,
		},
		{
			name:            "HTTP/1.1 overrides HTTP/2",
			curlCommand:     `curl --http2 --http1.1 https://httpbin.org/get`,
			expectedVersion: HTTPVersion11,
			expectedHTTP2:   false,
			wantErr:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := Parse(tt.curlCommand)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Parse() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Parse() unexpected error: %v", err)
				return
			}

			if c.HTTPVersion != tt.expectedVersion {
				t.Errorf("Parse() HTTPVersion = %v, want %v", c.HTTPVersion, tt.expectedVersion)
			}

			if c.HTTP2 != tt.expectedHTTP2 {
				t.Errorf("Parse() HTTP2 = %v, want %v", c.HTTP2, tt.expectedHTTP2)
			}
		})
	}
}

func TestHTTPVersionStringRepresentation(t *testing.T) {
	tests := []struct {
		version  HTTPVersion
		expected string
	}{
		{HTTPVersionAuto, "Auto"},
		{HTTPVersion10, "HTTP/1.0"},
		{HTTPVersion11, "HTTP/1.1"},
		{HTTPVersion2, "HTTP/2"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.version.String()
			if result != tt.expected {
				t.Errorf("HTTPVersion.String() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestHTTPVersionInDebugOutput(t *testing.T) {
	tests := []struct {
		name          string
		curlCommand   string
		expectInDebug string
	}{
		{
			name:          "HTTP/1.0 in debug",
			curlCommand:   `curl --http1.0 https://example.com`,
			expectInDebug: "HTTP Version: HTTP/1.0",
		},
		{
			name:          "HTTP/1.1 in debug",
			curlCommand:   `curl --http1.1 https://example.com`,
			expectInDebug: "HTTP Version: HTTP/1.1",
		},
		{
			name:          "HTTP/2 in debug",
			curlCommand:   `curl --http2 https://example.com`,
			expectInDebug: "HTTP Version: HTTP/2",
		},
		{
			name:          "Auto version in debug",
			curlCommand:   `curl https://example.com`,
			expectInDebug: "HTTP Version: Auto",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := Parse(tt.curlCommand)
			if err != nil {
				t.Errorf("Parse() error: %v", err)
				return
			}

			debugOutput := c.Debug()
			if !containsString(debugOutput, tt.expectInDebug) {
				t.Errorf("Debug() output missing expected string:\nWant: %s\nGot: %s",
					tt.expectInDebug, debugOutput)
			}
		})
	}
}

func TestHTTPVersionSessionConfiguration(t *testing.T) {
	tests := []struct {
		name        string
		curlCommand string
		version     HTTPVersion
	}{
		{
			name:        "HTTP/1.0 session config",
			curlCommand: `curl --http1.0 https://example.com`,
			version:     HTTPVersion10,
		},
		{
			name:        "HTTP/1.1 session config",
			curlCommand: `curl --http1.1 https://example.com`,
			version:     HTTPVersion11,
		},
		{
			name:        "HTTP/2 session config",
			curlCommand: `curl --http2 https://example.com`,
			version:     HTTPVersion2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := Parse(tt.curlCommand)
			if err != nil {
				t.Errorf("Parse() error: %v", err)
				return
			}

			// 验证会话创建不会崩溃
			session := c.CreateSession()
			if session == nil {
				t.Errorf("CreateSession() returned nil")
				return
			}

			// 验证配置方法被调用
			c.configureHTTPVersion(session)

			// 协议版本应该正确设置
			if c.HTTPVersion != tt.version {
				t.Errorf("HTTPVersion = %v, want %v", c.HTTPVersion, tt.version)
			}
		})
	}
}

func TestHTTPVersionComplexScenarios(t *testing.T) {
	tests := []struct {
		name            string
		curlCommand     string
		expectedVersion HTTPVersion
		expectedHTTP2   bool
	}{
		{
			name:            "HTTP/2 with POST data",
			curlCommand:     `curl --http2 -X POST -d '{"test":true}' -H "Content-Type: application/json" https://httpbin.org/post`,
			expectedVersion: HTTPVersion2,
			expectedHTTP2:   true,
		},
		{
			name:            "HTTP/1.1 with authentication",
			curlCommand:     `curl --http1.1 --digest user:pass https://httpbin.org/digest-auth/auth/user/pass`,
			expectedVersion: HTTPVersion11,
			expectedHTTP2:   false,
		},
		{
			name:            "HTTP/1.0 with headers",
			curlCommand:     `curl --http1.0 -H "Accept: application/json" -H "User-Agent: TestAgent" https://httpbin.org/get`,
			expectedVersion: HTTPVersion10,
			expectedHTTP2:   false,
		},
		{
			name:            "Protocol with timeout",
			curlCommand:     `curl --http2 --max-time 30s --connect-timeout 10s https://httpbin.org/delay/1`,
			expectedVersion: HTTPVersion2,
			expectedHTTP2:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := Parse(tt.curlCommand)
			if err != nil {
				t.Errorf("Parse() error: %v", err)
				return
			}

			if c.HTTPVersion != tt.expectedVersion {
				t.Errorf("HTTPVersion = %v, want %v", c.HTTPVersion, tt.expectedVersion)
			}

			if c.HTTP2 != tt.expectedHTTP2 {
				t.Errorf("HTTP2 = %v, want %v", c.HTTP2, tt.expectedHTTP2)
			}

			// 验证其他功能没有受到影响
			if c.ParsedURL == nil {
				t.Errorf("ParsedURL is nil")
			}

			// 验证可以创建会话
			session := c.CreateSession()
			if session == nil {
				t.Errorf("CreateSession() returned nil")
			}
		})
	}
}

// containsString 检查字符串是否包含子字符串
func containsString(str, substr string) bool {
	return len(str) >= len(substr) &&
		(len(substr) == 0 || findString(str, substr) >= 0)
}

// findString 查找子字符串位置
func findString(str, substr string) int {
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
