# termnews

A terminal news aggregator written in Rust. Browse RSS and Atom feeds in a beautiful TUI with tabbed windows — one per news source.

## Features

- 📰 **Multi-source news aggregation** — fetch multiple RSS/Atom feeds concurrently
- 🗂️ **Tabbed windows** — each tab shows a single news source
- ➕ **Add/remove sources** — manage news sources at runtime and persist them
- 📖 **Article detail view** — read the full article summary inline
- ⌨️ **Keyboard-driven** — fully navigable without a mouse

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

Example `config.toml`:

```toml
[[sources]]
name = "Reuters Top News"
url = "https://feeds.reuters.com/reuters/topNews"
feed_type = "rss"

[[sources]]
name = "My Atom Feed"
url = "https://example.com/feed.atom"
feed_type = "atom"
```

## Installation

```bash
cargo install --path .
```

## Building from source

```bash
cargo build --release
./target/release/termnews
```

