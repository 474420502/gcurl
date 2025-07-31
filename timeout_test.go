package gcurl

import (
	"strings"
	"testing"
	"time"
)

func TestEnhancedTimeoutSystem(t *testing.T) {
	t.Run("基本超时设置", func(t *testing.T) {
		curl := New()

		// 检查默认超时
		if curl.Timeout != 30*time.Second {
			t.Errorf("期望默认超时为30秒，得到 %v", curl.Timeout)
		}

		// 检查默认连接超时
		if curl.ConnectTimeout != 0 {
			t.Errorf("期望默认连接超时为0，得到 %v", curl.ConnectTimeout)
		}

		// 检查默认DNS超时
		if curl.DNSTimeout != 0 {
			t.Errorf("期望默认DNS超时为0，得到 %v", curl.DNSTimeout)
		}

		// 检查默认TLS握手超时
		if curl.TLSHandshakeTimeout != 0 {
			t.Errorf("期望默认TLS握手超时为0，得到 %v", curl.TLSHandshakeTimeout)
		}
	})

	t.Run("--max-time选项解析", func(t *testing.T) {
		tests := []struct {
			name     string
			cmd      string
			expected time.Duration
		}{
			{
				name:     "纯数字秒",
				cmd:      `curl --max-time 60 "https://example.com"`,
				expected: 60 * time.Second,
			},
			{
				name:     "带秒单位",
				cmd:      `curl --max-time 45s "https://example.com"`,
				expected: 45 * time.Second,
			},
			{
				name:     "分钟单位",
				cmd:      `curl --max-time 2m "https://example.com"`,
				expected: 2 * time.Minute,
			},
			{
				name:     "小时单位",
				cmd:      `curl --max-time 1h "https://example.com"`,
				expected: 1 * time.Hour,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				curl, err := Parse(tt.cmd)
				if err != nil {
					t.Fatalf("解析失败: %v", err)
				}

				if curl.Timeout != tt.expected {
					t.Errorf("期望超时为 %v，得到 %v", tt.expected, curl.Timeout)
				}
			})
		}
	})

	t.Run("--connect-timeout选项解析", func(t *testing.T) {
		tests := []struct {
			name     string
			cmd      string
			expected time.Duration
		}{
			{
				name:     "纯数字秒",
				cmd:      `curl --connect-timeout 10 "https://example.com"`,
				expected: 10 * time.Second,
			},
			{
				name:     "带秒单位",
				cmd:      `curl --connect-timeout 15s "https://example.com"`,
				expected: 15 * time.Second,
			},
			{
				name:     "分钟单位",
				cmd:      `curl --connect-timeout 1m "https://example.com"`,
				expected: 1 * time.Minute,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				curl, err := Parse(tt.cmd)
				if err != nil {
					t.Fatalf("解析失败: %v", err)
				}

				if curl.ConnectTimeout != tt.expected {
					t.Errorf("期望连接超时为 %v，得到 %v", tt.expected, curl.ConnectTimeout)
				}
			})
		}
	})

	t.Run("组合超时选项", func(t *testing.T) {
		curl, err := Parse(`curl --max-time 2m --connect-timeout 30s "https://example.com"`)
		if err != nil {
			t.Fatalf("解析失败: %v", err)
		}

		if curl.Timeout != 2*time.Minute {
			t.Errorf("期望总超时为2分钟，得到 %v", curl.Timeout)
		}

		if curl.ConnectTimeout != 30*time.Second {
			t.Errorf("期望连接超时为30秒，得到 %v", curl.ConnectTimeout)
		}
	})

	t.Run("无效超时值", func(t *testing.T) {
		invalidCmds := []string{
			`curl --max-time invalid "https://example.com"`,
			`curl --connect-timeout abc "https://example.com"`,
			`curl --max-time -1 "https://example.com"`,
		}

		for _, cmd := range invalidCmds {
			_, err := Parse(cmd)
			if err == nil {
				t.Errorf("期望命令 '%s' 解析失败，但成功了", cmd)
			}
		}
	})
}

func TestTimeoutDebugOutput(t *testing.T) {
	t.Run("Debug方法显示超时信息", func(t *testing.T) {
		curl, err := Parse(`curl --max-time 2m --connect-timeout 30s "https://example.com"`)
		if err != nil {
			t.Fatalf("解析失败: %v", err)
		}

		debug := curl.Debug()

		// 检查是否包含超时信息
		if !strings.Contains(debug, "Timeout: 2m0s") {
			t.Error("调试输出应该包含总超时信息")
		}

		if !strings.Contains(debug, "Connect Timeout: 30s") {
			t.Error("调试输出应该包含连接超时信息")
		}
	})
}

func TestTimeoutSessionConfiguration(t *testing.T) {
	t.Run("Session配置超时", func(t *testing.T) {
		curl, err := Parse(`curl --max-time 1m "https://httpbin.org/delay/1"`)
		if err != nil {
			t.Fatalf("解析失败: %v", err)
		}

		// 创建session并检查配置
		session := curl.CreateSession()
		if session == nil {
			t.Error("CreateSession返回nil")
		}

		// 注意：这里我们无法直接检查session的超时配置
		// 因为requests库没有暴露这些信息
		// 但我们可以确保代码正确运行而不出错
	})
}

func TestTimeoutTypeSafety(t *testing.T) {
	t.Run("类型安全检查", func(t *testing.T) {
		curl := New()

		// 设置各种超时
		curl.Timeout = 30 * time.Second
		curl.ConnectTimeout = 10 * time.Second
		curl.DNSTimeout = 5 * time.Second
		curl.TLSHandshakeTimeout = 15 * time.Second

		// 验证类型安全 - 编译时检查
		var _ time.Duration = curl.Timeout
		var _ time.Duration = curl.ConnectTimeout
		var _ time.Duration = curl.DNSTimeout
		var _ time.Duration = curl.TLSHandshakeTimeout

		// 验证值
		if curl.Timeout != 30*time.Second {
			t.Error("Timeout类型安全检查失败")
		}
		if curl.ConnectTimeout != 10*time.Second {
			t.Error("ConnectTimeout类型安全检查失败")
		}
		if curl.DNSTimeout != 5*time.Second {
			t.Error("DNSTimeout类型安全检查失败")
		}
		if curl.TLSHandshakeTimeout != 15*time.Second {
			t.Error("TLSHandshakeTimeout类型安全检查失败")
		}
	})
}
