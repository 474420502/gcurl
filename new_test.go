package gcurl

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"testing"
)

func TestNewT(t *testing.T) {
	main1()
}

type CurlOptions struct {
	Data       string
	FailFast   bool
	Help       string
	Include    bool
	Output     string
	RemoteName bool
	Silent     bool
	UploadFile string
	User       string
	UserAgent  string

	// Head    bool
	Headers StringSliceValue
}

type StringSliceValue []string

func newStringSliceValue() *StringSliceValue {
	return new(StringSliceValue)
}

func (ssv *StringSliceValue) Set(value string) error {
	*ssv = append(*ssv, value)
	return nil
}

func (ssv *StringSliceValue) String() string {
	return fmt.Sprintf("%v", *ssv)
}

func parseCurlOptions(cmdStr string) (*CurlOptions, error) {

	args := parseCurlCommandStr(cmdStr)
	opts := &CurlOptions{}

	fs := flag.NewFlagSet("curl", flag.ContinueOnError)
	fs.StringVar(&opts.Data, "d", "", "HTTP POST data")
	fs.StringVar(&opts.Data, "data", "", "HTTP POST data")

	fs.BoolVar(&opts.RemoteName, "O", false, "Write output to a file named as the remote file")
	fs.BoolVar(&opts.RemoteName, "remote-name", false, "Write output to a file named as the remote file")
	fs.BoolVar(&opts.Silent, "s", false, "Silent mode")
	fs.BoolVar(&opts.Silent, "silent", false, "Silent mode")
	fs.StringVar(&opts.UploadFile, "T", "", "Transfer local FILE to destination")
	fs.StringVar(&opts.UploadFile, "upload-file", "", "Transfer local FILE to destination")
	fs.StringVar(&opts.User, "u", "", "Server user and password")
	fs.StringVar(&opts.User, "user", "", "Server user and password")
	fs.StringVar(&opts.UserAgent, "A", "", "Send User-Agent <name> to server")
	fs.StringVar(&opts.UserAgent, "user-agent", "", "Send User-Agent <name> to server")
	// fs.Var(&opts.Headers, "H", "Pass custom header(s) to server")
	// fs.StringVar(&opts.Head, "header", "Pass custom header(s) to server")

	err := fs.Parse(args[1:])
	if err != nil {
		return nil, err
	}

	return opts, nil
}

type ArgCollect struct {
	Args      []string
	Arg       *strings.Builder
	OpenClose *rune
}

func (opc *ArgCollect) ResetArg() {
	opc.Arg.Reset()
}

func (opc *ArgCollect) ResetOpenClose() {
	opc.OpenClose = nil
}

func (opc *ArgCollect) WriteRune(r rune) (int, error) {
	return opc.Arg.WriteRune(r)
}

func (opc *ArgCollect) Collect() {
	if opc.Arg.Len() != 0 {
		opc.Args = append(opc.Args, opc.Arg.String())
		opc.Arg.Reset()
	}
	opc.OpenClose = nil
}

var strQuote1 = '\''
var strQuote2 = '"'

func parseCurlCommandStr(cmdstr string) []string {
	cmdstrbuf := []rune(cmdstr)
	buflen := len(cmdstrbuf)

	var cur = &ArgCollect{
		Arg: &strings.Builder{},
	}

	for i := 0; i < buflen; i++ {
		c := cmdstrbuf[i]
		// log.Println(string(c), c)
		if cur.OpenClose == nil {
			switch c {
			case ' ':
				cur.Collect()
				continue
			case strQuote1, strQuote2:
				cur.Collect()
				cur.OpenClose = &c
			case '\\':
				if i+1 < buflen {
					c2 := cmdstrbuf[i+1]
					switch c2 {

					default:
						cur.WriteRune(cmdstrbuf[i+1])
					}
					i++
				}
			// 直接跳过

			default:
				cur.WriteRune(c)
			}
		} else {

			if *cur.OpenClose == c {
				cur.Collect()
				continue
			}

			switch c {
			case '\\':
				if i+1 < buflen {
					c2 := cmdstrbuf[i+1]
					switch c2 {

					default:
						cur.WriteRune(cmdstrbuf[i+1])
					}
					i++
				}
				// 直接跳过
			default:
				cur.WriteRune(c)
			}

		}

	}

	return cur.Args
}

func main1() {

	cmdStr := `curl 'https://www.xxxxx.com/api-hk/heartbeat' \
	-d '{\"name\": \"Alice\"}' \
	-H 'accept: application/json, text/plain, */*' \
	-H 'accept-language: zh-CN,zh;q=0.9,en;q=0.8' \
	-H 'origin: https://www.xxxxx.com' \
	-H 'referer: https://www.xxxxx.com/' \
	-H 'sec-ch-ua: "Google Chrome";v="123", "Not:A-Brand";v="8", "Chromium";v="123"' \
	-H 'sec-ch-ua-mobile: ?0' \
	-H 'sec-ch-ua-platform: "Windows"' \
	-H 'sec-fetch-dest: empty' \
	-H 'sec-fetch-mode: cors' \
	-H 'sec-fetch-site: cross-site' \
	-H 'user-agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36'`

	opts, err := parseCurlOptions(cmdStr)
	log.Println(opts.Data)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Curl Options:")
	fmt.Printf("Data: %s\n", opts.Data)
	fmt.Printf("FailFast: %t\n", opts.FailFast)
	fmt.Printf("Help: %s\n", opts.Help)
	fmt.Printf("Include: %t\n", opts.Include)
	fmt.Printf("Output: %s\n", opts.Output)
	fmt.Printf("RemoteName: %t\n", opts.RemoteName)
	fmt.Printf("Silent: %t\n", opts.Silent)
	fmt.Printf("UploadFile: %s\n", opts.UploadFile)
	fmt.Printf("User: %s\n", opts.User)
	fmt.Printf("UserAgent: %s\n", opts.UserAgent)
}
