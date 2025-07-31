package gcurl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"time"
)

// Body 定义了统一的请求体接口，提供类型安全的Body操作
type Body interface {
	ContentType() string
	WriteTo(w io.Writer) (int64, error)
	Length() int64
	Type() string
}

// RawBody 表示原始字节数据
type RawBody struct {
	data        []byte
	contentType string
}

// NewRawBody 创建原始字节Body
func NewRawBody(data []byte, contentType string) *RawBody {
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	return &RawBody{
		data:        data,
		contentType: contentType,
	}
}

func (rb *RawBody) ContentType() string {
	return rb.contentType
}

func (rb *RawBody) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(rb.data)
	return int64(n), err
}

func (rb *RawBody) Length() int64 {
	return int64(len(rb.data))
}

func (rb *RawBody) Type() string {
	return "raw"
}

// FormBody 表示URL编码的表单数据
type FormBody struct {
	values url.Values
}

// NewFormBody 创建表单Body
func NewFormBody(values url.Values) *FormBody {
	return &FormBody{values: values}
}

func (fb *FormBody) ContentType() string {
	return "application/x-www-form-urlencoded"
}

func (fb *FormBody) WriteTo(w io.Writer) (int64, error) {
	data := fb.values.Encode()
	n, err := w.Write([]byte(data))
	return int64(n), err
}

func (fb *FormBody) Length() int64 {
	return int64(len(fb.values.Encode()))
}

func (fb *FormBody) Type() string {
	return "form"
}

// JSONBody 表示JSON数据
type JSONBody struct {
	data interface{}
}

// NewJSONBody 创建JSON Body
func NewJSONBody(data interface{}) *JSONBody {
	return &JSONBody{data: data}
}

func (jb *JSONBody) ContentType() string {
	return "application/json"
}

func (jb *JSONBody) WriteTo(w io.Writer) (int64, error) {
	jsonData, err := json.Marshal(jb.data)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal JSON: %w", err)
	}
	n, err := w.Write(jsonData)
	return int64(n), err
}

func (jb *JSONBody) Length() int64 {
	jsonData, err := json.Marshal(jb.data)
	if err != nil {
		return 0
	}
	return int64(len(jsonData))
}

func (jb *JSONBody) Type() string {
	return "json"
}

// TextBody 表示文本数据
type TextBody struct {
	text        string
	contentType string
}

// NewTextBody 创建文本Body
func NewTextBody(text, contentType string) *TextBody {
	if contentType == "" {
		contentType = "text/plain"
	}
	return &TextBody{
		text:        text,
		contentType: contentType,
	}
}

func (tb *TextBody) ContentType() string {
	return tb.contentType
}

func (tb *TextBody) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write([]byte(tb.text))
	return int64(n), err
}

func (tb *TextBody) Length() int64 {
	return int64(len(tb.text))
}

func (tb *TextBody) Type() string {
	return "text"
}

// MultipartBody 表示multipart/form-data
type MultipartBody struct {
	fields   []*FormField
	boundary string
}

// NewMultipartBody 创建multipart Body
func NewMultipartBody(fields []*FormField) *MultipartBody {
	return &MultipartBody{
		fields:   fields,
		boundary: generateBoundary(),
	}
}

func (mb *MultipartBody) ContentType() string {
	return fmt.Sprintf("multipart/form-data; boundary=%s", mb.boundary)
}

