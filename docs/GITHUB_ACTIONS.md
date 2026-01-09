# GitHub Actions Workflows

This project includes comprehensive GitHub Actions workflows modeled after the high-command projects, tailored for Go-based MCP servers.

## Workflows Overview

### 1. **tests.yml** - Continuous Integration Testing
- **Trigger**: Push to `main`/`develop`, Pull Requests
- **Matrix**: Tests on Ubuntu, macOS, Windows with Go 1.23, 1.24
- **Tasks**:
  - Run `go fmt` and `go vet` linters
  - Execute test suite with coverage reporting
  - Upload coverage to Codecov
  - Run Trivy security vulnerability scanner
  - Generate coverage HTML reports

### 2. **docker.yml** - Docker Build and Push to Harbor
- **Trigger**: Push to `main`/`develop`, Pull Requests
- **Runner**: Self-hosted runner (required for Harbor access)
- **Tasks**:
  - Login to Harbor registry using GitHub secrets
  - Pull base images from Harbor DockerHub cache
  - Build Docker image with layer caching
  - Push image to Harbor: `harbor.dataknife.net/library/unifi-protect-mcp:latest`
  - Tag and push with commit SHA for traceability
  - Test the built image

### 3. **auto-approve.yml** - Automatic PR Approval
- **Trigger**: When `tests.yml` or `docker.yml` complete successfully
- **Requirements**:
  - All checks must pass
  - PR must exist
  - PR author must be a trusted collaborator (MEMBER, OWNER, COLLABORATOR)
- **Action**: Automatically approves PR with comment

### 4. **auto-assign.yml** - Automatic PR Assignment
- **Trigger**: When a Pull Request is opened
- **Action**: Automatically assigns the PR to its author

### 5. **release.yml** - Build and Release
- **Trigger**: When a version tag (v*) is pushed
- **Runner**: Self-hosted runner
- **Matrix**: Builds for:
  - Linux (amd64, arm64)
  - macOS (amd64, arm64)
  - Windows (amd64)
- **Docker Job**: Builds and pushes Docker image to Harbor with version tag
- **Outputs**:
  - Compressed binaries (.tar.gz for Unix, .zip for Windows)
  - Docker image: `harbor.dataknife.net/library/unifi-protect-mcp:v*`
  - GitHub Release with all artifacts and auto-generated notes

## Usage

### Running Tests Locally
```bash
make test
```

### Building Docker Image Locally
```bash
docker build -t unifi-protect-mcp:latest .
```

### Creating a Release
```bash
# Tag a release
git tag -a v1.0.0 -m "Release version 1.0.0"
git push origin v1.0.0

# GitHub Actions will automatically:
# 1. Build binaries for all platforms
# 2. Create a GitHub Release
# 3. Attach compiled binaries as assets
```

## Harbor Registry Configuration

This project uses Harbor (harbor.dataknife.net) as the container registry. The setup includes:

### Required GitHub Secrets

To enable Docker builds and pushes to Harbor, configure the following secrets in your GitHub repository:

1. **HARBOR_USERNAME** - Harbor robot account username
   - Example: `robot$library+ci-builder`
   - Set in: Repository Settings → Secrets and variables → Actions → New repository secret

2. **HARBOR_PASSWORD** - Harbor robot account password/token
   - Set in: Repository Settings → Secrets and variables → Actions → New repository secret

### Adding GitHub Secrets

1. Go to your repository on GitHub
2. Navigate to **Settings** → **Secrets and variables** → **Actions**
3. Click **New repository secret**
4. Add each secret:
   - Name: `HARBOR_USERNAME`
   - Value: `robot$library+ci-builder` (or your Harbor robot account username)
   
   - Name: `HARBOR_PASSWORD`
   - Value: Your Harbor robot account password/token

### Harbor Image Location

- **Registry**: `harbor.dataknife.net`
- **Project**: `library`
- **Image**: `harbor.dataknife.net/library/unifi-protect-mcp`
- **Tags**: `latest`, `<commit-sha>`, `<version-tag>` (on releases)

### Local Build and Push

To build and push locally using the Makefile:

```bash
# Set environment variables
export HARBOR_USERNAME='robot$library+ci-builder'
export HARBOR_PASSWORD='your-harbor-password'

# Or use make commands
make docker-login          # Login to Harbor
make docker-pull-base      # Pull base images from Harbor cache
make docker-push           # Build and push to Harbor
```

### Docker Compose

The `docker-compose.yml` is configured to pull from Harbor:

```bash
docker-compose pull        # Pull latest image from Harbor
docker-compose up -d       # Start using Harbor image
```

## Security Considerations

- **Trusted Approvals**: Only PRs from collaborators are auto-approved
- **Vulnerability Scanning**: Trivy scans the codebase for security issues
- **Coverage Reports**: Code coverage is tracked via Codecov
- **Docker Security**: Image is built and tested before deployment
- **Harbor Secrets**: Credentials stored securely in GitHub Secrets
- **Self-hosted Runners**: Required for Harbor registry access

## Customization

To modify workflows:
1. Edit files in `.github/workflows/`
2. Commit and push changes
3. Workflows are active immediately

## Integration with High-Command Projects

These workflows follow the same patterns as the high-command projects:
- Consistent naming and structure
- Multi-OS testing strategy
- Security-first approach with automated scanning
- Auto-approval for trusted contributors
- Release automation with cross-platform builds
