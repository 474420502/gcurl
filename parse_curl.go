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
	return Parse(curlbash).CreateTemporary(nil).Execute()
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

// Parse curl_bash
func Parse(scurl string) (cURL *CURL) {
	executor := newPQueueExecute()
	curl := New()

	if len(scurl) <= 4 {
		panic("scurl error:" + scurl)
	}

	if scurl[0] == '"' && scurl[len(scurl)-1] == '"' {
		scurl = strings.Trim(scurl, `"`)
	} else if scurl[0] == '\'' && scurl[len(scurl)-1] == '\'' {
		scurl = strings.Trim(scurl, `'`)
	}

	scurl = strings.TrimSpace(scurl)
	scurl = strings.TrimLeft(scurl, "curl")

	pattern := regexp.MustCompile(
		`(-(?:O|L|I|s|k|C|4|6)([\n \t]|$))|` +
			`(--(?:remote-name|location|head|silent|insecure|continue-at|ipv4|ipv6|compressed)([\n \t]|$))|` +
			`(http.+(?:[\n \t]|$))|` +
			`(--data-binary +\$.+--\\r\\n'(?:[\n \t]|$))|` +
			`(--[^ ]+ +'[^']+'(?:[\n \t]|$))|` +
			`(--[^ ]+ +"[^"]+"(?:[\n \t]|$))|` +
			`(--[^ ]+ +[^ ]+)|` +
			`(-[A-Za-z] +'[^']+'(?:[\n \t]|$))|` +
			`(-[A-Za-z] +"[^"]+"(?:[\n \t]|$))|` +
			`(-[A-Za-z] +[^ ]+)|` +
			`([\n \t]'[^']+'(?:[\n \t]|$))|` +
			`([\n \t]"[^"]+"(?:[\n \t]|$))|` +
			`(--[a-z]+ {0,})`,
	)
	matches := pattern.FindAllStringSubmatch(scurl, -1)

	groupNames := map[int]string{
		1:  "short_no_arg",
		2:  "long_no_arg",
		3:  "http_https",
		4:  "data_binary",
		5:  "long_arg_quotes",
		6:  "long_arg_double_quotes",
		7:  "long_arg_no_quotes",
		8:  "short_arg_quotes",
		9:  "short_arg_double_quotes",
		10: "short_arg_no_quotes",
		11: "newline_quotes",
		12: "newline_double_quotes",
		13: "long_arg_no_arg",
	}

	for _, submatches := range matches {
		matchedGroup := ""
		matchedContent := ""
		for i, m := range submatches[1:] {
			if m != "" {
				matchedGroup = groupNames[i+1]
				matchedContent = m
				break
			}
		}
		matchedContent = strings.Trim(matchedContent, " \n\t")
		switch matchedGroup {
		case "http_https", "newline_quotes", "newline_double_quotes":
			purl, err := url.Parse(strings.Trim(matchedContent, `"'`))
			if err != nil {
				panic(err)
			}
			curl.ParsedURL = purl
		case "short_no_arg", "long_no_arg", "data_binary",
			"long_arg_quotes", "long_arg_double_quotes", "long_arg_no_quotes",
			"short_arg_quotes", "short_arg_double_quotes", "short_arg_no_quotes",
			"long_arg_no_arg":
			exec := judgeOptions(curl, matchedContent)
			if exec != nil {
				executor.Push(exec)
			}
		}
	}

	for executor.Len() > 0 {
		exec := executor.Pop()
		exec.Execute()
	}

	if curl.Method == "" {
		curl.Method = "GET"
	}

	return curl
}
