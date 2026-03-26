package tui

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	gh "github.com/ryoh827/revue/github"
)

// Tab represents a view tab.
type Tab int

const (
	TabReviewRequested Tab = iota
	TabMyPRs
)

var tabLabels = []string{"Review Requested", "My PRs"}

// Model is the Bubble Tea model for the revue TUI.
type Model struct {
	client *gh.Client
	width  int
	height int

	activeTab Tab
	cursors   [2]int // cursor per tab

	reviewItems []PRItem
	myPRItems   []PRItem

	loading bool
	err     error
}

// Messages

type fetchDoneMsg struct {
	reviewItems []PRItem
	myPRItems   []PRItem
	err         error
}

// New creates a new Model.
func New(client *gh.Client) Model {
	return Model{
		client:  client,
		loading: true,
	}
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	return m.fetchAll
}

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)

	case fetchDoneMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.reviewItems = msg.reviewItems
		m.myPRItems = msg.myPRItems
		m.err = nil
		return m, nil
	}

	return m, nil
}

// View implements tea.Model.
func (m Model) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("revue"))
	b.WriteString("\n")

	// Tabs
	var tabs []string
	for i, label := range tabLabels {
		if Tab(i) == m.activeTab {
			tabs = append(tabs, activeTabStyle.Render(label))
		} else {
			tabs = append(tabs, inactiveTabStyle.Render(label))
		}
	}
	b.WriteString(tabBarStyle.Render(strings.Join(tabs, " ")))
	b.WriteString("\n")

	if m.loading {
		b.WriteString(loadingStyle.Render("  Loading..."))
		b.WriteString("\n")
	} else if m.err != nil {
		b.WriteString(errorStyle.Render(fmt.Sprintf("  Error: %v", m.err)))
		b.WriteString("\n")
	} else {
		switch m.activeTab {
		case TabReviewRequested:
			b.WriteString(renderPRList(m.reviewItems, m.cursors[TabReviewRequested], m.width))
		case TabMyPRs:
			b.WriteString(renderPRList(m.myPRItems, m.cursors[TabMyPRs], m.width))
		}
	}

	b.WriteString(helpStyle.Render("  j/k: navigate  tab/h/l: switch tab  enter: open  r: refresh  q: quit"))
	b.WriteString("\n")

	return b.String()
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	action := parseKey(msg)
	switch action {
	case keyQuit:
		return m, tea.Quit
	case keyUp:
		if m.cursors[m.activeTab] > 0 {
			m.cursors[m.activeTab]--
		}
	case keyDown:
		max := m.currentListLen() - 1
		if max < 0 {
			max = 0
		}
		if m.cursors[m.activeTab] < max {
			m.cursors[m.activeTab]++
		}
	case keyTab:
		if m.activeTab == TabReviewRequested {
			m.activeTab = TabMyPRs
		} else {
			m.activeTab = TabReviewRequested
		}
	case keyEnter:
		if item := m.currentItem(); item != nil {
			openBrowser(item.PR.HTMLURL)
		}
	case keyRefresh:
		m.loading = true
		return m, m.fetchAll
	}
	return m, nil
}

func (m Model) currentListLen() int {
	switch m.activeTab {
	case TabReviewRequested:
		return len(m.reviewItems)
	case TabMyPRs:
		return len(m.myPRItems)
	}
	return 0
}

func (m Model) currentItem() *PRItem {
	var items []PRItem
	switch m.activeTab {
	case TabReviewRequested:
		items = m.reviewItems
	case TabMyPRs:
		items = m.myPRItems
	}
	idx := m.cursors[m.activeTab]
	if idx >= 0 && idx < len(items) {
		return &items[idx]
	}
	return nil
}

func (m Model) fetchAll() tea.Msg {
	reviewPRs, err := m.client.FetchReviewRequests()
	if err != nil {
		return fetchDoneMsg{err: err}
	}

	myPRs, err := m.client.FetchMyPRs()
	if err != nil {
		return fetchDoneMsg{err: err}
	}

	reviewItems := m.enrichWithCI(reviewPRs)
	myPRItems := m.enrichWithCI(myPRs)

	return fetchDoneMsg{
		reviewItems: reviewItems,
		myPRItems:   myPRItems,
	}
}

func (m Model) enrichWithCI(prs []gh.PullRequest) []PRItem {
	items := make([]PRItem, len(prs))
	for i, pr := range prs {
		items[i] = PRItem{PR: pr}

		// Need to fetch full PR details to get head SHA (search API doesn't include it)
		parts := strings.SplitN(pr.Repo, "/", 2)
		if len(parts) != 2 {
			continue
		}
		detail, err := m.client.FetchPRDetail(parts[0], parts[1], pr.Number)
		if err != nil {
			items[i].CIErr = err
			continue
		}
		items[i].PR.Head = detail.Head

		status, err := m.client.FetchCheckStatus(parts[0], parts[1], detail.Head.SHA)
		if err != nil {
			items[i].CIErr = err
			continue
		}
		items[i].CI = status
	}
	return items
}

func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	_ = cmd.Start()
}
