package github

import (
	"context"
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	out, err := exec.CommandContext(ctx, "gh", "auth", "token").Output()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("gh auth token timed out after 5s")
		}
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
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

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var user struct {
		Login string `json:"login"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return "", err
	}
	if user.Login == "" {
		return "", fmt.Errorf("authenticated user login is empty")
	}
	return user.Login, nil
}
