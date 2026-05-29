# termnews

A terminal news aggregator written in Go. Browse RSS and Atom feeds in a beautiful TUI with tabbed windows.

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
- 🎨 **Relative timestamps** — human-readable "2h ago" style timestamps

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

Example `config.toml`:

```toml
refresh_interval_mins = 5
max_articles_per_feed = 50
theme = "default"
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
| `theme` | Color theme (default/dark/light) | `default` |
| `show_read_status` | Track read/unread status | `true` |

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

