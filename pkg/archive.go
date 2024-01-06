package pkg

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"path"
)

func unTgzFileNamed(binaryName string, data io.Reader) ([]byte, error) {
	decompressed, err := gzip.NewReader(data)
	if err != nil {
		return nil, fmt.Errorf("error decompressing Gzipped data: %w", err)
	}

	tarReader := tar.NewReader(decompressed)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("error extracting from tar: %w", err)
		}

		if header.Typeflag == tar.TypeReg {
			_, archivedName := path.Split(header.Name)
			if archivedName != binaryName {
				continue
			}

			outData, err := io.ReadAll(tarReader)
			if err != nil {
				return nil, fmt.Errorf("error extracting file from tar: %w", err)
			}

			return outData, nil
		}
	}

	return nil, fmt.Errorf("no file named %q found in archive", binaryName)
}

func unZipFileNamed(binaryName string, data io.Reader) ([]byte, error) {
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
		if entry.Name == binaryName {
			fileReader, err := entry.Open()
			if err != nil {
				return nil, fmt.Errorf("error reading %q from Zip file: %w", binaryName, err)
			}

			outData, err := io.ReadAll(fileReader)
			if err != nil {
				return nil, fmt.Errorf("error reading %q from Zip file: %w", binaryName, err)
			}

			return outData, nil
		}
	}

	return nil, fmt.Errorf("no file named %q found in archive", binaryName)
}

func unGzip(data io.Reader) ([]byte, error) {
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
