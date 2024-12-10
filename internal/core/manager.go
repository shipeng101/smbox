package core

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/yourusername/singdns/internal/constant"
	"github.com/yourusername/singdns/internal/types"
)

type Manager struct {
	config  *types.ServiceConfig
	process *os.Process
	status  string
	mu      sync.RWMutex
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewManager(cfg *types.ServiceConfig) *Manager {
	ctx, cancel := context.WithCancel(context.Background())
	return &Manager{
		config: cfg,
		status: constant.StatusStopped,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (m *Manager) Start() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.status == constant.StatusRunning {
		return fmt.Errorf("service already running")
	}

	// 检查可执行文件
	execPath := filepath.Join(constant.CoreDir, m.config.Name)
	if _, err := os.Stat(execPath); os.IsNotExist(err) {
		return fmt.Errorf("executable not found: %s", execPath)
	}

	// 检查配置文件
	if _, err := os.Stat(m.config.Config); os.IsNotExist(err) {
		return fmt.Errorf("config file not found: %s", m.config.Config)
	}

	// 准备命令
	cmd := exec.Command(execPath)
	cmd.Args = append(cmd.Args, m.config.Args...)
	cmd.Env = os.Environ()
	for k, v := range m.config.Env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	// 启动进程
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start service: %v", err)
	}

	m.process = cmd.Process
	m.status = constant.StatusRunning

	// 监控进程
	go m.monitor()

	return nil
}

func (m *Manager) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.status != constant.StatusRunning {
		return nil
	}

	m.cancel()
	if err := m.process.Kill(); err != nil {
		return fmt.Errorf("failed to kill process: %v", err)
	}

	m.status = constant.StatusStopped
	return nil
}

func (m *Manager) Status() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.status
}

func (m *Manager) monitor() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-m.ctx.Done():
			return
		case <-ticker.C:
			if m.process == nil {
				continue
			}

			// 检查进程状态
			if err := m.process.Signal(syscall.Signal(0)); err != nil {
				m.mu.Lock()
				m.status = constant.StatusError
				m.process = nil
				m.mu.Unlock()

				// 尝试重启
				if m.config.Enabled {
					go func() {
						if err := m.Start(); err != nil {
							log.Printf("Failed to restart service %s: %v", m.config.Name, err)
						}
					}()
				}
			}
		}
	}
}
