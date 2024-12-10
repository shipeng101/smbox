package subscription

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/yourusername/singdns/internal/constant"
	"github.com/yourusername/singdns/internal/types"
)

type Manager struct {
	subscriptions map[string]*types.Subscription
	converter     *Converter
	mu            sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		subscriptions: make(map[string]*types.Subscription),
		converter:     NewConverter(),
	}
}

func (m *Manager) AddSubscription(sub *types.Subscription) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.subscriptions[sub.Name]; exists {
		return fmt.Errorf("subscription already exists: %s", sub.Name)
	}

	m.subscriptions[sub.Name] = sub
	return m.saveSubscriptions()
}

func (m *Manager) UpdateSubscription(name string) error {
	m.mu.Lock()
	sub, exists := m.subscriptions[name]
	m.mu.Unlock()

	if !exists {
		return fmt.Errorf("subscription not found: %s", name)
	}

	// 获取订阅内容
	resp, err := http.Get(sub.URL)
	if err != nil {
		return fmt.Errorf("failed to fetch subscription: %v", err)
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}

	// 转换配置
	config, err := m.converter.Convert(string(content), sub.Type)
	if err != nil {
		return fmt.Errorf("failed to convert subscription: %v", err)
	}

	// 保存配置
	configPath := filepath.Join(constant.SubscriptionDir, fmt.Sprintf("%s.json", name))
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return err
	}

	// 更新状态
	m.mu.Lock()
	sub.LastUpdate = time.Now().Format(time.RFC3339)
	m.mu.Unlock()

	return m.saveSubscriptions()
}

func (m *Manager) StartAutoUpdate() {
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			m.mu.RLock()
			subs := make([]*types.Subscription, 0, len(m.subscriptions))
			for _, sub := range m.subscriptions {
				subs = append(subs, sub)
			}
			m.mu.RUnlock()

			for _, sub := range subs {
				lastUpdate, _ := time.Parse(time.RFC3339, sub.LastUpdate)
				if time.Since(lastUpdate).Minutes() >= float64(sub.Interval) {
					if err := m.UpdateSubscription(sub.Name); err != nil {
						log.Printf("Failed to update subscription %s: %v", sub.Name, err)
					}
				}
			}
		}
	}()
}

func (m *Manager) saveSubscriptions() error {
	data, err := json.MarshalIndent(m.subscriptions, "", "  ")
	if err != nil {
		return err
	}

}
