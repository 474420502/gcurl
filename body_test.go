package gcurl

import (
	"bytes"
	"net/url"
	"strings"
	"testing"
)

func TestBodyInterfaces(t *testing.T) {
	t.Run("RawBody", func(t *testing.T) {
		data := []byte("Hello, World!")
		body := NewRawBody(data, "text/plain")

		if body.ContentType() != "text/plain" {
			t.Errorf("期望 ContentType 为 'text/plain', 得到 '%s'", body.ContentType())
		}

		if body.Length() != int64(len(data)) {
			t.Errorf("期望 Length 为 %d, 得到 %d", len(data), body.Length())
		}

		if body.Type() != "raw" {
			t.Errorf("期望 Type 为 'raw', 得到 '%s'", body.Type())
		}

		var buf bytes.Buffer
		n, err := body.WriteTo(&buf)
		if err != nil {
			t.Errorf("WriteTo 失败: %v", err)
		}
		if n != int64(len(data)) {
			t.Errorf("期望写入 %d 字节, 实际写入 %d 字节", len(data), n)
		}
		if !bytes.Equal(buf.Bytes(), data) {
			t.Errorf("写入的数据不匹配")
		}
	})

	t.Run("FormBody", func(t *testing.T) {
		values := url.Values{}
		values.Set("name", "test")
		values.Set("age", "25")

		body := NewFormBody(values)

		if body.ContentType() != "application/x-www-form-urlencoded" {
			t.Errorf("期望 ContentType 为 'application/x-www-form-urlencoded', 得到 '%s'", body.ContentType())
		}

		if body.Type() != "form" {
			t.Errorf("期望 Type 为 'form', 得到 '%s'", body.Type())
		}

		var buf bytes.Buffer
		_, err := body.WriteTo(&buf)
		if err != nil {
			t.Errorf("WriteTo 失败: %v", err)
		}

		result := buf.String()
		if !strings.Contains(result, "name=test") || !strings.Contains(result, "age=25") {
			t.Errorf("表单数据编码不正确: %s", result)
		}
	})

	t.Run("JSONBody", func(t *testing.T) {
		data := map[string]interface{}{
			"name": "test",
			"age":  25,
		}

		body := NewJSONBody(data)

		if body.ContentType() != "application/json" {
			t.Errorf("期望 ContentType 为 'application/json', 得到 '%s'", body.ContentType())
		}

		if body.Type() != "json" {
			t.Errorf("期望 Type 为 'json', 得到 '%s'", body.Type())
		}

		var buf bytes.Buffer
		_, err := body.WriteTo(&buf)
		if err != nil {
			t.Errorf("WriteTo 失败: %v", err)
		}

		result := buf.String()
		if !strings.Contains(result, "\"name\":\"test\"") {
			t.Errorf("JSON数据编码不正确: %s", result)
		}
	})

	t.Run("TextBody", func(t *testing.T) {
		text := "Hello, World!"
		body := NewTextBody(text, "text/plain; charset=utf-8")

		if body.ContentType() != "text/plain; charset=utf-8" {
			t.Errorf("期望 ContentType 为 'text/plain; charset=utf-8', 得到 '%s'", body.ContentType())
		}

		if body.Type() != "text" {
			t.Errorf("期望 Type 为 'text', 得到 '%s'", body.Type())
		}

		if body.Length() != int64(len(text)) {
			t.Errorf("期望 Length 为 %d, 得到 %d", len(text), body.Length())
		}
	})

	t.Run("MultipartBody", func(t *testing.T) {
		fields := []*FormField{
			{Name: "name", Value: "test", IsFile: false},
			{Name: "file", Value: "content", IsFile: true, Filename: "test.txt", MimeType: "text/plain"},
		}

		body := NewMultipartBody(fields)

		if !strings.HasPrefix(body.ContentType(), "multipart/form-data; boundary=") {
			t.Errorf("期望 ContentType 以 'multipart/form-data; boundary=' 开头, 得到 '%s'", body.ContentType())
		}

		if body.Type() != "multipart" {
			t.Errorf("期望 Type 为 'multipart', 得到 '%s'", body.Type())
		}

		var buf bytes.Buffer
		_, err := body.WriteTo(&buf)
		if err != nil {
			t.Errorf("WriteTo 失败: %v", err)
		}

		result := buf.String()
		if !strings.Contains(result, "name=\"name\"") || !strings.Contains(result, "name=\"file\"") {
			t.Errorf("multipart数据格式不正确: %s", result)
		}
	})
}

