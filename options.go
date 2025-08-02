package gcurl

import (
	"bytes"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/474420502/requests"
)

// OptionHandler 是一个处理函数的类型签名。
// 它接收当前的CURL对象和该选项所需的参数。
// 返回值是它消费掉的参数数量，以及一个error。
type OptionHandler func(c *CURL, args ...string) error

// OptionSpec 定义了一个选项的规范
type OptionSpec struct {
	Handler                OptionHandler // 处理逻辑
	NumArgs                int           // 需要的参数数量 (-1 表示可变)
	CanAppearMultipleTimes bool          // 是否可以出现多次（如 -H）
}

// optionRegistry 是所有支持的 cURL 选项的注册中心
var optionRegistry = make(map[string]OptionSpec)

// init 函数是注册选项的绝佳位置
func init() {
	// --header / -H
	headerSpec := OptionSpec{Handler: handleHeader, NumArgs: 1, CanAppearMultipleTimes: true}
	optionRegistry["-H"] = headerSpec
	optionRegistry["--header"] = headerSpec

	// --insecure / -k
	insecureSpec := OptionSpec{Handler: handleInsecure, NumArgs: 0}
	optionRegistry["-k"] = insecureSpec
	optionRegistry["--insecure"] = insecureSpec

	// --request / -X
	methodSpec := OptionSpec{Handler: handleMethod, NumArgs: 1}
	optionRegistry["-X"] = methodSpec
	optionRegistry["--request"] = methodSpec

	// --data / -d / --data-ascii
	dataSpec := OptionSpec{Handler: handleData, NumArgs: 1}
	optionRegistry["-d"] = dataSpec
	optionRegistry["--data"] = dataSpec
	optionRegistry["--data-ascii"] = dataSpec // --data-ascii 是 --data 的别名

	// ==========================================================
	//  在此处添加对 --data-binary 的注册
	// ==========================================================
	dataBinarySpec := OptionSpec{
		Handler: handleDataBinary,
		NumArgs: 1, // 需要1个参数 (数据或 @filename)
	}
	optionRegistry["--data-binary"] = dataBinarySpec

	compressedSpec := OptionSpec{
		Handler: handleCompressed,
		NumArgs: 0, // --compressed 选项本身不需要参数
	}
	optionRegistry["--compressed"] = compressedSpec

	socks5Spec := OptionSpec{
		Handler: handleSocks5,
		NumArgs: 1, // --socks5 需要一个参数（代理服务器地址）
	}
	optionRegistry["--socks5"] = socks5Spec

	connectTimeoutSpec := OptionSpec{
		Handler: handleConnectTimeout,
		NumArgs: 1, // 需要一个参数
	}
	optionRegistry["--connect-timeout"] = connectTimeoutSpec

	// ==========================================================
	//  在此处添加对 --data-urlencode 的注册
	// ==========================================================
	dataUrlencodeSpec := OptionSpec{
		Handler: handleDataUrlencode,
		NumArgs: 1, // 需要一个参数
	}
	optionRegistry["--data-urlencode"] = dataUrlencodeSpec

	// ==========================================================
	//  在此处添加对 --data-raw 的注册
	// ==========================================================
	dataRawSpec := OptionSpec{
		Handler: handleDataRaw,
		NumArgs: 1, // 需要一个参数
	}
	optionRegistry["--data-raw"] = dataRawSpec

	// --user / -u
	userSpec := OptionSpec{Handler: handleUser, NumArgs: 1}
	optionRegistry["-u"] = userSpec
	optionRegistry["--user"] = userSpec

	// --user-agent / -A
	userAgentSpec := OptionSpec{Handler: handleUserAgent, NumArgs: 1}
	optionRegistry["--user-agent"] = userAgentSpec
	optionRegistry["-A"] = userAgentSpec

	// --limit-rate
	limitRateSpec := OptionSpec{Handler: handleLimitRate, NumArgs: 1}
	optionRegistry["--limit-rate"] = limitRateSpec

	// --cookie / -b
	cookieSpec := OptionSpec{Handler: handleCookie, NumArgs: 1}
	optionRegistry["-b"] = cookieSpec
	optionRegistry["--cookie"] = cookieSpec

	// --head / -I
	headSpec := OptionSpec{Handler: handleHead, NumArgs: 0}
	optionRegistry["-I"] = headSpec
	optionRegistry["--head"] = headSpec

	// --form / -F (文件上传)
	formSpec := OptionSpec{Handler: handleForm, NumArgs: 1, CanAppearMultipleTimes: true}
	optionRegistry["-F"] = formSpec
	optionRegistry["--form"] = formSpec

	// --location / -L (重定向跟随)
	locationSpec := OptionSpec{Handler: handleLocation, NumArgs: 0}
	optionRegistry["-L"] = locationSpec
	optionRegistry["--location"] = locationSpec

	// --max-time (最大执行时间)
	maxTimeSpec := OptionSpec{Handler: handleMaxTime, NumArgs: 1}
	optionRegistry["--max-time"] = maxTimeSpec

	// --proxy (HTTP代理)
	proxySpec := OptionSpec{Handler: handleProxy, NumArgs: 1}
	optionRegistry["--proxy"] = proxySpec
	optionRegistry["-x"] = proxySpec
	// --proxy-user / -U (代理认证 用户:密码)
	proxyUserSpec := OptionSpec{Handler: handleProxyUser, NumArgs: 1}
	optionRegistry["-U"] = proxyUserSpec
	optionRegistry["--proxy-user"] = proxyUserSpec

	// --max-redirs (最大重定向次数)
	maxRedirsSpec := OptionSpec{Handler: handleMaxRedirs, NumArgs: 1}
	optionRegistry["--max-redirs"] = maxRedirsSpec

	// --cacert (自定义CA证书)
	cacertSpec := OptionSpec{Handler: handleCACert, NumArgs: 1}
	optionRegistry["--cacert"] = cacertSpec

	// --cert (客户端证书)
	certSpec := OptionSpec{Handler: handleClientCert, NumArgs: 1}
	optionRegistry["--cert"] = certSpec

	// --verbose / -v (详细输出)
	verboseSpec := OptionSpec{Handler: handleVerbose, NumArgs: 0}
	optionRegistry["-v"] = verboseSpec
	optionRegistry["--verbose"] = verboseSpec

	// --include / -i (包含响应头)
	includeSpec := OptionSpec{Handler: handleInclude, NumArgs: 0}
	optionRegistry["-i"] = includeSpec
	optionRegistry["--include"] = includeSpec

	// --silent / -s (静默模式)
	silentSpec := OptionSpec{Handler: handleSilent, NumArgs: 0}
	optionRegistry["-s"] = silentSpec
	optionRegistry["--silent"] = silentSpec

	// --trace
	traceSpec := OptionSpec{Handler: handleTrace, NumArgs: 0}
	optionRegistry["--trace"] = traceSpec

	// --digest (强制Digest认证)
	digestSpec := OptionSpec{Handler: handleDigest, NumArgs: 1}
	optionRegistry["--digest"] = digestSpec

	// --key (客户端私钥)
	keySpec := OptionSpec{Handler: handleClientKey, NumArgs: 1}
	optionRegistry["--key"] = keySpec

	// --http2 (强制HTTP/2)
	http2Spec := OptionSpec{Handler: handleHTTP2, NumArgs: 0}
	optionRegistry["--http2"] = http2Spec

	// --http1.1 (强制HTTP/1.1)
	http11Spec := OptionSpec{Handler: handleHTTP11, NumArgs: 0}
	optionRegistry["--http1.1"] = http11Spec

	// --http1.0 (强制HTTP/1.0)
	http10Spec := OptionSpec{Handler: handleHTTP10, NumArgs: 0}
	optionRegistry["--http1.0"] = http10Spec

	// 文件输出相关选项
	// -o/--output (指定输出文件)
	outputSpec := OptionSpec{Handler: handleOutput, NumArgs: 1}
	optionRegistry["-o"] = outputSpec
	optionRegistry["--output"] = outputSpec

	// -O/--remote-name (使用远程文件名)
	remoteNameSpec := OptionSpec{Handler: handleRemoteName, NumArgs: 0}
	optionRegistry["-O"] = remoteNameSpec
	optionRegistry["--remote-name"] = remoteNameSpec

	// --output-dir (指定输出目录)
	outputDirSpec := OptionSpec{Handler: handleOutputDir, NumArgs: 1}
	optionRegistry["--output-dir"] = outputDirSpec

	// --create-dirs (自动创建目录)
	createDirsSpec := OptionSpec{Handler: handleCreateDirs, NumArgs: 0}
	optionRegistry["--create-dirs"] = createDirsSpec

	// --remove-on-error (出错时删除文件)
	removeOnErrorSpec := OptionSpec{Handler: handleRemoveOnError, NumArgs: 0}
	optionRegistry["--remove-on-error"] = removeOnErrorSpec

	// -C/--continue-at (断点续传)
	continueAtSpec := OptionSpec{Handler: handleContinueAt, NumArgs: 1}
	optionRegistry["-C"] = continueAtSpec
	optionRegistry["--continue-at"] = continueAtSpec

	// 注册脚本与易用性增强选项
	// --write-out / -w (格式化输出)
	writeOutSpec := OptionSpec{Handler: handleWriteOut, NumArgs: 1}
	optionRegistry["-w"] = writeOutSpec
	optionRegistry["--write-out"] = writeOutSpec
	// --remote-header-name / -J (使用 Content-Disposition 中的文件名)
	remoteHeaderNameSpec := OptionSpec{Handler: handleRemoteHeaderName, NumArgs: 0}
	optionRegistry["-J"] = remoteHeaderNameSpec
	optionRegistry["--remote-header-name"] = remoteHeaderNameSpec
	// --fail / -f (脚本错误处理)
	failSpec := OptionSpec{Handler: handleFail, NumArgs: 0}
	optionRegistry["-f"] = failSpec
	optionRegistry["--fail"] = failSpec

	// --resolve (主机名解析映射)
	resolveSpec := OptionSpec{Handler: handleResolve, NumArgs: 1, CanAppearMultipleTimes: true}
	optionRegistry["--resolve"] = resolveSpec
}

