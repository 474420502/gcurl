package gcurl

import (
	"testing"
)

func TestRedirectHandling(t *testing.T) {
	tests := []struct {
		name            string
		curlCmd         string
		expectFollow    bool
		expectMaxRedirs int
	}{
		{
			name:            "No redirect flag",
			curlCmd:         `curl https://httpbin.org/get`,
			expectFollow:    false,
			expectMaxRedirs: -1,
		},
		{
			name:            "Location flag only",
			curlCmd:         `curl -L https://httpbin.org/redirect/3`,
			expectFollow:    true,
			expectMaxRedirs: 30, // 默认值
		},
		{
			name:            "Location with max-redirs",
			curlCmd:         `curl -L --max-redirs 5 https://httpbin.org/redirect/3`,
			expectFollow:    true,
			expectMaxRedirs: 5,
		},
		{
			name:            "Max-redirs without location flag",
			curlCmd:         `curl --max-redirs 10 https://httpbin.org/redirect/3`,
			expectFollow:    false,
			expectMaxRedirs: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			curl, err := Parse(tt.curlCmd)
			if err != nil {
				t.Fatalf("Failed to parse curl command: %v", err)
			}

			if curl.FollowRedirect != tt.expectFollow {
				t.Errorf("Expected FollowRedirect=%v, got %v", tt.expectFollow, curl.FollowRedirect)
			}

			if curl.MaxRedirs != tt.expectMaxRedirs {
				t.Errorf("Expected MaxRedirs=%d, got %d", tt.expectMaxRedirs, curl.MaxRedirs)
			}
		})
	}
}

func TestRedirectImprovement(t *testing.T) {
	// 测试新的重定向实现不再使用自定义Header
	curlCmd := `curl -L https://httpbin.org/redirect/3`
	curl, err := Parse(curlCmd)
	if err != nil {
		t.Fatalf("Failed to parse curl command: %v", err)
	}

	// 确保不再使用非标准的Header
	if curl.Header.Get("X-Gcurl-Follow-Redirects") != "" {
		t.Error("Should not use custom X-Gcurl-Follow-Redirects header")
	}

	// 确保使用了标准的重定向标志
	if !curl.FollowRedirect {
		t.Error("FollowRedirect should be true when -L is specified")
	}
}
