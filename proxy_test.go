package gcurl

import (
	"testing"
)

// 代理相关 handleProxy/handleSocks5
func TestHandleProxy_HTTP(t *testing.T) {
	scurl := `curl -x http://127.0.0.1:8888 http://example.com` // http代理
	curl, err := Parse(scurl)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if curl.Proxy == "" {
		t.Error("Proxy not parsed correctly")
	}
}

func TestHandleProxy_Socks5(t *testing.T) {
	scurl := `curl --socks5 127.0.0.1:1080 http://example.com` // socks5代理
	curl, err := Parse(scurl)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if curl.Proxy == "" {
		t.Error("Socks5 proxy not parsed correctly")
	}
}

// Test handling of proxy-user short option -U
func TestHandleProxyUser_Short(t *testing.T) {
	scurl := `curl -x http://127.0.0.1:8888 -U user:pass http://example.com`
	curl, err := Parse(scurl)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if curl.ProxyUser != "user" || curl.ProxyPassword != "pass" {
		t.Errorf("Proxy user credentials not parsed correctly: got %s:%s", curl.ProxyUser, curl.ProxyPassword)
	}
}

// Test handling of proxy-user long option --proxy-user
func TestHandleProxyUser_Long(t *testing.T) {
	scurl := `curl --proxy http://127.0.0.1:8888 --proxy-user alice:secret http://example.com`
	curl, err := Parse(scurl)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if curl.ProxyUser != "alice" || curl.ProxyPassword != "secret" {
		t.Errorf("Proxy user credentials not parsed correctly: got %s:%s", curl.ProxyUser, curl.ProxyPassword)
	}
}
