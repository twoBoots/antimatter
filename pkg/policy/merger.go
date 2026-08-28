package policy

import (
	"fmt"
	"time"
)

// MergePolicies merges an incoming policy into a base policy deterministically
func MergePolicies(base *Policy, incoming *Policy) (*Policy, error) {
	if base == nil && incoming == nil {
		return nil, fmt.Errorf("both base and incoming policies are nil")
	}
	if base == nil {
		return incoming, nil
	}
	if incoming == nil {
		return base, nil
	}

	result := &Policy{
		Version:     base.Version,
		Name:        base.Name,
		Description: base.Description,
		Rules:       make([]PolicyRule, 0, len(base.Rules)+len(incoming.Rules)),
		UpdatedAt:   time.Now().UTC(),
	}

	// Index base rules by ID
	ruleMap := make(map[string]PolicyRule)
	ruleOrder := make([]string, 0)

	for _, r := range base.Rules {
		ruleMap[r.ID] = r
		ruleOrder = append(ruleOrder, r.ID)
	}

	// Upsert incoming rules
	for _, r := range incoming.Rules {
		if _, exists := ruleMap[r.ID]; !exists {
			ruleOrder = append(ruleOrder, r.ID)
		}
		ruleMap[r.ID] = r
	}

	for _, id := range ruleOrder {
		result.Rules = append(result.Rules, ruleMap[id])
	}

	if err := Validate(result); err != nil {
		return nil, fmt.Errorf("merged policy validation failed: %w", err)
	}

	return result, nil
}
