package gcurl

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/474420502/requests"
)

// BodyData 表示不同类型的请求体数据
type BodyData struct {
	Type    string      // "raw", "form", "json", "urlencoded", "multipart"
	Content interface{} // 具体内容，根据Type而不同
}

// 向后兼容方法
func (bd *BodyData) Len() int {
	if bd == nil {
		return 0
	}
	switch bd.Type {
	case "raw":
		if buf, ok := bd.Content.(*bytes.Buffer); ok {
			return buf.Len()
		}
	case "form", "urlencoded":
		if str, ok := bd.Content.(string); ok {
			return len(str)
		}
	case "multipart":
		// multipart的长度需要在构建时计算，这里返回字段数量
		if fields, ok := bd.Content.([]*FormField); ok {
			return len(fields)
		}
	}
	return 0
}

func (bd *BodyData) Read(p []byte) (n int, err error) {
	if bd == nil {
		return 0, io.EOF
	}
	switch bd.Type {
	case "raw":
		if buf, ok := bd.Content.(*bytes.Buffer); ok {
			return buf.Read(p)
		}
	}
	return 0, io.EOF
}

func (bd *BodyData) String() string {
	if bd == nil {
		return ""
	}
	switch bd.Type {
	case "raw":
		if buf, ok := bd.Content.(*bytes.Buffer); ok {
			return buf.String()
		}
	case "form", "urlencoded":
		if str, ok := bd.Content.(string); ok {
			return str
		}
	case "multipart":
		if fields, ok := bd.Content.([]*FormField); ok {
			var parts []string
			for _, field := range fields {
				if field.IsFile {
					parts = append(parts, fmt.Sprintf("%s=@%s", field.Name, field.Value))
				} else {
					parts = append(parts, fmt.Sprintf("%s=%s", field.Name, field.Value))
				}
			}
			return strings.Join(parts, "&")
		}
	}
	return ""
}

func (bd *BodyData) WriteString(s string) (n int, err error) {
	if bd == nil {
		return 0, fmt.Errorf("body data is nil")
	}
	switch bd.Type {
	case "raw":
		if buf, ok := bd.Content.(*bytes.Buffer); ok {
			return buf.WriteString(s)
		}
	}
	return 0, fmt.Errorf("cannot write to body type: %s", bd.Type)
}

// setRawBody 设置原始字节数据类型的Body
func (c *CURL) setRawBody(data []byte) {
	c.Body = &BodyData{
		Type:    "raw",
		Content: bytes.NewBuffer(data),
	}
}

// setRawBodyString 设置原始字符串类型的Body
func (c *CURL) setRawBodyString(data string) {
	c.Body = &BodyData{
		Type:    "raw",
		Content: bytes.NewBufferString(data),
	}
}

// CURL 信息结构
type CURL struct {
	ParsedURL *url.URL
	Method    string
	Header    http.Header
	CookieJar http.CookieJar
	// --- 请做如下修改 ---
	// Cookies   []http.Cookie // 旧定义
	Cookies []*http.Cookie // 新定义：使用指针切片，符合标准库和requests库的用法
	// --- 修改结束 ---
	ContentType    string
	Body           *BodyData // 改为更灵活的结构
	Auth           *requests.BasicAuth
	Timeout        int // 对应 --max-time, 总超时
	ConnectTimeout int // 新增字段，对应 --connect-timeout
	Insecure       bool
	Proxy          string // 新增字段，用于存储代理地址
	LimitRate      string // 新增字段，用于存储传输速度限制

	// 新增SSL/TLS相关字段
	CACert     string // --cacert 自定义CA证书路径
	ClientCert string // --cert 客户端证书路径
	ClientKey  string // --key 客户端私钥路径

	// 新增HTTP协议相关字段
	HTTP2          bool // --http2 强制使用HTTP/2
	MaxRedirs      int  // --max-redirs 最大重定向次数 (-1表示无限制)
	FollowRedirect bool // -L/--location 是否跟随重定向
}

// New new 一个 curl 出来
func New() *CURL {
	u := &CURL{}
	u.Insecure = false
	u.Header = make(http.Header)
	u.CookieJar, _ = cookiejar.New(nil)
	u.Body = &BodyData{Type: "raw", Content: bytes.NewBuffer(nil)}
	u.Timeout = 30           // 默认总超时
	u.ConnectTimeout = 0     // 0 表示不设置，使用系统默认
	u.LimitRate = ""         // 默认不限速
	u.MaxRedirs = -1         // 默认无限制重定向
	u.HTTP2 = false          // 默认不强制HTTP/2
	u.FollowRedirect = false // 默认不跟随重定向（与curl默认行为一致）
	u.CACert = ""            // 默认无自定义CA证书
	u.ClientCert = ""        // 默认无客户端证书
	u.ClientKey = ""         // 默认无客户端私钥
	// --- 为了匹配新的字段类型，初始化也做相应调整 ---
	u.Cookies = make([]*http.Cookie, 0)
	return u
}

