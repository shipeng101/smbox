package firewall

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

type FirewallManager struct {
	Mode     string // "iptables" or "nftables"
	CorePath string // core目录的路径
}

func NewFirewallManager(mode string, corePath string) *FirewallManager {
	return &FirewallManager{
		Mode:     mode,
		CorePath: corePath,
	}
}

func (fm *FirewallManager) SetupTransparentProxy(port int) error {
	if fm.Mode == "iptables" {
		return fm.setupIptables(port)
	}
	return fm.setupNftables(port)
}

func (fm *FirewallManager) setupIptables(port int) error {
	rules := []string{
		// 创建新链
		"iptables -t nat -N singdns",
		// 本地地址不转发
		"iptables -t nat -A singdns -d 0.0.0.0/8 -j RETURN",
		"iptables -t nat -A singdns -d 127.0.0.0/8 -j RETURN",
		"iptables -t nat -A singdns -d 192.168.0.0/16 -j RETURN",
		// 转发到代理端口
		fmt.Sprintf("iptables -t nat -A singdns -p tcp -j REDIRECT --to-ports %d", port),
		// 加入到 PREROUTING 链
		"iptables -t nat -A PREROUTING -p tcp -j singdns",
	}

	for _, rule := range rules {
		cmd := exec.Command("sh", "-c", rule)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to setup iptables rule: %s, error: %v", rule, err)
		}
	}
	return nil
}

func (fm *FirewallManager) setupNftables(port int) error {
	nftScript := fmt.Sprintf(`
		table ip nat {
			chain singdns {
				ip daddr { 0.0.0.0/8, 127.0.0.0/8, 192.168.0.0/16 } return
				tcp dport != %d redirect to :%d
			}
			chain prerouting {
				type nat hook prerouting priority 0;
				jump singdns
			}
		}
	`, port, port)

	cmd := exec.Command("nft", "-f", "-")
	cmd.Stdin = strings.NewReader(nftScript)
	return cmd.Run()
}

func (fm *FirewallManager) CleanupRules() error {
	if fm.Mode == "iptables" {
		return fm.cleanupIptables()
	}
	return fm.cleanupNftables()
}

func (fm *FirewallManager) cleanupIptables() error {
	rules := []string{
		"iptables -t nat -D PREROUTING -p tcp -j singdns",
		"iptables -t nat -F singdns",
		"iptables -t nat -X singdns",
	}

	for _, rule := range rules {
		cmd := exec.Command("sh", "-c", rule)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to cleanup iptables rule: %s, error: %v", rule, err)
		}
	}
	return nil
}

func (fm *FirewallManager) cleanupNftables() error {
	cmd := exec.Command("nft", "delete", "table", "ip", "nat")
	return cmd.Run()
}

// 添加服务管理方法
func (fm *FirewallManager) StartService(service string) error {
	execPath := filepath.Join(fm.CorePath, service)
	cmd := exec.Command(execPath)
	return cmd.Start()
}

func (fm *FirewallManager) StopService(service string) error {
	cmd := exec.Command("pkill", "-f", service)
	return cmd.Run()
}
