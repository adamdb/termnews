use crate::app::{App, InputField, ViewMode};
use ratatui::{
    layout::{Alignment, Constraint, Direction, Layout, Rect},
    style::{Color, Modifier, Style},
    text::{Line, Span, Text},
    widgets::{
        Block, BorderType, Borders, Clear, List, ListItem, ListState, Paragraph, Tabs, Wrap,
    },
    Frame,
};

const HELP_LIST: &str =
    " Tab/→: next  Shift+Tab/←: prev  ↑/k: up  ↓/j: down  Enter: read  r: refresh  t: toggle auto-refresh  s: summary  a: add  d: delete  q: quit ";
const HELP_DETAIL: &str =
    " ↑/k: scroll up  ↓/j: scroll down  Esc/Backspace: back  q: quit ";
const HELP_SUMMARY: &str =
    " ↑/k: scroll up  ↓/j: scroll down  Esc/Backspace: back to list  q: quit ";
const HELP_ADD: &str =
    " Tab: next field  Enter: confirm  Esc: cancel ";

pub fn draw(f: &mut Frame, app: &App) {
    match app.view_mode {
        ViewMode::AddSource => draw_add_source(f, app, f.area()),
        ViewMode::Summary => draw_summary(f, app),
        _ => draw_main(f, app),
    }
}

fn draw_main(f: &mut Frame, app: &App) {
    let area = f.area();

    let chunks = Layout::default()
        .direction(Direction::Vertical)
        .constraints([
            Constraint::Length(3), // Tab bar
            Constraint::Min(0),    // Content
            Constraint::Length(1), // Help / status bar
        ])
        .split(area);

    draw_tabs(f, app, chunks[0]);
    draw_content(f, app, chunks[1]);
    draw_status(f, app, chunks[2]);
}

fn draw_tabs(f: &mut Frame, app: &App, area: Rect) {
    let tab_titles: Vec<Line> = app
        .tabs
        .iter()
        .enumerate()
        .map(|(i, tab)| {
            let style = if i == app.current_tab {
                Style::default()
                    .fg(Color::Yellow)
                    .add_modifier(Modifier::BOLD)
            } else {
                Style::default().fg(Color::White)
            };

            // Add source type prefix
            let prefix = match &tab.source.source_type {
                crate::config::SourceType::Feed { feed_type, .. } => {
                    match feed_type {
                        crate::config::FeedType::Rss => "[RSS] ",
                        crate::config::FeedType::Atom => "[Atom] ",
                    }
                }
                crate::config::SourceType::Irc { .. } => "[IRC] ",
            };

            let loading = if tab.loading { " ⟳" } else { "" };
            let err = if tab.error.is_some() { " ✗" } else { "" };
            Line::from(vec![Span::styled(
                format!(" {}{}{}{} ", prefix, tab.source.name, loading, err),
                style,
            )])
        })
        .collect();

    let tabs_widget = Tabs::new(tab_titles)
        .block(
            Block::default()
                .borders(Borders::ALL)
                .border_type(BorderType::Rounded)
                .title(" termnews ")
                .title_alignment(Alignment::Center)
                .style(Style::default().fg(Color::Cyan)),
        )
        .select(app.current_tab)
        .highlight_style(
            Style::default()
                .fg(Color::Yellow)
                .add_modifier(Modifier::BOLD),
        )
        .divider(Span::styled("|", Style::default().fg(Color::DarkGray)));

    f.render_widget(tabs_widget, area);
}

fn draw_content(f: &mut Frame, app: &App, area: Rect) {
    if app.tabs.is_empty() {
        let msg = Paragraph::new(
            "No news sources configured.\nPress 'a' to add a new source.",
        )
        .alignment(Alignment::Center)
        .block(
            Block::default()
                .borders(Borders::ALL)
                .border_type(BorderType::Rounded),
        );
        f.render_widget(msg, area);
        return;
    }

    let tab = &app.tabs[app.current_tab];

    if tab.loading {
        let loading = Paragraph::new("⟳ Fetching articles…")
            .alignment(Alignment::Center)
            .block(
                Block::default()
                    .borders(Borders::ALL)
                    .border_type(BorderType::Rounded)
                    .title(format!(" {} ", tab.source.name)),
            );
        f.render_widget(loading, area);
        return;
    }

    if let Some(err) = &tab.error {
        let error_msg = Paragraph::new(format!("✗ {}", err))
            .style(Style::default().fg(Color::Red))
            .wrap(Wrap { trim: true })
            .block(
                Block::default()
                    .borders(Borders::ALL)
                    .border_type(BorderType::Rounded)
                    .title(format!(" {} ", tab.source.name)),
            );
        f.render_widget(error_msg, area);
        return;
    }

    if tab.articles.is_empty() {
        let empty = Paragraph::new("No articles found.")
            .alignment(Alignment::Center)
            .block(
                Block::default()
                    .borders(Borders::ALL)
                    .border_type(BorderType::Rounded)
                    .title(format!(" {} ", tab.source.name)),
            );
        f.render_widget(empty, area);
        return;
    }

    match app.view_mode {
        ViewMode::Detail => draw_detail(f, app, area),
        _ => draw_list(f, app, area),
    }
}

