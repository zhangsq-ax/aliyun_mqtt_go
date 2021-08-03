package aliyun_mqtt_go

import (
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/zhangsq-ax/aliyun_mqtt_go/constants"
	"github.com/zhangsq-ax/aliyun_mqtt_go/options"
	"log"
	"strings"
	"time"
)

const TokenCheckCycle = 15 // token 检查周期，单位：s

// MQTTClient MQTT 操作对象
type MQTTClient struct {
	client           *mqtt.Client
	options          *options.MQTTClientOptions
	ticker           *time.Ticker
	tokenExpiredTime int64
	tokenUpdating    bool
}

// NewMQTTClient 创建新的 MQTT 客户端对象
func NewMQTTClient(opts *options.MQTTClientOptions) (*MQTTClient, error) {
	client, expiredTime, err := GetClient(opts)
	if err != nil {
		return nil, err
	}

	c := &MQTTClient{
		client:           client,
		options:          opts,
		tokenExpiredTime: expiredTime,
	}

	// 启动 token 过期检查
	go c.startCheckToken()

	return c, nil
}

func (c *MQTTClient) startCheckToken() {
	// 只有 Token 鉴权模式下才需要检查 token 过期
	if c.options.AuthType != constants.AuthTypeToken {
		return
	}

	// 创建定时器
	ticker := time.NewTicker(TokenCheckCycle * time.Second)
	c.ticker = ticker
	defer func() {
		c.ticker.Stop()
		c.ticker = nil
	}()

	for {
		select {
		case t := <-ticker.C:
			timestamp := t.UnixNano() / 1e6
			timeLeft := c.tokenExpiredTime - timestamp
			// 判断是否需要更新 token，剩余时间小于两个检查周期则认为需要更新 token
			if timeLeft < TokenCheckCycle*2*1000 {
				c.updateToken()
			}
		}
	}
}

func (c *MQTTClient) updateToken() {
	log.Println("Start upload token...")
	// 加锁防止更新 token 期间有 Pub 或 Sub 操作导致与 MQTT 的连接断开
	c.tokenUpdating = true

	// 获取新的 token
	if c.options.PasswordGetter == nil {
		log.Fatalf("Unable to update token: There is no PasswordGetter()")
	}
	token, expiredTime, err := c.options.PasswordGetter(c.options.ClientID)
	if err != nil {
		log.Fatalf("Failed to get new token by client %s: %v", c.options.ClientID, err)
	}

	// 分类型解析获取到的 token
	tokens := c.generateUpdateTokenPayload(token)
	for _, t := range tokens {
		// 不断开连接的情况下更新 Token
		err = c.publish(&options.PublishOptions{
			Topic:   "$SYS/uploadToken",
			Qos:     2,
			Payload: t,
		})
		if err != nil {
			log.Fatalf("Failed to update token: %v", err)
		}
	}

	// 记录新的 token 过期时间
	c.tokenExpiredTime = expiredTime
	// 解锁
	c.tokenUpdating = false
	log.Println("Update token complete")
}

// 构建更新 token 的数据
func (c *MQTTClient) generateUpdateTokenPayload(tokens string) []string {
	tmp := strings.Split(tokens, "|")
	result := make([]string, 0)
	for i := 0; i < len(tmp); i = i + 2 {
		t := map[string]string{
			"type":  tmp[i],
			"token": tmp[i+1],
		}
		byteT, _ := json.Marshal(t)
		result = append(result, string(byteT))
	}
	return result
}

// 不加锁的 publish() 方法用于更新 token
func (c *MQTTClient) publish(opts *options.PublishOptions) error {
	client := *c.client
	token := client.Publish(opts.Topic, opts.Qos, opts.Retained, opts.Payload)
	token.Wait()
	return token.Error()
}

// Publish 向 MQTT 服务发布消息
func (c *MQTTClient) Publish(opts *options.PublishOptions) error {
	// 更新 token 期间，等待 2s
	if c.tokenUpdating {
		time.Sleep(2 * time.Second)
	}
	return c.publish(opts)
}

// Subscribe 从 MQTT 服务订阅消息
func (c *MQTTClient) Subscribe(opts *options.SubscribeOptions) (chan *mqtt.Message, error) {
	// 更新 token 期间，等待 2s
	if c.tokenUpdating {
		time.Sleep(2 * time.Second)
	}
	client := *c.client
	ch := make(chan *mqtt.Message)
	token := client.Subscribe(opts.Topic, opts.Qos, func(client mqtt.Client, msg mqtt.Message) {
		ch <- &msg
	})
	token.Wait()
	err := token.Error()
	if err != nil {
		close(ch)
		return nil, err
	}
	return ch, nil
}
