package github

type Release struct {
	Name string `json:"name"`
	URL  string `json:"html_url"` //nolint:tagliatelle
}
