package pkg

import (
	"fmt"
	"io"
	"log/slog"

	"github.com/noizwaves/grab/pkg/github"
)

type Updater struct {
	GitHubClient github.Client
}

func (u *Updater) Update(context *Context, packageName string, out io.Writer) error {
	// Validate package name if specified
	if packageName != "" {
		if _, exists := context.Config.Packages[packageName]; !exists {
			return fmt.Errorf("package %q not found in configuration", packageName)
		}

		slog.Info("Updating specific package", "package", packageName)
	} else {
		// Check if any packages are configured
		if len(context.Config.Packages) == 0 {
			return fmt.Errorf("no packages configured in %s", context.ConfigPath)
		}

		slog.Info("Updating all configured packages")
	}

	dirty := false

	binariesToProcess := u.filterBinaries(context.Binaries, packageName)

	for _, binary := range binariesToProcess {
		latestRelease, err := u.GitHubClient.GetLatestRelease(binary.Org, binary.Repo)
		if err != nil {
			return fmt.Errorf("error fetching latest release for package %q: %w", binary.Name, err)
		}

		latestVersion, err := extractReleaseVersion(binary, latestRelease)
		if err != nil {
			return fmt.Errorf("error extracting version for package %q: %w", binary.Name, err)
		}

		if latestVersion == binary.PinnedVersion {
			fmt.Fprintf(out, "%s: %s is latest\n", binary.Name, binary.PinnedVersion)
		} else {
			fmt.Fprintf(out, "%s: %s -> %s (%s)\n", binary.Name, binary.PinnedVersion, latestVersion, latestRelease.URL)

			dirty = true

			setBinaryVersion(context.Config, binary.Name, latestVersion)
		}
	}

	if dirty {
		err := saveConfig(context.Config, context.ConfigPath)
		if err != nil {
			return fmt.Errorf("error updating config file: %w", err)
		}

		fmt.Fprintln(out, "\nUpdated config file. Now run `grab install`.")
	} else {
		slog.Debug("No config changes required, no versions were changed")
	}

	return nil
}

func (u *Updater) filterBinaries(binaries []*Binary, packageName string) []*Binary {
	if packageName == "" {
		return binaries
	}

	for _, binary := range binaries {
		if binary.Name == packageName {
			return []*Binary{binary}
		}
	}

	return []*Binary{}
}

func extractReleaseVersion(binary *Binary, release *github.Release) (string, error) {
	matches := binary.ReleaseRegex.FindStringSubmatch(release.Name)
	if len(matches) == 0 {
		return "", fmt.Errorf("release regex did not match name %q", release.Name)
	}

	return matches[0], nil
}

func setBinaryVersion(config *configRoot, binaryName, version string) {
	config.Packages[binaryName] = version
}
