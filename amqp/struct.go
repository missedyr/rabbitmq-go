package rabbitmqGo

const (
	ExchangeDirect  = "direct"
	ExchangeFanout  = "fanout"
	ExchangeTopic   = "topic"
	ExchangeHeaders = "headers"
)

// QueueExchange 定义队列交换机对象
type QueueExchange struct {
	ExchangeType string // 交换机类型 默认topic
	ExchangeName string // 交换机名称 (生产者和消费者 必须)
	RoutingKey   string // 路由key值 (生产者和消费者 必须)  注 支持通配符的场景
	QueueName    string // 队列名称 (生产者非必须  消费者必须）
}

// ConnectConf mq链接信息
type ConnectConf struct {
	InstanceId string // 实例ID  (实例ID存在时 自动使用阿里云AMQP用户名密码转译)
	Endpoint   string // Endpoint配置 或ip
	AccessKey  string // 用户名 或 阿里云AMQP-AccessKey
	SecretKey  string // 密码   或 阿里云AMQP-SecretKey
	Vhost      string
	Port       int64  // 端口号 非必须
}

