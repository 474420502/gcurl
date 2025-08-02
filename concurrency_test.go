package gcurl

import (
	"sync"
	"testing"
)

// 并发安全性测试
func TestCURLConcurrencySafety(t *testing.T) {
	curl, err := Parse(`curl -H 'X-Test: 1' http://example.com`)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			_ = curl.Header["X-Test"]
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			_ = curl.CreateRequest(nil)
		}
	}()

	wg.Wait()
}
