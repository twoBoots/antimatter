package policy_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twoBoots/antimatter/pkg/policy"
)

func TestPolicy_YAML_Serialization(t *testing.T) {
	p := policy.NewPolicy("dev-workstation", "Developer Workstation Baseline")
	p.Rules = append(p.Rules, policy.PolicyRule{
		ID:          "git-status",
		Type:        policy.RuleTypeCommand,
		Pattern:     "git status",
		MatchMode:   policy.MatchExact,
		Decision:    policy.DecisionAllow,
		Description: "Allow checking git status",
		Agents:      []string{"antigravity", "claude"},
	})

	yamlBytes, err := policy.EncodeYAML(p)
	require.NoError(t, err)
	assert.Contains(t, string(yamlBytes), "git status")
	assert.Contains(t, string(yamlBytes), "dev-workstation")

	decoded, err := policy.DecodeYAML(yamlBytes)
	require.NoError(t, err)
	assert.Equal(t, p.Name, decoded.Name)
	assert.Equal(t, p.Version, decoded.Version)
	require.Len(t, decoded.Rules, 1)
	assert.Equal(t, "git-status", decoded.Rules[0].ID)
	assert.Equal(t, policy.RuleTypeCommand, decoded.Rules[0].Type)
	assert.Equal(t, policy.MatchExact, decoded.Rules[0].MatchMode)
}

func TestPolicy_JSON_Serialization(t *testing.T) {
	p := policy.NewPolicy("ci-profile", "CI Profile")
	p.Rules = append(p.Rules, policy.PolicyRule{
		ID:        "npm-test",
		Type:      policy.RuleTypeCommand,
		Pattern:   "npm test",
		MatchMode: policy.MatchPrefix,
		Decision:  policy.DecisionAllow,
	})

	jsonBytes, err := policy.EncodeJSON(p)
	require.NoError(t, err)
	assert.Contains(t, string(jsonBytes), "npm test")

	decoded, err := policy.DecodeJSON(jsonBytes)
	require.NoError(t, err)
	assert.Equal(t, "ci-profile", decoded.Name)
	require.Len(t, decoded.Rules, 1)
}

func TestPolicy_Validate_Success(t *testing.T) {
	p := policy.NewPolicy("valid-policy", "Valid")
	p.Rules = []policy.PolicyRule{
		{
			ID:        "rule-1",
			Type:      policy.RuleTypeCommand,
			Pattern:   "git *",
			MatchMode: policy.MatchGlob,
			Decision:  policy.DecisionAllow,
		},
		{
			ID:        "rule-2",
			Type:      policy.RuleTypeCommand,
			Pattern:   "^go (test|build).*$",
			MatchMode: policy.MatchRegex,
			Decision:  policy.DecisionAllow,
		},
	}

	err := policy.Validate(p)
	assert.NoError(t, err)
}

func TestPolicy_Validate_Errors(t *testing.T) {
	t.Run("nil policy", func(t *testing.T) {
		err := policy.Validate(nil)
		assert.Error(t, err)
	})

	t.Run("missing name", func(t *testing.T) {
		p := &policy.Policy{Version: policy.CurrentSchemaVersion}
		err := policy.Validate(p)
		assert.ErrorContains(t, err, "name is required")
	})

	t.Run("invalid regex pattern", func(t *testing.T) {
		p := policy.NewPolicy("test", "Test")
		p.Rules = []policy.PolicyRule{
			{
				ID:        "broken-regex",
				Type:      policy.RuleTypeCommand,
				Pattern:   "[unclosed-bracket",
				MatchMode: policy.MatchRegex,
			},
		}
		err := policy.Validate(p)
		assert.ErrorContains(t, err, "invalid regex")
	})

	t.Run("duplicate rule IDs", func(t *testing.T) {
		p := policy.NewPolicy("test", "Test")
		p.Rules = []policy.PolicyRule{
			{ID: "dup-id", Type: policy.RuleTypeCommand, Pattern: "git status", MatchMode: policy.MatchExact},
			{ID: "dup-id", Type: policy.RuleTypeCommand, Pattern: "git diff", MatchMode: policy.MatchExact},
		}
		err := policy.Validate(p)
		assert.ErrorContains(t, err, "duplicate rule ID")
	})

	t.Run("invalid rule type", func(t *testing.T) {
		p := policy.NewPolicy("test", "Test")
		p.Rules = []policy.PolicyRule{
			{ID: "r1", Type: "invalid-type", Pattern: "echo hi", MatchMode: policy.MatchExact},
		}
		err := policy.Validate(p)
		assert.ErrorContains(t, err, "invalid rule type")
	})
}

