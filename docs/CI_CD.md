# CI/CD Guide

This document describes the Continuous Integration and Continuous Deployment (CI/CD) setup for TimeSeriesDB, including Docker image building and automated workflows.

## Overview

TimeSeriesDB uses GitHub Actions for automated CI/CD processes. The workflows automatically build, test, and deploy Docker images to GitHub Container Registry (GHCR) on specific events.

## Workflows

### 1. Docker Build (Simple) - `docker-simple.yml`

**Purpose:** Build and push Docker images for main merges and tags.

**Triggers:**
- Push to `main` or `master` branch
- Push of version tags (e.g., `v1.0.0`)
- Manual workflow dispatch

**Features:**
- Single platform build (linux/amd64)
- Automatic versioning
- Image testing
- GitHub Container Registry integration

**Output:**
- `ghcr.io/maher-naija-pro/my-timeserie:latest`
- `ghcr.io/maher-naija-pro/my-timeserie:{version}` (for tags)

### 2. Docker Build and Push - `docker-build.yml`

**Purpose:** Comprehensive Docker build workflow with multi-architecture support and security scanning.

**Triggers:**
- Push to `main` or `master` branch
- Push of version tags
- Pull requests (build only, no push)
- Manual workflow dispatch

**Features:**
- Multi-architecture builds (amd64, arm64, arm/v7)
- Security vulnerability scanning with Trivy
- Comprehensive image testing
- Discord notifications
- GitHub Security tab integration

**Output:**
- Multi-arch images for all supported platforms
- Security scan results in GitHub Security tab

## Docker Image Tags

### Automatic Tagging

The workflows automatically create appropriate tags based on the trigger:

- **Main branch:** `latest`, `main`
- **Version tags:** `v1.0.0`, `v1.0`, `v1`
- **Pull requests:** `dev` (build only)

### Manual Tagging

You can manually trigger builds with custom versions:

1. Go to Actions → Docker Build (Simple)
2. Click "Run workflow"
3. Enter your desired version (e.g., `v2.0.0-beta`)
4. Click "Run workflow"

## Using the Docker Images

### Pull Images

```bash
# Latest stable
docker pull ghcr.io/maher-naija-pro/my-timeserie:latest

# Specific version
docker pull ghcr.io/maher-naija-pro/my-timeserie:v1.0.0

# Main branch (development)
docker pull ghcr.io/maher-naija-pro/my-timeserie:main
```

### Run Containers

```bash
# Basic run
docker run -d \
  --name timeseriesdb \
  -p 8080:8080 \
  ghcr.io/maher-naija-pro/my-timeserie:latest

# With persistent storage
docker run -d \
  --name timeseriesdb \
  -p 8080:8080 \
  -v $(pwd)/data:/app/data \
  ghcr.io/maher-naija-pro/my-timeserie:latest
```

## Workflow Configuration

### Environment Variables

- `REGISTRY`: GitHub Container Registry (`ghcr.io`)
- `IMAGE_NAME`: Repository name (auto-detected)

### Required Secrets

- `GITHUB_TOKEN`: Automatically provided by GitHub
- `DISCORD_WEBHOOK`: For failure notifications (optional)

### Permissions

- `contents: read` - Read repository contents
- `packages: write` - Push to container registry
- `security-events: write` - Upload security scan results

## Build Process

### 1. Code Checkout
- Clones the repository with full history
- Sets up the build environment

### 2. Docker Setup
- Installs Docker Buildx for multi-platform builds
- Logs into GitHub Container Registry

### 3. Image Building
- Builds Docker image using the Dockerfile
- Applies appropriate tags and labels
- Uses GitHub Actions cache for faster builds

### 4. Image Testing
- Pulls the built image
- Starts a container
- Tests basic functionality
- Verifies health endpoint (if available)

### 5. Security Scanning (Full workflow only)
- Runs Trivy vulnerability scanner
- Uploads results to GitHub Security tab

### 6. Image Push
- Pushes images to GitHub Container Registry
- Only for main branch and tags (not PRs)

## Monitoring and Notifications

### Discord Notifications

The workflows send notifications to Discord on success/failure:

- **Success:** ✅ Docker Build Successful
- **Failure:** 🚨 Docker Build Failed

### GitHub Security Tab

Security scan results are automatically uploaded to:
- Repository → Security → Code scanning

### Workflow Status

Monitor workflow execution at:
- Repository → Actions → Workflows

## Troubleshooting

### Common Issues

1. **Build Failures**
   - Check the Actions tab for detailed logs
   - Verify Dockerfile syntax
   - Check for dependency issues

2. **Authentication Errors**
   - Ensure `GITHUB_TOKEN` has proper permissions
   - Check repository settings for Actions permissions

3. **Image Push Failures**
   - Verify container registry permissions
   - Check for duplicate tags

4. **Test Failures**
   - Review container startup logs
   - Check health endpoint configuration
   - Verify port binding

### Debug Mode

Enable debug logging by setting the secret:
```
ACTIONS_STEP_DEBUG=true
```

### Manual Testing

Test Docker builds locally:

```bash
# Build image
docker build -t timeseriesdb:test .

# Run container
docker run -d -p 8080:8080 timeseriesdb:test

# Test functionality
curl http://localhost:8080/health
```

## Best Practices

### 1. Version Management
- Use semantic versioning for releases
- Tag releases with `v` prefix (e.g., `v1.0.0`)
- Keep `latest` tag updated

### 2. Security
- Regularly update base images
- Monitor security scan results
- Review and address vulnerabilities

### 3. Testing
- Test images before pushing to production
- Verify all endpoints work correctly
- Check resource usage and performance

### 4. Monitoring
- Set up alerts for build failures
- Monitor image pull statistics
- Track security scan results

## Advanced Configuration

### Custom Build Arguments

Modify the Dockerfile to accept additional build arguments:

```dockerfile
ARG CUSTOM_FEATURE=false
ENV CUSTOM_FEATURE=${CUSTOM_FEATURE}
```

### Multi-Stage Builds

The current Dockerfile uses multi-stage builds for optimization:
- Builder stage compiles the Go application
- Runtime stage creates the final image

### Platform-Specific Builds

The full workflow supports multiple architectures:
- `linux/amd64` - Standard x86_64
- `linux/arm64` - ARM 64-bit
- `linux/arm/v7` - ARM 32-bit v7

## Support

For CI/CD related issues:

1. Check the Actions tab for workflow logs
2. Review this documentation
3. Check GitHub Actions documentation
4. Create an issue with workflow details

## References

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Docker Buildx](https://docs.docker.com/buildx/)
- [GitHub Container Registry](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry)
- [Trivy Security Scanner](https://aquasecurity.github.io/trivy/)
