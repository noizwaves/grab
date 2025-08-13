package importer

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/noizwaves/grab/pkg"
	"github.com/noizwaves/grab/pkg/github"
)

type Importer struct {
	githubClient github.Client
}

func NewImporter(githubClient github.Client) *Importer {
	return &Importer{
		githubClient: githubClient,
	}
}

func (i *Importer) Import(gCtx *pkg.GrabContext, url string, customPackageName string, out io.Writer) error {
	releaseURL, err := ParseGitHubReleaseURL(url)
	if err != nil {
		return err
	}

	ctx := context.Background()
	slog.InfoContext(ctx, "Importing GitHub Release", "org", releaseURL.Organization, "repo", releaseURL.Repository)

	release, err := i.githubClient.GetLatestRelease(releaseURL.Organization, releaseURL.Repository)
	if err != nil {
		return fmt.Errorf("failed to get release: %w", err)
	}

	// Detect the patterns from the Release
	packageName := releaseURL.Repository
	if customPackageName != "" {
		packageName = customPackageName
	}

	detectedPackage, err := detectPackage(
		i.githubClient,
		releaseURL.Organization,
		releaseURL.Repository,
		release,
		packageName,
	)
	if err != nil {
		return err
	}

	// Construct the new binary
	packageConfig := pkg.ConfigPackage{
		APIVersion: "grab.noizwaves.com/v1alpha1",
		Kind:       "Package",
		Metadata: pkg.ConfigPackageMetadata{
			Name: packageName,
		},
		Spec: pkg.ConfigPackageSpec{
			GitHubRelease: pkg.ConfigGitHubRelease{
				Org:          releaseURL.Organization,
				Repo:         releaseURL.Repository,
				Name:         detectedPackage.releaseName,
				VersionRegex: detectedPackage.versionRegex,
				FileName:     detectedPackage.assets,

				// Use detected embedded binary paths
				EmbeddedBinaryPath: detectedPackage.embeddedBinaryPaths,
			},
			Program: pkg.ConfigProgram{
				// Assume binary uses a --version flag and not a subcommand
				VersionArgs: []string{"--version"},
				// Assume binary printed version regex matches tag regex
				VersionRegex: detectedPackage.versionRegex,
			},
		},
	}

	packagePath, err := gCtx.SavePackage(&packageConfig)
	if err != nil {
		return fmt.Errorf("failed to save package: %w", err)
	}

	fmt.Fprintf(out, "Package %q saved to %s\n", packageName, packagePath)

	return nil
}

type detectedPackage struct {
	releaseName         string
	versionRegex        string
	assets              map[string]string
	embeddedBinaryPaths map[string]string
}

//nolint:funlen,lll
func detectPackage(ghClient github.Client, org, repo string, release *github.Release, packageName string) (*detectedPackage, error) {
	// Release name pattern
	releaseDetector := NewReleaseNamePatternDetector()

	releasePattern, err := releaseDetector.AnalyzeOne(release.TagName)
	if err != nil {
		return nil, err
	}

	releaseName := releasePattern.Value

	// Version regex
	versionRegex := DetectVersionRegex(release.TagName)

	// Current version value
	latestVersion, err := DetectVersionValue(versionRegex, release.TagName)
	if err != nil {
		return nil, fmt.Errorf("failed to detect version: %w", err)
	}

	fileNames := make(map[string]string)

	// analyze the asset names for all required platform+architecture pairs
	detectablePairs := getDetectablePairs()

	for _, pair := range detectablePairs {
		platform := pair[0]
		arch := pair[1]
		platformDetector := NewPlatformPatternDetector(platform)
		archDetector := NewArchitecturePatternDetector(arch)

		var result string

		for _, asset := range release.Assets {
			_, err := platformDetector.AnalyzeOne(asset.Name)
			if err != nil {
				continue
			}

			_, err = archDetector.AnalyzeOne(asset.Name)
			if err != nil {
				continue
			}

			// this asset name matches the current pair
			result = asset.Name

			break
		}

		if result == "" {
			return nil, fmt.Errorf("no matching asset name found for platform %s and architecture %s", platform, arch)
		}

		// Convert to a template string if needed
		result = UnrenderVersionValue(result, latestVersion)

		key := fmt.Sprintf("%s,%s", platform, arch)
		fileNames[key] = result
	}

	// Detect embedded binary paths for archive assets
	embeddedPaths, err := detectEmbeddedBinaryPaths(ghClient, org, repo, release, packageName, fileNames, latestVersion)
	if err != nil {
		slog.WarnContext(context.Background(), "Failed to detect embedded binary paths", "error", err)
		// Continue without embedded paths rather than failing completely
	}

	var embeddedBinaryPathsMap map[string]string
	if embeddedPaths != nil {
		embeddedBinaryPathsMap = *embeddedPaths
	}

	return &detectedPackage{
		releaseName:         releaseName,
		assets:              fileNames,
		versionRegex:        versionRegex,
		embeddedBinaryPaths: embeddedBinaryPathsMap,
	}, nil
}

