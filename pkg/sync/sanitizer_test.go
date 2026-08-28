package sync_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twoBoots/antimatter/pkg/policy"
	syncpkg "github.com/twoBoots/antimatter/pkg/sync"
)

func TestSanitizer_DetectsTokens(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		hasToken bool
	}{
		{"clean command", "git status", false},
		{"clean npm command", "npm test -- --coverage", false},
		{"anthropic key", "export ANTHROPIC_API_KEY=sk-ant-api03-abcdef1234567890abcdef1234567890", true},
		{"openai key", "curl -H 'Authorization: Bearer sk-1234567890abcdef1234567890' https://api.openai.com", true},
		{"github pat", "git clone https://ghp_0123456789abcdef0123456789abcdef0123@github.com/repo", true},
		{"google api key", "curl https://generativelanguage.googleapis.com?key=AIzaSyA1234567890abcdef123456789012345", true},
		{"jwt token", "curl -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.doNotLeakThisSignature'", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.hasToken, syncpkg.ContainsSecrets(tc.input))
		})
	}
}

func TestSanitizer_SanitizeText(t *testing.T) {
	raw := "curl -H 'Authorization: Bearer sk-ant-api03-1234567890abcdef1234567890' https://api.anthropic.com"
	clean := syncpkg.SanitizeText(raw)
	assert.NotContains(t, clean, "sk-ant-api03-1234567890abcdef1234567890")
	assert.Contains(t, clean, "[REDACTED_SECRET]")
}

func TestSanitizer_SanitizePolicy(t *testing.T) {
	p := policy.NewPolicy("secret-test", "Testing Policy Sanitizer")
	p.Rules = []policy.PolicyRule{
		{
			ID:          "rule-clean",
			Type:        policy.RuleTypeCommand,
			Pattern:     "git log -n 5",
			MatchMode:   policy.MatchExact,
			Description: "Safe rule",
		},
		{
			ID:          "rule-dirty",
			Type:        policy.RuleTypeCommand,
			Pattern:     "deploy --token=ghp_1234567890abcdef1234567890abcdef1234",
			MatchMode:   policy.MatchPrefix,
			Description: "Leak with AIzaSyA1234567890abcdef123456789012345 in desc",
		},
	}

	// ValidatePolicySecrets should detect the leak
	err := syncpkg.ValidatePolicySecrets(p)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "rule-dirty")

	// SanitizePolicy should clean it
	sanitized, count := syncpkg.SanitizePolicy(p)
	assert.Equal(t, 2, count)
	assert.Equal(t, "git log -n 5", sanitized.Rules[0].Pattern)
	assert.Contains(t, sanitized.Rules[1].Pattern, "[REDACTED_SECRET]")
	assert.Contains(t, sanitized.Rules[1].Description, "[REDACTED_SECRET]")

	// Validating sanitized policy should now succeed
	err = syncpkg.ValidatePolicySecrets(sanitized)
	assert.NoError(t, err)
}
