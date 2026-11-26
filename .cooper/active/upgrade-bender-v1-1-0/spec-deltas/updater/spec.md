# Self-Updater Capability Spec Delta

## Added Requirements

+ ### Requirement: SemVer 2.0.0 Standardized Comparison
+ The self-updater version comparison MUST adhere to SemVer 2.0.0 specification via `golang.org/x/mod/semver`, supporting pre-release tags, build metadata, and partial version strings.
+
+ #### Scenario: Evaluate Pre-Release and Partial Target Versions
+ - **GIVEN** a running binary at version `v1.0.0`
+ - **WHEN** user checks for updates against candidate pre-release or partial versions
+ - **THEN** updater correctly orders versions according to SemVer 2.0.0 precedence.
