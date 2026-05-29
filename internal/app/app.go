// Package app contains the application state and business logic.
package app

import (
	"sort"
	"time"

	"github.com/adamdb/termnews/internal/config"
	"github.com/adamdb/termnews/internal/feed"
)

// ViewMode represents the current view mode of the application.
type ViewMode int

const (
	ViewModeList ViewMode = iota
	ViewModeDetail
	ViewModeAddSource
	ViewModeSummary
	ViewModeSearch
	ViewModeHelp
)

// InputField represents which field is active in the add source form.
type InputField int

const (
	InputFieldName InputField = iota
	InputFieldURL
	InputFieldFeedType
	InputFieldCategory
)

// AddSourceState holds the state for the add source form.
type AddSourceState struct {
	ActiveField   InputField
	Name          string
	URL           string
	FeedTypeIndex int    // 0 = RSS, 1 = Atom, 2 = Auto
	Category      string
	Error         string
}

// NewAddSourceState creates a new add source state.
func NewAddSourceState() *AddSourceState {
	return &AddSourceState{
		ActiveField:   InputFieldName,
		FeedTypeIndex: 2, // Default to auto-detect
	}
}

// SelectedFeedType returns the selected feed type.
func (s *AddSourceState) SelectedFeedType() config.FeedType {
	switch s.FeedTypeIndex {
	case 0:
		return config.FeedTypeRSS
	case 1:
		return config.FeedTypeAtom
	default:
		return config.FeedTypeAuto
	}
}

// TabState holds the state for a single feed tab.
type TabState struct {
	Source       config.Source
	Articles     []feed.Article
	ScrollOffset int
	Selected     int
	Loading      bool
	Error        error
	LastFetched  time.Time
}

// NewTabState creates a new tab state for a source.
func NewTabState(source config.Source) *TabState {
	return &TabState{
		Source:   source,
		Articles: make([]feed.Article, 0),
		Selected: -1,
		Loading:  true,
	}
}

// App holds all application state.
type App struct {
	Config             *config.Config
	Tabs               []*TabState
	CurrentTab         int
	ViewMode           ViewMode
	AddSourceState     *AddSourceState
	ShouldQuit         bool
	StatusMessage      string
	StatusMessageTime  time.Time
	AutoRefreshEnabled bool
	LastRefreshTime    time.Time

	// Summary view state
	SummaryArticles     []SummaryArticle
	SummarySelected     int
	SummaryScrollOffset int

	// Search state
	SearchQuery   string
	SearchResults []SummaryArticle

	// Help state
	HelpScrollOffset int

	// Read tracking
	ReadArticles map[string]bool

	// Bookmarks
	Bookmarks []SummaryArticle
}

// SummaryArticle holds an article with its source name for summary/search views.
type SummaryArticle struct {
	SourceName string
	Article    feed.Article
}

// New creates a new App instance.
func New(cfg *config.Config) *App {
	tabs := make([]*TabState, len(cfg.Sources))
	for i, source := range cfg.Sources {
		tabs[i] = NewTabState(source)
	}

	return &App{
		Config:             cfg,
		Tabs:               tabs,
		CurrentTab:         0,
		ViewMode:           ViewModeList,
		AutoRefreshEnabled: false,
		LastRefreshTime:    time.Now(),
		ReadArticles:       make(map[string]bool),
		Bookmarks:          make([]SummaryArticle, 0),
	}
}

// NextTab switches to the next tab.
func (a *App) NextTab() {
	if len(a.Tabs) == 0 {
		return
	}
	a.CurrentTab = (a.CurrentTab + 1) % len(a.Tabs)
	a.ViewMode = ViewModeList
}

// PrevTab switches to the previous tab.
func (a *App) PrevTab() {
	if len(a.Tabs) == 0 {
		return
	}
	a.CurrentTab = (a.CurrentTab + len(a.Tabs) - 1) % len(a.Tabs)
	a.ViewMode = ViewModeList
}

// ScrollDown moves the selection down.
func (a *App) ScrollDown() {
	if len(a.Tabs) == 0 || a.CurrentTab >= len(a.Tabs) {
		return
	}
	tab := a.Tabs[a.CurrentTab]

	switch a.ViewMode {
	case ViewModeList:
		if len(tab.Articles) == 0 {
			return
		}
		maxIdx := len(tab.Articles) - 1
		if tab.Selected < maxIdx {
			tab.Selected++
		}
		// Adjust scroll offset for visibility
		if tab.Selected > tab.ScrollOffset+20 {
			tab.ScrollOffset = tab.Selected - 20
		}
	case ViewModeDetail:
		tab.ScrollOffset++
	case ViewModeHelp:
		a.HelpScrollOffset++
	}
}

