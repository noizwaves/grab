package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
)

type errorBody struct {
	Message string `json:"message"`
}

func parseRelease(data []byte) (*Release, error) {
	var output Release

	err := json.Unmarshal(data, &output)
	if err != nil {
		return nil, fmt.Errorf("error parsing response as JSON: %w", err)
	}

	return &output, nil
}

func parseError(data []byte) error {
	var output errorBody

	err := json.Unmarshal(data, &output)
	if err != nil {
		return fmt.Errorf("error parsing error response as JSON: %w", err)
	}

	return fmt.Errorf("%s", output.Message)
}

type Client interface {
	GetLatestRelease(org, repo string) (*Release, error)
	GetReleaseByTag(org, repo, tag string) (*Release, error)
	DownloadReleaseAsset(org, repo, releaseName, assetName string) ([]byte, error)
}

type ClientImpl struct {
	baseURL string
}

func NewClient() *ClientImpl {
	return &ClientImpl{
		baseURL: "https://api.github.com",
	}
}

func NewClientWithBaseURL(baseURL string) *ClientImpl {
	return &ClientImpl{
		baseURL: baseURL,
	}
}

func (g *ClientImpl) GetLatestRelease(org, repo string) (*Release, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/releases/latest", g.baseURL, org, repo)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating GET request: %w", err)
	}

	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("X-Github-Api-Version", "2022-11-28")

	if token := os.Getenv("GH_TOKEN"); token != "" {
		req.Header.Add("Authorization", "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return parseRelease(data)
	default:
		return nil, parseError(data)
	}
}

func (g *ClientImpl) GetReleaseByTag(org, repo, tag string) (*Release, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/releases/tags/%s", g.baseURL, org, repo, tag)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating GET request: %w", err)
	}

	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("X-Github-Api-Version", "2022-11-28")

	if token := os.Getenv("GH_TOKEN"); token != "" {
		req.Header.Add("Authorization", "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return parseRelease(data)
	default:
		return nil, parseError(data)
	}
}

func (g *ClientImpl) DownloadReleaseAsset(org, repo, release, asset string) ([]byte, error) {
	url := fmt.Sprintf("https://github.com/%s/%s/releases/download/%s/%s",
		org, repo, release, asset)

	slog.Debug("Downloading asset from GitHub", "url", url)

	return downloadArtifact(url)
}

func downloadArtifact(url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error requesting asset: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	return data, nil
}
