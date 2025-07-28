package gcurl

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// TestComprehensiveCurlFeatures 测试更多curl功能
func TestComprehensiveCurlFeatures(t *testing.T) {
	tests := []struct {
		name        string
		curlCmd     string
		expectError bool
		description string
	}{
		// HTTP方法测试
		{
			name:        "POST with data",
			curlCmd:     `curl -X POST -d "key=value&foo=bar" https://httpbin.org/post`,
			expectError: false,
			description: "POST请求带数据",
		},
		{
			name:        "PUT with data",
			curlCmd:     `curl -X PUT -d '{"name":"test"}' -H "Content-Type: application/json" https://httpbin.org/put`,
			expectError: false,
			description: "PUT请求带JSON数据",
		},
		{
			name:        "DELETE method",
			curlCmd:     `curl -X DELETE https://httpbin.org/delete`,
			expectError: false,
			description: "DELETE请求",
		},
		{
			name:        "PATCH method",
			curlCmd:     `curl -X PATCH -d '{"status":"updated"}' https://httpbin.org/patch`,
			expectError: false,
			description: "PATCH请求",
		},
		{
			name:        "HEAD method",
			curlCmd:     `curl -I https://httpbin.org/get`,
			expectError: false,
			description: "HEAD请求",
		},
		{
			name:        "OPTIONS method",
			curlCmd:     `curl -X OPTIONS https://httpbin.org/get`,
			expectError: false,
			description: "OPTIONS请求",
		},

		// 数据传输测试
		{
			name:        "Form data multiple",
			curlCmd:     `curl -d "name=John&age=30&city=New York" https://httpbin.org/post`,
			expectError: false,
			description: "表单数据多个字段",
		},
		{
			name:        "JSON data complex",
			curlCmd:     `curl -d '{"user":{"name":"John","details":{"age":30,"skills":["Go","Python"]}}}' -H "Content-Type: application/json" https://httpbin.org/post`,
			expectError: false,
			description: "复杂JSON数据",
		},
		{
			name:        "Raw binary data",
			curlCmd:     `curl --data-binary @/dev/null https://httpbin.org/post`,
			expectError: false,
			description: "二进制数据",
		},
		{
			name:        "URL encoded data",
			curlCmd:     `curl --data-urlencode "message=Hello World! @#$%^&*()" https://httpbin.org/post`,
			expectError: false,
			description: "URL编码数据",
		},

		// 文件上传测试
		{
			name:        "File upload",
			curlCmd:     `curl -F "file=@test.txt" -F "name=upload" https://httpbin.org/post`,
			expectError: false,
			description: "文件上传",
		},
		{
			name:        "Multiple files",
			curlCmd:     `curl -F "file1=@test1.txt" -F "file2=@test2.txt" https://httpbin.org/post`,
			expectError: false,
			description: "多文件上传",
		},

		// 认证测试
		{
			name:        "Basic auth",
			curlCmd:     `curl -u "username:password" https://httpbin.org/basic-auth/username/password`,
			expectError: false,
			description: "基本认证",
		},
		{
			name:        "Bearer token",
			curlCmd:     `curl -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9" https://httpbin.org/bearer`,
			expectError: false,
			description: "Bearer令牌",
		},
		{
			name:        "API key header",
			curlCmd:     `curl -H "X-API-Key: abc123def456" https://httpbin.org/get`,
			expectError: false,
			description: "API密钥头部",
		},

		// 复杂头部测试
		{
			name:        "Multiple custom headers",
			curlCmd:     `curl -H "X-Custom-Header: value1" -H "X-Another-Header: value2" -H "X-Third-Header: value3" https://httpbin.org/get`,
			expectError: false,
			description: "多个自定义头部",
		},
		{
			name:        "Headers with special chars",
			curlCmd:     `curl -H "X-Special: !@#$%^&*()_+-=[]{}|;':\",./<>?" https://httpbin.org/get`,
			expectError: false,
			description: "包含特殊字符的头部",
		},
		{
			name:        "Content negotiation",
			curlCmd:     `curl -H "Accept: application/json, application/xml;q=0.9, text/plain;q=0.8, */*;q=0.1" https://httpbin.org/get`,
			expectError: false,
			description: "内容协商",
		},
		{
			name:        "Language headers",
			curlCmd:     `curl -H "Accept-Language: zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7,ja;q=0.6" https://httpbin.org/get`,
			expectError: false,
			description: "语言头部",
		},

		// Cookie测试
		{
			name:        "Simple cookie",
			curlCmd:     `curl -b "session=abc123" https://httpbin.org/cookies`,
			expectError: false,
			description: "简单Cookie",
		},
		{
			name:        "Multiple cookies",
			curlCmd:     `curl -b "session=abc123; user=john; theme=dark" https://httpbin.org/cookies`,
			expectError: false,
			description: "多个Cookie",
		},
		{
			name:        "Complex cookie values",
			curlCmd:     `curl -b "data={\"user\":\"john\",\"id\":123}; token=eyJhbGciOiJIUzI1NiJ9" https://httpbin.org/cookies`,
			expectError: false,
			description: "复杂Cookie值",
		},

		// 重定向测试
		{
			name:        "Follow redirects",
			curlCmd:     `curl -L https://httpbin.org/redirect/3`,
			expectError: false,
			description: "跟随重定向",
		},
		{
			name:        "Max redirects",
			curlCmd:     `curl -L --max-redirs 5 https://httpbin.org/redirect/3`,
			expectError: false,
			description: "最大重定向次数",
		},

		// 超时测试
		{
			name:        "Connection timeout",
			curlCmd:     `curl --connect-timeout 10 https://httpbin.org/delay/2`,
			expectError: false,
			description: "连接超时",
		},
		{
			name:        "Max time",
			curlCmd:     `curl --max-time 30 https://httpbin.org/delay/1`,
			expectError: false,
			description: "最大执行时间",
		},

		// User-Agent测试
		{
			name:        "Custom user agent",
			curlCmd:     `curl -A "MyApp/1.0 (Linux; Android 10)" https://httpbin.org/user-agent`,
			expectError: false,
			description: "自定义User-Agent",
		},
		{
			name:        "Empty user agent",
			curlCmd:     `curl -A "" https://httpbin.org/user-agent`,
			expectError: false,
			description: "空User-Agent",
		},

		// 代理测试
		{
			name:        "HTTP proxy",
			curlCmd:     `curl --proxy http://proxy.example.com:8080 https://httpbin.org/get`,
			expectError: false,
			description: "HTTP代理",
		},
		{
			name:        "SOCKS proxy",
			curlCmd:     `curl --socks5 socks5://127.0.0.1:1080 https://httpbin.org/get`,
			expectError: false,
			description: "SOCKS代理",
		},

		// SSL/TLS测试
		{
			name:        "Skip SSL verification",
			curlCmd:     `curl -k https://self-signed.badssl.com/`,
			expectError: false,
			description: "跳过SSL验证",
		},
		{
			name:        "Custom CA cert",
			curlCmd:     `curl --cacert /dev/null https://httpbin.org/get`,
			expectError: false,
			description: "自定义CA证书",
		},
		{
			name:        "Client certificate",
			curlCmd:     `curl --cert /dev/null --key /dev/null https://httpbin.org/get`,
			expectError: false,
			description: "客户端证书",
		},

		// 压缩测试
		{
			name:        "Accept compression",
			curlCmd:     `curl --compressed https://httpbin.org/gzip`,
			expectError: false,
			description: "接受压缩",
		},

		// 范围请求测试
		{
			name:        "Range request",
			curlCmd:     `curl -H "Range: bytes=0-1023" https://httpbin.org/range/2048`,
			expectError: false,
			description: "范围请求",
		},

		// 条件请求测试
		{
			name:        "If-Modified-Since",
			curlCmd:     `curl -H "If-Modified-Since: Wed, 21 Oct 2015 07:28:00 GMT" https://httpbin.org/get`,
			expectError: false,
			description: "条件请求If-Modified-Since",
		},
		{
			name:        "If-None-Match",
			curlCmd:     `curl -H "If-None-Match: \"686897696a7c876b7e\"" https://httpbin.org/etag/test`,
			expectError: false,
			description: "条件请求If-None-Match",
		},

		// 错误处理测试
		{
			name:        "Invalid URL",
			curlCmd:     `curl "not-a-valid-url"`,
			expectError: true,
			description: "无效URL",
		},
		{
			name:        "Unsupported protocol",
			curlCmd:     `curl ftp://example.com/file.txt`,
			expectError: false, // gcurl可能支持或忽略
			description: "不支持的协议",
		},

		// 复杂查询参数测试
		{
			name:        "Complex query params",
			curlCmd:     `curl "https://httpbin.org/get?q=search%20term&limit=10&offset=0&sort=created_at&order=desc&filters[]=active&filters[]=verified"`,
			expectError: false,
			description: "复杂查询参数",
		},

		// 国际化测试
		{
			name:        "Unicode in URL",
			curlCmd:     `curl "https://httpbin.org/get?message=你好世界&emoji=🚀"`,
			expectError: false,
			description: "URL中包含Unicode",
		},
		{
			name:        "Unicode in headers",
			curlCmd:     `curl -H "X-Message: 你好世界 🌍" https://httpbin.org/get`,
			expectError: false,
			description: "头部中包含Unicode",
		},

		// 超长数据测试
		{
			name:        "Large header",
			curlCmd:     fmt.Sprintf(`curl -H "X-Large-Header: %s" https://httpbin.org/get`, strings.Repeat("a", 8192)),
			expectError: false,
			description: "超长头部",
		},
		{
			name:        "Large POST data",
			curlCmd:     fmt.Sprintf(`curl -d "%s" https://httpbin.org/post`, strings.Repeat("data", 1000)),
			expectError: false,
			description: "大量POST数据",
		},

		// 边界情况测试
		{
			name:        "Empty header value",
			curlCmd:     `curl -H "X-Empty:" https://httpbin.org/get`,
			expectError: false,
			description: "空头部值",
		},
		{
			name:        "Header with only spaces",
			curlCmd:     `curl -H "X-Spaces:     " https://httpbin.org/get`,
			expectError: false,
			description: "只包含空格的头部值",
		},
		{
			name:        "Multiple same headers",
			curlCmd:     `curl -H "X-Test: value1" -H "X-Test: value2" -H "X-Test: value3" https://httpbin.org/get`,
			expectError: false,
			description: "多个相同名称的头部",
		},

		// HTTP/2 和 HTTP/3 测试
		{
			name:        "HTTP/2",
			curlCmd:     `curl --http2 https://httpbin.org/get`,
			expectError: false,
			description: "强制使用HTTP/2",
		},

		// WebSocket升级测试（虽然curl不直接支持WebSocket，但可以测试升级头部）
		{
			name:        "WebSocket upgrade headers",
			curlCmd:     `curl -H "Upgrade: websocket" -H "Connection: Upgrade" -H "Sec-WebSocket-Key: dGhlIHNhbXBsZSBub25jZQ==" -H "Sec-WebSocket-Version: 13" https://httpbin.org/get`,
			expectError: false,
			description: "WebSocket升级头部",
		},

		// GraphQL测试
		{
			name:        "GraphQL query",
			curlCmd:     `curl -X POST -H "Content-Type: application/json" -d '{"query":"query { user(id: 1) { name email } }","variables":{"id":1}}' https://httpbin.org/post`,
			expectError: false,
			description: "GraphQL查询",
		},

		// CORS预检请求测试
		{
			name:        "CORS preflight",
			curlCmd:     `curl -X OPTIONS -H "Origin: https://example.com" -H "Access-Control-Request-Method: POST" -H "Access-Control-Request-Headers: X-Custom-Header" https://httpbin.org/post`,
			expectError: false,
			description: "CORS预检请求",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Printf("Testing: %s - %s\n", tt.name, tt.description)

			cu, err := Parse(tt.curlCmd)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for %s, but got none", tt.name)
				} else {
					fmt.Printf("✓ Expected error: %v\n", err)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error for %s: %v", tt.name, err)
				return
			}

			// 基本验证
			if cu.ParsedURL == nil {
				t.Errorf("URL not parsed for %s", tt.name)
				return
			}

			fmt.Printf("✓ Parsed URL: %s\n", cu.ParsedURL.String())
			fmt.Printf("✓ Method: %s\n", cu.Method)

			if len(cu.Header) > 0 {
				fmt.Printf("✓ Headers (%d):\n", len(cu.Header))
				for key, values := range cu.Header {
					for _, value := range values {
						fmt.Printf("  %s: %s\n", key, value)
					}
				}
			}

			if cu.Body != nil && cu.Body.Len() > 0 {
				fmt.Printf("✓ Body length: %d bytes\n", cu.Body.Len())
			}

			if len(cu.Cookies) > 0 {
				fmt.Printf("✓ Cookies (%d)\n", len(cu.Cookies))
			}

			fmt.Println()
		})
	}
}