// ScrollUp moves the selection up.
func (a *App) ScrollUp() {
	if len(a.Tabs) == 0 || a.CurrentTab >= len(a.Tabs) {
		return
	}
	tab := a.Tabs[a.CurrentTab]

	switch a.ViewMode {
	case ViewModeList:
		if len(tab.Articles) == 0 {
			return
		}
		if tab.Selected > 0 {
			tab.Selected--
		}
		if tab.Selected < tab.ScrollOffset {
			tab.ScrollOffset = tab.Selected
		}
	case ViewModeDetail:
		if tab.ScrollOffset > 0 {
			tab.ScrollOffset--
		}
	case ViewModeHelp:
		if a.HelpScrollOffset > 0 {
			a.HelpScrollOffset--
		}
	}
}

// SelectArticle enters detail view for the selected article.
func (a *App) SelectArticle() {
	if len(a.Tabs) == 0 || a.CurrentTab >= len(a.Tabs) {
		return
	}
	tab := a.Tabs[a.CurrentTab]
	if tab.Selected >= 0 && tab.Selected < len(tab.Articles) {
		tab.ScrollOffset = 0
		a.ViewMode = ViewModeDetail
		// Mark as read
		article := &tab.Articles[tab.Selected]
		if article.GUID != "" {
			a.ReadArticles[article.GUID] = true
			article.Read = true
		}
	}
}

// BackToList returns to the list view.
func (a *App) BackToList() {
	if a.ViewMode == ViewModeDetail {
		a.ViewMode = ViewModeList
		if len(a.Tabs) > 0 && a.CurrentTab < len(a.Tabs) {
			a.Tabs[a.CurrentTab].ScrollOffset = 0
		}
	}
}

// StartAddSource enters the add source mode.
func (a *App) StartAddSource() {
	a.AddSourceState = NewAddSourceState()
	a.ViewMode = ViewModeAddSource
}

// CancelAddSource cancels adding a source.
func (a *App) CancelAddSource() {
	a.AddSourceState = nil
	a.ViewMode = ViewModeList
}

// ConfirmAddSource confirms adding a new source.
func (a *App) ConfirmAddSource() bool {
	if a.AddSourceState == nil {
		return false
	}

	name := a.AddSourceState.Name
	url := a.AddSourceState.URL
	category := a.AddSourceState.Category

	if name == "" || url == "" {
		a.AddSourceState.Error = "Name and URL cannot be empty"
		return false
	}

	feedType := a.AddSourceState.SelectedFeedType()
	a.Config.AddSource(name, url, feedType, category)
	if err := a.Config.Save(); err != nil {
		a.AddSourceState.Error = "Failed to save config: " + err.Error()
		return false
	}

	source := config.Source{
		Name:     name,
		URL:      url,
		FeedType: feedType,
		Category: category,
	}
	tab := NewTabState(source)
	tab.Loading = true
	a.Tabs = append(a.Tabs, tab)

	a.AddSourceState = nil
	a.ViewMode = ViewModeList
	a.CurrentTab = len(a.Tabs) - 1

	return true
}

// RemoveCurrentSource removes the current source.
func (a *App) RemoveCurrentSource() {
	if len(a.Tabs) == 0 {
		return
	}

	a.Config.RemoveSource(a.CurrentTab)
	if err := a.Config.Save(); err != nil {
		a.SetStatusMessage("Failed to save config: " + err.Error())
		return
	}

	a.Tabs = append(a.Tabs[:a.CurrentTab], a.Tabs[a.CurrentTab+1:]...)
	if a.CurrentTab >= len(a.Tabs) && len(a.Tabs) > 0 {
		a.CurrentTab = len(a.Tabs) - 1
	}
	a.ViewMode = ViewModeList
}

// SetArticles sets the articles for a tab.
func (a *App) SetArticles(tabIndex int, articles []feed.Article, err error) {
	if tabIndex < 0 || tabIndex >= len(a.Tabs) {
		return
	}

	tab := a.Tabs[tabIndex]
	tab.Loading = false
	tab.LastFetched = time.Now()

	if err != nil {
		tab.Error = err
		return
	}

	// Apply max articles limit
	maxArticles := a.Config.MaxArticlesPerFeed
	if maxArticles > 0 && len(articles) > maxArticles {
		articles = articles[:maxArticles]
	}

	// Apply read status from tracking
	for i := range articles {
		if a.ReadArticles[articles[i].GUID] {
			articles[i].Read = true
		}
	}

	tab.Articles = articles
	tab.Error = nil
	if len(articles) > 0 && tab.Selected < 0 {
		tab.Selected = 0
	}
}

