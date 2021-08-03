package aliyun_mqtt_go

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/zhangsq-ax/aliyun_mqtt_go/constants"
	"github.com/zhangsq-ax/aliyun_mqtt_go/options"
)

// ConnectHelper MQTT 客户端连接助手
type ConnectHelper struct {
	options *options.MQTTClientOptions
}

// 生成 MQTT 客户端连接信息
func (helper *ConnectHelper) generateConnectOptions() (connOpts *options.ConnectOptions, expiredTime int64, err error) {
	var password string
	opts := helper.options
	port := helper.getPort()
	clientId := fmt.Sprintf("%s@@@%s", opts.GroupID, opts.ClientID)
	var brokers []string
	for _, endpoint := range opts.Endpoints {
		brokers = append(brokers, fmt.Sprintf("%s://%s:%d", opts.Protocol, endpoint, port))
	}
	username := fmt.Sprintf("%s|%s|%s", opts.AuthType, opts.Username, opts.InstanceID)
	password, expiredTime, err = helper.getPassword(clientId)
	if err != nil {
		return
	}

	connOpts = &options.ConnectOptions{
		Username: username,
		Password: password,
		Brokers:  brokers,
		ClientID: clientId,
	}

	return
}

// getPort 获取连接 MQTT 服务的端口
func (helper *ConnectHelper) getPort() int {
	if helper.options.Port == 0 {
		helper.options.Port = constants.ConnectPort[helper.options.Protocol]
	}
	return helper.options.Port
}

// getPassword 根据设置的鉴权模式获取相应的 password
func (helper *ConnectHelper) getPassword(clientId string) (password string, expiredTime int64, err error) {
	opts := helper.options
	switch opts.AuthType {
	case constants.AuthTypeSign, constants.AuthTypeDevice: // 签名鉴权或一机一密鉴权
		mac := hmac.New(sha1.New, []byte(opts.Password))
		mac.Write([]byte(clientId))
		password = base64.StdEncoding.EncodeToString(mac.Sum(nil))
	case constants.AuthTypeToken: // Token 鉴权
		if opts.PasswordGetter != nil {
			// 如果可能优先通过 PasswordGetter() 方法获取
			password, expiredTime, err = opts.PasswordGetter(opts.ClientID)
		}
		if opts.Password != "" {
			password = opts.Password
		}
	}

	return
}

// GetClient 获取 MQTT 客户端，客户端处理未连接状态，需要手动连接
func (helper *ConnectHelper) GetClient() (client *mqtt.Client, expiredTime int64, err error) {
	var connOpts *options.ConnectOptions
	opts := helper.options
	connOpts, expiredTime, err = helper.generateConnectOptions()
	if err != nil {
		return
	}
	mqOpts := connOpts.GetMQTTClientOptions()

	mqOpts.OnConnect = opts.OnConnect
	if opts.OnConnect == nil {
		mqOpts.OnConnect = onConnectDefault
	}

	mqOpts.OnConnectionLost = opts.OnConnectionLost
	if opts.OnConnectionLost == nil {
		mqOpts.OnConnectionLost = onConnectionLostDefault
	}

	c := mqtt.NewClient(mqOpts)
	client = &c

	return
}

// GetClient 获取 MQTT 客户端对象
func GetClient(opts *options.MQTTClientOptions) (client *mqtt.Client, expiredTime int64, err error) {
	helper := NewConnectHelperFromClientOptions(opts)
	client, expiredTime, err = helper.GetClient()
	if err != nil {
		return
	}
	c := *client

	token := c.Connect()
	token.Wait()
	err = token.Error()
	if err != nil {
		client = nil
		expiredTime = 0
		return
	}

	if !c.IsConnected() {
		client = nil
		expiredTime = 0
		err = errors.New("failed to connect MQTT service")
	}

	return
}

// NewConnectHelperFromClientOptions 根据 ClientOptions 创建 ConnectHelper 实例
func NewConnectHelperFromClientOptions(opts *options.MQTTClientOptions) *ConnectHelper {
	return &ConnectHelper{
		options: opts,
	}
}

func onConnectDefault(client mqtt.Client) {
	fmt.Println("---------- Connect to server success ----------")
}

func onConnectionLostDefault(client mqtt.Client, err error) {
	fmt.Printf("-----X----- Lost connection with server: %v\n", err)
}
