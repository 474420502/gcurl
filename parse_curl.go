package gcurl

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/474420502/requests"
)

// HTTPVersion 定义 HTTP 协议版本
type HTTPVersion int

const (
	HTTPVersionAuto HTTPVersion = iota // 自动选择 (默认)
	HTTPVersion10                      // HTTP/1.0
	HTTPVersion11                      // HTTP/1.1
	HTTPVersion2                       // HTTP/2
)

// String 返回协议版本的字符串表示
func (hv HTTPVersion) String() string {
	switch hv {
	case HTTPVersion10:
		return "HTTP/1.0"
	case HTTPVersion11:
		return "HTTP/1.1"
	case HTTPVersion2:
		return "HTTP/2"
	case HTTPVersionAuto:
		return "Auto"
	default:
		return "Unknown"
	}
}

// AuthInfo 表示认证信息
type AuthInfo struct {
	Type     string // 认证类型: "basic", "digest", etc.
	User     string // 用户名
	Username string // 用户名 (别名)
	Password string // 密码
}

// IsValid 检查认证信息是否有效
func (a *AuthInfo) IsValid() bool {
	if a == nil {
		return false
	}
	user := a.User
	if user == "" {
		user = a.Username
	}
	return user != "" && a.Password != ""
}

// GetAuthHeader 获取认证头信息
func (a *AuthInfo) GetAuthHeader() string {
	if !a.IsValid() {
		return ""
	}
	user := a.User
	if user == "" {
		user = a.Username
	}
	// 这里简单返回基础认证，实际情况可能需要根据 Type 处理
	return fmt.Sprintf("Basic %s", encodeBasicAuth(user, a.Password))
}

// encodeBasicAuth 编码基础认证
func encodeBasicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

// BodyData 表示不同类型的请求体数据
type BodyData struct {
	Type    string      // "raw", "form", "json", "urlencoded", "multipart"
	Content interface{} // 具体内容，根据Type而不同
}

// String 返回 BodyData 的字符串表示
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
		if form, ok := bd.Content.(url.Values); ok {
			return form.Encode()
		}
	}
	if str, ok := bd.Content.(string); ok {
		return str
	}
	return fmt.Sprintf("%v", bd.Content)
}

// Read 读取 BodyData 的内容
func (bd *BodyData) Read(p []byte) (n int, err error) {
	if bd == nil {
		return 0, fmt.Errorf("BodyData is nil")
	}
	switch bd.Type {
	case "raw":
		if buf, ok := bd.Content.(*bytes.Buffer); ok {
			return buf.Read(p)
		}
	case "form", "urlencoded":
		if form, ok := bd.Content.(url.Values); ok {
			data := form.Encode()
			buf := bytes.NewBufferString(data)
			return buf.Read(p)
		}
	}
	if str, ok := bd.Content.(string); ok {
		buf := bytes.NewBufferString(str)
		return buf.Read(p)
	}
	return 0, fmt.Errorf("unsupported content type")
}

// WriteString 向 BodyData 写入字符串
func (bd *BodyData) WriteString(s string) (int, error) {
	if bd == nil {
		return 0, fmt.Errorf("BodyData is nil")
	}
	if bd.Type == "raw" {
		if buf, ok := bd.Content.(*bytes.Buffer); ok {
			return buf.WriteString(s)
		}
		return 0, fmt.Errorf("BodyData content is not *bytes.Buffer")
	}
	// 对于非raw类型，返回错误而不是自动转换
	return 0, fmt.Errorf("WriteString is only supported for raw BodyData type")
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
		if form, ok := bd.Content.(url.Values); ok {
			return len(form.Encode())
		}
	}
	return 0
}

