package pkg

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/noizwaves/grab/pkg/github"
)

type Installer struct {
	GitHubClient github.Client
}

func (i *Installer) Install(context *Context, packageName string, out io.Writer) error {
	if packageName != "" {
		slog.Info("Installing specific package", "package", packageName)
	} else {
		slog.Info("Installing configured packages")
	}

	err := context.EnsureBinPathExists()
	if err != nil {
		return fmt.Errorf("bin path needs to exist before attempting install: %w", err)
	}

	binariesToProcess, err := i.getBinariesToProcess(context, packageName)
	if err != nil {
		return err
	}

	for _, binary := range binariesToProcess {
		err := i.installBinary(context, binary, out)
		if err != nil {
			return err
		}
	}

	return nil
}

func (i *Installer) getBinariesToProcess(context *Context, packageName string) ([]*Binary, error) {
	if packageName != "" {
		foundBinary := i.findBinaryByName(context.Binaries, packageName)
		if foundBinary == nil {
			return nil, errors.New("package definition for " + packageName + " not found")
		}

		return []*Binary{foundBinary}, nil
	}

	return context.Binaries, nil
}

func (i *Installer) findBinaryByName(binaries []*Binary, packageName string) *Binary {
	for _, binary := range binaries {
		if binary.Name == packageName {
			return binary
		}
	}

	return nil
}

func (i *Installer) installBinary(context *Context, binary *Binary, out io.Writer) error {
	destPath := path.Join(context.BinPath, binary.Name)

	// if destination file exists
	_, err := os.Stat(destPath)
	if err == nil {
		currentVersion, err := getCurrentVersion(destPath, binary)
		if err != nil {
			return fmt.Errorf("failed to determine current version of %q: %w", binary.Name, err)
		}

		if binary.ShouldReplace(currentVersion) {
			fmt.Fprintf(out, "%s: installing %s over %s...", binary.Name, binary.PinnedVersion, currentVersion)
		} else {
			fmt.Fprintf(out, "%s: %s already installed\n", binary.Name, currentVersion)

			return nil
		}
	} else {
		fmt.Fprintf(out, "%s: installing %s...", binary.Name, binary.PinnedVersion)
	}

	data, err := fetchExecutable(i.GitHubClient, context, binary)
	if err != nil {
		return fmt.Errorf("error executable binary for %s: %w", binary.Name, err)
	}

	err = writeToDisk(binary, &data, destPath)
	if err != nil {
		return err
	}

	fmt.Fprintln(out, " Done!")

	return nil
}

func fetchExecutable(ghClient github.Client, context *Context, binary *Binary) ([]byte, error) {
	slog.Info("Downloading asset", "binary", binary.Name, "version", binary.PinnedVersion)

	asset, err := binary.GetAssetFileName(context.Platform, context.Architecture)
	if err != nil {
		return nil, fmt.Errorf("error getting asset filename: %w", err)
	}

	embeddedBinaryPath, err := binary.GetEmbeddedBinaryPath(context.Platform, context.Architecture)
	if err != nil {
		return nil, fmt.Errorf("error getting embedded binary path: %w", err)
	}

	release, err := binary.GetReleaseName()
	if err != nil {
		return nil, fmt.Errorf("error getting asset filename: %w", err)
	}

	data, err := ghClient.DownloadReleaseAsset(binary.Org, binary.Repo, release, asset)
	if err != nil {
		return nil, fmt.Errorf("error downloading remote file: %w", err)
	}

	return extractExecutable(embeddedBinaryPath, asset, &data)
}

func extractExecutable(binary, asset string, data *[]byte) ([]byte, error) {
	switch {
	case strings.HasSuffix(asset, ".tar.gz") || strings.HasSuffix(asset, ".tgz"):
		executable, err := unTgzFileNamed(binary, bytes.NewBuffer(*data))
		if err != nil {
			return nil, fmt.Errorf("error extracting binary from tgz archive: %w", err)
		}

		return executable, nil
	case strings.HasSuffix(asset, ".tar.xz"):
		executable, err := unTarxzFileNamed(binary, bytes.NewBuffer(*data))
		if err != nil {
			return nil, fmt.Errorf("error extracting binary from xz archive: %w", err)
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
		return "", errors.New("version regex did not match command output")
	}

	return matches[0], nil
}

// Write the executable to disk as atomically as possible.
// First, it writes to a temporary file in the destination directory,
// then it moves the temporary file to the destination path.
func writeToDisk(binary *Binary, data *[]byte, destPath string) error {
	// Use dest instead of /tmp for temporary file writing; avoids the
	// "invalid cross-device link" error when /tmp is on a different device
	// i.e. memory mounted
	destDir := path.Dir(destPath)
	tempPath := path.Join(destDir, ".grab-temp-"+binary.Name)
	slog.Debug("Writing to temporary executable", "binary", binary.Name, "tempPath", tempPath)

	// Ensure temp path is clear
	err := removeFileIfPresent(tempPath)
	if err != nil {
		return fmt.Errorf("error removing temp file: %w", err)
	}

	//nolint:gosec,mnd
	err = os.WriteFile(tempPath, *data, 0o755)
	if err != nil {
		return fmt.Errorf("error writing executable to temp location: %w", err)
	}

	// Best-effort clean up if rename fails
	defer tryRemoveFromFilesystem(tempPath)

	err = os.Rename(tempPath, destPath)
	if err != nil {
		return fmt.Errorf("error moving temp to destination: %w", err)
	}

	return nil
}

// Best effort to remove a file or directory from filesystem, and warn on an error.
func tryRemoveFromFilesystem(path string) {
	_, err := os.Stat(path)
	if err == nil {
		err := os.Remove(path)
		if err != nil {
			slog.Warn("Failed to remove file", "path", path)
		}
	}
}

func removeFileIfPresent(path string) error {
	_, err := os.Stat(path)
	if err == nil {
		err := os.Remove(path)
		if err != nil {
			return fmt.Errorf("error removing file '%s': %w", path, err)
		}
	}

	return nil
}
