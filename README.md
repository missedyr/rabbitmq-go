# rabbitmq-go
rabbitmq精简使用包--go版本

### 安装

go get github.com/missedyr/rabbitmq-go

### 使用

```go
connectConf := rabbitmqGo.ConnectConf{
    InstanceId: "", // 实例ID (实例ID存在 自动使用阿里云AMQP用户名密码转译)
    Endpoint:   "", // Endpoint配置 或ip
    Port:       0,  // 端口号 非必须
    AccessKey:  "", // 用户名 或 阿里云AMQP-AccessKey
    SecretKey:  "", // 密码   或 阿里云AMQP-SecretKey
    Vhost:      "",
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
