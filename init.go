package gcurl

import (
	"compress/flate"
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
)

// testServer 使用 sync.Once 确保线程安全的单例模式
var (
	testServerOnce sync.Once
	testServer     *http.ServeMux
)

// getTestServer 返回测试用的HTTP服务器，线程安全
func getTestServer() *http.ServeMux {
	testServerOnce.Do(func() {
		testServer = http.NewServeMux()
		testServer.HandleFunc("/get/body-compressed", func(w http.ResponseWriter, r *http.Request) {
			var writer io.Writer = w

			encodings := r.Header.Get("Accept-Encoding")
			if strings.Contains(encodings, "gzip") {
				w.Header().Set("Content-Encoding", "gzip")
				writer = gzip.NewWriter(writer)
				writer.Write([]byte("hello compress"))
				defer writer.(*gzip.Writer).Close()
			} else if strings.Contains(encodings, "deflate") {
				w.Header().Set("Content-Encoding", "deflate")
				writer, err := flate.NewWriter(writer, flate.DefaultCompression)
				if err != nil {
					panic(err)
				}
				writer.Write([]byte("hello compress"))
				defer writer.Close()
			}
			// Use r.Body to read uncompressed request body
		})
	})
	return testServer
}

func init() {
	log.SetFlags(log.Llongfile)
}
