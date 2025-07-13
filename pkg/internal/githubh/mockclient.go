package githubh

import (
	"errors"

	"github.com/noizwaves/grab/pkg/github"
)

type MockGitHubClient struct {
	AssetData []byte
	Release   *github.Release
}

func (m *MockGitHubClient) DownloadReleaseAsset(_, _, _, _ string) ([]byte, error) {
	if len(m.AssetData) == 0 {
		return nil, errors.New("not implemented")
	}

	return m.AssetData, nil
}

func (m *MockGitHubClient) GetLatestRelease(_, _ string) (*github.Release, error) {
	if m.Release == nil {
		return nil, errors.New("not implemented")
	}

	return m.Release, nil
}