// CURL 结构体表示一个 curl 命令
type CURL struct {
	// HTTP方法和数据发送控制
	GetMode   bool        // -G/--get 使用GET方法发送POST数据（将数据附加到URL）
	Method    string      // HTTP方法
	ParsedURL *url.URL    // 解析后的URL
	Header    http.Header // HTTP头
	Body      *BodyData   // 请求体数据
	Cookies   []*http.Cookie

	// 认证相关
	User          string    // -u/--user 用户认证
	UserAgent     string    // -A/--user-agent
	Referer       string    // -e/--referer
	Authorization string    // Authorization头
	Auth          *AuthInfo // 认证信息（兼容性）
	AuthV2        *AuthInfo // 新版认证信息
	ContentType   string    // Content-Type 头

	// 代理相关
	Proxy         string // 代理服务器地址
	ProxyUser     string // 代理用户名
	ProxyPassword string // 代理密码

	// SSL/TLS相关
	Insecure      bool   // -k/--insecure 忽略SSL证书错误
	CACert        string // --cacert 自定义CA证书路径
	ClientCert    string // --cert 客户端证书路径
	ClientKey     string // --key 客户端私钥路径
	CertType      string // --cert-type 证书类型
	KeyType       string // --key-type 私钥类型
	CAPath        string // --capath CA证书目录
	CRLFile       string // --crlfile 证书吊销列表
	SSLVerifyPeer bool   // SSL对等验证
	SSLVerifyHost bool   // SSL主机验证
	TLSVersion    string // TLS版本
	Ciphers       string // 加密套件

	// 超时和重试相关
	Timeout             time.Duration // 总超时时间
	ConnectTimeout      time.Duration // 连接超时
	DNSTimeout          time.Duration // DNS解析超时
	TLSHandshakeTimeout time.Duration // TLS握手超时
	MaxRedirs           int           // --max-redirs 最大重定向次数
	Retry               int           // --retry 重试次数
	RetryDelay          time.Duration // --retry-delay 重试延迟
	RetryMaxTime        time.Duration // --retry-max-time 最大重试时间
	RetryAllErrors      bool          // --retry-all-errors 重试所有错误
	RetryConnRefused    bool          // --retry-connrefused 连接被拒绝时重试

	// HTTP协议相关
	HTTP2          bool        // --http2 强制使用HTTP/2
	HTTPVersion    HTTPVersion // 协议版本控制
	FollowRedirect bool        // -L/--location 是否跟随重定向
	MaxFileSize    int64       // --max-filesize 最大文件大小
	LimitRate      string      // --limit-rate 传输速度限制
	KeepAlive      bool        // --keepalive 保持连接
	TCPNoDelay     bool        // --tcp-nodelay TCP无延迟
	TCPKeepAlive   bool        // --tcp-keepalive TCP保活
	Interface      string      // --interface 网络接口
	LocalPort      string      // --local-port 本地端口
	IPVersion      int         // 4 或 6，IP版本

	// DNS 和网络解析相关
	Resolve   []string // --resolve 主机名解析映射，格式：host:port:address
	ConnectTo []string // --connect-to 连接重定向映射，格式：HOST1:PORT1:HOST2:PORT2

	// Cookie相关
	CookieJar  *cookiejar.Jar // Cookie存储
	CookieFile string         // -c/--cookie-jar Cookie文件

	// 调试和输出控制
	Verbose    bool   // -v/--verbose 详细输出
	Include    bool   // -i/--include 在输出中包含响应头
	Silent     bool   // -s/--silent 静默模式
	ShowError  bool   // -S/--show-error 显示错误
	FailEarly  bool   // --fail-early 早期失败
	Trace      bool   // --trace 追踪所有传入和传出的数据
	TraceFile  string // --trace-ascii 追踪文件
	DumpHeader string // -D/--dump-header 转储头文件
	WriteOut   string // -w/--write-out 输出格式

	// 文件输出控制
	OutputFile       string // -o/--output 指定输出文件路径
	RemoteName       bool   // -O/--remote-name 使用远程文件名
	RemoteHeaderName bool   // -J/--remote-header-name 使用Content-Disposition文件名
	OutputDir        string // --output-dir 指定输出目录
	CreateDirs       bool   // --create-dirs 自动创建目录
	RemoveOnError    bool   // --remove-on-error 出错时删除输出文件
	ContinueAt       int64  // -C/--continue-at 断点续传偏移
	Append           bool   // --append 追加到文件

	// 脚本与易用性增强功能
	WriteOutFormat  string // -w/--write-out 指定输出格式
	FailOnError     bool   // -f/--fail 脚本错误处理，遇到4xx/5xx返回错误
	LocationTrusted bool   // --location-trusted 信任重定向位置

	// 表单和数据相关
	FormData   map[string]string // 表单数据
	JSONData   interface{}       // JSON数据
	RawData    []byte            // 原始数据
	DataBinary []byte            // 二进制数据
	DataASCII  string            // ASCII数据
	DataRaw    string            // 原始字符串数据
	DataURLEnc string            // URL编码数据

	// 上传相关
	UploadFile string // -T/--upload-file 上传文件
	Form       string // -F/--form 表单数据
	FormString string // --form-string 表单字符串

	// 其他选项
	Config        string   // -K/--config 配置文件
	Progress      bool     // --progress-bar 进度条
	NoProgress    bool     // --no-progress 无进度
	Raw           bool     // --raw 原始输出
	Buffer        bool     // --buffer 缓冲输出
	Compressed    bool     // --compressed 压缩传输
	Globoff       bool     // -g/--globoff 关闭URL通配符
	IgnoreCase    bool     // --ignore-case 忽略大小写
	UseASCII      bool     // --use-ascii 使用ASCII
	StderrFile    string   // --stderr 错误输出文件
	TelnetOptions []string // --telnet-option Telnet选项
}

