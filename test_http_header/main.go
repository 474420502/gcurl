package main

import (
	"fmt"
	"net/http"

	"github.com/474420502/gcurl"
)

// handleRequest 是处理HTTP请求的处理器函数。
func handleRequest(w http.ResponseWriter, r *http.Request) {

	scurl := `curl 'https://www.futuhk.com/api-hk/heartbeat' \
	-H 'accept: application/json, text/plain, */*' \
	-H 'accept-language: zh-CN,zh;q=0.9,en;q=0.8' \
	-H 'origin: https://www.futunn.com' \
	-H 'referer: https://www.futunn.com/' \
	-H 'sec-ch-ua: "Google Chrome";v="123", "Not:A-Brand";v="8", "Chromium";v="123"' \
	-H 'sec-ch-ua-mobile: ?0' \
	-H 'sec-ch-ua-platform: "Windows"' \
	-H 'sec-fetch-dest: empty' \
	-H 'sec-fetch-mode: cors' \
	-H 'sec-fetch-site: cross-site' \
	-H 'user-agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36'`
	c, err := gcurl.Parse(scurl)
	if err != nil {
		panic(err)
	}
	if len(c.Header) != len(r.Header) {
		panic("len error")
	}
	for k, v := range r.Header {
		if v[0] != c.Header[k][0] {
			panic("")
		}
	}

	// 获取Sec-Ch-Ua-Platform头部的值
	userAgents := r.Header["Sec-Ch-Ua-Platform"]

	// 检查是否存在Sec-Ch-Ua-Platform头部
	if userAgents != nil {
		// 遍历并打印所有Sec-Ch-Ua-Platform头部的值
		for _, agent := range userAgents {
			fmt.Printf("Received Sec-Ch-Ua-Platform: %s\n", agent)
		}
	} else {
		fmt.Println("Sec-Ch-Ua-Platform header not found")
	}

	// 向客户端响应一条消息
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello, your request has been processed.")
}

func main() {
	// 设置监听的端口
	port := ":7070"

	// 使用http.HandleFunc注册处理函数
	http.HandleFunc("/", handleRequest)

	// 开始监听并在给定端口上提供服务
	fmt.Printf("Starting server on %s...\n", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
