// Termnews is a terminal news aggregator for RSS and Atom feeds.
package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/adamdb/termnews/internal/app"
	"github.com/adamdb/termnews/internal/config"
	"github.com/adamdb/termnews/internal/feed"
	"github.com/adamdb/termnews/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

// Model represents the application state for Bubble Tea.
type Model struct {
	app    *app.App
	width  int
	height int
	ctx    context.Context
	cancel context.CancelFunc
}

// feedResultMsg is a message containing the result of a feed fetch.
type feedResultMsg struct {
	index    int
	articles []feed.Article
	err      error
}

// tickMsg is a message for periodic updates.
type tickMsg time.Time

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	// Start fetching all feeds
	cmds := make([]tea.Cmd, len(m.app.Tabs))
	for i := range m.app.Tabs {
		cmds[i] = m.fetchFeed(i)
	}
	// Add tick command for auto-refresh and status clearing
	cmds = append(cmds, tickCmd())
	return tea.Batch(cmds...)
}

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		cmd := m.handleKey(msg)
		if m.app.ShouldQuit {
			m.cancel()
			return m, tea.Quit
		}
		return m, cmd

	case feedResultMsg:
		m.app.SetArticles(msg.index, msg.articles, msg.err)
		return m, nil

	case tickMsg:
		// Check for auto-refresh
		if m.app.ShouldAutoRefresh() {
			m.app.RefreshAll()
			cmds := make([]tea.Cmd, len(m.app.Tabs))
			for i := range m.app.Tabs {
				cmds[i] = m.fetchFeed(i)
			}
			cmds = append(cmds, tickCmd())
			return m, tea.Batch(cmds...)
		}
		return m, tickCmd()
	}

	return m, nil
}

// View implements tea.Model.
func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}
	return ui.Render(m.app, m.width, m.height)
}

// handleKey handles keyboard input.
func (m *Model) handleKey(msg tea.KeyMsg) tea.Cmd {
	switch m.app.ViewMode {
	case app.ViewModeAddSource:
		return m.handleAddSourceKey(msg)
	case app.ViewModeDetail:
		return m.handleDetailKey(msg)
	case app.ViewModeSummary:
		return m.handleSummaryKey(msg)
	case app.ViewModeSearch:
		return m.handleSearchKey(msg)
	case app.ViewModeHelp:
		return m.handleHelpKey(msg)
	default:
		return m.handleListKey(msg)
	}
}

func (m *Model) handleListKey(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "q", "Q", "ctrl+c":
		m.app.ShouldQuit = true
	case "tab", "right":
		m.app.NextTab()
	case "shift+tab", "left":
		m.app.PrevTab()
	case "down", "j":
		m.app.ScrollDown()
	case "up", "k":
		m.app.ScrollUp()
	case "enter":
		m.app.SelectArticle()
	case "r":
		return m.refreshCurrent()
	case "R":
		return m.refreshAll()
	case "t", "T":
		m.app.ToggleAutoRefresh()
	case "s", "S":
		m.app.EnterSummaryMode()
	case "a", "A":
		m.app.StartAddSource()
	case "d", "D":
		m.app.RemoveCurrentSource()
	case "b", "B":
		m.app.ToggleBookmark()
	case "m", "M":
		m.app.MarkAllRead()
	case "o", "O":
		if url := m.app.OpenInBrowser(); url != "" {
			openBrowser(url)
			m.app.SetStatusMessage("Opening in browser...")
		}
	case "/":
		m.app.EnterSearchMode()
	case "?":
		m.app.EnterHelpMode()
	}
	return nil
}

func (m *Model) handleDetailKey(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "q", "Q", "ctrl+c":
		m.app.ShouldQuit = true
	case "esc", "backspace":
		m.app.BackToList()
	case "down", "j":
		m.app.ScrollDown()
	case "up", "k":
		m.app.ScrollUp()
	case "b", "B":
		m.app.ToggleBookmark()
	case "o", "O":
		if url := m.app.OpenInBrowser(); url != "" {
			openBrowser(url)
			m.app.SetStatusMessage("Opening in browser...")
		}
	}
	return nil
}

func (m *Model) handleSummaryKey(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "q", "Q", "ctrl+c":
		m.app.ShouldQuit = true
	case "esc", "backspace":
		m.app.ExitSummaryMode()
	case "down", "j":
		m.app.ScrollSummaryDown()
	case "up", "k":
		m.app.ScrollSummaryUp()
	}
	return nil
}

func (m *Model) handleSearchKey(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "ctrl+c":
		m.app.ShouldQuit = true
	case "esc":
		m.app.ExitSearchMode()
	case "backspace":
		if len(m.app.SearchQuery) > 0 {
			m.app.SearchQuery = m.app.SearchQuery[:len(m.app.SearchQuery)-1]
			m.app.UpdateSearch()
		}
	case "enter":
		// TODO: Select search result
		m.app.ExitSearchMode()
	default:
		// Add character to search query
		if len(msg.String()) == 1 {
			m.app.SearchQuery += msg.String()
			m.app.UpdateSearch()
		}
	}
	return nil
}

