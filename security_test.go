package gcurl

import (
	"os"
	"testing"
)

// 安全风险测试
func TestParseDataBinaryFile(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "gcurl_test_data_bin_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())
	content := "test-binary-content-安全"
	if _, err := tmpfile.WriteString(content); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}
	tmpfile.Close()

	scurl := `curl --data-binary "@` + tmpfile.Name() + `" http://example.com`
	curl, err := Parse(scurl)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if curl.Body == nil || curl.Body.String() != content {
		t.Errorf("--data-binary file content not parsed correctly, got: %q, want: %q", curl.Body.String(), content)
	}
}

func TestParseDataCommandInjection(t *testing.T) {
	scurl := `curl --data '$(whoami)' http://example.com`
	curl, err := Parse(scurl)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if curl.Body == nil || curl.Body.String() != "$(whoami)" {
		t.Error("--data command injection not parsed as literal string, got: ", curl.Body.String())
	}
}
