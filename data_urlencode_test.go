package gcurl

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDataUrlencodeParser(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		desc     string
	}{
		{
			name:     "simple_content",
			input:    "Hello World!",
			expected: "Hello+World%21",
			desc:     "简单内容URL编码",
		},
		{
			name:     "equals_prefix",
			input:    "=Hello World!",
			expected: "Hello+World%21",
			desc:     "=content 格式",
		},
		{
			name:     "name_value",
			input:    "message=Hello World! @#$%^&*()",
			expected: "message=Hello+World%21+%40%23%24%25%5E%26%2A%28%29",
			desc:     "name=content 格式",
		},
		{
			name:     "special_chars",
			input:    "data=Special: chars & symbols!",
			expected: "data=Special%3A+chars+%26+symbols%21",
			desc:     "特殊字符编码",
		},
		{
			name:     "chinese_chars",
			input:    "text=你好世界",
			expected: "text=%E4%BD%A0%E5%A5%BD%E4%B8%96%E7%95%8C",
			desc:     "中文字符编码",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			curl := New()

			err := handleDataUrlencode(curl, tt.input)
			if err != nil {
				t.Errorf("handleDataUrlencode() error = %v", err)
				return
			}

			if curl.Body == nil {
				t.Error("Body is nil")
				return
			}

			bodyContent := curl.Body.String()
			if bodyContent != tt.expected {
				t.Errorf("handleDataUrlencode() got = %q, want = %q", bodyContent, tt.expected)
			}

			// 验证Content-Type设置
			if curl.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
				t.Errorf("Content-Type = %q, want %q", curl.Header.Get("Content-Type"), "application/x-www-form-urlencoded")
			}

			// 验证Method设置
			if curl.Method != "POST" {
				t.Errorf("Method = %q, want %q", curl.Method, "POST")
			}
		})
	}
}

func TestDataUrlencodeFile(t *testing.T) {
	// 创建临时文件
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	testContent := "Hello World!\nLine 2 with special chars: @#$%"

	err := os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name     string
		input    string
		expected string
		desc     string
	}{
		{
			name:     "file_only",
			input:    "@" + testFile,
			expected: url.QueryEscape(testContent),
			desc:     "@filename 格式",
		},
		{
			name:     "name_file",
			input:    "data@" + testFile,
			expected: "data=" + url.QueryEscape(testContent),
			desc:     "name@filename 格式",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			curl := New()

			err := handleDataUrlencode(curl, tt.input)
			if err != nil {
				t.Errorf("handleDataUrlencode() error = %v", err)
				return
			}

			if curl.Body == nil {
				t.Error("Body is nil")
				return
			}

			bodyContent := curl.Body.String()
			if bodyContent != tt.expected {
				t.Errorf("handleDataUrlencode() got = %q, want = %q", bodyContent, tt.expected)
			}
		})
	}
}

func TestDataUrlencodeNameAtFileFormat(t *testing.T) {
	// 创建临时文件用于测试name@filename格式
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "data.txt")
	testContent := "file content with special chars: !@#$%"

	err := os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	curl := New()
	err = handleDataUrlencode(curl, "field@"+testFile)
	if err != nil {
		t.Fatalf("handleDataUrlencode() error = %v", err)
	}

	expected := "field=" + url.QueryEscape(testContent)
	bodyContent := curl.Body.String()
	if bodyContent != expected {
		t.Errorf("name@filename format: got = %q, want = %q", bodyContent, expected)
	}
}

func TestDataUrlencodeFileError(t *testing.T) {
	curl := New()

	// 测试不存在的文件
	err := handleDataUrlencode(curl, "@nonexistent.txt")
	if err == nil {
		t.Error("Expected error for nonexistent file, got nil")
	}
	if !strings.Contains(err.Error(), "failed to read file") {
		t.Errorf("Expected 'failed to read file' error, got: %v", err)
	}
}

func TestDataUrlencodeMultiple(t *testing.T) {
	curl := New()

	// 第一次调用
	err := handleDataUrlencode(curl, "name1=value1")
	if err != nil {
		t.Fatalf("First handleDataUrlencode() error = %v", err)
	}

	// 第二次调用应该追加
	err = handleDataUrlencode(curl, "name2=value2")
	if err != nil {
		t.Fatalf("Second handleDataUrlencode() error = %v", err)
	}

	expected := "name1=value1&name2=value2"
	bodyContent := curl.Body.String()
	if bodyContent != expected {
		t.Errorf("Multiple handleDataUrlencode() got = %q, want = %q", bodyContent, expected)
	}
}

func TestDataUrlencodeInvalidFormat(t *testing.T) {
	tests := []struct {
		name  string
		input string
		desc  string
	}{
		{
			name:  "invalid_at_format",
			input: "name@@file",
			desc:  "无效的@格式",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			curl := New()

			err := handleDataUrlencode(curl, tt.input)
			if err == nil {
				t.Errorf("Expected error for invalid format %q, got nil", tt.input)
			}
		})
	}
}

func TestDataUrlencodeIntegration(t *testing.T) {
	// 测试完整的curl命令解析
	tests := []struct {
		name    string
		curlCmd string
		check   func(*CURL) error
		desc    string
	}{
		{
			name:    "basic_urlencode",
			curlCmd: `curl --data-urlencode "message=Hello World!" https://httpbin.org/post`,
			check: func(c *CURL) error {
				expected := "message=Hello+World%21"
				if c.Body.String() != expected {
					return fmt.Errorf("body = %q, want %q", c.Body.String(), expected)
				}
				return nil
			},
			desc: "基本URL编码测试",
		},
		{
			name:    "equals_format",
			curlCmd: `curl --data-urlencode "=Hello World!" https://httpbin.org/post`,
			check: func(c *CURL) error {
				expected := "Hello+World%21"
				if c.Body.String() != expected {
					return fmt.Errorf("body = %q, want %q", c.Body.String(), expected)
				}
				return nil
			},
			desc: "=content格式测试",
		},
		{
			name:    "multiple_data_urlencode",
			curlCmd: `curl --data-urlencode "name=John" --data-urlencode "city=New York" https://httpbin.org/post`,
			check: func(c *CURL) error {
				expected := "name=John&city=New+York"
				if c.Body.String() != expected {
					return fmt.Errorf("body = %q, want %q", c.Body.String(), expected)
				}
				return nil
			},
			desc: "多个--data-urlencode选项测试",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			curl, err := Parse(tt.curlCmd)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			if err := tt.check(curl); err != nil {
				t.Errorf("Check failed: %v", err)
			}
		})
	}
}