// TestStressAndEdgeCases 压力测试和边缘情况
func TestStressAndEdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		curlCmd string
		desc    string
	}{
		{
			name:    "Extremely long URL",
			curlCmd: fmt.Sprintf(`curl "https://httpbin.org/get?%s"`, strings.Repeat("param=value&", 1000)),
			desc:    "超长URL测试",
		},
		{
			name:    "Many headers",
			curlCmd: generateCurlWithManyHeaders(100),
			desc:    "大量头部测试",
		},
		{
			name:    "Nested quotes",
			curlCmd: `curl -H 'X-JSON: {"message": "He said \"Hello World!\" to me"}' https://httpbin.org/get`,
			desc:    "嵌套引号测试",
		},
		{
			name:    "Escaped characters",
			curlCmd: `curl -d "message=Line1\nLine2\tTabbed\r\nCRLF" https://httpbin.org/post`,
			desc:    "转义字符测试",
		},
		{
			name:    "Binary data in URL",
			curlCmd: `curl "https://httpbin.org/get?binary=%00%01%02%FF"`,
			desc:    "URL中的二进制数据",
		},
		{
			name:    "Multiple protocols",
			curlCmd: `curl "https://user:pass@httpbin.org:443/get?redirect=http://example.com"`,
			desc:    "多协议URL",
		},
		{
			name:    "IPv6 URL",
			curlCmd: `curl "http://[::1]:8080/test"`,
			desc:    "IPv6地址",
		},
		{
			name:    "Punycode domain",
			curlCmd: `curl "https://xn--nxasmq6b.xn--o3cw4h/test"`,
			desc:    "国际化域名",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Printf("Testing edge case: %s - %s\n", tt.name, tt.desc)

			cu, err := Parse(tt.curlCmd)
			if err != nil {
				fmt.Printf("❌ Parse error: %v\n", err)
				return
			}

			fmt.Printf("✓ Successfully parsed complex case\n")
			if cu.ParsedURL != nil {
				fmt.Printf("✓ URL: %s\n", cu.ParsedURL.String())
			}
			fmt.Println()
		})
	}
}

