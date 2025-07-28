package gcurl

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"

	"github.com/474420502/requests"
)

// CURL 信息结构
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
	Body           *bytes.Buffer
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
	HTTP2     bool // --http2 强制使用HTTP/2
	MaxRedirs int  // --max-redirs 最大重定向次数 (-1表示无限制)
}

// New new 一个 curl 出来
func New() *CURL {
	u := &CURL{}
	u.Insecure = false
	u.Header = make(http.Header)
	u.CookieJar, _ = cookiejar.New(nil)
	u.Body = bytes.NewBuffer(nil)
	u.Timeout = 30       // 默认总超时
	u.ConnectTimeout = 0 // 0 表示不设置，使用系统默认
	u.LimitRate = ""     // 默认不限速
	u.MaxRedirs = -1     // 默认无限制重定向
	u.HTTP2 = false      // 默认不强制HTTP/2
	u.CACert = ""        // 默认无自定义CA证书
	u.ClientCert = ""    // 默认无客户端证书
	u.ClientKey = ""     // 默认无客户端私钥
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
	return c.CreateTemporary(nil).Execute()
}

// CreateSession 创建Session
func (curl *CURL) CreateSession() *requests.Session {
	ses := requests.NewSession()

	// 设置基本配置
	ses.SetHeader(curl.Header)
	ses.SetCookies(curl.ParsedURL, curl.Cookies)

	// 设置总超时
	ses.Config().SetTimeout(curl.Timeout)

	// 设置认证
	if curl.Auth != nil {
		ses.Config().SetBasicAuth(curl.Auth)
	}

	// 设置跳过TLS验证
	if curl.Insecure {
		ses.Config().SetInsecure(curl.Insecure)
	}

	// 设置代理（包括SOCKS5）
	if curl.Proxy != "" {
		ses.Config().SetProxy(curl.Proxy)
	}

	// 注意：ConnectTimeout 的设置需要在requests库中添加支持
	// 目前我们只是解析和存储这个值，实际的连接超时设置
	// 需要requests库本身提供相应的接口
	// TODO: 如果requests库支持连接超时配置，在这里调用相应方法

	return ses
}

// CreateTemporary 根据Session 创建Temporary
func (curl *CURL) CreateTemporary(ses *requests.Session) *requests.Temporary {
	var wf *requests.Temporary

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
	wf.SetBody(curl.Body)
	return wf
}

// Temporary 根据自己CreateSession 创建Temporary
func (curl *CURL) Temporary() *requests.Temporary {
	return curl.CreateTemporary(curl.CreateSession())
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
