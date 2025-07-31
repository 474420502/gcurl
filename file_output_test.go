package gcurl

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
	"testing"
)

// TestFileOutputOptions 测试文件输出选项解析
func TestFileOutputOptions(t *testing.T) {
	tests := []struct {
		name        string
		command     string
		expectError bool
		checkFunc   func(*CURL) error
		description string
	}{
		{
			name:        "Output to file",
			command:     `curl -o output.txt https://httpbin.org/get`,
			expectError: false,
			checkFunc: func(c *CURL) error {
				if c.OutputFile != "output.txt" {
					return fmt.Errorf("expected OutputFile 'output.txt', got '%s'", c.OutputFile)
				}
				return nil
			},
			description: "指定输出文件名",
		},
		{
			name:        "Remote name",
			command:     `curl -O https://httpbin.org/status/200`,
			expectError: false,
			checkFunc: func(c *CURL) error {
				if !c.RemoteName {
					return fmt.Errorf("expected RemoteName to be true")
				}
				return nil
			},
			description: "使用远程文件名",
		},
		{
			name:        "Output directory",
			command:     `curl --output-dir /tmp https://httpbin.org/get`,
			expectError: false,
			checkFunc: func(c *CURL) error {
				if c.OutputDir != "/tmp" {
					return fmt.Errorf("expected OutputDir '/tmp', got '%s'", c.OutputDir)
				}
				return nil
			},
			description: "指定输出目录",
		},
		{
			name:        "Create directories",
			command:     `curl --create-dirs -o dir1/dir2/file.txt https://httpbin.org/get`,
			expectError: false,
			checkFunc: func(c *CURL) error {
				if !c.CreateDirs {
					return fmt.Errorf("expected CreateDirs to be true")
				}
				if c.OutputFile != "dir1/dir2/file.txt" {
					return fmt.Errorf("expected OutputFile 'dir1/dir2/file.txt', got '%s'", c.OutputFile)
				}
				return nil
			},
			description: "自动创建目录结构",
		},
		{
			name:        "Remove on error",
			command:     `curl --remove-on-error -o output.txt https://httpbin.org/get`,
			expectError: false,
			checkFunc: func(c *CURL) error {
				if !c.RemoveOnError {
					return fmt.Errorf("expected RemoveOnError to be true")
				}
				return nil
			},
			description: "出错时删除文件",
		},
		{
			name:        "Continue at offset",
			command:     `curl -C 1024 -o output.txt https://httpbin.org/get`,
			expectError: false,
			checkFunc: func(c *CURL) error {
				if c.ContinueAt != 1024 {
					return fmt.Errorf("expected ContinueAt 1024, got %d", c.ContinueAt)
				}
				return nil
			},
			description: "从指定字节偏移量继续下载",
		},
		{
			name:        "Continue auto",
			command:     `curl -C - -o output.txt https://httpbin.org/get`,
			expectError: false,
			checkFunc: func(c *CURL) error {
				if c.ContinueAt != -1 {
					return fmt.Errorf("expected ContinueAt -1 (auto), got %d", c.ContinueAt)
				}
				return nil
			},
			description: "自动检测断点续传",
		},
		{
			name:        "Combined options",
			command:     `curl -O --output-dir /tmp --create-dirs --remove-on-error https://httpbin.org/robots.txt`,
			expectError: false,
			checkFunc: func(c *CURL) error {
				if !c.RemoteName {
					return fmt.Errorf("expected RemoteName to be true")
				}
				if c.OutputDir != "/tmp" {
					return fmt.Errorf("expected OutputDir '/tmp', got '%s'", c.OutputDir)
				}
				if !c.CreateDirs {
					return fmt.Errorf("expected CreateDirs to be true")
				}
				if !c.RemoveOnError {
					return fmt.Errorf("expected RemoveOnError to be true")
				}
				return nil
			},
			description: "组合多个文件输出选项",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("测试命令: %s", tt.command)
			t.Logf("描述: %s", tt.description)

			curl, err := Parse(tt.command)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Parse error: %v", err)
				return
			}

			// 运行检查函数
			if err := tt.checkFunc(curl); err != nil {
				t.Error(err)
			} else {
				t.Logf("✓ 测试通过")
			}
		})
	}
}

// TestFileOutputPathDetermination 测试输出文件路径确定逻辑
func TestFileOutputPathDetermination(t *testing.T) {
	tests := []struct {
		name         string
		setupFunc    func() *CURL
		expectedPath string
		expectError  bool
		description  string
	}{
		{
			name: "Explicit output file",
			setupFunc: func() *CURL {
				c := New()
				c.OutputFile = "output.txt"
				return c
			},
			expectedPath: "output.txt",
			expectError:  false,
			description:  "明确指定输出文件",
		},
		{
			name: "Remote name with file",
			setupFunc: func() *CURL {
				c := New()
				c.RemoteName = true
				c.ParsedURL, _ = url.Parse("https://example.com/data/file.json")
				return c
			},
			expectedPath: "file.json",
			expectError:  false,
			description:  "使用远程文件名",
		},
		{
			name: "Remote name with root path",
			setupFunc: func() *CURL {
				c := New()
				c.RemoteName = true
				c.ParsedURL, _ = url.Parse("https://example.com/")
				return c
			},
			expectedPath: "index.html",
			expectError:  false,
			description:  "根路径时使用默认文件名",
		},
		{
			name: "Output directory with file",
			setupFunc: func() *CURL {
				c := New()
				c.OutputFile = "data.json"
				c.OutputDir = "/tmp/downloads"
				return c
			},
			expectedPath: "/tmp/downloads/data.json",
			expectError:  false,
			description:  "输出目录与文件名组合",
		},
		{
			name: "Remote name with output directory",
			setupFunc: func() *CURL {
				c := New()
				c.RemoteName = true
				c.OutputDir = "/tmp"
				c.ParsedURL, _ = url.Parse("https://example.com/archive.zip")
				return c
			},
			expectedPath: "/tmp/archive.zip",
			expectError:  false,
			description:  "远程文件名与输出目录组合",
		},
		{
			name: "No output specified",
			setupFunc: func() *CURL {
				c := New()
				return c
			},
			expectedPath: "",
			expectError:  false,
			description:  "未指定输出（输出到stdout）",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			curl := tt.setupFunc()

			path, err := curl.determineOutputPath()

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("determineOutputPath error: %v", err)
				return
			}

			if path != tt.expectedPath {
				t.Errorf("Expected path '%s', got '%s'", tt.expectedPath, path)
				return
			}

			t.Logf("✓ %s: '%s'", tt.description, path)
		})
	}
}

