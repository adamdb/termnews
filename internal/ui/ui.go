// Package ui contains the terminal UI components.
package ui

import (
	"fmt"
	"strings"

	"github.com/adamdb/termnews/internal/app"
	"github.com/charmbracelet/lipgloss"
)

// Styles for the UI
var (
	// Colors
	colorCyan      = lipgloss.Color("6")
	colorYellow    = lipgloss.Color("11")
	colorWhite     = lipgloss.Color("15")
	colorGray      = lipgloss.Color("8")
	colorDarkGray  = lipgloss.Color("240")
	colorBlue      = lipgloss.Color("4")
	colorGreen     = lipgloss.Color("2")
	colorRed       = lipgloss.Color("1")

	// Styles
	titleStyle = lipgloss.NewStyle().
			Foreground(colorCyan).
			Bold(true)

	selectedStyle = lipgloss.NewStyle().
			Foreground(colorYellow).
			Bold(true)

	normalStyle = lipgloss.NewStyle().
			Foreground(colorWhite)

	dimStyle = lipgloss.NewStyle().
			Foreground(colorDarkGray)

	metaStyle = lipgloss.NewStyle().
			Foreground(colorGray)

	linkStyle = lipgloss.NewStyle().
			Foreground(colorBlue).
			Underline(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(colorRed)

	statusStyle = lipgloss.NewStyle().
			Foreground(colorGreen)

	tabActiveStyle = lipgloss.NewStyle().
			Foreground(colorYellow).
			Bold(true)

	tabInactiveStyle = lipgloss.NewStyle().
			Foreground(colorWhite)

	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorCyan)

	helpStyle = lipgloss.NewStyle().
			Foreground(colorDarkGray)

	categoryStyle = lipgloss.NewStyle().
			Foreground(colorCyan).
			Italic(true)

	unreadStyle = lipgloss.NewStyle().
			Foreground(colorGreen).
			Bold(true)
)

// Help text constants
const (
	HelpList    = " Tab/→: next  Shift+Tab/←: prev  ↑/k: up  ↓/j: down  Enter: read  r: refresh  t: auto-refresh  s: summary  /: search  a: add  d: delete  b: bookmark  m: mark read  o: open  ?: help  q: quit "
	HelpDetail  = " ↑/k: scroll up  ↓/j: scroll down  o: open in browser  b: bookmark  Esc: back  q: quit "
	HelpSummary = " ↑/k: scroll up  ↓/j: scroll down  Esc: back to list  q: quit "
	HelpAdd     = " Tab: next field  ←/→: toggle type  Enter: confirm  Esc: cancel "
	HelpSearch  = " Type to search  Enter: select  Esc: cancel "
	HelpHelp    = " ↑/k: scroll up  ↓/j: scroll down  Esc/?: close help "
)

// Render renders the entire UI.
func Render(a *app.App, width, height int) string {
	switch a.ViewMode {
	case app.ViewModeAddSource:
		return renderAddSource(a, width, height)
	case app.ViewModeSummary:
		return renderSummary(a, width, height)
	case app.ViewModeSearch:
		return renderSearch(a, width, height)
	case app.ViewModeHelp:
		return renderHelp(a, width, height)
	default:
		return renderMain(a, width, height)
	}
}

func renderMain(a *app.App, width, height int) string {
	var b strings.Builder

	// Tab bar (3 lines)
	b.WriteString(renderTabs(a, width))
	b.WriteString("\n")

	// Content area (height - 4 for tabs and status)
	contentHeight := height - 4
	if contentHeight < 1 {
		contentHeight = 1
	}
	b.WriteString(renderContent(a, width, contentHeight))
	b.WriteString("\n")

	// Status bar (1 line)
	b.WriteString(renderStatus(a, width))

	return b.String()
}

func renderTabs(a *app.App, width int) string {
	if len(a.Tabs) == 0 {
		return borderStyle.Width(width - 2).Render(" termnews - No sources ")
	}

	var tabs []string
	for i, tab := range a.Tabs {
		name := tab.Source.Name
		if tab.Loading {
			name += " ⟳"
		}
		if tab.Error != nil {
			name += " ✗"
		}

		// Show unread count if tracking enabled
		unreadCount := 0
		for _, article := range tab.Articles {
			if !article.Read {
				unreadCount++
			}
		}
		if unreadCount > 0 && a.Config.ShowReadStatus {
			name += fmt.Sprintf(" (%d)", unreadCount)
		}

		if i == a.CurrentTab {
			tabs = append(tabs, tabActiveStyle.Render(" "+name+" "))
		} else {
			tabs = append(tabs, tabInactiveStyle.Render(" "+name+" "))
		}
	}

	tabBar := strings.Join(tabs, dimStyle.Render("|"))

	return borderStyle.Width(width - 2).Render(
		titleStyle.Render(" termnews ") + "\n" + tabBar,
	)
}

