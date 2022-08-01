package test

import (
	"fmt"
	rabbitmqGo "github.com/missedyr/rabbitmq-go/amqp"
	"time"
)

var connect = rabbitmqGo.ConnectConf{
	InstanceId: "",
	Endpoint:   "",
	Port:       0,
	AccessKey:  "",
	SecretKey:  "",
	Vhost:      "",
}
var exConf = rabbitmqGo.QueueExchange{
	ExchangeName: "miss",
	RoutingKey:   "miss-topic",
	QueueName:    "miss-xin",
}

func test() {

	// 发送
	rabbitmqGo.New(connect, exConf).Producer("miss-test")
	rabbitmqGo.New(connect, exConf).Producer("miss-test---11111")
	go runTime()

	// 消费
	rabbitmqGo.New(connect, exConf).Consumer(doFunc)

	fmt.Println("9999999999--")
	time.Sleep(1 * time.Minute)
}

func doFunc(msg string) error {
	fmt.Println("Consumer消费信息--", msg)
	return nil
}

func runTime() {
	C := time.Tick(5 * time.Second)
	for range C {
		msg := fmt.Sprintf(`miss-test--%d`, time.Now().UnixMilli())
		rabbitmqGo.New(connect, exConf).Producer(msg)
	}
}
