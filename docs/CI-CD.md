# CI/CD Pipeline - Auto PR Creator

## Overview
This CI/CD pipeline automatically creates pull requests from test branches to main when they don't already exist.

## How It Works

### Workflows
1. **`auto-pr-creator.yml`** - Full-featured workflow with detailed logging
2. **`auto-pr-creator-fast.yml`** - Optimized for speed with minimal steps

### Triggers
- **Automatic**: Pushes to branches matching patterns:
  - `test*` (e.g., `test`, `test-feature`, `test-bugfix`)
  - `test`
  - `feature*` (e.g., `feature/new-query`, `feature/optimization`)
  - `dev*` (e.g., `dev`, `dev-experimental`)
- **Manual**: Can be triggered manually via GitHub Actions

### Process
1. **Check**: Verifies if a PR already exists for the branch
2. **Create**: If no PR exists, creates a draft PR with:
   - Title: "Auto PR: {branch} → main"
   - Draft status (requires manual review)
   - Base: `main`
   - Head: current branch

### Features
- ✅ **Smart**: Only creates PRs when they don't exist
- ✅ **Fast**: Minimal checkout depth and efficient checks
- ✅ **Safe**: Creates draft PRs requiring manual review
- ✅ **Flexible**: Works with any test/feature/dev branch pattern

## Usage

### Automatic
Simply push to any matching branch:
```bash
git checkout -b test-new-feature
git push origin test-new-feature
# PR will be created automatically
```

### Manual Trigger
1. Go to GitHub Actions
2. Select "Auto PR Creator" workflow
3. Click "Run workflow"

## Configuration

### Branch Patterns
Edit the workflow files to modify which branches trigger PR creation:
```yaml
branches:
  - 'test*'      # All branches starting with 'test'
  - 'feature*'   # All branches starting with 'feature'
  - 'dev*'       # All branches starting with 'dev'
```

### PR Settings
- **Draft**: All PRs are created as drafts
- **Base**: Always targets `main` branch
- **Title**: Auto-generated with branch name
- **Body**: Simple template with checklist

## Requirements
- GitHub repository with Actions enabled
- `GITHUB_TOKEN` secret (automatically provided)
- GitHub CLI available in runner (included by default)

## Troubleshooting

### PR Not Created
- Check if branch name matches patterns
- Verify GitHub Actions are enabled
- Check workflow run logs for errors

### Duplicate PRs
- The workflow checks for existing PRs
- If duplicates occur, check branch naming conflicts

### Performance
- Use `auto-pr-creator-fast.yml` for maximum speed
- Minimal checkout depth reduces execution time
