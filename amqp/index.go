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
	ConnectConf         ConnectConf
	ProducerRoutingConf ProducerRoutingConf
	ConsumerQueueConf   ConsumerQueueConf
	Connection          *amqp.Connection
	Channel             *amqp.Channel
	Mu                  sync.RWMutex
	RabbitUrl           string
}

// New 创建一个新的操作对象
func New(connect ConnectConf, producerRoutingConf ProducerRoutingConf, consumerQueueConf ConsumerQueueConf) *RabbitMQ {
	if len(producerRoutingConf.ExchangeType) == 0 {
		producerRoutingConf.ExchangeType = ExchangeTopic
	}
	if len(consumerQueueConf.ExchangeType) == 0 {
		consumerQueueConf.ExchangeType = ExchangeTopic
	}
	return &RabbitMQ{
		ConnectConf:         connect,
		ProducerRoutingConf: producerRoutingConf,
		ConsumerQueueConf:   consumerQueueConf,
	}
}

// MqConnect 链接rabbitMQ
func (r *RabbitMQ) MqConnect() (*RabbitMQ, error) {
	var err error
	r.GetRabbitUrl()
	mqConn, err = amqp.Dial(r.RabbitUrl)
	r.Connection = mqConn // 赋值给RabbitMQ对象
	if err != nil {
		fmt.Printf("RabbitMQ链接失败url:%s ---- err:%s \n", r.RabbitUrl, err)
		return r, err
	}
	mqChan, err = mqConn.Channel()
	r.Channel = mqChan // 赋值给RabbitMQ对象
	if err != nil {
		fmt.Printf("RabbitMQ管道失败:%s \n", err)
	}
	return r, err
}

// MqClose 关闭RabbitMQ连接
func (r *RabbitMQ) MqClose() error {
	var err error
	// 先关闭管道,再关闭链接
	if r.Channel != nil {
		err = r.Channel.Close()
		if err != nil {
			fmt.Printf("RabbitMQ管道关闭失败:%s \n", err)
		}
	}
	if r.Connection != nil {
		err = r.Connection.Close()
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
	if r.Channel == nil {
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
	err = r.Channel.ExchangeDeclare(r.ProducerRoutingConf.ExchangeName, string(r.ProducerRoutingConf.ExchangeType), true, false, false, false, nil)
	if err != nil {
		fmt.Printf("RabbitMQ注册交换机失败:%s \n", err)
		return err
	}

	// 发送任务消息
	err = r.Channel.Publish(r.ProducerRoutingConf.ExchangeName, r.ProducerRoutingConf.RoutingKey, false, false, amqp.Publishing{
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
	// 验证链接是否正常,否则重新链接
	if r.Channel == nil {
		_, err := r.MqConnect()
		if err != nil {
			return
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
	err := r.Channel.ExchangeDeclare(r.ConsumerQueueConf.ExchangeName, string(r.ConsumerQueueConf.ExchangeType), true, false, false, false, nil)
	if err != nil {
		fmt.Printf("RabbitMQ注册交换机失败:%s \n", err)
		return
	}

	// 用于检查队列是否存在,已经存在不需要重复声明
	_, err = r.Channel.QueueDeclare(r.ConsumerQueueConf.QueueName, true, false, false, false, nil)
	if err != nil {
		fmt.Printf("RabbitMQ注册队列失败:%s \n", err)
		return
	}

	// 绑定d队列
	// name:队列名称;
	// durable:是否持久化,队列存盘,true服务重启后信息不会丢失,影响性能;
	// autoDelete:是否自动删除; 最后一个Consumer取消订阅后，Queue是否自动删除。
	// exclusive:是否设置排他
	// noWait:是否非阻塞, true为是,不等待RMQ返回信息;
	// args:参数,传nil即可;
	err = r.Channel.QueueBind(r.ConsumerQueueConf.QueueName, r.ConsumerQueueConf.RoutingKey, r.ConsumerQueueConf.ExchangeName, false, nil)
	if err != nil {
		fmt.Printf("RabbitMQ绑定队列失败:%s \n", err)
		return
	}
	// 获取消费通道,确保rabbitMQ一个一个发送消息
	err = r.Channel.Qos(1, 0, true)

	// 消息消费
	msgList, err := r.Channel.Consume(r.ConsumerQueueConf.QueueName, r.ConsumerQueueConf.ConsumerTag, false, false, false, false, nil)
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
				err = msg.Nack(true, r.ConsumerQueueConf.Requeue)
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
	userName := r.ConnectConf.UserName
	password := r.ConnectConf.Password
	if len(userName) == 0 || len(password) == 0 {
		userName = GetUserName(r.ConnectConf.AccessKey, r.ConnectConf.InstanceId)
		password = GetPassword(r.ConnectConf.SecretKey)
		if len(userName) == 0 || len(password) == 0 {
			fmt.Printf("RabbitMQ--GetRabbitUrl--使用AccessKey和SecretKey动态获取用户名密码失败 \n")
			return
		}
	}

	r.RabbitUrl = fmt.Sprintf("amqp://%s:%s@%s", userName, password, r.ConnectConf.Endpoint)
	if r.ConnectConf.Port != 0 {
		r.RabbitUrl = fmt.Sprintf(`%s:%d`, r.RabbitUrl, r.ConnectConf.Port)
	}
	if len(r.ConnectConf.Vhost) == 0 {
		r.ConnectConf.Vhost = "default"
	}

	r.RabbitUrl = fmt.Sprintf(`%s/%s`, r.RabbitUrl, r.ConnectConf.Vhost)
	r.RabbitUrl += "?heartbeat=5"
}
