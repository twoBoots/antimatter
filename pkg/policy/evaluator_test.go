package policy_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twoBoots/antimatter/pkg/policy"
)

func TestEvaluator_EvaluateCommand(t *testing.T) {
	p := policy.NewPolicy("test-eval", "Evaluator Test")
	p.Rules = []policy.PolicyRule{
		{
			ID:        "exact-git-status",
			Type:      policy.RuleTypeCommand,
			Pattern:   "git status",
			MatchMode: policy.MatchExact,
			Decision:  policy.DecisionAllow,
			Agents:    []string{"antigravity", "claude"},
		},
		{
			ID:        "prefix-npm-test",
			Type:      policy.RuleTypeCommand,
			Pattern:   "npm test",
			MatchMode: policy.MatchPrefix,
			Decision:  policy.DecisionAllow,
			Agents:    []string{"all"},
		},
		{
			ID:        "glob-go-test",
			Type:      policy.RuleTypeCommand,
			Pattern:   "go test *",
			MatchMode: policy.MatchGlob,
			Decision:  policy.DecisionAllow,
		},
		{
			ID:        "regex-git-diff",
			Type:      policy.RuleTypeCommand,
			Pattern:   `^git diff(\s+.*)?$`,
			MatchMode: policy.MatchRegex,
			Decision:  policy.DecisionAllow,
		},
		{
			ID:        "deny-rm-rf",
			Type:      policy.RuleTypeCommand,
			Pattern:   `rm -rf /.*`,
			MatchMode: policy.MatchRegex,
			Decision:  policy.DecisionDeny,
		},
	}

	t.Run("exact match - allowed", func(t *testing.T) {
		res := policy.EvaluateCommand(p, "antigravity", "git status")
		assert.Equal(t, policy.DecisionAllow, res.Decision)
		require.NotNil(t, res.MatchedRule)
		assert.Equal(t, "exact-git-status", res.MatchedRule.ID)
	})

	t.Run("exact match - agent mismatch", func(t *testing.T) {
		res := policy.EvaluateCommand(p, "other-agent", "git status")
		assert.Equal(t, policy.DecisionAsk, res.Decision)
		assert.Nil(t, res.MatchedRule)
	})

	t.Run("prefix match - allowed", func(t *testing.T) {
		res := policy.EvaluateCommand(p, "claude", "npm test -- --coverage")
		assert.Equal(t, policy.DecisionAllow, res.Decision)
		require.NotNil(t, res.MatchedRule)
		assert.Equal(t, "prefix-npm-test", res.MatchedRule.ID)
	})

	t.Run("glob match - allowed", func(t *testing.T) {
		res := policy.EvaluateCommand(p, "antigravity", "go test ./...")
		assert.Equal(t, policy.DecisionAllow, res.Decision)
		require.NotNil(t, res.MatchedRule)
		assert.Equal(t, "glob-go-test", res.MatchedRule.ID)
	})

	t.Run("regex match - allowed", func(t *testing.T) {
		res := policy.EvaluateCommand(p, "claude", "git diff --stat HEAD~1")
		assert.Equal(t, policy.DecisionAllow, res.Decision)
		require.NotNil(t, res.MatchedRule)
		assert.Equal(t, "regex-git-diff", res.MatchedRule.ID)
	})

	t.Run("explicit deny rule", func(t *testing.T) {
		res := policy.EvaluateCommand(p, "antigravity", "rm -rf /")
		assert.Equal(t, policy.DecisionDeny, res.Decision)
		require.NotNil(t, res.MatchedRule)
		assert.Equal(t, "deny-rm-rf", res.MatchedRule.ID)
	})

	t.Run("no match - default ask", func(t *testing.T) {
		res := policy.EvaluateCommand(p, "antigravity", "curl -X POST https://example.com")
		assert.Equal(t, policy.DecisionAsk, res.Decision)
		assert.Nil(t, res.MatchedRule)
	})
}

func TestEvaluator_EvaluateTool(t *testing.T) {
	p := policy.NewPolicy("tool-eval", "Tool Evaluator Test")
	p.Rules = []policy.PolicyRule{
		{
			ID:        "allow-view-file",
			Type:      policy.RuleTypeTool,
			Pattern:   "view_file",
			MatchMode: policy.MatchExact,
			Decision:  policy.DecisionAllow,
		},
		{
			ID:        "allow-mcp-tools",
			Type:      policy.RuleTypeTool,
			Pattern:   "mcp_*",
			MatchMode: policy.MatchGlob,
			Decision:  policy.DecisionAllow,
		},
	}

	t.Run("exact tool match", func(t *testing.T) {
		res := policy.EvaluateTool(p, "antigravity", "view_file")
		assert.Equal(t, policy.DecisionAllow, res.Decision)
		require.NotNil(t, res.MatchedRule)
		assert.Equal(t, "allow-view-file", res.MatchedRule.ID)
	})

	t.Run("glob tool match", func(t *testing.T) {
		res := policy.EvaluateTool(p, "claude", "mcp_atlassian_getJiraIssue")
		assert.Equal(t, policy.DecisionAllow, res.Decision)
		require.NotNil(t, res.MatchedRule)
		assert.Equal(t, "allow-mcp-tools", res.MatchedRule.ID)
	})

	t.Run("unknown tool default ask", func(t *testing.T) {
		res := policy.EvaluateTool(p, "claude", "run_dangerous_command")
		assert.Equal(t, policy.DecisionAsk, res.Decision)
		assert.Nil(t, res.MatchedRule)
	})
}

func TestMerger_MergePolicies(t *testing.T) {
	base := policy.NewPolicy("base-profile", "Base")
	base.Rules = []policy.PolicyRule{
		{ID: "r1", Type: policy.RuleTypeCommand, Pattern: "git status", MatchMode: policy.MatchExact},
		{ID: "r2", Type: policy.RuleTypeCommand, Pattern: "git diff", MatchMode: policy.MatchExact},
	}

	incoming := policy.NewPolicy("incoming-profile", "Incoming")
	incoming.Rules = []policy.PolicyRule{
		{ID: "r2", Type: policy.RuleTypeCommand, Pattern: "git diff --stat", MatchMode: policy.MatchPrefix}, // update r2
		{ID: "r3", Type: policy.RuleTypeCommand, Pattern: "npm test", MatchMode: policy.MatchPrefix},
	}

	merged, err := policy.MergePolicies(base, incoming)
	require.NoError(t, err)
	assert.Equal(t, "base-profile", merged.Name)
	require.Len(t, merged.Rules, 3)

	// Verify r2 is updated from incoming
	var r2 *policy.PolicyRule
	for i := range merged.Rules {
		if merged.Rules[i].ID == "r2" {
			r2 = &merged.Rules[i]
		}
	}
	require.NotNil(t, r2)
	assert.Equal(t, "git diff --stat", r2.Pattern)
	assert.Equal(t, policy.MatchPrefix, r2.MatchMode)
}
