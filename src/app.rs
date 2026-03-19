use crate::config::{Config, FeedType, Source};
use crate::feed::{Article, FetchError};

#[derive(Debug, Clone, PartialEq)]
pub enum ViewMode {
    List,
    Detail,
    AddSource,
}

#[derive(Debug, Clone, PartialEq)]
pub enum InputField {
    Name,
    Url,
    FeedType,
}

#[derive(Debug)]
pub struct AddSourceState {
    pub active_field: InputField,
    pub name: String,
    pub url: String,
    pub feed_type_index: usize, // 0 = RSS, 1 = Atom
    pub error: Option<String>,
}

impl AddSourceState {
    pub fn new() -> Self {
        Self {
            active_field: InputField::Name,
            name: String::new(),
            url: String::new(),
            feed_type_index: 0,
            error: None,
        }
    }

    pub fn selected_feed_type(&self) -> FeedType {
        if self.feed_type_index == 0 {
            FeedType::Rss
        } else {
            FeedType::Atom
        }
    }
}

pub struct TabState {
    pub source: Source,
    pub articles: Vec<Article>,
    pub scroll_offset: usize,
    pub selected: Option<usize>,
    pub loading: bool,
    pub error: Option<FetchError>,
}

impl TabState {
    pub fn new(source: Source) -> Self {
        Self {
            source,
            articles: Vec::new(),
            scroll_offset: 0,
            selected: None,
            loading: true,
            error: None,
        }
    }
}

pub struct App {
    pub config: Config,
    pub tabs: Vec<TabState>,
    pub current_tab: usize,
    pub view_mode: ViewMode,
    pub add_source_state: Option<AddSourceState>,
    pub should_quit: bool,
    pub status_message: Option<String>,
}

impl App {
    pub fn new(config: Config) -> Self {
        let tabs = config
            .sources
            .iter()
            .map(|s| TabState::new(s.clone()))
            .collect();

        Self {
            config,
            tabs,
            current_tab: 0,
            view_mode: ViewMode::List,
            add_source_state: None,
            should_quit: false,
            status_message: None,
        }
    }

    pub fn next_tab(&mut self) {
        if !self.tabs.is_empty() {
            self.current_tab = (self.current_tab + 1) % self.tabs.len();
            self.view_mode = ViewMode::List;
        }
    }

    pub fn prev_tab(&mut self) {
        if !self.tabs.is_empty() {
            self.current_tab =
                (self.current_tab + self.tabs.len().saturating_sub(1)) % self.tabs.len();
            self.view_mode = ViewMode::List;
        }
    }

    pub fn scroll_down(&mut self) {
        if let Some(tab) = self.tabs.get_mut(self.current_tab) {
            match self.view_mode {
                ViewMode::List => {
                    if !tab.articles.is_empty() {
                        let max = tab.articles.len().saturating_sub(1);
                        let sel = tab.selected.unwrap_or(0);
                        let new_sel = (sel + 1).min(max);
                        tab.selected = Some(new_sel);
                        // Adjust scroll offset for visibility
                        tab.scroll_offset = tab.scroll_offset.max(new_sel.saturating_sub(20));
                    }
                }
                ViewMode::Detail => {
                    tab.scroll_offset = tab.scroll_offset.saturating_add(1);
                }
                _ => {}
            }
        }
    }

    pub fn scroll_up(&mut self) {
        if let Some(tab) = self.tabs.get_mut(self.current_tab) {
            match self.view_mode {
                ViewMode::List => {
                    if !tab.articles.is_empty() {
                        let sel = tab.selected.unwrap_or(0);
                        let new_sel = sel.saturating_sub(1);
                        tab.selected = Some(new_sel);
                        if new_sel < tab.scroll_offset {
                            tab.scroll_offset = new_sel;
                        }
                    }
                }
                ViewMode::Detail => {
                    tab.scroll_offset = tab.scroll_offset.saturating_sub(1);
                }
                _ => {}
            }
        }
    }

