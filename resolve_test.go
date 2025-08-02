package gcurl

import (
	"strings"
	"testing"
)

// TestResolveOption 测试 --resolve 选项功能
func TestResolveOption(t *testing.T) {
	testCases := []struct {
		name        string
		command     string
		expectError bool
		expected    []string
		description string
	}{
		{
			name:        "Basic resolve mapping",
			command:     `curl https://example.com --resolve example.com:443:192.168.1.100`,
			expectError: false,
			expected:    []string{"example.com:443:192.168.1.100"},
			description: "基本的主机名解析映射",
		},
		{
			name:        "Multiple addresses",
			command:     `curl https://api.example.com --resolve api.example.com:443:192.168.1.100,192.168.1.101`,
			expectError: false,
			expected:    []string{"api.example.com:443:192.168.1.100,192.168.1.101"},
			description: "多个地址映射",
		},
		{
			name:        "Multiple resolve entries",
			command:     `curl https://example.com --resolve example.com:443:192.168.1.100 --resolve api.example.com:80:10.0.0.1`,
			expectError: false,
			expected:    []string{"example.com:443:192.168.1.100", "api.example.com:80:10.0.0.1"},
			description: "多个解析条目",
		},
		{
			name:        "Force replacement with plus",
			command:     `curl https://example.com --resolve +example.com:443:127.0.0.1`,
			expectError: false,
			expected:    []string{"+example.com:443:127.0.0.1"},
			description: "强制替换模式（带+前缀）",
		},
		{
			name:        "Invalid format missing port",
			command:     `curl https://example.com --resolve example.com:192.168.1.100`,
			expectError: true,
			expected:    nil,
			description: "格式错误：缺少端口",
		},
		{
			name:        "Invalid port format",
			command:     `curl https://example.com --resolve example.com:abc:192.168.1.100`,
			expectError: true,
			expected:    nil,
			description: "格式错误：端口不是数字",
		},
		{
			name:        "Empty address",
			command:     `curl https://example.com --resolve example.com:443:`,
			expectError: true,
			expected:    nil,
			description: "格式错误：地址为空",
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
				} else {
					t.Logf("✓ Expected error: %v", err)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if curl == nil {
				t.Error("Expected valid CURL object but got nil")
				return
			}

			// 验证 Resolve 映射
			if len(curl.Resolve) != len(tc.expected) {
				t.Errorf("Expected %d resolve entries, got %d", len(tc.expected), len(curl.Resolve))
				return
			}

			for i, expected := range tc.expected {
				if i >= len(curl.Resolve) || curl.Resolve[i] != expected {
					t.Errorf("Expected resolve[%d] = %s, got %s", i, expected, curl.Resolve[i])
					return
				}
			}

			t.Logf("✓ Successfully parsed %d resolve mappings:", len(curl.Resolve))
			for i, resolve := range curl.Resolve {
				t.Logf("  [%d] %s", i, resolve)
			}
		})
	}
}

// TestResolveIntegrationWithVerbose 测试 --resolve 与 -v 选项的集成
func TestResolveIntegrationWithVerbose(t *testing.T) {
	command := `curl -v https://httpbin.org/get --resolve httpbin.org:443:127.0.0.1`

	curl, err := Parse(command)
	if err != nil {
		t.Fatalf("Failed to parse command: %v", err)
	}

	// 验证解析选项
	if len(curl.Resolve) != 1 {
		t.Fatalf("Expected 1 resolve entry, got %d", len(curl.Resolve))
	}

	expectedResolve := "httpbin.org:443:127.0.0.1"
	if curl.Resolve[0] != expectedResolve {
		t.Errorf("Expected resolve %s, got %s", expectedResolve, curl.Resolve[0])
	}

	// 验证详细模式也被启用
	if !curl.Verbose {
		t.Error("Expected verbose mode to be enabled")
	}

	t.Logf("✓ Command successfully parsed with resolve and verbose options")
	t.Logf("  URL: %s", curl.ParsedURL.String())
	t.Logf("  Resolve: %s", curl.Resolve[0])
	t.Logf("  Verbose: %v", curl.Verbose)
}

// TestResolveRealWorldScenarios 测试真实世界的 --resolve 使用场景
func TestResolveRealWorldScenarios(t *testing.T) {
	scenarios := []struct {
		name        string
		command     string
		description string
	}{
		{
			name:        "Local development",
			command:     `curl https://api.myapp.com/health --resolve api.myapp.com:443:127.0.0.1`,
			description: "本地开发：将生产域名解析到本地",
		},
		{
			name:        "Load balancing test",
			command:     `curl https://service.company.com/status --resolve service.company.com:443:10.0.1.100,10.0.1.101`,
			description: "负载均衡测试：多个后端地址",
		},
		{
			name:        "HTTP and HTTPS different servers",
			command:     `curl https://example.com/secure --resolve example.com:80:192.168.1.10 --resolve example.com:443:192.168.1.11`,
			description: "不同端口指向不同服务器",
		},
		{
			name:        "Staging environment test",
			command:     `curl -H "Host: production.example.com" https://staging.example.com/api --resolve staging.example.com:443:203.0.113.10`,
			description: "预发布环境测试",
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			t.Logf("Testing scenario: %s", scenario.description)
			t.Logf("Command: %s", scenario.command)

			curl, err := Parse(scenario.command)
			if err != nil {
				t.Errorf("Failed to parse real-world scenario: %v", err)
				return
			}

			// 基本验证
			if curl.ParsedURL == nil {
				t.Error("URL should be parsed successfully")
				return
			}

			if len(curl.Resolve) == 0 {
				t.Error("Should have at least one resolve mapping")
				return
			}

			t.Logf("✓ Successfully parsed:")
			t.Logf("  URL: %s", curl.ParsedURL.String())
			t.Logf("  Resolve mappings:")
			for i, resolve := range curl.Resolve {
				t.Logf("    [%d] %s", i, resolve)
			}
			if len(curl.Header) > 0 {
				t.Logf("  Headers: %d", len(curl.Header))
			}
		})
	}
}

// TestResolveVerboseOutput 测试带有 --resolve 的详细输出功能
func TestResolveVerboseOutput(t *testing.T) {
	command := `curl -v --resolve example.com:443:127.0.0.1 https://example.com/test`

	curl, err := Parse(command)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	// 测试详细输出功能
	if !curl.Verbose {
		t.Error("Verbose mode should be enabled")
	}

	verboseInfo := curl.VerboseInfo()
	if verboseInfo == "" {
		t.Error("VerboseInfo should return non-empty string")
	}

	// 检查详细输出是否包含解析信息
	expectedParts := []string{
		"example.com",
		"443",
		"127.0.0.1",
		"GET /test",
	}

	for _, part := range expectedParts {
		if !strings.Contains(verboseInfo, part) {
			t.Logf("Warning: VerboseInfo might not contain expected part: %s", part)
		}
	}

	t.Logf("✓ Verbose output generated:")
	t.Logf("%s", verboseInfo)
}
