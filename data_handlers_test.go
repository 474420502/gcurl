package gcurl

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/474420502/requests"
)

func TestDataHandlers(t *testing.T) {
	tests := []struct {
		name        string
		curlCmd     string
		expectedLen int
		contains    []string
		desc        string
	}{
		{
			name:        "single_data",
			curlCmd:     `curl --data "name=John" https://httpbin.org/post`,
			expectedLen: 9,
			contains:    []string{"name=John"},
			desc:        "单个--data选项",
		},
		{
			name:        "multiple_data",
			curlCmd:     `curl --data "name=John" --data "age=30" https://httpbin.org/post`,
			expectedLen: 16,
			contains:    []string{"name=John", "age=30"},
			desc:        "多个--data选项",
		},
		{
			name:        "single_data_raw",
			curlCmd:     `curl --data-raw "name=John" https://httpbin.org/post`,
			expectedLen: 9,
			contains:    []string{"name=John"},
			desc:        "单个--data-raw选项",
		},
		{
			name:        "multiple_data_raw",
			curlCmd:     `curl --data-raw "name=John" --data-raw "age=30" https://httpbin.org/post`,
			expectedLen: 16,
			contains:    []string{"name=John", "age=30"},
			desc:        "多个--data-raw选项",
		},
		{
			name:        "single_data_binary",
			curlCmd:     `curl --data-binary "Hello World" https://httpbin.org/post`,
			expectedLen: 11,
			contains:    []string{"Hello World"},
			desc:        "单个--data-binary选项",
		},
		{
			name:        "multiple_data_binary",
			curlCmd:     `curl --data-binary "Hello" --data-binary " World" https://httpbin.org/post`,
			expectedLen: 11,
			contains:    []string{"Hello World"},
			desc:        "多个--data-binary选项",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			curl, err := Parse(tt.curlCmd)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			if curl.Body == nil {
				t.Error("Body is nil")
				return
			}

			bodyContent := curl.Body.String()
			if len(bodyContent) != tt.expectedLen {
				t.Errorf("Body length = %d, want %d. Body: %q", len(bodyContent), tt.expectedLen, bodyContent)
			}

			for _, expected := range tt.contains {
				if !strings.Contains(bodyContent, expected) {
					t.Errorf("Body does not contain %q. Body: %q", expected, bodyContent)
				}
			}

			// 验证Method设置
			if curl.Method != "POST" {
				t.Errorf("Method = %q, want %q", curl.Method, "POST")
			}
		})
	}
}

func TestDataHandlersFile(t *testing.T) {
	// 创建临时文件
	tmpDir := t.TempDir()
	testFile1 := filepath.Join(tmpDir, "data1.txt")
	testFile2 := filepath.Join(tmpDir, "data2.txt")

	content1 := "name=John\nage=30"
	content2 := "city=NYC\ncountry=USA"

	err := os.WriteFile(testFile1, []byte(content1), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	err = os.WriteFile(testFile2, []byte(content2), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name     string
		handler  func(*CURL, ...string) error
		inputs   []string
		expected string
		desc     string
	}{
		{
			name:     "data_from_file",
			handler:  handleData,
			inputs:   []string{"@" + testFile1},
			expected: "name=Johnage=30", // handleData removes newlines
			desc:     "--data从文件读取（删除换行符）",
		},
		{
			name:     "data_binary_from_file",
			handler:  handleDataBinary,
			inputs:   []string{"@" + testFile1},
			expected: content1, // handleDataBinary preserves newlines
			desc:     "--data-binary从文件读取（保留换行符）",
		},
		{
			name:     "multiple_data_files",
			handler:  handleData,
			inputs:   []string{"@" + testFile1, "@" + testFile2},
			expected: "name=Johnage=30&city=NYCcountry=USA",
			desc:     "多个--data文件",
		},
		{
			name:     "multiple_data_binary_files",
			handler:  handleDataBinary,
			inputs:   []string{"@" + testFile1, "@" + testFile2},
			expected: content1 + content2,
			desc:     "多个--data-binary文件",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			curl := New()

			// 依次调用handler处理每个输入
			for _, input := range tt.inputs {
				err := tt.handler(curl, input)
				if err != nil {
					t.Fatalf("handler() error = %v", err)
				}
			}

			if curl.Body == nil {
				t.Error("Body is nil")
				return
			}

			bodyContent := curl.Body.String()
			if bodyContent != tt.expected {
				t.Errorf("Body content = %q, want %q", bodyContent, tt.expected)
			}
		})
	}
}

