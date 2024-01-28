package pkg

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/noizwaves/grab/pkg/github"
)

func fetchExecutable(ghClient github.Client, context *Context, binary *Binary) ([]byte, error) {
	slog.Info("Downloading asset", "binary", binary.Name, "version", binary.Version)

	asset, err := binary.GetAssetFileName(context.Platform, context.Architecture)
	if err != nil {
		return nil, fmt.Errorf("error getting asset filename: %w", err)
	}

	release, err := binary.GetReleaseName()
	if err != nil {
		return nil, fmt.Errorf("error getting asset filename: %w", err)
	}

	data, err := ghClient.DownloadReleaseAsset(binary.Org, binary.Repo, release, asset)
	if err != nil {
		return nil, fmt.Errorf("error downloading remote file: %w", err)
	}

	return extractExecutable(binary.Name, asset, &data)
}

func extractExecutable(binary, asset string, data *[]byte) ([]byte, error) {
	switch {
	case strings.HasSuffix(asset, ".tar.gz") || strings.HasSuffix(asset, ".tgz"):
		executable, err := unTgzFileNamed(binary, bytes.NewBuffer(*data))
		if err != nil {
			return nil, fmt.Errorf("error extracting binary from tgz archive: %w", err)
		}

		return executable, nil
	case strings.HasSuffix(asset, ".gz"):
		executable, err := unGzip(bytes.NewBuffer(*data))
		if err != nil {
			return nil, fmt.Errorf("error extracting binary from gzip archive: %w", err)
		}

		return executable, nil
	case strings.HasSuffix(asset, ".zip"):
		executable, err := unZipFileNamed(binary, bytes.NewBuffer(*data))
		if err != nil {
			return nil, fmt.Errorf("error extracting binary from zip archive: %w", err)
		}

		return executable, nil
	}

	return *data, nil
}

func getCurrentVersion(destPath string, binary *Binary) (string, error) {
	//nolint:gosec
	cmd := exec.Command(destPath, binary.VersionArgs...)

	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error executing binary to find version: %w", err)
	}

	matches := binary.VersionRegex.FindStringSubmatch(string(out))
	if len(matches) == 0 {
		return "", fmt.Errorf("version regex did not match command output")
	}

	return matches[0], nil
}

func writeToDisk(binary *Binary, data *[]byte, destPath string) error {
	// Explicitly remove an existing binary file. Although this could leave the file system in a bad state,
	// it is required for grab to be able to update itself.
	if _, err := os.Stat(destPath); err == nil {
		slog.Info("Removing existing binary", "name", binary.Name)

		err = os.Remove(destPath)
		if err != nil {
			return fmt.Errorf("error removing existing binary: %w", err)
		}
	}

	//nolint:gosec,gomnd
	err := os.WriteFile(destPath, *data, 0o755)
	if err != nil {
		return fmt.Errorf("error writing binary to disk: %w", err)
	}

	return nil
}

type Installer struct {
	GitHubClient github.Client
}

func (i *Installer) Install(context *Context, out io.Writer) error {
	slog.Info("Installing configured packages")

	err := context.EnsureBinPathExists()
	if err != nil {
		return fmt.Errorf("bin path needs to exist before attempting install: %w", err)
	}

	for _, binary := range context.Binaries {
		destPath := path.Join(context.BinPath, binary.Name)
		// if destination file exists
		if _, err := os.Stat(destPath); err == nil {
			currentVersion, err := getCurrentVersion(destPath, binary)
			if err != nil {
				return fmt.Errorf("failed to determine current version of %q: %w", binary.Name, err)
			}

			if binary.ShouldReplace(currentVersion) {
				fmt.Fprintf(out, "%s: installing %s over %s...", binary.Name, binary.Version, currentVersion)
			} else {
				fmt.Fprintf(out, "%s: %s already installed\n", binary.Name, currentVersion)

				continue
			}
		} else {
			fmt.Fprintf(out, "%s: installing %s...", binary.Name, binary.Version)
		}

		data, err := fetchExecutable(i.GitHubClient, context, binary)
		if err != nil {
			return fmt.Errorf("error executable binary for %s: %w", binary.Name, err)
		}

		if err := writeToDisk(binary, &data, destPath); err != nil {
			return err
		}

		fmt.Fprintln(out, " Done!")
	}

	return nil
}
