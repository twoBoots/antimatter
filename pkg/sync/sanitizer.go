package sync

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/twoBoots/antimatter/pkg/policy"
)

const RedactedSecretPlaceholder = "[REDACTED_SECRET]"

var secretPatterns = []*regexp.Regexp{
	// Anthropic API Keys (sk-ant-...)
	regexp.MustCompile(`sk-ant-[a-zA-Z0-9_\-]{20,}`),
	// OpenAI API Keys (sk-...)
	regexp.MustCompile(`\bsk-[a-zA-Z0-9]{20,}\b`),
	// GitHub Tokens (ghp_..., github_pat_...)
	regexp.MustCompile(`\bghp_[a-zA-Z0-9]{36}\b`),
	regexp.MustCompile(`\bgithub_pat_[a-zA-Z0-9_]{22,}\b`),
	// Google API Keys (AIzaSy...)
	regexp.MustCompile(`\bAIza[0-9A-Za-z\-_]{30,}\b`),
	// JWT Tokens (eyJ...)
	regexp.MustCompile(`eyJ[a-zA-Z0-9_\-]+\.eyJ[a-zA-Z0-9_\-]+\.[a-zA-Z0-9_\-]+`),
}

// ContainsSecrets checks if a string contains any known API token or sensitive key patterns
func ContainsSecrets(text string) bool {
	for _, re := range secretPatterns {
		if re.MatchString(text) {
			return true
		}
	}
	return false
}

// SanitizeText replaces sensitive key patterns with a redaction placeholder
func SanitizeText(text string) string {
	result := text
	for _, re := range secretPatterns {
		result = re.ReplaceAllString(result, RedactedSecretPlaceholder)
	}
	return result
}

// ValidatePolicySecrets checks if any rule in the policy contains unredacted secrets
func ValidatePolicySecrets(p *policy.Policy) error {
	if p == nil {
		return nil
	}

	var leakedRules []string
	for _, r := range p.Rules {
		if ContainsSecrets(r.Pattern) || ContainsSecrets(r.Description) {
			leakedRules = append(leakedRules, r.ID)
		}
	}

	if len(leakedRules) > 0 {
		return fmt.Errorf("sensitive tokens or API keys detected in rules: %s", strings.Join(leakedRules, ", "))
	}
	return nil
}

// SanitizePolicy produces a clean copy of the policy with all secrets redacted
func SanitizePolicy(p *policy.Policy) (*policy.Policy, int) {
	if p == nil {
		return nil, 0
	}

	clean := &policy.Policy{
		Version:     p.Version,
		Name:        p.Name,
		Description: SanitizeText(p.Description),
		Rules:       make([]policy.PolicyRule, len(p.Rules)),
		UpdatedAt:   p.UpdatedAt,
	}

	redactionCount := 0
	if ContainsSecrets(p.Description) {
		redactionCount++
	}

	for i, r := range p.Rules {
		cleanRule := r
		if ContainsSecrets(r.Pattern) {
			cleanRule.Pattern = SanitizeText(r.Pattern)
			redactionCount++
		}
		if ContainsSecrets(r.Description) {
			cleanRule.Description = SanitizeText(r.Description)
			redactionCount++
		}
		clean.Rules[i] = cleanRule
	}

	return clean, redactionCount
}