// --- 具体的 Handler 实现 ---

func handleHeader(c *CURL, args ...string) error {
	headerValue := args[0]

	// 忽略空头部（与 curl 行为一致）
	if strings.TrimSpace(headerValue) == "" {
		return nil
	}

	key, value, err := parseHTTPHeaderKeyValue(headerValue)
	if err != nil {
		return fmt.Errorf("invalid header format: %w", err)
	}

	lkey := strings.ToLower(key)
	switch lkey {
	case "cookie":
		// 1. 仍然添加原始Header，保持与curl行为一致
		c.Header.Add(key, value)
		// 2. 调用 GetRawCookies 解析，并暂存到CURL对象中
		//    此时不操作CookieJar，避免URL未解析的风险
		parsedCookies := GetRawCookies(value, "") // GetRawCookies 在 cookie.go 中
		c.Cookies = append(c.Cookies, parsedCookies...)

	case "content-type":
		// 对Content-Type使用Set而不是Add，因为它应该是唯一的
		c.Header.Set(key, value)
		c.ContentType = value

	default:
		c.Header.Add(key, value)
	}

	return nil
}
func handleInsecure(c *CURL, args ...string) error {
	c.Insecure = true
	return nil
}

func handleHead(c *CURL, args ...string) error {
	c.Method = "HEAD"
	return nil
}

