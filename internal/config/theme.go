// Package config handles application configuration loading, saving, and defaults.
package config

// Theme defines the complete visual appearance of the application.
// Inspired by classic terminal applications like BitchX, with ANSI art
// and highly customizable colors and visual elements.
type Theme struct {
	// Theme name identifier
	Name string `toml:"name,omitempty"`

	// Color palette for all UI elements (ANSI color codes 0-255 or hex)
	Colors ThemeColors `toml:"colors,omitempty"`

	// Border and box drawing characters
	Borders ThemeBorders `toml:"borders,omitempty"`

	// Text decorations and markers
	Markers ThemeMarkers `toml:"markers,omitempty"`

	// Header/banner configuration
	Header ThemeHeader `toml:"header,omitempty"`

	// Status bar configuration
	StatusBar ThemeStatusBar `toml:"status_bar,omitempty"`

	// Article list styling
	ArticleList ThemeArticleList `toml:"article_list,omitempty"`

	// Tab bar styling
	TabBar ThemeTabBar `toml:"tab_bar,omitempty"`
}

// ThemeColors defines the color palette for all UI elements.
type ThemeColors struct {
	// Primary accent color (tabs, titles, borders)
	Primary string `toml:"primary,omitempty"`
	// Secondary accent color (highlights, selections)
	Secondary string `toml:"secondary,omitempty"`
	// Normal text color
	Text string `toml:"text,omitempty"`
	// Muted/dim text color
	TextMuted string `toml:"text_muted,omitempty"`
	// Background color (empty string for terminal default)
	Background string `toml:"background,omitempty"`
	// Error/warning color
	Error string `toml:"error,omitempty"`
	// Success/confirmation color
	Success string `toml:"success,omitempty"`
	// Link color
	Link string `toml:"link,omitempty"`
	// Unread indicator color
	Unread string `toml:"unread,omitempty"`
	// Bookmarked indicator color
	Bookmark string `toml:"bookmark,omitempty"`
	// Category/tag color
	Category string `toml:"category,omitempty"`
	// Border color
	Border string `toml:"border,omitempty"`
	// Header/banner color
	Header string `toml:"header,omitempty"`
	// Status bar background
	StatusBg string `toml:"status_bg,omitempty"`
	// Status bar text
	StatusFg string `toml:"status_fg,omitempty"`
}

// ThemeBorders defines box drawing and border characters.
type ThemeBorders struct {
	// Border style: "rounded", "sharp", "double", "heavy", "ascii", "none", or "custom"
	Style string `toml:"style,omitempty"`
	// Custom border characters (used when Style is "custom")
	TopLeft     string `toml:"top_left,omitempty"`
	TopRight    string `toml:"top_right,omitempty"`
	BottomLeft  string `toml:"bottom_left,omitempty"`
	BottomRight string `toml:"bottom_right,omitempty"`
	Horizontal  string `toml:"horizontal,omitempty"`
	Vertical    string `toml:"vertical,omitempty"`
	// Separator line character
	Separator string `toml:"separator,omitempty"`
}

// ThemeMarkers defines text decorations and indicator symbols.
type ThemeMarkers struct {
	// Unread article indicator (e.g., "●", "•", "*", "NEW")
	Unread string `toml:"unread,omitempty"`
	// Bookmarked indicator (e.g., "★", "*", "[B]")
	Bookmark string `toml:"bookmark,omitempty"`
	// Loading indicator (e.g., "⟳", "...", "[~]")
	Loading string `toml:"loading,omitempty"`
	// Error indicator (e.g., "✗", "X", "[!]")
	Error string `toml:"error,omitempty"`
	// Selected/cursor indicator (e.g., "▶", ">", "->", "»")
	Selected string `toml:"selected,omitempty"`
	// Tab separator (e.g., "|", "│", "┃", " ")
	TabSeparator string `toml:"tab_separator,omitempty"`
	// Active tab brackets (e.g., ["[", "]"], ["<", ">"], ["«", "»"])
	TabActiveLeft  string `toml:"tab_active_left,omitempty"`
	TabActiveRight string `toml:"tab_active_right,omitempty"`
}

// ThemeHeader defines the application header/banner configuration.
type ThemeHeader struct {
	// Show the header banner
	Show bool `toml:"show,omitempty"`
	// ASCII art banner (multiline string)
	Banner string `toml:"banner,omitempty"`
	// Show the application name
	ShowName bool `toml:"show_name,omitempty"`
	// Custom application name/title
	Title string `toml:"title,omitempty"`
	// Alignment: "left", "center", "right"
	Align string `toml:"align,omitempty"`
}

