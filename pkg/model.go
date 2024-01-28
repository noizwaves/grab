package pkg

import (
	"bytes"
	"fmt"
	"log/slog"
	"regexp"
	"text/template"
)

type Binary struct {
	Name    string
	Version string

	// source
	Org  string
	Repo string
	// release name template
	ReleaseName  string
	ReleaseRegex *regexp.Regexp
	// (platform,arch) -> filename template
	FileName map[string]string

	// program related fields
	VersionArgs  []string
	VersionRegex *regexp.Regexp
}

func NewBinary(name, version string, config configPackage) (*Binary, error) {
	versionRegex, err := regexp.Compile(config.Spec.Program.VersionRegex)
	if err != nil {
		return nil, fmt.Errorf("version regex does not compile: %w", err)
	}

	releaseRegex, err := regexp.Compile(config.Spec.GitHubRelease.VersionRegex)
	if err != nil {
		return nil, fmt.Errorf("release regex does not compile: %w", err)
	}

	return &Binary{
		Name:    name,
		Version: version,
		// package
		Org:          config.Spec.GitHubRelease.Org,
		Repo:         config.Spec.GitHubRelease.Repo,
		ReleaseName:  config.Spec.GitHubRelease.Name,
		ReleaseRegex: releaseRegex,
		FileName:     config.Spec.GitHubRelease.FileName,
		// program
		VersionArgs:  config.Spec.Program.VersionArgs,
		VersionRegex: versionRegex,
	}, nil
}

func (b *Binary) GetAssetFileName(platform, arch string) (string, error) {
	key := platform + "," + arch
	fileNameTmplStr, ok := b.FileName[key]
	if !ok {
		return "", fmt.Errorf("filename missing for platform,arch of %q", key)
	}

	tmpl, err := template.New("filename:" + b.Name).Parse(fileNameTmplStr)
	if err != nil {
		return "", fmt.Errorf("error parsing asset filename template: %w", err)
	}

	vm := newURLViewModel(b)

	var output bytes.Buffer
	err = tmpl.Execute(&output, vm)
	if err != nil {
		return "", fmt.Errorf("error rendering asset filename template: %w", err)
	}

	return output.String(), nil
}

func (b *Binary) GetReleaseName() (string, error) {
	tmpl, err := template.New("releaseName:" + b.Name).Parse(b.ReleaseName)
	if err != nil {
		return "", fmt.Errorf("error parsing release name template: %w", err)
	}

	vm := newURLViewModel(b)

	var output bytes.Buffer
	err = tmpl.Execute(&output, vm)
	if err != nil {
		return "", fmt.Errorf("error rendering release name template: %w", err)
	}

	return output.String(), nil
}

func (b *Binary) ShouldReplace(currentVersion string) bool {
	result := b.Version != currentVersion
	slog.Info("Checking if installed binary should be replaced", "name", b.Name, "replace", result)
	slog.Debug("Version information", "name", b.Name, "current", currentVersion, "desired", b.Version)

	return result
}

type urlViewModel struct {
	Version string
}

func newURLViewModel(binary *Binary) urlViewModel {
	return urlViewModel{
		Version: binary.Version,
	}
}
