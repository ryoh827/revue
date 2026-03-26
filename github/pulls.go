package github

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

// PullRequest represents a GitHub pull request.
type PullRequest struct {
	Number    int       `json:"number"`
	Title     string    `json:"title"`
	HTMLURL   string    `json:"html_url"`
	State     string    `json:"state"`
	Draft     bool      `json:"draft"`
	User      User      `json:"user"`
	UpdatedAt time.Time `json:"updated_at"`
	Repo      string    // "owner/repo" — populated after fetch
	Head      Head      `json:"head"`
}

// User represents a GitHub user.
type User struct {
	Login string `json:"login"`
}

// Head represents the head ref of a PR.
type Head struct {
	SHA string `json:"sha"`
	Ref string `json:"ref"`
}

type searchResult struct {
	Items []PullRequest `json:"items"`
}

// FetchReviewRequests returns open PRs where the authenticated user is requested as reviewer.
func (c *Client) FetchReviewRequests() ([]PullRequest, error) {
	query := fmt.Sprintf("type:pr state:open review-requested:%s", c.username)
	return c.searchPRs(query)
}

// FetchMyPRs returns open PRs authored by the authenticated user.
func (c *Client) FetchMyPRs() ([]PullRequest, error) {
	query := fmt.Sprintf("type:pr state:open author:%s", c.username)
	return c.searchPRs(query)
}

func (c *Client) searchPRs(query string) ([]PullRequest, error) {
	var all []PullRequest
	for page := 1; ; page++ {
		u := fmt.Sprintf("https://api.github.com/search/issues?q=%s&sort=updated&order=desc&per_page=50&page=%d",
			url.QueryEscape(query), page)

		resp, err := c.get(u)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != 200 {
			resp.Body.Close()
			return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
		}

		var result searchResult
		err = json.NewDecoder(resp.Body).Decode(&result)
		resp.Body.Close()
		if err != nil {
			return nil, err
		}

		for i := range result.Items {
			result.Items[i].Repo = extractRepo(result.Items[i].HTMLURL)
		}
		all = append(all, result.Items...)

		if len(result.Items) < 50 {
			break
		}
	}
	return all, nil
}

// extractRepo extracts "owner/repo" from a PR HTML URL.
func extractRepo(htmlURL string) string {
	u, err := url.Parse(htmlURL)
	if err != nil {
		return ""
	}
	// path: /owner/repo/pull/123
	parts := splitPath(u.Path)
	if len(parts) >= 2 {
		return parts[0] + "/" + parts[1]
	}
	return ""
}

func splitPath(path string) []string {
	var parts []string
	for _, p := range splitString(path, '/') {
		if p != "" {
			parts = append(parts, p)
		}
	}
	return parts
}

func splitString(s string, sep byte) []string {
	var result []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == sep {
			result = append(result, s[start:i])
			start = i + 1
		}
	}
	result = append(result, s[start:])
	return result
}
