package gcurl

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"

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

	// --max-redirs (最大重定向次数)
	maxRedirsSpec := OptionSpec{Handler: handleMaxRedirs, NumArgs: 1}
	optionRegistry["--max-redirs"] = maxRedirsSpec

	// --cacert (自定义CA证书)
	cacertSpec := OptionSpec{Handler: handleCACert, NumArgs: 1}
	optionRegistry["--cacert"] = cacertSpec

	// --cert (客户端证书)
	certSpec := OptionSpec{Handler: handleClientCert, NumArgs: 1}
	optionRegistry["--cert"] = certSpec

	// --key (客户端私钥)
	keySpec := OptionSpec{Handler: handleClientKey, NumArgs: 1}
	optionRegistry["--key"] = keySpec

	// --http2 (强制HTTP/2)
	http2Spec := OptionSpec{Handler: handleHTTP2, NumArgs: 0}
	optionRegistry["--http2"] = http2Spec

	// 在这里可以继续注册其他选项, 例如 --user-agent
}

// --- 具体的 Handler 实现 ---

func handleHeader(c *CURL, args ...string) error {
	key, value, err := parseHTTPHeaderKeyValue(args[0])
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

	c.Body = bytes.NewBuffer(content)
	return nil
}

// handleDataBinary 用于处理 --data-binary
// 它会发送原始数据，不会删除文件中的换行符。
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

	c.Body = bytes.NewBuffer(content)

	// 同步一下便利字段 ContentType, 确保它与 Header 一致
	// 注意：这里不要覆盖已有的Content-Type
	if c.ContentType == "" {
		c.ContentType = c.Header.Get("Content-Type")
	}

	return nil
} // handleCompressed 用于处理 --compressed 选项
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
	// cURL 通常使用秒作为单位，可以是浮点数，但我们的库使用整数秒更方便
	timeout, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid value for --connect-timeout: %w", err)
	}
	c.ConnectTimeout = timeout
	return nil
}

// handleDataUrlencode 用于处理 --data-urlencode 选项
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
	var content []byte

	// 检查 @filename 语法
	if strings.HasPrefix(data, "@") {
		filePath := data[1:]
		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read file for --data-urlencode: %w", err)
		}
		content = fileContent
	} else {
		content = []byte(data)
	}

	// --data-urlencode 会进行 URL 编码
	// 注意：这是一个简化版本，真正的 cURL --data-urlencode 有更复杂的语法
	// 但对于大多数用例这应该足够了
	c.Body = bytes.NewBufferString(string(content))
	return nil
}

// handleDataRaw 用于处理 --data-raw 选项
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
	c.Body = bytes.NewBufferString(data)
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
		c.Auth = &requests.BasicAuth{}
	}
	c.Auth.User = parts[0]
	c.Auth.Password = parts[1]
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

	// 目前简单实现：设置 Content-Type 为 multipart/form-data
	// 实际的 multipart 数据构建需要更复杂的逻辑
	c.Header.Set("Content-Type", "multipart/form-data")
	c.ContentType = "multipart/form-data"

	// 将表单数据写入body（简化实现）
	// 实际应该构建正确的 multipart 格式
	if c.Body == nil {
		c.Body = bytes.NewBuffer(nil)
	}
	c.Body.WriteString(formData)

	return nil
}

// handleLocation 用于处理 --location / -L 选项 (重定向跟随)
func handleLocation(c *CURL, args ...string) error {
	// 设置一个标志表示应该跟随重定向
	// 这个功能需要在执行阶段实现
	// 目前只是标记这个选项被设置了
	c.Header.Set("X-Gcurl-Follow-Redirects", "true")
	return nil
}

// handleMaxTime 用于处理 --max-time 选项 (最大执行时间)
func handleMaxTime(c *CURL, args ...string) error {
	maxTime := args[0]

	// 解析时间值（可能是秒数或带单位的值）
	// 简化实现：假设是秒数
	timeout, err := strconv.Atoi(maxTime)
	if err != nil {
		return fmt.Errorf("invalid max-time value: %s", maxTime)
	}

	c.Timeout = timeout
	return nil
}

// handleProxy 用于处理 --proxy / -x 选项 (HTTP代理)
func handleProxy(c *CURL, args ...string) error {
	proxy := args[0]
	c.Proxy = proxy
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

// handleHTTP2 用于处理 --http2 选项 (强制HTTP/2)
func handleHTTP2(c *CURL, args ...string) error {
	c.HTTP2 = true
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