    pub fn select_article(&mut self) {
        if let Some(tab) = self.tabs.get(self.current_tab) {
            if tab.selected.is_some() && !tab.articles.is_empty() {
                // Reset scroll for detail view
                let current = self.current_tab;
                if let Some(t) = self.tabs.get_mut(current) {
                    t.scroll_offset = 0;
                }
                self.view_mode = ViewMode::Detail;
            }
        }
    }

    pub fn back_to_list(&mut self) {
        if self.view_mode == ViewMode::Detail {
            self.view_mode = ViewMode::List;
            if let Some(tab) = self.tabs.get_mut(self.current_tab) {
                tab.scroll_offset = 0;
            }
        }
    }

    pub fn start_add_source(&mut self) {
        self.add_source_state = Some(AddSourceState::new());
        self.view_mode = ViewMode::AddSource;
    }

    pub fn cancel_add_source(&mut self) {
        self.add_source_state = None;
        self.view_mode = ViewMode::List;
    }

    pub fn confirm_add_source(&mut self) -> bool {
        if let Some(state) = &self.add_source_state {
            let name = state.name.trim().to_string();
            let url = state.url.trim().to_string();
            if name.is_empty() || url.is_empty() {
                if let Some(s) = &mut self.add_source_state {
                    s.error = Some("Name and URL cannot be empty".to_string());
                }
                return false;
            }
            let feed_type = state.selected_feed_type();
            let source = Source {
                name: name.clone(),
                url: url.clone(),
                feed_type,
            };
            self.config.add_source(name, url, source.feed_type.clone());
            let _ = self.config.save();
            self.tabs.push(TabState::new(source));
            let new_idx = self.tabs.len() - 1;
            self.tabs[new_idx].loading = true;
            self.add_source_state = None;
            self.view_mode = ViewMode::List;
            self.current_tab = new_idx;
            return true;
        }
        false
    }

    pub fn remove_current_source(&mut self) {
        if self.tabs.is_empty() {
            return;
        }
        let idx = self.current_tab;
        self.config.remove_source(idx);
        let _ = self.config.save();
        self.tabs.remove(idx);
        if self.current_tab >= self.tabs.len() && !self.tabs.is_empty() {
            self.current_tab = self.tabs.len() - 1;
        }
        self.view_mode = ViewMode::List;
    }

    pub fn set_articles(
        &mut self,
        tab_index: usize,
        result: Result<Vec<Article>, FetchError>,
    ) {
        if let Some(tab) = self.tabs.get_mut(tab_index) {
            tab.loading = false;
            match result {
                Ok(articles) => {
                    tab.articles = articles;
                    tab.error = None;
                    if !tab.articles.is_empty() {
                        tab.selected = Some(0);
                    }
                }
                Err(e) => {
                    tab.error = Some(e);
                }
            }
        }
    }

