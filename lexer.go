package gcurl

import (
	"errors"
	"strings"
	"unicode"
)

type LexerState int

const (
	StateDefault LexerState = iota
	StateInArg
	StateInSingleQuotes
	StateInDoubleQuotes
	StateInAnsiCQuotes // 新增状态用于处理 $'...' 格式
)

type Lexer struct {
	input    []rune
	pos      int
	state    LexerState
	builder  strings.Builder
	Tokens   []string
	inQuotes bool // 标记当前token是否来自引号（包括空引号）
}

// NewLexer now includes the pre-processing step.
func NewLexer(input string) *Lexer {
	// 【关键步骤1】在这里进行预处理
	processedInput := preProcessCurlString(input)
	return &Lexer{
		input:  []rune(processedInput),
		Tokens: make([]string, 0),
		state:  StateDefault,
	}
}

// preProcessCurlString is a helper to clean up line continuations.
func preProcessCurlString(input string) string {
	// Order matters: replace CRLF first, then LF.
	processed := strings.ReplaceAll(input, "\\\r\n", " ")
	processed = strings.ReplaceAll(processed, "\\\n", " ")
	return processed
}

func (l *Lexer) Parse() error {
	for l.pos < len(l.input) {
		r := l.input[l.pos]

		switch l.state {
		case StateDefault:
			// 【关键步骤2】任何空白字符都只是被跳过
			if !unicode.IsSpace(r) {
				l.state = StateInArg
				// 不回退，因为当前字符需要被处理
				// 通过 continue 来重新进入循环处理当前字符
				continue
			}

		case StateInArg:
			switch r {
			case '\'':
				// 检查是否是 $' 格式（ANSI-C 引用）
				if l.builder.Len() > 0 && l.builder.String()[l.builder.Len()-1] == '$' {
					// 移除最后的 '$' 字符，因为我们要进入 ANSI-C 引用状态
					content := l.builder.String()
					l.builder.Reset()
					l.builder.WriteString(content[:len(content)-1])
					l.state = StateInAnsiCQuotes
					l.inQuotes = true
				} else {
					l.state = StateInSingleQuotes
					l.inQuotes = true
				}
			case '"':
				l.state = StateInDoubleQuotes
				l.inQuotes = true
			case '\\': // 处理转义 (现在只处理字面转义，不是行连续)
				l.pos++
				if l.pos < len(l.input) {
					l.builder.WriteRune(l.input[l.pos])
				}
			default:
				// 【关键步骤3】任何空白字符都会结束当前参数
				if unicode.IsSpace(r) {
					l.finalizeToken()
					l.state = StateDefault
				} else {
					l.builder.WriteRune(r)
				}
			}

		case StateInSingleQuotes:
			if r == '\'' {
				l.state = StateInArg
			} else {
				// 单引号内，所有字符原样保留，包括 \n
				l.builder.WriteRune(r)
			}

		case StateInDoubleQuotes:
			switch r {
			case '"':
				l.state = StateInArg
			case '\\':
				l.pos++
				if l.pos < len(l.input) {
					l.builder.WriteRune(l.input[l.pos])
				}
			default:
				// 双引号内，除了转义的外，所有字符原样保留，包括 \n
				l.builder.WriteRune(r)
			}

		case StateInAnsiCQuotes:
			if r == '\'' {
				l.state = StateInArg
			} else if r == '\\' && l.pos+1 < len(l.input) {
				// 处理 ANSI-C 转义序列
				l.pos++
				nextChar := l.input[l.pos]
				switch nextChar {
				case 'r':
					l.builder.WriteByte('\r')
				case 'n':
					l.builder.WriteByte('\n')
				case 't':
					l.builder.WriteByte('\t')
				case '\'':
					l.builder.WriteByte('\'')
				case '\\':
					l.builder.WriteByte('\\')
				default:
					// 对于未知的转义序列，保留原样
					l.builder.WriteByte('\\')
					l.builder.WriteRune(nextChar)
				}
			} else {
				// ANSI-C 引号内的普通字符
				l.builder.WriteRune(r)
			}
		}
		l.pos++
	}

	if l.builder.Len() > 0 {
		l.finalizeToken()
	}

	if l.state == StateInSingleQuotes || l.state == StateInDoubleQuotes || l.state == StateInAnsiCQuotes {
		return errors.New("command has unclosed quotes")
	}

	return nil
}

func (l *Lexer) finalizeToken() {
	// 即使是空字符串，如果来自引号也要添加
	if l.builder.Len() > 0 || l.inQuotes {
		l.Tokens = append(l.Tokens, l.builder.String())
		l.builder.Reset()
		l.inQuotes = false // 重置引号标记
	}
}
