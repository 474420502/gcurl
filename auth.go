package gcurl

import "fmt"

// AuthType 定义认证类型
type AuthType int

const (
	AuthBasic AuthType = iota
	AuthDigest
	AuthBearer
	AuthNTLM
)

// String 返回认证类型的字符串表示
func (at AuthType) String() string {
	switch at {
	case AuthBasic:
		return "Basic"
	case AuthDigest:
		return "Digest"
	case AuthBearer:
		return "Bearer"
	case AuthNTLM:
		return "NTLM"
	default:
		return "Unknown"
	}
}

// Authentication 统一的认证结构
type Authentication struct {
	Type     AuthType
	Username string
	Password string
	Token    string // 用于Bearer认证

	// Digest认证的特殊字段
	Realm     string
	Nonce     string
	URI       string
	Response  string
	Algorithm string
	QOP       string // quality of protection
	NC        string // nonce count
	CNonce    string // client nonce

	// 配置选项
	Digest bool // 是否强制使用Digest认证
}

// NewBasicAuth 创建基本认证
func NewBasicAuth(username, password string) *Authentication {
	return &Authentication{
		Type:     AuthBasic,
		Username: username,
		Password: password,
	}
}

// NewDigestAuth 创建摘要认证
func NewDigestAuth(username, password string) *Authentication {
	return &Authentication{
		Type:     AuthDigest,
		Username: username,
		Password: password,
		Digest:   true,
	}
}

// NewBearerAuth 创建Bearer令牌认证
func NewBearerAuth(token string) *Authentication {
	return &Authentication{
		Type:  AuthBearer,
		Token: token,
	}
}

// IsValid 检查认证信息是否有效
func (auth *Authentication) IsValid() bool {
	if auth == nil {
		return false
	}

	switch auth.Type {
	case AuthBasic, AuthDigest:
		return auth.Username != "" && auth.Password != ""
	case AuthBearer:
		return auth.Token != ""
	case AuthNTLM:
		return auth.Username != "" && auth.Password != ""
	default:
		return false
	}
}

// GetAuthHeader 生成认证头部值
func (auth *Authentication) GetAuthHeader() string {
	if !auth.IsValid() {
		return ""
	}

	switch auth.Type {
	case AuthBasic:
		// Basic认证在requests库中处理
		return ""
	case AuthBearer:
		return "Bearer " + auth.Token
	case AuthDigest:
		// Digest认证需要服务器挑战后构建
		// 这里返回空，实际处理在requests库中
		return ""
	default:
		return ""
	}
}

// String 返回认证信息的字符串表示（用于调试）
func (auth *Authentication) String() string {
	if auth == nil {
		return "None"
	}

	switch auth.Type {
	case AuthBasic:
		return fmt.Sprintf("Basic (%s:***)", auth.Username)
	case AuthDigest:
		return fmt.Sprintf("Digest (%s:***)", auth.Username)
	case AuthBearer:
		return "Bearer (***)"
	case AuthNTLM:
		return fmt.Sprintf("NTLM (%s:***)", auth.Username)
	default:
		return "Unknown"
	}
}
