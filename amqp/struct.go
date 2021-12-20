package rabbitmqGo

const (
	ExchangeDirect  = "direct"
	ExchangeFanout  = "fanout"
	ExchangeTopic   = "topic"
	ExchangeHeaders = "headers"
)

// RabbitMQQueueExchange 定义队列交换机对象
type RabbitMQQueueExchange struct {
	ExchangeType string // 交换机类型 默认topic
	ExchangeName string // 交换机名称 (生产者和消费者 必须)
	RoutingKey   string // 路由key值 (生产者和消费者 必须)  注 支持通配符的场景
	QueueName    string // 队列名称 (生产者非必须  消费者必须）
}

// RabbitMQConnectConf mq链接信息
type RabbitMQConnectConf struct {
	InstanceId string // 实例ID  (实例ID存在  默认转化阿里云AMQP 用户名密码转译)
	Endpoint   string // Endpoint配置 或ip
	AccessKey  string
	SecretKey  string
	Vhost      string
	Port       int64  // 端口号 非必须
}

