package aliyun_mqtt_go

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/zhangsq-ax/aliyun_mqtt_go/constants"
)

type ConnectHelper struct {
	AuthType   constants.AuthType        // 鉴权类型，非必须设置，不设置时需要在 GetClient() 或 GetConnectedClient() 方法中指定
	Protocol   constants.ConnectProtocol // 连接协议
	InstanceID string                    // 服务实例标识
	Endpoints  []string                  // 服务接入点
	Port       int                       // 服务接入点端口，非必须设置，缺省设置时将使用 Protocol 对应的默认端口
	Username   string                    // 用户名，签名鉴权和 Token 鉴权模式下为管理员分配的 AccessKeyId，一机一密鉴权模式下使用鉴权服务分发的 DeviceAccessKeyId
	Password   string                    // 密码，签名鉴权模式下使用管理分发的 AccessKeyId, Token 鉴权模式下使用鉴权服务分发的 Token, 一机一密鉴权模式下使用鉴权服务分发的 DeviceAccessKeySecret
}

func (helper *ConnectHelper) generateConnectOptions(opts *constants.ClientOptions) *constants.ConnectOptions {
	port := helper.getPort()
	clientId := fmt.Sprintf("%s@@@%s", opts.GroupID, opts.ClientID)
	var brokers []string
	for _, endpoint := range helper.Endpoints {
		brokers = append(brokers, fmt.Sprintf("%s://%s:%d", helper.Protocol, endpoint, port))
	}
	username := fmt.Sprintf("%s|%s|%s", helper.AuthType, helper.Username, helper.InstanceID)
	password := helper.getPassword(opts.AuthType, clientId)

	return &constants.ConnectOptions{
		Username: username,
		Password: password,
		Brokers:  brokers,
		ClientID: clientId,
	}
}

// getPort 获取连接 MQTT 服务的端口
func (helper *ConnectHelper) getPort() int {
	if helper.Port == 0 {
		helper.Port = constants.ConnectPort[helper.Protocol]
	}
	return helper.Port
}

// getPassword 根据设置的鉴权模式获取相应的 password
func (helper *ConnectHelper) getPassword(authType *constants.AuthType, clientId string) string {
	if authType == nil {
		authType = &helper.AuthType
	}
	var password string
	switch *authType {
	case constants.AuthTypeSign, constants.AuthTypeDevice: // 签名鉴权或一机一密鉴权
		mac := hmac.New(sha1.New, []byte(helper.Password))
		mac.Write([]byte(clientId))
		password = base64.StdEncoding.EncodeToString(mac.Sum(nil))
	case constants.AuthTypeToken: // Token 鉴权
		password = helper.Password
	}

	return password
}

// GetClient 获取 MQTT 客户端，客户端处理未连接状态，需要手动连接
func (helper *ConnectHelper) GetClient(opts *constants.ClientOptions) mqtt.Client {
	connectOpts := helper.generateConnectOptions(opts)
	mqOpts := connectOpts.GetMQTTClientOptions()

	mqOpts.OnConnect = opts.OnConnect
	if opts.OnConnect == nil {
		mqOpts.OnConnect = onConnectDefault
	}

	mqOpts.OnConnectionLost = opts.OnConnectionLost
	if opts.OnConnectionLost == nil {
		mqOpts.OnConnectionLost = onConnectionLostDefault
	}

	return mqtt.NewClient(mqOpts)
}

// GetConnectedClient 获取已连接的 MQTT 客户端
func (helper *ConnectHelper) GetConnectedClient(opts *constants.ClientOptions) (mqtt.Client, error) {
	client := helper.GetClient(opts)
	token := client.Connect()
	token.Wait()
	err := token.Error()
	if err != nil {
		return nil, err
	}

	if client.IsConnected() {
		return client, nil
	}
	return nil, errors.New("failed to connect")
}

func onConnectDefault(client mqtt.Client) {
	fmt.Println("---------- Connect to server success ----------")
}

func onConnectionLostDefault(client mqtt.Client, err error) {
	fmt.Printf("-----X----- Lost connection with server: %v\n", err)
}
