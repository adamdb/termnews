// Package config handles application configuration loading, saving, and defaults.
package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// FeedType represents the type of feed (RSS or Atom).
type FeedType string

const (
	FeedTypeRSS  FeedType = "rss"
	FeedTypeAtom FeedType = "atom"
	FeedTypeAuto FeedType = "auto" // Auto-detect feed type
)

// Source represents a single news feed source.
type Source struct {
	Name     string   `toml:"name"`
	URL      string   `toml:"url"`
	FeedType FeedType `toml:"feed_type"`
	Category string   `toml:"category,omitempty"` // Optional category for organizing feeds
}

// Keybinding represents a single keybinding configuration.
type Keybinding struct {
	Key    string `toml:"key"`              // Primary key (e.g., "q", "ctrl+c", "enter")
	Alt    string `toml:"alt,omitempty"`    // Alternative key binding
	Action string `toml:"action,omitempty"` // Action name (for documentation)
}

// Keybindings holds all configurable key bindings.
type Keybindings struct {
	Quit          Keybinding `toml:"quit,omitempty"`
	NextTab       Keybinding `toml:"next_tab,omitempty"`
	PrevTab       Keybinding `toml:"prev_tab,omitempty"`
	ScrollDown    Keybinding `toml:"scroll_down,omitempty"`
	ScrollUp      Keybinding `toml:"scroll_up,omitempty"`
	Select        Keybinding `toml:"select,omitempty"`
	Back          Keybinding `toml:"back,omitempty"`
	Refresh       Keybinding `toml:"refresh,omitempty"`
	RefreshAll    Keybinding `toml:"refresh_all,omitempty"`
	ToggleAuto    Keybinding `toml:"toggle_auto,omitempty"`
	Summary       Keybinding `toml:"summary,omitempty"`
	Search        Keybinding `toml:"search,omitempty"`
	AddSource     Keybinding `toml:"add_source,omitempty"`
	DeleteSource  Keybinding `toml:"delete_source,omitempty"`
	Bookmark      Keybinding `toml:"bookmark,omitempty"`
	MarkRead      Keybinding `toml:"mark_read,omitempty"`
	OpenBrowser   Keybinding `toml:"open_browser,omitempty"`
	Help          Keybinding `toml:"help,omitempty"`
}

// DefaultKeybindings returns the default keybinding configuration.
func DefaultKeybindings() Keybindings {
	return Keybindings{
		Quit:          Keybinding{Key: "q", Alt: "ctrl+c", Action: "Quit application"},
		NextTab:       Keybinding{Key: "tab", Alt: "right", Action: "Next tab"},
		PrevTab:       Keybinding{Key: "shift+tab", Alt: "left", Action: "Previous tab"},
		ScrollDown:    Keybinding{Key: "down", Alt: "j", Action: "Scroll down"},
		ScrollUp:      Keybinding{Key: "up", Alt: "k", Action: "Scroll up"},
		Select:        Keybinding{Key: "enter", Action: "Select/open article"},
		Back:          Keybinding{Key: "esc", Alt: "backspace", Action: "Go back"},
		Refresh:       Keybinding{Key: "r", Action: "Refresh current feed"},
		RefreshAll:    Keybinding{Key: "R", Action: "Refresh all feeds"},
		ToggleAuto:    Keybinding{Key: "t", Action: "Toggle auto-refresh"},
		Summary:       Keybinding{Key: "s", Action: "Summary view"},
		Search:        Keybinding{Key: "/", Action: "Search articles"},
		AddSource:     Keybinding{Key: "a", Action: "Add new source"},
		DeleteSource:  Keybinding{Key: "d", Action: "Delete current source"},
		Bookmark:      Keybinding{Key: "b", Action: "Toggle bookmark"},
		MarkRead:      Keybinding{Key: "m", Action: "Mark all as read"},
		OpenBrowser:   Keybinding{Key: "o", Action: "Open in browser"},
		Help:          Keybinding{Key: "?", Action: "Show help"},
	}
}

// MergeKeybindings merges user keybinding overrides with defaults.
// Any keybinding with an empty Key in the override uses the base value.
func MergeKeybindings(base, override Keybindings) Keybindings {
	result := base
	if override.Quit.Key != "" {
		result.Quit = override.Quit
	}
	if override.NextTab.Key != "" {
		result.NextTab = override.NextTab
	}
	if override.PrevTab.Key != "" {
		result.PrevTab = override.PrevTab
	}
	if override.ScrollDown.Key != "" {
		result.ScrollDown = override.ScrollDown
	}
	if override.ScrollUp.Key != "" {
		result.ScrollUp = override.ScrollUp
	}
	if override.Select.Key != "" {
		result.Select = override.Select
	}
	if override.Back.Key != "" {
		result.Back = override.Back
	}
	if override.Refresh.Key != "" {
		result.Refresh = override.Refresh
	}
	if override.RefreshAll.Key != "" {
		result.RefreshAll = override.RefreshAll
	}
	if override.ToggleAuto.Key != "" {
		result.ToggleAuto = override.ToggleAuto
	}
	if override.Summary.Key != "" {
		result.Summary = override.Summary
	}
	if override.Search.Key != "" {
		result.Search = override.Search
	}
	if override.AddSource.Key != "" {
		result.AddSource = override.AddSource
	}
	if override.DeleteSource.Key != "" {
		result.DeleteSource = override.DeleteSource
	}
	if override.Bookmark.Key != "" {
		result.Bookmark = override.Bookmark
	}
	if override.MarkRead.Key != "" {
		result.MarkRead = override.MarkRead
	}
	if override.OpenBrowser.Key != "" {
		result.OpenBrowser = override.OpenBrowser
	}
	if override.Help.Key != "" {
		result.Help = override.Help
	}
	return result
}

