package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

// Client is a GitHub API client authenticated via gh CLI.
type Client struct {
	httpClient *http.Client
	token      string
	username   string
}

// NewClient creates a new GitHub API client using gh auth token.
func NewClient() (*Client, error) {
	token, err := ghAuthToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get gh auth token: %w\nMake sure gh CLI is installed and authenticated (gh auth login)", err)
	}

	c := &Client{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		token:      strings.TrimSpace(token),
	}

	username, err := c.fetchUsername()
	if err != nil {
		return nil, fmt.Errorf("failed to get authenticated user: %w", err)
	}
	c.username = username

	return c, nil
}

// Username returns the authenticated user's login.
func (c *Client) Username() string {
	return c.username
}

func ghAuthToken() (string, error) {
	out, err := exec.Command("gh", "auth", "token").Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	return c.httpClient.Do(req)
}

func (c *Client) get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return c.do(req)
}

func (c *Client) fetchUsername() (string, error) {
	resp, err := c.get("https://api.github.com/user")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var user struct {
		Login string `json:"login"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return "", err
	}
	return user.Login, nil
}