// TestRealWorldScenarios 真实世界场景测试
func TestRealWorldScenarios(t *testing.T) {
	scenarios := []struct {
		name    string
		curlCmd string
		desc    string
	}{
		{
			name: "GitHub API",
			curlCmd: `curl -H "Accept: application/vnd.github.v3+json" \
				-H "Authorization: token ghp_xxxxxxxxxxxxxxxxxxxx" \
				"https://api.github.com/user/repos?type=private&sort=updated&per_page=50"`,
			desc: "GitHub API调用",
		},
		{
			name: "Docker Registry",
			curlCmd: `curl -H "Accept: application/vnd.docker.distribution.manifest.v2+json" \
				-H "Authorization: Bearer token" \
				"https://registry.hub.docker.com/v2/library/nginx/manifests/latest"`,
			desc: "Docker Registry API",
		},
		{
			name: "Elasticsearch Query",
			curlCmd: `curl -X POST "http://localhost:9200/logs/_search?pretty" \
				-H "Content-Type: application/json" \
				-d '{"query":{"bool":{"must":[{"match":{"level":"ERROR"}},{"range":{"@timestamp":{"gte":"2023-01-01","lte":"2023-12-31"}}}]}},"size":100}'`,
			desc: "Elasticsearch查询",
		},
		{
			name: "AWS S3 API",
			curlCmd: `curl -X PUT "https://mybucket.s3.amazonaws.com/myfile.txt" \
				-H "Host: mybucket.s3.amazonaws.com" \
				-H "Date: Wed, 01 Mar 2023 12:00:00 GMT" \
				-H "Authorization: AWS4-HMAC-SHA256 Credential=AKIAIOSFODNN7EXAMPLE/20230301/us-east-1/s3/aws4_request, SignedHeaders=host;range;x-amz-date, Signature=signature" \
				-H "x-amz-content-sha256: UNSIGNED-PAYLOAD" \
				--data-binary @file.txt`,
			desc: "AWS S3上传",
		},
		{
			name: "Webhook with signature",
			curlCmd: `curl -X POST "https://api.example.com/webhook" \
				-H "Content-Type: application/json" \
				-H "X-Hub-Signature-256: sha256=1234567890abcdef" \
				-H "X-GitHub-Event: push" \
				-H "X-GitHub-Delivery: 12345678-1234-1234-1234-123456789abc" \
				-d '{"ref":"refs/heads/main","commits":[{"id":"abc123","message":"Update README"}]}'`,
			desc: "带签名的Webhook",
		},
		{
			name: "OAuth token exchange",
			curlCmd: `curl -X POST "https://oauth2.googleapis.com/token" \
				-H "Content-Type: application/x-www-form-urlencoded" \
				-d "grant_type=authorization_code&code=4/P7q7W91a-oMsCeLvIaQm6bTrgtp7&client_id=your_client_id&client_secret=your_client_secret&redirect_uri=https://oauth2.example.com/code"`,
			desc: "OAuth令牌交换",
		},
		{
			name: "Kubernetes API",
			curlCmd: `curl -X GET "https://kubernetes.default.svc/api/v1/namespaces/default/pods" \
				-H "Authorization: Bearer $(cat /var/run/secrets/kubernetes.io/serviceaccount/token)" \
				-H "Accept: application/json" \
				--cacert /var/run/secrets/kubernetes.io/serviceaccount/ca.crt`,
			desc: "Kubernetes API调用",
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			fmt.Printf("Testing real-world scenario: %s - %s\n", scenario.name, scenario.desc)

			cu, err := Parse(scenario.curlCmd)
			if err != nil {
				fmt.Printf("❌ Parse error: %v\n", err)
				return
			}

			fmt.Printf("✓ Successfully parsed real-world scenario\n")
			fmt.Printf("✓ URL: %s\n", cu.ParsedURL.String())
			fmt.Printf("✓ Method: %s\n", cu.Method)
			fmt.Printf("✓ Headers: %d\n", len(cu.Header))
			if cu.Body != nil && cu.Body.Len() > 0 {
				fmt.Printf("✓ Has body data\n")
			}
			fmt.Println()
		})
	}
}

