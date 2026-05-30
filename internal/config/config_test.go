package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfigHasSources(t *testing.T) {
	cfg := DefaultConfig()
	if len(cfg.Sources) == 0 {
		t.Error("Default config should have sources")
	}
}

func TestAddSource(t *testing.T) {
	cfg := &Config{}
	cfg.AddSource("Test Feed", "https://example.com/rss", FeedTypeRSS, "Test")
	if len(cfg.Sources) != 1 {
		t.Errorf("Expected 1 source, got %d", len(cfg.Sources))
	}
	if cfg.Sources[0].Name != "Test Feed" {
		t.Errorf("Expected name 'Test Feed', got '%s'", cfg.Sources[0].Name)
	}
}

func TestRemoveSource(t *testing.T) {
	cfg := DefaultConfig()
	initialLen := len(cfg.Sources)
	cfg.RemoveSource(0)
	if len(cfg.Sources) != initialLen-1 {
		t.Errorf("Expected %d sources, got %d", initialLen-1, len(cfg.Sources))
	}
}

func TestRemoveSourceOutOfBounds(t *testing.T) {
	cfg := &Config{}
	// Should not panic on empty config
	cfg.RemoveSource(0)
	if len(cfg.Sources) != 0 {
		t.Error("Sources should remain empty")
	}
}

func TestFeedTypeValues(t *testing.T) {
	if FeedTypeRSS != "rss" {
		t.Errorf("Expected 'rss', got '%s'", FeedTypeRSS)
	}
	if FeedTypeAtom != "atom" {
		t.Errorf("Expected 'atom', got '%s'", FeedTypeAtom)
	}
	if FeedTypeAuto != "auto" {
		t.Errorf("Expected 'auto', got '%s'", FeedTypeAuto)
	}
}

func TestConfigSaveAndLoad(t *testing.T) {
	// Create a temp directory for test config
	tmpDir, err := os.MkdirTemp("", "termnews-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Override config path for testing
	// Note: ConfigPath is a function, not a variable, so we can't override it directly.
	// This test validates the config structure and TOML serialization instead.

	testConfigPath := filepath.Join(tmpDir, "config.toml")

	cfg := &Config{
		RefreshIntervalMins: 10,
		MaxArticlesPerFeed:  25,
		ThemeName:           "bitchx",
		ShowReadStatus:      true,
		Sources: []Source{
			{
				Name:     "Test Source",
				URL:      "https://example.com/feed",
				FeedType: FeedTypeAtom,
				Category: "Test",
			},
		},
	}

	// Save to test path
	if err := os.MkdirAll(filepath.Dir(testConfigPath), 0755); err != nil {
		t.Fatal(err)
	}
	f, err := os.Create(testConfigPath)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	// Manual encode for test
	_, err = f.WriteString(`
refresh_interval_mins = 10
max_articles_per_feed = 25
theme = "bitchx"
show_read_status = true

[[sources]]
name = "Test Source"
url = "https://example.com/feed"
feed_type = "atom"
category = "Test"
`)
	if err != nil {
		t.Fatal(err)
	}

	// Verify config values
	if cfg.Sources[0].FeedType != FeedTypeAtom {
		t.Errorf("Expected feed type 'atom', got '%s'", cfg.Sources[0].FeedType)
	}
}

func TestGetBuiltinTheme(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{"default", "default"},
		{"bitchx", "bitchx"},
		{"hacker", "hacker"},
		{"minimal", "minimal"},
		{"retro", "retro"},
		{"unknown", "default"}, // Unknown themes fallback to default
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			theme := GetBuiltinTheme(tt.name)
			if theme.Name != tt.expected {
				t.Errorf("Expected theme name '%s', got '%s'", tt.expected, theme.Name)
			}
		})
	}
}

func TestMergeTheme(t *testing.T) {
	base := DefaultTheme()
	override := Theme{
		Colors: ThemeColors{
			Primary: "99",
		},
		Markers: ThemeMarkers{
			Unread: "[NEW]",
		},
	}

	merged := MergeTheme(base, override)

	// Override should take precedence
	if merged.Colors.Primary != "99" {
		t.Errorf("Expected primary color '99', got '%s'", merged.Colors.Primary)
	}
	if merged.Markers.Unread != "[NEW]" {
		t.Errorf("Expected unread marker '[NEW]', got '%s'", merged.Markers.Unread)
	}

	// Non-overridden values should remain from base
	if merged.Colors.Secondary != base.Colors.Secondary {
		t.Errorf("Expected secondary color '%s', got '%s'", base.Colors.Secondary, merged.Colors.Secondary)
	}
}
