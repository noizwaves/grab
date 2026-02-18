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

### Adding new packages

The quickest way to add a new package is with `grab get`:

```sh
grab get https://github.com/junegunn/fzf
```

This will automatically import the package spec, add the latest version to the config, and install the binary.

A custom package name can be specified using the `--name`/`-n` option:

```sh
grab get --name my-tool https://github.com/org/repo
```

For more control, packages can be added step-by-step using `grab import` to import GitHub releases into the local repository:

```sh
# Any of these formats work:
grab import https://github.com/junegunn/fzf
grab import https://github.com/junegunn/fzf/releases/latest
grab import https://github.com/junegunn/fzf/issues
```

The import command will automatically fetch the latest release from the repository and generate the appropriate package configuration.
Packages will be named after the repository slug by default, and can be overridden using the `--name`/`-n` option.
After importing, add the desired version to `~/.grab/config.yml` and run `grab install`.

Packages can also be defined manually by editing package YML files at `~/.grab/repository/*.yml`.

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
      linux,amd64: "{{ .Version }}ed/path/to/binary"
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
- `embeddedBinaryPath`: _(Optional)_ Platform-specific path to binary within the archive (Go templated string, with `Version` available)

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

#### Complex package with version specific embedded binary path

```yaml
apiVersion: grab.noizwaves.com/v1alpha1
kind: Package
metadata:
  name: crush
spec:
  gitHubRelease:
    org: charmbracelet
    repo: crush
    name: "v{{ .Version }}"
    versionRegex: \d+\.\d+\.\d+
    fileName:
      darwin,amd64: crush_{{ .Version }}_Darwin_x86_64.tar.gz
      darwin,arm64: crush_{{ .Version }}_Darwin_arm64.tar.gz
      linux,amd64: crush_{{ .Version }}_Linux_x86_64.tar.gz
      linux,arm64: crush_{{ .Version }}_Linux_arm64.tar.gz
    embeddedBinaryPath:
      darwin,amd64: crush_{{ .Version }}_Darwin_x86_64/crush
      darwin,arm64: crush_{{ .Version }}_Darwin_arm64/crush
      linux,amd64: crush_{{ .Version }}_Linux_x86_64/crush
      linux,arm64: crush_{{ .Version }}_Linux_arm64/crush
  program:
    versionArgs: [--version]
    versionRegex: \d+\.\d+\.\d+
```

## Development

1.  [Install Mise](https://mise.jdx.dev/installing-mise.html)
1.  Run `./setup.sh`
1.  Generate a GitHub token and set it as the `GH_TOKEN` environment variable
