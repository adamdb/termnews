package app

import (
	"testing"

	"github.com/adamdb/termnews/internal/config"
	"github.com/adamdb/termnews/internal/feed"
)

func makeArticle(title string) feed.Article {
	return feed.Article{
		Title:       title,
		Link:        "https://example.com",
		Description: "<p>Content</p>",
		Author:      "Author",
		GUID:        title,
	}
}

func makeAppWithTabs(count int) *App {
	cfg := &config.Config{}
	for i := 0; i < count; i++ {
		cfg.AddSource("Source "+string(rune('0'+i)), "https://example.com/rss", config.FeedTypeRSS, "")
	}
	return New(cfg)
}

func TestNextTabWraps(t *testing.T) {
	app := makeAppWithTabs(3)
	if app.CurrentTab != 0 {
		t.Errorf("Expected current tab 0, got %d", app.CurrentTab)
	}
	app.NextTab()
	if app.CurrentTab != 1 {
		t.Errorf("Expected current tab 1, got %d", app.CurrentTab)
	}
	app.NextTab()
	if app.CurrentTab != 2 {
		t.Errorf("Expected current tab 2, got %d", app.CurrentTab)
	}
	app.NextTab()
	if app.CurrentTab != 0 {
		t.Errorf("Expected current tab 0 (wrapped), got %d", app.CurrentTab)
	}
}

func TestPrevTabWraps(t *testing.T) {
	app := makeAppWithTabs(3)
	if app.CurrentTab != 0 {
		t.Errorf("Expected current tab 0, got %d", app.CurrentTab)
	}
	app.PrevTab()
	if app.CurrentTab != 2 {
		t.Errorf("Expected current tab 2 (wrapped), got %d", app.CurrentTab)
	}
	app.PrevTab()
	if app.CurrentTab != 1 {
		t.Errorf("Expected current tab 1, got %d", app.CurrentTab)
	}
}

func TestScrollDownSelectsArticles(t *testing.T) {
	app := makeAppWithTabs(1)
	articles := []feed.Article{
		makeArticle("Article 1"),
		makeArticle("Article 2"),
		makeArticle("Article 3"),
	}
	app.SetArticles(0, articles, nil)
	if app.Tabs[0].Selected != 0 {
		t.Errorf("Expected selected 0, got %d", app.Tabs[0].Selected)
	}
	app.ScrollDown()
	if app.Tabs[0].Selected != 1 {
		t.Errorf("Expected selected 1, got %d", app.Tabs[0].Selected)
	}
	app.ScrollDown()
	if app.Tabs[0].Selected != 2 {
		t.Errorf("Expected selected 2, got %d", app.Tabs[0].Selected)
	}
	// Should not go past last
	app.ScrollDown()
	if app.Tabs[0].Selected != 2 {
		t.Errorf("Expected selected 2 (no change), got %d", app.Tabs[0].Selected)
	}
}

func TestScrollUpDoesNotUnderflow(t *testing.T) {
	app := makeAppWithTabs(1)
	articles := []feed.Article{
		makeArticle("Article 1"),
		makeArticle("Article 2"),
	}
	app.SetArticles(0, articles, nil)
	if app.Tabs[0].Selected != 0 {
		t.Errorf("Expected selected 0, got %d", app.Tabs[0].Selected)
	}
	app.ScrollUp()
	if app.Tabs[0].Selected != 0 {
		t.Errorf("Expected selected 0 (no change), got %d", app.Tabs[0].Selected)
	}
}

func TestSelectArticleEntersDetailMode(t *testing.T) {
	app := makeAppWithTabs(1)
	articles := []feed.Article{makeArticle("Article 1")}
	app.SetArticles(0, articles, nil)
	if app.ViewMode != ViewModeList {
		t.Errorf("Expected ViewModeList, got %d", app.ViewMode)
	}
	app.SelectArticle()
	if app.ViewMode != ViewModeDetail {
		t.Errorf("Expected ViewModeDetail, got %d", app.ViewMode)
	}
}

func TestBackToList(t *testing.T) {
	app := makeAppWithTabs(1)
	articles := []feed.Article{makeArticle("A")}
	app.SetArticles(0, articles, nil)
	app.SelectArticle()
	if app.ViewMode != ViewModeDetail {
		t.Errorf("Expected ViewModeDetail, got %d", app.ViewMode)
	}
	app.BackToList()
	if app.ViewMode != ViewModeList {
		t.Errorf("Expected ViewModeList, got %d", app.ViewMode)
	}
}

func TestAddSourceFlow(t *testing.T) {
	app := makeAppWithTabs(1)
	app.StartAddSource()
	if app.ViewMode != ViewModeAddSource {
		t.Errorf("Expected ViewModeAddSource, got %d", app.ViewMode)
	}
	app.AddSourceState.Name = "New Feed"
	app.AddSourceState.URL = "https://new.example.com/rss"
	result := app.ConfirmAddSource()
	if !result {
		t.Error("Expected ConfirmAddSource to return true")
	}
	if len(app.Tabs) != 2 {
		t.Errorf("Expected 2 tabs, got %d", len(app.Tabs))
	}
	if app.CurrentTab != 1 {
		t.Errorf("Expected current tab 1, got %d", app.CurrentTab)
	}
	if app.ViewMode != ViewModeList {
		t.Errorf("Expected ViewModeList, got %d", app.ViewMode)
	}
}