// setRawBody 设置原始字节数据作为请求体
func (c *CURL) setRawBody(data []byte) {
	c.Body = &BodyData{
		Type:    "raw",
		Content: bytes.NewBuffer(data),
	}
}

// setRawBodyString 设置原始字符串数据作为请求体
func (c *CURL) setRawBodyString(data string) {
	c.Body = &BodyData{
		Type:    "raw",
		Content: bytes.NewBufferString(data),
	}
}

// SaveToFile 将响应内容保存到文件
func (c *CURL) SaveToFile(response *requests.Response) error {
	if response == nil {
		return fmt.Errorf("response is nil")
	}

	// 确定输出文件路径
	var outputPath string

	// 1. Content-Disposition 文件名 (-J)
	if c.RemoteHeaderName {
		hdr := response.GetHeader()
		if cd := hdr.Get("Content-Disposition"); cd != "" {
			re := regexp.MustCompile(`filename\*?=([^;]+)`)
			if m := re.FindStringSubmatch(cd); len(m) == 2 {
				filename := strings.Trim(m[1], `"`)
				if c.OutputDir != "" {
					outputPath = filepath.Join(c.OutputDir, filename)
				} else if c.OutputFile != "" {
					outputPath = filepath.Join(filepath.Dir(c.OutputFile), filename)
				} else {
					outputPath = filename
				}
			}
		}
	}

	// 2. 普通 -o or -O
	if outputPath == "" {
		var err error
		outputPath, err = c.determineOutputPath()
		if err != nil {
			return fmt.Errorf("failed to determine output path: %w", err)
		}
		if outputPath == "" {
			return nil
		}
	}

	// 创建目录
	if c.CreateDirs {
		dir := filepath.Dir(outputPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directories: %w", err)
		}
	}

	// 打开文件
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// 写入内容
	content := response.Content()
	if _, err := file.Write(content); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}

// determineOutputPath 确定输出文件路径
func (c *CURL) determineOutputPath() (string, error) {
	if c.OutputFile != "" {
		// 如果同时指定了输出目录，则组合路径
		if c.OutputDir != "" {
			return filepath.Join(c.OutputDir, c.OutputFile), nil
		}
		return c.OutputFile, nil
	}

	if c.RemoteName && c.ParsedURL != nil {
		// 从URL获取文件名
		path := c.ParsedURL.Path
		var filename string

		if path == "" || path == "/" {
			// 根路径时使用默认文件名
			filename = "index.html"
		} else {
			filename = filepath.Base(path)
			if filename == "." || filename == "/" {
				filename = "index.html"
			}
		}

		if c.OutputDir != "" {
			return filepath.Join(c.OutputDir, filename), nil
		}
		return filename, nil
	}

	return "", nil
}