// ThemeStatusBar defines the status bar appearance.
type ThemeStatusBar struct {
	// Show the status bar
	Show bool `toml:"show,omitempty"`
	// Position: "top" or "bottom"
	Position string `toml:"position,omitempty"`
	// Custom help text format (supports placeholders)
	HelpFormat string `toml:"help_format,omitempty"`
	// Show key bindings in status bar
	ShowKeys bool `toml:"show_keys,omitempty"`
	// Separator between key hints
	KeySeparator string `toml:"key_separator,omitempty"`
}

// ThemeArticleList defines article list appearance.
type ThemeArticleList struct {
	// Show article numbers/indices
	ShowNumbers bool `toml:"show_numbers,omitempty"`
	// Number format (e.g., "%d.", "%02d)", "[%d]")
	NumberFormat string `toml:"number_format,omitempty"`
	// Show horizontal separators between articles
	ShowSeparators bool `toml:"show_separators,omitempty"`
	// Separator style (character)
	SeparatorChar string `toml:"separator_char,omitempty"`
	// Show article metadata (date, author)
	ShowMeta bool `toml:"show_meta,omitempty"`
	// Meta format
	MetaFormat string `toml:"meta_format,omitempty"`
	// Compact mode (single line per article)
	Compact bool `toml:"compact,omitempty"`
}

// ThemeTabBar defines the tab bar appearance.
type ThemeTabBar struct {
	// Position: "top" or "bottom"
	Position string `toml:"position,omitempty"`
	// Show unread counts in tabs
	ShowUnreadCount bool `toml:"show_unread_count,omitempty"`
	// Unread count format (e.g., "(%d)", " [%d]", " •%d")
	UnreadFormat string `toml:"unread_format,omitempty"`
	// Max tab name length (0 = no limit)
	MaxTabWidth int `toml:"max_tab_width,omitempty"`
	// Truncation suffix
	TruncateSuffix string `toml:"truncate_suffix,omitempty"`
}

// DefaultTheme returns the default modern theme.
func DefaultTheme() Theme {
	return Theme{
		Name: "default",
		Colors: ThemeColors{
			Primary:   "6",   // Cyan
			Secondary: "11",  // Yellow
			Text:      "15",  // White
			TextMuted: "240", // Dark gray
			Error:     "1",   // Red
			Success:   "2",   // Green
			Link:      "4",   // Blue
			Unread:    "2",   // Green
			Bookmark:  "11",  // Yellow
			Category:  "6",   // Cyan
			Border:    "6",   // Cyan
			Header:    "6",   // Cyan
			StatusBg:  "",    // Default
			StatusFg:  "240", // Dark gray
		},
		Borders: ThemeBorders{
			Style:     "rounded",
			Separator: "─",
		},
		Markers: ThemeMarkers{
			Unread:         "●",
			Bookmark:       "★",
			Loading:        "⟳",
			Error:          "✗",
			Selected:       "▶",
			TabSeparator:   "|",
			TabActiveLeft:  " ",
			TabActiveRight: " ",
		},
		Header: ThemeHeader{
			Show:     true,
			ShowName: true,
			Title:    "termnews",
			Align:    "left",
		},
		StatusBar: ThemeStatusBar{
			Show:         true,
			Position:     "bottom",
			ShowKeys:     true,
			KeySeparator: "  ",
		},
		ArticleList: ThemeArticleList{
			ShowNumbers:    false,
			ShowSeparators: true,
			SeparatorChar:  "─",
			ShowMeta:       true,
			Compact:        false,
		},
		TabBar: ThemeTabBar{
			Position:        "top",
			ShowUnreadCount: true,
			UnreadFormat:    " (%d)",
			MaxTabWidth:     20,
			TruncateSuffix:  "…",
		},
	}
}

