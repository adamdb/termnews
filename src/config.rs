use serde::{Deserialize, Serialize};
use std::fs;
use std::io::Write;
use std::path::PathBuf;

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
#[serde(rename_all = "lowercase")]
pub enum FeedType {
    Rss,
    Atom,
}

impl std::fmt::Display for FeedType {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            FeedType::Rss => write!(f, "rss"),
            FeedType::Atom => write!(f, "atom"),
        }
    }
}

/// Represents the type of source: either a Feed (RSS/Atom) or IRC channel
#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
#[serde(rename_all = "lowercase", tag = "source_type")]
pub enum SourceType {
    Feed {
        url: String,
        feed_type: FeedType,
    },
    Irc {
        server: String,
        port: u16,
        channel: String,
        #[serde(default = "default_nick")]
        nick: String,
        #[serde(default)]
        use_ssl: bool,
    },
}

fn default_nick() -> String {
    "termnews_user".to_string()
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Source {
    pub name: String,
    #[serde(flatten)]
    pub source_type: SourceType,
}

impl Source {
    /// Returns true if this is a feed source (RSS or Atom)
    pub fn is_feed(&self) -> bool {
        matches!(self.source_type, SourceType::Feed { .. })
    }

    /// Returns true if this is an IRC source
    pub fn is_irc(&self) -> bool {
        matches!(self.source_type, SourceType::Irc { .. })
    }

    /// Legacy accessor for feed URL (for backward compatibility)
    #[allow(dead_code)]
    pub fn url(&self) -> Option<&str> {
        match &self.source_type {
            SourceType::Feed { url, .. } => Some(url),
            _ => None,
        }
    }

    /// Legacy accessor for feed type (for backward compatibility)
    #[allow(dead_code)]
    pub fn feed_type(&self) -> Option<&FeedType> {
        match &self.source_type {
            SourceType::Feed { feed_type, .. } => Some(feed_type),
            _ => None,
        }
    }
}

#[derive(Debug, Serialize, Deserialize, Default)]
pub struct Config {
    pub sources: Vec<Source>,
}

impl Config {
    pub fn config_path() -> PathBuf {
        let mut path = dirs::config_dir().unwrap_or_else(|| PathBuf::from("."));
        path.push("termnews");
        path.push("config.toml");
        path
    }

    pub fn load() -> Self {
        let path = Self::config_path();
        if path.exists() {
            let content = fs::read_to_string(&path).unwrap_or_default();
            toml::from_str(&content).unwrap_or_default()
        } else {
            Self::default_config()
        }
    }

    pub fn save(&self) -> std::io::Result<()> {
        let path = Self::config_path();
        if let Some(parent) = path.parent() {
            fs::create_dir_all(parent)?;
        }
        let content = toml::to_string_pretty(self).unwrap_or_default();
        let mut file = fs::File::create(&path)?;
        file.write_all(content.as_bytes())?;
        Ok(())
    }

    pub fn add_source(&mut self, name: String, url: String, feed_type: FeedType) {
        self.sources.push(Source {
            name,
            source_type: SourceType::Feed { url, feed_type },
        });
    }

    /// Add an IRC source to the configuration
    ///
    /// # Future Use
    /// This method is reserved for future IRC functionality.
    /// IRC sources will be displayed with [IRC] prefix in the UI.
    #[allow(dead_code)]
    pub fn add_irc_source(
        &mut self,
        name: String,
        server: String,
        port: u16,
        channel: String,
        nick: Option<String>,
        use_ssl: bool,
    ) {
        self.sources.push(Source {
            name,
            source_type: SourceType::Irc {
                server,
                port,
                channel,
                nick: nick.unwrap_or_else(default_nick),
                use_ssl,
            },
        });
    }

    pub fn remove_source(&mut self, index: usize) {
        if index < self.sources.len() {
            self.sources.remove(index);
        }
    }

    fn default_config() -> Self {
        Config {
            sources: vec![
                Source {
                    name: "Reuters Top News".to_string(),
                    source_type: SourceType::Feed {
                        url: "http://feeds.reuters.com/reuters/topNews".to_string(),
                        feed_type: FeedType::Rss,
                    },
                },
                Source {
                    name: "BBC News".to_string(),
                    source_type: SourceType::Feed {
                        url: "http://feeds.bbci.co.uk/news/rss.xml".to_string(),
                        feed_type: FeedType::Rss,
                    },
                },
                Source {
                    name: "Hacker News".to_string(),
                    source_type: SourceType::Feed {
                        url: "https://news.ycombinator.com/rss".to_string(),
                        feed_type: FeedType::Rss,
                    },
                },
                Source {
                    name: "NPR News".to_string(),
                    source_type: SourceType::Feed {
                        url: "https://feeds.npr.org/1001/rss.xml".to_string(),
                        feed_type: FeedType::Rss,
                    },
                },
            ],
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_default_config_has_sources() {
        let config = Config::default_config();
        assert!(!config.sources.is_empty());
    }

    #[test]
    fn test_add_source() {
        let mut config = Config::default();
        config.add_source(
            "Test Feed".to_string(),
            "https://example.com/rss".to_string(),
            FeedType::Rss,
        );
        assert_eq!(config.sources.len(), 1);
        assert_eq!(config.sources[0].name, "Test Feed");
    }

    #[test]
    fn test_remove_source() {
        let mut config = Config::default_config();
        let initial_len = config.sources.len();
        config.remove_source(0);
        assert_eq!(config.sources.len(), initial_len - 1);
    }

    #[test]
    fn test_remove_source_out_of_bounds() {
        let mut config = Config::default();
        // Should not panic on empty config
        config.remove_source(0);
        assert!(config.sources.is_empty());
    }

    #[test]
    fn test_feed_type_display() {
        assert_eq!(FeedType::Rss.to_string(), "rss");
        assert_eq!(FeedType::Atom.to_string(), "atom");
    }

    #[test]
    fn test_config_roundtrip() {
        let mut config = Config::default();
        config.add_source(
            "News Source".to_string(),
            "https://example.com/feed".to_string(),
            FeedType::Atom,
        );
        let serialized = toml::to_string_pretty(&config).unwrap();
        let deserialized: Config = toml::from_str(&serialized).unwrap();
        assert_eq!(deserialized.sources.len(), 1);
        assert_eq!(
            deserialized.sources[0].feed_type(),
            Some(&FeedType::Atom)
        );
    }

    #[test]
    fn test_source_type_feed() {
        let source = Source {
            name: "Test Feed".to_string(),
            source_type: SourceType::Feed {
                url: "https://example.com/rss".to_string(),
                feed_type: FeedType::Rss,
            },
        };
        assert!(source.is_feed());
        assert!(!source.is_irc());
        assert_eq!(source.url(), Some("https://example.com/rss"));
        assert_eq!(source.feed_type(), Some(&FeedType::Rss));
    }

    #[test]
    fn test_source_type_irc() {
        let source = Source {
            name: "Test IRC".to_string(),
            source_type: SourceType::Irc {
                server: "irc.libera.chat".to_string(),
                port: 6667,
                channel: "#rust".to_string(),
                nick: "testbot".to_string(),
                use_ssl: false,
            },
        };
        assert!(source.is_irc());
        assert!(!source.is_feed());
        assert_eq!(source.url(), None);
        assert_eq!(source.feed_type(), None);
    }

    #[test]
    fn test_add_irc_source() {
        let mut config = Config::default();
        config.add_irc_source(
            "IRC Test".to_string(),
            "irc.libera.chat".to_string(),
            6667,
            "#test".to_string(),
            None,
            false,
        );
        assert_eq!(config.sources.len(), 1);
        assert!(config.sources[0].is_irc());
    }

    #[test]
    fn test_toml_serialization_format() {
        let mut config = Config::default();
        config.sources.clear();

        // Add an RSS feed
        config.add_source(
            "Reuters".to_string(),
            "http://feeds.reuters.com/reuters/topNews".to_string(),
            FeedType::Rss,
        );

        // Add an IRC source
        config.add_irc_source(
            "Libera #rust".to_string(),
            "irc.libera.chat".to_string(),
            6697,
            "#rust".to_string(),
            Some("testuser".to_string()),
            true,
        );

        let toml_str = toml::to_string_pretty(&config).unwrap();

        // Verify both source types are present in serialized format
        assert!(toml_str.contains("source_type = \"feed\""));
        assert!(toml_str.contains("source_type = \"irc\""));
        assert!(toml_str.contains("feed_type = \"rss\""));
        assert!(toml_str.contains("server = \"irc.libera.chat\""));
        assert!(toml_str.contains("channel = \"#rust\""));

        // Print for manual verification
        println!("Serialized config:\n{}", toml_str);
    }
}
