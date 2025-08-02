package gcurl

import (
	"strings"
	"testing"
)

// TestConnectToOption 测试 --connect-to 选项功能
func TestConnectToOption(t *testing.T) {
	testCases := []struct {
		name        string
		command     string
		expectError bool
		expected    []string
		description string
	}{
		{
			name:        "Basic connect-to mapping",
			command:     `curl https://example.com --connect-to example.com:443:127.0.0.1:8443`,
			expectError: false,
			expected:    []string{"example.com:443:127.0.0.1:8443"},
			description: "基本的连接重定向映射",
		},
		{
			name:        "Multiple connect-to entries",
			command:     `curl https://example.com --connect-to example.com:443:127.0.0.1:8443 --connect-to api.example.com:80:localhost:8080`,
			expectError: false,
			expected:    []string{"example.com:443:127.0.0.1:8443", "api.example.com:80:localhost:8080"},
			description: "多个连接重定向条目",
		},
		{
			name:        "Wildcard source host",
			command:     `curl https://example.com --connect-to ::proxy.example.com:8080`,
			expectError: false,
			expected:    []string{"::proxy.example.com:8080"},
			description: "通配符源主机（任意主机端口到代理）",
		},
		{
			name:        "Invalid format - missing parts",
			command:     `curl https://example.com --connect-to example.com:443:127.0.0.1`,
			expectError: true,
			expected:    nil,
			description: "格式错误：缺少端口部分",
		},
		{
			name:        "Invalid source port",
			command:     `curl https://example.com --connect-to example.com:abc:127.0.0.1:8443`,
			expectError: true,
			expected:    nil,
			description: "格式错误：源端口不是数字",
		},
		{
			name:        "Invalid target port",
			command:     `curl https://example.com --connect-to example.com:443:127.0.0.1:xyz`,
			expectError: true,
			expected:    nil,
			description: "格式错误：目标端口不是数字",
		},
		{
			name:        "Empty target host",
			command:     `curl https://example.com --connect-to example.com:443::8443`,
			expectError: true,
			expected:    nil,
			description: "格式错误：目标主机为空",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("Testing: %s", tc.description)
			t.Logf("Command: %s", tc.command)

			curl, err := Parse(tc.command)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error but parsing succeeded")
					return
				}
				t.Logf("✓ Expected error: %v", err)
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if len(curl.ConnectTo) != len(tc.expected) {
				t.Errorf("Expected %d connect-to mappings, got %d", len(tc.expected), len(curl.ConnectTo))
				return
			}

			t.Logf("✓ Successfully parsed %d connect-to mappings:", len(curl.ConnectTo))
			for i, connectTo := range curl.ConnectTo {
				if i < len(tc.expected) && connectTo != tc.expected[i] {
					t.Errorf("Expected [%d] %s, got %s", i, tc.expected[i], connectTo)
					return
				}
				t.Logf("  [%d] %s", i, connectTo)
			}
		})
	}
}

// TestConnectToIntegrationWithVerbose 测试 --connect-to 与 -v 选项的集成
func TestConnectToIntegrationWithVerbose(t *testing.T) {
	cmd := `curl -v https://httpbin.org/get --connect-to httpbin.org:443:127.0.0.1:8443`

	curl, err := Parse(cmd)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	t.Logf("✓ Command successfully parsed with connect-to and verbose options")
	t.Logf("  URL: %s", curl.ParsedURL.String())
	t.Logf("  ConnectTo: %s", strings.Join(curl.ConnectTo, ", "))
	t.Logf("  Verbose: %t", curl.Verbose)

	// 测试详细输出包含连接重定向信息
	verbose := curl.VerboseInfo()
	if !strings.Contains(verbose, "Connection redirects") {
		t.Error("Verbose output should contain connection redirect information")
	}
	if !strings.Contains(verbose, "httpbin.org:443 -> 127.0.0.1:8443") {
		t.Error("Verbose output should show specific connection redirect mapping")
	}
}

// TestConnectToRealWorldScenarios 测试实际使用场景
func TestConnectToRealWorldScenarios(t *testing.T) {
	scenarios := []struct {
		name        string
		command     string
		description string
	}{
		{
			name:        "Local development",
			command:     `curl https://api.myapp.com/health --connect-to api.myapp.com:443:localhost:3000`,
			description: "本地开发：将生产API重定向到本地开发服务器",
		},
		{
			name:        "Testing with staging",
			command:     `curl https://production.example.com/api --connect-to production.example.com:443:staging.internal.com:443`,
			description: "测试环境：将生产域名重定向到内部测试服务器",
		},
		{
			name:        "Proxy all connections",
			command:     `curl https://example.com/test --connect-to ::proxy.company.com:8080`,
			description: "代理设置：所有连接都通过代理服务器",
		},
		{
			name:        "Load balancer testing",
			command:     `curl https://service.com/status --connect-to service.com:443:backend1.internal:443`,
			description: "负载均衡测试：直接连接到特定后端服务器",
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			t.Logf("Testing scenario: %s", scenario.description)
			t.Logf("Command: %s", scenario.command)

			curl, err := Parse(scenario.command)
			if err != nil {
				t.Errorf("Parse failed: %v", err)
				return
			}

			t.Logf("✓ Successfully parsed:")
			t.Logf("  URL: %s", curl.ParsedURL.String())
			t.Logf("  Connect-to mappings:")
			for i, connectTo := range curl.ConnectTo {
				t.Logf("    [%d] %s", i, connectTo)
			}
		})
	}
}

// TestConnectToVerboseOutput 测试连接重定向的详细输出
func TestConnectToVerboseOutput(t *testing.T) {
	cmd := `curl -v https://example.com/test --connect-to example.com:443:127.0.0.1:8443 --connect-to ::proxy.local:8080`

	curl, err := Parse(cmd)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	verbose := curl.VerboseInfo()
	t.Logf("✓ Verbose output generated:")
	t.Log(verbose)

	// 验证输出包含连接重定向信息
	if !strings.Contains(verbose, "Connection redirects") {
		t.Error("Should contain connection redirects section")
	}
	if !strings.Contains(verbose, "example.com:443 -> 127.0.0.1:8443") {
		t.Error("Should show first connection redirect")
	}
	if !strings.Contains(verbose, "*:* -> proxy.local:8080") {
		t.Error("Should show wildcard connection redirect")
	}
}