func (curl *CURL) String() string {
	if curl != nil {
		return fmt.Sprintf("Method: %s\nParsedURL: %s\nHeader: %s\nCookie: %s",
			curl.Method, curl.ParsedURL.String(), curl.Header, curl.Cookies)
	}
	return ""
}

// Execute 直接执行curlbash
func Execute(curlbash string) (*requests.Response, error) {
	c, err := ParseBash(curlbash)
	if err != nil {
		return nil, err
	}
	return c.CreateRequest(nil).Execute()
}

// CreateSession 创建Session
func (curl *CURL) CreateSession() *requests.Session {
	ses := requests.NewSession()

	// 设置基本配置
	ses.SetHeader(curl.Header)
	ses.SetCookies(curl.ParsedURL, curl.Cookies)

	// 设置总超时
	if curl.Timeout > 0 {
		ses.Config().SetTimeout(time.Duration(curl.Timeout) * time.Second)
	}

	// 设置认证
	if curl.Auth != nil {
		ses.Config().SetBasicAuth(curl.Auth.User, curl.Auth.Password)
	}

	// 设置跳过TLS验证
	if curl.Insecure {
		ses.Config().SetInsecure(curl.Insecure)
	}

	// 设置代理（包括SOCKS5）
	if curl.Proxy != "" {
		ses.Config().SetProxy(curl.Proxy)
	}

	// 设置连接超时（如果指定了）
	if curl.ConnectTimeout > 0 {
		// 目前先使用总超时作为连接超时的替代方案
		// 理想情况下应该在requests库中添加专门的SetConnectTimeout方法
		// TODO: 在requests库中实现真正的连接超时设置
		ses.Config().SetTimeout(time.Duration(curl.ConnectTimeout) * time.Second)
	}

	// 设置重定向策略
	if curl.FollowRedirect {
		// 设置最大重定向次数
		maxRedirs := curl.MaxRedirs
		if maxRedirs < 0 {
			maxRedirs = 30 // 默认值
		}
		// TODO: 在requests库中添加SetRedirectPolicy方法
		// ses.Config().SetRedirectPolicy(maxRedirs)
		// 目前先用注释标记这个功能需要实现
	}

	// 设置TLS/SSL证书配置
	if curl.CACert != "" {
		// TODO: 在requests库中添加SetCACert方法
		// ses.Config().SetCACert(curl.CACert)
	}

	if curl.ClientCert != "" && curl.ClientKey != "" {
		// TODO: 在requests库中添加SetClientCerts方法
		// ses.Config().SetClientCerts(curl.ClientCert, curl.ClientKey)
	}

	return ses
}

// CreateRequest 根据Session 创建Request
func (curl *CURL) CreateRequest(ses *requests.Session) *requests.Request {
	var wf *requests.Request

	if ses == nil {
		ses = curl.CreateSession()
	}

	curl.Method = strings.ToUpper(curl.Method)

	switch curl.Method {
	case "HEAD":
		wf = ses.Head(curl.ParsedURL.String())
	case "GET", "":
		wf = ses.Get(curl.ParsedURL.String())
	case "POST":
		wf = ses.Post(curl.ParsedURL.String())
	case "PUT":
		wf = ses.Put(curl.ParsedURL.String())
	case "PATCH":
		wf = ses.Patch(curl.ParsedURL.String())
	case "OPTIONS":
		wf = ses.Options(curl.ParsedURL.String())
	case "DELETE":
		wf = ses.Delete(curl.ParsedURL.String())
	default:
		panic("curl.Method is not UNKNOWN")
	}

	wf.SetContentType(curl.ContentType)

	// 根据Body类型设置不同的请求体
	if curl.Body != nil {
		switch curl.Body.Type {
		case "raw":
			if buf, ok := curl.Body.Content.(*bytes.Buffer); ok {
				wf.SetBody(buf)
			}
		case "multipart":
			if fields, ok := curl.Body.Content.([]*FormField); ok {
				// 将FormField转换为requests库支持的格式
				formData := make(map[string]interface{})
				for _, field := range fields {
					if field.IsFile {
						// 文件上传
						formData[field.Name] = field.Value // requests库会处理文件路径
					} else {
						// 普通字段
						formData[field.Name] = field.Value
					}
				}
				wf.SetBodyFormData(formData)
			}
		case "form", "urlencoded":
			if str, ok := curl.Body.Content.(string); ok {
				wf.SetBody(strings.NewReader(str))
			}
		}
	}

	return wf
}