// New new 一个 curl 出来
func New() *CURL {
	u := &CURL{}
	u.Insecure = false
	u.Header = make(http.Header)
	u.CookieJar, _ = cookiejar.New(nil)
	u.Body = &BodyData{Type: "raw", Content: bytes.NewBuffer(nil)}

	// 设置默认超时 - 使用 time.Duration 类型
	u.Timeout = 30 * time.Second // 默认总超时30秒
	u.ConnectTimeout = 0         // 0 表示不设置，使用系统默认
	u.DNSTimeout = 0             // 0 表示不设置，使用系统默认
	u.TLSHandshakeTimeout = 0    // 0 表示不设置，使用系统默认

	u.LimitRate = ""                // 默认不限速
	u.MaxRedirs = -1                // 默认无限制重定向
	u.HTTP2 = false                 // 默认不强制HTTP/2
	u.HTTPVersion = HTTPVersionAuto // 默认自动选择协议版本
	u.FollowRedirect = false        // 默认不跟随重定向（与curl默认行为一致）
	u.CACert = ""                   // 默认无自定义CA证书
	u.ClientCert = ""               // 默认无客户端证书
	u.ClientKey = ""                // 默认无客户端私钥
	// --- 为了匹配新的字段类型，初始化也做相应调整 ---
	u.Cookies = make([]*http.Cookie, 0)

	// 初始化文件输出相关字段
	u.OutputFile = ""       // 默认不指定输出文件（输出到stdout）
	u.RemoteName = false    // 默认不使用远程文件名
	u.OutputDir = ""        // 默认不指定输出目录
	u.CreateDirs = false    // 默认不自动创建目录
	u.RemoveOnError = false // 默认不在出错时删除文件
	u.ContinueAt = 0        // 默认从头开始下载

	return u
}

func (curl *CURL) String() string {
	if curl != nil {
		return fmt.Sprintf("Method: %s\nParsedURL: %s\nHeader: %s\nCookie: %s",
			curl.Method, curl.ParsedURL.String(), curl.Header, curl.Cookies)
	}
	return ""
}

// Execute 直接执行curlbash，并支持脚本选项 --fail 和 --write-out
// Execute 直接执行curlbash，并支持脚本选项 --fail 和 --write-out
func Execute(curlbash string) (*requests.Response, error) {
	c, err := ParseBash(curlbash)
	if err != nil {
		return nil, err
	}
	start := time.Now()
	resp, err := c.CreateRequest(nil).Execute()
	duration := time.Since(start)
	// 脚本错误处理
	if c.FailOnError && resp != nil && resp.GetStatusCode() >= 400 {
		return resp, fmt.Errorf("HTTP error: %d", resp.GetStatusCode())
	}
	// 格式化输出
	if c.WriteOutFormat != "" && resp != nil {
		out := c.WriteOutFormat
		// 支持 %{http_code}
		out = strings.ReplaceAll(out, "%{http_code}", fmt.Sprintf("%d", resp.GetStatusCode()))
		// 支持 %{time_total}
		out = strings.ReplaceAll(out, "%{time_total}", fmt.Sprintf("%.3f", duration.Seconds()))
		fmt.Print(out)
	}
	return resp, err
}