func getDetectablePairs() [][]string {
	return [][]string{
		{"linux", "amd64"},
		{"linux", "arm64"},
		{"darwin", "amd64"},
		{"darwin", "arm64"},
	}
}

func detectEmbeddedBinaryPaths(
	ghClient github.Client,
	org, repo string,
	release *github.Release,
	packageName string,
	detectedAssets map[string]string,
	versionLiteral string,
) (*map[string]string, error) {
	ctx := context.Background()
	slog.InfoContext(ctx, "Detecting embedded binary paths", "package", packageName)

	embeddedPaths := make(map[string]string)

	for platformArch, assetName := range detectedAssets {
		// Render the asset name template with version
		renderedAssetName, err := renderAssetNameTemplate(assetName, versionLiteral)
		if err != nil {
			return nil, fmt.Errorf("failed to render asset name template %s: %w", assetName, err)
		}
		// Skip non-archive assets - they don't need embedded paths
		if !isArchiveAsset(renderedAssetName) {
			continue
		}

		slog.DebugContext(ctx, "Analyzing archive asset", "platformArch", platformArch, "asset", renderedAssetName)

		// Download the asset
		data, err := ghClient.DownloadReleaseAsset(org, repo, release.TagName, renderedAssetName)
		if err != nil {
			return nil, fmt.Errorf("failed to download asset %s for binary detection: %w", renderedAssetName, err)
		}

		// List archive contents
		files, err := listArchiveContents(renderedAssetName, bytes.NewBuffer(data))
		if err != nil {
			return nil, fmt.Errorf("failed to list archive contents for %s: %w", renderedAssetName, err)
		}

		// Find binary matching package name
		binaryPath, err := findBinaryInArchive(files, packageName)
		if err != nil {
			return nil, fmt.Errorf("failed to find binary in asset %s: %w", renderedAssetName, err)
		}

		// Skip if binary path is just the package name (default path)
		if binaryPath == packageName {
			continue
		}

		// Template the binary path by replacing version literals
		templatedPath := UnrenderVersionValue(binaryPath, versionLiteral)
		embeddedPaths[platformArch] = templatedPath
		slog.InfoContext(ctx, "Detected embedded binary path", "platformArch", platformArch, "path", binaryPath)
	}

	// Return nil if no embedded paths were detected
	if len(embeddedPaths) == 0 {
		//nolint:nilnil
		return nil, nil
	}

	return &embeddedPaths, nil
}

func isArchiveAsset(assetName string) bool {
	return strings.HasSuffix(assetName, ".tar.gz") ||
		strings.HasSuffix(assetName, ".tgz") ||
		strings.HasSuffix(assetName, ".tar.xz") ||
		strings.HasSuffix(assetName, ".zip")
}

func findBinaryInArchive(files []string, packageName string) (string, error) {
	var candidates []string

	for _, path := range files {
		// Skip directories
		if strings.HasSuffix(path, "/") {
			continue
		}

		// Get fileName without directory path
		fileName := filepath.Base(path)

		// Exact match takes priority
		if fileName == packageName {
			return path, nil
		}

		// Collect potential matches for fallback
		if strings.Contains(fileName, packageName) {
			candidates = append(candidates, path)
		}
	}

	// If no exact match, return first candidate
	if len(candidates) > 0 {
		return candidates[0], nil
	}

	return "", errors.New("no binary found matching package name \"" + packageName + "\"")
}

//nolint:wrapcheck
func listArchiveContents(assetName string, data *bytes.Buffer) ([]string, error) {
	switch {
	case strings.HasSuffix(assetName, ".tar.gz") || strings.HasSuffix(assetName, ".tgz"):
		return pkg.ListTgzContents(data)
	case strings.HasSuffix(assetName, ".tar.xz"):
		return pkg.ListTarxzContents(data)
	case strings.HasSuffix(assetName, ".zip"):
		return pkg.ListZipContents(data)
	default:
		return nil, errors.New("unsupported archive format: " + assetName)
	}
}

type templateViewModel struct {
	Version string
}

func renderAssetNameTemplate(assetNameTemplate, versionLiteral string) (string, error) {
	tmpl, err := template.New("assetName").Parse(assetNameTemplate)
	if err != nil {
		return "", fmt.Errorf("error parsing asset name template: %w", err)
	}

	viewModel := templateViewModel{
		Version: versionLiteral,
	}

	var output bytes.Buffer

	err = tmpl.Execute(&output, viewModel)
	if err != nil {
		return "", fmt.Errorf("error rendering asset name template: %w", err)
	}

	return output.String(), nil
}
