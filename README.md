# termnews

A terminal news aggregator written in Go. Browse RSS and Atom feeds in a beautiful TUI with tabbed windows and highly customizable themes.

## Features

- 📰 **Multi-source news aggregation** — fetch multiple RSS/Atom feeds concurrently
- 🗂️ **Tabbed windows** — each tab shows a single news source
- ➕ **Add/remove sources** — manage news sources at runtime and persist them
- 📖 **Article detail view** — read the full article summary inline
- ⌨️ **Keyboard-driven** — fully navigable without a mouse
- 🔍 **Search** — search across all loaded articles
- 📚 **Read tracking** — track which articles you've already read
- ⭐ **Bookmarks** — bookmark important articles for later
- 🏷️ **Categories** — organize feeds by category
- 🔄 **Auto-refresh** — configurable automatic feed refresh
- 📊 **Summary view** — see top articles from all feeds at once
- 🌐 **Open in browser** — quickly open articles in your default browser
- 🎨 **Highly customizable themes** — BitchX-inspired ANSI styling with built-in presets
- ✨ **Custom colors and markers** — configure every visual element

## Themes

termnews includes several built-in themes inspired by classic terminal applications:

| Theme | Description |
|-------|-------------|
| `default` | Modern, clean appearance with rounded borders |
| `bitchx` | Classic IRC client style with ASCII art banners and double borders |
| `hacker` | Matrix-style green-on-black with ASCII art |
| `minimal` | Clean, distraction-free interface with no borders |
| `retro` | BBS-style with bold colors and heavy borders |

### Setting a Theme

In your `config.toml`:

```toml
theme = "bitchx"
```

### Custom Theme Configuration

You can customize any visual element by adding a `[custom_theme]` section:

```toml
theme = "default"  # Start with a base theme

[custom_theme.colors]
primary = "14"       # Bright cyan
secondary = "10"     # Bright green
text = "7"           # Silver
text_muted = "8"     # Dark gray
error = "9"          # Bright red
success = "10"       # Bright green
link = "12"          # Bright blue
unread = "10"        # Bright green
bookmark = "11"      # Bright yellow
category = "13"      # Bright magenta
border = "14"        # Bright cyan
header = "11"        # Bright yellow
status_bg = "4"      # Blue background
status_fg = "15"     # White text

[custom_theme.borders]
style = "double"     # Options: rounded, sharp, double, heavy, ascii, none, custom
separator = "═"

[custom_theme.markers]
unread = "[+]"
bookmark = "[*]"
loading = "[~]"
error = "[!]"
selected = "»»"
tab_separator = "│"
tab_active_left = "«"
tab_active_right = "»"

[custom_theme.header]
show = true
show_name = true
title = "TERMNEWS"
align = "center"     # Options: left, center, right
banner = """
╔══════════════════════════════════════════════════════╗
║  ▀█▀ █▀▀ █▀█ █▀▄▀█ █▄░█ █▀▀ █░█░█ █▀                 ║
║  ░█░ ██▄ █▀▄ █░▀░█ █░▀█ ██▄ ▀▄▀▄▀ ▄█  v1.0          ║
╚══════════════════════════════════════════════════════╝
"""

[custom_theme.status_bar]
show = true
position = "bottom"
show_keys = true
key_separator = " │ "

[custom_theme.article_list]
show_numbers = true
number_format = "%02d»"
show_separators = true
separator_char = "─"
show_meta = true
compact = false

[custom_theme.tab_bar]
position = "top"
show_unread_count = true
unread_format = "[%d]"
max_tab_width = 15
truncate_suffix = ".."
```

## Keybindings

### List View
| Key | Action |
|-----|--------|
| `Tab` / `→` | Next source tab |
| `Shift+Tab` / `←` | Previous source tab |
| `↓` / `j` | Move selection down |
| `↑` / `k` | Move selection up |
| `Enter` | Open selected article |
| `r` | Refresh current feed |
| `R` | Refresh all feeds |
| `t` | Toggle auto-refresh |
| `s` | Summary view |
| `/` | Search articles |
| `a` | Add a new news source |
| `d` | Delete the current news source |
| `b` | Toggle bookmark |
| `m` | Mark all as read |
| `o` | Open in browser |
| `?` | Show help |
| `q` | Quit |

### Detail View
| Key | Action |
|-----|--------|
| `↓` / `j` | Scroll down |
| `↑` / `k` | Scroll up |
| `b` | Toggle bookmark |
| `o` | Open in browser |
| `Esc` / `Backspace` | Back to article list |
| `q` | Quit |

## Configuration

On first run, a configuration file is created at:

- **Linux/macOS**: `~/.config/termnews/config.toml`
- **Windows**: `%APPDATA%\termnews\config.toml`

### Example `config.toml`

```toml
refresh_interval_mins = 5
max_articles_per_feed = 50
theme = "bitchx"
show_read_status = true

[[sources]]
name = "Reuters Top News"
url = "http://feeds.reuters.com/reuters/topNews"
feed_type = "rss"
category = "World News"

[[sources]]
name = "Hacker News"
url = "https://news.ycombinator.com/rss"
feed_type = "rss"
category = "Tech"

[[sources]]
name = "My Atom Feed"
url = "https://example.com/feed.atom"
feed_type = "atom"
```

### Configuration Options

| Option | Description | Default |
|--------|-------------|---------|
| `refresh_interval_mins` | Auto-refresh interval in minutes | `5` |
| `max_articles_per_feed` | Maximum articles shown per feed | `50` |
| `theme` | Built-in theme name | `default` |
| `show_read_status` | Track read/unread status | `true` |

### Available Themes

| Theme | Style |
|-------|-------|
| `default` | Modern with rounded borders and subtle colors |
| `bitchx` | Classic IRC client with ASCII art and double borders |
| `hacker` | Matrix-style green terminal aesthetic |
| `minimal` | Clean, borderless, distraction-free |
| `retro` | BBS-style with bold colors and heavy borders |

### Feed Types

- `rss` — RSS feed
- `atom` — Atom feed  
- `auto` — Auto-detect feed type (recommended)

## Installation

### Using Go

```bash
go install github.com/adamdb/termnews@latest
```

### Building from source

```bash
git clone https://github.com/adamdb/termnews.git
cd termnews
go build -o termnews .
./termnews
```

### Building a release binary

```bash
go build -ldflags="-s -w" -o termnews .
```

## Requirements

- Go 1.21 or later

## Updating

If you've already installed termnews and want to pull the latest changes:

```bash
# Pull the latest changes from the repository
git pull

# Rebuild
go build -o termnews .
```

## Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) — TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) — Styling
- [gofeed](https://github.com/mmcdole/gofeed) — RSS/Atom parsing
- [BurntSushi/toml](https://github.com/BurntSushi/toml) — TOML configuration

