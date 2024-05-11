package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"time"

	"github.com/474420502/gcurl"
)

func check(pfunc func(scurl string) (curl *gcurl.CURL, err error), r *http.Request, scurl string) {
	c, err := pfunc(scurl)
	if err != nil {
		panic(err)
	}

	if len(c.Header) != len(r.Header) {
		for k := range r.Header {
			if _, ok := c.Header[k]; !ok {
				log.Println(k, "is not exists")
			}
		}
		panic("len error")
	}
	for k, v := range r.Header {
		myHeader := c.Header[k]
		if v[0] != myHeader[0] {
			log.Println(v[0])
			log.Println(myHeader[0])
			panic("")
		}
	}
}

// handleRequest1 是处理HTTP请求的处理器函数。
func handleRequest1(w http.ResponseWriter, r *http.Request) {

	scurl := `curl 'http://localhost:7070/api-hk/heartbeat' \
	-H 'accept: application/json, text/plain, */*' \
	-H 'accept-language: zh-CN,zh;q=0.9,en;q=0.8' \
	-H 'origin: http://www.futunn.com' \
	-H 'referer: http://www.futunn.com/' \
	-H 'sec-ch-ua: "Google Chrome";v="123", "Not:A-Brand";v="8", "Chromium";v="123"' \
	-H 'sec-ch-ua-mobile: ?0' \
	-H 'sec-ch-ua-platform: "Windows"' \
	-H 'sec-fetch-dest: empty' \
	-H 'sec-fetch-mode: cors' \
	-H 'sec-fetch-site: cross-site' \
	-H 'user-agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36'`
	check(gcurl.ParseBash, r, scurl)

	scurl2 := `curl "http://localhost:7070/api-hk/heartbeat" ^
	-H "accept: application/json, text/plain, */*" ^
	-H "accept-language: zh-CN,zh;q=0.9,en;q=0.8" ^
	-H "origin: http://www.futunn.com" ^
	-H "referer: http://www.futunn.com/" ^
	-H ^"sec-ch-ua: ^\^"Google Chrome^\^";v=^\^"123^\^", ^\^"Not:A-Brand^\^";v=^\^"8^\^", ^\^"Chromium^\^";v=^\^"123^\^"^" ^
	-H "sec-ch-ua-mobile: ?0" ^
	-H ^"sec-ch-ua-platform: ^\^"Windows^\^"^" ^
	-H "sec-fetch-dest: empty" ^
	-H "sec-fetch-mode: cors" ^
	-H "sec-fetch-site: cross-site" ^
	-H "user-agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36"`
	check(gcurl.ParseCmd, r, scurl2)

	// 向客户端响应一条消息
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello, your request has been processed.")
}

