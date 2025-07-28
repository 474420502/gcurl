package gcurl

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// isValidURL 实现更严格的URL验证
func isValidURL(urlStr string) bool {
	if urlStr == "" || len(strings.TrimSpace(urlStr)) == 0 {
		return false
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return false
	}

	// 检查是否有scheme
	if parsedURL.Scheme == "" {
		return false
	}

	// 对于http/https，检查是否有host
	if (parsedURL.Scheme == "http" || parsedURL.Scheme == "https") && parsedURL.Host == "" {
		return false
	}

	return true
}

// in parse_options.go or a new parser.go

func buildFromArgs(args []string) (*CURL, error) {
	curl := New() // New() 初始化一个空的 CURL 对象

	i := 0
	for i < len(args) {
		arg := args[i]

		// 第一个不以'-'开头的参数通常是 "curl" 本身或 URL
		if !strings.HasPrefix(arg, "-") {
			if arg == "curl" {
				i++
				continue
			}
			// 假定这是URL，只处理一次
			if curl.ParsedURL == nil {
				// 使用增强的URL验证
				if !isValidURL(arg) {
					return nil, fmt.Errorf("invalid or malformed URL: %s", arg)
				}

				purl, err := url.Parse(arg)
				if err != nil {
					return nil, fmt.Errorf("invalid URL: %s", arg)
				}
				curl.ParsedURL = purl
			} else {
				return nil, fmt.Errorf("multiple URLs provided or misplaced argument: %s", arg)
			}
			i++
			continue
		}

		// 在注册表中查找选项
		spec, found := optionRegistry[arg]
		if !found {
			// 检查是否在跳过列表中
			skipType := checkInSkipList(arg)
			if skipType != ST_NotSkipType {
				// 这是一个已知但跳过的选项，根据类型决定跳过多少个参数
				if skipType == ST_OnlyOption {
					// 只跳过选项本身
					i++
				} else if skipType == ST_WithValue {
					// 跳过选项和它的值
					i += 2
				}
				continue
			}
			return nil, fmt.Errorf("unsupported or unknown option: %s", arg)
		}

		// 检查是否有足够的参数
		if i+1+spec.NumArgs > len(args) {
			return nil, fmt.Errorf("option %s requires %d argument(s), but not enough provided", arg, spec.NumArgs)
		}

		// 提取参数并调用处理器
		handlerArgs := args[i+1 : i+1+spec.NumArgs]
		if err := spec.Handler(curl, handlerArgs...); err != nil {
			// 标准化错误返回！
			return nil, fmt.Errorf("error processing option %s: %w", arg, err)
		}

		// 更新循环索引，跳过已消费的选项和参数
		i += 1 + spec.NumArgs
	}

	// 1. 检查URL是否存在
	if curl.ParsedURL == nil {
		return nil, errors.New("no URL specified in command")
	}

	// 2. 在确认URL存在后，安全地将所有解析到的Cookies添加到CookieJar
	if len(curl.Cookies) > 0 {
		curl.CookieJar.SetCookies(curl.ParsedURL, curl.Cookies)
	}

	// 设置默认请求方法
	if curl.Method == "" {
		if curl.Body != nil && curl.Body.Len() > 0 {
			curl.Method = "POST"
		} else {
			curl.Method = "GET"
		}
	}

	return curl, nil
}