fn draw_list(f: &mut Frame, app: &App, area: Rect) {
    let tab = &app.tabs[app.current_tab];
    let total = tab.articles.len();

    let items: Vec<ListItem> = tab
        .articles
        .iter()
        .enumerate()
        .map(|(i, article)| {
            let is_selected = tab.selected == Some(i);
            let title_style = if is_selected {
                Style::default()
                    .fg(Color::Yellow)
                    .add_modifier(Modifier::BOLD)
            } else {
                Style::default().fg(Color::White)
            };
            let meta_style = Style::default().fg(Color::DarkGray);

            let mut lines = vec![Line::from(vec![Span::styled(
                format!(" {} ", article.title),
                title_style,
            )])];

            let mut meta_parts = vec![];
            if !article.pub_date.is_empty() {
                meta_parts.push(article.pub_date.clone());
            }
            if !article.author.is_empty() {
                meta_parts.push(format!("by {}", article.author));
            }
            if !meta_parts.is_empty() {
                lines.push(Line::from(vec![Span::styled(
                    format!("   {}", meta_parts.join(" · ")),
                    meta_style,
                )]));
            }
            // Separator
            lines.push(Line::from(Span::styled(
                "─".repeat(area.width.saturating_sub(4) as usize),
                Style::default().fg(Color::DarkGray),
            )));

            ListItem::new(Text::from(lines))
        })
        .collect();

    let mut list_state = ListState::default();
    list_state.select(tab.selected);

    let list = List::new(items)
        .block(
            Block::default()
                .borders(Borders::ALL)
                .border_type(BorderType::Rounded)
                .title(format!(
                    " {} ({} articles) ",
                    tab.source.name, total
                ))
                .title_style(Style::default().fg(Color::Cyan)),
        )
        .highlight_style(
            Style::default()
                .bg(Color::DarkGray)
                .add_modifier(Modifier::BOLD),
        );

    f.render_stateful_widget(list, area, &mut list_state);
}

fn draw_detail(f: &mut Frame, app: &App, area: Rect) {
    let tab = &app.tabs[app.current_tab];
    if let Some(idx) = tab.selected {
        if let Some(article) = tab.articles.get(idx) {
            let chunks = Layout::default()
                .direction(Direction::Vertical)
                .constraints([
                    Constraint::Length(5), // Title block
                    Constraint::Min(0),    // Body
                ])
                .split(area);

            // Title section
            let title_text = vec![
                Line::from(vec![Span::styled(
                    &article.title,
                    Style::default()
                        .fg(Color::Yellow)
                        .add_modifier(Modifier::BOLD),
                )]),
                Line::from(vec![Span::styled(
                    format!(
                        "{} {}",
                        if article.pub_date.is_empty() {
                            String::new()
                        } else {
                            format!("{} ", article.pub_date.clone())
                        },
                        if article.author.is_empty() {
                            String::new()
                        } else {
                            format!("| by {}", article.author)
                        }
                    )
                    .trim()
                    .to_string(),
                    Style::default().fg(Color::DarkGray),
                )]),
                Line::from(vec![Span::styled(
                    &article.link,
                    Style::default()
                        .fg(Color::Blue)
                        .add_modifier(Modifier::UNDERLINED),
                )]),
            ];

            let title_para = Paragraph::new(title_text)
                .wrap(Wrap { trim: true })
                .block(
                    Block::default()
                        .borders(Borders::ALL)
                        .border_type(BorderType::Rounded)
                        .title(format!(" {} ", tab.source.name))
                        .title_style(Style::default().fg(Color::Cyan)),
                );
            f.render_widget(title_para, chunks[0]);

            // Body / summary
            let summary = article.summary();
            let body_para = Paragraph::new(summary.as_str())
                .wrap(Wrap { trim: true })
                .scroll((tab.scroll_offset as u16, 0))
                .block(
                    Block::default()
                        .borders(Borders::ALL)
                        .border_type(BorderType::Rounded)
                        .title(" Article ")
                        .title_style(Style::default().fg(Color::Cyan)),
                );
            f.render_widget(body_para, chunks[1]);
        }
    }
}

