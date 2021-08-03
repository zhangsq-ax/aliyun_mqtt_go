package options

// PublishOptions 发送消息设置
type PublishOptions struct {
	Topic    string
	Qos      byte
	Retained bool
	Payload  interface{}
}