func renderContent(a *app.App, width, height int) string {
	if len(a.Tabs) == 0 {
		msg := "No news sources configured.\nPress 'a' to add a new source."
		return borderStyle.Width(width - 2).Height(height).Render(msg)
	}

	tab := a.Tabs[a.CurrentTab]

	if tab.Loading {
		return borderStyle.Width(width-2).Height(height).Render(
			fmt.Sprintf(" %s \n\n⟳ Fetching articles…", tab.Source.Name),
		)
	}

	if tab.Error != nil {
		return borderStyle.Width(width-2).Height(height).Render(
			fmt.Sprintf(" %s \n\n%s", tab.Source.Name, errorStyle.Render("✗ "+tab.Error.Error())),
		)
	}

	if len(tab.Articles) == 0 {
		return borderStyle.Width(width-2).Height(height).Render(
			fmt.Sprintf(" %s \n\nNo articles found.", tab.Source.Name),
		)
	}

	switch a.ViewMode {
	case app.ViewModeDetail:
		return renderDetail(a, tab, width, height)
	default:
		return renderList(a, tab, width, height)
	}
}

func renderList(a *app.App, tab *app.TabState, width, height int) string {
	var b strings.Builder
	contentWidth := width - 4

	title := fmt.Sprintf(" %s (%d articles) ", tab.Source.Name, len(tab.Articles))
	if tab.Source.Category != "" {
		title += categoryStyle.Render("[" + tab.Source.Category + "]")
	}

	// Calculate visible articles
	visibleHeight := height - 2 // Account for border
	linesPerArticle := 3        // title + meta + separator
	visibleArticles := visibleHeight / linesPerArticle
	if visibleArticles < 1 {
		visibleArticles = 1
	}

	startIdx := tab.ScrollOffset
	endIdx := startIdx + visibleArticles
	if endIdx > len(tab.Articles) {
		endIdx = len(tab.Articles)
	}

	for i := startIdx; i < endIdx; i++ {
		article := tab.Articles[i]
		isSelected := i == tab.Selected

		// Title line
		titleText := article.Title
		if len(titleText) > contentWidth-2 {
			titleText = titleText[:contentWidth-5] + "..."
		}

		// Add unread/bookmark indicators
		prefix := " "
		if !article.Read && a.Config.ShowReadStatus {
			prefix = unreadStyle.Render("●") + " "
		}
		if article.Bookmarked {
			prefix = "★ "
		}

		if isSelected {
			b.WriteString(selectedStyle.Render(prefix + titleText))
		} else {
			b.WriteString(normalStyle.Render(prefix + titleText))
		}
		b.WriteString("\n")

		// Meta line
		var metaParts []string
		if formattedDate := article.FormattedDate(); formattedDate != "" {
			metaParts = append(metaParts, formattedDate)
		}
		if article.Author != "" {
			metaParts = append(metaParts, "by "+article.Author)
		}
		if len(metaParts) > 0 {
			b.WriteString(metaStyle.Render("   " + strings.Join(metaParts, " · ")))
		}
		b.WriteString("\n")

		// Separator
		b.WriteString(dimStyle.Render(strings.Repeat("─", contentWidth)))
		b.WriteString("\n")
	}

	return borderStyle.Width(width-2).Height(height).Render(
		titleStyle.Render(title) + "\n" + b.String(),
	)
}

func renderDetail(a *app.App, tab *app.TabState, width, height int) string {
	if tab.Selected < 0 || tab.Selected >= len(tab.Articles) {
		return ""
	}

	article := tab.Articles[tab.Selected]
	contentWidth := width - 4
	var b strings.Builder

	// Title section
	b.WriteString(selectedStyle.Render(article.Title))
	b.WriteString("\n")

	// Meta line
	var metaParts []string
	if formattedDate := article.FormattedDate(); formattedDate != "" {
		metaParts = append(metaParts, formattedDate)
	}
	if article.Author != "" {
		metaParts = append(metaParts, "by "+article.Author)
	}
	if len(metaParts) > 0 {
		b.WriteString(metaStyle.Render(strings.Join(metaParts, " · ")))
		b.WriteString("\n")
	}

	// Link
	b.WriteString(linkStyle.Render(article.Link))
	b.WriteString("\n")
	b.WriteString(dimStyle.Render(strings.Repeat("─", contentWidth)))
	b.WriteString("\n\n")

	// Article summary
	summary := article.Summary()
	lines := strings.Split(summary, "\n")

	// Apply scroll offset
	startLine := tab.ScrollOffset
	if startLine >= len(lines) {
		startLine = 0
	}

	// Calculate visible lines
	visibleHeight := height - 8 // Account for header and border
	endLine := startLine + visibleHeight
	if endLine > len(lines) {
		endLine = len(lines)
	}

	for i := startLine; i < endLine; i++ {
		line := lines[i]
		if len(line) > contentWidth {
			line = line[:contentWidth]
		}
		b.WriteString(normalStyle.Render(line))
		b.WriteString("\n")
	}

	return borderStyle.Width(width-2).Height(height).Render(
		titleStyle.Render(" "+tab.Source.Name+" ") + "\n" + b.String(),
	)
}

