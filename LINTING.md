# Enhanced Linting & Pre-commit Hooks

This document describes the comprehensive linting system and pre-commit hooks implemented for the TimeSeriesDB project to ensure code quality, consistency, and security.

## 🚀 Overview

The enhanced linting system provides multiple layers of code quality assurance:

- **Pre-commit Hooks**: Automatic checks before each commit
- **Comprehensive Linting**: Multiple linters with extensive rules
- **Security Scanning**: Vulnerability detection and prevention
- **Code Quality**: Complexity, duplication, and style enforcement
- **Documentation**: Automated documentation and spelling checks

## 🛠️ Tools & Linters

### Core Go Linters

| Tool | Purpose | Configuration |
|------|---------|---------------|
| **golangci-lint** | Comprehensive static analysis | `.golangci.yml` |
| **gosec** | Security vulnerability scanning | Built-in rules |
| **gocyclo** | Code complexity analysis | Max complexity: 15 |
| **misspell** | Spelling mistake detection | Error-level reporting |
| **goimports** | Import formatting | Auto-fix enabled |

### Additional Quality Tools

| Tool | Purpose | Usage |
|------|---------|-------|
| **pre-commit** | Git hooks management | `.pre-commit-config.yaml` |
| **commitlint** | Commit message validation | `.commitlintrc.js` |
| **shellcheck** | Shell script validation | Warning-level checks |
| **markdownlint** | Markdown formatting | Auto-fix enabled |

## 📋 Pre-commit Hooks

### Installation

```bash
# Install pre-commit
pip install pre-commit

# Install hooks
make install-hooks

# Or manually
pre-commit install --install-hooks
```

### Available Hooks

#### Go-specific Hooks
- **go-fmt**: Code formatting
- **go-imports**: Import organization
- **go-vet**: Common Go mistakes
- **go-build**: Compilation check
- **go-test**: Test execution
- **go-test-race**: Race condition detection
- **go-cover**: Coverage reporting
- **go-mod-tidy**: Module cleanup
- **golangci-lint**: Static analysis
- **go-critic**: Code criticism
- **go-cyclo**: Complexity check
- **go-errcheck**: Error handling
- **go-gosec**: Security scan
- **go-misspell**: Spelling check
- **go-simple**: Code simplification
- **go-staticcheck**: Static analysis
- **go-unused**: Unused code detection

#### General Quality Hooks
- **trailing-whitespace**: Remove trailing spaces
- **end-of-file-fixer**: Ensure files end with newline
- **check-yaml**: YAML validation
- **check-json**: JSON validation
- **check-added-large-files**: Prevent large file commits
- **check-merge-conflict**: Detect merge conflicts
- **check-case-conflict**: Case sensitivity issues
- **debug-statements**: Remove debug code
- **name-tests-test**: Test file naming

#### Formatting Hooks
- **prettier**: Code formatting (YAML, JSON, Markdown)
- **markdownlint**: Markdown validation
- **shellcheck**: Shell script validation
- **yapf**: Python formatting

#### Security Hooks
- **detect-secrets**: Secret detection
- **license-eye**: License header validation
- **hadolint**: Dockerfile validation

#### Custom Local Hooks
- **go-test-coverage**: Coverage threshold enforcement
- **go-benchmark-check**: Performance regression detection
- **security-scan**: Security vulnerability scan
- **dependency-check**: Outdated dependency detection
- **code-complexity-check**: Function complexity validation
- **test-structure-check**: Test file validation
- **documentation-check**: Required docs validation
- **makefile-targets-check**: Makefile target validation

## 🔧 Configuration Files

### `.golangci.yml`

Comprehensive golangci-lint configuration with:

```yaml
run:
  timeout: 10m
  go: "1.20"
  modules-download-mode: readonly
  allow-parallel-runners: true
  allow-serial-runners: true

linters:
  enable:
    # Code formatting and style
    - gofmt, goimports, whitespace
    - gocognit, gocyclo, gomnd
    - goprintffuncname, goconst, dupl
    - funlen, lll, nestif, wsl
    
    # Code correctness
    - govet, errcheck, staticcheck
    - gosimple, ineffassign, unused
    - deadcode, varcheck, structcheck
    - typecheck, unparam, unconvert
    - gosec, noctx, nolintlint, nakedret
    
    # Performance
    - prealloc, makezero, nilnil
    - tenv, usestdlibvars
    
    # Best practices
    - gocritic, godox, godot, goheader
    - gomnd, gomodguard, goprintffuncname
```

