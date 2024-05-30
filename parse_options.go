package gcurl

import (
	"fmt"
	"log"
	"net/url"
)

func ParseOptions(cp *CommandParser) (curl *CURL, err error) {
	curl = New()
	if len(cp.Args) == 0 {
		return nil, fmt.Errorf("args len is 0")
	}

	var urlStr string = cp.Args[0]
	if urlStr == "curl" {
		if len(cp.Args) < 2 {
			return nil, fmt.Errorf("urlstr is not exists")
		}
		urlStr = cp.Args[1]
	}

	purl, err := url.Parse(urlStr)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	curl.ParsedURL = purl

	for key, values := range cp.ArgOptionMap {
		switch key {
		case "-H", "--header":
			for _, v := range values {
				parseHeader(curl, v)
			}
		case "-X":
			for _, v := range values {
				parseMethod(curl, v)
			}
		case "-A", "--user-agent":
			for _, v := range values {
				parseUserAgent(curl, v)
			}
		case "-u", "--user":
			for _, v := range values {
				parseUser(curl, v)
			}
		case "-k", "--insecure":
			for _, v := range values {
				parseInsecure(curl, v)
			}
		case "-d", "--data", "--data-ascii":
			for _, v := range values {
				parseBodyASCII(curl, v)
			}

		case "--data-urlencode":
			for _, v := range values {
				parseBodyURLEncode(curl, v)
			}

		case "--data-binary":
			for _, v := range values {
				parseBodyBinary(curl, v)
			}

		case "--data-raw":
			for _, v := range values {
				parseBodyRaw(curl, v)
			}

		case "--connect-timeout":
			for _, v := range values {
				parseTimeout(curl, v)
			}
		}
	}

	if curl.Method == "" {
		curl.Method = "GET"
	}

	return curl, nil
}
