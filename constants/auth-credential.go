package constants

// AuthType 鉴权类型
type AuthType string

const (
	AuthTypeSign   AuthType = "Signature"        // 签名鉴权模式
	AuthTypeToken  AuthType = "Token"            // Token 鉴权模式
	AuthTypeDevice AuthType = "DeviceCredential" // 一机一密鉴权模
)

// AuthCredential 鉴权凭证
type AuthCredential struct {
	Type  string
	Token string
}