    pub fn refresh_current(&mut self) {
        if let Some(tab) = self.tabs.get_mut(self.current_tab) {
            tab.loading = true;
            tab.articles.clear();
            tab.selected = None;
            tab.scroll_offset = 0;
            tab.error = None;
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::config::{Config, FeedType};
    use crate::feed::Article;

    fn make_article(title: &str) -> Article {
        Article {
            title: title.to_string(),
            link: "https://example.com".to_string(),
            description: "<p>Content</p>".to_string(),
            pub_date: "2024-01-01".to_string(),
            author: "Author".to_string(),
        }
    }

    fn make_app_with_tabs(count: usize) -> App {
        let mut config = Config::default();
        for i in 0..count {
            config.add_source(
                format!("Source {}", i),
                "https://example.com/rss".to_string(),
                FeedType::Rss,
            );
        }
        App::new(config)
    }

    #[test]
    fn test_next_tab_wraps() {
        let mut app = make_app_with_tabs(3);
        assert_eq!(app.current_tab, 0);
        app.next_tab();
        assert_eq!(app.current_tab, 1);
        app.next_tab();
        assert_eq!(app.current_tab, 2);
        app.next_tab();
        assert_eq!(app.current_tab, 0); // wraps
    }

    #[test]
    fn test_prev_tab_wraps() {
        let mut app = make_app_with_tabs(3);
        assert_eq!(app.current_tab, 0);
        app.prev_tab();
        assert_eq!(app.current_tab, 2); // wraps to last
        app.prev_tab();
        assert_eq!(app.current_tab, 1);
    }

    #[test]
    fn test_scroll_down_selects_articles() {
        let mut app = make_app_with_tabs(1);
        let articles = vec![
            make_article("Article 1"),
            make_article("Article 2"),
            make_article("Article 3"),
        ];
        app.set_articles(0, Ok(articles));
        // Initial selection is 0
        assert_eq!(app.tabs[0].selected, Some(0));
        app.scroll_down();
        assert_eq!(app.tabs[0].selected, Some(1));
        app.scroll_down();
        assert_eq!(app.tabs[0].selected, Some(2));
        // Should not go past last
        app.scroll_down();
        assert_eq!(app.tabs[0].selected, Some(2));
    }

    #[test]
    fn test_scroll_up_does_not_underflow() {
        let mut app = make_app_with_tabs(1);
        let articles = vec![make_article("Article 1"), make_article("Article 2")];
        app.set_articles(0, Ok(articles));
        assert_eq!(app.tabs[0].selected, Some(0));
        app.scroll_up();
        assert_eq!(app.tabs[0].selected, Some(0)); // stays at 0
    }

    #[test]
    fn test_select_article_enters_detail_mode() {
        let mut app = make_app_with_tabs(1);
        let articles = vec![make_article("Article 1")];
        app.set_articles(0, Ok(articles));
        assert_eq!(app.view_mode, ViewMode::List);
        app.select_article();
        assert_eq!(app.view_mode, ViewMode::Detail);
    }

    #[test]
    fn test_back_to_list() {
        let mut app = make_app_with_tabs(1);
        let articles = vec![make_article("A")];
        app.set_articles(0, Ok(articles));
        app.select_article();
        assert_eq!(app.view_mode, ViewMode::Detail);
        app.back_to_list();
        assert_eq!(app.view_mode, ViewMode::List);
    }

    #[test]
    fn test_add_source_flow() {
        let mut app = make_app_with_tabs(1);
        app.start_add_source();
        assert_eq!(app.view_mode, ViewMode::AddSource);
        if let Some(state) = &mut app.add_source_state {
            state.name = "New Feed".to_string();
            state.url = "https://new.example.com/rss".to_string();
        }
        let result = app.confirm_add_source();
        assert!(result);
        assert_eq!(app.tabs.len(), 2);
        assert_eq!(app.current_tab, 1);
        assert_eq!(app.view_mode, ViewMode::List);
    }

    #[test]
    fn test_confirm_add_source_fails_on_empty() {
        let mut app = make_app_with_tabs(1);
        app.start_add_source();
        let result = app.confirm_add_source();
        assert!(!result);
        assert!(app.add_source_state.as_ref().unwrap().error.is_some());
    }

    #[test]
    fn test_remove_current_source() {
        let mut app = make_app_with_tabs(3);
        assert_eq!(app.tabs.len(), 3);
        app.remove_current_source();
        assert_eq!(app.tabs.len(), 2);
    }

    #[test]
    fn test_refresh_current() {
        let mut app = make_app_with_tabs(1);
        let articles = vec![make_article("A"), make_article("B")];
        app.set_articles(0, Ok(articles));
        assert_eq!(app.tabs[0].articles.len(), 2);
        app.refresh_current();
        assert!(app.tabs[0].loading);
        assert!(app.tabs[0].articles.is_empty());
    }

    #[test]
    fn test_set_articles_error() {
        let mut app = make_app_with_tabs(1);
        use crate::feed::FetchError;
        app.set_articles(0, Err(FetchError::Network("timeout".to_string())));
        assert!(!app.tabs[0].loading);
        assert!(app.tabs[0].error.is_some());
    }

    #[test]
    fn test_add_source_state_feed_type_toggle() {
        let mut state = AddSourceState::new();
        assert_eq!(state.feed_type_index, 0);
        assert_eq!(state.selected_feed_type(), FeedType::Rss);
        state.feed_type_index = 1;
        assert_eq!(state.selected_feed_type(), FeedType::Atom);
    }
}
