package gcurl

import "strings"

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
	ArgOptionMap map[string][]string
	CurArgKey    *string
	ArgBuilder   *strings.Builder
	OpenClose    *rune
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

func (opc *CommandParser) Collect() {
	if opc.ArgBuilder.Len() != 0 {
		arg := opc.ArgBuilder.String()
		if arg[0] == '-' {
			if _, ok := opc.ArgOptionMap[arg]; !ok {
				opc.ArgOptionMap[arg] = []string{}
			}

			opc.CurArgKey = &arg
		} else {
			if opc.CurArgKey == nil {
				opc.Args = append(opc.Args, arg)
			} else {

				opc.ArgOptionMap[*opc.CurArgKey] = append(opc.ArgOptionMap[*opc.CurArgKey], arg)
				opc.CurArgKey = nil
			}
		}
		opc.ArgBuilder.Reset()
	}
	opc.OpenClose = nil
}

var strQuote1 = '\''
var strQuote2 = '"'

func parseCurlCommandStr(cmdstr string) *CommandParser {
	cmdstrbuf := []rune(cmdstr)
	buflen := len(cmdstrbuf)

	var cur = &CommandParser{
		ArgOptionMap: make(map[string][]string),
		ArgBuilder:   &strings.Builder{},
	}

	for i := 0; i < buflen; i++ {
		c := cmdstrbuf[i]
		// log.Println(string(c), c)
		if cur.OpenClose == nil {
			switch c {
			case ' ', '\t', '\n':
				cur.Collect()
				continue
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
