# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**grab** is a Go-based CLI package manager that downloads terminal programs directly from GitHub releases. It's designed as a fast, sudo-less alternative to traditional package managers, focusing specifically on binary executables from GitHub repositories.

## Common Development Commands

### Building and Testing
```bash
# Build the application
go build -o ./grab main.go

# Run all tests
go test -count=1 ./...

# Run tests for specific package
go test -count=1 ./pkg/github

# Run with coverage
go test -cover ./...
```

### Code Quality
```bash
# Run linter (this is run automatically on pre-commit)
golangci-lint run ./...

# Format code
gofumpt -w .
```

### Environment Setup
```bash
# Set up development environment with Mise
./setup.sh

# Set GitHub token to avoid API rate limits
export GH_TOKEN=your_github_token_here
```

## Architecture

### Core Components

- **`cmd/`**: CLI command implementations using Cobra framework
  - `install.go`: Installs/updates packages based on configuration
  - `update.go`: Updates config with latest upstream versions
  - `root.go`: Root command setup and configuration loading

- **`pkg/`**: Core business logic
  - `installer.go`: Package installation and archive extraction logic
  - `updater.go`: Version checking and config updating
  - `config.go`: YAML configuration file handling
  - `model.go`: Package definitions and URL templating
  - `github/`: GitHub API client for release information
  - `archive.go`: Multi-format archive extraction (tar.gz, tar.xz, zip, gzip)

### Package Definition System

Uses Kubernetes-style YAML manifests in `~/.grab/repository/` to define packages:

```yaml
apiVersion: grab.noizwaves.com/v1alpha1
kind: Package
metadata:
  name: package-name
spec:
  gitHubRelease:
    org: github-org
    repo: github-repo
    name: "v{{ .Version }}"
    versionRegex: \d+\.\d+\.\d+
    fileName:
      darwin,amd64: "filename-{{ .Version }}-darwin_amd64.tar.gz"
    embeddedBinaryPath:
      darwin,amd64: "path/to/binary"
  program:
    versionArgs: [--version]
    versionRegex: \d+\.\d+\.\d+
```

User configuration in `~/.grab/config.yml` maps package names to desired versions.

### Key Dependencies

- **Cobra**: CLI framework (`github.com/spf13/cobra`)
- **Viper**: Configuration management (`github.com/spf13/viper`)
- **Testify**: Testing framework (`github.com/stretchr/testify`)
- **Archive handling**: Support for tar.gz, tar.xz, zip, gzip formats

### Testing Strategy

- Unit tests for all core functionality
- Mock GitHub client for testing (`pkg/internal/github/`)
- Test helpers in `pkg/internal/assert/`
- Pre-commit hooks run tests and linting automatically via Lefthook

### Template System

Uses Go templates for dynamic URL and filename generation based on:
- `{{ .Version }}`: Package version
- Platform detection (darwin/linux, amd64/arm64)
- Release name patterns
- Binary path templates for custom archive structures

## Important Notes

- All binaries install to `~/.local/bin/` (no sudo required)
- Binary paths within archives are configurable via `embeddedBinaryPath` (defaults to package name)
- GitHub API rate limits apply (60 requests/hour without token)
- Supports cross-platform builds for macOS and Linux on amd64/arm64
- Uses semantic version comparison for updates
- Configuration files use YAML format exclusively

## Task Master AI Instructions
**Import Task Master's development workflow commands and guidelines, treat as if import is in the main CLAUDE.md file.**
@./.taskmaster/CLAUDE.md

## Development Guidelines

- Do not test private functions directly. Instead, call the public functions that use the private functions.