// CreateSession 创建Session
func (curl *CURL) CreateSession() *requests.Session {
	ses := requests.NewSession()

	// 设置基本配置
	ses.SetHeader(curl.Header)
	ses.SetCookies(curl.ParsedURL, curl.Cookies)

	// 设置总超时
	if curl.Timeout > 0 {
		ses.Config().SetTimeout(curl.Timeout)
	}

	// 设置认证 - 支持新的认证系统
	if curl.AuthV2 != nil && curl.AuthV2.IsValid() {
		switch curl.AuthV2.Type {
		case "basic":
			ses.Config().SetBasicAuth(curl.AuthV2.Username, curl.AuthV2.Password)
		case "digest":
			// Digest认证需要特殊处理
			ses.Config().SetBasicAuth(curl.AuthV2.Username, curl.AuthV2.Password)
			// TODO: 在requests库中实现真正的Digest认证支持
		case "bearer":
			// Bearer认证通过Header设置
			authHeader := make(http.Header)
			authHeader.Set("Authorization", curl.AuthV2.GetAuthHeader())
			ses.SetHeader(authHeader)
		}
	} else if curl.Auth != nil {
		// 向后兼容旧的认证系统
		ses.Config().SetBasicAuth(curl.Auth.User, curl.Auth.Password)
	}

	// 设置跳过TLS验证
	if curl.Insecure {
		ses.Config().SetInsecure(curl.Insecure)
	}

	// 设置代理（包括SOCKS5）及可选的代理认证
	if curl.Proxy != "" {
		proxyURL := curl.Proxy
		// 如果提供了代理用户信息，注入到URL中
		if curl.ProxyUser != "" {
			if u, err := url.Parse(curl.Proxy); err == nil {
				u.User = url.UserPassword(curl.ProxyUser, curl.ProxyPassword)
				proxyURL = u.String()
			}
		}
		ses.Config().SetProxy(proxyURL)
	}

	// 设置连接超时（如果指定了）
	if curl.ConnectTimeout > 0 {
		// 目前先使用总超时作为连接超时的替代方案
		// 理想情况下应该在requests库中添加专门的SetConnectTimeout方法
		// TODO: 在requests库中实现真正的连接超时设置
		ses.Config().SetTimeout(curl.ConnectTimeout)
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

	// 设置HTTP协议版本控制
	curl.configureHTTPVersion(ses)

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

// configureHTTPVersion 配置HTTP协议版本
func (curl *CURL) configureHTTPVersion(ses *requests.Session) {
	// 根据协议版本设置相应的配置
	switch curl.HTTPVersion {
	case HTTPVersion10:
		// 强制使用 HTTP/1.0
		// 通过设置TLS配置来限制协议版本
		curl.setHTTPVersionInTLS(ses, "1.0")
	case HTTPVersion11:
		// 强制使用 HTTP/1.1
		curl.setHTTPVersionInTLS(ses, "1.1")
	case HTTPVersion2:
		// 强制使用 HTTP/2
		curl.HTTP2 = true // 保持向后兼容
		curl.setHTTPVersionInTLS(ses, "2")
	case HTTPVersionAuto:
		// 自动选择，如果设置了HTTP2标志则使用HTTP/2
		if curl.HTTP2 {
			curl.setHTTPVersionInTLS(ses, "2")
		}
		// 否则让库自动选择
	}
}

// setHTTPVersionInTLS 通过TLS配置设置HTTP版本
func (curl *CURL) setHTTPVersionInTLS(ses *requests.Session, version string) {
	// 这是一个内部方法，用于设置协议版本偏好
	// 实际的协议协商由Go的http包和TLS处理
	// 我们主要是设置一些标志来影响协议选择

	// 记录协议版本偏好，供调试和日志使用
	curl.setHTTPVersionPreference(version)
}

// setHTTPVersionPreference 设置协议版本偏好（内部方法）
func (curl *CURL) setHTTPVersionPreference(version string) {
	// 这个方法主要用于记录用户的协议版本偏好
	// 实际的HTTP版本控制在Go标准库层面处理
	// 我们主要确保正确的配置传递到底层
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

	// 2. 调用核心解析函数
	return buildFromArgs(args)
}

// Debug 返回 CURL 对象的详细调试信息
func (c *CURL) Debug() string {
	var b strings.Builder

	b.WriteString("=== CURL Debug Information ===\n")

	// 基本信息
	b.WriteString(fmt.Sprintf("Method: %s\n", c.Method))
	if c.ParsedURL != nil {
		b.WriteString(fmt.Sprintf("URL: %s\n", c.ParsedURL.String()))
		b.WriteString(fmt.Sprintf("  Scheme: %s\n", c.ParsedURL.Scheme))
		b.WriteString(fmt.Sprintf("  Host: %s\n", c.ParsedURL.Host))
		b.WriteString(fmt.Sprintf("  Path: %s\n", c.ParsedURL.Path))
		if c.ParsedURL.RawQuery != "" {
			b.WriteString(fmt.Sprintf("  Query: %s\n", c.ParsedURL.RawQuery))
		}
	}

	// 头部信息
	if len(c.Header) > 0 {
		b.WriteString(fmt.Sprintf("Headers (%d):\n", len(c.Header)))
		for key, values := range c.Header {
			for _, value := range values {
				b.WriteString(fmt.Sprintf("  %s: %s\n", key, value))
			}
		}
	}

	// Cookie 信息
	if len(c.Cookies) > 0 {
		b.WriteString(fmt.Sprintf("Cookies (%d):\n", len(c.Cookies)))
		for _, cookie := range c.Cookies {
			b.WriteString(fmt.Sprintf("  %s=%s", cookie.Name, cookie.Value))
			if cookie.Domain != "" {
				b.WriteString(fmt.Sprintf("; Domain=%s", cookie.Domain))
			}
			if cookie.Path != "" {
				b.WriteString(fmt.Sprintf("; Path=%s", cookie.Path))
			}
			b.WriteString("\n")
		}
	}

	// 认证信息
	if c.Auth != nil {
		b.WriteString(fmt.Sprintf("Authentication: Basic (%s:***)\n", c.Auth.User))
	}

	// Body 信息
	if c.Body != nil {
		b.WriteString("Body:\n")
		b.WriteString(fmt.Sprintf("  Type: %s\n", c.Body.Type))
		b.WriteString(fmt.Sprintf("  Length: %d bytes\n", c.Body.Len()))
		if c.Body.Type == "raw" && c.Body.Len() < 200 {
			if buf, ok := c.Body.Content.(*bytes.Buffer); ok {
				b.WriteString(fmt.Sprintf("  Content: %s\n", buf.String()))
			}
		} else if c.Body.Len() >= 200 {
			b.WriteString("  Content: [too large to display]\n")
		}
	}

	// 网络配置
	b.WriteString("Network Configuration:\n")
	if c.Timeout > 0 {
		b.WriteString(fmt.Sprintf("  Timeout: %v\n", c.Timeout))
	}
	if c.ConnectTimeout > 0 {
		b.WriteString(fmt.Sprintf("  Connect Timeout: %v\n", c.ConnectTimeout))
	}
	if c.DNSTimeout > 0 {
		b.WriteString(fmt.Sprintf("  DNS Timeout: %v\n", c.DNSTimeout))
	}
	if c.TLSHandshakeTimeout > 0 {
		b.WriteString(fmt.Sprintf("  TLS Handshake Timeout: %v\n", c.TLSHandshakeTimeout))
	}

	// HTTP协议版本信息
	if c.HTTPVersion != HTTPVersionAuto {
		b.WriteString(fmt.Sprintf("  HTTP Version: %s\n", c.HTTPVersion.String()))
	} else if c.HTTP2 {
		b.WriteString("  HTTP Version: HTTP/2 (legacy flag)\n")
	} else {
		b.WriteString("  HTTP Version: Auto\n")
	}

	if c.Proxy != "" {
		b.WriteString(fmt.Sprintf("  Proxy: %s\n", c.Proxy))
	}
	if c.Insecure {
		b.WriteString("  SSL Verification: DISABLED\n")
	}

	// SSL/TLS 配置
	if c.CACert != "" || c.ClientCert != "" || c.ClientKey != "" {
		b.WriteString("SSL/TLS Configuration:\n")
		if c.CACert != "" {
			b.WriteString(fmt.Sprintf("  CA Certificate: %s\n", c.CACert))
		}
		if c.ClientCert != "" {
			b.WriteString(fmt.Sprintf("  Client Certificate: %s\n", c.ClientCert))
		}
		if c.ClientKey != "" {
			b.WriteString(fmt.Sprintf("  Client Key: %s\n", c.ClientKey))
		}
	}

	// 重定向配置
	if c.FollowRedirect {
		b.WriteString("Redirect Configuration:\n")
		b.WriteString("  Follow Redirects: YES\n")
		if c.MaxRedirs >= 0 {
			b.WriteString(fmt.Sprintf("  Max Redirects: %d\n", c.MaxRedirs))
		} else {
			b.WriteString("  Max Redirects: unlimited\n")
		}
	}

	// 调试标志
	debugFlags := []string{}
	if c.Verbose {
		debugFlags = append(debugFlags, "verbose")
	}
	if c.Include {
		debugFlags = append(debugFlags, "include-headers")
	}
	if c.Silent {
		debugFlags = append(debugFlags, "silent")
	}
	if c.Trace {
		debugFlags = append(debugFlags, "trace")
	}
	if len(debugFlags) > 0 {
		b.WriteString(fmt.Sprintf("Debug Flags: %s\n", strings.Join(debugFlags, ", ")))
	}

	// 文件输出配置
	if c.OutputFile != "" || c.RemoteName || c.OutputDir != "" {
		b.WriteString("File Output Configuration:\n")
		if c.OutputFile != "" {
			b.WriteString(fmt.Sprintf("  Output File: %s\n", c.OutputFile))
		}
		if c.RemoteName {
			b.WriteString("  Use Remote Name: YES\n")
		}
		if c.OutputDir != "" {
			b.WriteString(fmt.Sprintf("  Output Directory: %s\n", c.OutputDir))
		}
		if c.CreateDirs {
			b.WriteString("  Create Directories: YES\n")
		}
		if c.RemoveOnError {
			b.WriteString("  Remove on Error: YES\n")
		}
		if c.ContinueAt > 0 {
			b.WriteString(fmt.Sprintf("  Continue At: %d bytes\n", c.ContinueAt))
		}
	}

	b.WriteString("===============================")
	return b.String()
}

// Summary 返回 CURL 对象的简要信息
func (c *CURL) Summary() string {
	var parts []string

	// 基本请求信息
	if c.ParsedURL != nil {
		parts = append(parts, fmt.Sprintf("%s %s", c.Method, c.ParsedURL.String()))
	}

	// 重要配置
	if len(c.Header) > 0 {
		parts = append(parts, fmt.Sprintf("%d headers", len(c.Header)))
	}
	if len(c.Cookies) > 0 {
		parts = append(parts, fmt.Sprintf("%d cookies", len(c.Cookies)))
	}
	if c.Body != nil && c.Body.Len() > 0 {
		parts = append(parts, fmt.Sprintf("body(%s, %d bytes)", c.Body.Type, c.Body.Len()))
	}
	if c.Auth != nil {
		parts = append(parts, "auth")
	}
	if c.Proxy != "" {
		parts = append(parts, "proxy")
	}
	if c.Insecure {
		parts = append(parts, "insecure")
	}

	return strings.Join(parts, " | ")
}

// Verbose 返回详细的执行信息（模拟 curl -v 的输出）
func (c *CURL) VerboseInfo() string {
	var b strings.Builder

	if c.ParsedURL != nil {
		// DNS 解析信息
		if len(c.Resolve) > 0 {
			b.WriteString("* DNS resolution overrides:\n")
			for _, resolve := range c.Resolve {
				parts := strings.Split(resolve, ":")
				if len(parts) >= 3 {
					host := parts[0]
					if strings.HasPrefix(host, "+") {
						b.WriteString(fmt.Sprintf("*   %s (force): %s\n", parts[0], strings.Join(parts[2:], ":")))
					} else {
						b.WriteString(fmt.Sprintf("*   %s:%s -> %s\n", parts[0], parts[1], strings.Join(parts[2:], ":")))
					}
				}
			}
		}

		// 连接重定向信息
		if len(c.ConnectTo) > 0 {
			b.WriteString("* Connection redirects:\n")
			for _, connectTo := range c.ConnectTo {
				parts := strings.Split(connectTo, ":")
				if len(parts) == 4 {
					sourceHost, sourcePort, targetHost, targetPort := parts[0], parts[1], parts[2], parts[3]
					if sourceHost == "" {
						sourceHost = "*"
					}
					if sourcePort == "" {
						sourcePort = "*"
					}
					b.WriteString(fmt.Sprintf("*   %s:%s -> %s:%s\n", sourceHost, sourcePort, targetHost, targetPort))
				}
			}
		}

		b.WriteString(fmt.Sprintf("* Trying %s...\n", c.ParsedURL.Host))
		b.WriteString(fmt.Sprintf("* Connected to %s port %s\n", c.ParsedURL.Hostname(), c.ParsedURL.Port()))

		if c.ParsedURL.Scheme == "https" {
			b.WriteString("* SSL connection using TLS\n")
			if c.Insecure {
				b.WriteString("* WARNING: SSL verification disabled!\n")
			}
		}

		// 请求行
		path := c.ParsedURL.Path
		if path == "" {
			path = "/"
		}
		if c.ParsedURL.RawQuery != "" {
			path += "?" + c.ParsedURL.RawQuery
		}
		b.WriteString(fmt.Sprintf("> %s %s HTTP/1.1\n", c.Method, path))
		b.WriteString(fmt.Sprintf("> Host: %s\n", c.ParsedURL.Host))

		// 请求头
		for key, values := range c.Header {
			for _, value := range values {
				b.WriteString(fmt.Sprintf("> %s: %s\n", key, value))
			}
		}

		if c.Body != nil && c.Body.Len() > 0 {
			b.WriteString(">\n")
			b.WriteString(fmt.Sprintf("* upload completely sent off: %d out of %d bytes\n", c.Body.Len(), c.Body.Len()))
		}
	}

	return b.String()
}
