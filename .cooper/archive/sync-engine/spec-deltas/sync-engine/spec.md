# Sync Engine Spec Delta

## Added Requirements

+ ### Requirement: Git-Backed Sync Engine & Sanitization
+ The sync engine MUST manage a Git dotfiles repository, sanitize credentials from policies, merge remote rules, and push updates upstream.
+
+ #### Scenario: Synchronize Local Policy with Remote Git Store
+ - **GIVEN** a configured remote Git repository URL
+ - **WHEN** user executes `antimatter sync`
+ - **THEN** remote changes are pulled and merged into local policy, local additions are sanitized and committed, and changes are pushed upstream.
+
+ #### Scenario: Prevent Secret Leaks
+ - **GIVEN** a candidate command whitelist rule containing a detected API token (e.g. `sk-ant-12345`)
+ - **WHEN** sanitizer evaluates the policy before commit
+ - **THEN** sanitizer flags and rejects or masks the token with an error report.
+
+ #### Scenario: Inspect Remote Diff
+ - **GIVEN** local policy has uncommitted rule additions
+ - **WHEN** user executes `antimatter diff`
+ - **THEN** output displays added and modified rules compared to the remote sync HEAD.
