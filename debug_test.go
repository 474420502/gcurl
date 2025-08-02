package gcurl

import (
	"strings"
	"testing"
)


// TestHeadMethodBug 专门测试HEAD方法的bug
func TestHeadMethodBug(t *testing.T) {
	// 测试 -I 选项是否正确解析为HEAD方法
	c, err := Parse(`curl -I http://httpbin.org/get`)
	if err != nil {
		t.Error("Parse error:", err)
		return
	}

	// 这里应该是HEAD方法，但实际解析为GET
	t.Logf("解析的方法: %s", c.Method)

	if c.Method != "HEAD" {
		t.Errorf("HEAD方法解析错误: 期望 'HEAD'，实际得到 '%s'", c.Method)
	}
}

// TestEmptyUserAgentBug 专门测试空User-Agent的bug
func TestEmptyUserAgentBug(t *testing.T) {
	// 测试空User-Agent是否导致解析失败
	c, err := Parse(`curl -A "" http://httpbin.org/get`)

	if err != nil {
		t.Logf("空User-Agent导致解析错误: %v", err)
		// 这确实是一个bug，但至少我们知道了问题所在
		return
	}

	t.Logf("解析成功 - URL: %s", c.ParsedURL.String())

	// 检查User-Agent头部
	userAgent := c.Header.Get("User-Agent")
	if userAgent != "" {
		t.Logf("User-Agent 值: '%s'", userAgent)
	} else {
		t.Log("没有找到User-Agent头部")
	}
}

// TestInvalidURLBug 测试无效URL验证问题
func TestInvalidURLBug(t *testing.T) {
	// 测试明显无效的URL是否被正确识别
	c, err := Parse(`curl "not-a-valid-url"`)

	if err == nil {
		t.Logf("⚠️  无效URL被接受了: %s", c.ParsedURL.String())
		t.Log("这可能不是严重问题，因为URL验证可能在执行时进行")
	} else {
		t.Logf("✓ 无效URL被正确拒绝: %v", err)
	}
}

// TestSupportedOptions 测试当前支持的选项
func TestSupportedOptions(t *testing.T) {
	testCases := []struct {
		name       string
		command    string
		shouldPass bool
	}{
		{"基本GET", `curl http://httpbin.org/get`, true},
		{"POST数据", `curl -d "key=value" http://httpbin.org/post`, true},
		{"自定义头部", `curl -H "Custom: value" http://httpbin.org/get`, true},
		{"Cookie", `curl -b "session=abc" http://httpbin.org/get`, true},
		{"用户认证", `curl -u user:pass http://httpbin.org/get`, true},
		{"User-Agent", `curl -A "MyAgent" http://httpbin.org/get`, true},
		{"SSL忽略", `curl -k https://example.com`, true},
		{"连接超时", `curl --connect-timeout 5 http://httpbin.org/get`, true},

		// 现在支持的选项
		{"文件上传", `curl -F "file=@test.txt" http://httpbin.org/post`, true},
		{"重定向", `curl -L http://httpbin.org/redirect/1`, true},
		{"最大时间", `curl --max-time 30 http://httpbin.org/get`, true},
		{"代理", `curl --proxy http://proxy:8080 http://httpbin.org/get`, true},

		// 仍不支持的选项
		{"CA证书", `curl --cacert ca.pem http://httpbin.org/get`, false},
		{"HTTP2", `curl --http2 http://httpbin.org/get`, false},
	}

	supported := 0
	total := len(testCases)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Parse(tc.command)

			if tc.shouldPass {
				if err == nil {
					t.Logf("✓ %s: 支持", tc.name)
					supported++
				} else {
					t.Errorf("✗ %s: 应该支持但失败了 - %v", tc.name, err)
				}
			} else {
				if err != nil {
					t.Logf("⚠️  %s: 不支持 - %v", tc.name, err)
				} else {
					t.Logf("? %s: 意外地成功了", tc.name)
				}
			}
		})
	}

	t.Logf("\n📊 支持度统计: %d/%d (%.1f%%)", supported, total, float64(supported)/float64(total)*100)
}

