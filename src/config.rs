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

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Source {
    pub name: String,
    pub url: String,
    pub feed_type: FeedType,
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
            url,
            feed_type,
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
                    url: "https://feeds.reuters.com/reuters/topNews".to_string(),
                    feed_type: FeedType::Rss,
                },
                Source {
                    name: "BBC News".to_string(),
                    url: "http://feeds.bbci.co.uk/news/rss.xml".to_string(),
                    feed_type: FeedType::Rss,
                },
                Source {
                    name: "Hacker News".to_string(),
                    url: "https://news.ycombinator.com/rss".to_string(),
                    feed_type: FeedType::Rss,
                },
                Source {
                    name: "NPR News".to_string(),
                    url: "https://feeds.npr.org/1001/rss.xml".to_string(),
                    feed_type: FeedType::Rss,
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
        assert_eq!(deserialized.sources[0].feed_type, FeedType::Atom);
    }
}