func TestMixedDataHandlers(t *testing.T) {
	// 测试混合使用不同的data选项
	curl := New()

	// 先用--data添加一些数据
	err := handleData(curl, "name=John")
	if err != nil {
		t.Fatalf("handleData() error = %v", err)
	}

	// 再用--data-raw添加数据
	err = handleDataRaw(curl, "age=30")
	if err != nil {
		t.Fatalf("handleDataRaw() error = %v", err)
	}

	expected := "name=John&age=30"
	bodyContent := curl.Body.String()
	if bodyContent != expected {
		t.Errorf("Mixed data handlers: got = %q, want = %q", bodyContent, expected)
	}

	// 验证Content-Type
	if curl.Header.Get("Content-Type") != requests.TypeURLENCODED {
		t.Errorf("Content-Type = %q, want %q", curl.Header.Get("Content-Type"), requests.TypeURLENCODED)
	}
}

func TestDataHandlerContentType(t *testing.T) {
	tests := []struct {
		name           string
		handler        func(*CURL, ...string) error
		setContentType string
		expectedType   string
		desc           string
	}{
		{
			name:         "data_default_content_type",
			handler:      handleData,
			expectedType: requests.TypeURLENCODED,
			desc:         "--data默认Content-Type",
		},
		{
			name:         "data_raw_default_content_type",
			handler:      handleDataRaw,
			expectedType: requests.TypeURLENCODED,
			desc:         "--data-raw默认Content-Type",
		},
		{
			name:           "data_preserve_content_type",
			handler:        handleData,
			setContentType: "application/json",
			expectedType:   "application/json",
			desc:           "--data保留现有Content-Type",
		},
		{
			name:           "data_binary_preserve_content_type",
			handler:        handleDataBinary,
			setContentType: "application/octet-stream",
			expectedType:   "application/octet-stream",
			desc:           "--data-binary保留现有Content-Type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			curl := New()

			// 如果需要，预设Content-Type
			if tt.setContentType != "" {
				curl.Header.Set("Content-Type", tt.setContentType)
			}

			err := tt.handler(curl, "test data")
			if err != nil {
				t.Fatalf("handler() error = %v", err)
			}

			actualType := curl.Header.Get("Content-Type")
			if actualType != tt.expectedType {
				t.Errorf("Content-Type = %q, want %q", actualType, tt.expectedType)
			}
		})
	}
}

func TestDataHandlerIntegration(t *testing.T) {
	// 测试复杂的curl命令解析
	tests := []struct {
		name    string
		curlCmd string
		check   func(*CURL) error
		desc    string
	}{
		{
			name:    "complex_form_data",
			curlCmd: `curl --data "username=john" --data "password=secret" --data "remember=true" https://example.com/login`,
			check: func(c *CURL) error {
				expected := "username=john&password=secret&remember=true"
				if c.Body.String() != expected {
					return fmt.Errorf("body = %q, want %q", c.Body.String(), expected)
				}
				if c.Header.Get("Content-Type") != requests.TypeURLENCODED {
					return fmt.Errorf("content-type = %q, want %q", c.Header.Get("Content-Type"), requests.TypeURLENCODED)
				}
				return nil
			},
			desc: "复杂表单数据",
		},
		{
			name:    "binary_data_chain",
			curlCmd: `curl --data-binary "chunk1" --data-binary "chunk2" --data-binary "chunk3" https://example.com/upload`,
			check: func(c *CURL) error {
				expected := "chunk1chunk2chunk3"
				if c.Body.String() != expected {
					return fmt.Errorf("body = %q, want %q", c.Body.String(), expected)
				}
				return nil
			},
			desc: "二进制数据链",
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
