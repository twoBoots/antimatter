package policy

import (
	"path/filepath"
	"regexp"
	"strings"
)

// EvaluationResult represents the outcome of evaluating a candidate action
type EvaluationResult struct {
	Decision    Decision    `json:"decision"`
	MatchedRule *PolicyRule `json:"matched_rule,omitempty"`
	Reason      string      `json:"reason"`
}

// EvaluateCommand checks whether a candidate CLI command is allowed, denied, or requires prompt
func EvaluateCommand(p *Policy, agent string, cmd string) EvaluationResult {
	if p == nil {
		return EvaluationResult{Decision: DecisionAsk, Reason: "no active policy"}
	}

	trimmedCmd := strings.TrimSpace(cmd)
	if trimmedCmd == "" {
		return EvaluationResult{Decision: DecisionAsk, Reason: "empty command"}
	}

	for i := range p.Rules {
		rule := &p.Rules[i]
		if rule.Type != RuleTypeCommand {
			continue
		}

		if !isAgentApplicable(rule.Agents, agent) {
			continue
		}

		if matchesPattern(rule.MatchMode, rule.Pattern, trimmedCmd) {
			return EvaluationResult{
				Decision:    rule.Decision,
				MatchedRule: rule,
				Reason:      "matched command rule: " + rule.ID,
			}
		}
	}

	return EvaluationResult{
		Decision: DecisionAsk,
		Reason:   "no matching whitelist rule found",
	}
}

// EvaluateTool checks whether a tool invocation (e.g. MCP tool or native tool) is permitted
func EvaluateTool(p *Policy, agent string, toolName string) EvaluationResult {
	if p == nil {
		return EvaluationResult{Decision: DecisionAsk, Reason: "no active policy"}
	}

	trimmedTool := strings.TrimSpace(toolName)
	if trimmedTool == "" {
		return EvaluationResult{Decision: DecisionAsk, Reason: "empty tool name"}
	}

	for i := range p.Rules {
		rule := &p.Rules[i]
		if rule.Type != RuleTypeTool {
			continue
		}

		if !isAgentApplicable(rule.Agents, agent) {
			continue
		}

		if matchesPattern(rule.MatchMode, rule.Pattern, trimmedTool) {
			return EvaluationResult{
				Decision:    rule.Decision,
				MatchedRule: rule,
				Reason:      "matched tool rule: " + rule.ID,
			}
		}
	}

	return EvaluationResult{
		Decision: DecisionAsk,
		Reason:   "no matching tool rule found",
	}
}

func isAgentApplicable(ruleAgents []string, candidateAgent string) bool {
	if len(ruleAgents) == 0 {
		return true
	}
	candidateAgent = strings.ToLower(strings.TrimSpace(candidateAgent))
	for _, a := range ruleAgents {
		a = strings.ToLower(strings.TrimSpace(a))
		if a == "all" || a == "*" || a == candidateAgent {
			return true
		}
	}
	return false
}

func matchesPattern(mode MatchMode, pattern string, candidate string) bool {
	switch mode {
	case MatchExact:
		return candidate == strings.TrimSpace(pattern)
	case MatchPrefix:
		prefix := strings.TrimSpace(pattern)
		return strings.HasPrefix(candidate, prefix)
	case MatchGlob:
		// Convert shell glob pattern (* and ?) into regex for full command line matching
		regexPattern := globToRegex(pattern)
		re, err := regexp.Compile(regexPattern)
		if err != nil {
			// Fallback to filepath.Match
			matched, mErr := filepath.Match(pattern, candidate)
			return mErr == nil && matched
		}
		return re.MatchString(candidate)
	case MatchRegex:
		re, err := regexp.Compile(pattern)
		if err != nil {
			return false
		}
		return re.MatchString(candidate)
	default:
		return candidate == pattern
	}
}

func globToRegex(glob string) string {
	var b strings.Builder
	b.WriteString("^")
	for i := 0; i < len(glob); i++ {
		ch := glob[i]
		switch ch {
		case '*':
			b.WriteString(".*")
		case '?':
			b.WriteString(".")
		case '.', '+', '(', ')', '[', ']', '{', '}', '^', '$', '|', '\\':
			b.WriteString("\\")
			b.WriteByte(ch)
		default:
			b.WriteByte(ch)
		}
	}
	b.WriteString("$")
	return b.String()
}
