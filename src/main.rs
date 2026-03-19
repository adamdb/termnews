mod app;
mod config;
mod feed;
mod ui;

use app::{App, InputField, ViewMode};
use config::Config;
use crossterm::{
    event::{self, DisableMouseCapture, EnableMouseCapture, Event, KeyCode, KeyModifiers},
    execute,
    terminal::{disable_raw_mode, enable_raw_mode, EnterAlternateScreen, LeaveAlternateScreen},
};
use ratatui::{backend::CrosstermBackend, Terminal};
use std::{
    io,
    sync::{Arc, Mutex},
    thread,
    time::Duration,
};

fn main() -> io::Result<()> {
    // Load configuration
    let config = Config::load();

    // Save defaults if no config file exists
    if !Config::config_path().exists() {
        let _ = config.save();
    }

    let app = Arc::new(Mutex::new(App::new(config)));

    // Spawn background threads to fetch feeds
    {
        let app_clone = Arc::clone(&app);
        let sources: Vec<_> = {
            let a = app_clone.lock().unwrap();
            a.config.sources.clone()
        };
        for (i, source) in sources.into_iter().enumerate() {
            let app_ref = Arc::clone(&app_clone);
            thread::spawn(move || {
                let result = feed::fetch_feed(&source);
                let mut app = app_ref.lock().unwrap();
                app.set_articles(i, result);
            });
        }
    }

    // Set up terminal
    enable_raw_mode()?;
    let mut stdout = io::stdout();
    execute!(stdout, EnterAlternateScreen, EnableMouseCapture)?;
    let backend = CrosstermBackend::new(stdout);
    let mut terminal = Terminal::new(backend)?;

    // Run the event loop
    let result = run_app(&mut terminal, Arc::clone(&app));

    // Restore terminal
    disable_raw_mode()?;
    execute!(
        terminal.backend_mut(),
        LeaveAlternateScreen,
        DisableMouseCapture
    )?;
    terminal.show_cursor()?;

    if let Err(e) = result {
        eprintln!("Error: {}", e);
    }

    Ok(())
}

fn run_app(
    terminal: &mut ratatui::Terminal<CrosstermBackend<io::Stdout>>,
    app: Arc<Mutex<App>>,
) -> io::Result<()> {
    loop {
        {
            let app_lock = app.lock().unwrap();
            terminal.draw(|f| ui::draw(f, &app_lock))?;
        }

        // Check if auto-refresh is needed
        {
            let mut app_lock = app.lock().unwrap();
            if app_lock.should_auto_refresh() {
                app_lock.refresh_all();
                let sources: Vec<_> = app_lock
                    .tabs
                    .iter()
                    .map(|t| t.source.clone())
                    .collect();
                drop(app_lock);

                // Spawn threads to fetch all feeds
                for (i, source) in sources.into_iter().enumerate() {
                    let app_ref = Arc::clone(&app);
                    thread::spawn(move || {
                        let result = feed::fetch_feed(&source);
                        let mut a = app_ref.lock().unwrap();
                        a.set_articles(i, result);
                    });
                }
            }
        }

        // Poll for events with a short timeout so we can re-render on feed updates
        if event::poll(Duration::from_millis(200))? {
            if let Event::Key(key) = event::read()? {
                let mut app_lock = app.lock().unwrap();

                if app_lock.should_quit {
                    return Ok(());
                }

                let mode = app_lock.view_mode.clone();
                match mode {
                    ViewMode::AddSource => handle_add_source_input(&mut app_lock, key, Arc::clone(&app)),
                    ViewMode::Detail => handle_detail_input(&mut app_lock, key),
                    ViewMode::Summary => handle_summary_input(&mut app_lock, key),
                    ViewMode::List => handle_list_input(&mut app_lock, key, Arc::clone(&app)),
                }

                if app_lock.should_quit {
                    return Ok(());
                }
            }
        }
    }
}

