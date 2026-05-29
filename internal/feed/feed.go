// Package feed handles fetching and parsing RSS and Atom feeds.
package feed

import (
	"context"
	"fmt"
	"html"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/adamdb/termnews/internal/config"
	"github.com/mmcdole/gofeed"
)

// Article represents a single news article.
type Article struct {
	Title       string
	Link        string
	Description string
	PubDate     time.Time
	Author      string
	Categories  []string // Categories/tags for the article
	Read        bool     // Whether the article has been read
	Bookmarked  bool     // Whether the article is bookmarked
	GUID        string   // Unique identifier for tracking read status
}

// Summary returns a plain text summary of the article description.
func (a *Article) Summary() string {
	return htmlToPlainText(a.Description)
}

// FormattedDate returns a human-readable date string.
func (a *Article) FormattedDate() string {
	if a.PubDate.IsZero() {
		return ""
	}
	now := time.Now()
	diff := now.Sub(a.PubDate)

	switch {
	case diff < time.Hour:
		mins := int(diff.Minutes())
		if mins <= 1 {
			return "just now"
		}
		return fmt.Sprintf("%dm ago", mins)
	case diff < 24*time.Hour:
		hours := int(diff.Hours())
		if hours == 1 {
			return "1h ago"
		}
		return fmt.Sprintf("%dh ago", hours)
	case diff < 7*24*time.Hour:
		days := int(diff.Hours() / 24)
		if days == 1 {
			return "1d ago"
		}
		return fmt.Sprintf("%dd ago", days)
	default:
		return a.PubDate.Format("Jan 2, 2006")
	}
}

// FetchError represents an error that occurred during feed fetching.
type FetchError struct {
	Type    string // "network" or "parse"
	Message string
}

func (e FetchError) Error() string {
	return fmt.Sprintf("%s error: %s", e.Type, e.Message)
}

// FetchResult holds the result of fetching a feed.
type FetchResult struct {
	SourceIndex int
	Articles    []Article
	Error       error
}

// Fetch fetches articles from a feed source.
func Fetch(ctx context.Context, source *config.Source) ([]Article, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, source.URL, nil)
	if err != nil {
		return nil, FetchError{Type: "network", Message: err.Error()}
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; termnews/1.0; +https://github.com/adamdb/termnews)")
	req.Header.Set("Accept", "application/rss+xml, application/atom+xml, application/xml, text/xml")

	resp, err := client.Do(req)
	if err != nil {
		return nil, FetchError{Type: "network", Message: err.Error()}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, FetchError{Type: "network", Message: fmt.Sprintf("HTTP %d: %s", resp.StatusCode, resp.Status)}
	}

	parser := gofeed.NewParser()
	feed, err := parser.Parse(resp.Body)
	if err != nil {
		return nil, FetchError{Type: "parse", Message: err.Error()}
	}

	articles := make([]Article, 0, len(feed.Items))
	for _, item := range feed.Items {
		article := Article{
			Title:       strings.TrimSpace(item.Title),
			Link:        strings.TrimSpace(item.Link),
			Description: item.Description,
			Categories:  item.Categories,
		}

		// Set publication date
		if item.PublishedParsed != nil {
			article.PubDate = *item.PublishedParsed
		} else if item.UpdatedParsed != nil {
			article.PubDate = *item.UpdatedParsed
		}

		// Set author
		if item.Author != nil {
			article.Author = strings.TrimSpace(item.Author.Name)
		} else if len(item.Authors) > 0 && item.Authors[0] != nil {
			article.Author = strings.TrimSpace(item.Authors[0].Name)
		}

		// Use content if description is empty
		if article.Description == "" && item.Content != "" {
			article.Description = item.Content
		}

		// Set GUID for tracking
		if item.GUID != "" {
			article.GUID = item.GUID
		} else {
			article.GUID = item.Link
		}

		articles = append(articles, article)
	}

	return articles, nil
}

// FetchAsync fetches a feed asynchronously and sends the result to a channel.
func FetchAsync(ctx context.Context, source *config.Source, index int, results chan<- FetchResult) {
	articles, err := Fetch(ctx, source)
	results <- FetchResult{
		SourceIndex: index,
		Articles:    articles,
		Error:       err,
	}
}

// htmlToPlainText converts HTML to plain text.
func htmlToPlainText(htmlStr string) string {
	// Unescape HTML entities
	text := html.UnescapeString(htmlStr)

	// Remove script and style elements
	reScript := regexp.MustCompile(`(?is)<script[^>]*>.*?</script>`)
	text = reScript.ReplaceAllString(text, "")

	reStyle := regexp.MustCompile(`(?is)<style[^>]*>.*?</style>`)
	text = reStyle.ReplaceAllString(text, "")

	// Replace common block elements with newlines
	blockTags := []string{"br", "p", "div", "li", "tr", "h1", "h2", "h3", "h4", "h5", "h6"}
	for _, tag := range blockTags {
		reTag := regexp.MustCompile(`(?i)<` + tag + `[^>]*>`)
		text = reTag.ReplaceAllString(text, "\n")
		reTagClose := regexp.MustCompile(`(?i)</` + tag + `>`)
		text = reTagClose.ReplaceAllString(text, "\n")
	}

	// Remove all other HTML tags
	reTags := regexp.MustCompile(`<[^>]+>`)
	text = reTags.ReplaceAllString(text, "")

	// Clean up whitespace
	lines := strings.Split(text, "\n")
	var cleanLines []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			cleanLines = append(cleanLines, line)
		}
	}

	return strings.Join(cleanLines, "\n")
}
