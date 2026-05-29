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

// Config holds all application configuration.
type Config struct {
	Sources             []Source `toml:"sources"`
	RefreshIntervalMins int      `toml:"refresh_interval_mins,omitempty"` // Auto-refresh interval in minutes
	MaxArticlesPerFeed  int      `toml:"max_articles_per_feed,omitempty"` // Max articles to show per feed
	Theme               string   `toml:"theme,omitempty"`                 // Color theme: "default", "dark", "light"
	ShowReadStatus      bool     `toml:"show_read_status,omitempty"`      // Track read/unread status
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
	if cfg.Theme == "" {
		cfg.Theme = "default"
	}

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
	return &Config{
		RefreshIntervalMins: 5,
		MaxArticlesPerFeed:  50,
		Theme:               "default",
		ShowReadStatus:      true,
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
}
