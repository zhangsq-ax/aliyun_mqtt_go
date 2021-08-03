package options

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/zhangsq-ax/aliyun_mqtt_go/constants"
)

// MQTTClientOptions 创建 MQTTClient 的设置
type MQTTClientOptions struct {
	AuthType         constants.AuthType                                                // 鉴权类型，必填
	Protocol         constants.ConnectProtocol                                         // 连接协议，必填
	InstanceID       string                                                            // 服务实例标识，必填
	Endpoints        []string                                                          // 服务接入点，必填
	Port             int                                                               // 服务接入点端口，非必填，缺省设置时将使用 Protocol 对应的默认端口
	Username         string                                                            // 用户名，必填，签名鉴权和 Token 鉴权模式下为管理员分配的 AccessKeyId，一机一密鉴权模式下使用鉴权服务分发的 DeviceAccessKeyId
	Password         string                                                            // 密码，非必填，签名鉴权模式下使用管理分发的 AccessKeyId, Token 鉴权模式下使用鉴权服务分发的 Token, 一机一密鉴权模式下使用鉴权服务分发的 DeviceAccessKeySecret
	PasswordGetter   func(clientId string) (token string, expireTime int64, err error) // 动态获取密码方法，非必填，只在 Token 鉴权模式下设置有效，优先级高于 Password。token 鉴权模式下 PasswordGetter 和 Password 必填设置一个
	GroupID          string                                                            // 客户端分组标识，必填
	ClientID         string                                                            // 客户端标识，必填
	OnConnect        mqtt.OnConnectHandler                                             // 连接后的回调方法，非必填
	OnConnectionLost mqtt.ConnectionLostHandler                                        // 连接断开后的回调方法，非必填
}
