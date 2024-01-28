package pkg

import (
	"fmt"
	"io"
	"log/slog"

	"github.com/noizwaves/grab/pkg/github"
)

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

func Update(context *Context, out io.Writer) error {
	slog.Info("Updating configured packages")

	dirty := false
	for _, binary := range context.Binaries {
		latestRelease, err := github.GetLatestRelease(binary.Org, binary.Repo)
		if err != nil {
			return fmt.Errorf("error fetching latest release from GitHub: %w", err)
		}

		latestVersion, err := extractReleaseVersion(binary, latestRelease)
		if err != nil {
			return fmt.Errorf("error extracting version from latest version: %w", err)
		}

		if latestVersion == binary.Version {
			fmt.Fprintf(out, "%s: %s is latest\n", binary.Name, binary.Version)
		} else {
			fmt.Fprintf(out, "%s: %s -> %s (%s)\n", binary.Name, binary.Version, latestVersion, latestRelease.URL)
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