func handleMethod(c *CURL, args ...string) error {
	c.Method = strings.ToUpper(args[0])
	return nil
}

// handleData 用于处理 -d, --data, --data-ascii
// 根据 cURL 文档，这些选项在从文件读取时会删除换行符。
// 多个--data选项会被连接起来
func handleData(c *CURL, args ...string) error {
	if c.Method == "" {
		c.Method = "POST"
	}
	// --data 通常意味着内容类型为 application/x-www-form-urlencoded
	if c.Header.Get("Content-Type") == "" {
		c.Header.Set("Content-Type", requests.TypeURLENCODED)
		c.ContentType = requests.TypeURLENCODED
	}

	data := args[0]
	var content []byte

	// 检查 @filename 语法
	if strings.HasPrefix(data, "@") {
		filePath := data[1:]
		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read file for --data: %w", err)
		}
		// cURL 的 --data 会删除换行符
		content = bytes.ReplaceAll(fileContent, []byte("\n"), []byte(""))
		content = bytes.ReplaceAll(content, []byte("\r"), []byte(""))
	} else {
		content = []byte(data)
	}

	// 检查是否已经有body数据，如果有则追加
	if c.Body != nil && c.Body.Type == "raw" {
		// 获取现有内容
		existingData := c.Body.String()
		if existingData != "" {
			// 对于form数据，使用&连接
			if c.ContentType == requests.TypeURLENCODED {
				content = []byte(existingData + "&" + string(content))
			} else {
				// 对于其他类型，直接连接
				content = []byte(existingData + string(content))
			}
		}
	}

	c.setRawBody(content)
	return nil
}

