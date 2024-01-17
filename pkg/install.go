package pkg

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
)

func downloadArtifact(sourceURL string) ([]byte, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, sourceURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error requesting binary: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	return data, nil
}

func fetchBinaryData(binaryName string, sourceURL string) ([]byte, error) {
	slog.Info("Downloading artifact", "url", sourceURL)
	data, err := downloadArtifact(sourceURL)
	if err != nil {
		return nil, fmt.Errorf("error downloading remote file: %w", err)
	}

	switch {
	case strings.HasSuffix(sourceURL, ".tar.gz") || strings.HasSuffix(sourceURL, ".tgz"):
		data, err = unTgzFileNamed(binaryName, bytes.NewBuffer(data))

		if err != nil {
			return nil, fmt.Errorf("error extracting binary from tgz archive: %w", err)
		}
	case strings.HasSuffix(sourceURL, ".gz"):
		data, err = unGzip(bytes.NewBuffer(data))

		if err != nil {
			return nil, fmt.Errorf("error extracting binary from gzip archive: %w", err)
		}
	case strings.HasSuffix(sourceURL, ".zip"):
		data, err = unZipFileNamed(binaryName, bytes.NewBuffer(data))

		if err != nil {
			return nil, fmt.Errorf("error extracting binary from zip archive: %w", err)
		}
	}

	return data, nil
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

func Install(context *Context) error {
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
				fmt.Printf("%s: installing %s over %s...", binary.Name, binary.Version, currentVersion)
			} else {
				fmt.Printf("%s: %s already installed\n", binary.Name, currentVersion)

				continue
			}
		} else {
			fmt.Printf("%s: installing %s...", binary.Name, binary.Version)
		}

		// download and extract target URL
		sourceURL, err := binary.GetURL(context.Platform, context.Architecture)
		if err != nil {
			return fmt.Errorf("error getting source url for %s: %w", binary.Name, err)
		}

		data, err := fetchBinaryData(binary.Name, sourceURL)
		if err != nil {
			return fmt.Errorf("error downloading binary for %s: %w", binary.Name, err)
		}

		// write binary as executable to file system
		//nolint:gosec,gomnd
		err = os.WriteFile(destPath, data, 0o755)
		if err != nil {
			return fmt.Errorf("error writing binary to disk: %w", err)
		}

		fmt.Println(" Done!")
	}

	return nil
}
