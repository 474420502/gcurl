package gcurl

import (
	"compress/flate"
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

var gserver *http.ServeMux

func init() {

	gserver = http.NewServeMux()
	gserver.HandleFunc("/get/body-compressed", func(w http.ResponseWriter, r *http.Request) {
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

}