// handleDataBinary 用于处理 --data-binary
// 它会发送原始数据，不会删除文件中的换行符。
// 多个--data-binary选项会被连接起来
func handleDataBinary(c *CURL, args ...string) error {
	if c.Method == "" {
		c.Method = "POST"
	}

	// 从参数中获取原始数据字符串
	data := args[0]
	var content []byte

	// 检查 @filename 语法
	if strings.HasPrefix(data, "@") {
		filePath := data[1:]
		var err error
		content, err = os.ReadFile(filePath) // 直接读取文件原始字节
		if err != nil {
			return fmt.Errorf("failed to read file for --data-binary: %w", err)
		}
		// 当从文件读取时，不假设 Content-Type，让用户通过 -H 自行指定。
	} else {
		content = []byte(data)
		// 为了兼容性，如果用户未通过 -H 设置 Content-Type，我们不主动设置。
		// 在 multipart 的场景下，Header里已经有带 boundary 的 Content-Type 了。
		// 对于其他情况，让 net/http 自动检测或保持为空。
	}

	// 检查是否已经有body数据，如果有则追加
	if c.Body != nil && c.Body.Type == "raw" {
		// 获取现有内容
		existingData := c.Body.String()
		if existingData != "" {
			// 对于二进制数据，直接连接（不添加&）
			content = []byte(existingData + string(content))
		}
	}

	c.setRawBody(content)

	// 同步一下便利字段 ContentType, 确保它与 Header 一致
	// 注意：这里不要覆盖已有的Content-Type
	if c.ContentType == "" {
		c.ContentType = c.Header.Get("Content-Type")
	}

	return nil
}

// handleCompressed 用于处理 --compressed 选项
func handleCompressed(c *CURL, args ...string) error {
	// 设置标准的 Accept-Encoding 头，告诉服务器我们接受这些压缩格式
	// 你的 requests 库支持 gzip, deflate, 和 br
	c.Header.Set("Accept-Encoding", "gzip, deflate, br")
	return nil
}

// handleSocks5 用于处理 --socks5 选项
func handleSocks5(c *CURL, args ...string) error {
	proxyAddr := args[0]
	// 你的 SetProxy 函数需要一个完整的 URL（例如 socks5://localhost:1080）
	// 为用户着想，如果他们只提供了 host:port，我们自动添加 scheme
	if !strings.HasPrefix(proxyAddr, "socks5://") {
		proxyAddr = "socks5://" + proxyAddr
	}
	c.Proxy = proxyAddr
	return nil
}

