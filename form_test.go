package gcurl

import (
	"fmt"
	"strings"
	"testing"
)

func TestFormParser(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *FormField
		wantErr  bool
	}{
		{
			name:  "Simple field",
			input: "name=John",
			expected: &FormField{
				Name:   "name",
				Value:  "John",
				IsFile: false,
			},
		},
		{
			name:  "File upload basic",
			input: "file=@test.txt",
			expected: &FormField{
				Name:     "file",
				Value:    "test.txt",
				IsFile:   true,
				Filename: "test.txt",
				MimeType: "text/plain",
			},
		},
		{
			name:  "File upload with type",
			input: "file=@image.jpg;type=image/jpeg",
			expected: &FormField{
				Name:     "file",
				Value:    "image.jpg",
				IsFile:   true,
				Filename: "image.jpg",
				MimeType: "image/jpeg",
			},
		},
		{
			name:  "File upload with filename",
			input: "file=@document.pdf;filename=report.pdf",
			expected: &FormField{
				Name:     "file",
				Value:    "document.pdf",
				IsFile:   true,
				Filename: "report.pdf",
				MimeType: "application/pdf",
			},
		},
		{
			name:  "File upload with type and filename",
			input: "upload=@data.bin;type=application/octet-stream;filename=binary.dat",
			expected: &FormField{
				Name:     "upload",
				Value:    "data.bin",
				IsFile:   true,
				Filename: "binary.dat",
				MimeType: "application/octet-stream",
			},
		},
		{
			name:    "Invalid format - no equals",
			input:   "invalidformat",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseFormData(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result.Name != tt.expected.Name {
				t.Errorf("Name: expected %q, got %q", tt.expected.Name, result.Name)
			}

			if result.Value != tt.expected.Value {
				t.Errorf("Value: expected %q, got %q", tt.expected.Value, result.Value)
			}

			if result.IsFile != tt.expected.IsFile {
				t.Errorf("IsFile: expected %v, got %v", tt.expected.IsFile, result.IsFile)
			}

			if result.IsFile {
				if result.Filename != tt.expected.Filename {
					t.Errorf("Filename: expected %q, got %q", tt.expected.Filename, result.Filename)
				}

				// 对于MIME类型，只检查主要部分，忽略charset等参数
				expectedMime := tt.expected.MimeType
				actualMime := result.MimeType
				if !strings.HasPrefix(actualMime, expectedMime) {
					t.Errorf("MimeType: expected to start with %q, got %q", expectedMime, actualMime)
				}
			}
		})
	}
}

func TestFormHandling(t *testing.T) {
	tests := []struct {
		name     string
		curlCmd  string
		validate func(*CURL) error
	}{
		{
			name:    "Single form field",
			curlCmd: `curl -F "name=John" https://httpbin.org/post`,
			validate: func(c *CURL) error {
				if c.Body == nil || c.Body.Type != "multipart" {
					return fmt.Errorf("expected multipart body")
				}

				fields, ok := c.Body.Content.([]*FormField)
				if !ok || len(fields) != 1 {
					return fmt.Errorf("expected 1 form field")
				}

				if fields[0].Name != "name" || fields[0].Value != "John" {
					return fmt.Errorf("unexpected field content")
				}

				return nil
			},
		},
		{
			name:    "Multiple form fields",
			curlCmd: `curl -F "name=John" -F "age=30" https://httpbin.org/post`,
			validate: func(c *CURL) error {
				if c.Body == nil || c.Body.Type != "multipart" {
					return fmt.Errorf("expected multipart body")
				}

				fields, ok := c.Body.Content.([]*FormField)
				if !ok || len(fields) != 2 {
					return fmt.Errorf("expected 2 form fields, got %d", len(fields))
				}

				return nil
			},
		},
		{
			name:    "File upload field",
			curlCmd: `curl -F "file=@/dev/null" https://httpbin.org/post`,
			validate: func(c *CURL) error {
				if c.Body == nil || c.Body.Type != "multipart" {
					return fmt.Errorf("expected multipart body")
				}

				fields, ok := c.Body.Content.([]*FormField)
				if !ok || len(fields) != 1 {
					return fmt.Errorf("expected 1 form field")
				}

				if !fields[0].IsFile {
					return fmt.Errorf("expected file field")
				}

				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			curl, err := Parse(tt.curlCmd)
			if err != nil {
				t.Fatalf("Failed to parse curl command: %v", err)
			}

			if err := tt.validate(curl); err != nil {
				t.Errorf("Validation failed: %v", err)
			}
		})
	}
}
