package gcurl

import (
	"bytes"
	"io"
	"testing"
)

func TestBodyFromLegacy_AllBranches(t *testing.T) {
	// nil
	if BodyFromLegacy(nil) != nil {
		t.Error("BodyFromLegacy(nil) should return nil")
	}
	// raw bytes.Buffer
	buf := bytes.NewBuffer([]byte("abc"))
	bd := &BodyData{Type: "raw", Content: buf}
	b := BodyFromLegacy(bd)
	if b == nil || b.Type() != "raw" {
		t.Error("raw buffer failed")
	}
	// raw string
	bd = &BodyData{Type: "raw", Content: "xyz"}
	b = BodyFromLegacy(bd)
	if b == nil || b.Type() != "text" {
		t.Error("raw string failed")
	}
	// form/urlencoded
	bd = &BodyData{Type: "form", Content: "a=1&b=2"}
	b = BodyFromLegacy(bd)
	if b == nil || b.Type() != "form" {
		t.Error("form/urlencoded failed")
	}
	// json valid
	js := `{"a":1}`
	bd = &BodyData{Type: "json", Content: js}
	b = BodyFromLegacy(bd)
	if b == nil || b.Type() != "json" {
		t.Error("json valid failed")
	}
	// json invalid
	bd = &BodyData{Type: "json", Content: "notjson"}
	b = BodyFromLegacy(bd)
	if b == nil || b.Type() != "text" {
		t.Error("json invalid fallback failed")
	}
	// multipart
	fields := []*FormField{{Name: "f", Value: "v", IsFile: false}}
	bd = &BodyData{Type: "multipart", Content: fields}
	b = BodyFromLegacy(bd)
	if b == nil || b.Type() != "multipart" {
		t.Error("multipart failed")
	}
	// unknown type
	bd = &BodyData{Type: "unknown", Content: "x"}
	b = BodyFromLegacy(bd)
	if b != nil {
		t.Error("unknown type should return nil")
	}
}

func TestLegacyFromBody_AllBranches(t *testing.T) {
	if LegacyFromBody(nil) != nil {
		t.Error("LegacyFromBody(nil) should return nil")
	}
	// raw
	raw := NewRawBody([]byte("abc"), "application/octet-stream")
	bd := LegacyFromBody(raw)
	if bd == nil || bd.Type != "raw" {
		t.Error("raw failed")
	}
	// form
	f := NewFormBody(map[string][]string{"a": {"1"}})
	bd = LegacyFromBody(f)
	if bd == nil || bd.Type != "urlencoded" {
		t.Error("form failed")
	}
	// json
	j := NewJSONBody(map[string]interface{}{"a": 1})
	bd = LegacyFromBody(j)
	if bd == nil || bd.Type != "json" {
		t.Error("json failed")
	}
	// text
	txt := NewTextBody("abc", "text/plain")
	bd = LegacyFromBody(txt)
	if bd == nil || bd.Type != "raw" {
		t.Error("text failed")
	}
	// multipart
	m := NewMultipartBody([]*FormField{{Name: "f", Value: "v", IsFile: false}})
	bd = LegacyFromBody(m)
	if bd == nil || bd.Type != "multipart" {
		t.Error("multipart failed")
	}
	// unknown type
	b := &mockBody{}
	bd = LegacyFromBody(b)
	if bd != nil {
		t.Error("unknown type should return nil")
	}
}

type mockBody struct{}

func (m *mockBody) ContentType() string                { return "mock" }
func (m *mockBody) WriteTo(w io.Writer) (int64, error) { return 0, nil }
func (m *mockBody) Len() int                           { return 0 }
func (m *mockBody) Length() int64                      { return 0 } // implement Body interface
func (m *mockBody) Type() string                       { return "mock" }

func TestRawBody_Length(t *testing.T) {
	rb := NewRawBody([]byte("abc"), "text/plain")
	if rb.Length() != 3 {
		t.Errorf("expected length 3, got %d", rb.Length())
	}
}

func TestFormBody_Length(t *testing.T) {
	fb := NewFormBody(nil)
	if fb.Length() != 0 {
		t.Errorf("expected length 0, got %d", fb.Length())
	}
}

func TestTextBody_WriteTo(t *testing.T) {
	tb := NewTextBody("hello", "text/plain")
	var buf bytes.Buffer
	n, err := tb.WriteTo(&buf)
	if err != nil || n == 0 {
		t.Error("WriteTo failed")
	}
}

func TestMultipartBody_Length(t *testing.T) {
	mb := NewMultipartBody(nil)
	if mb.Length() <= 0 {
		t.Errorf("expected length > 0, got %d", mb.Length())
	}
}

func TestBodyFromLegacy(t *testing.T) {
	_ = BodyFromLegacy(nil)
}

func TestLegacyFromBody(t *testing.T) {
	_ = LegacyFromBody(nil)
}

func TestTextBody_ContentType(t *testing.T) {
	tb := NewTextBody("hi", "text/plain")
	if tb.ContentType() != "text/plain" {
		t.Errorf("expected text/plain, got %s", tb.ContentType())
	}
}

func TestTextBody_Type(t *testing.T) {
	tb := NewTextBody("hi", "text/plain")
	if tb.Type() != "text" {
		t.Errorf("expected text, got %s", tb.Type())
	}
}

func TestTextBody_Length(t *testing.T) {
	tb := NewTextBody("hi", "text/plain")
	if tb.Length() != 2 {
		t.Errorf("expected 2, got %d", tb.Length())
	}
}

func TestTextBody_WriteTo_Empty(t *testing.T) {
	tb := NewTextBody("", "text/plain")
	var buf bytes.Buffer
	n, err := tb.WriteTo(&buf)
	if err != nil || n != 0 {
		t.Error("expected 0 bytes written for empty text body")
	}
}

func TestRawBody_WriteTo(t *testing.T) {
	rb := NewRawBody([]byte("abc"), "text/plain")
	var buf bytes.Buffer
	n, err := rb.WriteTo(&buf)
	if err != nil || n != 3 {
		t.Error("expected 3 bytes written for raw body")
	}
}

func TestFormBody_WriteTo(t *testing.T) {
	fb := NewFormBody(nil)
	var buf bytes.Buffer
	n, err := fb.WriteTo(&buf)
	if err != nil || n != 0 {
		t.Error("expected 0 bytes written for empty form body")
	}
}

func TestJSONBody_WriteTo(t *testing.T) {
	jb := NewJSONBody(map[string]string{"a": "b"})
	var buf bytes.Buffer
	n, err := jb.WriteTo(&buf)
	if err != nil || n == 0 {
		t.Error("expected non-zero bytes written for JSON body")
	}
}

func TestMultipartBody_WriteTo(t *testing.T) {
	mb := NewMultipartBody(nil)
	var buf bytes.Buffer
	n, err := mb.WriteTo(&buf)
	if err != nil {
		t.Error("WriteTo failed for multipart body")
	}
	_ = n
}

func TestGenerateBoundary(t *testing.T) {
	_ = generateBoundary()
}
