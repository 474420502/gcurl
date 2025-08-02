package gcurl

import (
	"testing"
)

// TestResolveFeatureDemo 演示 --resolve 功能的完整使用场景
func TestResolveFeatureDemo(t *testing.T) {
	t.Log("🎯 --resolve 功能演示")
	t.Log("这是测试流程中的一个巨大痛点的解决方案!")

	// 场景1：本地开发环境
	t.Run("Local Development", func(t *testing.T) {
		t.Log("🏠 场景：本地开发环境测试")
		cmd := `curl -v https://api.production.com/health --resolve api.production.com:443:127.0.0.1`

		curl, err := Parse(cmd)
		if err != nil {
			t.Fatalf("解析失败: %v", err)
		}

		t.Logf("✓ 成功将生产域名 %s 解析到本地 127.0.0.1", curl.ParsedURL.Host)
		t.Logf("✓ 解析映射: %v", curl.Resolve)

		// 显示详细输出
		verbose := curl.VerboseInfo()
		t.Logf("🔍 详细输出:\n%s", verbose)
	})

	// 场景2：负载均衡测试
	t.Run("Load Balancing Test", func(t *testing.T) {
		t.Log("⚖️ 场景：负载均衡多后端测试")
		cmd := `curl https://service.example.com/status --resolve service.example.com:443:10.0.1.100,10.0.1.101,10.0.1.102`

		curl, err := Parse(cmd)
		if err != nil {
			t.Fatalf("解析失败: %v", err)
		}

		t.Logf("✓ 负载均衡测试配置完成")
		t.Logf("✓ 后端服务器: %s", curl.Resolve[0])

		// 验证包含多个IP地址
		resolveEntry := curl.Resolve[0]
		if !findSubstring(resolveEntry, "10.0.1.100") || !findSubstring(resolveEntry, "10.0.1.101") {
			t.Error("应该包含多个后端IP地址")
		} else {
			t.Log("✓ 多后端IP配置正确")
		}
	})

	// 场景3：预发布环境测试
	t.Run("Staging Environment", func(t *testing.T) {
		t.Log("🚀 场景：预发布环境验证")
		cmd := `curl -H "X-Environment: staging" https://api.myapp.com/version --resolve api.myapp.com:443:staging-server.internal.com`

		curl, err := Parse(cmd)
		if err != nil {
			t.Fatalf("解析失败: %v", err)
		}

		t.Logf("✓ 预发布环境配置完成")
		t.Logf("✓ 目标服务器: %s", curl.Resolve[0])
		t.Logf("✓ 环境标识头部: %s", curl.Header.Get("X-Environment"))
	})

	// 场景4：强制解析覆盖
	t.Run("Force Resolution Override", func(t *testing.T) {
		t.Log("🔄 场景：强制覆盖DNS解析")
		cmd := `curl --resolve +problematic.service.com:443:127.0.0.1 https://problematic.service.com/debug`

		curl, err := Parse(cmd)
		if err != nil {
			t.Fatalf("解析失败: %v", err)
		}

		t.Logf("✓ 强制解析覆盖配置完成")
		if !findSubstring(curl.Resolve[0], "+problematic.service.com") {
			t.Error("应该包含强制覆盖标记 '+'")
		} else {
			t.Log("✓ 强制覆盖标记正确")
		}
	})

	// 场景5：多端口服务测试
	t.Run("Multi-Port Service", func(t *testing.T) {
		t.Log("🌐 场景：多端口服务分别测试")
		cmd := `curl https://service.com/api --resolve service.com:80:192.168.1.10 --resolve service.com:443:192.168.1.11`

		curl, err := Parse(cmd)
		if err != nil {
			t.Fatalf("解析失败: %v", err)
		}

		t.Logf("✓ 多端口解析配置完成")
		t.Logf("✓ HTTP (80): %s", curl.Resolve[0])
		t.Logf("✓ HTTPS (443): %s", curl.Resolve[1])

		if len(curl.Resolve) != 2 {
			t.Error("应该有两个解析条目")
		} else {
			t.Log("✓ 双端口配置正确")
		}
	})

	t.Log("🎉 --resolve 功能演示完成！")
	t.Log("这个功能极大地简化了开发和测试流程：")
	t.Log("  • 本地开发时无需修改 /etc/hosts")
	t.Log("  • 负载均衡测试轻松配置")
	t.Log("  • 预发布环境验证更简单")
	t.Log("  • 问题排查时快速重定向")
	t.Log("  • 与 -v 选项完美集成，清晰显示解析信息")
}

// findSubstring 辅助函数
func findSubstring(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
