package options

import (
	"github.com/eclipse/paho.mqtt.golang"
)

// ConnectOptions MQTT 连接设置
type ConnectOptions struct {
	Username string
	Password string
	Brokers  []string
	ClientID string
}

// GetMQTTClientOptions 获取 MQTT 客户端设置
func (co *ConnectOptions) GetMQTTClientOptions() *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	for _, broker := range co.Brokers {
		opts.AddBroker(broker)
	}
	opts.SetClientID(co.ClientID)
	opts.SetUsername(co.Username)
	opts.SetPassword(co.Password)

	return opts
}
