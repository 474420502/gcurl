package gcurl

import (
	"strings"
	"testing"
)

// TLS/SSL相关 handleCACert/handleClientCert/handleClientKey/handleInsecure
func TestHandleCACert(t *testing.T) {
	scurl := `curl --cacert /tmp/ca.pem https://example.com` // CA证书
	curl, err := Parse(scurl)
	if err != nil && !strings.Contains(err.Error(), "CA certificate file not found") {
		t.Fatalf("Parse failed: %v", err)
	}
	if curl != nil && curl.CACert != "/tmp/ca.pem" {
		t.Error("CA cert not parsed correctly")
	}
}

func TestHandleClientCert(t *testing.T) {
	scurl := `curl --cert /tmp/client.pem https://example.com`
	curl, err := Parse(scurl)
	if err != nil && !strings.Contains(err.Error(), "client certificate file not found") {
		t.Fatalf("Parse failed: %v", err)
	}
	if curl != nil && curl.ClientCert != "/tmp/client.pem" {
		t.Error("Client cert not parsed correctly")
	}
}

func TestHandleClientKey(t *testing.T) {
	scurl := `curl --key /tmp/client.key https://example.com`
	curl, err := Parse(scurl)
	if err != nil && !strings.Contains(err.Error(), "client key file not found") {
		t.Fatalf("Parse failed: %v", err)
	}
	if curl != nil && curl.ClientKey != "/tmp/client.key" {
		t.Error("Client key not parsed correctly")
	}
}

func TestHandleInsecure(t *testing.T) {
	scurl := `curl -k https://example.com`
	curl, err := Parse(scurl)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if !curl.Insecure {
		t.Error("Insecure not parsed correctly")
	}
}
