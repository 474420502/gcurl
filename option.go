package gcurl

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/474420502/requests"
)

func init() {
	optionTrie = newTrie()
	// oelist := []*optionExecute{
	// 	{"-H", 10, parseHeader, nil},
	// 	{"-X", 10, parseMethod, nil},
	// 	{"-A", 15, parseUserAgent, &extract{re: "^-A +(.+)", execute: extractData}},
	// 	{"-d", 10, parseBodyASCII, &extract{re: "^-d +(.+)", execute: extractData}},
	// 	{"-u", 15, parseUser, &extract{re: "^-u +(.+)", execute: extractData}},
	// 	{"-k", 15, parseInsecure, nil},

	// 	{"--data", 10, parseBodyASCII, &extract{re: "--data +(.+)", execute: extractData}},
	// 	{"--data-urlencode", 10, parseBodyURLEncode, &extract{re: "--data-urlencode +(.+)", execute: extractData}},
	// 	{"--data-binary", 10, parseBodyBinary, &extract{re: "--data-binary +(\\{0,1}.+)", execute: extractData}},
	// 	{"--data-ascii", 10, parseBodyASCII, &extract{re: "--data-ascii +(.+)", execute: extractData}},
	// 	{"--data-raw", 10, parseBodyRaw, &extract{re: "--data-raw +(.+)", execute: extractData}},
	// 	//"--"
	// 	{"--header", 10, parseHeader, nil},
	// 	{"--insecure", 15, parseInsecure, nil},
	// 	{"--user-agent", 15, parseUserAgent, &extract{re: "--user-agent +(.+)", execute: extractData}},
	// 	{"--user", 15, parseUser, &extract{re: "--user +(.+)", execute: extractData}},
	// 	{"--connect-timeout", 15, parseTimeout, &extract{re: "--connect-timeout +(.+)", execute: extractData}},
	// 	// 自定义
	// 	// {"--task", 10, parseITask, &extract{re: "--task +(.+)", execute: extractData}},
	// 	// {"--crontab", 10, parseCrontab, &extract{re: "--crontab +(.+)", execute: extractData}},
	// 	// {"--name", 10, parseName, &extract{re: "--name +(.+)", execute: extractData}},
	// }

	// for _, oe := range oelist {
	// 	optionTrie.Insert(oe)
	// }

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

func judgeOptions(u *CURL, soption *OptionValue) *parseFunction {
	word := trieStrWord(soption.String())
	if ioe := optionTrie.SearchDepth(&word); ioe != nil {
		oe := ioe.(*optionExecute)
		return oe.BuildFunction(u, soption.String())
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

func parseTimeout(u *CURL, value *OptionValue) error {
	timeout, err := strconv.Atoi(value.String())
	if err != nil {
		log.Println(err)
		return err
	}
	u.Timeout = timeout
	return nil
}

func parseInsecure(u *CURL, soption *OptionValue) error {
	u.Insecure = true
	return nil
}

func parseUser(u *CURL, soption *OptionValue) error {
	auth := strings.Split(soption.String(), ":")
	if len(auth) != 2 {
		err := fmt.Errorf("error: parseUser soption = %s", soption)
		log.Println(err)
		return err
	}
	u.Auth = &requests.BasicAuth{User: auth[0], Password: auth[1]}
	return nil
}

func parseUserAgent(u *CURL, value *OptionValue) error {
	u.Header.Add("User-Agent", value.String())
	return nil
}

func parseMethod(u *CURL, soption *OptionValue) error {
	u.Method = soption.String()
	return nil
}

func parseBodyURLEncode(u *CURL, data *OptionValue) error {
	if u.Method != "" {
		u.Method = "POST"
	}

	u.ContentType = requests.TypeURLENCODED
	u.Body = data.Buffer()
	return nil
}

func parseBodyRaw(u *CURL, data *OptionValue) error {
	if u.Method != "" {
		u.Method = "POST"
	}

	u.ContentType = requests.TypeURLENCODED
	u.Body = data.Buffer()
	return nil
}

func parseBodyASCII(u *CURL, data *OptionValue) error {
	if u.Method != "" {
		u.Method = "POST"
	}

	u.ContentType = requests.TypeURLENCODED

	u.Body = data.Buffer()

	return nil
}

// 处理@ 并且替/r/n符号
func parseBodyBinary(u *CURL, data *OptionValue) error {

	if u.Method == "" {
		u.Method = "POST"
	}

	u.ContentType = requests.TypeURLENCODED
	u.Body = data.Buffer()
	// firstchar := data[0]
	// switch firstchar {
	// case '@':
	// 	f, err := os.Open(data[1:])
	// 	if err != nil {
	// 		log.Println(err)
	// 		return err
	// 	}
	// 	defer f.Close()
	// 	bdata, err := ioutil.ReadAll(f)
	// 	if err != nil {
	// 		log.Println(err)
	// 		return err
	// 	}
	// 	bdata = regexp.MustCompile("\n|\r").ReplaceAll(bdata, []byte(""))
	// 	u.Body = bytes.NewBuffer(bdata)
	// case '$':
	// 	data = strings.ReplaceAll(data[2:], `\r\n`, "\r\n")
	// 	u.Body = bytes.NewBufferString(data)
	// 	// boundary parse
	// 	// bindex := strings.Index(data, `\r\n`)
	// 	// boundary := data[4:bindex] // '$--(len=4) build function 已经Trim 末尾'

	// 	// log.Println(fmt.Sprintf(`\r\n--%s--\r\n`, boundary))
	// 	// blastindex := strings.LastIndex(data, fmt.Sprintf(`\r\n--%s--\r\n`, boundary))
	// 	// data = data[bindex+4 : blastindex]
	// 	// strings.Split(data, fmt.Sprintf(`\r\n--%s\r\n`, boundary))
	// 	// log.Println(data)
	// default:
	// 	u.Body = bytes.NewBufferString(data)
	// }
	return nil
}

func parseHTTPHeaderKeyValue(soption string) (hkey string, hvalue string, err error) {
	var i = 0
	charlen := len(soption)

	var keyBuilder = &strings.Builder{}
	var valueBuilder = &strings.Builder{}
	// get key string
	for ; i < charlen; i++ {
		c := soption[i]
		if c == ':' {
			i++
			break
		}
		err := keyBuilder.WriteByte(c)
		if err != nil {
			return "", "", err
		}
	}

	// skip space
	for ; i < charlen; i++ {
		c := soption[i]
		if c != ' ' {
			break
		}
	}

	var stringMark byte
	c := soption[i]
	if c == '\'' || c == '"' {
		stringMark = c
		if soption[charlen-1] != stringMark {
			return "", "", fmt.Errorf("%b != %b", soption[charlen-1], stringMark)
		}
		i++
		charlen--
	}

	for ; i < charlen; i++ {
		c := soption[i]
		err := valueBuilder.WriteByte(c)
		if err != nil {
			return "", "", err
		}
	}

	return keyBuilder.String(), valueBuilder.String(), nil
}

func parseHeader(u *CURL, soption *OptionValue) error {

	key, value, err := parseHTTPHeaderKeyValue(soption.String())
	if err != nil {
		return err
	}

	lkey := strings.ToLower(key)
	u.Header.Add(key, value)

	switch lkey {
	case "cookie":
		u.Cookies = GetRawCookies(value, "")
		u.CookieJar.SetCookies(u.ParsedURL, u.Cookies)
	case "content-type":
		u.ContentType = value
	}
	return nil
}