fn draw_status(f: &mut Frame, app: &App, area: Rect) {
    let help_text = match app.view_mode {
        ViewMode::Detail => HELP_DETAIL,
        ViewMode::AddSource => HELP_ADD,
        ViewMode::Summary => HELP_SUMMARY,
        ViewMode::List => HELP_LIST,
    };

    let text = if let Some(msg) = &app.status_message {
        Span::styled(msg.clone(), Style::default().fg(Color::Green))
    } else {
        Span::styled(help_text, Style::default().fg(Color::DarkGray))
    };

    let status = Paragraph::new(Line::from(text));
    f.render_widget(status, area);
}

fn draw_add_source(f: &mut Frame, app: &App, area: Rect) {
    // Dim the background
    f.render_widget(Clear, area);

    // Render the main view behind
    draw_main_behind(f, app, area);

    // Popup dimensions
    let popup_width = 60u16.min(area.width.saturating_sub(4));
    let popup_height = 14u16.min(area.height.saturating_sub(4));
    let popup_x = (area.width.saturating_sub(popup_width)) / 2;
    let popup_y = (area.height.saturating_sub(popup_height)) / 2;
    let popup_area = Rect::new(
        area.x + popup_x,
        area.y + popup_y,
        popup_width,
        popup_height,
    );

    f.render_widget(Clear, popup_area);

    if let Some(state) = &app.add_source_state {
        let chunks = Layout::default()
            .direction(Direction::Vertical)
            .constraints([
                Constraint::Length(3), // Name
                Constraint::Length(3), // URL
                Constraint::Length(3), // Feed type
                Constraint::Min(0),    // Error / help
            ])
            .margin(1)
            .split(popup_area);

        let block = Block::default()
            .borders(Borders::ALL)
            .border_type(BorderType::Rounded)
            .title(" Add News Source ")
            .title_alignment(Alignment::Center)
            .style(Style::default().fg(Color::Cyan));
        f.render_widget(block, popup_area);

        // Name field
        let name_style = if state.active_field == InputField::Name {
            Style::default()
                .fg(Color::Yellow)
                .add_modifier(Modifier::BOLD)
        } else {
            Style::default().fg(Color::White)
        };
        let name_label = if state.active_field == InputField::Name {
            "▶ Name: "
        } else {
            "  Name: "
        };
        let name_para = Paragraph::new(format!("{}{}", name_label, state.name))
            .style(name_style)
            .block(Block::default().borders(Borders::ALL).border_type(BorderType::Plain));
        f.render_widget(name_para, chunks[0]);

        // URL field
        let url_style = if state.active_field == InputField::Url {
            Style::default()
                .fg(Color::Yellow)
                .add_modifier(Modifier::BOLD)
        } else {
            Style::default().fg(Color::White)
        };
        let url_label = if state.active_field == InputField::Url {
            "▶ URL:  "
        } else {
            "  URL:  "
        };
        let url_para = Paragraph::new(format!("{}{}", url_label, state.url))
            .style(url_style)
            .block(Block::default().borders(Borders::ALL).border_type(BorderType::Plain));
        f.render_widget(url_para, chunks[1]);

        // Feed type field
        let ft_style = if state.active_field == InputField::FeedType {
            Style::default()
                .fg(Color::Yellow)
                .add_modifier(Modifier::BOLD)
        } else {
            Style::default().fg(Color::White)
        };
        let ft_label = if state.active_field == InputField::FeedType {
            "▶ Type: "
        } else {
            "  Type: "
        };
        let ft_value = if state.feed_type_index == 0 {
            "[ RSS ] / Atom"
        } else {
            "  RSS  / [Atom]"
        };
        let ft_para = Paragraph::new(format!("{}{}", ft_label, ft_value))
            .style(ft_style)
            .block(Block::default().borders(Borders::ALL).border_type(BorderType::Plain));
        f.render_widget(ft_para, chunks[2]);

        // Error / help
        let bottom_text = if let Some(err) = &state.error {
            Span::styled(err.clone(), Style::default().fg(Color::Red))
        } else {
            Span::styled(HELP_ADD, Style::default().fg(Color::DarkGray))
        };
        let bottom = Paragraph::new(Line::from(bottom_text)).alignment(Alignment::Center);
        f.render_widget(bottom, chunks[3]);
    }
}