// handleConnectTimeout 用于处理 --connect-timeout 选项
func handleConnectTimeout(c *CURL, args ...string) error {
	connectTimeoutStr := args[0]

	// 支持单位：s(秒), m(分), h(小时)，如 "30s", "5m", "1h"
	if strings.HasSuffix(connectTimeoutStr, "s") || strings.HasSuffix(connectTimeoutStr, "m") || strings.HasSuffix(connectTimeoutStr, "h") {
		duration, err := time.ParseDuration(connectTimeoutStr)
		if err != nil {
			return fmt.Errorf("invalid value for --connect-timeout: %w", err)
		}
		if duration < 0 {
			return fmt.Errorf("connect-timeout must be non-negative: %s", connectTimeoutStr)
		}
		c.ConnectTimeout = duration
	} else {
		// 纯数字，按秒处理
		timeout, err := strconv.Atoi(connectTimeoutStr)
		if err != nil {
			return fmt.Errorf("invalid value for --connect-timeout: %w", err)
		}
		if timeout < 0 {
			return fmt.Errorf("connect-timeout must be non-negative: %d", timeout)
		}
		c.ConnectTimeout = time.Duration(timeout) * time.Second
	}
	return nil
} // handleDataUrlencode 用于处理 --data-urlencode 选项
// 支持以下语法格式：
// content - URL编码内容
// =content - URL编码内容（与上面相同）
// name=content - URL编码内容并添加name=前缀
// @filename - URL编码文件内容
// name@filename - URL编码文件内容并添加name=前缀
func handleDataUrlencode(c *CURL, args ...string) error {
	if c.Method == "" {
		c.Method = "POST"
	}
	// --data-urlencode 通常意味着内容类型为 application/x-www-form-urlencoded
	if c.Header.Get("Content-Type") == "" {
		c.Header.Set("Content-Type", requests.TypeURLENCODED)
		c.ContentType = requests.TypeURLENCODED
	}

	data := args[0]
	var result string

	// 解析不同的语法格式
	if strings.HasPrefix(data, "=") {
		// =content 格式
		content := data[1:]
		result = url.QueryEscape(content)
	} else if strings.HasPrefix(data, "@") {
		// @filename 格式
		filename := data[1:]

		// 读取文件内容
		fileContent, err := os.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("failed to read file for --data-urlencode: %w", err)
		}

		result = url.QueryEscape(string(fileContent))
	} else if !strings.Contains(data, "=") && strings.Contains(data, "@") {
		// name@filename 格式 (没有=符号)
		parts := strings.SplitN(data, "@", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid --data-urlencode format: %s", data)
		}

		name := parts[0]
		filename := parts[1]

		// 读取文件内容
		fileContent, err := os.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("failed to read file for --data-urlencode: %w", err)
		}

		result = name + "=" + url.QueryEscape(string(fileContent))
	} else if strings.Contains(data, "=") {
		// name=content 格式
		parts := strings.SplitN(data, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid --data-urlencode format: %s", data)
		}

		name := parts[0]
		content := parts[1]
		result = name + "=" + url.QueryEscape(content)
	} else {
		// 普通content格式
		result = url.QueryEscape(data)
	} // 检查是否已经有body数据，如果有则追加
	if c.Body != nil && c.Body.Type == "raw" {
		// 获取现有内容
		existingData := c.Body.String()
		if existingData != "" {
			result = existingData + "&" + result
		}
	}

	c.setRawBodyString(result)
	return nil
}

// handleDataRaw 用于处理 --data-raw 选项
// 多个--data-raw选项会被连接起来
func handleDataRaw(c *CURL, args ...string) error {
	if c.Method == "" {
		c.Method = "POST"
	}
	// --data-raw 通常意味着内容类型为 application/x-www-form-urlencoded
	if c.Header.Get("Content-Type") == "" {
		c.Header.Set("Content-Type", requests.TypeURLENCODED)
		c.ContentType = requests.TypeURLENCODED
	}

	// --data-raw 直接使用提供的数据，不支持 @filename 语法
	data := args[0]

	// 检查是否已经有body数据，如果有则追加
	if c.Body != nil && c.Body.Type == "raw" {
		// 获取现有内容
		existingData := c.Body.String()
		if existingData != "" {
			// 对于form数据，使用&连接
			if c.ContentType == requests.TypeURLENCODED {
				data = existingData + "&" + data
			} else {
				// 对于其他类型，直接连接
				data = existingData + data
			}
		}
	}

	c.setRawBodyString(data)
	return nil
}

// handleUser 用于处理 -u / --user 选项
func handleUser(c *CURL, args ...string) error {
	userpass := args[0]
	// 使用 strings.SplitN 分割，避免密码中包含 ':' 导致的问题
	parts := strings.SplitN(userpass, ":", 2)
	if len(parts) < 2 {
		// cURL 在这种情况下可能会提示输入密码，但作为一个库，我们要求格式必须完整
		return fmt.Errorf("invalid user format. Expected 'user:password', got '%s'", userpass)
	}
	if c.Auth == nil {
		c.Auth = &AuthInfo{Type: "basic"}
	}
	c.Auth.User = parts[0]
	c.Auth.Password = parts[1]
	return nil
}

