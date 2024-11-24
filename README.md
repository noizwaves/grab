# grab

A fast, sudo-less package manager for your terminal programs. Downloads directly from GitHub. Supports macOS, Linux, and containers.

## Installation

Install the latest published version:
```
bash <(curl --silent https://raw.githubusercontent.com/noizwaves/grab/main/install.sh)>
```

Or install from source:
```
go build -o ~/.local/bin/grab main.go
```

## Usage

1. Add a package to the local repository. For example, `~/.grab/repository/fzf.yml`:
```yaml
apiVersion: grab.noizwaves.com/v1alpha1
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
```

2. Add the version to install to `~/.grab/config.yml`:
```yaml
packages:
  fzf: "0.45.0"
```

3. Run `go build -o ~/.local/bin/grab main.go`
4. Run `grab install` to install all programs.
5. Use the installed program:
```sh
❯ which fzf
/home/adam/.local/bin/fzf
```

### Updating versions

1. Run `grab update` to update the config file with the latest upstream versions.
1. Run `grab install` to install the updated versions.

> [!IMPORTANT]
> `update` uses the GitHub API which has a low rate limit of 60 requests/hour for anonymous users. To avoid the rate limit, [generate a token with public read-only permission](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-fine-grained-personal-access-token) and set the value via the `GH_TOKEN` environment variable.

## Development

1. Install `goenv`
1. Run `./setup.sh`
1. Generate a GitHub token and set it as the `GH_TOKEN` environment variable