// TestWithLocalServer 使用本地服务器测试实际HTTP请求
func TestWithLocalServer(t *testing.T) {
	// 创建测试服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 记录请求详情
		fmt.Printf("=== Server received ===\n")
		fmt.Printf("Method: %s\n", r.Method)
		fmt.Printf("URL: %s\n", r.URL.String())
		fmt.Printf("Headers:\n")
		for name, values := range r.Header {
			for _, value := range values {
				fmt.Printf("  %s: %s\n", name, value)
			}
		}

		// 读取body
		body := make([]byte, r.ContentLength)
		r.Body.Read(body)
		if len(body) > 0 {
			fmt.Printf("Body: %s\n", string(body))
		}

		// 响应
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok", "received": true}`))
	}))
	defer server.Close()

	tests := []string{
		fmt.Sprintf(`curl "%s/test"`, server.URL),
		fmt.Sprintf(`curl -X POST "%s/post" -d "key=value"`, server.URL),
		fmt.Sprintf(`curl -H "Custom-Header: test-value" "%s/headers"`, server.URL),
		fmt.Sprintf(`curl -b "session=abc123" "%s/cookies"`, server.URL),
		fmt.Sprintf(`curl -X PUT -H "Content-Type: application/json" -d '{"test": true}' "%s/json"`, server.URL),
	}

	for i, curlCmd := range tests {
		t.Run(fmt.Sprintf("LocalServer_%d", i+1), func(t *testing.T) {
			fmt.Printf("Testing with local server: %s\n", curlCmd)

			cu, err := Parse(curlCmd)
			if err != nil {
				t.Errorf("Parse error: %v", err)
				return
			}

			// 尝试执行请求
			resp, err := cu.Temporary().Execute()
			if err != nil {
				fmt.Printf("❌ Execute error: %v\n", err)
				return
			}

			fmt.Printf("✓ Response Status: %d\n", resp.GetStatusCode())
			fmt.Printf("✓ Response Body: %s\n", resp.ContentString())
			fmt.Println()
		})
	}
}