### `.pre-commit-config.yaml`

Pre-commit hooks configuration with:

```yaml
repos:
  # Go-specific hooks
  - repo: https://github.com/dnephin/pre-commit-golang
    rev: v0.5.1
    hooks:
      - id: go-fmt
      - id: golangci-lint
      - id: go-gosec
      # ... more hooks

  # General quality hooks
  - repo: https://github.com/pre-commit/pre-commit-hooks
    rev: v4.5.0
    hooks:
      - id: trailing-whitespace
      - id: end-of-file-fixer
      # ... more hooks

  # Custom local hooks
  - repo: local
    hooks:
      - id: go-test-coverage
      - id: security-scan
      # ... more hooks
```

### `.commitlintrc.js`

Commit message validation rules:

```javascript
module.exports = {
  extends: ['@commitlint/config-conventional'],
  rules: {
    'type-enum': [2, 'always', [
      'feat', 'fix', 'docs', 'style', 'refactor',
      'perf', 'test', 'chore', 'ci', 'build',
      'revert', 'security', 'deps', 'wip',
      'hotfix', 'release'
    ]],
    'subject-max-length': [2, 'always', 72],
    'scope-enum': [2, 'always', [
      'api', 'auth', 'build', 'ci', 'cli',
      'config', 'core', 'db', 'docs', 'feat',
      'fix', 'lint', 'perf', 'refactor',
      'security', 'storage', 'test', 'ui',
      'utils', 'web', 'deps', 'release'
    ]]
  }
};
```

## 🚦 Usage

### Command Line

```bash
# Run comprehensive linting
make lint

# Run fast linting
make lint-fast

# Fix formatting issues
make lint-fix

# Install pre-commit hooks
make install-hooks

# Run hooks on all files
make run-hooks

# Run hooks on staged files
make run-hooks-staged

# Security scan
make security

# Code complexity analysis
make complexity

# Spell checking
make spell
```

### Pre-commit Hooks

```bash
# Install hooks
pre-commit install --install-hooks

# Run on staged files (default)
pre-commit run

# Run on all files
pre-commit run --all-files

# Run specific hook
pre-commit run golangci-lint

# Update hooks
pre-commit autoupdate

# Uninstall hooks
pre-commit uninstall
```

### Manual Tool Usage

```bash
# golangci-lint
golangci-lint run ./...

# gosec security scan
gosec ./...

# gocyclo complexity
gocyclo -over 15 ./...

# misspell
misspell -error .

# goimports
goimports -w .
```

## 📊 Quality Metrics

### Coverage Thresholds

- **Minimum**: 80% (enforced)
- **Target**: 90%+ (recommended)
- **Failure**: Below 80% blocks commit

### Complexity Limits

- **Function Complexity**: Maximum 15 (cyclomatic)
- **File Complexity**: Maximum 100 (duplication)
- **Nesting Depth**: Maximum 5 levels

### Performance Standards

- **Benchmark Regression**: Detected automatically
- **Memory Usage**: Tracked per operation
- **Response Time**: Monitored for endpoints

## 🔍 Custom Hooks

### Test Coverage Hook

```yaml
- id: go-test-coverage
  name: Go Test Coverage
  entry: bash -c 'go test -coverprofile=coverage.out ./test/... && go tool cover -func=coverage.out | grep total | awk "{print \$3}" | sed "s/%//" | awk "{exit (\$1 < 80)}"'
  language: system
  pass_filenames: false
  always_run: true
  stages: [commit]
  description: "Ensure test coverage is above 80%"
```

### Security Scan Hook

```yaml
- id: security-scan
  name: Security Scan
  entry: bash -c 'if command -v gosec >/dev/null 2>&1; then gosec ./...; else echo "gosec not installed, skipping security scan"; fi'
  language: system
  pass_filenames: false
  always_run: true
  stages: [commit]
  description: "Run security vulnerability scan"
```

## 🚨 Error Handling

### Common Issues

1. **Coverage Below Threshold**
   - Add tests for uncovered code paths
   - Review test data generation
   - Check for conditional branches

2. **Complexity Violations**
   - Break down complex functions
   - Extract helper functions
   - Reduce nesting levels

3. **Security Vulnerabilities**
   - Review gosec findings
   - Update dependencies if needed
   - Implement secure coding practices

4. **Formatting Issues**
   - Run `make lint-fix`
   - Use `go fmt ./...`
   - Check import organization

### Bypassing Hooks (Emergency)