/// Render the main view (list/tabs) as the background for the popup.
fn draw_main_behind(f: &mut Frame, app: &App, area: Rect) {
    let chunks = Layout::default()
        .direction(Direction::Vertical)
        .constraints([
            Constraint::Length(3),
            Constraint::Min(0),
            Constraint::Length(1),
        ])
        .split(area);

    draw_tabs(f, app, chunks[0]);
    // Render the list content dimmed
    let tab = app.tabs.get(app.current_tab);
    if let Some(tab) = tab {
        let dimmed = Paragraph::new(
            tab.articles
                .iter()
                .map(|a| Line::from(Span::styled(
                    format!(" {}", a.title),
                    Style::default().fg(Color::DarkGray),
                )))
                .collect::<Vec<_>>(),
        )
        .block(
            Block::default()
                .borders(Borders::ALL)
                .border_type(BorderType::Rounded)
                .title(format!(" {} ", tab.source.name))
                .style(Style::default().fg(Color::DarkGray)),
        );
        f.render_widget(dimmed, chunks[1]);
    }
}

fn draw_summary(f: &mut Frame, app: &App) {
    let area = f.area();

    let chunks = Layout::default()
        .direction(Direction::Vertical)
        .constraints([
            Constraint::Length(3), // Title
            Constraint::Min(0),    // Content
            Constraint::Length(1), // Help / status bar
        ])
        .split(area);

    // Draw title
    let title = Block::default()
        .borders(Borders::ALL)
        .border_type(BorderType::Rounded)
        .title(" Summary - Top 5 Articles from All Feeds ")
        .title_alignment(Alignment::Center)
        .style(Style::default().fg(Color::Cyan));
    f.render_widget(title, chunks[0]);

    // Draw articles list
    if app.summary_articles.is_empty() {
        let empty = Paragraph::new("No articles available for summary.\nRefresh feeds to see content.")
            .alignment(Alignment::Center)
            .block(
                Block::default()
                    .borders(Borders::ALL)
                    .border_type(BorderType::Rounded),
            );
        f.render_widget(empty, chunks[1]);
    } else {
        let total = app.summary_articles.len();
        let items: Vec<ListItem> = app
            .summary_articles
            .iter()
            .enumerate()
            .map(|(i, (source_name, article))| {
                let is_selected = app.summary_selected == Some(i);
                let title_style = if is_selected {
                    Style::default()
                        .fg(Color::Yellow)
                        .add_modifier(Modifier::BOLD)
                } else {
                    Style::default().fg(Color::White)
                };
                let meta_style = Style::default().fg(Color::DarkGray);
                let source_style = Style::default()
                    .fg(Color::Cyan)
                    .add_modifier(Modifier::ITALIC);

                let mut lines = vec![
                    Line::from(vec![Span::styled(
                        format!(" {} ", article.title),
                        title_style,
                    )]),
                ];

                // Add source name
                lines.push(Line::from(vec![Span::styled(
                    format!("   Source: {}", source_name),
                    source_style,
                )]));

                // Add metadata
                let mut meta_parts = vec![];
                if !article.pub_date.is_empty() {
                    meta_parts.push(article.pub_date.clone());
                }
                if !article.author.is_empty() {
                    meta_parts.push(format!("by {}", article.author));
                }
                if !meta_parts.is_empty() {
                    lines.push(Line::from(vec![Span::styled(
                        format!("   {}", meta_parts.join(" · ")),
                        meta_style,
                    )]));
                }

                // Separator
                lines.push(Line::from(Span::styled(
                    "─".repeat(area.width.saturating_sub(4) as usize),
                    Style::default().fg(Color::DarkGray),
                )));

                ListItem::new(Text::from(lines))
            })
            .collect();

        let mut list_state = ListState::default();
        list_state.select(app.summary_selected);

        let list = List::new(items)
            .block(
                Block::default()
                    .borders(Borders::ALL)
                    .border_type(BorderType::Rounded)
                    .title(format!(" Summary ({} articles) ", total))
                    .title_style(Style::default().fg(Color::Cyan)),
            )
            .highlight_style(
                Style::default()
                    .bg(Color::DarkGray)
                    .add_modifier(Modifier::BOLD),
            );

        f.render_stateful_widget(list, chunks[1], &mut list_state);
    }

    draw_status(f, app, chunks[2]);
}
