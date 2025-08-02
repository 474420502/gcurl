package gcurl

import (
	"testing"
)

// 跨平台兼容性测试
func TestParseWindowsCmdQuotes(t *testing.T) {
	scurl := `curl "https://example.com/api" ^
-H "accept: application/json" ^
-H "user-agent: test-agent"`
	_, err := ParseCmd(scurl)
	if err != nil {
		t.Error("Windows cmd format should parse without error")
	}
}

func TestParseBashQuotes(t *testing.T) {
	scurl := `curl 'https://example.com/api' -H 'accept: application/json' -H 'user-agent: test-agent'`
	_, err := Parse(scurl)
	if err != nil {
		t.Error("Bash format should parse without error")
	}
}
