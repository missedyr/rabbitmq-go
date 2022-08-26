# rabbitmq-go
rabbitmq精简使用包--go版本

### 安装

go get github.com/missedyr/rabbitmq-go

### 使用

```go
connectConf := rabbitmqGo.ConnectConf{
    Endpoint:   "", // Endpoint配置 或ip
    UserName:   "", // 用户名 		非必需 (注* 用户名密码登录时为 必须)
    Password:   "", // 密码 		非必需 (注* 用户名密码登录时为 必须)
    InstanceId: "", // 实例ID 		非必需 (注* key和密钥登录时为 必须)
    AccessKey:  "", // AccessKey	非必需 (注* key和密钥登录时为 必须)
    SecretKey:  "", // SecretKey	非必需 (注* key和密钥登录时为 必须)
    Vhost:      "", // 非必需 默认值 default
    Port:       0   // 端口号 非必须
}
queueConf := rabbitmqGo.QueueExchange{
    ExchangeName: "", // 交换机名称 (生产者和消费者 必须)
    RoutingKey:   "", // 路由key值 (生产者和消费者 必须) 注 支持通配符的场景
    QueueName:    "", // 队列名称 (生产者非必须  消费者必须）
    ExchangeType: "", // 交换机类型 (非必须 默认topic模式)
}
```

### 发送消息 (发送string)
```go
rabbitmqGo.New(connectConf, queueConf).Producer("miss-test")
```

### 接收消息 (传入自定义的消费消息的处理方法)
```go
func doFunc(msg string) error {
    fmt.Println("Consumer消费信息--", msg)
    return nil
}

rabbitmqGo.New(connectConf, queueConf).Consumer(doFunc)
```
