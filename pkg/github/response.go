package github

// Release represents a GitHub release with asset information.
type Release struct {
	Name    string  `json:"name"`
	URL     string  `json:"html_url"` //nolint:tagliatelle
	TagName string  `json:"tag_name"` //nolint:tagliatelle
	Assets  []Asset `json:"assets"`
}

// Asset represents a GitHub release asset.
type Asset struct {
	Name        string `json:"name"`
	Size        int64  `json:"size"`
	DownloadURL string `json:"browser_download_url"` //nolint:tagliatelle
}
