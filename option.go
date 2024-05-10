package gcurl

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/474420502/requests"
)

func init() {
	optionTrie = newTrie()
	oelist := []*optionExecute{
		{"-H", 10, parseHeader, nil},
		{"-X", 10, parseOptX, nil},
		{"-A", 15, parseUserAgent, &extract{re: "^-A +(.+)", execute: extractData}},
		{"-d", 10, parseBodyASCII, &extract{re: "^-d +(.+)", execute: extractData}},
		{"-u", 15, parseUser, &extract{re: "^-u +(.+)", execute: extractData}},
		{"-k", 15, parseInsecure, nil},
		// Body
		{"--data", 10, parseBodyASCII, &extract{re: "--data +(.+)", execute: extractData}},
		{"--data-urlencode", 10, parseBodyURLEncode, &extract{re: "--data-urlencode +(.+)", execute: extractData}},
		{"--data-binary", 10, parseBodyBinary, &extract{re: "--data-binary +(\\${0,1}.+)", execute: extractData}},
		{"--data-ascii", 10, parseBodyASCII, &extract{re: "--data-ascii +(.+)", execute: extractData}},
		{"--data-raw", 10, parseBodyRaw, &extract{re: "--data-raw +(.+)", execute: extractData}},
		//"--"
		{"--header", 10, parseHeader, nil},
		{"--insecure", 15, parseInsecure, nil},
		{"--user-agent", 15, parseUserAgent, &extract{re: "--user-agent +(.+)", execute: extractData}},
		{"--user", 15, parseUser, &extract{re: "--user +(.+)", execute: extractData}},
		{"--connect-timeout", 15, parseTimeout, &extract{re: "--connect-timeout +(.+)", execute: extractData}},
		// 自定义
		// {"--task", 10, parseITask, &extract{re: "--task +(.+)", execute: extractData}},
		// {"--crontab", 10, parseCrontab, &extract{re: "--crontab +(.+)", execute: extractData}},
		// {"--name", 10, parseName, &extract{re: "--name +(.+)", execute: extractData}},
	}

	for _, oe := range oelist {
		optionTrie.Insert(oe)
	}

	// log.Println("support options:", optionTrie.AllWords())
}

// extract 用于提取设置的数据
type extract struct {
	re      string
	execute func(re, soption string) string
}

func (et *extract) Execute(soption string) string {
	return et.execute(et.re, soption)
}

// OptionTrie 设置的前缀树
var optionTrie *hTrie

type optionExecute struct {
	Prefix string

	Priority int

	Parse   func(*CURL, string) error // 执行函数
	Extract *extract                  // 提取的方法结构与参数
}

func (oe *optionExecute) GetWord() string {
	return oe.Prefix + " "
}

func (oe *optionExecute) BuildFunction(curl *CURL, soption string) *parseFunction {
	data := soption
	if oe.Extract != nil {
		data = oe.Extract.Execute(data)
	}
	return &parseFunction{ParamCURL: curl, ParamData: data, ExecuteFunction: oe.Parse, Priority: oe.Priority}
}

func judgeOptions(u *CURL, soption string) *parseFunction {
	word := trieStrWord(soption)
	if ioe := optionTrie.SearchDepth(&word); ioe != nil {
		oe := ioe.(*optionExecute)
		return oe.BuildFunction(u, soption)
	}

	return nil
}

// 提取 被' or " 被包裹 Value值
func extractData(re, soption string) string {
	datas := regexp.MustCompile(re).FindStringSubmatch(soption)
	if len(datas) < 2 {
		log.Printf("error: extractData soption %s", soption)
		return ""
	}
	return strings.Trim(datas[1], "'\"")
}

// func parseName(u *CURL, value string) error {
// 	u.Name = value
// }

// func parseCrontab(u *CURL, value string) error {
// 	u.Crontab = value
// }

// func parseITask(u *CURL, value string) error {
// 	u.iTask = value
// }

func parseTimeout(u *CURL, value string) error {
	timeout, err := strconv.Atoi(value)
	if err != nil {
		log.Println(err)
		return err
	}
	u.Timeout = timeout
	return nil
}

func parseInsecure(u *CURL, soption string) error {
	u.Insecure = true
	return nil
}

