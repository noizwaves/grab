# grab

A fast, sudo-less package manager for your terminal programs. Downloads directly from GitHub. Supports macOS, Linux, and containers.

## Installation

Install the latest published version:
```sh
curl https://raw.githubusercontent.com/noizwaves/grab/main/install.sh | bash
```

Or from source:
```sh
go build -o ~/.local/bin/grab main.go
```

## Usage

1.  Install grab
1.  Add a package to the local repository. For example, `~/.grab/repository/fzf.yml`:
    ```yaml
    apiVersion: grab.noizwaves.com/v1alpha1
    kind: Package
    metadata:
      name: fzf
    spec:
      gitHubRelease:
        org: junegunn
        repo: fzf
        name: "v{{ .Version }}"
        versionRegex: \d+\.\d+\.\d+
        fileName:
          darwin,amd64: fzf-{{ .Version }}-darwin_amd64.tar.gz
          darwin,arm64: fzf-{{ .Version }}-darwin_arm64.tar.gz
          linux,amd64: fzf-{{ .Version }}-linux_amd64.tar.gz
          linux,arm64: fzf-{{ .Version }}-linux_arm64.tar.gz
      program:
        versionArgs: [--version]
        versionRegex: \d+\.\d+\.\d+
    ```
1.  Add the version to install to `~/.grab/config.yml`:
    ```yaml
    packages:
      fzf: "0.45.0"
    ```
4.  Run `grab install` to install all programs.
5.  Use the installed program:
    ```sh
    ❯ which fzf
    /home/adam/.local/bin/fzf
    ```

### Updating versions

1.  Run `grab update` to update the config file with the latest upstream versions.
1.  Run `grab install` to install the updated versions.

> [!IMPORTANT]
> `update` uses the GitHub API which has a low rate limit of 60 requests/hour for anonymous users. To avoid the rate limit, [generate a token with public read-only permission](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-fine-grained-personal-access-token) and set the value via the `GH_TOKEN` environment variable.

## Configuration Reference

### Package Definition Reference

Package definitions are stored as YAML files in `~/.grab/repository/` and follow this structure:

```yaml
apiVersion: grab.noizwaves.com/v1alpha1
kind: Package
metadata:
  name: package-name
spec:
  gitHubRelease:
    org: github-org
    repo: github-repo
    name: release-name-template
    versionRegex: \d+\.\d+\.\d+
    fileName:
      darwin,amd64: "archive-{{ .Version }}-darwin_amd64.tar.gz"
      darwin,arm64: "archive-{{ .Version }}-darwin_arm64.tar.gz"
      linux,amd64: "archive-{{ .Version }}-linux_amd64.tar.gz"
      linux,arm64: "archive-{{ .Version }}-linux_arm64.tar.gz"
    embeddedBinaryPath:
      darwin,amd64: "path/to/binary"
      linux,amd64: "bin/binary"
  program:
    versionArgs: [--version]
    versionRegex: \d+\.\d+\.\d+
```

**Metadata**
- `name`: Unique identifier for the package

**GitHub Release Configuration**
- `org`: GitHub organization or username
- `repo`: GitHub repository name
- `name`: Release name template (Go templated string, with `Version` available)
- `versionRegex`: Regular expression to extract version numbers from release names
- `fileName`: Platform-specific asset archive filenames (Go templated string, with `Version` available)
- `embeddedBinaryPath`: _(Optional)_ Platform-specific path to binary within the archive

**Program Configuration**
- `versionArgs`: Command-line arguments to retrieve the program's version
- `versionRegex`: Regular expression to extract version from program output

### User Configuration Reference

User configuration is stored in `~/.grab/config.yml` and specifies which versions to install:

```yaml
packages:
  package-name: "1.2.3"
  another-package: "2.0.1"
```

### Supported Platforms

- `darwin,amd64`: macOS on Intel processors
- `darwin,arm64`: macOS on Apple Silicon
- `linux,amd64`: Linux on x86_64 processors
- `linux,arm64`: Linux on ARM64 processors

### Advanced Examples

#### Complex Package with Embedded Binary Path

```yaml
apiVersion: grab.noizwaves.com/v1alpha1
kind: Package
metadata:
  name: yazi
spec:
  gitHubRelease:
    org: sxyazi
    repo: yazi
    name: "v{{ .Version }}"
    versionRegex: \d+\.\d+\.\d+
    fileName:
      darwin,amd64: "yazi-x86_64-apple-darwin.zip"
      darwin,arm64: "yazi-aarch64-apple-darwin.zip"
      linux,amd64: "yazi-x86_64-unknown-linux-gnu.zip"
      linux,arm64: "yazi-aarch64-unknown-linux-gnu.zip"
    embeddedBinaryPath:
      darwin,amd64: "yazi-x86_64-apple-darwin/yazi"
      darwin,arm64: "yazi-aarch64-apple-darwin/yazi"
      linux,amd64: "yazi-x86_64-unknown-linux-gnu/yazi"
      linux,arm64: "yazi-aarch64-unknown-linux-gnu/yazi"
  program:
    versionArgs: [--version]
    versionRegex: \d+\.\d+\.\d+
```

## Development

1.  [Install Mise](https://mise.jdx.dev/installing-mise.html)
1.  Run `./setup.sh`
1.  Generate a GitHub token and set it as the `GH_TOKEN` environment variable
