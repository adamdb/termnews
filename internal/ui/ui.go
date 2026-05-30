// Package ui contains the terminal UI components.
package ui

import (
	"fmt"
	"strings"

	"github.com/adamdb/termnews/internal/app"
	"github.com/adamdb/termnews/internal/config"
	"github.com/charmbracelet/lipgloss"
)

// Renderer holds the theme-aware styles for rendering.
type Renderer struct {
	theme config.Theme

	// Computed styles
	titleStyle       lipgloss.Style
	selectedStyle    lipgloss.Style
	normalStyle      lipgloss.Style
	dimStyle         lipgloss.Style
	metaStyle        lipgloss.Style
	linkStyle        lipgloss.Style
	errorStyle       lipgloss.Style
	statusStyle      lipgloss.Style
	tabActiveStyle   lipgloss.Style
	tabInactiveStyle lipgloss.Style
	borderStyle      lipgloss.Style
	helpStyle        lipgloss.Style
	categoryStyle    lipgloss.Style
	unreadStyle      lipgloss.Style
	bookmarkStyle    lipgloss.Style
	headerStyle      lipgloss.Style
	statusBarStyle   lipgloss.Style

	// Border configuration
	border lipgloss.Border
}

// NewRenderer creates a new renderer with the given theme.
func NewRenderer(theme config.Theme) *Renderer {
	r := &Renderer{theme: theme}
	r.initStyles()
	return r
}

// initStyles initializes all lipgloss styles from theme configuration.
func (r *Renderer) initStyles() {
	t := r.theme

	// Parse colors
	primary := lipgloss.Color(t.Colors.Primary)
	secondary := lipgloss.Color(t.Colors.Secondary)
	text := lipgloss.Color(t.Colors.Text)
	textMuted := lipgloss.Color(t.Colors.TextMuted)
	errorColor := lipgloss.Color(t.Colors.Error)
	successColor := lipgloss.Color(t.Colors.Success)
	linkColor := lipgloss.Color(t.Colors.Link)
	unreadColor := lipgloss.Color(t.Colors.Unread)
	bookmarkColor := lipgloss.Color(t.Colors.Bookmark)
	categoryColor := lipgloss.Color(t.Colors.Category)
	borderColor := lipgloss.Color(t.Colors.Border)
	headerColor := lipgloss.Color(t.Colors.Header)
	statusBg := lipgloss.Color(t.Colors.StatusBg)
	statusFg := lipgloss.Color(t.Colors.StatusFg)

	// Initialize styles
	r.titleStyle = lipgloss.NewStyle().
		Foreground(primary).
		Bold(true)

	r.selectedStyle = lipgloss.NewStyle().
		Foreground(secondary).
		Bold(true)

	r.normalStyle = lipgloss.NewStyle().
		Foreground(text)

	r.dimStyle = lipgloss.NewStyle().
		Foreground(textMuted)

	r.metaStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")) // Gray for metadata

	r.linkStyle = lipgloss.NewStyle().
		Foreground(linkColor).
		Underline(true)

	r.errorStyle = lipgloss.NewStyle().
		Foreground(errorColor)

	r.statusStyle = lipgloss.NewStyle().
		Foreground(successColor)

	r.tabActiveStyle = lipgloss.NewStyle().
		Foreground(secondary).
		Bold(true)

	r.tabInactiveStyle = lipgloss.NewStyle().
		Foreground(text)

	r.helpStyle = lipgloss.NewStyle().
		Foreground(textMuted)

	r.categoryStyle = lipgloss.NewStyle().
		Foreground(categoryColor).
		Italic(true)

	r.unreadStyle = lipgloss.NewStyle().
		Foreground(unreadColor).
		Bold(true)

	r.bookmarkStyle = lipgloss.NewStyle().
		Foreground(bookmarkColor).
		Bold(true)

	r.headerStyle = lipgloss.NewStyle().
		Foreground(headerColor).
		Bold(true)

	// Status bar style
	r.statusBarStyle = lipgloss.NewStyle().
		Foreground(statusFg)
	if t.Colors.StatusBg != "" {
		r.statusBarStyle = r.statusBarStyle.Background(statusBg)
	}

	// Configure border based on theme
	r.border = r.getBorder()

	r.borderStyle = lipgloss.NewStyle().
		Border(r.border).
		BorderForeground(borderColor)
}

