package gcurl

import (
	"fmt"
	"net/http"
	"testing"
)

func TestQuoteHandling(t *testing.T) {
	// 测试引号处理问题
	scurl := `curl 'https://example.com' -H 'sec-ch-ua: "Chromium";v="124", "Google Chrome";v="124", "Not-A.Brand";v="99"' -H 'sec-ch-ua-platform: "Windows"'`

	cu, err := Parse(scurl)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// 检查解析出的头部值
	fmt.Println("=== gcurl 解析结果 ===")
	for key, values := range cu.Header {
		for _, value := range values {
			fmt.Printf("Header %s: %s\n", key, value)
		}
	}

	// 模拟真实的 curl 行为，看看应该是什么样的
	fmt.Println("\n=== 预期的 curl 行为 ===")
	fmt.Printf("Header sec-ch-ua: %s\n", `"Chromium";v="124", "Google Chrome";v="124", "Not-A.Brand";v="99"`)
	fmt.Printf("Header sec-ch-ua-platform: %s\n", `"Windows"`)
}

func TestQuoteHandlingWithServer(t *testing.T) {
	// 创建测试服务器来检查实际收到的头部
	server := &http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Printf("=== 服务器收到的头部 ===\n")
			fmt.Printf("sec-ch-ua: %s\n", r.Header.Get("sec-ch-ua"))
			fmt.Printf("sec-ch-ua-platform: %s\n", r.Header.Get("sec-ch-ua-platform"))
			w.WriteHeader(http.StatusOK)
		}),
	}

	go server.ListenAndServe()
	defer server.Close()

	// 使用 gcurl 发送请求
	scurl := `curl 'http://localhost:8080' -H 'sec-ch-ua: "Chromium";v="124", "Google Chrome";v="124", "Not-A.Brand";v="99"' -H 'sec-ch-ua-platform: "Windows"'`

	cu, err := Parse(scurl)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	_, err = cu.Temporary().Execute()
	if err != nil {
		t.Logf("Execute failed (expected for test): %v", err)
	}
}