func renderStatus(a *app.App, width int) string {
	var helpText string
	switch a.ViewMode {
	case app.ViewModeDetail:
		helpText = HelpDetail
	case app.ViewModeAddSource:
		helpText = HelpAdd
	case app.ViewModeSummary:
		helpText = HelpSummary
	case app.ViewModeSearch:
		helpText = HelpSearch
	case app.ViewModeHelp:
		helpText = HelpHelp
	default:
		helpText = HelpList
	}

	// Show status message if present
	a.ClearStatusMessage()
	if a.StatusMessage != "" {
		return statusStyle.Render(a.StatusMessage)
	}

	if len(helpText) > width {
		helpText = helpText[:width-3] + "..."
	}

	return helpStyle.Render(helpText)
}

func renderAddSource(a *app.App, width, height int) string {
	if a.AddSourceState == nil {
		return ""
	}

	state := a.AddSourceState

	var b strings.Builder
	b.WriteString(titleStyle.Render(" Add News Source "))
	b.WriteString("\n\n")

	// Name field
	nameLabel := "  Name: "
	if state.ActiveField == app.InputFieldName {
		nameLabel = "▶ Name: "
		b.WriteString(selectedStyle.Render(nameLabel + state.Name + "▌"))
	} else {
		b.WriteString(normalStyle.Render(nameLabel + state.Name))
	}
	b.WriteString("\n\n")

	// URL field
	urlLabel := "  URL:  "
	if state.ActiveField == app.InputFieldURL {
		urlLabel = "▶ URL:  "
		b.WriteString(selectedStyle.Render(urlLabel + state.URL + "▌"))
	} else {
		b.WriteString(normalStyle.Render(urlLabel + state.URL))
	}
	b.WriteString("\n\n")

	// Feed type field
	typeLabel := "  Type: "
	if state.ActiveField == app.InputFieldFeedType {
		typeLabel = "▶ Type: "
	}
	var typeValue string
	switch state.FeedTypeIndex {
	case 0:
		typeValue = "[RSS] / Atom / Auto"
	case 1:
		typeValue = "RSS / [Atom] / Auto"
	default:
		typeValue = "RSS / Atom / [Auto]"
	}
	if state.ActiveField == app.InputFieldFeedType {
		b.WriteString(selectedStyle.Render(typeLabel + typeValue))
	} else {
		b.WriteString(normalStyle.Render(typeLabel + typeValue))
	}
	b.WriteString("\n\n")

	// Category field
	catLabel := "  Category: "
	if state.ActiveField == app.InputFieldCategory {
		catLabel = "▶ Category: "
		b.WriteString(selectedStyle.Render(catLabel + state.Category + "▌"))
	} else {
		b.WriteString(normalStyle.Render(catLabel + state.Category))
	}
	b.WriteString("\n\n")

	// Error message
	if state.Error != "" {
		b.WriteString(errorStyle.Render("Error: " + state.Error))
		b.WriteString("\n")
	}

	// Help
	b.WriteString("\n")
	b.WriteString(helpStyle.Render(HelpAdd))

	popupWidth := 60
	if width < 64 {
		popupWidth = width - 4
	}
	popupHeight := 16

	return borderStyle.Width(popupWidth).Height(popupHeight).Render(b.String())
}

