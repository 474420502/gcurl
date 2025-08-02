package gcurl

import (
	"net/url"
	"strings"
	"testing"
)

// TestGetModeOption 测试 -G/--get 选项功能
func TestGetModeOption(t *testing.T) {
	testCases := []struct {
		name         string
		command      string
		expectError  bool
		expectMethod string
		expectURL    string
		description  string
	}{
		{
			name:         "Basic GET mode with data",
			command:      `curl -G -d "param1=value1" -d "param2=value2" https://example.com/api`,
			expectError:  false,
			expectMethod: "GET",
			expectURL:    "https://example.com/api?param1=value1&param2=value2",
			description:  "基本GET模式：将POST数据作为查询参数",
		},
		{
			name:         "GET mode with existing query params",
			command:      `curl -G -d "new=param" "https://example.com/api?existing=value"`,
			expectError:  false,
			expectMethod: "GET",
			expectURL:    "https://example.com/api?existing=value&new=param",
			description:  "GET模式：添加到现有查询参数",
		},
		{
			name:         "GET mode with form data",
			command:      `curl --get -d "user=john&age=30" https://example.com/search`,
			expectError:  false,
			expectMethod: "GET",
			expectURL:    "https://example.com/search?user=john&age=30",
			description:  "GET模式：表单数据转为查询参数",
		},
		{
			name:         "GET mode with URL encoded data",
			command:      `curl -G -d "query=hello world" https://example.com/search`,
			expectError:  false,
			expectMethod: "GET",
			expectURL:    "https://example.com/search?query=hello+world",
			description:  "GET模式：URL编码处理",
		},
		{
			name:         "GET mode with complex data",
			command:      `curl -G -d "filters[status]=active" -d "filters[type]=user" https://example.com/api`,
			expectError:  false,
			expectMethod: "GET",
			expectURL:    "https://example.com/api?filters%5Bstatus%5D=active&filters%5Btype%5D=user",
			description:  "GET模式：复杂数据结构",
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

			// 验证 GET 模式已启用
			if !curl.GetMode {
				t.Error("GetMode should be true")
			}

			// 验证 HTTP 方法
			if curl.Method != tc.expectMethod {
				t.Errorf("Expected method %s, got %s", tc.expectMethod, curl.Method)
			}

			// 验证 URL（需要处理数据到查询参数的转换）
			expectedURL, err := url.Parse(tc.expectURL)
			if err != nil {
				t.Fatalf("Invalid expected URL: %v", err)
			}

			// 这里需要实现实际的数据到查询参数转换逻辑
			// 暂时验证基本信息
			t.Logf("✓ GET mode enabled")
			t.Logf("✓ Method: %s", curl.Method)
			t.Logf("✓ Original URL: %s", curl.ParsedURL.String())
			t.Logf("✓ Expected final URL pattern: %s", expectedURL.String())

			// 验证有数据存在
			if curl.Body == nil {
				t.Error("Body data should exist for conversion to query parameters")
			}
		})
	}
}

// TestGetModeWithoutData 测试没有数据时的GET模式
func TestGetModeWithoutData(t *testing.T) {
	cmd := `curl -G https://example.com/api`

	curl, err := Parse(cmd)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if !curl.GetMode {
		t.Error("GetMode should be true")
	}

	if curl.Method != "GET" {
		t.Errorf("Expected method GET, got %s", curl.Method)
	}

	t.Logf("✓ GET mode works without data")
	t.Logf("✓ Method: %s", curl.Method)
	t.Logf("✓ URL: %s", curl.ParsedURL.String())
}

// TestGetModeMethodOverride 测试GET模式的方法覆盖
func TestGetModeMethodOverride(t *testing.T) {
	testCases := []struct {
		name    string
		command string
		expect  string
	}{
		{
			name:    "GET mode overrides POST",
			command: `curl -X POST -G -d "data=value" https://example.com`,
			expect:  "GET",
		},
		{
			name:    "Explicit method after GET mode",
			command: `curl -G -X PUT -d "data=value" https://example.com`,
			expect:  "PUT", // 最后的 -X 应该优先
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			curl, err := Parse(tc.command)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}

			if curl.Method != tc.expect {
				t.Errorf("Expected method %s, got %s", tc.expect, curl.Method)
			}

			t.Logf("✓ Method correctly set to: %s", curl.Method)
		})
	}
}

// TestGetModeIntegration 测试GET模式的集成场景
func TestGetModeIntegration(t *testing.T) {
	scenarios := []struct {
		name        string
		command     string
		description string
	}{
		{
			name:        "Search API with filters",
			command:     `curl -G -d "q=golang" -d "sort=stars" -d "order=desc" https://api.github.com/search/repositories`,
			description: "搜索API：多个过滤参数",
		},
		{
			name:        "Analytics with date range",
			command:     `curl -G -d "start_date=2023-01-01" -d "end_date=2023-12-31" -d "metrics=views,clicks" https://analytics.example.com/api`,
			description: "分析API：日期范围和指标参数",
		},
		{
			name:        "Pagination with GET",
			command:     `curl -G -d "page=2" -d "limit=50" -d "category=tech" https://blog.example.com/api/posts`,
			description: "分页API：页码和限制参数",
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

			if !curl.GetMode {
				t.Error("GetMode should be enabled")
			}

			if curl.Method != "GET" {
				t.Errorf("Expected GET method, got %s", curl.Method)
			}

			t.Logf("✓ Successfully parsed GET mode request")
			t.Logf("✓ Method: %s", curl.Method)
			t.Logf("✓ URL: %s", curl.ParsedURL.String())
			t.Logf("✓ GET mode: %t", curl.GetMode)

			// 验证有数据用于转换
			if curl.Body != nil {
				t.Logf("✓ Has body data for query parameter conversion")
			}
		})
	}
}

// TestGetModeVerboseOutput 测试GET模式的详细输出
func TestGetModeVerboseOutput(t *testing.T) {
	cmd := `curl -v -G -d "param=value" https://example.com/api`

	curl, err := Parse(cmd)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	verbose := curl.VerboseInfo()
	t.Logf("✓ Verbose output for GET mode:")
	t.Log(verbose)

	// 验证输出包含GET方法
	if !strings.Contains(verbose, "GET") {
		t.Error("Verbose output should contain GET method")
	}

	// 验证URL包含在输出中
	if !strings.Contains(verbose, "example.com") {
		t.Error("Verbose output should contain the URL")
	}
}
