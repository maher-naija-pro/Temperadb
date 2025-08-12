# CI/CD Pipeline Documentation

This document describes the Continuous Integration and Continuous Deployment (CI/CD) pipeline for TimeSeriesDB.

## Overview

The CI/CD system automatically builds, tests, and publishes packages for multiple platforms whenever code changes are pushed or tags are created.

## Workflows

### 1. Build and Publish Packages (`build-packages.yml`)

**Triggers:**
- Push to `main` or `master` branch
- Push of a tag (e.g., `v1.0.0`)
- Pull request to `main` or `master`
- Manual workflow dispatch

**Jobs:**

#### Build Job
- **Matrix Strategy**: Builds for multiple platforms simultaneously
- **Supported Targets**:
  - Linux AMD64
  - Linux ARM64
  - Windows AMD64
  - macOS AMD64
  - macOS ARM64 (Apple Silicon)

- **Artifacts**: Creates platform-specific binaries and archives
- **Docker**: Builds Docker images for Linux platforms

#### Release Job
- **Condition**: Only runs on tag pushes
- **Action**: Creates GitHub release with all built artifacts
- **Files**: Uploads `.tar.gz` and `.zip` files for each platform

#### Publish Docker Job
- **Condition**: Only runs on tag pushes
- **Action**: Publishes Docker images to GitHub Container Registry
- **Tags**: Creates versioned and `latest` tags

#### Test Binaries Job
- **Condition**: Only runs on tag pushes
- **Action**: Tests the built binaries to ensure they work correctly

#### Security Scan Job
- **Condition**: Only runs on tag pushes
- **Action**: Runs Trivy vulnerability scanner
- **Output**: Uploads results to GitHub Security tab

## Local Development

### Prerequisites

- Go 1.20 or later
- Docker (optional, for Docker builds)
- Make

### Building Locally

```bash
# Basic build for current platform
make build

# Platform-specific builds
make build-linux      # Linux AMD64
make build-windows    # Windows AMD64
make build-darwin     # macOS AMD64

# Build for all platforms
make build-all

# Docker build
make build-docker

# Clean build artifacts
make clean-build

# Show build help
make build-help
```

### Testing the Build System

```bash
# Run the build test script
./scripts/build-test.sh
```

This script will:
1. Test all build targets
2. Verify platform-specific builds
3. Test Docker builds (if available)
4. Clean up artifacts

## GitHub Packages

### Container Registry

Docker images are automatically published to:
```
ghcr.io/yourusername/timeseriesdb:latest
ghcr.io/yourusername/timeseriesdb:v1.0.0
ghcr.io/yourusername/timeseriesdb:v1.0.0-amd64
ghcr.io/yourusername/timeseriesdb:v1.0.0-arm64
```

### Releases

GitHub releases are automatically created with:
- Release notes generated from commits
- Downloadable binaries for all platforms
- Source code archives

## Configuration

### Environment Variables

The build system uses these environment variables:
- `VERSION`: Version string (auto-detected from git tags)
- `BUILD_TIME`: Build timestamp
- `COMMIT_HASH`: Git commit hash

### Build Arguments

Docker builds accept these arguments:
- `VERSION`: Version to build
- `GOOS`: Target operating system
- `GOARCH`: Target architecture

## Workflow Customization

### Adding New Platforms

To add a new platform, modify the matrix in `build-packages.yml`:

```yaml
- os: newos
  arch: newarch
  goos: newos
  goarch: newarch
  docker_platform: newos/newarch
```

### Modifying Build Process

The build process can be customized by:
1. Modifying the `Makefile` build targets
2. Updating the Dockerfile
3. Changing the GitHub Actions workflow steps

## Troubleshooting

### Common Issues

1. **Build Failures**: Check Go version compatibility
2. **Docker Build Issues**: Ensure Docker is running and accessible
3. **Permission Errors**: Check GitHub repository permissions for packages

### Debug Mode

Enable debug output in GitHub Actions by setting the secret `ACTIONS_STEP_DEBUG` to `true`.

### Local Debugging

Use the build test script to verify your local setup:
```bash
./scripts/build-test.sh
```

## Best Practices

1. **Versioning**: Use semantic versioning (e.g., `v1.0.0`)
2. **Testing**: Always test builds locally before pushing
3. **Security**: Regularly review security scan results
4. **Documentation**: Update this document when making changes

## Support

For issues with the CI/CD pipeline:
1. Check GitHub Actions logs
2. Verify local builds work
3. Review workflow configuration
4. Check GitHub repository settings