fn handle_list_input(
    app: &mut App,
    key: crossterm::event::KeyEvent,
    app_arc: Arc<Mutex<App>>,
) {
    match key.code {
        KeyCode::Char('q') | KeyCode::Char('Q') => {
            app.should_quit = true;
        }
        KeyCode::Tab => {
            app.next_tab();
        }
        KeyCode::BackTab => {
            app.prev_tab();
        }
        KeyCode::Right => {
            app.next_tab();
        }
        KeyCode::Left => {
            app.prev_tab();
        }
        KeyCode::Down | KeyCode::Char('j') => {
            app.scroll_down();
        }
        KeyCode::Up | KeyCode::Char('k') => {
            app.scroll_up();
        }
        KeyCode::Enter => {
            app.select_article();
        }
        KeyCode::Char('t') | KeyCode::Char('T') => {
            app.toggle_auto_refresh();
        }
        KeyCode::Char('s') | KeyCode::Char('S') => {
            app.enter_summary_mode();
        }
        KeyCode::Char('r') | KeyCode::Char('R') => {
            let idx = app.current_tab;
            if app.tabs.is_empty() {
                return;
            }
            let source = app.tabs[idx].source.clone();
            app.refresh_current();
            // Drop the mutable reference to `app` by returning, then
            // the caller re-acquires the lock for re-render.
            // We spawn a thread here before dropping — the Arc is cloned first.
            let app_ref = Arc::clone(&app_arc);
            thread::spawn(move || {
                let result = feed::fetch_feed(&source);
                let mut a = app_ref.lock().unwrap();
                a.set_articles(idx, result);
            });
        }
        KeyCode::Char('a') | KeyCode::Char('A') => {
            app.start_add_source();
        }
        KeyCode::Char('d') | KeyCode::Char('D') => {
            if !app.tabs.is_empty() {
                app.remove_current_source();
            }
        }
        KeyCode::Char('c') if key.modifiers.contains(KeyModifiers::CONTROL) => {
            app.should_quit = true;
        }
        _ => {}
    }
}

fn handle_detail_input(app: &mut App, key: crossterm::event::KeyEvent) {
    match key.code {
        KeyCode::Char('q') | KeyCode::Char('Q') => {
            app.should_quit = true;
        }
        KeyCode::Esc | KeyCode::Backspace => {
            app.back_to_list();
        }
        KeyCode::Down | KeyCode::Char('j') => {
            app.scroll_down();
        }
        KeyCode::Up | KeyCode::Char('k') => {
            app.scroll_up();
        }
        KeyCode::Char('c') if key.modifiers.contains(KeyModifiers::CONTROL) => {
            app.should_quit = true;
        }
        _ => {}
    }
}

fn handle_summary_input(app: &mut App, key: crossterm::event::KeyEvent) {
    match key.code {
        KeyCode::Char('q') | KeyCode::Char('Q') => {
            app.should_quit = true;
        }
        KeyCode::Esc | KeyCode::Backspace => {
            app.exit_summary_mode();
        }
        KeyCode::Down | KeyCode::Char('j') => {
            app.scroll_summary_down();
        }
        KeyCode::Up | KeyCode::Char('k') => {
            app.scroll_summary_up();
        }
        KeyCode::Char('c') if key.modifiers.contains(KeyModifiers::CONTROL) => {
            app.should_quit = true;
        }
        _ => {}
    }
}

fn handle_add_source_input(
    app: &mut App,
    key: crossterm::event::KeyEvent,
    app_arc: Arc<Mutex<App>>,
) {
    match key.code {
        KeyCode::Esc => {
            app.cancel_add_source();
        }
        KeyCode::Tab => {
            if let Some(state) = &mut app.add_source_state {
                state.error = None;
                state.active_field = match state.active_field {
                    InputField::Name => InputField::Url,
                    InputField::Url => InputField::FeedType,
                    InputField::FeedType => InputField::Name,
                };
            }
        }
        KeyCode::Enter => {
            let new_tab_idx = app.tabs.len();
            if app.confirm_add_source() {
                // Fetch articles for the newly added source
                if let Some(tab) = app.tabs.get(new_tab_idx) {
                    let source = tab.source.clone();
                    // Mark as loading
                    if let Some(t) = app.tabs.get_mut(new_tab_idx) {
                        t.loading = true;
                    }
                    let app_ref = Arc::clone(&app_arc);
                    thread::spawn(move || {
                        let result = feed::fetch_feed(&source);
                        let mut a = app_ref.lock().unwrap();
                        a.set_articles(new_tab_idx, result);
                    });
                }
                app.status_message = None;
            }
        }
        KeyCode::Backspace => {
            if let Some(state) = &mut app.add_source_state {
                match state.active_field {
                    InputField::Name => {
                        state.name.pop();
                    }
                    InputField::Url => {
                        state.url.pop();
                    }
                    _ => {}
                }
                state.error = None;
            }
        }
        KeyCode::Left | KeyCode::Right => {
            if let Some(state) = &mut app.add_source_state {
                if state.active_field == InputField::FeedType {
                    state.feed_type_index = 1 - state.feed_type_index;
                }
            }
        }
        KeyCode::Char(c) => {
            if let Some(state) = &mut app.add_source_state {
                match state.active_field {
                    InputField::Name => state.name.push(c),
                    InputField::Url => state.url.push(c),
                    InputField::FeedType => {
                        if c == ' ' {
                            state.feed_type_index = 1 - state.feed_type_index;
                        }
                    }
                }
                state.error = None;
            }
        }
        _ => {}
    }
}

