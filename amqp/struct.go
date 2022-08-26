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
	Endpoint   string // Endpoint配置 或ip
	UserName   string // 用户名 		非必需 (注* 用户名密码登录时为 必须)
	Password   string // 密码 		非必需 (注* 用户名密码登录时为 必须)
	InstanceId string // 实例ID 		非必需 (注* key和密钥登录时为 必须)
	AccessKey  string // AccessKey	非必需 (注* key和密钥登录时为 必须)
	SecretKey  string // SecretKey	非必需 (注* key和密钥登录时为 必须)
	Vhost      string // 非必需 默认值 default
	Port       int64  // 端口号 非必须
}
