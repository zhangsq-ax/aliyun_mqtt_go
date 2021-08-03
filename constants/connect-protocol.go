package constants

// ConnectProtocol 连接协议
type ConnectProtocol string

const (
	ConnectProtocolTcp ConnectProtocol = "tcp" // Port 1883
	ConnectProtocolSsl ConnectProtocol = "ssl" // Port 8883
	ConnectProtocolWs  ConnectProtocol = "ws"  // Port 80
	ConnectProtocolWss ConnectProtocol = "wss" // Port 443
)

// ConnectPort 连接协议对应的默认端口
var ConnectPort = map[ConnectProtocol]int{
	ConnectProtocolTcp: 1883,
	ConnectProtocolSsl: 8883,
	ConnectProtocolWs:  80,
	ConnectProtocolWss: 443,
}