func (m *Model) handleHelpKey(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "q", "Q", "ctrl+c":
		m.app.ShouldQuit = true
	case "esc", "?":
		m.app.ExitHelpMode()
	case "down", "j":
		m.app.ScrollDown()
	case "up", "k":
		m.app.ScrollUp()
	}
	return nil
}

func (m *Model) handleAddSourceKey(msg tea.KeyMsg) tea.Cmd {
	state := m.app.AddSourceState
	if state == nil {
		return nil
	}

	switch msg.String() {
	case "ctrl+c":
		m.app.ShouldQuit = true
	case "esc":
		m.app.CancelAddSource()
	case "tab":
		state.Error = ""
		switch state.ActiveField {
		case app.InputFieldName:
			state.ActiveField = app.InputFieldURL
		case app.InputFieldURL:
			state.ActiveField = app.InputFieldFeedType
		case app.InputFieldFeedType:
			state.ActiveField = app.InputFieldCategory
		case app.InputFieldCategory:
			state.ActiveField = app.InputFieldName
		}
	case "shift+tab":
		state.Error = ""
		switch state.ActiveField {
		case app.InputFieldName:
			state.ActiveField = app.InputFieldCategory
		case app.InputFieldURL:
			state.ActiveField = app.InputFieldName
		case app.InputFieldFeedType:
			state.ActiveField = app.InputFieldURL
		case app.InputFieldCategory:
			state.ActiveField = app.InputFieldFeedType
		}
	case "enter":
		if m.app.ConfirmAddSource() {
			// Fetch the newly added feed
			newIdx := len(m.app.Tabs) - 1
			return m.fetchFeed(newIdx)
		}
	case "backspace":
		switch state.ActiveField {
		case app.InputFieldName:
			if len(state.Name) > 0 {
				state.Name = state.Name[:len(state.Name)-1]
			}
		case app.InputFieldURL:
			if len(state.URL) > 0 {
				state.URL = state.URL[:len(state.URL)-1]
			}
		case app.InputFieldCategory:
			if len(state.Category) > 0 {
				state.Category = state.Category[:len(state.Category)-1]
			}
		}
		state.Error = ""
	case "left", "right":
		if state.ActiveField == app.InputFieldFeedType {
			if msg.String() == "right" {
				state.FeedTypeIndex = (state.FeedTypeIndex + 1) % 3
			} else {
				state.FeedTypeIndex = (state.FeedTypeIndex + 2) % 3
			}
		}
	case " ":
		if state.ActiveField == app.InputFieldFeedType {
			state.FeedTypeIndex = (state.FeedTypeIndex + 1) % 3
		} else {
			// Add space to text fields
			switch state.ActiveField {
			case app.InputFieldName:
				state.Name += " "
			case app.InputFieldURL:
				state.URL += " "
			case app.InputFieldCategory:
				state.Category += " "
			}
		}
	default:
		// Add character to current field
		if len(msg.String()) == 1 {
			switch state.ActiveField {
			case app.InputFieldName:
				state.Name += msg.String()
			case app.InputFieldURL:
				state.URL += msg.String()
			case app.InputFieldCategory:
				state.Category += msg.String()
			}
			state.Error = ""
		}
	}
	return nil
}

func (m *Model) refreshCurrent() tea.Cmd {
	if len(m.app.Tabs) == 0 {
		return nil
	}
	idx := m.app.CurrentTab
	m.app.RefreshCurrent()
	return m.fetchFeed(idx)
}

func (m *Model) refreshAll() tea.Cmd {
	m.app.RefreshAll()
	cmds := make([]tea.Cmd, len(m.app.Tabs))
	for i := range m.app.Tabs {
		cmds[i] = m.fetchFeed(i)
	}
	return tea.Batch(cmds...)
}

func (m *Model) fetchFeed(index int) tea.Cmd {
	if index < 0 || index >= len(m.app.Tabs) {
		return nil
	}
	source := m.app.Tabs[index].Source
	ctx := m.ctx

	return func() tea.Msg {
		articles, err := feed.Fetch(ctx, &source)
		return feedResultMsg{
			index:    index,
			articles: articles,
			err:      err,
		}
	}
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*30, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// openBrowser opens the given URL in the default browser.
func openBrowser(url string) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	default:
		return
	}

	// Run in background, ignore errors
	_ = cmd.Start()
}

func main() {
	// Load configuration
	cfg := config.Load()

	// Save defaults if no config file exists
	configPath := config.ConfigPath()
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := cfg.Save(); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not save config: %v\n", err)
		}
	}

	// Create app state
	application := app.New(cfg)

	// Create context for cancellation
	ctx, cancel := context.WithCancel(context.Background())

	// Create model
	model := Model{
		app:    application,
		ctx:    ctx,
		cancel: cancel,
	}

	// Run the Bubble Tea program
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running termnews: %v\n", err)
		os.Exit(1)
	}
}
