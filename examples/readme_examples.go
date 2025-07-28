package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/474420502/gcurl"
)

// Example from README.md - Example 1: Basic GET request with headers
func exampleBasicGET() {
	fmt.Println("=== Example 1: Basic GET Request ===")

	surl := `http://httpbin.org/get -H 'Connection: keep-alive' -H 'Accept-Encoding: gzip, deflate' -H 'Accept-Language: zh-CN,zh;q=0.9'`
	curl, err := gcurl.Parse(surl)
	if err != nil {
		log.Fatal(err)
	}

	ses := curl.CreateSession()
	tp := curl.CreateTemporary(ses)

	fmt.Println("Headers:", ses.GetHeader())

	resp, err := tp.Execute()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Response:", string(resp.Content()))
	fmt.Println()
}

// Example from README.md - Example 8: Custom Session Configuration
func exampleCustomSession() {
	fmt.Println("=== Example 8: Custom Session Configuration ===")

	curl, err := gcurl.Parse(`curl "http://httpbin.org/headers"`)
	if err != nil {
		log.Fatal(err)
	}

	ses := curl.CreateSession()

	customHeaders := make(http.Header)
	customHeaders.Set("X-Custom-Header", "MyValue")
	customHeaders.Set("User-Agent", "MyApp/1.0")
	ses.SetHeader(customHeaders)

	ses.Config().SetTimeout(5)

	tp := curl.CreateTemporary(ses)

	resp, err := tp.Execute()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Status: %d\n", resp.GetStatusCode())
	fmt.Println("Response:", string(resp.Content()))
	fmt.Println()
}

// Example from README.md - Example 3: POST Request with JSON Data
func examplePostJSON() {
	fmt.Println("=== Example 3: POST Request with JSON Data ===")

	scurl := `curl -X POST "http://httpbin.org/post" -H "Content-Type: application/json" -d '{"name":"test","age":25}'`

	curl, err := gcurl.Parse(scurl)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Method: %s\n", curl.Method)
	fmt.Printf("Content-Type: %s\n", curl.ContentType)

	resp, err := curl.Temporary().Execute()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Response:", string(resp.Content()))
	fmt.Println()
}

func main() {
	fmt.Println("gcurl README Examples Demonstration")
	fmt.Println("=====================================")

	exampleBasicGET()
	exampleCustomSession()
	examplePostJSON()

	fmt.Println("All examples completed successfully! 🎉")
}