func TestPolicy_FileIO(t *testing.T) {
	tmpDir := t.TempDir()
	yamlPath := filepath.Join(tmpDir, "policy.yaml")
	jsonPath := filepath.Join(tmpDir, "policy.json")

	p := policy.NewPolicy("file-test", "File IO Test")
	p.Rules = []policy.PolicyRule{
		{
			ID:        "rule-io",
			Type:      policy.RuleTypeCommand,
			Pattern:   "ls -la",
			MatchMode: policy.MatchExact,
		},
	}

	// YAML Save and Load
	err := policy.SaveToFile(p, yamlPath)
	require.NoError(t, err)
	assert.FileExists(t, yamlPath)

	loadedYAML, err := policy.LoadFromFile(yamlPath)
	require.NoError(t, err)
	assert.Equal(t, p.Name, loadedYAML.Name)
	require.Len(t, loadedYAML.Rules, 1)

	// JSON Save and Load
	err = policy.SaveToFile(p, jsonPath)
	require.NoError(t, err)
	assert.FileExists(t, jsonPath)

	loadedJSON, err := policy.LoadFromFile(jsonPath)
	require.NoError(t, err)
	assert.Equal(t, p.Name, loadedJSON.Name)
	require.Len(t, loadedJSON.Rules, 1)

	// Non-existent file
	_, err = policy.LoadFromFile(filepath.Join(tmpDir, "non-existent.yaml"))
	assert.Error(t, err)

	// Save invalid policy
	invalidPolicy := &policy.Policy{}
	err = policy.SaveToFile(invalidPolicy, filepath.Join(tmpDir, "invalid.yaml"))
	assert.Error(t, err)
}

func TestPolicy_Codec_ErrorCases(t *testing.T) {
	t.Run("invalid YAML unmarshal", func(t *testing.T) {
		_, err := policy.DecodeYAML([]byte(":\ninvalid_yaml: ["))
		assert.Error(t, err)
	})

	t.Run("invalid JSON unmarshal", func(t *testing.T) {
		_, err := policy.DecodeJSON([]byte("{not_json"))
		assert.Error(t, err)
	})

	t.Run("encode invalid policy", func(t *testing.T) {
		_, err := policy.EncodeYAML(&policy.Policy{})
		assert.Error(t, err)

		_, err = policy.EncodeJSON(&policy.Policy{})
		assert.Error(t, err)
	})

	t.Run("defaults applied on empty match_mode and decision", func(t *testing.T) {
		p := &policy.Policy{
			Name: "test-defaults",
			Rules: []policy.PolicyRule{
				{
					ID:      "r1",
					Type:    policy.RuleTypeCommand,
					Pattern: "git status",
				},
			},
		}
		err := policy.Validate(p)
		require.NoError(t, err)
		assert.Equal(t, policy.CurrentSchemaVersion, p.Version)
		assert.Equal(t, policy.MatchExact, p.Rules[0].MatchMode)
		assert.Equal(t, policy.DecisionAllow, p.Rules[0].Decision)
	})

	t.Run("empty rule ID", func(t *testing.T) {
		p := policy.NewPolicy("test", "test")
		p.Rules = []policy.PolicyRule{
			{Type: policy.RuleTypeCommand, Pattern: "echo"},
		}
		err := policy.Validate(p)
		assert.ErrorContains(t, err, "ID is required")
	})

	t.Run("empty pattern", func(t *testing.T) {
		p := policy.NewPolicy("test", "test")
		p.Rules = []policy.PolicyRule{
			{ID: "r1", Type: policy.RuleTypeCommand, Pattern: ""},
		}
		err := policy.Validate(p)
		assert.ErrorContains(t, err, "pattern cannot be empty")
	})

	t.Run("invalid decision", func(t *testing.T) {
		p := policy.NewPolicy("test", "test")
		p.Rules = []policy.PolicyRule{
			{ID: "r1", Type: policy.RuleTypeCommand, Pattern: "echo", Decision: "unknown"},
		}
		err := policy.Validate(p)
		assert.ErrorContains(t, err, "invalid decision")
	})

	t.Run("invalid match mode", func(t *testing.T) {
		p := policy.NewPolicy("test", "test")
		p.Rules = []policy.PolicyRule{
			{ID: "r1", Type: policy.RuleTypeCommand, Pattern: "echo", MatchMode: "fuzzy"},
		}
		err := policy.Validate(p)
		assert.ErrorContains(t, err, "invalid match mode")
	})
}
