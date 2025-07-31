# gcurl 改进计划：从解析器到完整 cURL 替代方案

## 当前状态分析

### ✅ 已实现的优势
- 完整的 cURL 命令解析
- 基本 HTTP 方法支持 (GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS)
- 头部处理完善
- Cookie 基础支持
- 文件上传 (multipart)
- 基本认证
- SSL 配置 (--insecure, --cacert, --cert, --key)
- 代理支持
- 超时控制
- 重定向处理
- 数据编码 (--data-urlencode, --data-raw, --data-binary)

### 🔍 待改进的关键问题

#### 1. 类型安全性 (Type Safety)
- **BodyData**: 当前使用 `interface{}` 缺乏类型安全
- **Cookies**: 已使用正确的 `[]*http.Cookie` 类型

#### 2. 用户体验 (Developer Experience) 
- **缺少调试输出**: 没有 `-v/--verbose` 支持
- **缺少响应头显示**: 没有 `-i/--include` 支持  
- **文档不完整**: 缺少清晰的 API 文档

#### 3. 功能完整性 (Feature Completeness)
- **认证**: 缺少 `--digest` 摘要认证
- **协议控制**: 缺少 `--http1.1/--http1.0` 强制版本
- **网络调试**: 缺少 `--trace` 详细跟踪
- **输出控制**: 缺少 `-o/--output` 文件输出

## 阶段性升级计划

### 🚀 阶段一：核心改进与类型安全 (立即开始)

#### 1.1 重构 Body 系统
```go
// 当前问题：interface{} 不够类型安全
type BodyData struct {
    Type    string
    Content interface{} // 这里缺乏类型安全
}

// 改进方案：定义明确的 Body 接口
type Body interface {
    ContentType() string
    WriteTo(w io.Writer) (int64, error)
    Length() int64
}

type RawBody struct {
    Data        []byte
    contentType string
}

type FormBody struct {
    Values url.Values
}

type MultipartBody struct {
    Fields []*FormField
    boundary string
}

type JSONBody struct {
    Data interface{}
}
```

#### 1.2 添加调试支持 (最高优先级)
```go
// 在 CURL 结构体中添加
type CURL struct {
    // ... 现有字段
    Verbose bool     // -v/--verbose
    Include bool     // -i/--include  
    Trace   bool     // --trace
    Silent  bool     // -s/--silent
}

// 添加调试方法
func (c *CURL) Debug() string
func (c *CURL) Verbose() string
func (c *CURL) Summary() string
```

#### 1.3 完善超时系统
```go
type CURL struct {
    // ... 现有字段
    Timeout           time.Duration // 总超时
    ConnectTimeout    time.Duration // 连接超时  
    DNSTimeout        time.Duration // DNS解析超时
    TLSHandshakeTimeout time.Duration // TLS握手超时
}
```

### 📈 阶段二：功能扩展与 cURL 对齐 (中期)

#### 2.1 认证系统扩展
```go
type AuthType int
const (
    AuthBasic AuthType = iota
    AuthDigest
    AuthBearer
    AuthNTLM
)

type Authentication struct {
    Type     AuthType
    Username string
    Password string
    Token    string
    // 摘要认证的特殊字段
    Realm    string
    Nonce    string
}
```

#### 2.2 协议控制
```go
type CURL struct {
    // ... 现有字段
    HTTPVersion string // "1.0", "1.1", "2", "auto"
    ForceIPv4   bool   // -4/--ipv4
    ForceIPv6   bool   // -6/--ipv6
    Resolve     map[string]string // --resolve host:port:addr
}
```

#### 2.3 响应处理增强
```go
type Response struct {
    *requests.Response
    Headers    http.Header
    StatusLine string
    Verbose    []string // 详细日志
}

// 添加响应处理方法
func (r *Response) IncludeHeaders() string
func (r *Response) SaveToFile(filename string) error
func (r *Response) TraceInfo() []string
```

### 🔧 阶段三：深度集成与高级功能 (长期)

#### 3.1 与 requests 深度集成
- 暴露中间件系统
- 支持连接池配置
- 支持自定义传输层

#### 3.2 高级网络功能
- `--interface` 指定网络接口
- `--dns-servers` 自定义DNS服务器
- `--happy-eyeballs-timeout` IPv6优先级控制

#### 3.3 性能与监控
- 请求时间统计
- 连接复用统计
- 内存使用优化

## 实施优先级

### 🔥 高优先级 (立即实施)
1. **调试输出** (`-v/--verbose`)
2. **响应头显示** (`-i/--include`)
3. **Body 类型安全重构**
4. **完善 API 文档**

### 🔥 中优先级 (近期实施)  
1. **摘要认证** (`--digest`)
2. **协议版本控制** (`--http1.1/--http1.0`)
3. **文件输出** (`-o/--output`)
4. **详细跟踪** (`--trace`)

### 🔥 低优先级 (长期规划)
1. **网络接口控制**
2. **性能监控**
3. **高级代理功能**

## 代码示例：改进后的用法

```go
// 基础用法 (向后兼容)
curl, _ := gcurl.Parse(`curl -v "https://httpbin.org/get"`)
resp, _ := curl.Request().Execute()

// 高级调试用法
curl.SetVerbose(true)
fmt.Println(curl.Debug()) // 显示解析的详细信息

// 响应处理
resp.IncludeHeaders() // 包含响应头
resp.SaveToFile("response.json") // 保存到文件

// 类型安全的 Body 构建
body := NewJSONBody(map[string]interface{}{
    "name": "test",
    "age":  25,
})
curl.SetBody(body)
```

## 成功指标

### 用户体验指标
- [ ] 100% cURL 命令解析成功率
- [ ] 详细错误信息和调试输出
- [ ] 零学习成本的 API 设计

### 功能覆盖指标  
- [ ] 覆盖 80% 的常用 cURL 选项
- [ ] 支持所有主要认证方式
- [ ] 完整的协议版本控制

### 性能指标
- [ ] 解析性能 < 1ms per command
- [ ] 内存使用 < 1MB for typical use
- [ ] 零内存泄漏

## 下一步行动

1. **立即开始**: 实施 `--verbose` 支持
2. **本周完成**: Body 类型安全重构  
3. **本月目标**: 完成阶段一的所有改进
4. **持续迭代**: 根据用户反馈调整优先级

这个计划将把 gcurl 从一个优秀的 cURL 解析器升级为一个功能完整、类型安全、用户友好的 cURL 替代方案。