// RefreshCurrent marks the current tab for refresh.
func (a *App) RefreshCurrent() {
	if len(a.Tabs) == 0 || a.CurrentTab >= len(a.Tabs) {
		return
	}
	tab := a.Tabs[a.CurrentTab]
	tab.Loading = true
	tab.Articles = nil
	tab.Selected = -1
	tab.ScrollOffset = 0
	tab.Error = nil
}

// RefreshAll marks all tabs for refresh.
func (a *App) RefreshAll() {
	a.LastRefreshTime = time.Now()
	for _, tab := range a.Tabs {
		tab.Loading = true
		tab.Articles = nil
		tab.Selected = -1
		tab.ScrollOffset = 0
		tab.Error = nil
	}
}

// ToggleAutoRefresh toggles auto-refresh.
func (a *App) ToggleAutoRefresh() {
	a.AutoRefreshEnabled = !a.AutoRefreshEnabled
	a.LastRefreshTime = time.Now()
	if a.AutoRefreshEnabled {
		a.SetStatusMessage("Auto-refresh enabled")
	} else {
		a.SetStatusMessage("Auto-refresh disabled")
	}
}

// ShouldAutoRefresh returns true if auto-refresh should occur.
func (a *App) ShouldAutoRefresh() bool {
	if !a.AutoRefreshEnabled {
		return false
	}
	elapsed := time.Since(a.LastRefreshTime)
	return elapsed >= time.Duration(a.Config.RefreshIntervalMins)*time.Minute
}

// SetStatusMessage sets a temporary status message.
func (a *App) SetStatusMessage(msg string) {
	a.StatusMessage = msg
	a.StatusMessageTime = time.Now()
}

// ClearStatusMessage clears the status message if it's old enough.
func (a *App) ClearStatusMessage() {
	if a.StatusMessage != "" && time.Since(a.StatusMessageTime) > 3*time.Second {
		a.StatusMessage = ""
	}
}

// EnterSummaryMode enters the summary view.
func (a *App) EnterSummaryMode() {
	a.ViewMode = ViewModeSummary
	a.BuildSummary()
}

// ExitSummaryMode exits the summary view.
func (a *App) ExitSummaryMode() {
	a.ViewMode = ViewModeList
	a.SummarySelected = 0
	a.SummaryScrollOffset = 0
}

// BuildSummary builds the summary view with top articles from all feeds.
func (a *App) BuildSummary() {
	a.SummaryArticles = nil

	for _, tab := range a.Tabs {
		if tab.Error != nil || tab.Loading {
			continue
		}
		// Take top 5 articles from each feed
		count := 5
		if len(tab.Articles) < count {
			count = len(tab.Articles)
		}
		for i := 0; i < count; i++ {
			a.SummaryArticles = append(a.SummaryArticles, SummaryArticle{
				SourceName: tab.Source.Name,
				Article:    tab.Articles[i],
			})
		}
	}

	// Sort by publication date (newest first)
	sort.Slice(a.SummaryArticles, func(i, j int) bool {
		return a.SummaryArticles[i].Article.PubDate.After(a.SummaryArticles[j].Article.PubDate)
	})

	if len(a.SummaryArticles) > 0 {
		a.SummarySelected = 0
	} else {
		a.SummarySelected = -1
	}
	a.SummaryScrollOffset = 0
}

// ScrollSummaryDown moves the summary selection down.
func (a *App) ScrollSummaryDown() {
	if len(a.SummaryArticles) == 0 {
		return
	}
	maxIdx := len(a.SummaryArticles) - 1
	if a.SummarySelected < maxIdx {
		a.SummarySelected++
	}
	if a.SummarySelected > a.SummaryScrollOffset+20 {
		a.SummaryScrollOffset = a.SummarySelected - 20
	}
}

// ScrollSummaryUp moves the summary selection up.
func (a *App) ScrollSummaryUp() {
	if len(a.SummaryArticles) == 0 {
		return
	}
	if a.SummarySelected > 0 {
		a.SummarySelected--
	}
	if a.SummarySelected < a.SummaryScrollOffset {
		a.SummaryScrollOffset = a.SummarySelected
	}
}

// EnterSearchMode enters search mode.
func (a *App) EnterSearchMode() {
	a.ViewMode = ViewModeSearch
	a.SearchQuery = ""
	a.SearchResults = nil
}

// ExitSearchMode exits search mode.
func (a *App) ExitSearchMode() {
	a.ViewMode = ViewModeList
	a.SearchQuery = ""
	a.SearchResults = nil
}

