package feed

import (
	"testing"
	"time"
)

func TestArticleFormattedDate(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		pubDate  time.Time
		expected string
	}{
		{
			name:     "just now",
			pubDate:  now.Add(-30 * time.Second),
			expected: "just now",
		},
		{
			name:     "minutes ago",
			pubDate:  now.Add(-15 * time.Minute),
			expected: "15m ago",
		},
		{
			name:     "1 hour ago",
			pubDate:  now.Add(-1 * time.Hour),
			expected: "1h ago",
		},
		{
			name:     "hours ago",
			pubDate:  now.Add(-5 * time.Hour),
			expected: "5h ago",
		},
		{
			name:     "1 day ago",
			pubDate:  now.Add(-24 * time.Hour),
			expected: "1d ago",
		},
		{
			name:     "days ago",
			pubDate:  now.Add(-3 * 24 * time.Hour),
			expected: "3d ago",
		},
		{
			name:     "zero time",
			pubDate:  time.Time{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			article := Article{PubDate: tt.pubDate}
			got := article.FormattedDate()
			if got != tt.expected {
				t.Errorf("FormattedDate() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestHtmlToPlainText(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected string
	}{
		{
			name:     "simple text",
			html:     "Hello World",
			expected: "Hello World",
		},
		{
			name:     "paragraph tags",
			html:     "<p>First paragraph</p><p>Second paragraph</p>",
			expected: "First paragraph\nSecond paragraph",
		},
		{
			name:     "html entities ampersand",
			html:     "&amp;",
			expected: "&",
		},
		{
			name:     "script removal",
			html:     "Before<script>alert('test')</script>After",
			expected: "BeforeAfter",
		},
		{
			name:     "nested tags",
			html:     "<div><p><strong>Bold</strong> text</p></div>",
			expected: "Bold text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := htmlToPlainText(tt.html)
			if got != tt.expected {
				t.Errorf("htmlToPlainText() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestArticleSummary(t *testing.T) {
	article := Article{
		Description: "<p>This is a <strong>test</strong> article.</p>",
	}
	summary := article.Summary()
	if summary != "This is a test article." {
		t.Errorf("Summary() = %q, want %q", summary, "This is a test article.")
	}
}

func TestFetchErrorString(t *testing.T) {
	err := FetchError{Type: "network", Message: "connection timeout"}
	expected := "network error: connection timeout"
	if err.Error() != expected {
		t.Errorf("Error() = %q, want %q", err.Error(), expected)
	}
}
