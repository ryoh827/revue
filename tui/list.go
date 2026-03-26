package tui

import (
	"fmt"
	"strings"
	"time"

	gh "github.com/ryoh827/revue/github"
)

// PRItem wraps a PullRequest with its CI status for display.
type PRItem struct {
	PR     gh.PullRequest
	CI     gh.CheckStatus
	CIErr  error
}

func renderPRList(items []PRItem, cursor int, width int) string {
	if len(items) == 0 {
		return repoStyle.Render("  No pull requests found.")
	}

	var b strings.Builder
	for i, item := range items {
		prefix := "  "
		style := normalStyle
		if i == cursor {
			prefix = "▸ "
			style = selectedStyle
		}

		ci := renderCI(item.CI)
		repo := repoStyle.Render(item.PR.Repo)
		title := item.PR.Title
		if item.PR.Draft {
			title = draftStyle.Render("[draft] " + title)
		}
		age := formatAge(item.PR.UpdatedAt)

		line := fmt.Sprintf("%s%s %s #%d %s  %s",
			prefix, ci, repo, item.PR.Number, style.Render(title), repoStyle.Render(age))

		// Truncate to terminal width
		if width > 0 && len(line) > width {
			line = line[:width-1] + "…"
		}

		b.WriteString(line)
		b.WriteString("\n")
	}
	return b.String()
}

func renderCI(status gh.CheckStatus) string {
	icon := status.String()
	switch status {
	case gh.CheckStatusSuccess:
		return ciSuccessStyle.Render(icon)
	case gh.CheckStatusFailure:
		return ciFailureStyle.Render(icon)
	case gh.CheckStatusPending:
		return ciPendingStyle.Render(icon)
	default:
		return ciNoneStyle.Render(icon)
	}
}

func formatAge(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	default:
		days := int(d.Hours() / 24)
		if days == 1 {
			return "1d ago"
		}
		return fmt.Sprintf("%dd ago", days)
	}
}