// UpdateSearch updates search results based on query.
func (a *App) UpdateSearch() {
	if a.SearchQuery == "" {
		a.SearchResults = nil
		return
	}

	query := a.SearchQuery
	a.SearchResults = nil

	for _, tab := range a.Tabs {
		if tab.Error != nil || tab.Loading {
			continue
		}
		for _, article := range tab.Articles {
			// Simple case-insensitive search in title and description
			titleLower := article.Title
			descLower := article.Summary()
			queryLower := query

			if containsIgnoreCase(titleLower, queryLower) || containsIgnoreCase(descLower, queryLower) {
				a.SearchResults = append(a.SearchResults, SummaryArticle{
					SourceName: tab.Source.Name,
					Article:    article,
				})
			}
		}
	}

	// Sort by relevance (title matches first) then by date
	sort.Slice(a.SearchResults, func(i, j int) bool {
		iTitle := containsIgnoreCase(a.SearchResults[i].Article.Title, query)
		jTitle := containsIgnoreCase(a.SearchResults[j].Article.Title, query)
		if iTitle != jTitle {
			return iTitle
		}
		return a.SearchResults[i].Article.PubDate.After(a.SearchResults[j].Article.PubDate)
	})
}

// containsIgnoreCase checks if s contains substr (case-insensitive).
func containsIgnoreCase(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		len(substr) == 0 ||
		(len(s) > 0 && containsIgnoreCaseHelper(s, substr)))
}

func containsIgnoreCaseHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if equalIgnoreCase(s[i:i+len(substr)], substr) {
			return true
		}
	}
	return false
}

func equalIgnoreCase(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		ca, cb := a[i], b[i]
		if ca >= 'A' && ca <= 'Z' {
			ca += 'a' - 'A'
		}
		if cb >= 'A' && cb <= 'Z' {
			cb += 'a' - 'A'
		}
		if ca != cb {
			return false
		}
	}
	return true
}

// ToggleBookmark toggles bookmark status for current article.
func (a *App) ToggleBookmark() {
	if len(a.Tabs) == 0 || a.CurrentTab >= len(a.Tabs) {
		return
	}
	tab := a.Tabs[a.CurrentTab]
	if tab.Selected < 0 || tab.Selected >= len(tab.Articles) {
		return
	}

	article := &tab.Articles[tab.Selected]
	article.Bookmarked = !article.Bookmarked

	if article.Bookmarked {
		a.Bookmarks = append(a.Bookmarks, SummaryArticle{
			SourceName: tab.Source.Name,
			Article:    *article,
		})
		a.SetStatusMessage("Article bookmarked")
	} else {
		// Remove from bookmarks
		for i, b := range a.Bookmarks {
			if b.Article.GUID == article.GUID {
				a.Bookmarks = append(a.Bookmarks[:i], a.Bookmarks[i+1:]...)
				break
			}
		}
		a.SetStatusMessage("Bookmark removed")
	}
}

// EnterHelpMode enters help view.
func (a *App) EnterHelpMode() {
	a.ViewMode = ViewModeHelp
	a.HelpScrollOffset = 0
}

// ExitHelpMode exits help view.
func (a *App) ExitHelpMode() {
	a.ViewMode = ViewModeList
}

// OpenInBrowser returns the URL of the currently selected article for opening.
func (a *App) OpenInBrowser() string {
	if len(a.Tabs) == 0 || a.CurrentTab >= len(a.Tabs) {
		return ""
	}
	tab := a.Tabs[a.CurrentTab]
	if tab.Selected < 0 || tab.Selected >= len(tab.Articles) {
		return ""
	}
	return tab.Articles[tab.Selected].Link
}

// GetCurrentArticle returns the currently selected article.
func (a *App) GetCurrentArticle() *feed.Article {
	if len(a.Tabs) == 0 || a.CurrentTab >= len(a.Tabs) {
		return nil
	}
	tab := a.Tabs[a.CurrentTab]
	if tab.Selected < 0 || tab.Selected >= len(tab.Articles) {
		return nil
	}
	return &tab.Articles[tab.Selected]
}

// MarkAllRead marks all articles in the current tab as read.
func (a *App) MarkAllRead() {
	if len(a.Tabs) == 0 || a.CurrentTab >= len(a.Tabs) {
		return
	}
	tab := a.Tabs[a.CurrentTab]
	for i := range tab.Articles {
		tab.Articles[i].Read = true
		if tab.Articles[i].GUID != "" {
			a.ReadArticles[tab.Articles[i].GUID] = true
		}
	}
	a.SetStatusMessage("All articles marked as read")
}
