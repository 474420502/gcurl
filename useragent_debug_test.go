package gcurl

import (
	"fmt"
	"testing"
)

// TestUserAgentDebug 详细调试User-Agent解析问题
func TestUserAgentDebug(t *testing.T) {
	testCases := []struct {
		name        string
		command     string
		expectError bool
	}{
		{
			name:        "正常User-Agent",
			command:     `curl -A "MyAgent" http://httpbin.org/get`,
			expectError: false,
		},
		{
			name:        "空User-Agent",
			command:     `curl -A "" http://httpbin.org/get`,
			expectError: false,
		},
		{
			name:        "空格User-Agent",
			command:     `curl -A " " http://httpbin.org/get`,
			expectError: false,
		},
		{
			name:        "不带引号的空User-Agent",
			command:     `curl -A  http://httpbin.org/get`,
			expectError: true, // 这可能会把URL当作User-Agent参数
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("测试命令: %s", tc.command)

			c, err := Parse(tc.command)

			if tc.expectError {
				if err != nil {
					t.Logf("✓ 预期错误: %v", err)
				} else {
					t.Errorf("✗ 应该出错但没有出错")
				}
			} else {
				if err != nil {
					t.Errorf("✗ 不应该出错但出错了: %v", err)
				} else {
					t.Logf("✓ 解析成功")
					t.Logf("  URL: %s", c.ParsedURL.String())
					t.Logf("  Method: %s", c.Method)

					userAgent := c.Header.Get("User-Agent")
					t.Logf("  User-Agent: '%s'", userAgent)
				}
			}
		})
	}
}

// TestCommandLineParsing 测试命令行解析的边界情况
func TestCommandLineParsing(t *testing.T) {
	// 测试几个可能导致混乱的命令
	testCases := []string{
		`curl -A "" http://httpbin.org/get`,
		`curl -A '' http://httpbin.org/get`,
		`curl -A"" http://httpbin.org/get`,
		`curl -A'' http://httpbin.org/get`,
		`curl --user-agent "" http://httpbin.org/get`,
		`curl --user-agent='' http://httpbin.org/get`,
	}

	for i, cmd := range testCases {
		t.Run(fmt.Sprintf("Case_%d", i+1), func(t *testing.T) {
			t.Logf("解析命令: %s", cmd)

			c, err := Parse(cmd)
			if err != nil {
				t.Logf("解析错误: %v", err)
			} else {
				t.Logf("解析成功 - URL: %s, User-Agent: '%s'",
					c.ParsedURL.String(), c.Header.Get("User-Agent"))
			}
		})
	}
}
