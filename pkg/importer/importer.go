package importer

import (
	"context"
	"fmt"
	"io"
	"log/slog"

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

func (i *Importer) Import(gCtx *pkg.GrabContext, url string, out io.Writer) error {
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
	detectedPackage, err := detectPackage(release)
	if err != nil {
		return err
	}

	// Assume binary name matches repository name verbatim
	packageName := releaseURL.Repository

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
				VersionRegex: "\\d+\\.\\d+\\.\\d+",
				FileName:     detectedPackage.assets,

				// Assume binary is at the root of any archive file
				EmbeddedBinaryPath: nil,
			},
			Program: pkg.ConfigProgram{
				// Assume binary uses a --version flag
				VersionArgs:  []string{"--version"},
				VersionRegex: "\\d+\\.\\d+\\.\\d+",
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
	releaseName string
	assets      map[string]string
}

func detectPackage(release *github.Release) (*detectedPackage, error) {
	releaseDetector := NewReleaseNamePatternDetector()

	releasePattern, err := releaseDetector.AnalyzeOne(release.TagName)
	if err != nil {
		return nil, err
	}

	releaseName := releasePattern.Value

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

			// this asset name matches the current pair!
			result = asset.Name

			break
		}

		if result == "" {
			return nil, fmt.Errorf("no matching asset name found for platform %s and architecture %s", platform, arch)
		}

		// parse this asset name for a version number

		key := fmt.Sprintf("%s,%s", platform, arch)
		fileNames[key] = result
	}

	return &detectedPackage{
		releaseName: releaseName,
		assets:      fileNames,
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
