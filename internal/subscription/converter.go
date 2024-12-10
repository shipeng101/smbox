package subscription

import (
	"fmt"
	"io"
	"net/http"
)

type SingboxConfig struct {
	Outbounds []Outbound `json:"outbounds"`
	Route     Route      `json:"route"`
}

type Outbound struct {
	Type     string      `json:"type"`
	Tag      string      `json:"tag"`
	Settings interface{} `json:"settings"`
}

type Route struct {
	Rules []Rule `json:"rules"`
}

type Rule struct {
	Type        string   `json:"type"`
	OutboundTag string   `json:"outbound_tag"`
	Domain      []string `json:"domain,omitempty"`
	IPCIDRs     []string `json:"ip_cidr,omitempty"`
}

type Converter struct {
	client *http.Client
}

func NewConverter() *Converter {
	return &Converter{
		client: &http.Client{},
	}
}

func (c *Converter) Convert(url string, format string) (*SingboxConfig, error) {
	// 获取订阅内容
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch subscription: %v", err)
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	// 根据格式转换
	switch format {
	case "clash":
		return c.convertFromClash(content)
	case "v2ray":
		return c.convertFromV2ray(content)
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}

func (c *Converter) convertFromClash(content []byte) (*SingboxConfig, error) {
	// 实现Clash配置转换逻辑
	return nil, nil
}

func (c *Converter) convertFromV2ray(content []byte) (*SingboxConfig, error) {
	// 实现V2ray配置转换逻辑
	return nil, nil
}