// TestFileOutputInDebug 测试调试输出中的文件输出信息
func TestFileOutputInDebug(t *testing.T) {
	curl, err := Parse(`curl -O --output-dir /tmp --create-dirs --remove-on-error https://httpbin.org/robots.txt`)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	debug := curl.Debug()
	t.Logf("Debug output:\n%s", debug)

	// 检查调试输出是否包含文件输出配置
	requiredStrings := []string{
		"File Output Configuration:",
		"Use Remote Name: YES",
		"Output Directory: /tmp",
		"Create Directories: YES",
		"Remove on Error: YES",
	}

	for _, required := range requiredStrings {
		if !strings.Contains(debug, required) {
			t.Errorf("Debug output should contain '%s'", required)
		}
	}
}

// TestFileSaveSimulation 测试文件保存的模拟（不实际创建文件）
func TestFileSaveSimulation(t *testing.T) {
	// 创建临时目录
	tmpDir := t.TempDir()

	tests := []struct {
		name        string
		setupFunc   func() *CURL
		checkFunc   func(string) error
		description string
	}{
		{
			name: "Save with output file",
			setupFunc: func() *CURL {
				c := New()
				c.OutputFile = filepath.Join(tmpDir, "test_output.txt")
				return c
			},
			checkFunc: func(expectedPath string) error {
				if !strings.Contains(expectedPath, "test_output.txt") {
					return fmt.Errorf("expected path to contain 'test_output.txt'")
				}
				return nil
			},
			description: "保存到指定文件",
		},
		{
			name: "Save with create dirs",
			setupFunc: func() *CURL {
				c := New()
				c.OutputFile = filepath.Join(tmpDir, "subdir", "nested", "file.txt")
				c.CreateDirs = true
				return c
			},
			checkFunc: func(expectedPath string) error {
				if !strings.Contains(expectedPath, "subdir/nested/file.txt") {
					return fmt.Errorf("expected nested path")
				}
				return nil
			},
			description: "创建嵌套目录结构",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			curl := tt.setupFunc()

			// 测试路径确定逻辑
			path, err := curl.determineOutputPath()
			if err != nil {
				t.Errorf("determineOutputPath error: %v", err)
				return
			}

			if err := tt.checkFunc(path); err != nil {
				t.Error(err)
				return
			}

			t.Logf("✓ %s: 路径 '%s'", tt.description, path)
		})
	}
}

// TestFileOutputErrors 测试文件输出相关的错误处理
func TestFileOutputErrors(t *testing.T) {
	tests := []struct {
		name        string
		command     string
		expectError bool
		errorCheck  func(error) bool
		description string
	}{
		{
			name:        "Invalid continue offset",
			command:     `curl -C invalid https://httpbin.org/get`,
			expectError: true,
			errorCheck: func(err error) bool {
				return strings.Contains(err.Error(), "invalid continue-at value")
			},
			description: "无效的断点续传偏移量",
		},
		{
			name:        "Negative continue offset",
			command:     `curl -C -100 https://httpbin.org/get`,
			expectError: true,
			errorCheck: func(err error) bool {
				return strings.Contains(err.Error(), "continue-at offset must be non-negative")
			},
			description: "负数断点续传偏移量",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.command)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
					return
				}

				if !tt.errorCheck(err) {
					t.Errorf("Error check failed for: %v", err)
					return
				}

				t.Logf("✓ %s: 正确处理错误 '%v'", tt.description, err)
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
					return
				}
			}
		})
	}
}

// TestFileOutputIntegration 测试文件输出功能的集成
func TestFileOutputIntegration(t *testing.T) {
	t.Run("Complete file output workflow", func(t *testing.T) {
		// 创建临时目录
		tmpDir := t.TempDir()
		outputFile := filepath.Join(tmpDir, "test_download.json")

		// 解析包含文件输出选项的命令
		curlCmd := fmt.Sprintf(`curl -o %s --create-dirs --remove-on-error https://httpbin.org/json`, outputFile)
		curl, err := Parse(curlCmd)
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}

		// 验证解析结果
		if curl.OutputFile != outputFile {
			t.Errorf("Expected OutputFile '%s', got '%s'", outputFile, curl.OutputFile)
		}

		if !curl.CreateDirs {
			t.Error("Expected CreateDirs to be true")
		}

		if !curl.RemoveOnError {
			t.Error("Expected RemoveOnError to be true")
		}

		// 测试路径确定
		path, err := curl.determineOutputPath()
		if err != nil {
			t.Errorf("determineOutputPath error: %v", err)
		}

		if path != outputFile {
			t.Errorf("Expected path '%s', got '%s'", outputFile, path)
		}

		t.Logf("✓ 完整的文件输出工作流程测试通过")
		t.Logf("  输出文件: %s", curl.OutputFile)
		t.Logf("  创建目录: %v", curl.CreateDirs)
		t.Logf("  出错删除: %v", curl.RemoveOnError)
	})
}
