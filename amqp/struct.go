package rabbitmqGo

type ExchangeType string

const (
	ExchangeDirect  ExchangeType = "direct"
	ExchangeFanout  ExchangeType = "fanout"
	ExchangeTopic   ExchangeType = "topic"
	ExchangeHeaders ExchangeType = "headers"
)

// ConnectConf mq链接信息
type ConnectConf struct {
	Endpoint   string // Endpoint配置 或ip
	UserName   string // 用户名 		非必需 (注* 用户名密码登录时为 必须)
	Password   string // 密码 		非必需 (注* 用户名密码登录时为 必须)
	InstanceId string // 实例ID 		非必需 (注* key和密钥登录时为 必须)
	AccessKey  string // AccessKey	非必需 (注* key和密钥登录时为 必须)
	SecretKey  string // SecretKey	非必需 (注* key和密钥登录时为 必须)
	Vhost      string // 非必需 默认值 default
	Port       int64  // 端口号 非必须
}

// ProducerRoutingConf 生产者配置
type ProducerRoutingConf struct {
	ExchangeType ExchangeType // 交换机类型 默认topic
	ExchangeName string       // 交换机名称 必须
	RoutingKey   string       // 路由key值 必须  注 支持通配符的场景
}

// ConsumerQueueConf 消费者配置
type ConsumerQueueConf struct {
	ConsumerTag  string       // 标签	非必须
	ExchangeType ExchangeType // 交换机类型 默认topic
	ExchangeName string       // 交换机名称 必须
	RoutingKey   string       // 路由key值 必须  注 支持通配符的场景
	QueueName    string       // 队列名称  必须
	Requeue      bool         // 是否重排任务 Nack->requeue
}