// getBorder returns the appropriate lipgloss border based on theme config.
func (r *Renderer) getBorder() lipgloss.Border {
	switch r.theme.Borders.Style {
	case "sharp":
		return lipgloss.NormalBorder()
	case "double":
		return lipgloss.DoubleBorder()
	case "heavy":
		return lipgloss.ThickBorder()
	case "ascii":
		return lipgloss.Border{
			Top:         "-",
			Bottom:      "-",
			Left:        "|",
			Right:       "|",
			TopLeft:     "+",
			TopRight:    "+",
			BottomLeft:  "+",
			BottomRight: "+",
		}
	case "none":
		return lipgloss.HiddenBorder()
	case "custom":
		b := r.theme.Borders
		return lipgloss.Border{
			Top:         b.Horizontal,
			Bottom:      b.Horizontal,
			Left:        b.Vertical,
			Right:       b.Vertical,
			TopLeft:     b.TopLeft,
			TopRight:    b.TopRight,
			BottomLeft:  b.BottomLeft,
			BottomRight: b.BottomRight,
		}
	default: // "rounded"
		return lipgloss.RoundedBorder()
	}
}

// Global renderer instance - initialized on first use
var globalRenderer *Renderer

// Render renders the entire UI.
func Render(a *app.App, width, height int) string {
	// Initialize or update renderer if theme changed
	if globalRenderer == nil || globalRenderer.theme.Name != a.Config.ResolvedTheme.Name {
		globalRenderer = NewRenderer(a.Config.ResolvedTheme)
	}

	return globalRenderer.render(a, width, height)
}

// render handles the actual rendering with the renderer's styles.
func (r *Renderer) render(a *app.App, width, height int) string {
	switch a.ViewMode {
	case app.ViewModeAddSource:
		return r.renderAddSource(a, width, height)
	case app.ViewModeSummary:
		return r.renderSummary(a, width, height)
	case app.ViewModeSearch:
		return r.renderSearch(a, width, height)
	case app.ViewModeHelp:
		return r.renderHelp(a, width, height)
	default:
		return r.renderMain(a, width, height)
	}
}

func (r *Renderer) renderMain(a *app.App, width, height int) string {
	var b strings.Builder

	// Calculate banner height
	bannerHeight := 0
	if r.theme.Header.Show && r.theme.Header.Banner != "" {
		bannerHeight = strings.Count(r.theme.Header.Banner, "\n") + 2
	}

	// Render banner if configured
	if r.theme.Header.Show && r.theme.Header.Banner != "" {
		b.WriteString(r.renderBanner(width))
		b.WriteString("\n")
	}

	// Tab bar (3 lines)
	b.WriteString(r.renderTabs(a, width))
	b.WriteString("\n")

	// Content area (height - 4 for tabs and status - banner)
	contentHeight := height - 4 - bannerHeight
	if contentHeight < 1 {
		contentHeight = 1
	}
	b.WriteString(r.renderContent(a, width, contentHeight))
	b.WriteString("\n")

	// Status bar (1 line)
	b.WriteString(r.renderStatus(a, width))

	return b.String()
}

func (r *Renderer) renderBanner(width int) string {
	banner := r.theme.Header.Banner
	lines := strings.Split(banner, "\n")

	var b strings.Builder
	for _, line := range lines {
		// Center or align based on theme
		switch r.theme.Header.Align {
		case "center":
			padding := (width - len(line)) / 2
			if padding > 0 {
				b.WriteString(strings.Repeat(" ", padding))
			}
		case "right":
			padding := width - len(line)
			if padding > 0 {
				b.WriteString(strings.Repeat(" ", padding))
			}
		}
		b.WriteString(r.headerStyle.Render(line))
		b.WriteString("\n")
	}
	return b.String()
}

