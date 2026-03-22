use crate::config::{FeedType, Source, SourceType};

#[derive(Debug, Clone)]
pub struct Article {
    pub title: String,
    pub link: String,
    pub description: String,
    pub pub_date: String,
    pub author: String,
}

impl Article {
    pub fn summary(&self) -> String {
        let plain = html2text::from_read(self.description.as_bytes(), 80);
        // Trim excessive whitespace/newlines
        plain
            .lines()
            .map(|l| l.trim())
            .filter(|l| !l.is_empty())
            .collect::<Vec<_>>()
            .join("\n")
    }
}

#[derive(Debug, Clone)]
pub enum FetchError {
    Network(String),
    Parse(String),
}

impl std::fmt::Display for FetchError {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            FetchError::Network(e) => write!(f, "Network error: {}", e),
            FetchError::Parse(e) => write!(f, "Parse error: {}", e),
        }
    }
}

pub fn fetch_feed(source: &Source) -> Result<Vec<Article>, FetchError> {
    match &source.source_type {
        SourceType::Feed { url, feed_type } => fetch_rss_atom(url, feed_type),
        SourceType::Irc { .. } => {
            // IRC sources are not yet supported for fetching
            Err(FetchError::Network(
                "IRC source fetching is not yet implemented".to_string(),
            ))
        }
    }
}

fn fetch_rss_atom(url: &str, feed_type: &FeedType) -> Result<Vec<Article>, FetchError> {
    // Build a client with proper configuration for feed fetching
    let client = reqwest::blocking::Client::builder()
        .user_agent("Mozilla/5.0 (compatible; termnews/0.1; +https://github.com/adamdb/newsterm)")
        .timeout(std::time::Duration::from_secs(30))
        .redirect(reqwest::redirect::Policy::limited(10))
        .build()
        .map_err(|e| FetchError::Network(format!("Failed to build HTTP client: {}", e)))?;

    let response = client
        .get(url)
        .send()
        .map_err(|e| FetchError::Network(e.to_string()))?;

    let body = response
        .text()
        .map_err(|e| FetchError::Network(e.to_string()))?;

    match feed_type {
        FeedType::Rss => parse_rss(&body),
        FeedType::Atom => parse_atom(&body),
    }
}

fn parse_rss(body: &str) -> Result<Vec<Article>, FetchError> {
    let channel = body
        .parse::<rss::Channel>()
        .map_err(|e| FetchError::Parse(e.to_string()))?;

    let articles = channel
        .items()
        .iter()
        .map(|item| Article {
            title: item.title().unwrap_or("(no title)").to_string(),
            link: item.link().unwrap_or("").to_string(),
            description: item.description().unwrap_or("").to_string(),
            pub_date: item.pub_date().unwrap_or("").to_string(),
            author: item
                .author()
                .or_else(|| {
                    item.dublin_core_ext()
                        .and_then(|dc| dc.creators().first().map(|s| s.as_str()))
                })
                .unwrap_or("")
                .to_string(),
        })
        .collect();

    Ok(articles)
}

fn parse_atom(body: &str) -> Result<Vec<Article>, FetchError> {
    let feed: atom_syndication::Feed = body
        .parse()
        .map_err(|e: atom_syndication::Error| FetchError::Parse(e.to_string()))?;

    let articles = feed
        .entries()
        .iter()
        .map(|entry| {
            let description = entry
                .content()
                .and_then(|c| c.value())
                .map(|s| s.to_string())
                .or_else(|| entry.summary().map(|s| s.value.clone()))
                .unwrap_or_default();

            let link = entry
                .links()
                .first()
                .map(|l| l.href())
                .unwrap_or("")
                .to_string();

            let pub_date = entry
                .published()
                .map(|d| d.to_rfc3339())
                .unwrap_or_else(|| entry.updated().to_rfc3339());

            let author = entry
                .authors()
                .first()
                .map(|a| a.name())
                .unwrap_or("")
                .to_string();

            Article {
                title: entry.title().value.to_string(),
                link,
                description,
                pub_date,
                author,
            }
        })
        .collect();

    Ok(articles)
}
