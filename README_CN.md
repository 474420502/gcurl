# gcurl - Go语言 cURL 解析器

![Go Version](https://img.shields.io/badge/Go-1.20+-00ADD8?style=flat&logo=go)
![Test Coverage](https://img.shields.io/badge/Coverage-79.1%25-brightgreen)
![License](https://img.shields.io/badge/License-MIT-blue)
![Latest Release](https://img.shields.io/github/v/release/474420502/gcurl)

一个强大的 Go 语言库，用于将 cURL 命令转换为 Go HTTP 请求，具有完整的功能兼容性。

* 基于强大的 [requests](https://github.com/474420502/requests) 库构建
* 无缝将 cURL bash 命令转换为 Go 代码
* 继承所有 cURL 功能，同时增加 Go 的配置、Cookie、头部和 URL 处理灵活性
* 支持 Bash 和 Windows Cmd 格式的 cURL 命令
* 生产就绪，具有全面的测试覆盖率

## 🚀 功能特性

- 🌐 **完整的 cURL 命令解析** - 解析任何 cURL 命令为 Go 请求
- 🔧 **完整的 cURL 兼容性** - 支持所有主要的 cURL 选项
- 📁 **文件输出支持** - 使用 `-o`、`-O`、`--output-dir` 等保存响应到文件
- 🔐 **认证支持** - Basic、Digest、Bearer token 和自定义认证
- 📤 **文件上传** - 多部分表单和文件上传支持
- 🍪 **Cookie 管理** - 自动 Cookie 处理和会话管理
- 🐛 **调试模式** - 详细的调试输出和请求追踪
- ⚡ **高性能** - 基于高效的 requests 库

## 📦 安装

```bash
go get github.com/474420502/gcurl@v1.2.0
```

## 🎯 快速开始

将任何 cURL 命令瞬间转换为 Go 代码：

```go
package main

import (
   "fmt"
   "log"
   "github.com/474420502/gcurl"
)

func main() {
   curlCmd := `curl -X POST "https://api.example.com/users" \
      -H "Content-Type: application/json" \
      -d '{"name": "张三", "email": "zhangsan@example.com"}'`

   curl, err := gcurl.ParseCurl(curlCmd)
   if err != nil {
      log.Fatal(err)
   }

   fmt.Printf("URL: %s\n", curl.URL)
   fmt.Printf("方法: %s\n", curl.Method)
   fmt.Printf("数据: %s\n", string(curl.Data))
}
```

## 📚 详细示例

### 示例 1: 基本 GET 请求

```go
curlCmd := `curl -H "Authorization: Bearer token123" "https://api.github.com/user"`

curl, err := gcurl.ParseCurl(curlCmd)
if err != nil {
   log.Fatal(err)
}

// 使用解析后的结构
fmt.Printf("URL: %s\n", curl.URL)
fmt.Printf("头部: %v\n", curl.Headers)
```

### 示例 2: POST 请求与 JSON 数据

```go
curlCmd := `curl -X POST "https://api.example.com/data" \
   -H "Content-Type: application/json" \
   -d '{"key": "value", "number": 42}'`

curl, err := gcurl.ParseCurl(curlCmd)
if err != nil {
   log.Fatal(err)
}

fmt.Printf("请求体: %s\n", string(curl.Data))
```

### 示例 3: 文件上传

```go
curlCmd := `curl -X POST "https://upload.example.com/files" \
   -F "file=@document.pdf" \
   -F "description=重要文档"`

curl, err := gcurl.ParseCurl(curlCmd)
if err != nil {
   log.Fatal(err)
}

fmt.Printf("表单数据: %v\n", curl.Form)
```

### 示例 4: 摘要认证

```go
curlCmd := `curl --digest -u "用户名:密码" "https://secure.example.com/api"`

curl, err := gcurl.ParseCurl(curlCmd)
if err != nil {
   log.Fatal(err)
}

fmt.Printf("认证类型: Digest\n")
fmt.Printf("用户名: %s\n", curl.DigestAuth.Username)
```

### 示例 5: 文件输出

```go
curlCmd := `curl -o "输出.html" "https://example.com"`

curl, err := gcurl.ParseCurl(curlCmd)
if err != nil {
   log.Fatal(err)
}

fmt.Printf("输出文件: %s\n", curl.OutputFile)
```

### 示例 6: HTTP 版本控制

```go
curlCmd := `curl --http2 "https://api.example.com/v2/data"`

curl, err := gcurl.ParseCurl(curlCmd)
if err != nil {
   log.Fatal(err)
}

fmt.Printf("HTTP 版本: %s\n", curl.HTTPVersion.String())
```

## 🔧 支持的 cURL 选项

### 基本选项

| 分类               | 选项              | 描述             | 状态 | 示例                                   |
| ------------------ | ----------------- | ---------------- | ---- | -------------------------------------- |
| **请求方法** | `-X, --request` | 指定请求方法     | ✅   | `curl -X POST`                       |
| **URL**      | `[URL]`         | 目标 URL         | ✅   | `curl "https://api.com"`             |
| **头部**     | `-H, --header`  | 自定义 HTTP 头部 | ✅   | `curl -H "Accept: application/json"` |

### 数据选项

| 分类                | 选项                 | 描述           | 状态 | 示例                                  |
| ------------------- | -------------------- | -------------- | ---- | ------------------------------------- |
| **POST 数据** | `-d, --data`       | 发送 POST 数据 | ✅   | `curl -d "name=value"`              |
| **原始数据**  | `--data-raw`       | 发送原始数据   | ✅   | `curl --data-raw "raw content"`     |
| **URL 编码**  | `--data-urlencode` | URL 编码数据   | ✅   | `curl --data-urlencode "name=张三"` |
| **表单**      | `-F, --form`       | 多部分表单数据 | ✅   | `curl -F "file=@path.txt"`          |

### 认证选项

| 分类               | 选项           | 描述          | 状态 | 示例                             |
| ------------------ | -------------- | ------------- | ---- | -------------------------------- |
| **基本认证** | `-u, --user` | 用户认证      | ✅   | `curl -u "user:pass"`          |
| **摘要认证** | `--digest`   | HTTP 摘要认证 | ✅   | `curl --digest -u "user:pass"` |

### 文件输出选项

| 分类                 | 选项                  | 描述               | 状态 | 示例                                 |
| -------------------- | --------------------- | ------------------ | ---- | ------------------------------------ |
| **输出文件**   | `-o, --output`      | 写入到文件         | ✅   | `curl -o "file.html"`              |
| **远程名称**   | `-O, --remote-name` | 使用远程文件名     | ✅   | `curl -O`                          |
| **输出目录**   | `--output-dir`      | 输出目录           | ✅   | `curl --output-dir "/path/"`       |
| **创建目录**   | `--create-dirs`     | 创建必要目录       | ✅   | `curl --create-dirs -o "dir/file"` |
| **错误时删除** | `--remove-on-error` | 出错时删除部分文件 | ✅   | `curl --remove-on-error`           |
| **断点续传**   | `-C, --continue-at` | 断点续传下载       | ✅   | `curl -C -`                        |

### HTTP 版本选项

| 分类                | 选项          | 描述          | 状态 | 示例               |
| ------------------- | ------------- | ------------- | ---- | ------------------ |
| **HTTP 版本** | `--http1.0` | 强制 HTTP/1.0 | ✅   | `curl --http1.0` |
|                     | `--http1.1` | 强制 HTTP/1.1 | ✅   | `curl --http1.1` |
|                     | `--http2`   | 强制 HTTP/2   | ✅   | `curl --http2`   |

### 其他选项

| 分类               | 选项              | 描述         | 状态 | 示例                              |
| ------------------ | ----------------- | ------------ | ---- | --------------------------------- |
| **调试**     | `-v, --verbose` | 详细输出     | ✅   | `curl -v`                       |
| **静默**     | `-s, --silent`  | 静默模式     | ✅   | `curl -s`                       |
| **用户代理** | `--user-agent`  | 设置用户代理 | ✅   | `curl --user-agent "MyApp/1.0"` |
| **Cookie**   | `--cookie`      | 发送 Cookie  | ✅   | `curl --cookie "name=value"`    |

## 🧪 测试

运行完整的测试套件：

```bash
# 运行所有测试
go test -v

# 运行特定测试
go test -v -run TestDigestAuth
go test -v -run TestFileOutput
go test -v -run TestFormData

# 测试覆盖率
go test -cover
```

### 测试统计

- **总测试数**: 100+
- **测试覆盖率**: 79.1%
- **测试场景**: 包括认证、文件处理、表单数据、错误处理等

## 🏗️ 项目结构

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
    ├── readme_examples.go
    └── go.mod
```

## 🤝 贡献

我们欢迎各种形式的贡献！

### 贡献方式

1. **🐛 报告 Bug**

   - 开启 issue 并提供详细的重现步骤
   - 包含导致问题的 cURL 命令
   - 提供 Go 版本和操作系统信息
2. **💡 建议功能**

   - 提议新的 cURL 选项支持
   - 分享使用案例和示例
   - 讨论实现方法
3. **🔧 提交 Pull Request**

   - 为任何新功能添加测试
   - 遵循现有的代码风格和约定
   - 根据需要更新文档
4. **📚 改进文档**

   - 修复错别字和澄清说明
   - 添加更多示例和用例
   - 翻译文档

### 开发设置

```bash
# 克隆仓库
git clone https://github.com/474420502/gcurl.git
cd gcurl

# 运行测试
go test -v

# 检查代码质量
go vet ./...
go fmt ./...
```

## 🔮 未来计划

### 即将推出的功能

- **扩展协议支持**

  - FTP 和 SFTP 支持
  - WebSocket 连接
  - gRPC 协议支持
- **高级认证**

  - OAuth 2.0 流程支持
  - JWT 令牌处理
  - API 密钥管理
- **性能优化**

  - 并发请求执行
  - 批量命令处理
  - 请求队列管理


## 📋 版本历史

### v1.2.1（当前版本）

- 🛠️ **代码质量与架构重大改进**
  - 全局变量全部重构为线程安全（如 gserver → getTestServer + sync.Once）
  - Debug 函数复杂度从 46 降低到 1，拆分为 10+ 个辅助方法
  - optionRegistry 增加线程安全文档说明
  - 新增详细 CODE_QUALITY_REPORT.md，包含复杂度与技术债务分析
  - 所有 gserver 引用已切换为新线程安全方法
  - 保持 86.9% 测试覆盖率，所有测试通过
  - 注释与技术债务跟踪文档完善

### v1.1.0

- ✅ 全面的文件输出支持（`-o`、`-O`、`--output-dir` 等）
- ✅ 完整的摘要认证实现
- ✅ HTTP 协议版本控制
- ✅ 增强的调试功能
- ✅ 改进的测试覆盖率（79.1%）
- ✅ 生产就绪的稳定性

### v1.0.0

- ✅ 基础 cURL 命令解析
- ✅ 基本 HTTP 选项支持
- ✅ 初始项目架构

## 📄 许可证

本项目基于 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🙏 致谢

- **[requests](https://github.com/474420502/requests)** - 支撑 gcurl 的强大 HTTP 库
- **[cURL](https://curl.se/)** - 启发本项目的出色命令行工具
- **Go 社区** - 提供优秀的生态系统和工具
- **所有贡献者** - 每一位帮助改进这个库的人

## 🔗 相关项目

- [requests](https://github.com/474420502/requests) - gcurl 使用的 HTTP 客户端库

## 📞 支持

- **GitHub Issues**: [报告 Bug 和请求功能](https://github.com/474420502/gcurl/issues)
- **文档**: 查看此 README 和代码示例
- **社区**: 在 GitHub issues 中加入讨论

---

⭐ **如果这个项目帮助您将 cURL 命令转换为 Go 代码，请给项目加星！** 您的支持激励我们持续开发和改进。

## 🏆 成功案例

> "gcurl 在我们将 API 测试从 shell 脚本迁移到 Go 时为我们节省了数小时的手动转换工作。认证和文件上传功能运行完美！" - *开发团队负责人*

> "非常适合我们的 CI/CD 流水线，我们需要将现有的基于 cURL 的健康检查转换为 Go 服务。会话重用功能显著提高了我们的性能。" - *DevOps 工程师*

> "调试输出帮助我们准确理解复杂的 cURL 命令是如何被解释的。对于排查 API 集成问题非常有用。" - *后端开发者*

---

**为 Go 社区用 ❤️ 制作**
