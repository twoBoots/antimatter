package adapters

import (
	"fmt"
	"sync"

	"github.com/twoBoots/antimatter/pkg/policy"
)

// Manager coordinates registered agent adapters
type Manager struct {
	mu       sync.RWMutex
	adapters map[string]Adapter
	order    []string
}

// NewManager creates a new empty adapter Manager
func NewManager() *Manager {
	return &Manager{
		adapters: make(map[string]Adapter),
		order:    make([]string, 0),
	}
}

// Register registers an adapter
func (m *Manager) Register(a Adapter) {
	m.mu.Lock()
	defer m.mu.Unlock()

	name := a.Name()
	if _, exists := m.adapters[name]; !exists {
		m.order = append(m.order, name)
	}
	m.adapters[name] = a
}

// Get returns an adapter by name
func (m *Manager) Get(name string) (Adapter, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	a, ok := m.adapters[name]
	return a, ok
}

// List returns all registered adapters in registration order
func (m *Manager) List() []Adapter {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]Adapter, 0, len(m.order))
	for _, name := range m.order {
		result = append(result, m.adapters[name])
	}
	return result
}

// ExportAll exports the canonical policy to all registered and detected adapters
func (m *Manager) ExportAll(p *policy.Policy) map[string]error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	errors := make(map[string]error)
	for _, name := range m.order {
		a := m.adapters[name]
		if err := a.Export(p); err != nil {
			errors[name] = err
		}
	}
	return errors
}

// ImportFrom imports rules from a specific agent adapter
func (m *Manager) ImportFrom(name string) (*policy.Policy, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	a, ok := m.adapters[name]
	if !ok {
		return nil, fmt.Errorf("adapter '%s' not registered", name)
	}

	return a.Import()
}

// GetStatuses returns inspection statuses for all registered adapters
func (m *Manager) GetStatuses() []AdapterStatus {
	m.mu.RLock()
	defer m.mu.RUnlock()

	statuses := make([]AdapterStatus, 0, len(m.order))
	for _, name := range m.order {
		a := m.adapters[name]
		status := AdapterStatus{
			Name:        a.Name(),
			DisplayName: a.DisplayName(),
			Detected:    a.Detected(),
			ConfigPath:  a.ConfigPath(),
		}

		backups, err := ListBackups(name)
		if err == nil && len(backups) > 0 {
			status.LastBackup = backups[0]
		}

		statuses = append(statuses, status)
	}
	return statuses
}
