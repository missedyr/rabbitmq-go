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
	Vhost:      "primary",
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
	go runTime()
	//for i := 1; i <= 10000; i++ {
	//	runTime()
	//	fmt.Printf("启动弟 %d\n", i)
	//}

	// 消费
	rabbitmqGo.NewConsumer(connect, exCsmConf).Consumer(doFunc1)
	rabbitmqGo.NewConsumer(connect, exCsmConf).Consumer(doFunc2)

	fmt.Println("9999999999--")
	time.Sleep(10 * time.Minute)
}

func doFunc1(msg string) error {
	var err error
	fmt.Println("Consumer消费信息--1111---", msg)
	return err
}

func doFunc2(msg string) error {
	var err error
	fmt.Println("Consumer消费信息--2222---", msg)
	return err
}

func runTime() {
	cline := rabbitmqGo.New(connect, exPrdConf, exCsmConf)
	st := time.Now().UnixMilli()
	for i := 1; i <= 2; i++ {
		msg := fmt.Sprintf(`miss-test--%d`, time.Now().UnixMilli())
		cline.Producer(msg, -1)
	}
	et := time.Now().UnixMilli()
	fmt.Println("耗时-----", et-st)
	cline.MqClose()
}