// handleDigest 用于处理 --digest 选项
func handleDigest(c *CURL, args ...string) error {
	userpass := args[0]
	// 使用 strings.SplitN 分割，避免密码中包含 ':' 导致的问题
	parts := strings.SplitN(userpass, ":", 2)
	if len(parts) < 2 {
		// cURL 在这种情况下可能会提示输入密码，但作为一个库，我们要求格式必须完整
		return fmt.Errorf("invalid digest format. Expected 'user:password', got '%s'", userpass)
	}
	// 使用新的 AuthV2 系统
	auth := &AuthInfo{
		Type:     "digest",
		User:     parts[0],
		Username: parts[0],
		Password: parts[1],
	}
	c.AuthV2 = auth
	return nil
}

// handleUserAgent 用于处理 --user-agent 选项
func handleUserAgent(c *CURL, args ...string) error {
	userAgent := args[0]
	c.Header.Set("User-Agent", userAgent)
	return nil
}

// handleLimitRate 用于处理 --limit-rate 选项
func handleLimitRate(c *CURL, args ...string) error {
	// --limit-rate 选项用于限制传输速度
	// 参数可以是 bytes/second 或带单位的值 (如 200K, 1M)
	limitRate := args[0]

	// 存储限速设置，具体的限速逻辑由底层的 requests 库处理
	c.LimitRate = limitRate

	return nil
}

// handleCookie 用于处理 -b / --cookie 选项
func handleCookie(c *CURL, args ...string) error {
	cookieValue := args[0]

	// 检查是否为文件路径（以@开头）
	if strings.HasPrefix(cookieValue, "@") {
		// 从文件读取 cookies
		filePath := cookieValue[1:]
		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read cookie file: %w", err)
		}

		// 解析文件中的 cookies（支持 Netscape cookie jar 格式和简单的键值对格式）
		cookieValue = strings.TrimSpace(string(fileContent))
	}

	// 解析 cookie 字符串并添加到 Header 中
	// 对于多个 cookies，可以是分号分隔的格式：name1=value1; name2=value2
	// 或者是单个 cookie：name=value

	// 直接设置 Cookie header，让现有的 handleHeader 逻辑处理解析
	key := "Cookie"

	// 检查是否已有 Cookie header，如果有则追加
	existingCookie := c.Header.Get("Cookie")
	if existingCookie != "" {
		// 用分号和空格连接多个 cookie
		cookieValue = existingCookie + "; " + cookieValue
	}

	c.Header.Set(key, cookieValue)

	// 调用 GetRawCookies 解析，并暂存到CURL对象中
	parsedCookies := GetRawCookies(cookieValue, "") // GetRawCookies 在 cookie.go 中
	c.Cookies = append(c.Cookies, parsedCookies...)

	return nil
}

// handleForm 用于处理 --form / -F 选项 (文件上传)
func handleForm(c *CURL, args ...string) error {
	if c.Method == "" {
		c.Method = "POST"
	}

	formData := args[0]

	// 解析form数据
	field, err := parseFormData(formData)
	if err != nil {
		return fmt.Errorf("failed to parse form data: %w", err)
	}

	// 初始化multipart body数据结构（如果还没有）
	if c.Body == nil || c.Body.Type != "multipart" {
		c.Body = &BodyData{
			Type:    "multipart",
			Content: make([]*FormField, 0),
		}
		// 设置Content-Type（boundary会在执行时生成）
		c.Header.Set("Content-Type", "multipart/form-data")
		c.ContentType = "multipart/form-data"
	}

	// 添加字段到multipart数据
	if fields, ok := c.Body.Content.([]*FormField); ok {
		c.Body.Content = append(fields, field)
	} else {
		// 如果类型不匹配，重新创建
		c.Body.Content = []*FormField{field}
	}

	return nil
}

// handleLocation 用于处理 --location / -L 选项 (重定向跟随)
func handleLocation(c *CURL, args ...string) error {
	c.FollowRedirect = true
	// 如果MaxRedirs还是默认值，设置为一个合理的默认重定向次数
	if c.MaxRedirs == -1 {
		c.MaxRedirs = 30 // curl的默认值
	}
	return nil
}

