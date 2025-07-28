package gcurl

import (
	"testing"
)

// TestNewFeatures 测试新添加的功能
func TestNewFeatures(t *testing.T) {
	t.Run("文件上传功能", func(t *testing.T) {
		c, err := Parse(`curl -F "name=value" -F "file=@test.txt" http://httpbin.org/post`)
		if err != nil {
			t.Fatalf("解析错误: %v", err)
		}

		if c.Method != "POST" {
			t.Errorf("期望方法是POST，得到: %s", c.Method)
		}

		contentType := c.Header.Get("Content-Type")
		if contentType != "multipart/form-data" {
			t.Errorf("期望Content-Type是multipart/form-data，得到: %s", contentType)
		}

		t.Logf("✓ 文件上传功能正常 - Method: %s, Content-Type: %s", c.Method, contentType)
	})

	t.Run("重定向跟随功能", func(t *testing.T) {
		c, err := Parse(`curl -L http://httpbin.org/redirect/3`)
		if err != nil {
			t.Fatalf("解析错误: %v", err)
		}

		if !c.FollowRedirect {
			t.Errorf("期望设置重定向标志为true，得到: %v", c.FollowRedirect)
		}

		t.Logf("✓ 重定向跟随功能正常 - URL: %s, Follow: %v", c.ParsedURL.String(), c.FollowRedirect)
	})

	t.Run("最大时间功能", func(t *testing.T) {
		c, err := Parse(`curl --max-time 30 http://httpbin.org/get`)
		if err != nil {
			t.Fatalf("解析错误: %v", err)
		}

		if c.Timeout != 30 {
			t.Errorf("期望超时时间是30秒，得到: %d", c.Timeout)
		}

		t.Logf("✓ 最大时间功能正常 - Timeout: %d秒", c.Timeout)
	})

	t.Run("代理功能", func(t *testing.T) {
		c, err := Parse(`curl --proxy http://proxy.example.com:8080 http://httpbin.org/get`)
		if err != nil {
			t.Fatalf("解析错误: %v", err)
		}

		if c.Proxy != "http://proxy.example.com:8080" {
			t.Errorf("期望代理地址是http://proxy.example.com:8080，得到: %s", c.Proxy)
		}

		t.Logf("✓ 代理功能正常 - Proxy: %s", c.Proxy)
	})

	t.Run("短选项测试", func(t *testing.T) {
		// 测试 -F 短选项
		c1, err := Parse(`curl -F "test=value" http://httpbin.org/post`)
		if err != nil {
			t.Fatalf("解析-F选项错误: %v", err)
		}
		if c1.Header.Get("Content-Type") != "multipart/form-data" {
			t.Errorf("-F选项未正确设置Content-Type")
		}

		// 测试 -L 短选项
		c2, err := Parse(`curl -L http://httpbin.org/redirect/1`)
		if err != nil {
			t.Fatalf("解析-L选项错误: %v", err)
		}
		if !c2.FollowRedirect {
			t.Errorf("-L选项未正确设置重定向标志")
		}

		// 测试 -x 短选项（proxy的别名）
		c3, err := Parse(`curl -x socks5://proxy:1080 http://httpbin.org/get`)
		if err != nil {
			t.Fatalf("解析-x选项错误: %v", err)
		}
		if c3.Proxy != "socks5://proxy:1080" {
			t.Errorf("-x选项未正确设置代理: %s", c3.Proxy)
		}

		t.Logf("✓ 所有短选项正常工作")
	})

	t.Run("多个选项组合", func(t *testing.T) {
		c, err := Parse(`curl -L --max-time 60 --proxy http://proxy:8080 -F "data=test" http://httpbin.org/post`)
		if err != nil {
			t.Fatalf("解析组合选项错误: %v", err)
		}

		// 验证所有选项都正确设置
		if c.Method != "POST" {
			t.Errorf("期望方法是POST，得到: %s", c.Method)
		}

		if c.Timeout != 60 {
			t.Errorf("期望超时时间是60秒，得到: %d", c.Timeout)
		}

		if c.Proxy != "http://proxy:8080" {
			t.Errorf("期望代理地址是http://proxy:8080，得到: %s", c.Proxy)
		}

		if !c.FollowRedirect {
			t.Errorf("期望设置重定向标志")
		}

		if c.Header.Get("Content-Type") != "multipart/form-data" {
			t.Errorf("期望Content-Type是multipart/form-data，得到: %s", c.Header.Get("Content-Type"))
		}

		t.Logf("✓ 多选项组合正常工作")
	})
}
