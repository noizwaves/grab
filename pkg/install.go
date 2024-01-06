package pkg

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

const localBinPath = ".local/bin"

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

func Install(context Context) error {
	for _, binary := range context.Binaries {
		// if destination file exists
		destPath := path.Join(context.HomeDir, localBinPath, binary.Name)
		if _, err := os.Stat(destPath); err == nil {
			fmt.Printf("%s already installed\n", binary.Name)

			continue
		}

		fmt.Printf("Installing %s...\n", binary.Name)

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
		err = os.WriteFile(destPath, data, 0o744)
		if err != nil {
			return fmt.Errorf("error writing binary to disk: %w", err)
		}

		fmt.Printf("%s has been installed\n", binary.Name)
	}

	return nil
}