// handleMaxTime 用于处理 --max-time 选项 (最大执行时间)
func handleMaxTime(c *CURL, args ...string) error {
	maxTime := args[0]

	// 解析时间值（可能是秒数或带单位的值）
	// 支持单位：s(秒), m(分), h(小时)，如 "30s", "5m", "1h"
	if strings.HasSuffix(maxTime, "s") || strings.HasSuffix(maxTime, "m") || strings.HasSuffix(maxTime, "h") {
		duration, err := time.ParseDuration(maxTime)
		if err != nil {
			return fmt.Errorf("invalid max-time value: %s", maxTime)
		}
		if duration < 0 {
			return fmt.Errorf("max-time must be non-negative: %s", maxTime)
		}
		c.Timeout = duration
	} else {
		// 纯数字，按秒处理
		timeout, err := strconv.Atoi(maxTime)
		if err != nil {
			return fmt.Errorf("invalid max-time value: %s", maxTime)
		}
		if timeout < 0 {
			return fmt.Errorf("max-time must be non-negative: %d", timeout)
		}
		c.Timeout = time.Duration(timeout) * time.Second
	}

	return nil
}

// handleProxy 用于处理 --proxy / -x 选项 (HTTP代理)
func handleProxy(c *CURL, args ...string) error {
	proxy := args[0]
	c.Proxy = proxy
	return nil
}

// handleProxyUser 用于处理 --proxy-user / -U 选项 (代理认证 用户:密码)
func handleProxyUser(c *CURL, args ...string) error {
	cred := args[0]
	parts := strings.SplitN(cred, ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid proxy-user format: expected user:password")
	}
	c.ProxyUser = parts[0]
	c.ProxyPassword = parts[1]
	return nil
}

// handleMaxRedirs 用于处理 --max-redirs 选项 (最大重定向次数)
func handleMaxRedirs(c *CURL, args ...string) error {
	maxRedirs := args[0]

	redirs, err := strconv.Atoi(maxRedirs)
	if err != nil {
		return fmt.Errorf("invalid max-redirs value: %s", maxRedirs)
	}

	if redirs < 0 {
		return fmt.Errorf("max-redirs must be non-negative: %d", redirs)
	}

	c.MaxRedirs = redirs
	return nil
}

// handleCACert 用于处理 --cacert 选项 (自定义CA证书)
func handleCACert(c *CURL, args ...string) error {
	certPath := args[0]

	// 检查证书文件是否存在
	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		return fmt.Errorf("CA certificate file not found: %s", certPath)
	}

	c.CACert = certPath
	return nil
}

// handleClientCert 用于处理 --cert 选项 (客户端证书)
func handleClientCert(c *CURL, args ...string) error {
	certPath := args[0]

	// 检查证书文件是否存在
	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		return fmt.Errorf("client certificate file not found: %s", certPath)
	}

	c.ClientCert = certPath
	return nil
}

// handleClientKey 用于处理 --key 选项 (客户端私钥)
func handleClientKey(c *CURL, args ...string) error {
	keyPath := args[0]

	// 检查私钥文件是否存在
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		return fmt.Errorf("client key file not found: %s", keyPath)
	}

	c.ClientKey = keyPath
	return nil
}

// handleWriteOut 用于处理 -w/--write-out 选项 (格式化输出)
func handleWriteOut(c *CURL, args ...string) error {
	c.WriteOutFormat = args[0]
	return nil
}

// handleRemoteHeaderName 用于处理 -J/--remote-header-name 选项 (使用 Content-Disposition 中的文件名)
func handleRemoteHeaderName(c *CURL, args ...string) error {
	c.RemoteHeaderName = true
	return nil
}

// handleFail 用于处理 -f/--fail 选项 (HTTP 错误时返回非零退出码)
func handleFail(c *CURL, args ...string) error {
	c.FailOnError = true
	return nil
}

// handleHTTP2 用于处理 --http2 选项 (强制HTTP/2)
func handleHTTP2(c *CURL, args ...string) error {
	c.HTTP2 = true
	c.HTTPVersion = HTTPVersion2
	return nil
}

// handleHTTP11 用于处理 --http1.1 选项 (强制HTTP/1.1)
func handleHTTP11(c *CURL, args ...string) error {
	c.HTTP2 = false
	c.HTTPVersion = HTTPVersion11
	return nil
}

// handleHTTP10 用于处理 --http1.0 选项 (强制HTTP/1.0)
func handleHTTP10(c *CURL, args ...string) error {
	c.HTTP2 = false
	c.HTTPVersion = HTTPVersion10
	return nil
}

