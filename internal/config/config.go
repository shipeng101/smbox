package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/yourusername/singdns/internal/constant"
)

type Config struct {
	Core struct {
		Path string `json:"path"`
		Port int    `json:"port"`
	} `json:"core"`

	DNS struct {
		Enable    bool `json:"enable"`
		Port      int  `json:"port"`
		AdBlock   bool `json:"adblock"`
		CacheSize int  `json:"cache_size"`
	} `json:"dns"`

	Proxy struct {
		Enable       bool   `json:"enable"`
		Mode         string `json:"mode"`
		Subscription string `json:"subscription"`
	} `json:"proxy"`

	Remote struct {
		Enable bool   `json:"enable"`
		Host   string `json:"host"`
		Port   int    `json:"port"`
	} `json:"remote"`

	Web struct {
		Port     int    `json:"port"`
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"web"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// 设置默认值
	if cfg.Core.Port == 0 {
		cfg.Core.Port = constant.DefaultPort
	}
	if cfg.DNS.Port == 0 {
		cfg.DNS.Port = constant.DefaultDnsPort
	}
	if cfg.Web.Port == 0 {
		cfg.Web.Port = constant.DefaultWebPort
	}

	return &cfg, nil
}

func (c *Config) Save(path string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
