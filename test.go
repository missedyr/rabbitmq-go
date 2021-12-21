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
	mq := rabbitmqGo.New(connect, exConf)
	mq.Producer("miss-test")
	mq.Consumer(doFunc)
}

func doFunc(msg string) error {
	fmt.Println("Consumer消费信息--", msg)
	return nil
}