func parseUser(u *CURL, soption string) error {
	auth := strings.Split(soption, ":")
	if len(auth) != 2 {
		err := fmt.Errorf("error: parseUser soption = %s", soption)
		log.Println(err)
		return err
	}
	u.Auth = &requests.BasicAuth{User: auth[0], Password: auth[1]}
	return nil
}

func parseUserAgent(u *CURL, value string) error {
	u.Header.Add("User-Agent", value)
	return nil
}

func parseOptX(u *CURL, soption string) error {
	matches := regexp.MustCompile("-X +(.+)").FindStringSubmatch(soption)
	if len(matches) < 2 {
		err := fmt.Errorf("error:parseOptX soption = %s", soption)
		log.Println(err)
		return err
	}
	method := strings.Trim(matches[1], "'")
	u.Method = method
	return nil
}

func parseBodyURLEncode(u *CURL, data string) error {
	if u.Method != "" {
		u.Method = "POST"
	}

	u.ContentType = requests.TypeURLENCODED
	u.Body = bytes.NewBufferString(data)
	return nil
}

func parseBodyRaw(u *CURL, data string) error {
	if u.Method != "" {
		u.Method = "POST"
	}

	u.ContentType = requests.TypeURLENCODED
	u.Body = bytes.NewBufferString(data)
	return nil
}

func parseBodyASCII(u *CURL, data string) error {
	if u.Method != "" {
		u.Method = "POST"
	}

	u.ContentType = requests.TypeURLENCODED

	if data[0] != '@' {
		u.Body = bytes.NewBufferString(data)
	} else {
		f, err := os.Open(data[1:])
		if err != nil {
			log.Println(err)
			return err
		}
		defer f.Close()

		bdata, err := ioutil.ReadAll(f)
		if err != nil {
			log.Println(err)
			return err
		}
		u.Body = bytes.NewBuffer(bdata)
	}

	return nil
}

// 处理@ 并且替/r/n符号
func parseBodyBinary(u *CURL, data string) error {

	if len(data) == 0 {
		err := fmt.Errorf("error: parseBodyBinary data len is 0")
		log.Println(err)
		return err
	}

	if u.Method == "" {
		u.Method = "POST"
	}

	u.ContentType = requests.TypeURLENCODED

	firstchar := data[0]
	switch firstchar {
	case '@':
		f, err := os.Open(data[1:])
		if err != nil {
			log.Println(err)
			return err
		}
		defer f.Close()
		bdata, err := ioutil.ReadAll(f)
		if err != nil {
			log.Println(err)
			return err
		}
		bdata = regexp.MustCompile("\n|\r").ReplaceAll(bdata, []byte(""))
		u.Body = bytes.NewBuffer(bdata)
	case '$':
		data = strings.ReplaceAll(data[2:], `\r\n`, "\r\n")
		u.Body = bytes.NewBufferString(data)
		// boundary parse
		// bindex := strings.Index(data, `\r\n`)
		// boundary := data[4:bindex] // '$--(len=4) build function 已经Trim 末尾'

		// log.Println(fmt.Sprintf(`\r\n--%s--\r\n`, boundary))
		// blastindex := strings.LastIndex(data, fmt.Sprintf(`\r\n--%s--\r\n`, boundary))
		// data = data[bindex+4 : blastindex]
		// strings.Split(data, fmt.Sprintf(`\r\n--%s\r\n`, boundary))
		// log.Println(data)
	default:
		u.Body = bytes.NewBufferString(data)
	}
	return nil
}

func parseHeader(u *CURL, soption string) error {
	result := regexp.MustCompile(`['"]([^:]+): ([^'"]+)['"]`).FindAllStringSubmatch(soption, 1)
	if len(result) == 0 {
		err := fmt.Errorf("error: parseHeader soption = %s", soption)
		log.Println(err)
		return err
	}

	matches := result[0]
	key := matches[1]
	lkey := strings.ToLower(key)
	value := matches[2]

	switch lkey {
	case "cookie":
		u.Cookies = GetRawCookies(value, "")
		u.CookieJar.SetCookies(u.ParsedURL, u.Cookies)
	case "content-type":
		u.ContentType = value
	default:
		u.Header.Add(key, value)
	}
	return nil
}
