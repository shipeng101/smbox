package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/yourusername/singdns/internal/firewall"
	"github.com/yourusername/singdns/internal/service"
)

func main() {
	var (
		corePath string
		port     int
	)

	flag.StringVar(&corePath, "core", "./core", "Path to core executables")
	flag.IntVar(&port, "port", 7890, "Transparent proxy port")
	flag.Parse()

	// 确保 core 目录存在
	if _, err := os.Stat(corePath); os.IsNotExist(err) {
		log.Fatalf("Core directory not found: %s", corePath)
	}

	// 初始化服务管理器
	sm := service.NewServiceManager(corePath)

	// 初始化防火墙管理器
	fm := firewall.NewFirewallManager("iptables", corePath)

	// 启动核心服务
	if err := sm.StartSingbox(); err != nil {
		log.Printf("Failed to start singbox: %v", err)
	}

	if err := sm.StartMosdns(); err != nil {
		log.Printf("Failed to start mosdns: %v", err)
	}

	// 设置透明代理
	if err := fm.SetupTransparentProxy(port); err != nil {
		log.Printf("Failed to setup transparent proxy: %v", err)
	}

	// 处理退出信号
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	<-sigCh
	log.Println("Shutting down...")

	// 清理防火墙规则
	if err := fm.CleanupRules(); err != nil {
		log.Printf("Failed to cleanup firewall rules: %v", err)
	}

	// 停止服务
	if err := sm.StopService(service.ServiceSingbox); err != nil {
		log.Printf("Failed to stop singbox: %v", err)
	}

	if err := sm.StopService(service.ServiceMosdns); err != nil {
		log.Printf("Failed to stop mosdns: %v", err)
	}
}
