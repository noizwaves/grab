package pkg

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"log/slog"

	"github.com/ulikunitz/xz"
)

func unTgzFileNamed(binaryPath string, data io.Reader) ([]byte, error) {
	slog.Info("Extracting file from tgz archive", "path", binaryPath)

	decompressed, err := gzip.NewReader(data)
	if err != nil {
		return nil, fmt.Errorf("error decompressing Gzipped data: %w", err)
	}

	return unTar(binaryPath, decompressed) //golint:nowrap
}

func unZipFileNamed(binaryPath string, data io.Reader) ([]byte, error) {
	slog.Info("Extracting file from zip archive", "path", binaryPath)

	raw, err := io.ReadAll(data)
	if err != nil {
		return nil, fmt.Errorf("error reading raw Zip file: %w", err)
	}

	reader := bytes.NewReader(raw)

	decompressed, err := zip.NewReader(reader, int64(len(raw)))
	if err != nil {
		return nil, fmt.Errorf("error decompressing Zipped data: %w", err)
	}

	for _, entry := range decompressed.File {
		if entry.Name != binaryPath {
			slog.Debug("Skipping inner file on path mismatch", "innerName", entry.Name)

			continue
		}

		slog.Info("Found file in zip", "path", entry.Name)

		fileReader, err := entry.Open()
		if err != nil {
			return nil, fmt.Errorf("error reading %q from Zip file: %w", binaryPath, err)
		}

		outData, err := io.ReadAll(fileReader)
		if err != nil {
			return nil, fmt.Errorf("error reading %q from Zip file: %w", binaryPath, err)
		}

		return outData, nil
	}

	return nil, fmt.Errorf("no file named %q found in archive", binaryPath)
}

func unTarxzFileNamed(binaryPath string, data io.Reader) ([]byte, error) {
	slog.Info("Extracting file from xz archive", "path", binaryPath)

	decompressed, err := xz.NewReader(data)
	if err != nil {
		return nil, fmt.Errorf("error decompressing xz data: %w", err)
	}

	return unTar(binaryPath, decompressed) //golint:nowrap
}

func unGzip(data io.Reader) ([]byte, error) {
	slog.Info("Extracting contents of gz archive")

	decompressed, err := gzip.NewReader(data)
	if err != nil {
		return nil, fmt.Errorf("error decompressing Gzipped data: %w", err)
	}

	outData, err := io.ReadAll(decompressed)
	if err != nil {
		return nil, fmt.Errorf("error decompressing Gzipped data: %w", err)
	}

	return outData, nil
}

func unTar(binaryPath string, data io.Reader) ([]byte, error) {
	tarReader := tar.NewReader(data)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("error extracting from tar: %w", err)
		}

		if header.Typeflag == tar.TypeReg {
			if header.Name != binaryPath {
				slog.Debug("Skipping inner file on path mismatch", "innerName", header.Name)

				continue
			}

			slog.Info("Found file in tar", "path", header.Name)

			outData, err := io.ReadAll(tarReader)
			if err != nil {
				return nil, fmt.Errorf("error extracting file from tar: %w", err)
			}

			return outData, nil
		}
	}

	return nil, fmt.Errorf("no file named %q found in archive", binaryPath)
}