```bash
# Skip pre-commit hooks for one commit
git commit -m "Emergency fix" --no-verify

# Skip specific hook
SKIP=golangci-lint git commit -m "Skip linting"

# Run hooks manually later
pre-commit run --all-files
```

## 📈 Performance Optimization

### Fast Mode

```bash
# Quick linting check
make lint-fast

# Fast golangci-lint
golangci-lint run --fast ./...

# Skip expensive checks
SKIP=gocyclo,dupl pre-commit run
```

### Parallel Execution

```yaml
# .golangci.yml
run:
  allow-parallel-runners: true
  allow-serial-runners: true
  timeout: 10m
```

### Caching

```bash
# Pre-commit cache
pre-commit run --all-files

# golangci-lint cache
golangci-lint run --new-from-rev=HEAD~1 ./...
```

## 🔧 Integration

### CI/CD Pipeline

The linting system integrates with GitHub Actions:

```yaml
# .github/workflows/test.yml
- name: Run linter
  run: |
    golangci-lint run ./...
    
- name: Security scan
  run: |
    gosec ./...
    
- name: Code formatting
  run: |
    go fmt ./...
    goimports -w .
```

### IDE Integration

#### VS Code

```json
{
  "go.lintTool": "golangci-lint",
  "go.lintFlags": ["--fast"],
  "go.formatTool": "goimports",
  "go.useLanguageServer": true
}
```

#### GoLand

- Enable golangci-lint integration
- Configure pre-commit hooks
- Set up file watchers

### Git Hooks

```bash
# Pre-commit hook
#!/bin/sh
pre-commit run

# Commit-msg hook
#!/bin/sh
npx commitlint --edit $1
```

## 📚 Best Practices

### 1. Progressive Implementation

- Start with basic hooks (fmt, vet)
- Add security scanning gradually
- Implement complexity checks last

### 2. Team Training

- Document linting rules
- Provide examples of good/bad code
- Regular code review sessions

### 3. Performance Monitoring

- Track linting time
- Monitor false positives
- Optimize rule sets

### 4. Maintenance

- Regular hook updates
- Rule set reviews
- Performance optimization

## 🚀 Advanced Features

### Custom Rule Sets

```yaml
# .golangci.yml
linters-settings:
  gocritic:
    enabled-tags:
      - diagnostic
      - experimental
      - opinionated
      - performance
      - style
    disabled-checks:
      - commentFormatting
      - hugeParam
      - ifElseChain
```

### Conditional Hooks

```yaml
# .pre-commit-config.yaml
- id: security-scan
  name: Security Scan
  entry: bash -c 'if [ "$CI" = "true" ]; then gosec ./...; else echo "Skipping security scan in local environment"; fi'
  language: system
  pass_filenames: false
  always_run: true
  stages: [commit]
```

### Performance Profiling

```bash
# Profile linting performance
time golangci-lint run ./...

# Memory usage
/usr/bin/time -v golangci-lint run ./...

# CPU profiling
go tool pprof golangci-lint
```

## 🔍 Troubleshooting

### Common Problems

1. **Hook Installation Failures**
   ```bash
   # Clear pre-commit cache
   pre-commit clean
   
   # Reinstall hooks
   pre-commit install --install-hooks
   ```

2. **Performance Issues**
   ```bash
   # Use fast mode
   golangci-lint run --fast ./...
   
   # Skip expensive checks
   SKIP=gocyclo,dupl pre-commit run
   ```

3. **False Positives**
   ```go
   //nolint:gocyclo
   func complexFunction() {
       // Complex logic here
   }
   ```

### Debug Mode

```bash
# Enable debug output
DEBUG=1 pre-commit run

# Verbose golangci-lint
golangci-lint run -v ./...

# Show hook details
pre-commit run --verbose
```

## 📚 Resources

### Documentation
- [golangci-lint](https://golangci-lint.run/)
- [pre-commit](https://pre-commit.com/)
- [commitlint](https://commitlint.js.org/)
- [gosec](https://github.com/securecodewarrior/gosec)

### Best Practices
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go Security Best Practices](https://golang.org/doc/security)

### Community
- [golangci-lint Issues](https://github.com/golangci/golangci-lint/issues)
- [pre-commit Discussions](https://github.com/pre-commit/pre-commit/discussions)
- [Go Security Working Group](https://groups.google.com/g/golang-security)

The enhanced linting system provides comprehensive code quality assurance while maintaining development velocity. Regular use ensures consistent, secure, and maintainable code.