// parseHTTPHeaderKeyValue 解析HTTP头部的键值对
func parseHTTPHeaderKeyValue(headerValue string) (key string, value string, err error) {
	colonIndex := strings.Index(headerValue, ":")
	if colonIndex == -1 {
		return "", "", fmt.Errorf("invalid header format: %s", headerValue)
	}

	key = strings.TrimSpace(headerValue[:colonIndex])
	value = strings.TrimSpace(headerValue[colonIndex+1:])

	// 保持与 curl 一致的引号处理行为
	// 不自动去掉引号，因为引号可能是值的一部分（如 sec-ch-ua 头部）
	// curl 会保留用户明确指定的引号

	return key, value, nil
}

// handleVerbose 用于处理 -v, --verbose 选项 (详细输出)
func handleVerbose(c *CURL, args ...string) error {
	c.Verbose = true
	return nil
}

// handleInclude 用于处理 -i, --include 选项 (包含响应头)
func handleInclude(c *CURL, args ...string) error {
	c.Include = true
	return nil
}

// handleSilent 用于处理 -s, --silent 选项 (静默模式)
func handleSilent(c *CURL, args ...string) error {
	c.Silent = true
	return nil
}

// handleTrace 用于处理 --trace 选项 (追踪模式)
func handleTrace(c *CURL, args ...string) error {
	c.Trace = true
	return nil
}

// handleOutput 用于处理 -o, --output 选项 (指定输出文件)
func handleOutput(c *CURL, args ...string) error {
	path := args[0]
	// 如果路径存在且为目录，则视为输出目录
	if info, err := os.Stat(path); err == nil && info.IsDir() {
		c.OutputDir = path
	} else {
		c.OutputFile = path
	}
	return nil
}

// handleRemoteName 用于处理 -O, --remote-name 选项 (使用远程文件名)
func handleRemoteName(c *CURL, args ...string) error {
	c.RemoteName = true
	return nil
}

// handleOutputDir 用于处理 --output-dir 选项 (指定输出目录)
func handleOutputDir(c *CURL, args ...string) error {
	c.OutputDir = args[0]
	return nil
}

// handleCreateDirs 用于处理 --create-dirs 选项 (自动创建目录)
func handleCreateDirs(c *CURL, args ...string) error {
	c.CreateDirs = true
	return nil
}

// handleRemoveOnError 用于处理 --remove-on-error 选项 (出错时删除文件)
func handleRemoveOnError(c *CURL, args ...string) error {
	c.RemoveOnError = true
	return nil
}

// handleContinueAt 用于处理 -C, --continue-at 选项 (断点续传)
func handleContinueAt(c *CURL, args ...string) error {
	continueAt := args[0]

	// 支持 "-" 表示自动检测文件大小
	if continueAt == "-" {
		// 自动检测模式，设置为-1作为标记
		c.ContinueAt = -1
		return nil
	}

	// 解析字节偏移量
	offset, err := strconv.ParseInt(continueAt, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid continue-at value: %s", continueAt)
	}

	if offset < 0 {
		return fmt.Errorf("continue-at offset must be non-negative: %d", offset)
	}

	c.ContinueAt = offset
	return nil
}

// handleResolve 用于处理 --resolve 选项 (主机名解析映射)
// 格式：--resolve [+]host:port:address[,address]...
// 例如：--resolve example.com:80:192.168.1.100
//
//	--resolve example.com:443:192.168.1.100,192.168.1.101
//	--resolve +example.com:80:192.168.1.100  (强制替换)
func handleResolve(c *CURL, args ...string) error {
	resolveMapping := args[0]

	// 基本格式验证：host:port:address
	parts := strings.Split(resolveMapping, ":")
	if len(parts) < 3 {
		return fmt.Errorf("invalid --resolve format, expected host:port:address, got: %s", resolveMapping)
	}

	// 验证端口是有效数字
	port := parts[1]
	if _, err := strconv.Atoi(port); err != nil {
		return fmt.Errorf("invalid port in --resolve: %s", port)
	}

	// 验证至少有一个地址
	addressPart := strings.Join(parts[2:], ":")
	if strings.TrimSpace(addressPart) == "" {
		return fmt.Errorf("missing address in --resolve: %s", resolveMapping)
	}

	// 将解析映射添加到CURL对象
	c.Resolve = append(c.Resolve, resolveMapping)
	return nil
}