func (r *Renderer) renderTabs(a *app.App, width int) string {
	t := r.theme
	title := t.Header.Title
	if title == "" {
		title = "termnews"
	}

	if len(a.Tabs) == 0 {
		return r.borderStyle.Width(width - 2).Render(" " + title + " - No sources ")
	}

	var tabs []string
	for i, tab := range a.Tabs {
		name := tab.Source.Name

		// Truncate if necessary
		if t.TabBar.MaxTabWidth > 0 && len(name) > t.TabBar.MaxTabWidth {
			name = name[:t.TabBar.MaxTabWidth-len(t.TabBar.TruncateSuffix)] + t.TabBar.TruncateSuffix
		}

		if tab.Loading {
			name += " " + t.Markers.Loading
		}
		if tab.Error != nil {
			name += " " + t.Markers.Error
		}

		// Show unread count if tracking enabled
		unreadCount := 0
		for _, article := range tab.Articles {
			if !article.Read {
				unreadCount++
			}
		}
		if unreadCount > 0 && a.Config.ShowReadStatus && t.TabBar.ShowUnreadCount {
			name += fmt.Sprintf(t.TabBar.UnreadFormat, unreadCount)
		}

		if i == a.CurrentTab {
			tabText := t.Markers.TabActiveLeft + name + t.Markers.TabActiveRight
			tabs = append(tabs, r.tabActiveStyle.Render(tabText))
		} else {
			tabs = append(tabs, r.tabInactiveStyle.Render(" "+name+" "))
		}
	}

	tabBar := strings.Join(tabs, r.dimStyle.Render(t.Markers.TabSeparator))

	return r.borderStyle.Width(width - 2).Render(
		r.titleStyle.Render(" "+title+" ") + "\n" + tabBar,
	)
}

func (r *Renderer) renderContent(a *app.App, width, height int) string {
	if len(a.Tabs) == 0 {
		msg := "No news sources configured.\nPress 'a' to add a new source."
		return r.borderStyle.Width(width - 2).Height(height).Render(msg)
	}

	tab := a.Tabs[a.CurrentTab]
	t := r.theme

	if tab.Loading {
		return r.borderStyle.Width(width-2).Height(height).Render(
			fmt.Sprintf(" %s \n\n%s Fetching articles…", tab.Source.Name, t.Markers.Loading),
		)
	}

	if tab.Error != nil {
		return r.borderStyle.Width(width-2).Height(height).Render(
			fmt.Sprintf(" %s \n\n%s", tab.Source.Name, r.errorStyle.Render(t.Markers.Error+" "+tab.Error.Error())),
		)
	}

	if len(tab.Articles) == 0 {
		return r.borderStyle.Width(width-2).Height(height).Render(
			fmt.Sprintf(" %s \n\nNo articles found.", tab.Source.Name),
		)
	}

	switch a.ViewMode {
	case app.ViewModeDetail:
		return r.renderDetail(a, tab, width, height)
	default:
		return r.renderList(a, tab, width, height)
	}
}

