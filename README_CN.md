# gcurl - Go语言 cURL 解析器

[![Go Version](https://img.shields.io/badge/Go-1.16+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Test Coverage](https://img.shields.io/badge/Coverage-79.1%25-yellow.svg)](tests)

`gcurl` 是一个强大的 Go 语言库，用于解析 cURL 命令并将其转换为结构化的 Go 数据类型。它支持大多数常用的 cURL 选项，包括认证、文件上传、输出处理等功能。

## 功能特性

- ✅ **全面的选项支持**: 支持大部分常用的 cURL 命令行选项
- ✅ **摘要认证**: 完整的 HTTP Digest Authentication 支持
- ✅ **文件处理**: 文件上传、下载和输出管理
- ✅ **协议控制**: HTTP/HTTPS 协议版本控制
- ✅ **调试模式**: 详细的调试信息输出
- ✅ **高测试覆盖率**: 79.1% 的测试覆盖率

## 快速开始

### 安装

```bash
go get github.com/474420502/gcurl
```

### 基本用法

```go
package main

import (
    "fmt"
    "log"
    "github.com/474420502/gcurl"
)

func main() {
    // 解析简单的 GET 请求
    curlCmd := `curl -X GET "https://api.example.com/users" -H "Accept: application/json"`
    
    curl, err := gcurl.ParseCurl(curlCmd)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("URL: %s\n", curl.URL)
    fmt.Printf("Method: %s\n", curl.Method)
    fmt.Printf("Headers: %v\n", curl.Headers)
}
```

### 摘要认证示例

```go
curlCmd := `curl --digest -u "user:pass" "https://example.com/protected"`
curl, err := gcurl.ParseCurl(curlCmd)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Auth Type: Digest\n")
fmt.Printf("Username: %s\n", curl.DigestAuth.Username)
```

### 文件输出示例

```go
curlCmd := `curl -o "output.html" "https://example.com"`
curl, err := gcurl.ParseCurl(curlCmd)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Output File: %s\n", curl.OutputFile)
```

## 支持的 cURL 选项

### 认证选项
- `--digest`: HTTP 摘要认证
- `-u, --user`: 用户认证信息

### 文件输出选项
- `-o, --output`: 指定输出文件
- `-O, --remote-name`: 使用远程文件名
- `--output-dir`: 输出目录
- `--create-dirs`: 创建必要的目录
- `--remove-on-error`: 错误时删除文件
- `-C, --continue-at`: 续传下载

### HTTP 选项
- `-X, --request`: 请求方法
- `-H, --header`: 自定义头部
- `-d, --data`: POST 数据
- `--data-raw`: 原始数据
- `--data-urlencode`: URL 编码数据

### 其他选项
- `-v, --verbose`: 详细输出
- `-s, --silent`: 静默模式
- `--user-agent`: 用户代理
- `--cookie`: Cookie 数据

## 测试

运行所有测试：

```bash
go test -v
```

运行特定测试：

```bash
go test -v -run TestDigestAuth
go test -v -run TestFileOutput
```

## 项目结构

```
gcurl/
├── parse_curl.go          # 核心解析逻辑
├── options.go             # 选项处理器
├── lexer.go              # 词法分析器
├── form_parser.go        # 表单解析器
├── cookie.go             # Cookie 处理
├── skip_options.go       # 跳过的选项
├── *_test.go             # 测试文件
└── examples/             # 使用示例
```

## 贡献

欢迎提交 Issue 和 Pull Request！

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 许可证

本项目基于 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 作者

- **474420502** - [GitHub](https://github.com/474420502)

## 更新日志

### v0.2.0 (当前版本)
- ✅ 添加摘要认证支持
- ✅ 实现文件输出功能
- ✅ 改进协议控制
- ✅ 增加全面的测试覆盖

### v0.1.0
- ✅ 基础 cURL 命令解析
- ✅ 基本 HTTP 选项支持
- ✅ 初始项目结构
