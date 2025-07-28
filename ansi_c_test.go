package gcurl

import (
	"strings"
	"testing"
)

func TestAnsiCQuoteProcessing(t *testing.T) {
	// 测试基本的 ANSI-C 引用处理
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "CRLF sequence",
			input:    `--data-binary $'hello\r\nworld'`,
			expected: "hello\r\nworld",
		},
		{
			name:     "LF sequence",
			input:    `--data-binary $'hello\nworld'`,
			expected: "hello\nworld",
		},
		{
			name:     "Tab sequence",
			input:    `--data-binary $'hello\tworld'`,
			expected: "hello\tworld",
		},
		{
			name:     "Escaped quote",
			input:    `--data-binary $'hello\'world'`,
			expected: "hello'world",
		},
		{
			name:     "Escaped backslash",
			input:    `--data-binary $'hello\\world'`,
			expected: "hello\\world",
		},
		{
			name:     "Regular single quotes (no ANSI-C)",
			input:    `--data-binary 'hello\r\nworld'`,
			expected: "hello\\r\\nworld", // 应保持字面量
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			curlCmd := "curl 'http://httpbin.org/post' " + tc.input

			curl, err := Parse(curlCmd)
			if err != nil {
				t.Fatalf("Failed to parse curl command: %v", err)
			}

			if curl.Body == nil {
				t.Fatal("Body is nil")
			}

			bodyContent := make([]byte, curl.Body.Len())
			curl.Body.Read(bodyContent)
			actualBody := string(bodyContent)

			if actualBody != tc.expected {
				t.Errorf("Expected body %q, got %q", tc.expected, actualBody)
				// 帮助调试：显示字节级别的差异
				t.Logf("Expected bytes: %v", []byte(tc.expected))
				t.Logf("Actual bytes:   %v", []byte(actualBody))
			}
		})
	}
}

func TestComplexMultipartData(t *testing.T) {
	// 测试复杂的 multipart/form-data
	curlCmd := `curl 'http://httpbin.org/post' \
		-H 'content-type: multipart/form-data; boundary=----WebKitFormBoundary3bCA1lzvhj4kBR4Q' \
		--data-binary $'------WebKitFormBoundary3bCA1lzvhj4kBR4Q\r\nContent-Disposition: form-data; name="keyType"\r\n\r\n0\r\n------WebKitFormBoundary3bCA1lzvhj4kBR4Q\r\nContent-Disposition: form-data; name="body"\r\n\r\n{"test":true}\r\n------WebKitFormBoundary3bCA1lzvhj4kBR4Q--\r\n'`

	curl, err := Parse(curlCmd)
	if err != nil {
		t.Fatalf("Failed to parse curl command: %v", err)
	}

	if curl.Body == nil {
		t.Fatal("Body is nil")
	}

	bodyContent := make([]byte, curl.Body.Len())
	curl.Body.Read(bodyContent)
	actualBody := string(bodyContent)

	// 验证 body 包含正确的 CRLF 序列
	if !strings.Contains(actualBody, "\r\n") {
		t.Error("Body should contain CRLF sequences")
		t.Logf("Actual body: %q", actualBody)
	}

	// 验证具体的字段存在
	if !strings.Contains(actualBody, `name="keyType"`) {
		t.Error("Body should contain keyType field")
	}
	if !strings.Contains(actualBody, `name="body"`) {
		t.Error("Body should contain body field")
	}
	if !strings.Contains(actualBody, `{"test":true}`) {
		t.Error("Body should contain JSON data")
	}
}
