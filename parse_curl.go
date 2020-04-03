package gcurl

import (
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
	Body      *requests.Body
	Auth      *requests.BasicAuth
	Timeout   int // second
	Insecure  bool

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
	u.Body = requests.NewBody()
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

// CreateSession 创建Session
func (curl *CURL) CreateSession() *requests.Session {
	ses := requests.NewSession()
	ses.SetHeader(curl.Header)
	ses.SetCookies(curl.ParsedURL, curl.Cookies)
	ses.SetConfig(requests.CRequestTimeout, curl.Timeout)

	if curl.Auth != nil {
		ses.SetConfig(requests.CBasicAuth, curl.Auth)
	}

	if curl.Insecure {
		ses.SetConfig(requests.CInsecure, curl.Insecure)
	}

	return ses
}

// CreateWorkflow 根据Session 创建Workflow
func (curl *CURL) CreateWorkflow(ses *requests.Session) *requests.Workflow {
	var wf *requests.Workflow

	if ses == nil {
		ses = curl.CreateSession()
	}

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

	wf.SetBody(curl.Body)
	return wf
}

// Workflow 根据自己CreateSession 创建Workflow
func (curl *CURL) Workflow() *requests.Workflow {
	return curl.CreateWorkflow(curl.CreateSession())
}

// ParseRawCURL curl_bash
func ParseRawCURL(scurl string) (cURL *CURL) {
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

	mathches := regexp.MustCompile(
		`--[^ ]+ +'[^']+' |`+
			`--[^ ]+ +"[^"]+" |`+
			`--[^ ]+ +[^ ]+|`+

			`-[A-Za-z] +'[^']+' |`+
			`-[A-Za-z] +"[^"]+" |`+
			`-[A-Za-z] +[^ ]+|`+

			` '[^']+' |`+
			` "[^"]+" |`+
			`--[a-z]+ {0,}`,
	).FindAllString(scurl, -1)
	for _, m := range mathches {
		m = strings.TrimSpace(m)
		switch v := m[0]; v {
		case '\'':
			purl, err := url.Parse(m[1 : len(m)-1])
			if err != nil {
				panic(err)
			}
			curl.ParsedURL = purl
		case '"':
			purl, err := url.Parse(m[1 : len(m)-1])
			if err != nil {
				panic(err)
			}
			curl.ParsedURL = purl
		case '-':
			exec := judgeOptions(curl, m)
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