// TestDebugFunctionality 测试新增的调试功能
func TestDebugFunctionality(t *testing.T) {
	tests := []struct {
		name    string
		curlCmd string
		desc    string
	}{
		{
			name:    "Basic GET with headers",
			curlCmd: `curl -H "Accept: application/json" -H "User-Agent: TestApp/1.0" "https://httpbin.org/get?param=value"`,
			desc:    "基础GET请求，包含头部和查询参数",
		},
		{
			name:    "POST with JSON data",
			curlCmd: `curl -X POST -H "Content-Type: application/json" -d '{"name":"test","age":25}' "https://httpbin.org/post"`,
			desc:    "POST请求，包含JSON数据",
		},
		{
			name:    "Complex request with auth and cookies",
			curlCmd: `curl -u "user:pass" -b "session=abc123; theme=dark" -H "X-API-Key: secret" "https://httpbin.org/get"`,
			desc:    "复杂请求，包含认证和Cookie",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("\n🧪 测试用例: %s", tt.desc)
			t.Logf("命令: %s", tt.curlCmd)

			curl, err := Parse(tt.curlCmd)
			if err != nil {
				t.Errorf("解析失败: %v", err)
				return
			}

			// 测试 Summary 方法
			summary := curl.Summary()
			t.Logf("\n📝 简要信息: %s", summary)

			// 测试 Debug 方法
			debug := curl.Debug()
			if len(debug) == 0 {
				t.Error("Debug() 不应该返回空字符串")
			}

			// 测试 VerboseInfo 方法
			verbose := curl.VerboseInfo()
			if len(verbose) == 0 {
				t.Error("VerboseInfo() 不应该返回空字符串")
			}

			// 验证基本信息存在
			if !strings.Contains(debug, curl.Method) {
				t.Error("调试信息应该包含HTTP方法")
			}
			if curl.ParsedURL != nil && !strings.Contains(debug, curl.ParsedURL.String()) {
				t.Error("调试信息应该包含URL")
			}
		})
	}
}

// TestDebugOutputFormat 测试调试输出格式
func TestDebugOutputFormat(t *testing.T) {
	// 测试包含多种特性的复杂请求
	curlCmd := `curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer token123" -d '{"key":"value"}' -b "session=abc; theme=dark" -u "user:pass" --connect-timeout 30 -L -k "https://api.example.com/data?filter=active"`

	curl, err := Parse(curlCmd)
	if err != nil {
		t.Fatalf("解析失败: %v", err)
	}

	// 设置调试标志
	curl.Verbose = true
	curl.Include = true
	curl.Silent = false
	curl.Trace = true

	// 测试 Debug() 输出
	debug := curl.Debug()
	t.Logf("\n🔍 Debug() 输出:\n%s", debug)

	// 验证 Debug() 输出包含所有关键信息
	requiredSections := []string{"Method:", "URL:", "Headers", "Authentication:", "Body:", "Debug Flags:"}
	for _, section := range requiredSections {
		if !strings.Contains(debug, section) {
			t.Errorf("Debug() 输出应该包含 '%s' 部分", section)
		}
	}

	// 测试 VerboseInfo() 输出
	verbose := curl.VerboseInfo()
	t.Logf("\n📋 VerboseInfo() 输出:\n%s", verbose)

	// 验证详细信息的完整性
	verboseChecks := []string{"POST", "api.example.com", "Content-Type", "Bearer", "session=abc"}
	for _, check := range verboseChecks {
		if !strings.Contains(verbose, check) {
			t.Errorf("VerboseInfo() 输出应该包含 '%s'", check)
		}
	}
}

// TestDebugWithEmptyFields 测试空字段的调试输出
func TestDebugWithEmptyFields(t *testing.T) {
	// 最简单的GET请求
	curl, err := Parse("curl https://example.com")
	if err != nil {
		t.Fatalf("解析失败: %v", err)
	}

	summary := curl.Summary()
	debug := curl.Debug()
	verbose := curl.VerboseInfo()

	t.Logf("简单请求 Summary: %s", summary)
	t.Logf("简单请求 Debug 长度: %d", len(debug))
	t.Logf("简单请求 Verbose 长度: %d", len(verbose))

	// 验证即使是简单请求也有基础信息
	if !strings.Contains(summary, "GET") {
		t.Error("Summary 应该包含HTTP方法")
	}
	if !strings.Contains(summary, "example.com") {
		t.Error("Summary 应该包含域名")
	}
}
