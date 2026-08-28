package adapters

import "github.com/twoBoots/antimatter/pkg/policy"

// Adapter defines the contract for synchronizing with an AI agent environment
type Adapter interface {
	Name() string
	DisplayName() string
	Detected() bool
	ConfigPath() string
	Backup() (string, error)
	Export(p *policy.Policy) error
	Import() (*policy.Policy, error)
}

// AdapterStatus represents the inspection metadata for a configured adapter
type AdapterStatus struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Detected    bool   `json:"detected"`
	ConfigPath  string `json:"configPath"`
	RuleCount   int    `json:"ruleCount,omitempty"`
	LastBackup  string `json:"lastBackup,omitempty"`
	Error       string `json:"error,omitempty"`
}
