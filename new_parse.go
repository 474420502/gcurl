package gcurl

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

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
	// Headers StringSliceValue
}

type CommandParser struct {
	Args         []string
	ArgOptionMap map[string][]*OptionValue
	CurArgKey    *string
	CurSign      *rune
	ArgBuilder   *strings.Builder
	OpenClose    *rune
}

func (opc *CommandParser) compare(other *CommandParser) bool {
	if len(opc.Args) != len(other.Args) {
		return false
	}

	for i, v1 := range opc.Args {
		if v1 != other.Args[i] {
			return false
		}
	}

	for k1, args1 := range opc.ArgOptionMap {
		args2, ok := other.ArgOptionMap[k1]
		if !ok {
			return false
		}

		for i, arg2 := range args2 {
			if arg2 != args1[i] {
				return false
			}
		}
	}

	return true
}

func (opc *CommandParser) ResetArg() {
	opc.ArgBuilder.Reset()
}

func (opc *CommandParser) ResetOpenClose() {
	opc.OpenClose = nil
}

func (opc *CommandParser) WriteRune(r rune) (int, error) {
	return opc.ArgBuilder.WriteRune(r)
}

// collect arg 和 opt
func (opc *CommandParser) Collect() {
	if opc.ArgBuilder.Len() != 0 {
		arg := opc.ArgBuilder.String()
		if opc.CurArgKey != nil {

			optvalue := &OptionValue{}
			if opc.CurSign != nil {
				optvalue.sign = opc.CurSign
				optvalue.expression = arg
			} else {
				optvalue.value = bytes.NewBufferString(arg)
			}

			opc.ArgOptionMap[*opc.CurArgKey] = append(opc.ArgOptionMap[*opc.CurArgKey], optvalue)
			opc.CurArgKey = nil
			opc.CurSign = nil
		} else {
			if arg[0] == '-' {
				if _, ok := opc.ArgOptionMap[arg]; !ok {
					opc.ArgOptionMap[arg] = []*OptionValue{}
				}
				opc.CurArgKey = &arg
			} else {
				opc.Args = append(opc.Args, arg)
			}
		}

		opc.ArgBuilder.Reset()
	}
	opc.OpenClose = nil

}

type OptionValue struct {
	value *bytes.Buffer

	expression string // 字符串值
	sign       *rune  // 标记符号
}

func (optv *OptionValue) check() error {

	if optv.value == nil {
		if optv.sign == nil {
			return fmt.Errorf("Value and Sign is nil")
		}

		if optv.sign != nil {
			sign := *optv.sign
			switch sign {
			case '@':

				f, err := os.Open(optv.expression)
				if err != nil {

					return err
				}
				defer f.Close()
				bdata, err := ioutil.ReadAll(f)
				if err != nil {

					return err
				}
				bdata = regexp.MustCompile("\n|\r").ReplaceAll(bdata, []byte(""))
				optv.value = bytes.NewBuffer(bdata)
				// u.Body = bytes.NewBuffer(bdata)
			case '$':

				optv.value = bytes.NewBufferString(strings.ReplaceAll(optv.expression, `\r\n`, "\r\n"))

			default:

				return fmt.Errorf("unknown sign %b", sign)
			}
		}
	}

	return nil
}

func (optv *OptionValue) Buffer() *bytes.Buffer {
	err := optv.check()
	if err != nil {
		log.Println(err)
	}
	return optv.value
}

func (optv *OptionValue) String() string {
	err := optv.check()
	if err != nil {
		log.Println(err)
	}

	return optv.value.String()
}

func newCommandParser() *CommandParser {
	return &CommandParser{
		ArgOptionMap: make(map[string][]*OptionValue),
		ArgBuilder:   &strings.Builder{},
	}
}

var strQuote1 = '\''
var strQuote2 = '"'

func parseCurlCommandStr(cmdstr string) *CommandParser {
	cmdstrbuf := []rune(cmdstr)
	buflen := len(cmdstrbuf)

	var cur = newCommandParser()

	for i := 0; i < buflen; i++ {
		c := cmdstrbuf[i]
		// log.Println(string(c), c)
		if cur.OpenClose == nil {
			switch c {
			case ' ', '\t', '\n':
				cur.Collect()
				continue
			case '$', '@':
				cur.Collect()
				cur.CurSign = &c
			case strQuote1, strQuote2:
				cur.Collect()
				cur.OpenClose = &c
			case '\\':
				if i+1 < buflen {
					c2 := cmdstrbuf[i+1]
					switch c2 {

					case strQuote1, strQuote2:
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
	cur.Collect()
	return cur
}

func parseCommandArgs(curlStr string) []string {
	cmd := exec.Command("bash", "-c", "go run get_args/main.go "+curlStr)

	data, err := cmd.Output()
	if err != nil {
		log.Println(curlStr)
		panic(err)
	}
	var buf = bytes.NewBuffer(data)
	var args []string
	err = gob.NewDecoder(buf).Decode(&args)
	if err != nil {
		panic(err)
	}

	return args
}

func parseCommandArgsEx(curlStr string) *CommandParser {
	args := parseCommandArgs(curlStr)
	maxsize := len(args)

	result := newCommandParser()

	for i := 0; i < maxsize; i++ {
		arg := args[i]
		if len(arg) > 0 {
			if arg[0] != '-' {
				result.Args = append(result.Args, arg)
			} else {

				for {
					nextIndex := i + 1

					if nextIndex >= maxsize {
						result.ArgOptionMap[arg] = nil
						break
					}

					nextArg := args[nextIndex]
					if len(nextArg) == 0 {
						i++
						continue
					}

					if nextArg[0] == '-' {
						result.ArgOptionMap[arg] = nil
						break
					}

					result.ArgOptionMap[arg] = append(result.ArgOptionMap[arg], &OptionValue{
						value: bytes.NewBufferString(nextArg),
					})
					i = nextIndex
					break
				}
			}
		}
	}

	return result
}
