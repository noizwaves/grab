# garb

A fast sudo-less package manager for your terminal programs. Supports macOS, Linux, and containers.

## Usage

1. Add a package to the local repository. For example, `~/.garb/repository/fzf.yml`:
```yaml
apiVersion: garb.noizwaves.com/v1alpha1
kind: Package
metadata:
  name: fzf
spec:
  gitHubRelease:
    org: junegunn
    repo: fzf
    name: "{{ .Version }}"
    versionRegex: \d+\.\d+\.\d+
    fileName:
      darwin,amd64: fzf-{{ .Version }}-darwin_amd64.zip
      darwin,arm64: fzf-{{ .Version }}-darwin_arm64.zip
      linux,amd64: fzf-{{ .Version }}-linux_amd64.tar.gz
      linux,arm64: fzf-{{ .Version }}-linux_arm64.tar.gz
  program:
    versionArgs: [--version]
    versionRegex: \d+\.\d+\.\d+
...
```

2. Add the version to install to `~/.garb/config.yml`:
```yaml
packages:
  fzf: "0.45.0"
```

3. Run `go build -o ~/.local/bin/garb main.go`
4. Run `garb install` to install all programs.
5. Use the installed program:
```sh
❯ which fzf
/home/adam/.local/bin/fzf
```

### Upgrading package versions

1. Run `garb upgrade` to update the config file
1. Run `garb install` to install the updated programs

> [!IMPORTANT]
> `upgrade` uses the GitHub API which has a low rate limit of 60 requests/hour. To avoid the rate limit, [generate a token with public read-only permission](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-fine-grained-personal-access-token) and set the value via the `GH_TOKEN` environment variable.

## Development

1. Install `goenv` and specified Go version
1. Install `golangci-lint`: `go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2`
1. Install `gofumpt`: `go install mvdan.cc/gofumpt@latest`
1. Generate a GitHub token and set it as the `GH_TOKEN` environment variable
