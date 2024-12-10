package constant

// 服务名称
const (
	ServiceSingbox = "singbox"
	ServiceMosdns  = "mosdns"
	ServiceSingdns = "singdns"
)

// 默认配置
const (
	DefaultPort    = 7890
	DefaultDnsPort = 53
	DefaultWebPort = 9090
)

// 配置文件路径
const (
	ConfigDir       = "/etc/singdns"
	CoreDir         = "/usr/local/singdns/core"
	LogDir          = "/var/log/singdns"
	RuleDir         = "/etc/singdns/rules"
	SubscriptionDir = "/etc/singdns/subscriptions"
)

// 服务状态
const (
	StatusRunning = "running"
	StatusStopped = "stopped"
	StatusError   = "error"
)
