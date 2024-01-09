package pkg

import (
	"fmt"

	"github.com/noizwaves/dotlocalbin/pkg/github"
)

func extractReleaseVersion(binary Binary, release github.Release) (string, error) {
	matches := binary.ReleaseRegex.FindStringSubmatch(release.Name)
	if len(matches) == 0 {
		return "", fmt.Errorf("release regex did not match name %q", release.Name)
	}

	return matches[0], nil
}

func Upgrade(context Context) error {
	for _, binary := range context.Binaries {
		latestRelease, err := github.GetLatestRelease(binary.Org, binary.Repo)
		if err != nil {
			return err
		}

		latestVersion, err := extractReleaseVersion(binary, *latestRelease)
		if err != nil {
			return err
		}

		if latestVersion == binary.Version {
			fmt.Printf("%s: up-to-date %s is latest\n", binary.Name, binary.Version)
		} else {
			fmt.Printf("%s: upgrade %s -> %s (%s)\n", binary.Name, binary.Version, latestVersion, latestRelease.URL)
		}
	}

	return nil
}