func (r *Renderer) renderList(a *app.App, tab *app.TabState, width, height int) string {
	var b strings.Builder
	contentWidth := width - 4
	t := r.theme

	title := fmt.Sprintf(" %s (%d articles) ", tab.Source.Name, len(tab.Articles))
	if tab.Source.Category != "" {
		title += r.categoryStyle.Render("[" + tab.Source.Category + "]")
	}

	// Calculate visible articles based on display mode
	visibleHeight := height - 2 // Account for border
	
	// Determine lines per article: base + meta (if not compact) + separator (if shown)
	// Compact mode: title only (1 line)
	// Normal mode: title + meta (2 lines)
	// With separator: add 1 line
	linesPerArticle := 1 // Always have at least the title line
	if !t.ArticleList.Compact && t.ArticleList.ShowMeta {
		linesPerArticle++ // Meta line
	}
	if t.ArticleList.ShowSeparators {
		linesPerArticle++ // Separator line
	}
	// Minimum 1 line per article is guaranteed by starting at 1

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

		// Title line - reserve space for prefix (selection marker, number, indicators)
		// Truncate safely by limiting to available width minus prefix and ellipsis
		titleText := article.Title
		maxTitleLen := contentWidth - 10 // Reserve space for prefix elements and ellipsis
		if maxTitleLen < 10 {
			maxTitleLen = 10 // Minimum readable title length
		}
		if len(titleText) > maxTitleLen {
			titleText = titleText[:maxTitleLen-3] + "..."
		}

		// Build prefix with optional article number
		prefix := ""
		if t.ArticleList.ShowNumbers {
			prefix = fmt.Sprintf(t.ArticleList.NumberFormat, i+1) + " "
		}

		// Add selection marker
		if isSelected {
			prefix = t.Markers.Selected + " " + prefix
		} else {
			prefix = "  " + prefix
		}

		// Add unread/bookmark indicators
		if article.Bookmarked {
			prefix += r.bookmarkStyle.Render(t.Markers.Bookmark) + " "
		} else if !article.Read && a.Config.ShowReadStatus {
			prefix += r.unreadStyle.Render(t.Markers.Unread) + " "
		}

		if isSelected {
			b.WriteString(r.selectedStyle.Render(prefix + titleText))
		} else {
			b.WriteString(r.normalStyle.Render(prefix + titleText))
		}
		b.WriteString("\n")

		// Meta line (if not compact)
		if !t.ArticleList.Compact && t.ArticleList.ShowMeta {
			var metaParts []string
			if formattedDate := article.FormattedDate(); formattedDate != "" {
				metaParts = append(metaParts, formattedDate)
			}
			if article.Author != "" {
				metaParts = append(metaParts, "by "+article.Author)
			}
			if len(metaParts) > 0 {
				b.WriteString(r.metaStyle.Render("   " + strings.Join(metaParts, " · ")))
			}
			b.WriteString("\n")
		}

		// Separator
		if t.ArticleList.ShowSeparators {
			sep := t.ArticleList.SeparatorChar
			if sep == "" {
				sep = t.Borders.Separator
			}
			b.WriteString(r.dimStyle.Render(strings.Repeat(sep, contentWidth)))
			b.WriteString("\n")
		}
	}

	return r.borderStyle.Width(width-2).Height(height).Render(
		r.titleStyle.Render(title) + "\n" + b.String(),
	)
}

func (r *Renderer) renderDetail(a *app.App, tab *app.TabState, width, height int) string {
	if tab.Selected < 0 || tab.Selected >= len(tab.Articles) {
		return ""
	}

	article := tab.Articles[tab.Selected]
	contentWidth := width - 4
	var b strings.Builder
	t := r.theme

	// Title section
	b.WriteString(r.selectedStyle.Render(article.Title))
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
		b.WriteString(r.metaStyle.Render(strings.Join(metaParts, " · ")))
		b.WriteString("\n")
	}

	// Link
	b.WriteString(r.linkStyle.Render(article.Link))
	b.WriteString("\n")
	b.WriteString(r.dimStyle.Render(strings.Repeat(t.Borders.Separator, contentWidth)))
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
		b.WriteString(r.normalStyle.Render(line))
		b.WriteString("\n")
	}

	return r.borderStyle.Width(width-2).Height(height).Render(
		r.titleStyle.Render(" "+tab.Source.Name+" ") + "\n" + b.String(),
	)
}

