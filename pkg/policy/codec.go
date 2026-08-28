package policy

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"gopkg.in/yaml.v3"
)

// EncodeYAML encodes a Policy into formatted YAML
func EncodeYAML(p *Policy) ([]byte, error) {
	if err := Validate(p); err != nil {
		return nil, fmt.Errorf("invalid policy: %w", err)
	}
	return yaml.Marshal(p)
}

// DecodeYAML decodes and validates YAML bytes into a Policy
func DecodeYAML(data []byte) (*Policy, error) {
	var p Policy
	if err := yaml.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("failed to parse policy YAML: %w", err)
	}
	if err := Validate(&p); err != nil {
		return nil, fmt.Errorf("policy validation failed: %w", err)
	}
	return &p, nil
}

// EncodeJSON encodes a Policy into indented JSON
func EncodeJSON(p *Policy) ([]byte, error) {
	if err := Validate(p); err != nil {
		return nil, fmt.Errorf("invalid policy: %w", err)
	}
	return json.MarshalIndent(p, "", "  ")
}

// DecodeJSON decodes and validates JSON bytes into a Policy
func DecodeJSON(data []byte) (*Policy, error) {
	var p Policy
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("failed to parse policy JSON: %w", err)
	}
	if err := Validate(&p); err != nil {
		return nil, fmt.Errorf("policy validation failed: %w", err)
	}
	return &p, nil
}

// Validate performs structural and syntactic validation on a Policy
func Validate(p *Policy) error {
	if p == nil {
		return errors.New("policy cannot be nil")
	}
	if p.Version == "" {
		p.Version = CurrentSchemaVersion
	}
	if p.Name == "" {
		return errors.New("policy name is required")
	}

	seenIDs := make(map[string]bool)

	for i, rule := range p.Rules {
		if rule.ID == "" {
			return fmt.Errorf("rule[%d]: ID is required", i)
		}
		if seenIDs[rule.ID] {
			return fmt.Errorf("rule[%d]: duplicate rule ID '%s'", i, rule.ID)
		}
		seenIDs[rule.ID] = true

		switch rule.Type {
		case RuleTypeCommand, RuleTypeTool, RuleTypePath:
		default:
			return fmt.Errorf("rule[%d] '%s': invalid rule type '%s'", i, rule.ID, rule.Type)
		}

		if rule.Pattern == "" {
			return fmt.Errorf("rule[%d] '%s': pattern cannot be empty", i, rule.ID)
		}

		switch rule.MatchMode {
		case MatchExact, MatchPrefix, MatchGlob:
			// valid modes
		case MatchRegex:
			if _, err := regexp.Compile(rule.Pattern); err != nil {
				return fmt.Errorf("rule[%d] '%s': invalid regex pattern '%s': %w", i, rule.ID, rule.Pattern, err)
			}
		case "":
			p.Rules[i].MatchMode = MatchExact
		default:
			return fmt.Errorf("rule[%d] '%s': invalid match mode '%s'", i, rule.ID, rule.MatchMode)
		}

		switch rule.Decision {
		case DecisionAllow, DecisionDeny, DecisionAsk:
		case "":
			p.Rules[i].Decision = DecisionAllow
		default:
			return fmt.Errorf("rule[%d] '%s': invalid decision '%s'", i, rule.ID, rule.Decision)
		}
	}

	return nil
}

// LoadFromFile reads and decodes a policy from a YAML or JSON file
func LoadFromFile(filePath string) (*Policy, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read policy file '%s': %w", filePath, err)
	}

	ext := filepath.Ext(filePath)
	if ext == ".json" {
		return DecodeJSON(data)
	}
	return DecodeYAML(data)
}

// SaveToFile writes an encoded policy to disk atomically
func SaveToFile(p *Policy, filePath string) error {
	if err := Validate(p); err != nil {
		return fmt.Errorf("cannot save invalid policy: %w", err)
	}

	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory '%s': %w", dir, err)
	}

	var data []byte
	var err error
	if filepath.Ext(filePath) == ".json" {
		data, err = EncodeJSON(p)
	} else {
		data, err = EncodeYAML(p)
	}
	if err != nil {
		return err
	}

	// Write atomically via temporary file
	tmpFile := filePath + ".tmp"
	if err := os.WriteFile(tmpFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write temporary file '%s': %w", tmpFile, err)
	}

	if err := os.Rename(tmpFile, filePath); err != nil {
		_ = os.Remove(tmpFile)
		return fmt.Errorf("failed to atomically rename '%s' to '%s': %w", tmpFile, filePath, err)
	}

	return nil
}
