# garb

A fast sudo-less package manager for your terminal programs. Supports macOS, Linux, and containers.

## Usage

1. Create a config file at `~/.garb.yml` like:
```yaml
binaries:
  - name: fzf
    source: https://github.com/junegunn/fzf/releases/download/{{ .Version }}/fzf-{{ .Version }}-linux_amd64.tar.gz
    version: "0.45.0"
...
```

2. Run `go build -o ~/.local/bin/garb main.go`
3. Run `garb install` to install the binaries

### Upgrading binary versions

1. Run `garb upgrade` to update the config file
1. Run `garb install` to install the updated binaries

> [!IMPORTANT]
> `upgrade` uses the GitHub API which has a low rate limit of 60 requests/hour. To avoid the rate limit, [generate a token with public read-only permission](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-fine-grained-personal-access-token) and set the value via the `GH_TOKEN` environment variable.

## Configuration Format

- `binaries[*].source` (required): a Go template compatible string with access to variables:
  - `Version`
  - `Platform`: defaults to current `GOOS`
  - `Arch`: defaults to current `GOARCH`
  - `Ext`: defaults to `""` (empty string)
- `binaries[*].platforms` (optional): per platform (`GOOS`)/architecture (`GOARCH`)/extension overrides for `Platform`, `Arch` and `Ext` values in source URL. Takes the form:
  ```yaml
  linux:
    amd64: [unknown-linux-musl, x86_64, tar.gz]
    arm64: [unknown-linux-gnu, aarch64, tar.gz]
  darwin:
    amd64: [apple-darwin, x86_64, zip]
    arm64: [apple-darwin, aarch64, zip]
  ```

## Development

1. Install `goenv` and specified Go version
1. Install `golangci-lint`: `go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2`
1. Install `gofumpt`: `go install mvdan.cc/gofumpt@latest`