func (r *Renderer) renderStatus(a *app.App, width int) string {
	t := r.theme
	keySep := t.StatusBar.KeySeparator
	if keySep == "" {
		keySep = "  "
	}

	var helpParts []string
	kb := a.Config.Keybindings

	switch a.ViewMode {
	case app.ViewModeDetail:
		helpParts = []string{
			kb.ScrollUp.Key + "/" + kb.ScrollUp.Alt + ": up",
			kb.ScrollDown.Key + "/" + kb.ScrollDown.Alt + ": down",
			kb.OpenBrowser.Key + ": browser",
			kb.Bookmark.Key + ": bookmark",
			kb.Back.Key + ": back",
			kb.Quit.Key + ": quit",
		}
	case app.ViewModeAddSource:
		helpParts = []string{"Tab: next field", "←/→: toggle type", "Enter: confirm", "Esc: cancel"}
	case app.ViewModeSummary:
		helpParts = []string{
			kb.ScrollUp.Key + "/" + kb.ScrollUp.Alt + ": up",
			kb.ScrollDown.Key + "/" + kb.ScrollDown.Alt + ": down",
			kb.Back.Key + ": back",
			kb.Quit.Key + ": quit",
		}
	case app.ViewModeSearch:
		helpParts = []string{"Type to search", "Enter: select", "Esc: cancel"}
	case app.ViewModeHelp:
		helpParts = []string{
			kb.ScrollUp.Key + "/" + kb.ScrollUp.Alt + ": up",
			kb.ScrollDown.Key + "/" + kb.ScrollDown.Alt + ": down",
			kb.Help.Key + "/Esc: close",
		}
	default:
		helpParts = []string{
			kb.NextTab.Key + ": next",
			kb.ScrollDown.Key + "/" + kb.ScrollDown.Alt + ": down",
			kb.Select.Key + ": read",
			kb.Refresh.Key + ": refresh",
			kb.Summary.Key + ": summary",
			kb.Search.Key + ": search",
			kb.AddSource.Key + ": add",
			kb.Bookmark.Key + ": bookmark",
			kb.Help.Key + ": help",
			kb.Quit.Key + ": quit",
		}
	}

	helpText := strings.Join(helpParts, keySep)

	// Show status message if present
	a.ClearStatusMessage()
	if a.StatusMessage != "" {
		return r.statusStyle.Render(a.StatusMessage)
	}

	if len(helpText) > width {
		helpText = helpText[:width-3] + "..."
	}

	return r.statusBarStyle.Render(helpText)
}

func (r *Renderer) renderAddSource(a *app.App, width, height int) string {
	if a.AddSourceState == nil {
		return ""
	}

	state := a.AddSourceState
	t := r.theme

	var b strings.Builder
	b.WriteString(r.titleStyle.Render(" Add News Source "))
	b.WriteString("\n\n")

	sel := t.Markers.Selected

	// Name field
	nameLabel := "  Name: "
	if state.ActiveField == app.InputFieldName {
		nameLabel = sel + " Name: "
		b.WriteString(r.selectedStyle.Render(nameLabel + state.Name + "▌"))
	} else {
		b.WriteString(r.normalStyle.Render(nameLabel + state.Name))
	}
	b.WriteString("\n\n")

	// URL field
	urlLabel := "  URL:  "
	if state.ActiveField == app.InputFieldURL {
		urlLabel = sel + " URL:  "
		b.WriteString(r.selectedStyle.Render(urlLabel + state.URL + "▌"))
	} else {
		b.WriteString(r.normalStyle.Render(urlLabel + state.URL))
	}
	b.WriteString("\n\n")

	// Feed type field
	typeLabel := "  Type: "
	if state.ActiveField == app.InputFieldFeedType {
		typeLabel = sel + " Type: "
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
		b.WriteString(r.selectedStyle.Render(typeLabel + typeValue))
	} else {
		b.WriteString(r.normalStyle.Render(typeLabel + typeValue))
	}
	b.WriteString("\n\n")

	// Category field
	catLabel := "  Category: "
	if state.ActiveField == app.InputFieldCategory {
		catLabel = sel + " Category: "
		b.WriteString(r.selectedStyle.Render(catLabel + state.Category + "▌"))
	} else {
		b.WriteString(r.normalStyle.Render(catLabel + state.Category))
	}
	b.WriteString("\n\n")

	// Error message
	if state.Error != "" {
		b.WriteString(r.errorStyle.Render("Error: " + state.Error))
		b.WriteString("\n")
	}

	// Help
	b.WriteString("\n")
	b.WriteString(r.helpStyle.Render(" Tab: next field  ←/→: toggle type  Enter: confirm  Esc: cancel "))

	popupWidth := 60
	if width < 64 {
		popupWidth = width - 4
	}
	popupHeight := 16

	return r.borderStyle.Width(popupWidth).Height(popupHeight).Render(b.String())
}

