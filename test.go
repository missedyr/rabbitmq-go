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
	Vhost:      "devices",
	Port:       0,
}

var exPrdConf = rabbitmqGo.ProducerRoutingConf{
	ExchangeName: "test-debug",
	RoutingKey:   "debug-topic",
}

var exCsmConf = rabbitmqGo.ConsumerQueueConf{
	ConsumerTag:  "test--miss--debug",
	ExchangeName: "test-debug",
	RoutingKey:   "debug-topic",
	QueueName:    "debug-xin",
	Requeue:      true,
}

func test() {

	// 发送
	rabbitmqGo.New(connect, exPrdConf, exCsmConf).Producer("miss-test---222")
	//go runTime()

	// 消费
	rabbitmqGo.New(connect, exPrdConf, exCsmConf).Consumer(doFunc)

	fmt.Println("9999999999--")
	time.Sleep(10 * time.Minute)
}

func doFunc(msg string) error {
	var err error
	fmt.Println("Consumer消费信息--", msg)
	return err
}

func runTime() {
	C := time.Tick(5 * time.Second)
	for range C {
		msg := fmt.Sprintf(`miss-test--%d`, time.Now().UnixMilli())
		rabbitmqGo.New(connect, exPrdConf, exCsmConf).Producer(msg)
	}
}
