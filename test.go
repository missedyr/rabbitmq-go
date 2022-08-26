package test

import (
	"fmt"
	rabbitmqGo "github.com/missedyr/rabbitmq-go/amqp"
	"time"
)

var connect = rabbitmqGo.ConnectConf{
	Endpoint:   "",
	UserName:   "",
	Password:   "",
	InstanceId: "",
	AccessKey:  "",
	SecretKey:  "",
	Vhost:      "",
	Port:       0,
}

var exConf = rabbitmqGo.QueueExchange{
	ExchangeName: "test-debug",
	RoutingKey:   "debug-topic",
	QueueName:    "debug-xin",
}

func test() {

	// 发送
	rabbitmqGo.New(connect, exConf).Producer("miss-test---11111")
	go runTime()

	// 消费
	rabbitmqGo.New(connect, exConf).Consumer(doFunc)

	fmt.Println("9999999999--")
	time.Sleep(10 * time.Minute)
}

func doFunc(msg string) error {
	fmt.Println("Consumer消费信息--", msg)
	return nil
}

func runTime() {
	C := time.Tick(1 * time.Second)
	for range C {
		msg := fmt.Sprintf(`miss-test--%d`, time.Now().UnixMilli())
		rabbitmqGo.New(connect, exConf).Producer(msg)
	}
}
