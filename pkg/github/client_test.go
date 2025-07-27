package github

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetReleaseByTag_Success(t *testing.T) {
	// Mock GitHub API response
	mockRelease := Release{
		Name: "v1.2.3",
		URL:  "https://github.com/owner/repo/releases/tag/v1.2.3",
		Assets: []Asset{
			{
				Name:        "app-linux-amd64.tar.gz",
				Size:        1024,
				DownloadURL: "https://github.com/owner/repo/releases/download/v1.2.3/app-linux-amd64.tar.gz",
			},
			{
				Name:        "app-darwin-amd64.tar.gz",
				Size:        2048,
				DownloadURL: "https://github.com/owner/repo/releases/download/v1.2.3/app-darwin-amd64.tar.gz",
			},
		},
	}

	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		// Verify request
		if request.URL.Path != "/repos/owner/repo/releases/tags/v1.2.3" {
			t.Errorf("Expected path '/repos/owner/repo/releases/tags/v1.2.3', got '%s'", request.URL.Path)
		}

		if request.Header.Get("Accept") != "application/vnd.github+json" {
			t.Errorf("Expected Accept header 'application/vnd.github+json', got '%s'", request.Header.Get("Accept"))
		}

		// Return mock response
		responseWriter.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(responseWriter).Encode(mockRelease)
	}))
	defer server.Close()

	// Create client pointing to mock server
	client := NewClientWithBaseURL(server.URL)

	result, err := client.GetReleaseByTag("owner", "repo", "v1.2.3")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result.Name != "v1.2.3" {
		t.Errorf("Expected release name 'v1.2.3', got '%s'", result.Name)
	}

	if len(result.Assets) != 2 {
		t.Errorf("Expected 2 assets, got %d", len(result.Assets))
	}

	if result.Assets[0].Name != "app-linux-amd64.tar.gz" {
		t.Errorf("Expected first asset name 'app-linux-amd64.tar.gz', got '%s'", result.Assets[0].Name)
	}

	if result.Assets[0].Size != 1024 {
		t.Errorf("Expected first asset size 1024, got %d", result.Assets[0].Size)
	}
}

func TestGetReleaseByTag_NotFound(t *testing.T) {
	// Create mock server that returns 404
	server := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, _ *http.Request) {
		responseWriter.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(responseWriter).Encode(map[string]string{
			"message": "Not Found",
		})
	}))
	defer server.Close()

	// Create client pointing to mock server
	client := NewClientWithBaseURL(server.URL)

	_, err := client.GetReleaseByTag("owner", "repo", "nonexistent")
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err.Error() != "Not Found" {
		t.Errorf("Expected error message 'Not Found', got '%s'", err.Error())
	}
}

func TestGetLatestRelease_Success(t *testing.T) {
	// Mock GitHub API response
	mockRelease := Release{
		Name: "v2.1.0",
		URL:  "https://github.com/owner/repo/releases/tag/v2.1.0",
		Assets: []Asset{
			{
				Name:        "app-latest-linux-amd64.tar.gz",
				Size:        2048,
				DownloadURL: "https://github.com/owner/repo/releases/download/v2.1.0/app-latest-linux-amd64.tar.gz",
			},
		},
	}

	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		// Verify request
		if request.URL.Path != "/repos/owner/repo/releases/latest" {
			t.Errorf("Expected path '/repos/owner/repo/releases/latest', got '%s'", request.URL.Path)
		}

		if request.Header.Get("Accept") != "application/vnd.github+json" {
			t.Errorf("Expected Accept header 'application/vnd.github+json', got '%s'", request.Header.Get("Accept"))
		}

		// Return mock response
		responseWriter.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(responseWriter).Encode(mockRelease)
	}))
	defer server.Close()

	// Create client pointing to mock server
	client := NewClientWithBaseURL(server.URL)

	result, err := client.GetLatestRelease("owner", "repo")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result.Name != "v2.1.0" {
		t.Errorf("Expected release name 'v2.1.0', got '%s'", result.Name)
	}

	if len(result.Assets) != 1 {
		t.Errorf("Expected 1 asset, got %d", len(result.Assets))
	}

	if result.Assets[0].Name != "app-latest-linux-amd64.tar.gz" {
		t.Errorf("Expected asset name 'app-latest-linux-amd64.tar.gz', got '%s'", result.Assets[0].Name)
	}

	if result.Assets[0].Size != 2048 {
		t.Errorf("Expected asset size 2048, got %d", result.Assets[0].Size)
	}
}

func TestGetLatestRelease_NoReleases(t *testing.T) {
	// Create mock server that returns 404 for no releases
	server := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, _ *http.Request) {
		responseWriter.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(responseWriter).Encode(map[string]string{
			"message": "Not Found",
		})
	}))
	defer server.Close()

	// Create client pointing to mock server
	client := NewClientWithBaseURL(server.URL)

	_, err := client.GetLatestRelease("owner", "repo")
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err.Error() != "Not Found" {
		t.Errorf("Expected error message 'Not Found', got '%s'", err.Error())
	}
}