// BitchXTheme returns a theme inspired by the classic BitchX IRC client.
func BitchXTheme() Theme {
	return Theme{
		Name: "bitchx",
		Colors: ThemeColors{
			Primary:   "14",  // Bright cyan
			Secondary: "10",  // Bright green
			Text:      "7",   // Silver/light gray
			TextMuted: "8",   // Dark gray
			Error:     "9",   // Bright red
			Success:   "10",  // Bright green
			Link:      "12",  // Bright blue
			Unread:    "10",  // Bright green
			Bookmark:  "11",  // Bright yellow
			Category:  "13",  // Bright magenta
			Border:    "14",  // Bright cyan
			Header:    "11",  // Bright yellow
			StatusBg:  "4",   // Blue background
			StatusFg:  "15",  // White text
		},
		Borders: ThemeBorders{
			Style:     "double",
			Separator: "═",
		},
		Markers: ThemeMarkers{
			Unread:         "[+]",
			Bookmark:       "[*]",
			Loading:        "[~]",
			Error:          "[!]",
			Selected:       "»»",
			TabSeparator:   "│",
			TabActiveLeft:  "«",
			TabActiveRight: "»",
		},
		Header: ThemeHeader{
			Show:     true,
			ShowName: true,
			Title:    "TERMNEWS",
			Align:    "center",
			Banner: `╔══════════════════════════════════════════════════════╗
║  ▀█▀ █▀▀ █▀█ █▀▄▀█ █▄░█ █▀▀ █░█░█ █▀                 ║
║  ░█░ ██▄ █▀▄ █░▀░█ █░▀█ ██▄ ▀▄▀▄▀ ▄█  v1.0          ║
╚══════════════════════════════════════════════════════╝`,
		},
		StatusBar: ThemeStatusBar{
			Show:         true,
			Position:     "bottom",
			ShowKeys:     true,
			KeySeparator: " │ ",
		},
		ArticleList: ThemeArticleList{
			ShowNumbers:    true,
			NumberFormat:   "%02d»",
			ShowSeparators: true,
			SeparatorChar:  "─",
			ShowMeta:       true,
			Compact:        false,
		},
		TabBar: ThemeTabBar{
			Position:        "top",
			ShowUnreadCount: true,
			UnreadFormat:    "[%d]",
			MaxTabWidth:     15,
			TruncateSuffix:  "..",
		},
	}
}

// HackerTheme returns a hacker/matrix style theme.
func HackerTheme() Theme {
	return Theme{
		Name: "hacker",
		Colors: ThemeColors{
			Primary:   "10",  // Bright green (matrix green)
			Secondary: "2",   // Green
			Text:      "10",  // Bright green
			TextMuted: "22",  // Dark green
			Error:     "9",   // Bright red
			Success:   "10",  // Bright green
			Link:      "46",  // Lime green
			Unread:    "118", // Light green
			Bookmark:  "226", // Yellow
			Category:  "82",  // Green
			Border:    "22",  // Dark green
			Header:    "10",  // Bright green
			StatusBg:  "22",  // Dark green
			StatusFg:  "10",  // Bright green
		},
		Borders: ThemeBorders{
			Style:     "ascii",
			Separator: "-",
		},
		Markers: ThemeMarkers{
			Unread:         "*",
			Bookmark:       "#",
			Loading:        "...",
			Error:          "ERR",
			Selected:       ">>",
			TabSeparator:   "|",
			TabActiveLeft:  "[",
			TabActiveRight: "]",
		},
		Header: ThemeHeader{
			Show:     true,
			ShowName: true,
			Title:    "TERMNEWS",
			Align:    "left",
			Banner: `
  _____                                         
 |_   _|__ _ __ _ __ ___  _ __   _____      _____
   | |/ _ \ '__| '_ ' _ \| '_ \ / _ \ \ /\ / / __|
   | |  __/ |  | | | | | | | | |  __/\ V  V /\__ \
   |_|\___|_|  |_| |_| |_|_| |_|\___| \_/\_/ |___/
`,
		},
		StatusBar: ThemeStatusBar{
			Show:         true,
			Position:     "bottom",
			ShowKeys:     true,
			KeySeparator: " | ",
		},
		ArticleList: ThemeArticleList{
			ShowNumbers:    true,
			NumberFormat:   "[%d]",
			ShowSeparators: true,
			SeparatorChar:  "-",
			ShowMeta:       true,
			Compact:        false,
		},
		TabBar: ThemeTabBar{
			Position:        "top",
			ShowUnreadCount: true,
			UnreadFormat:    ":%d",
			MaxTabWidth:     18,
			TruncateSuffix:  "~",
		},
	}
}

