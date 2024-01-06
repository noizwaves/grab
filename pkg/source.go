package pkg

import (
	"bytes"
	"fmt"
	"runtime"
	"text/template"
)

type urlViewModel struct {
	Version  string
	Platform string
	Arch     string
	Ext      string
}

func newUrlViewModel(binary configBinary) urlViewModel {
	version := binary.Version
	platform := runtime.GOOS
	arch := runtime.GOARCH
	ext := ""

	if pOverrides, ok := binary.Platforms[runtime.GOOS]; ok {
		if aOverrides, ok := pOverrides[runtime.GOARCH]; ok {
			platform = aOverrides[0]
			arch = aOverrides[1]
			ext = aOverrides[2]
		}
	}

	return urlViewModel{
		Version:  version,
		Platform: platform,
		Arch:     arch,
		Ext:      ext,
	}
}

func renderSourceUrl(binary configBinary) (string, error) {
	tmpl, err := template.New("sourceUrl" + binary.Name).Parse(binary.Source)
	if err != nil {
		return "", fmt.Errorf("Error parsing Source as template: %w", err)
	}

	vm := newUrlViewModel(binary)

	var output bytes.Buffer

	err = tmpl.Execute(&output, vm)
	if err != nil {
		return "", fmt.Errorf("Error rendering Source as template: %w", err)
	}

	return output.String(), nil
}