func TestConfirmAddSourceFailsOnEmpty(t *testing.T) {
	app := makeAppWithTabs(1)
	app.StartAddSource()
	result := app.ConfirmAddSource()
	if result {
		t.Error("Expected ConfirmAddSource to return false")
	}
	if app.AddSourceState.Error == "" {
		t.Error("Expected error message")
	}
}

func TestRemoveCurrentSource(t *testing.T) {
	app := makeAppWithTabs(3)
	if len(app.Tabs) != 3 {
		t.Errorf("Expected 3 tabs, got %d", len(app.Tabs))
	}
	app.RemoveCurrentSource()
	if len(app.Tabs) != 2 {
		t.Errorf("Expected 2 tabs, got %d", len(app.Tabs))
	}
}

func TestRefreshCurrent(t *testing.T) {
	app := makeAppWithTabs(1)
	articles := []feed.Article{makeArticle("A"), makeArticle("B")}
	app.SetArticles(0, articles, nil)
	if len(app.Tabs[0].Articles) != 2 {
		t.Errorf("Expected 2 articles, got %d", len(app.Tabs[0].Articles))
	}
	app.RefreshCurrent()
	if !app.Tabs[0].Loading {
		t.Error("Expected loading to be true")
	}
	if len(app.Tabs[0].Articles) != 0 {
		t.Errorf("Expected 0 articles, got %d", len(app.Tabs[0].Articles))
	}
}

func TestSetArticlesError(t *testing.T) {
	app := makeAppWithTabs(1)
	err := feed.FetchError{Type: "network", Message: "timeout"}
	app.SetArticles(0, nil, err)
	if app.Tabs[0].Loading {
		t.Error("Expected loading to be false")
	}
	if app.Tabs[0].Error == nil {
		t.Error("Expected error to be set")
	}
}

func TestAddSourceStateFeedTypeToggle(t *testing.T) {
	state := NewAddSourceState()
	if state.FeedTypeIndex != 2 {
		t.Errorf("Expected feed type index 2 (auto), got %d", state.FeedTypeIndex)
	}
	if state.SelectedFeedType() != config.FeedTypeAuto {
		t.Errorf("Expected FeedTypeAuto, got %s", state.SelectedFeedType())
	}
	state.FeedTypeIndex = 0
	if state.SelectedFeedType() != config.FeedTypeRSS {
		t.Errorf("Expected FeedTypeRSS, got %s", state.SelectedFeedType())
	}
	state.FeedTypeIndex = 1
	if state.SelectedFeedType() != config.FeedTypeAtom {
		t.Errorf("Expected FeedTypeAtom, got %s", state.SelectedFeedType())
	}
}

func TestToggleAutoRefresh(t *testing.T) {
	app := makeAppWithTabs(1)
	if app.AutoRefreshEnabled {
		t.Error("Expected auto refresh to be disabled initially")
	}
	app.ToggleAutoRefresh()
	if !app.AutoRefreshEnabled {
		t.Error("Expected auto refresh to be enabled")
	}
	app.ToggleAutoRefresh()
	if app.AutoRefreshEnabled {
		t.Error("Expected auto refresh to be disabled")
	}
}

func TestContainsIgnoreCase(t *testing.T) {
	tests := []struct {
		s      string
		substr string
		want   bool
	}{
		{"Hello World", "world", true},
		{"Hello World", "HELLO", true},
		{"Hello World", "xyz", false},
		{"", "", true},
		{"abc", "", true},
		{"", "a", false},
	}

	for _, tt := range tests {
		got := containsIgnoreCase(tt.s, tt.substr)
		if got != tt.want {
			t.Errorf("containsIgnoreCase(%q, %q) = %v, want %v", tt.s, tt.substr, got, tt.want)
		}
	}
}

func TestToggleBookmark(t *testing.T) {
	app := makeAppWithTabs(1)
	articles := []feed.Article{makeArticle("Article 1")}
	app.SetArticles(0, articles, nil)

	app.ToggleBookmark()
	if !app.Tabs[0].Articles[0].Bookmarked {
		t.Error("Expected article to be bookmarked")
	}
	if len(app.Bookmarks) != 1 {
		t.Errorf("Expected 1 bookmark, got %d", len(app.Bookmarks))
	}

	app.ToggleBookmark()
	if app.Tabs[0].Articles[0].Bookmarked {
		t.Error("Expected article to not be bookmarked")
	}
	if len(app.Bookmarks) != 0 {
		t.Errorf("Expected 0 bookmarks, got %d", len(app.Bookmarks))
	}
}