// MinimalTheme returns a clean, minimal theme.
func MinimalTheme() Theme {
	return Theme{
		Name: "minimal",
		Colors: ThemeColors{
			Primary:   "7",   // White
			Secondary: "15",  // Bright white
			Text:      "7",   // White
			TextMuted: "8",   // Gray
			Error:     "1",   // Red
			Success:   "2",   // Green
			Link:      "4",   // Blue
			Unread:    "7",   // White (subtle)
			Bookmark:  "3",   // Yellow
			Category:  "8",   // Gray
			Border:    "8",   // Gray
			Header:    "15",  // Bright white
			StatusBg:  "",    // Default
			StatusFg:  "8",   // Gray
		},
		Borders: ThemeBorders{
			Style:     "none",
			Separator: " ",
		},
		Markers: ThemeMarkers{
			Unread:         "·",
			Bookmark:       "★",
			Loading:        "...",
			Error:          "×",
			Selected:       ">",
			TabSeparator:   " ",
			TabActiveLeft:  "[",
			TabActiveRight: "]",
		},
		Header: ThemeHeader{
			Show:     true,
			ShowName: true,
			Title:    "termnews",
			Align:    "left",
		},
		StatusBar: ThemeStatusBar{
			Show:         true,
			Position:     "bottom",
			ShowKeys:     true,
			KeySeparator: " ",
		},
		ArticleList: ThemeArticleList{
			ShowNumbers:    false,
			ShowSeparators: false,
			ShowMeta:       true,
			Compact:        true,
		},
		TabBar: ThemeTabBar{
			Position:        "top",
			ShowUnreadCount: true,
			UnreadFormat:    " %d",
			MaxTabWidth:     0,
			TruncateSuffix:  "…",
		},
	}
}

// RetroTheme returns a retro BBS-style theme.
func RetroTheme() Theme {
	return Theme{
		Name: "retro",
		Colors: ThemeColors{
			Primary:   "14",  // Bright cyan
			Secondary: "11",  // Bright yellow
			Text:      "15",  // White
			TextMuted: "7",   // Silver
			Error:     "12",  // Bright red
			Success:   "10",  // Bright green
			Link:      "13",  // Bright magenta
			Unread:    "11",  // Bright yellow
			Bookmark:  "9",   // Bright red
			Category:  "14",  // Bright cyan
			Border:    "11",  // Bright yellow
			Header:    "14",  // Bright cyan
			StatusBg:  "1",   // Red background
			StatusFg:  "15",  // White
		},
		Borders: ThemeBorders{
			Style:     "heavy",
			Separator: "━",
		},
		Markers: ThemeMarkers{
			Unread:         "■",
			Bookmark:       "♦",
			Loading:        "◌",
			Error:          "▲",
			Selected:       "►",
			TabSeparator:   "┃",
			TabActiveLeft:  "◄",
			TabActiveRight: "►",
		},
		Header: ThemeHeader{
			Show:     true,
			ShowName: true,
			Title:    ":: TERMNEWS ::",
			Align:    "center",
			Banner: `
┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓
┃  ████████╗███████╗██████╗ ███╗   ███╗███╗   ██╗      ┃
┃  ╚══██╔══╝██╔════╝██╔══██╗████╗ ████║████╗  ██║      ┃
┃     ██║   █████╗  ██████╔╝██╔████╔██║██╔██╗ ██║      ┃
┃     ██║   ██╔══╝  ██╔══██╗██║╚██╔╝██║██║╚██╗██║      ┃
┃     ██║   ███████╗██║  ██║██║ ╚═╝ ██║██║ ╚████║EWS   ┃
┃     ╚═╝   ╚══════╝╚═╝  ╚═╝╚═╝     ╚═╝╚═╝  ╚═══╝      ┃
┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛
`,
		},
		StatusBar: ThemeStatusBar{
			Show:         true,
			Position:     "bottom",
			ShowKeys:     true,
			KeySeparator: " ┃ ",
		},
		ArticleList: ThemeArticleList{
			ShowNumbers:    true,
			NumberFormat:   "·%d·",
			ShowSeparators: true,
			SeparatorChar:  "─",
			ShowMeta:       true,
			Compact:        false,
		},
		TabBar: ThemeTabBar{
			Position:        "top",
			ShowUnreadCount: true,
			UnreadFormat:    "·%d·",
			MaxTabWidth:     16,
			TruncateSuffix:  "»",
		},
	}
}

// GetBuiltinTheme returns a built-in theme by name.
func GetBuiltinTheme(name string) Theme {
	switch name {
	case "bitchx":
		return BitchXTheme()
	case "hacker":
		return HackerTheme()
	case "minimal":
		return MinimalTheme()
	case "retro":
		return RetroTheme()
	default:
		return DefaultTheme()
	}
}

// MergeTheme merges user theme overrides with a base theme.
// Non-empty/non-zero values in override take precedence.
func MergeTheme(base, override Theme) Theme {
	result := base

	// Override name
	if override.Name != "" {
		result.Name = override.Name
	}

	// Merge colors
	result.Colors = mergeColors(base.Colors, override.Colors)

	// Merge borders
	result.Borders = mergeBorders(base.Borders, override.Borders)

	// Merge markers
	result.Markers = mergeMarkers(base.Markers, override.Markers)

	// Merge header
	result.Header = mergeHeader(base.Header, override.Header)

	// Merge status bar
	result.StatusBar = mergeStatusBar(base.StatusBar, override.StatusBar)

	// Merge article list
	result.ArticleList = mergeArticleList(base.ArticleList, override.ArticleList)

	// Merge tab bar
	result.TabBar = mergeTabBar(base.TabBar, override.TabBar)

	return result
}

