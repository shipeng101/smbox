package service

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/yourusername/singdns/internal/constant"
)

type ServiceManager struct {
	CorePath string
	services map[string]*Service
}

type Service struct {
	Name    string
	Status  string
	Process *os.Process
}

func NewServiceManager(corePath string) *ServiceManager {
	return &ServiceManager{
		CorePath: corePath,
		services: make(map[string]*Service),
	}
}

func (sm *ServiceManager) StartService(name string) error {
	execPath := filepath.Join(sm.CorePath, name)
	if _, err := os.Stat(execPath); os.IsNotExist(err) {
		return fmt.Errorf("service executable not found: %s", execPath)
	}

	cmd := exec.Command(execPath)
	if err := cmd.Start(); err != nil {
		return err
	}

	sm.services[name] = &Service{
		Name:    name,
		Status:  constant.StatusRunning,
		Process: cmd.Process,
	}

	return nil
}

func (sm *ServiceManager) StopService(name string) error {
	if service, ok := sm.services[name]; ok {
		if err := service.Process.Kill(); err != nil {
			return err
		}
		service.Status = constant.StatusStopped
		delete(sm.services, name)
	}

	// 确保进程被终止
	cmd := exec.Command("pkill", "-f", name)
	return cmd.Run()
}

func (sm *ServiceManager) GetServiceStatus(name string) string {
	if service, ok := sm.services[name]; ok {
		return service.Status
	}
	return constant.StatusStopped
}
