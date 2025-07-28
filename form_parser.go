package gcurl

import (
	"fmt"
	"mime"
	"path/filepath"
	"strings"
)

// FormField 表示一个表单字段
type FormField struct {
	Name     string // 字段名
	Value    string // 字段值
	IsFile   bool   // 是否是文件上传
	Filename string // 文件名（用于文件上传）
	MimeType string // MIME类型（用于文件上传）
}

// parseFormData 解析 curl -F 参数
// 支持以下格式：
// - name=value
// - name=@filename
// - name=@filename;type=mime/type
// - name=@filename;filename=newname
// - name=@filename;type=mime/type;filename=newname
func parseFormData(formData string) (*FormField, error) {
	// 找到第一个等号，分离字段名和值部分
	eqIndex := strings.Index(formData, "=")
	if eqIndex == -1 {
		return nil, fmt.Errorf("invalid form data format: missing '=' in %s", formData)
	}

	field := &FormField{
		Name: formData[:eqIndex],
	}

	valuePart := formData[eqIndex+1:]

	// 检查是否是文件上传（以@开头）
	if strings.HasPrefix(valuePart, "@") {
		field.IsFile = true

		// 移除@符号
		fileSpec := valuePart[1:]

		// 解析文件规格：filename[;type=mime/type][;filename=newname]
		parts := strings.Split(fileSpec, ";")
		field.Value = parts[0] // 实际文件路径

		// 设置默认文件名（基于路径）
		field.Filename = filepath.Base(field.Value)

		// 设置默认MIME类型
		field.MimeType = mime.TypeByExtension(filepath.Ext(field.Value))
		if field.MimeType == "" {
			field.MimeType = "application/octet-stream"
		}

		// 解析额外的参数
		for i := 1; i < len(parts); i++ {
			part := strings.TrimSpace(parts[i])
			if strings.HasPrefix(part, "type=") {
				field.MimeType = part[5:]
			} else if strings.HasPrefix(part, "filename=") {
				field.Filename = part[9:]
			}
		}
	} else {
		// 普通字段值
		field.IsFile = false
		field.Value = valuePart
	}

	return field, nil
}