// Config holds all application configuration.
type Config struct {
	Sources             []Source    `toml:"sources"`
	RefreshIntervalMins int         `toml:"refresh_interval_mins,omitempty"` // Auto-refresh interval in minutes
	MaxArticlesPerFeed  int         `toml:"max_articles_per_feed,omitempty"` // Max articles to show per feed
	ThemeName           string      `toml:"theme,omitempty"`                 // Built-in theme: "default", "bitchx", "hacker", "minimal", "retro"
	ShowReadStatus      bool        `toml:"show_read_status,omitempty"`      // Track read/unread status
	CustomTheme         Theme       `toml:"custom_theme,omitempty"`          // Custom theme overrides
	Keybindings         Keybindings `toml:"keybindings,omitempty"`           // Custom keybindings

	// Computed/resolved theme (not persisted)
	ResolvedTheme Theme `toml:"-"`
}

// ConfigPath returns the path to the config file.
func ConfigPath() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = "."
	}
	return filepath.Join(configDir, "termnews", "config.toml")
}

// Load loads the configuration from disk or returns default config.
func Load() *Config {
	path := ConfigPath()
	cfg := &Config{}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return DefaultConfig()
	}

	if _, err := toml.DecodeFile(path, cfg); err != nil {
		return DefaultConfig()
	}

	// Set defaults for missing values
	if cfg.RefreshIntervalMins <= 0 {
		cfg.RefreshIntervalMins = 5
	}
	if cfg.MaxArticlesPerFeed <= 0 {
		cfg.MaxArticlesPerFeed = 50
	}
	if cfg.ThemeName == "" {
		cfg.ThemeName = "default"
	}

	// Resolve the theme: start with built-in, merge with custom overrides
	baseTheme := GetBuiltinTheme(cfg.ThemeName)
	cfg.ResolvedTheme = MergeTheme(baseTheme, cfg.CustomTheme)

	// Merge default keybindings with user overrides
	// This ensures any unset keybindings get default values
	cfg.Keybindings = MergeKeybindings(DefaultKeybindings(), cfg.Keybindings)

	return cfg
}

// Save writes the configuration to disk.
func (c *Config) Save() error {
	path := ConfigPath()

	// Create parent directories if needed
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := toml.NewEncoder(f)
	return encoder.Encode(c)
}

// AddSource adds a new feed source to the configuration.
func (c *Config) AddSource(name, url string, feedType FeedType, category string) {
	c.Sources = append(c.Sources, Source{
		Name:     name,
		URL:      url,
		FeedType: feedType,
		Category: category,
	})
}

// RemoveSource removes a feed source at the given index.
func (c *Config) RemoveSource(index int) {
	if index >= 0 && index < len(c.Sources) {
		c.Sources = append(c.Sources[:index], c.Sources[index+1:]...)
	}
}

// DefaultConfig returns the default configuration with popular news sources.
func DefaultConfig() *Config {
	cfg := &Config{
		RefreshIntervalMins: 5,
		MaxArticlesPerFeed:  50,
		ThemeName:           "default",
		ShowReadStatus:      true,
		Keybindings:         DefaultKeybindings(),
		Sources: []Source{
			{
				Name:     "Reuters Top News",
				URL:      "http://feeds.reuters.com/reuters/topNews",
				FeedType: FeedTypeRSS,
				Category: "World News",
			},
			{
				Name:     "BBC News",
				URL:      "http://feeds.bbci.co.uk/news/rss.xml",
				FeedType: FeedTypeRSS,
				Category: "World News",
			},
			{
				Name:     "Hacker News",
				URL:      "https://news.ycombinator.com/rss",
				FeedType: FeedTypeRSS,
				Category: "Tech",
			},
			{
				Name:     "NPR News",
				URL:      "https://feeds.npr.org/1001/rss.xml",
				FeedType: FeedTypeRSS,
				Category: "World News",
			},
		},
	}
	// Resolve the default theme
	cfg.ResolvedTheme = DefaultTheme()
	return cfg
}
