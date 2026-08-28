# Adapters Spec Delta

## Added Requirements

+ ### Requirement: Bidirectional Antigravity Synchronization
+ The Antigravity adapter MUST export canonical policies into Antigravity target stores and import local whitelisted rules.
+
+ #### Scenario: Export to Antigravity
+ - **GIVEN** a canonical policy with command and tool rules
+ - **WHEN** adapter export is executed
+ - **THEN** a backup is created in `~/.antimatter/backups/antigravity/` and settings are safely written.
+
+ #### Scenario: Import from Antigravity
+ - **GIVEN** existing Antigravity settings containing approved tool calls
+ - **WHEN** adapter import is executed
+ - **THEN** valid rules are parsed and converted into canonical `PolicyRule` records.
+
+ ### Requirement: Bidirectional Claude Synchronization
+ The Claude adapter MUST export canonical policies to Claude Code configurations and import approved tools/commands.
+
+ #### Scenario: Export to Claude Code
+ - **GIVEN** a canonical policy
+ - **WHEN** Claude adapter export is executed
+ - **THEN** `allowedTools` and command approvals are updated with prior backup creation.
