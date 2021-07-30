package constants

import (
	"github.com/eclipse/paho.mqtt.golang"
)

type ConnectOptions struct {
	Username string
	Password string
	Brokers  []string
	ClientID string
}

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
