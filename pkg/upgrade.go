package pkg

import (
	"fmt"

	"github.com/noizwaves/garb/pkg/github"
)

func extractReleaseVersion(binary Binary, release github.Release) (string, error) {
	matches := binary.ReleaseRegex.FindStringSubmatch(release.Name)
	if len(matches) == 0 {
		return "", fmt.Errorf("release regex did not match name %q", release.Name)
	}

	return matches[0], nil
}

func setBinaryVersion(config *configRoot, binaryName, version string) error {
	for i, binary := range config.Binaries {
		if binary.Name == binaryName {
			binary.Version = version
			config.Binaries[i] = binary

			return nil
		}
	}

	return fmt.Errorf("%q missing from config, unable to update version", binaryName)
}

func Upgrade(context Context) error {
	dirty := false
	for _, binary := range context.Binaries {
		latestRelease, err := github.GetLatestRelease(binary.Org, binary.Repo)
		if err != nil {
			return fmt.Errorf("error fetching latest release from GitHub: %w", err)
		}

		latestVersion, err := extractReleaseVersion(binary, *latestRelease)
		if err != nil {
			return fmt.Errorf("error extracting version from latest version: %w", err)
		}

		if latestVersion == binary.Version {
			fmt.Printf("%s: %s is latest\n", binary.Name, binary.Version)
		} else {
			fmt.Printf("%s: %s -> %s (%s)\n", binary.Name, binary.Version, latestVersion, latestRelease.URL)
			dirty = true
			if err = setBinaryVersion(context.Config, binary.Name, latestVersion); err != nil {
				return err
			}
		}
	}

	if dirty {
		err := saveConfig(context.Config, context.ConfigPath)
		if err != nil {
			return fmt.Errorf("error updating config file: %w", err)
		}

		fmt.Println("\nUpdated config file. Now run `garb install`.")
	}

	return nil
}