func (mb *MultipartBody) WriteTo(w io.Writer) (int64, error) {
	var totalBytes int64

	for _, field := range mb.fields {
		// 写入分隔符
		boundary := fmt.Sprintf("--%s\r\n", mb.boundary)
		n, err := w.Write([]byte(boundary))
		totalBytes += int64(n)
		if err != nil {
			return totalBytes, err
		}

		// 写入字段头部
		if field.Filename != "" {
			header := fmt.Sprintf("Content-Disposition: form-data; name=\"%s\"; filename=\"%s\"\r\n", field.Name, field.Filename)
			n, err = w.Write([]byte(header))
			totalBytes += int64(n)
			if err != nil {
				return totalBytes, err
			}

			if field.MimeType != "" {
				contentType := fmt.Sprintf("Content-Type: %s\r\n", field.MimeType)
				n, err = w.Write([]byte(contentType))
				totalBytes += int64(n)
				if err != nil {
					return totalBytes, err
				}
			}
		} else {
			header := fmt.Sprintf("Content-Disposition: form-data; name=\"%s\"\r\n", field.Name)
			n, err = w.Write([]byte(header))
			totalBytes += int64(n)
			if err != nil {
				return totalBytes, err
			}
		}

		// 空行
		n, err = w.Write([]byte("\r\n"))
		totalBytes += int64(n)
		if err != nil {
			return totalBytes, err
		}

		// 写入内容
		if field.Value != "" {
			n, err = w.Write([]byte(field.Value))
			totalBytes += int64(n)
			if err != nil {
				return totalBytes, err
			}
		}

		// 结束行
		n, err = w.Write([]byte("\r\n"))
		totalBytes += int64(n)
		if err != nil {
			return totalBytes, err
		}
	}

	// 最终分隔符
	finalBoundary := fmt.Sprintf("--%s--\r\n", mb.boundary)
	n, err := w.Write([]byte(finalBoundary))
	totalBytes += int64(n)
	if err != nil {
		return totalBytes, err
	}

	return totalBytes, nil
}

func (mb *MultipartBody) Length() int64 {
	var buf bytes.Buffer
	n, _ := mb.WriteTo(&buf)
	return n
}

func (mb *MultipartBody) Type() string {
	return "multipart"
}

// generateBoundary 生成multipart边界字符串
func generateBoundary() string {
	return fmt.Sprintf("gcurl-boundary-%d", time.Now().UnixNano())
}

// BodyFromLegacy 从旧的BodyData创建新的Body接口实现
func BodyFromLegacy(bd *BodyData) Body {
	if bd == nil {
		return nil
	}

	switch bd.Type {
	case "raw":
		if buf, ok := bd.Content.(*bytes.Buffer); ok {
			return NewRawBody(buf.Bytes(), "application/octet-stream")
		}
		if str, ok := bd.Content.(string); ok {
			return NewTextBody(str, "text/plain")
		}
	case "form", "urlencoded":
		if str, ok := bd.Content.(string); ok {
			values, _ := url.ParseQuery(str)
			return NewFormBody(values)
		}
	case "json":
		if str, ok := bd.Content.(string); ok {
			var data interface{}
			if err := json.Unmarshal([]byte(str), &data); err == nil {
				return NewJSONBody(data)
			}
			return NewTextBody(str, "application/json")
		}
	case "multipart":
		if fields, ok := bd.Content.([]*FormField); ok {
			return NewMultipartBody(fields)
		}
	}

	return nil
}

// LegacyFromBody 从新的Body接口创建旧的BodyData (向后兼容)
func LegacyFromBody(body Body) *BodyData {
	if body == nil {
		return nil
	}

	switch body.Type() {
	case "raw":
		var buf bytes.Buffer
		body.WriteTo(&buf)
		return &BodyData{
			Type:    "raw",
			Content: &buf,
		}
	case "form":
		var buf bytes.Buffer
		body.WriteTo(&buf)
		return &BodyData{
			Type:    "urlencoded",
			Content: buf.String(),
		}
	case "json":
		var buf bytes.Buffer
		body.WriteTo(&buf)
		return &BodyData{
			Type:    "json",
			Content: buf.String(),
		}
	case "text":
		var buf bytes.Buffer
		body.WriteTo(&buf)
		return &BodyData{
			Type:    "raw",
			Content: &buf,
		}
	case "multipart":
		// 这需要从MultipartBody中提取fields
		if mb, ok := body.(*MultipartBody); ok {
			return &BodyData{
				Type:    "multipart",
				Content: mb.fields,
			}
		}
	}

	return nil
}
