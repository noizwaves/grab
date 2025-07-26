package githubh

import (
	"errors"

	"github.com/noizwaves/grab/pkg/github"
)

type MockGitHubClient struct {
	AssetData []byte
	Release   *github.Release

	// Call tracking
	GetLatestReleaseCalls []GetLatestReleaseCall
}

type GetLatestReleaseCall struct {
	Org  string
	Repo string
}

func (m *MockGitHubClient) DownloadReleaseAsset(_, _, _, _ string) ([]byte, error) {
	if len(m.AssetData) == 0 {
		return nil, errors.New("not implemented")
	}

	return m.AssetData, nil
}

func (m *MockGitHubClient) GetLatestRelease(org, repo string) (*github.Release, error) {
	// Track the call
	m.GetLatestReleaseCalls = append(m.GetLatestReleaseCalls, GetLatestReleaseCall{
		Org:  org,
		Repo: repo,
	})

	if m.Release == nil {
		return nil, errors.New("not implemented")
	}

	return m.Release, nil
}