func TestBodyLegacyCompatibility(t *testing.T) {
	t.Run("BodyFromLegacy", func(t *testing.T) {
		// 测试从旧的BodyData转换到新的Body接口

		// 测试raw类型
		rawData := &BodyData{
			Type:    "raw",
			Content: bytes.NewBufferString("test data"),
		}
		body := BodyFromLegacy(rawData)
		if body == nil {
			t.Error("BodyFromLegacy 返回了 nil")
		}
		if body.Type() != "raw" {
			t.Errorf("期望 Type 为 'raw', 得到 '%s'", body.Type())
		}

		// 测试form类型
		formData := &BodyData{
			Type:    "urlencoded",
			Content: "name=test&age=25",
		}
		formBody := BodyFromLegacy(formData)
		if formBody == nil {
			t.Error("BodyFromLegacy 返回了 nil for form data")
		}
		if formBody.Type() != "form" {
			t.Errorf("期望 Type 为 'form', 得到 '%s'", formBody.Type())
		}
	})

	t.Run("LegacyFromBody", func(t *testing.T) {
		// 测试从新的Body接口转换到旧的BodyData

		// 测试JSONBody
		jsonBody := NewJSONBody(map[string]string{"name": "test"})
		legacy := LegacyFromBody(jsonBody)
		if legacy == nil {
			t.Error("LegacyFromBody 返回了 nil")
		}
		if legacy.Type != "json" {
			t.Errorf("期望 Type 为 'json', 得到 '%s'", legacy.Type)
		}

		// 测试FormBody
		values := url.Values{}
		values.Set("name", "test")
		formBody := NewFormBody(values)
		formLegacy := LegacyFromBody(formBody)
		if formLegacy == nil {
			t.Error("LegacyFromBody 返回了 nil for form")
		}
		if formLegacy.Type != "urlencoded" {
			t.Errorf("期望 Type 为 'urlencoded', 得到 '%s'", formLegacy.Type)
		}
	})
}

func TestBodyDefaultValues(t *testing.T) {
	t.Run("RawBody默认ContentType", func(t *testing.T) {
		body := NewRawBody([]byte("test"), "")
		if body.ContentType() != "application/octet-stream" {
			t.Errorf("期望默认 ContentType 为 'application/octet-stream', 得到 '%s'", body.ContentType())
		}
	})

	t.Run("TextBody默认ContentType", func(t *testing.T) {
		body := NewTextBody("test", "")
		if body.ContentType() != "text/plain" {
			t.Errorf("期望默认 ContentType 为 'text/plain', 得到 '%s'", body.ContentType())
		}
	})
}

func TestBodyAdvancedFeatures(t *testing.T) {
	t.Run("JSON序列化错误处理", func(t *testing.T) {
		// 创建一个不能序列化的对象
		invalidData := make(chan int)
		body := NewJSONBody(invalidData)

		var buf bytes.Buffer
		_, err := body.WriteTo(&buf)
		if err == nil {
			t.Error("期望JSON序列化失败，但没有返回错误")
		}

		if body.Length() != 0 {
			t.Errorf("期望Length为0（序列化失败），得到 %d", body.Length())
		}
	})

	t.Run("Multipart边界生成", func(t *testing.T) {
		body1 := NewMultipartBody([]*FormField{{Name: "test", Value: "value"}})
		body2 := NewMultipartBody([]*FormField{{Name: "test", Value: "value"}})

		// 每次生成的边界应该不同
		if body1.boundary == body2.boundary {
			t.Error("多次生成的multipart边界应该不同")
		}
	})
}
