# aliyun_mqtt_go
MQTT connection auxiliary package based on Aliyun



## Usage

### Installation

```
go get github.com/zhangsq-ax/aliyun_mqtt_go
```

### Example

Publish message

```go
package main

import (
  "github.com/zhangsq-ax/aliyun_mqtt_go"
  "github.com/zhangsq-ax/aliyun_mqtt_go/constants"
  "github.com/zhangsq-ax/aliyun_mqtt_go/options"
)

func main() {
  // Get MQTT client
  client, err := aliyun_mqtt_go.NewMQTTClient(&options.MQTTClientOptions{
    AuthType:   constants.AuthTypeSign,
		Protocol:   constants.ConnectProtocolSsl,
		InstanceID: "xxxx-xx-xxxxxxxxxxx",
		Endpoints:  []string{"xxxx-xx-xxxxxxxxxx.mqtt.aliyuncs.com"},
		Username:   "xxxxxxxxxxxxxxxxx",
		Password:   "xxxxxxxxxxxxxxxxx",
		GroupID:        "GID_xxxx",
		ClientID:       "xxxxxxxxxx",
  })
  if err != nil {
    panic(err)
  }
  
  // Publish message
  err = client.Publish(&options.PublishOptions{
    Topic: "test",
    Qos: 2,
    Payload: "hello mqtt",
  })
  if err != nil {
    panic(err)
  }
}
```



Subscribe Message

```go
package main

import (
  "fmt"
  "github.com/zhangsq-ax/aliyun_mqtt_go"
  "github.com/zhangsq-ax/aliyun_mqtt_go/constants"
  "github.com/zhangsq-ax/aliyun_mqtt_go/options"
)

func main() {
  // Get MQTT client
  client, err := aliyun_mqtt_go.NewMQTTClient(&options.MQTTClientOptions{
    AuthType:   constants.AuthTypeSign,
		Protocol:   constants.ConnectProtocolSsl,
		InstanceID: "xxxx-xx-xxxxxxxxxxx",
		Endpoints:  []string{"xxxx-xx-xxxxxxxxxx.mqtt.aliyuncs.com"},
		Username:   "xxxxxxxxxxxxxxxxx",
		Password:   "xxxxxxxxxxxxxxxxx",
		GroupID:        "GID_xxxx",
		ClientID:       "xxxxxxxxxx",
  })
  if err != nil {
    panic(err)
  }
  
  // Subscribe message
  msgChan, err := client.Subscribe(&options.SubscribeOptions{
    Topic: "test",
    Qos: 2,
  })
  if err != nil {
    panic(err)
  }
  
  for msg := range msgChan {
    fmt.Println(msg.Payload())
  }
}
```

