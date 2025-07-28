package gcurl

import (
	"fmt"
	"testing"
)

// TestSpecificIssueReproduction 重现发现的具体问题
func TestSpecificIssueReproduction(t *testing.T) {
	fmt.Println("\n🔬 重现具体问题:")

	tests := []struct {
		name    string
		curlCmd string
		issue   string
	}{
		{
			name:    "HEAD方法问题",
			curlCmd: `curl -I https://httpbin.org/get`,
			issue:   "应该是HEAD方法，但解析为GET",
		},
		{
			name:    "空User-Agent问题",
			curlCmd: `curl -A "" https://httpbin.org/get`,
			issue:   "空User-Agent导致URL解析失败",
		},
		{
			name:    "无效URL问题",
			curlCmd: `curl "not-a-valid-url"`,
			issue:   "应该报错但没有",
		},
		{
			name:    "文件上传不支持",
			curlCmd: `curl -F "file=@test.txt" https://httpbin.org/post`,
			issue:   "不支持-F选项",
		},
		{
			name:    "重定向不支持",
			curlCmd: `curl -L https://httpbin.org/redirect/3`,
			issue:   "不支持-L选项",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Printf("\n🧪 测试: %s\n", tt.name)
			fmt.Printf("命令: %s\n", tt.curlCmd)
			fmt.Printf("预期问题: %s\n", tt.issue)

			cu, err := Parse(tt.curlCmd)
			if err != nil {
				fmt.Printf("❌ 解析错误: %v\n", err)
			} else {
				fmt.Printf("✓ 解析成功")
				if cu.ParsedURL != nil {
					fmt.Printf(" - URL: %s", cu.ParsedURL.String())
				}
				if cu.Method != "" {
					fmt.Printf(" - 方法: %s", cu.Method)
				}
				fmt.Println()
			}
		})
	}
}
