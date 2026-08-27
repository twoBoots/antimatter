# Policy Engine Spec Delta

## Added Requirements

+ ### Requirement: Policy Engine Domain Core & Matching
+ The policy engine MUST parse canonical YAML policies and evaluate match rules against candidate commands and tool calls.
+
+ #### Scenario: Exact & Prefix Command Matching
+ - **GIVEN** a policy with command prefix rule `npm test`
+ - **WHEN** candidate `npm test -- --watch=false` is evaluated
+ - **THEN** evaluator returns `Decision: Allow` with matched rule ID.
+
+ #### Scenario: Regex Command Matching & Validation
+ - **GIVEN** a policy with regex rule `^git (status|diff|log.*)$`
+ - **WHEN** candidate `git status` is evaluated
+ - **THEN** evaluator returns `Decision: Allow`.
+
+ #### Scenario: Unsafe Regex Rejection
+ - **GIVEN** an invalid regex string `*git (`
+ - **WHEN** policy engine attempts to parse or validate the rule
+ - **THEN** validation fails with a syntax error.
