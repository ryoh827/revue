package github

import (
	"encoding/json"
	"fmt"
)

// CheckStatus represents the aggregated CI status of a PR.
type CheckStatus int

const (
	CheckStatusNone    CheckStatus = iota // No checks
	CheckStatusPending                    // Some checks still running
	CheckStatusSuccess                    // All checks passed
	CheckStatusFailure                    // Some checks failed
)

// String returns a display icon for the check status.
func (s CheckStatus) String() string {
	switch s {
	case CheckStatusSuccess:
		return "✓"
	case CheckStatusFailure:
		return "✗"
	case CheckStatusPending:
		return "⏳"
	default:
		return "–"
	}
}

type checkRunsResponse struct {
	TotalCount int        `json:"total_count"`
	CheckRuns  []checkRun `json:"check_runs"`
}

type checkRun struct {
	Status     string  `json:"status"`     // queued, in_progress, completed
	Conclusion *string `json:"conclusion"` // success, failure, neutral, cancelled, skipped, timed_out, action_required
}

// FetchCheckStatus returns the aggregated CI check status for a given commit SHA.
func (c *Client) FetchCheckStatus(owner, repo, sha string) (CheckStatus, error) {
	u := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits/%s/check-runs?per_page=100", owner, repo, sha)

	resp, err := c.get(u)
	if err != nil {
		return CheckStatusNone, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return CheckStatusNone, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var result checkRunsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return CheckStatusNone, err
	}

	if result.TotalCount == 0 {
		return CheckStatusNone, nil
	}

	hasFailure := false
	hasPending := false
	for _, cr := range result.CheckRuns {
		if cr.Status != "completed" {
			hasPending = true
			continue
		}
		if cr.Conclusion != nil {
			switch *cr.Conclusion {
			case "failure", "timed_out", "action_required":
				hasFailure = true
			case "cancelled":
				hasFailure = true
			}
		}
	}

	if hasFailure {
		return CheckStatusFailure, nil
	}
	if hasPending {
		return CheckStatusPending, nil
	}
	return CheckStatusSuccess, nil
}

// FetchPRDetail fetches the full PR object (including head SHA) from the pulls API.
func (c *Client) FetchPRDetail(owner, repo string, number int) (*PullRequest, error) {
	u := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls/%d", owner, repo, number)

	resp, err := c.get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var pr PullRequest
	if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
		return nil, err
	}
	pr.Repo = owner + "/" + repo
	return &pr, nil
}
