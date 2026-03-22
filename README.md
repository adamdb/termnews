# termnews

A terminal news aggregator written in Rust. Browse RSS and Atom feeds in a beautiful TUI with tabbed windows — one per news source. Future support for IRC channels is planned.

## Features

- 📰 **Multi-source news aggregation** — fetch multiple RSS/Atom feeds concurrently
- 🗂️ **Tabbed windows** — each tab shows a single news source with type indicator ([RSS], [Atom], [IRC])
- ➕ **Add/remove sources** — manage news sources at runtime and persist them
- 📖 **Article detail view** — read the full article summary inline
- ⌨️ **Keyboard-driven** — fully navigable without a mouse
- 🔌 **Extensible architecture** — ready for IRC channel support (coming soon)

## Keybindings

| Key | Action |
|-----|--------|
| `Tab` / `→` | Next source tab |
| `Shift+Tab` / `←` | Previous source tab |
| `↓` / `j` | Move selection down |
| `↑` / `k` | Move selection up |
| `Enter` | Open selected article |
| `Esc` / `Backspace` | Back to article list |
| `r` | Refresh current feed |
| `a` | Add a new news source |
| `d` | Delete the current news source |
| `q` | Quit |

## Configuration

On first run, a configuration file is created at:

- **Linux/macOS**: `~/.config/termnews/config.toml`
- **Windows**: `%APPDATA%\termnews\config.toml`

### RSS/Atom Feed Configuration

Example `config.toml` for RSS/Atom feeds:

```toml
[[sources]]
name = "Reuters Top News"
source_type = "feed"
url = "http://feeds.reuters.com/reuters/topNews"
feed_type = "rss"

[[sources]]
name = "My Atom Feed"
source_type = "feed"
url = "https://example.com/feed.atom"
feed_type = "atom"
```

### IRC Configuration (Placeholder - Not Yet Implemented)

Future IRC channel support will use this configuration format:

```toml
[[sources]]
name = "Libera #rust"
source_type = "irc"
server = "irc.libera.chat"
port = 6697
channel = "#rust"
nick = "termnews_user"
use_ssl = true
```

**Note:** IRC functionality is currently a placeholder and not yet implemented. The UI will display an error if you try to add an IRC source.

## Installation

```bash
cargo install --path .
```

## Updating

If you've already installed termnews and want to pull the latest changes:

```bash
# Pull the latest changes from the repository
git pull

# Reinstall with the updated code
cargo install --path .
```

## Building from source

```bash
cargo build --release
./target/release/termnews
```