// handleRequest1 是处理HTTP请求的处理器函数。
func handleRequest2(w http.ResponseWriter, r *http.Request) {

	scurl := `curl 'http://localhost:7070/shakespeare/notes/94447551/included_collections?page=1&count=7' \
	-H 'accept: application/json' \
	-H 'accept-language: zh-CN,zh;q=0.9,en;q=0.8' \
	-H 'cookie: locale=zh-CN; Hm_lvt_0c0e9d9b1e7d617b3e6842e85b9fb068=1712547029; _ga=GA1.2.273290375.1712547029; _ga_Y1EKTCT110=GS1.2.1712567179.2.0.1712567179.0.0.0; read_mode=day; default_font=font2; Hm_lpvt_0c0e9d9b1e7d617b3e6842e85b9fb068=1712567196; signin_redirect=http://localhost:7070/p/99941d7b8368; sensorsdata2015jssdkcross=%7B%22distinct_id%22%3A%2218eb693c6da528-007f4e1dd043a2-26001a51-2073600-18eb693c6db1c8c%22%2C%22first_id%22%3A%22%22%2C%22props%22%3A%7B%22%24latest_traffic_source_type%22%3A%22%E7%9B%B4%E6%8E%A5%E6%B5%81%E9%87%8F%22%2C%22%24latest_search_keyword%22%3A%22%E6%9C%AA%E5%8F%96%E5%88%B0%E5%80%BC_%E7%9B%B4%E6%8E%A5%E6%89%93%E5%BC%80%22%2C%22%24latest_referrer%22%3A%22%22%7D%2C%22%24device_id%22%3A%2218eb693c6da528-007f4e1dd043a2-26001a51-2073600-18eb693c6db1c8c%22%7D' \
	-H 'if-none-match: W/"c76a1ce3db1e3d9de4516a5cd05b8f6f"' \
	-H 'referer: http://localhost:7070/p/99941d7b8368' \
	-H 'sec-ch-ua: "Google Chrome";v="123", "Not:A-Brand";v="8", "Chromium";v="123"' \
	-H 'sec-ch-ua-mobile: ?0' \
	-H 'sec-ch-ua-platform: "Windows"' \
	-H 'sec-fetch-dest: empty' \
	-H 'sec-fetch-mode: cors' \
	-H 'sec-fetch-site: same-origin' \
	-H 'user-agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36'`
	check(gcurl.ParseBash, r, scurl)

	scurl2 := `curl "http://localhost:7070/shakespeare/notes/94447551/included_collections?page=1&count=7" ^
	-H "accept: application/json" ^
	-H "accept-language: zh-CN,zh;q=0.9,en;q=0.8" ^
	-H ^"cookie: locale=zh-CN; Hm_lvt_0c0e9d9b1e7d617b3e6842e85b9fb068=1712547029; _ga=GA1.2.273290375.1712547029; _ga_Y1EKTCT110=GS1.2.1712567179.2.0.1712567179.0.0.0; read_mode=day; default_font=font2; Hm_lpvt_0c0e9d9b1e7d617b3e6842e85b9fb068=1712567196; signin_redirect=http://localhost:7070/p/99941d7b8368; sensorsdata2015jssdkcross=^%^7B^%^22distinct_id^%^22^%^3A^%^2218eb693c6da528-007f4e1dd043a2-26001a51-2073600-18eb693c6db1c8c^%^22^%^2C^%^22first_id^%^22^%^3A^%^22^%^22^%^2C^%^22props^%^22^%^3A^%^7B^%^22^%^24latest_traffic_source_type^%^22^%^3A^%^22^%^E7^%^9B^%^B4^%^E6^%^8E^%^A5^%^E6^%^B5^%^81^%^E9^%^87^%^8F^%^22^%^2C^%^22^%^24latest_search_keyword^%^22^%^3A^%^22^%^E6^%^9C^%^AA^%^E5^%^8F^%^96^%^E5^%^88^%^B0^%^E5^%^80^%^BC_^%^E7^%^9B^%^B4^%^E6^%^8E^%^A5^%^E6^%^89^%^93^%^E5^%^BC^%^80^%^22^%^2C^%^22^%^24latest_referrer^%^22^%^3A^%^22^%^22^%^7D^%^2C^%^22^%^24device_id^%^22^%^3A^%^2218eb693c6da528-007f4e1dd043a2-26001a51-2073600-18eb693c6db1c8c^%^22^%^7D^" ^
	-H ^"if-none-match: W/^\^"c76a1ce3db1e3d9de4516a5cd05b8f6f^\^"^" ^
	-H "referer: http://localhost:7070/p/99941d7b8368" ^
	-H ^"sec-ch-ua: ^\^"Google Chrome^\^";v=^\^"123^\^", ^\^"Not:A-Brand^\^";v=^\^"8^\^", ^\^"Chromium^\^";v=^\^"123^\^"^" ^
	-H "sec-ch-ua-mobile: ?0" ^
	-H ^"sec-ch-ua-platform: ^\^"Windows^\^"^" ^
	-H "sec-fetch-dest: empty" ^
	-H "sec-fetch-mode: cors" ^
	-H "sec-fetch-site: same-origin" ^
	-H "user-agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36"`
	check(gcurl.ParseCmd, r, scurl2)

	// 向客户端响应一条消息
	w.WriteHeader(http.StatusOK)
	log.Println("Hello, your request has been processed.")
	fmt.Fprintf(w, "Hello, your request has been processed.")
}

func main() {
	// 设置监听的端口
	port := ":7070"

	// 使用http.HandleFunc注册处理函数
	http.HandleFunc("/", handleRequest2)

	// 开始监听并在给定端口上提供服务
	fmt.Printf("Starting server on %s...\n", port)
	time.AfterFunc(time.Second*1, func() {
		cmd := exec.Command("bash", "./test2.sh")
		err := cmd.Run()
		if err != nil {
			println(err.Error())
			return
		}
	})
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
