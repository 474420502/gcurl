package gcurl

import (
	"bytes"
	"net/url"
	"testing"
)

func TestBodyData_WriteString(t *testing.T) {
	// nil bd
	var nilbd *BodyData
	n, err := nilbd.WriteString("abc")
	if err == nil {
		t.Error("expected error for nil BodyData")
	}
	if n != 0 {
		t.Errorf("expected 0, got %d", n)
	}

	// 非raw类型
	bd := &BodyData{Type: "json"}
	n, err = bd.WriteString("abc")
	if err == nil || n != 0 {
		t.Error("expected error for non-raw BodyData")
	}

	// raw类型但Content不是*bytes.Buffer
	bd = &BodyData{Type: "raw", Content: "notbuffer"}
	n, err = bd.WriteString("abc")
	if err == nil || n != 0 {
		t.Error("expected error for raw BodyData with wrong Content type")
	}

	// 正常raw类型
	buf := new(bytes.Buffer)
	bd = &BodyData{Type: "raw", Content: buf}
	n, err = bd.WriteString("abc")
	if err != nil {
		t.Errorf("WriteString error: %v", err)
	}
	if n != 3 {
		t.Errorf("expected 3, got %d", n)
	}
}

func TestCURL_String(t *testing.T) {
	c := New()
	c.Method = "GET"
	c.ParsedURL, _ = url.Parse("http://example.com")
	c.Header.Set("K", "V")
	c.Cookies = nil
	_ = c.String()
}

func TestCURL_SaveToFile(t *testing.T) {
	c := New()
	err := c.SaveToFile(nil)
	if err == nil {
		t.Error("expected error for nil response")
	}
}

func TestCURL_setHTTPVersionPreference(t *testing.T) {
	c := New()
	c.setHTTPVersionPreference("2")
}
