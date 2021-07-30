package constants

type AuthType string

const (
	AuthTypeSign   AuthType = "Signature"        // 签名鉴权模式
	AuthTypeToken  AuthType = "Token"            // Token 鉴权模式
	AuthTypeDevice AuthType = "DeviceCredential" // 一机一密鉴权模
)

type AuthCredential struct {
	Type  string
	Token string
}
