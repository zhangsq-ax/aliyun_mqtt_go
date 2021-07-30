# aliyun_mqtt_go
MQTT connection auxiliary package based on Aliyun



## Usage

### Installation

```
go get github.com/zhangsq-ax/aliyun_mqtt_go
```

### Example

```go
package main

import (
  "github.com/zhangsq-ax/aliyun_mqtt_go"
  "github.com/zhangsq-ax/aliyun_mqtt_go/constants"
)

func main() {
  // Get helper instance
  helper := &aliyun_mqtt_go.ConnectHelper{
    AuthType:   constants.AuthTypeSign,                             // Authentication type. Refer to https://github.com/zhangsq-ax/aliyun_mqtt_go/blob/main/constants/auth-credential.go
		Protocol:   constants.ConnectProtocolSsl,                       // The protocol of connect to MQTT server, Refer to https://github.com/zhangsq-ax/aliyun_mqtt_go/blob/main/constants/connect-protocol.go
		InstanceID: "xxxx-xx-xxxxxxxxxxx",                              // The MQTT server instance ID
		Endpoints:  []string{"xxxx-xx-xxxxxxxxxxx.mqtt.aliyuncs.com"},  // The connect endpoints of MQTT server
		Username:   "xxxxxxxxxxxxxxxx",
		Password:   "xxxxxxxxxxxxxxxx",
  }
  
  // Get connected client
  client, err := helper.GetConnectedClient(&constants.ClientOptions{
    GroupID: "GID_test",              // The client group ID
    ClientID: "client_test",          // Client ID
  })
  if err != nil {
    panic(err)
  }
  
  // Publish message. Refer to https://github.com/eclipse/paho.mqtt.golang
  token := client.Publish("test", 2, false, "hello world")
  token.Wait()
  err = token.Error()
  if err != nil {
    panic(err)
  }
}
```