func mergeColors(base, override ThemeColors) ThemeColors {
	if override.Primary != "" {
		base.Primary = override.Primary
	}
	if override.Secondary != "" {
		base.Secondary = override.Secondary
	}
	if override.Text != "" {
		base.Text = override.Text
	}
	if override.TextMuted != "" {
		base.TextMuted = override.TextMuted
	}
	if override.Background != "" {
		base.Background = override.Background
	}
	if override.Error != "" {
		base.Error = override.Error
	}
	if override.Success != "" {
		base.Success = override.Success
	}
	if override.Link != "" {
		base.Link = override.Link
	}
	if override.Unread != "" {
		base.Unread = override.Unread
	}
	if override.Bookmark != "" {
		base.Bookmark = override.Bookmark
	}
	if override.Category != "" {
		base.Category = override.Category
	}
	if override.Border != "" {
		base.Border = override.Border
	}
	if override.Header != "" {
		base.Header = override.Header
	}
	if override.StatusBg != "" {
		base.StatusBg = override.StatusBg
	}
	if override.StatusFg != "" {
		base.StatusFg = override.StatusFg
	}
	return base
}

func mergeBorders(base, override ThemeBorders) ThemeBorders {
	if override.Style != "" {
		base.Style = override.Style
	}
	if override.TopLeft != "" {
		base.TopLeft = override.TopLeft
	}
	if override.TopRight != "" {
		base.TopRight = override.TopRight
	}
	if override.BottomLeft != "" {
		base.BottomLeft = override.BottomLeft
	}
	if override.BottomRight != "" {
		base.BottomRight = override.BottomRight
	}
	if override.Horizontal != "" {
		base.Horizontal = override.Horizontal
	}
	if override.Vertical != "" {
		base.Vertical = override.Vertical
	}
	if override.Separator != "" {
		base.Separator = override.Separator
	}
	return base
}

func mergeMarkers(base, override ThemeMarkers) ThemeMarkers {
	if override.Unread != "" {
		base.Unread = override.Unread
	}
	if override.Bookmark != "" {
		base.Bookmark = override.Bookmark
	}
	if override.Loading != "" {
		base.Loading = override.Loading
	}
	if override.Error != "" {
		base.Error = override.Error
	}
	if override.Selected != "" {
		base.Selected = override.Selected
	}
	if override.TabSeparator != "" {
		base.TabSeparator = override.TabSeparator
	}
	if override.TabActiveLeft != "" {
		base.TabActiveLeft = override.TabActiveLeft
	}
	if override.TabActiveRight != "" {
		base.TabActiveRight = override.TabActiveRight
	}
	return base
}

func mergeHeader(base, override ThemeHeader) ThemeHeader {
	// For booleans, we check if the override has explicitly set them
	// by looking at the overall structure
	if override.Banner != "" {
		base.Banner = override.Banner
	}
	if override.Title != "" {
		base.Title = override.Title
	}
	if override.Align != "" {
		base.Align = override.Align
	}
	// Note: Show and ShowName booleans are always overridden when explicitly set
	// Since we can't distinguish between "not set" and "set to false" in TOML,
	// the user theme will override these values
	return base
}

func mergeStatusBar(base, override ThemeStatusBar) ThemeStatusBar {
	if override.Position != "" {
		base.Position = override.Position
	}
	if override.HelpFormat != "" {
		base.HelpFormat = override.HelpFormat
	}
	if override.KeySeparator != "" {
		base.KeySeparator = override.KeySeparator
	}
	return base
}

func mergeArticleList(base, override ThemeArticleList) ThemeArticleList {
	if override.NumberFormat != "" {
		base.NumberFormat = override.NumberFormat
	}
	if override.SeparatorChar != "" {
		base.SeparatorChar = override.SeparatorChar
	}
	if override.MetaFormat != "" {
		base.MetaFormat = override.MetaFormat
	}
	return base
}

func mergeTabBar(base, override ThemeTabBar) ThemeTabBar {
	if override.Position != "" {
		base.Position = override.Position
	}
	if override.UnreadFormat != "" {
		base.UnreadFormat = override.UnreadFormat
	}
	if override.MaxTabWidth != 0 {
		base.MaxTabWidth = override.MaxTabWidth
	}
	if override.TruncateSuffix != "" {
		base.TruncateSuffix = override.TruncateSuffix
	}
	return base
}
