package gcurl

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"

	"github.com/474420502/requests"
)

// CURL 信息结构
type CURL struct {
	ParsedURL *url.URL
	Method    string
	Header    http.Header
	CookieJar http.CookieJar
	Cookies   []*http.Cookie

	ContentType string
	Body        *bytes.Buffer

	Auth     *requests.BasicAuth
	Timeout  int // second
	Insecure bool

	// ITask   string
	// Crontab string
	// Name    string
}

// New new 一个 curl 出来
func New() *CURL {

	u := &CURL{}
	u.Insecure = false

	u.Header = make(http.Header)
	u.CookieJar, _ = cookiejar.New(nil)
	u.Body = bytes.NewBuffer(nil)
	u.Timeout = 30

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
	ses.SetHeader(curl.Header)
	ses.SetCookies(curl.ParsedURL, curl.Cookies)

	ses.Config().SetTimeout(curl.Timeout)

	if curl.Auth != nil {
		ses.Config().SetBasicAuth(curl.Auth)
	}

	if curl.Insecure {
		ses.Config().SetInsecure(curl.Insecure)
	}

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
	case "GET":
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

func checkCmdForamt2(scurl string) bool {
	i := 0
	count := 0
	for i < len(scurl) {
		c := scurl[i]
		if c == '^' {
			if i+3 < len(scurl) && scurl[i+1] == '\\' && scurl[i+2] == '^' {
				// 处理 ^\\^"
				count += 4
				i += 4
			} else if i+2 < len(scurl) && scurl[i+2] == '^' {
				// ^%^ 处理这种字符串转换
				count += 3
				i += 3
			} else if i+1 < len(scurl) {
				// 处理 ^" 特殊的把符号转换为regexp能识别的格式
				count += 2
				i += 2
			} else {
				i++
			}
		} else {
			i++
		}
	}

	return count > 0
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

// ParseBash curl bash  *Supports copying as cURL command only (Bash)
func ParseBash(scurl string) (curl *CURL, err error) {
	opts := parseCurlCommandStr(scurl)
	log.Println(scurl, opts)
	executor := newPQueueExecute()

	if len(scurl) <= 4 {
		err = fmt.Errorf("scurl error: %s", scurl)
		log.Println(err)
		return nil, err
	}

	if scurl[0] == '"' && scurl[len(scurl)-1] == '"' {
		scurl = strings.Trim(scurl, `"`)
	} else if scurl[0] == '\'' && scurl[len(scurl)-1] == '\'' {
		scurl = strings.Trim(scurl, `'`)
	}

	scurl = strings.TrimSpace(scurl)
	scurl = strings.TrimLeft(scurl, "curl")

	pattern := `((?:http|https)://[^\n\s]+(?:[\n \t]|$))|` +
		`(-(?:O|L|I|s|k|C|4|6)(?:[\n \t]|$))|` +
		`(--(?:remote-name|location|head|silent|insecure|continue-at|ipv4|ipv6|compressed)(?:[\n \t]|$))|` +
		`(--data-binary +\$.+--\\r\\n'(?:[\n \t]|$))|` +
		`(--[^ ]+ +'[^']+'(?:[\n \t]|$))|` +
		`(--[^ ]+ +"[^"]+"(?:[\n \t]|$))|` +
		`(--[^ ]+ +[^ ]+)|` +
		`(-[A-Za-z] +'[^']+'(?:[\n \t]|$))|` +
		`(-[A-Za-z] +"[^"]+"(?:[\n \t]|$))|` +
		`(-[A-Za-z] +[^ ]+)|` +
		`([\n \t]'[^']+'(?:[\n \t]|$))|` +
		`([\n \t]"[^"]+"(?:[\n \t]|$))|` +
		`(--[a-z]+ {0,})`

	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(scurl, -1)
	if len(matches) != 0 {
		curl = New()
	}
	// args := parseCurlCommandStr(scurl)
	// log.Println(args)
	for _, match := range matches {
		for i, matchedContent := range match[1:] {
			// 忽略空字符串
			if matchedContent == "" {
				continue
			}
			matchedContent = strings.Trim(matchedContent, " \n\t")

			// 使用 MatchGroup 常量替换 matchedGroup 字符串
			switch MatchGroup(i) {
			case HTTPHTTPS, NewlineQuotes, NewlineDoubleQuotes:
				purl, err := url.Parse(strings.Trim(matchedContent, `"'`))
				if err != nil {
					log.Println(err)
					return nil, err
				}
				curl.ParsedURL = purl

			case DataBinary,
				LongArgQuotes, LongArgDoubleQuotes, LongArgNoQuotes,
				ShortArgQuotes, ShortArgDoubleQuotes, ShortArgNoQuotes,
				LongArgNoArg:
				exec := judgeOptions(curl, matchedContent)
				if exec != nil {
					executor.Push(exec)
				}
			case ShortNoArg, LongNoArgSpecial:
				switch matchedContent {
				case "-I", "--head":
					curl.Method = "HEAD"
				default:
					log.Println(matchedContent, "this option is invalid.")
				}
			}
		}
	}

	for executor.Len() > 0 {
		exec := executor.Pop()
		if err = exec.Execute(); err != nil {
			return nil, err
		}
	}

	if curl.Method == "" {
		curl.Method = "GET"
	}

	return curl, nil
}