// 辅助函数：生成包含大量头部的curl命令
func generateCurlWithManyHeaders(count int) string {
	var headers []string
	for i := 0; i < count; i++ {
		headers = append(headers, fmt.Sprintf(`-H "X-Header-%d: value-%d"`, i, i))
	}
	return fmt.Sprintf(`curl %s https://httpbin.org/get`, strings.Join(headers, " "))
}

// TestPerformanceImpact 性能影响测试
func TestPerformanceImpact(t *testing.T) {
	// 测试解析大量curl命令的性能
	commands := []string{
		`curl "https://httpbin.org/get"`,
		`curl -X POST -d "data=test" "https://httpbin.org/post"`,
		`curl -H "Authorization: Bearer token" "https://httpbin.org/bearer"`,
		`curl -b "session=abc123" "https://httpbin.org/cookies"`,
		`curl -X PUT -H "Content-Type: application/json" -d '{"key":"value"}' "https://httpbin.org/put"`,
	}

	start := time.Now()
	for i := 0; i < 1000; i++ {
		for _, cmd := range commands {
			_, err := Parse(cmd)
			if err != nil {
				t.Errorf("Parse error at iteration %d: %v", i, err)
				return
			}
		}
	}
	duration := time.Since(start)

	fmt.Printf("Performance test: Parsed %d commands in %v\n", 1000*len(commands), duration)
	fmt.Printf("Average time per command: %v\n", duration/time.Duration(1000*len(commands)))
}
