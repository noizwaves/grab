package pkg

import (
	"bytes"
	"fmt"
	"regexp"
	"text/template"
)

type Override struct {
	Platform     string
	Architecture string
	Extension    string
}

type Binary struct {
	Name         string
	Version      string
	TemplateURL  string
	Overrides    map[string]Override
	VersionFlags []string
	VersionRegex *regexp.Regexp
}

func NewBinary(config configBinary) (Binary, error) {
	overrides := make(map[string]Override)
	for platform, pOver := range config.Platforms {
		for arch, over := range pOver {
			key := fmt.Sprintf("%s,%s", platform, arch)
			overrides[key] = Override{
				Platform:     over[0],
				Architecture: over[1],
				Extension:    over[2],
			}
		}
	}

	versionRegex, err := regexp.Compile(config.VersionRegex)
	if err != nil {
		return Binary{}, fmt.Errorf("version regex does not compile: %w", err)
	}

	return Binary{
		Name:         config.Name,
		Version:      config.Version,
		TemplateURL:  config.Source,
		Overrides:    overrides,
		VersionFlags: config.VersionFlags,
		VersionRegex: versionRegex,
	}, nil
}

func (b *Binary) GetURL(platform, arch string) (string, error) {
	tmpl, err := template.New("sourceUrl:" + b.Name).Parse(b.TemplateURL)
	if err != nil {
		return "", fmt.Errorf("error parsing source template: %w", err)
	}

	vm := newURLViewModel(b, platform, arch)

	var output bytes.Buffer
	err = tmpl.Execute(&output, vm)
	if err != nil {
		return "", fmt.Errorf("error rendering source template: %w", err)
	}

	return output.String(), nil
}

func (b *Binary) getOveride(platform, arch string) (Override, bool) {
	key := fmt.Sprintf("%s,%s", platform, arch)

	if over, ok := b.Overrides[key]; ok {
		return over, true
	}

	return Override{}, false
}

type urlViewModel struct {
	Version  string
	Platform string
	Arch     string
	Ext      string
}

func newURLViewModel(binary *Binary, platform, arch string) urlViewModel {
	ext := ""
	if over, ok := binary.getOveride(platform, arch); ok {
		platform = over.Platform
		arch = over.Architecture
		ext = over.Extension
	}

	return urlViewModel{
		Version:  binary.Version,
		Platform: platform,
		Arch:     arch,
		Ext:      ext,
	}
}