func (r *Renderer) renderSummary(a *app.App, width, height int) string {
	var b strings.Builder
	contentWidth := width - 4
	t := r.theme

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

			prefix := "  "
			if isSelected {
				prefix = t.Markers.Selected + " "
			}

			if isSelected {
				b.WriteString(r.selectedStyle.Render(prefix + titleText))
			} else {
				b.WriteString(r.normalStyle.Render(prefix + titleText))
			}
			b.WriteString("\n")

			// Source
			b.WriteString(r.categoryStyle.Render("   Source: " + sa.SourceName))
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
				b.WriteString(r.metaStyle.Render("   " + strings.Join(metaParts, " · ")))
			}
			b.WriteString("\n")

			// Separator
			b.WriteString(r.dimStyle.Render(strings.Repeat(t.Borders.Separator, contentWidth)))
			b.WriteString("\n")
		}
	}

	content := r.borderStyle.Width(width-2).Height(height-2).Render(
		r.titleStyle.Render(title) + "\n" + b.String(),
	)

	return content + "\n" + r.helpStyle.Render(" ↑/k: scroll up  ↓/j: scroll down  Esc: back to list  q: quit ")
}

func (r *Renderer) renderSearch(a *app.App, width, height int) string {
	var b strings.Builder
	t := r.theme

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
			b.WriteString(r.normalStyle.Render(" " + titleText))
			b.WriteString("\n")
			b.WriteString(r.categoryStyle.Render("   " + sr.SourceName))
			b.WriteString("\n")
			b.WriteString(r.dimStyle.Render(strings.Repeat(t.Borders.Separator, contentWidth)))
			b.WriteString("\n")
		}
	}

	content := r.borderStyle.Width(width-2).Height(height-2).Render(
		r.titleStyle.Render(title) + "\n" + b.String(),
	)

	return content + "\n" + r.helpStyle.Render(" Type to search  Enter: select  Esc: cancel ")
}

func (r *Renderer) renderHelp(a *app.App, width, height int) string {
	t := r.theme
	title := t.Header.Title
	if title == "" {
		title = "TERMNEWS"
	}

	// Format help content with title and underline decoration
	// First %s: application title (e.g., "TERMNEWS")
	// Second %s: underline decoration (═ repeated to match title length + 5)
	helpContent := fmt.Sprintf(`
%s HELP
%s

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
  o              Open article in browser

VIEWS
  ?              Show/hide this help
  q              Quit application

THEMING
  termnews supports extensive customization through themes.
  Available built-in themes: default, bitchx, hacker, minimal, retro

  Set theme in config.toml:
    theme = "bitchx"

  Or customize individual elements:
    [custom_theme.colors]
    primary = "14"
    secondary = "10"
    
    [custom_theme.markers]
    unread = "[+]"
    bookmark = "[*]"

FEATURES
  • RSS and Atom feed support with auto-detection
  • Persistent configuration stored in config file
  • Read/unread tracking across sessions
  • Bookmarking for important articles
  • Category organization for feeds
  • Auto-refresh at configurable intervals
  • Search across all loaded articles
  • Highly configurable themes and visual styles

CONFIGURATION
  Config file location:
    Linux/macOS: ~/.config/termnews/config.toml
    Windows:     %%APPDATA%%\termnews\config.toml

Press Esc or ? to close this help.
`, title, strings.Repeat("═", len(title)+5))

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

	content := r.borderStyle.Width(width-2).Height(height-2).Render(b.String())
	return content + "\n" + r.helpStyle.Render(" ↑/k: scroll up  ↓/j: scroll down  Esc/?: close help ")
}
