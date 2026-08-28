package policy

import (
	"time"
)

// CurrentSchemaVersion defines the active policy schema version
const CurrentSchemaVersion = "antimatter/v1"

// RuleType categorizes the subject of a whitelist rule
type RuleType string

const (
	// RuleTypeCommand applies to shell CLI commands
	RuleTypeCommand RuleType = "command"
	// RuleTypeTool applies to agent / MCP tool calls
	RuleTypeTool RuleType = "tool"
	// RuleTypePath applies to filesystem path access
	RuleTypePath RuleType = "path"
)

// MatchMode defines how pattern matching is evaluated
type MatchMode string

const (
	// MatchExact requires identical string match
	MatchExact MatchMode = "exact"
	// MatchPrefix requires prefix match
	MatchPrefix MatchMode = "prefix"
	// MatchGlob requires wildcard/glob match
	MatchGlob MatchMode = "glob"
	// MatchRegex requires regular expression match
	MatchRegex MatchMode = "regex"
)

// Decision defines the evaluation outcome
type Decision string

const (
	// DecisionAllow permits execution
	DecisionAllow Decision = "allow"
	// DecisionDeny blocks execution
	DecisionDeny Decision = "deny"
	// DecisionAsk prompts for interactive user approval
	DecisionAsk Decision = "ask"
)

// PolicyRule defines an individual permission rule
type PolicyRule struct {
	ID          string    `json:"id" yaml:"id"`
	Type        RuleType  `json:"type" yaml:"type"`
	Pattern     string    `json:"pattern" yaml:"pattern"`
	MatchMode   MatchMode `json:"match_mode" yaml:"match_mode"`
	Decision    Decision  `json:"decision" yaml:"decision"`
	Description string    `json:"description,omitempty" yaml:"description,omitempty"`
	Agents      []string  `json:"agents,omitempty" yaml:"agents,omitempty"` // e.g. ["antigravity", "claude", "all"]
	CreatedAt   time.Time `json:"created_at,omitempty" yaml:"created_at,omitempty"`
}

// Policy is the root container for version-controlled agent rules
type Policy struct {
	Version     string       `json:"version" yaml:"version"`
	Name        string       `json:"name" yaml:"name"`
	Description string       `json:"description,omitempty" yaml:"description,omitempty"`
	Rules       []PolicyRule `json:"rules" yaml:"rules"`
	UpdatedAt   time.Time    `json:"updated_at,omitempty" yaml:"updated_at,omitempty"`
}

// NewPolicy creates a new Policy with defaults initialized
func NewPolicy(name, description string) *Policy {
	return &Policy{
		Version:     CurrentSchemaVersion,
		Name:        name,
		Description: description,
		Rules:       make([]PolicyRule, 0),
		UpdatedAt:   time.Now().UTC(),
	}
}