func renderSummary(a *app.App, width, height int) string {
	var b strings.Builder
	contentWidth := width - 4

	title := fmt.Sprintf(" Summary - Top Articles from All Feeds (%d total) ", len(a.SummaryArticles))

	if len(a.SummaryArticles) == 0 {
		b.WriteString("No articles available for summary.\n")
		b.WriteString("Refresh feeds to see content.")
	} else {
		// Calculate visible articles
		visibleHeight := height - 4
		linesPerArticle := 4 // title + source + meta + separator
		visibleArticles := visibleHeight / linesPerArticle
		if visibleArticles < 1 {
			visibleArticles = 1
		}

		startIdx := a.SummaryScrollOffset
		endIdx := startIdx + visibleArticles
		if endIdx > len(a.SummaryArticles) {
			endIdx = len(a.SummaryArticles)
		}

		for i := startIdx; i < endIdx; i++ {
			sa := a.SummaryArticles[i]
			isSelected := i == a.SummarySelected

			// Title
			titleText := sa.Article.Title
			if len(titleText) > contentWidth-2 {
				titleText = titleText[:contentWidth-5] + "..."
			}

			if isSelected {
				b.WriteString(selectedStyle.Render(" " + titleText))
			} else {
				b.WriteString(normalStyle.Render(" " + titleText))
			}
			b.WriteString("\n")

			// Source
			b.WriteString(categoryStyle.Render("   Source: " + sa.SourceName))
			b.WriteString("\n")

			// Meta
			var metaParts []string
			if formattedDate := sa.Article.FormattedDate(); formattedDate != "" {
				metaParts = append(metaParts, formattedDate)
			}
			if sa.Article.Author != "" {
				metaParts = append(metaParts, "by "+sa.Article.Author)
			}
			if len(metaParts) > 0 {
				b.WriteString(metaStyle.Render("   " + strings.Join(metaParts, " · ")))
			}
			b.WriteString("\n")

			// Separator
			b.WriteString(dimStyle.Render(strings.Repeat("─", contentWidth)))
			b.WriteString("\n")
		}
	}

	content := borderStyle.Width(width-2).Height(height-2).Render(
		titleStyle.Render(title) + "\n" + b.String(),
	)

	return content + "\n" + helpStyle.Render(HelpSummary)
}

func renderSearch(a *app.App, width, height int) string {
	var b strings.Builder

	title := " Search Articles "
	b.WriteString("Search: " + a.SearchQuery + "▌")
	b.WriteString("\n\n")

	if len(a.SearchResults) == 0 {
		if a.SearchQuery == "" {
			b.WriteString("Type to search across all feeds...")
		} else {
			b.WriteString("No results found.")
		}
	} else {
		b.WriteString(fmt.Sprintf("Found %d results:\n\n", len(a.SearchResults)))
		contentWidth := width - 4

		visibleHeight := height - 8
		linesPerResult := 3
		visibleResults := visibleHeight / linesPerResult
		if visibleResults < 1 {
			visibleResults = 1
		}
		if visibleResults > len(a.SearchResults) {
			visibleResults = len(a.SearchResults)
		}

		for i := 0; i < visibleResults; i++ {
			sr := a.SearchResults[i]
			titleText := sr.Article.Title
			if len(titleText) > contentWidth-4 {
				titleText = titleText[:contentWidth-7] + "..."
			}
			b.WriteString(normalStyle.Render(" " + titleText))
			b.WriteString("\n")
			b.WriteString(categoryStyle.Render("   " + sr.SourceName))
			b.WriteString("\n")
			b.WriteString(dimStyle.Render(strings.Repeat("─", contentWidth)))
			b.WriteString("\n")
		}
	}

	content := borderStyle.Width(width-2).Height(height-2).Render(
		titleStyle.Render(title) + "\n" + b.String(),
	)

	return content + "\n" + helpStyle.Render(HelpSearch)
}

func renderHelp(a *app.App, width, height int) string {
	helpContent := `
TERMNEWS HELP
═════════════

NAVIGATION
  Tab, →         Next feed tab
  Shift+Tab, ←   Previous feed tab
  ↓, j           Move selection down / scroll down
  ↑, k           Move selection up / scroll up
  Enter          Open selected article
  Esc, Backspace Back to article list

ACTIONS
  r              Refresh current feed
  R              Refresh all feeds
  t              Toggle auto-refresh
  s              Summary view (top articles from all feeds)
  /              Search across all articles
  a              Add a new news source
  d              Delete current news source
  b              Toggle bookmark on selected article
  m              Mark all articles in current feed as read
  o              Open article in browser (copies URL)

VIEWS
  ?              Show/hide this help
  q              Quit application

FEATURES
  • RSS and Atom feed support with auto-detection
  • Persistent configuration stored in config file
  • Read/unread tracking across sessions
  • Bookmarking for important articles
  • Category organization for feeds
  • Auto-refresh at configurable intervals
  • Search across all loaded articles

CONFIGURATION
  Config file location:
    Linux/macOS: ~/.config/termnews/config.toml
    Windows:     %APPDATA%\termnews\config.toml

  Example config:
    [[sources]]
    name = "Hacker News"
    url = "https://news.ycombinator.com/rss"
    feed_type = "auto"
    category = "Tech"

Press Esc or ? to close this help.
`

	lines := strings.Split(helpContent, "\n")
	var b strings.Builder

	visibleHeight := height - 4
	startLine := a.HelpScrollOffset
	endLine := startLine + visibleHeight
	if endLine > len(lines) {
		endLine = len(lines)
	}

	for i := startLine; i < endLine; i++ {
		b.WriteString(lines[i])
		b.WriteString("\n")
	}

	content := borderStyle.Width(width-2).Height(height-2).Render(b.String())
	return content + "\n" + helpStyle.Render(HelpHelp)
}
