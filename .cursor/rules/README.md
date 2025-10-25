# Cursor Rules — Bemax Backend

This directory contains Cursor rules to standardize PRs and assist with code reviews.

## 🚀 Quick Start - Available Shortcuts

### 📝 `!pr` - Automatic PR Generation
Generates a complete Pull Request description based on your current changes.

**What it does:**
- Automatically analyzes workspace changes
- Detects change type from branch prefix (feature/, hotfix/, migration/)
- Detects scope based on modified files
- Generates standardized description ready to copy/paste
- Includes compliance analysis against project standards
- Optional Mermaid diagram (request explicitly)

### 🔍 `!review` - Automated Code Review
Performs comprehensive code review analysis against project standards.

**What it analyzes:**
- ✅ Clean Architecture compliance
- 📝 Go standards (goimports, naming, error handling)
- 📚 English-only GoDoc documentation
- 💎 Code quality (SOLID, DRY, patterns)
- 🧪 Testing standards

**Output:**
- Compliance score (0-100)
- Detailed report with specific files and line numbers
- Improvement suggestions
- Priority actions

## 📋 What the rules do
- **Standardization**: Consistent PR titles and descriptions
- **Automation**: Automatic PR text generation when requested
- **Quality**: Review guidance with project-based checklist
- **Standards**: Consolidation of architecture/Go/DDD specific to this repository
- **Compliance**: Automatic verification of code standards

## 🔧 Detailed Functionality

### 📝 `!pr` Command Deep Dive
The `!pr` command triggers the main PR generation system powered by [pr-rules.mdc](mdc:.cursor/rules/pr-rules.mdc).

**How it works:**
1. **Branch Analysis**: Detects PR type from branch prefix:
- `feature/` → `feature(feature): description`
- `hotfix/` → `hotfix(hotfix): description`
- `migration/` → `migration(migration): description`

2. **Scope Detection**: Maps modified files to scopes:
```
   internal/adapters/auth/ → auth
   internal/adapters/handlers/agent* → agent
   internal/core/domain/ → domain
   migrations/ → migration
   ```

3. **Template Generation**: Creates standardized PR description with:
- Motivation and context
- Compliance analysis against project standards
- Technical details and architecture impact
- Testing checklist
- Deployment notes

4. **Alternative Triggers**: Also responds to keywords like "PR text", "pull request description" (multilingual support)

### 🔍 `!review` Command Deep Dive
The `!review` command triggers comprehensive code analysis powered by [code-review-analysis.mdc](mdc:.cursor/rules/code-review-analysis.mdc).

**Analysis Process:**
1. **File Discovery**: Identifies all modified files in current branch
2. **Standards Validation**: Checks against:
- [project-standards.mdc](mdc:.cursor/rules/project-standards.mdc) - Architecture & Go standards
- [godoc-english-only.mdc](mdc:.cursor/rules/godoc-english-only.mdc) - Documentation language
3. **Quality Assessment**: Evaluates:
- Clean Architecture layer separation
- Go coding standards (goimports, naming, error handling)
- Documentation completeness and language
- SOLID principles adherence
- Testing patterns
4. **Report Generation**: Produces detailed markdown report with compliance score

## 📁 Rule Files and Functions

### 🎯 Always Applied Rules
#### [pr-rules.mdc](mdc:.cursor/rules/pr-rules.mdc)
**Description:** Main rule for PR standardization and automatic generation
- **Always applied:** ✅ Yes
- **Functions:**
- Standardized template for PR descriptions
- Automatic type detection via branch prefix
- Scope mapping based on file paths
- Compliance analysis with project standards
- Multilingual support (Portuguese, English, Spanish)
- Optional Mermaid diagram integration

### 🛠️ Fetchable Rules (Loaded on Demand)
#### [project-standards.mdc](mdc:.cursor/rules/project-standards.mdc)
**Description:** Project-specific standards and best practices
- **Always applied:** ❌ No (fetchable)
- **Scope:** `**/*.go`, `**/*.sql`, `**/*.md`
- **Content:**
- Detailed Clean Architecture guidelines
- Project-specific Go standards
- Layer dependency rules
- Structured logging with go-core
- Testing and mocking standards
- Naming conventions and directory structure

#### [godoc-english-only.mdc](mdc:.cursor/rules/godoc-english-only.mdc)
**Description:** Language standard for GoDoc documentation
- **Always applied:** ❌ No (fetchable)
- **Scope:** `*.go`
- **Rule:** All GoDoc documentation MUST be written in English only
- **Applies to:**
- Exported function comments
- Exported type comments
- Exported constants and variables comments
- Package comments

#### [pr-review-guidelines.mdc](mdc:.cursor/rules/pr-review-guidelines.mdc)
**Description:** Checklist and guidelines for PR reviews
- **Always applied:** ❌ No (fetchable)
- **Functions:**
- Architecture and design checklist
- Go code quality validation
- Domain rules verification
- HTTP API standards
- Security and performance checks
- Branch-specific constraints

#### [code-review-analysis.mdc](mdc:.cursor/rules/code-review-analysis.mdc)
**Description:** Automated code review analysis tool
- **Always applied:** ❌ No (fetchable)
- **Trigger:** `!review`
- **Functions:**
- Architecture compliance analysis
- Go standards verification
- Documentation quality assessment
- Compliance score (0-100)
- Detailed report with specific suggestions

## 🏗️ How the Rules Work

### 1. **Always Applied Rules (alwaysApply: true)**
- Loaded automatically in all sessions
- Define default Cursor behaviors
- Currently: only `pr-rules.mdc`

### 2. **Fetchable Rules (alwaysApply: false)**
- Loaded only when needed
- Save resources and improve performance
- Activated by specific triggers or commands

### 3. **Trigger System**
#### PR Triggers:
- Commands: `!pr`, `!PR`, `!pr text`
- Keywords: "PR text", "pull request description", etc.
- Action: Loads `pr-rules.mdc` + analyzes workspace

#### Review Triggers:
- Command: `!review`
- Action: Loads `code-review-analysis.mdc` + executes complete analysis

### 4. **Automatic Detection**
#### PR Type (via branch prefix):
```
feature/ → feature(feature): description
hotfix/  → hotfix(hotfix): description
migration/ → migration(migration): description
```

#### Scope (via modified paths):
```
internal/adapters/auth/ → auth
internal/adapters/handlers/agent* → agent
internal/core/domain/ → domain
migrations/ → migration
```

### 5. **Compliance and Validation**
- Automatic analysis against `project-standards.mdc`
- GoDoc language verification via `godoc-english-only.mdc`
- Branch-specific constraints:
- `migration/*`: only files in `migrations/` allowed
- `feature/*` and `hotfix/*`: no specific restrictions

## 📚 Repository References
- [README.md](mdc:README.md) - Main documentation
- [CODING_GUIDELINES.md](mdc:CODING_GUIDELINES.md) - Coding guidelines
- [CONTRIBUTING.md](mdc:CONTRIBUTING.md) - Contribution guide

## 🔧 Maintenance
- **Scopes/types:** Update mappings in `pr-rules.mdc` if new directories appear
- **Checklists:** Adjust lists in `pr-review-guidelines.mdc` as the team evolves
- **Technical standards:** Keep `project-standards.mdc` updated with current practices
- **Triggers:** Add new commands in `code-review-analysis.mdc` as needed