// Request 根据自己CreateSession 创建Request
func (curl *CURL) Request() *requests.Request {
	return curl.CreateRequest(curl.CreateSession())
}

// Temporary 向后兼容方法，内部调用Request()
// Deprecated: 使用 Request() 方法替代
func (curl *CURL) Temporary() *requests.Request {
	return curl.Request()
}

// CreateTemporary 向后兼容方法，内部调用CreateRequest()
// Deprecated: 使用 CreateRequest() 方法替代
func (curl *CURL) CreateTemporary(ses *requests.Session) *requests.Request {
	return curl.CreateRequest(ses)
}

type MatchGroup int

const (
	HTTPHTTPS MatchGroup = iota
	ShortNoArg
	LongNoArgSpecial
	DataBinary
	LongArgQuotes
	LongArgDoubleQuotes
	LongArgNoQuotes
	ShortArgQuotes
	ShortArgDoubleQuotes
	ShortArgNoQuotes
	NewlineQuotes
	NewlineDoubleQuotes
	LongArgNoArg
)

// cmdformat2bash cmdformat2bash
func cmdformat2bash(scurl string) string {
	builder := &strings.Builder{}
	i := 0
	for i < len(scurl) {
		c := scurl[i]
		if c == '^' {
			if i+3 < len(scurl) && scurl[i+1] == '\\' && scurl[i+2] == '^' {
				// 处理 ^\\^"
				// log.Println(scurl[i:i+4], string(scurl[i+3]))
				builder.WriteByte(scurl[i+3])
				i += 4
			} else if i+2 < len(scurl) && scurl[i+2] == '^' {
				// ^%^ 处理这种字符串转换
				// log.Println(scurl[i:i+3], string(scurl[i+1]))
				builder.WriteByte(scurl[i+1])
				i += 3
			} else if i+1 < len(scurl) && scurl[i+1] == '"' {
				// log.Println(scurl[i:i+2], string(scurl[i+1]))
				builder.WriteByte('\'')
				i += 2
			} else if i+1 < len(scurl) {
				// 处理 ^\\n 处理通用的转意格式
				// 处理 ^" 特殊的把符号转换为regexp能识别的格式
				// log.Println(scurl[i:i+2], string(scurl[i+1]))
				builder.WriteByte(scurl[i+1])
				i += 2
			} else {
				builder.WriteByte(c)
				i++
			}
		} else {
			builder.WriteByte(c)
			i++
		}
	}
	return builder.String()
}

// ParseCmd curl cmd  *Supports copying as cURL command (Cmd)
func ParseCmd(scurl string) (curl *CURL, err error) {
	return ParseBash(cmdformat2bash(scurl))
}

// (-H \\^\"|\\^\n|\\^\\\\\\^|\\^%\\^)
var recheckCmdFormat = regexp.MustCompile("(-H \\^\"|\\^\n|\\^\\\\\\^|\\^%\\^)")

// CheckCmdForamt CheckCmdFormat checks if a curl string is in the cmd format.
func CheckCmdForamt(scurl string) bool {
	// x := recheckCmdFormat.FindAllString(scurl, -1)
	// log.Println(x)
	return recheckCmdFormat.MatchString(scurl)
}

// Parse This method is compatible with both cmd and bash formats
// but it merely forcibly converts cmd to bash.
// It's recommended to use ParseBash instead.
// If you encounter any issues, please submit an issue so that I can fix it.
func Parse(scurl string) (curl *CURL, err error) {
	if CheckCmdForamt(scurl) {
		return ParseCmd(scurl)
	}
	return ParseBash(scurl)
}

func ParseBash(scurl string) (*CURL, error) {
	// 1. 使用新的纯Go分词器
	lexer := NewLexer(scurl)
	if err := lexer.Parse(); err != nil {
		return nil, fmt.Errorf("failed to tokenize curl command: %w", err)
	}
	args := lexer.Tokens

	// 2. 将分词结果传递给选项处理器
	// (此部分将在下一节中重构)
	return buildFromArgs(args)
}
