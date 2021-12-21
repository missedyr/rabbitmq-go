package test

import (
	"fmt"
	"github.com/xuexin520/rabbitmq-go/amqp"
)

func test() {
	connect := rabbitmqGo.RabbitMQConnectConf{
		InstanceId: "",
		Endpoint:   "",
		Port:       0,
		AccessKey:  "",
		SecretKey:  "",
		Vhost:      "",
	}
	exConf := rabbitmqGo.RabbitMQQueueExchange{
		ExchangeName: "miss",
		RoutingKey:   "miss-topic",
		QueueName:    "miss-xin",
	}

	// 发送
	rabbitmqGo.New(connect, exConf).Producer("miss-test")

	// 消费
	rabbitmqGo.New(connect, exConf).Consumer(doFunc)
}

func doFunc(msg string) error {
	fmt.Println("Consumer消费信息--", msg)
	return nil
}
