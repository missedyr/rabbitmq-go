package rabbitmqGo

import (
	"fmt"
	"github.com/streadway/amqp"
	"sync"
)

// 定义全局变量,指针类型
var mqConn *amqp.Connection
var mqChan *amqp.Channel

// RabbitMQ 定义RabbitMQ对象
type RabbitMQ struct {
	connectConf ConnectConf
	queueConf   QueueExchange // 队列配置
	connection  *amqp.Connection
	channel     *amqp.Channel
	mu          sync.RWMutex
	rabbitUrl   string
}

// New 创建一个新的操作对象
func New(connect ConnectConf, queueConf QueueExchange) *RabbitMQ {
	if queueConf.ExchangeType == "" {
		queueConf.ExchangeType = ExchangeTopic
	}
	return &RabbitMQ{
		connectConf: connect,
		queueConf:   queueConf,
	}
}

// MqConnect 链接rabbitMQ
func (r *RabbitMQ) MqConnect() (*RabbitMQ, error) {
	var err error
	r.GetRabbitUrl()
	mqConn, err = amqp.Dial(r.rabbitUrl)
	r.connection = mqConn // 赋值给RabbitMQ对象
	if err != nil {
		fmt.Printf("RabbitMQ链接失败url:%s ---- err:%s \n", r.rabbitUrl, err)
	}
	mqChan, err = mqConn.Channel()
	r.channel = mqChan // 赋值给RabbitMQ对象
	if err != nil {
		fmt.Printf("RabbitMQ管道失败:%s \n", err)
	}
	return r, err
}

// MqClose 关闭RabbitMQ连接
func (r *RabbitMQ) MqClose() error {
	var err error
	// 先关闭管道,再关闭链接
	if r.channel != nil {
		err = r.channel.Close()
		if err != nil {
			fmt.Printf("RabbitMQ管道关闭失败:%s \n", err)
		}
	}
	if r.connection != nil {
		err = r.connection.Close()
		if err != nil {
			fmt.Printf("RabbitMQ链接关闭失败:%s \n", err)
		}
	}
	return err
}

// Producer 发送任务 送指定队列指定路由的生产者
func (r *RabbitMQ) Producer(msg string) error {
	var err error
	// 处理结束关闭链接
	defer r.MqClose()
	if err != nil {
		return err
	}

	// 验证链接是否正常,否则重新链接
	if r.channel == nil {
		_, err = r.MqConnect()
		if err != nil {
			return err
		}
	}

	// 注册交换机
	// name:交换机名称,
	// kind:交换机类型,
	// durable:是否持久化,队列存盘,true服务重启后信息不会丢失,影响性能;
	// autoDelete:是否自动删除;
	// internal:是否为内部
	// noWait:是否非阻塞, true为是,不等待RMQ返回信息;
	// args:参数,传nil即可;
	err = r.channel.ExchangeDeclare(r.queueConf.ExchangeName, r.queueConf.ExchangeType, true, false, false, false, nil)
	if err != nil {
		fmt.Printf("RabbitMQ注册交换机失败:%s \n", err)
		return err
	}

	// 发送任务消息
	err = r.channel.Publish(r.queueConf.ExchangeName, r.queueConf.RoutingKey, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(msg),
	})
	if err != nil {
		fmt.Printf("RabbitMQ消息发送失败:%s \n", err)
		return err
	}
	fmt.Printf("RabbitMQ-Producer-消息发送成功\n")
	return err
}

// Consumer 接收任务消费消息 接收指定队列指定路由的数据接收者
func (r *RabbitMQ) Consumer(doFunc func(string) error) {
	// 处理结束关闭链接
	//defer r.MqClose()

	// 验证链接是否正常
	if r.channel == nil {
		r.MqConnect()
	}

	// 注册交换机
	// name:交换机名称,
	// kind:交换机类型,
	// durable:是否持久化,队列存盘,true服务重启后信息不会丢失,影响性能;
	// autoDelete:是否自动删除;
	// internal:是否为内部
	// noWait:是否非阻塞, true为是,不等待RMQ返回信息;
	// args:参数,传nil即可;
	err := r.channel.ExchangeDeclare(r.queueConf.ExchangeName, r.queueConf.ExchangeType, true, false, false, false, nil)
	if err != nil {
		fmt.Printf("RabbitMQ注册交换机失败:%s \n", err)
		return
	}

	// 用于检查队列是否存在,已经存在不需要重复声明
	_, err = r.channel.QueueDeclare(r.queueConf.QueueName, true, false, false, false, nil)
	if err != nil {
		fmt.Printf("RabbitMQ注册队列失败:%s \n", err)
		return
	}

	// 绑定d队列
	// name:队列名称;
	// durable:是否持久化,队列存盘,true服务重启后信息不会丢失,影响性能;
	// autoDelete:是否自动删除;
	// exclusive:是否设置排他
	// noWait:是否非阻塞, true为是,不等待RMQ返回信息;
	// args:参数,传nil即可;
	err = r.channel.QueueBind(r.queueConf.QueueName, r.queueConf.RoutingKey, r.queueConf.ExchangeName, false, nil)
	if err != nil {
		fmt.Printf("RabbitMQ绑定队列失败:%s \n", err)
		return
	}
	// 获取消费通道,确保rabbitMQ一个一个发送消息
	err = r.channel.Qos(1, 0, true)

	// 消息消费
	msgList, err := r.channel.Consume(r.queueConf.QueueName, "miss-c-n", false, false, false, false, nil)
	if err != nil {
		fmt.Printf("RabbitMQ获取消费通道异常:%s \n", err)
		return
	}

	// 处理数据
	go func() {
		for msg := range msgList {
			// 处理数据
			err := doFunc(string(msg.Body))
			fmt.Printf("RabbitMQ--Consume--doFunc监听消息执行结果 err:%s \n", err)
			if err != nil {
				err = msg.Ack(true)
				if err != nil {
					fmt.Printf("确认消息未完成异常:%s \n", err)
				}
			} else {
				// 确认消息,必须为false
				err = msg.Ack(false)
				if err != nil {
					fmt.Printf("确认消息完成异常:%s \n", err)
				}
			}
		}
	}()
	return
}

func (r *RabbitMQ) GetRabbitUrl() {
	//rabbitUrl := fmt.Sprintf("amqp://%s:%s@%s:%d/", "guest", "guest", "******", 5673)
	userName := r.connectConf.AccessKey
	password := r.connectConf.SecretKey
	if r.connectConf.InstanceId != "" { // 实例ID存在  默认转化阿里云AMQP 用户名密码转译
		userName = GetUserName(r.connectConf.AccessKey, r.connectConf.InstanceId)
		password = GetPassword(r.connectConf.SecretKey)
	}
	r.rabbitUrl = fmt.Sprintf("amqp://%s:%s@%s", userName, password, r.connectConf.Endpoint)
	if r.connectConf.Port != 0 {
		r.rabbitUrl = fmt.Sprintf(`%s:%d`, r.rabbitUrl, r.connectConf.Port)
	}
	if r.connectConf.Vhost != "" {
		r.rabbitUrl = fmt.Sprintf(`%s/%s`, r.rabbitUrl, r.connectConf.Vhost)
	}

	r.rabbitUrl += "?heartbeat=5"
}
