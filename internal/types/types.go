package types

// 服务配置
type ServiceConfig struct {
	Name    string            `json:"name"`
	Enabled bool              `json:"enabled"`
	Config  string            `json:"config"` // 配置文件路径
	Args    []string          `json:"args"`   // 启动参数
	Env     map[string]string `json:"env"`    // 环境变量
}

// 订阅信息
type Subscription struct {
	Name       string `json:"name"`
	URL        string `json:"url"`
	Type       string `json:"type"`     // clash, v2ray, ss, ssr, trojan
	Interval   int    `json:"interval"` // 更新间隔(分钟)
	LastUpdate string `json:"last_update"`
}

// 节点信息
type Node struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	Server     string `json:"server"`
	Port       int    `json:"port"`
	Password   string `json:"password,omitempty"`
	Method     string `json:"method,omitempty"`
	Plugin     string `json:"plugin,omitempty"`
	PluginOpts string `json:"plugin_opts,omitempty"`
}

// 系统状态
type SystemStatus struct {
	CPU        float64           `json:"cpu"`
	Memory     float64           `json:"memory"`
	Uptime     int64             `json:"uptime"`
	Version    string            `json:"version"`
	CoreStatus map[string]string `json:"core_status"`
}
