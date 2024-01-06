package pkg

import (
	"bytes"
	"fmt"
	"text/template"
)

type Override struct {
	Platform     string
	Architecture string
	Extension    string
}

type Binary struct {
	Name        string
	Version     string
	TemplateUrl string
	Overrides   map[string]Override
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

	return Binary{
		Name:        config.Name,
		Version:     config.Version,
		TemplateUrl: config.Source,
		Overrides:   overrides,
	}, nil
}

func (b *Binary) GetUrl(platform, arch string) (string, error) {
	tmpl, err := template.New("sourceUrl:" + b.Name).Parse(b.TemplateUrl)
	if err != nil {
		return "", fmt.Errorf("Error parsing source template: %w", err)
	}

	vm := newUrlViewModel(b, platform, arch)

	var output bytes.Buffer
	err = tmpl.Execute(&output, vm)
	if err != nil {
		return "", fmt.Errorf("Error rendering source template: %w", err)
	}

	return output.String(), nil
}

func (b *Binary) getOveride(platform, arch string) (Override, bool) {
	key := fmt.Sprintf("%s,%s", platform, arch)

	if over, ok := b.Overrides[key]; ok {
		return over, true
	} else {
		return Override{}, false
	}
}

type urlViewModel struct {
	Version  string
	Platform string
	Arch     string
	Ext      string
}

func newUrlViewModel(binary *Binary, platform, arch string) urlViewModel {
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
