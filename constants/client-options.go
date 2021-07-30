package constants

import mqtt "github.com/eclipse/paho.mqtt.golang"

type ClientOptions struct {
	AuthType         *AuthType                  // 鉴权模式，非必须设置，缺省时使用 ConnectHelper 的设置
	GroupID          string                     // 客户端分组标识
	ClientID         string                     // 客户端标识
	OnConnect        mqtt.OnConnectHandler      // 连接后的回调方法
	OnConnectionLost mqtt.ConnectionLostHandler // 连接断开后的回调方法，可以在这里重新连接
}
