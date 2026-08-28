# CLI Spec Delta

## Added Requirements

+ ### Requirement: Baseline CLI Binary Structure
+ The CLI MUST provide root and core subcommands with standard `--help` and `-v` flags.
+
+ #### Scenario: Inspect Version
+ - **GIVEN** compiled `antimatter` binary
+ - **WHEN** user executes `antimatter version` or `antimatter -v`
+ - **THEN** output displays version, commit SHA, and build timestamp.
+
+ #### Scenario: Check Command Whitelist Match
+ - **GIVEN** an active policy loaded from config
+ - **WHEN** user executes `antimatter whitelist check "git status"`
+ - **THEN** output indicates whether the command is ALLOWED or DENIED with matching rule details.
