# mdview — Go CLI Markdown Checklist Viewer

## Context

Build an interactive terminal markdown viewer focused on checklists with rich checkbox rendering. This is the CLI component of a two-part system (Go CLI viewer + future Bun web app) that shares `.md` files as the source of truth. Editing is done in Vim; this app is view-only with search and filter capabilities.

## Architecture Overview

```
mdview [optional-path]
┌─────────────────────────────────────────────────────┐
│ File Browser (left)  │  Markdown Viewer (right)      │
│                      │                               │
│ 📁 projects/         │  # Deploy Checklist           │
│   todo.md            │                               │
│   Heading.. 2h ago   │  ✅ Set up CI pipeline        │
│ > deploy.md          │  🔄 Write migration scripts   │
│   Deploy..  5m ago   │  ⬜ Run load tests            │
│   notes.md           │  🔥 Fix auth vulnerability    │
│   Notes..   1d ago   │  ❌ Old deploy method         │
│                      │                               │
├─────────────────────────────────────────────────────┤
│ ~/notes/projects > deploy.md  | q:quit Tab:switch   │
└─────────────────────────────────────────────────────┘
```

## Tech Stack

- **Language:** Go
- **TUI framework:** [Bubble Tea](https://github.com/charmbracelet/bubbletea) (Elm-architecture TUI)
- **Markdown rendering:** [Glamour](https://github.com/charmbracelet/glamour) with custom checkbox preprocessor
- **Styling:** [Lip Gloss](https://github.com/charmbracelet/lipgloss)
- **File watching:** [fsnotify](https://github.com/fsnotify/fsnotify)
- **Fuzzy search:** [go-fuzzyfinder](https://github.com/sahilm/fuzzy) or similar
- **Config:** YAML via `gopkg.in/yaml.v3`

## Features & Specifications

### 1. Two-Pane Layout

- **Left pane:** File browser showing `.md` files and directories
- **Right pane:** Rendered markdown content
- Pane ratio starts at ~30/70, resizable at runtime with `+`/`-`
- Tab key switches focus between panes
- When file browser has focus, cursor movement auto-previews the selected file in the viewer pane

### 2. File Browser

- Shows only `.md` files and directories (directories always listed first)
- Each entry displays: **filename**, **first `#` heading** (truncated), **last modified** (relative time: "5m ago", "2h ago", "1d ago")
- Sort order: alphabetical by default, configurable (toggle with a hotkey, persist in config)
- Empty directories show "No markdown files found" in the viewer pane, but still allow navigating to subdirectories
- **Navigation:**
  - `j`/`k` — move cursor up/down
  - `enter` — enter a subdirectory
  - `backspace` — go up a directory
  - Moving cursor auto-previews the file in the viewer pane

### 3. Markdown Viewer

- Full markdown rendering via Glamour (headings, bold, italic, links, code blocks, lists, tables, etc.)
- **Custom checkbox rendering** — preprocessor transforms checkbox syntax before Glamour renders
- When focused (via Tab):
  - `j`/`k` — scroll up/down
  - `g`/`G` — top/bottom of file

### 4. Checkbox States (5 states)

| Syntax | State       | Emoji Theme | Colored Symbol Theme |
|--------|-------------|-------------|---------------------|
| `[ ]`  | Todo        | ⬜          | □ (gray)            |
| `[x]`  | Done        | ✅          | ■ (green)           |
| `[/]`  | In-progress | 🔄          | ▶ (blue)            |
| `[-]`  | Cancelled   | ❌          | – (dim/strikethrough)|
| `[!]`  | Urgent      | 🔥          | ! (red bold)        |

- Theme is selectable: `emoji`, `symbols`, or `custom`
- Custom themes defined in config file with user-specified glyphs/colors per state

### 5. Command System

Colon prefix (vim-style). Press `:` to open command input in the status bar. Type command + Enter to execute. Esc to cancel.

#### `:s <query>` — Fuzzy Search
- Searches filenames and file content across all `.md` files in the current folder and all subfolders (recursive)
- Fuzzy matching (not exact substring)
- Results displayed as a **full-screen overlay** list
- Each result shows: file path, matching line/context
- `j`/`k` to navigate results, `enter` to jump to that file, `Esc` to dismiss

#### `:f <status> [status...]` — Filter by Checkbox Status
- Filters the currently displayed content by checkbox state
- Accepts status names: `todo`, `done`, `inprogress`, `cancelled`, `urgent`
- Multiple statuses combinable: `:f todo urgent` shows both `[ ]` and `[!]` items
- `:f clear` removes all filters
- When active, filter indicator shown in status bar

### 6. Configuration (`~/.mdview.yaml`)

```yaml
# Theme: emoji | symbols | custom
theme: emoji

# Custom theme example (only used when theme: custom)
custom_theme:
  todo: "○"
  done: "●"
  inprogress: "◐"
  cancelled: "⊘"
  urgent: "⚠"

# Keybinding style: vim | arrows
keybindings: vim

# Default sort: alpha | modified
sort: alpha

# Sidebar width percentage (10-50)
sidebar_width: 30
```

### 7. File Watching (Auto-reload)

- Use fsnotify to watch the currently displayed file
- On file change, re-read and re-render automatically
- Also watch current directory for new/deleted/renamed `.md` files to update the file browser
- Debounce rapid changes (e.g. 100ms)

### 8. CLI Usage

```
mdview              # opens in current directory
mdview ~/notes      # opens in specified directory
mdview ~/notes/todo.md  # opens with that file selected
```

### 9. Keybinding Summary

| Key        | Context       | Action                           |
|------------|---------------|----------------------------------|
| `j`/`k`    | File browser  | Move cursor up/down              |
| `j`/`k`    | Viewer        | Scroll up/down                   |
| `enter`    | File browser  | Enter directory                  |
| `backspace`| File browser  | Go up a directory                |
| `Tab`      | Global        | Switch focus between panes       |
| `+`/`-`    | Global        | Resize panes                     |
| `g`/`G`    | Viewer        | Jump to top/bottom               |
| `:`        | Global        | Open command input               |
| `q`        | Global        | Quit                             |
| `Esc`      | Command/Overlay| Cancel/dismiss                  |

### 10. Status Bar (Bottom)

Single bottom bar showing:
- Left: current path as breadcrumb (e.g. `~/notes/projects > deploy.md`)
- Right: contextual hotkey hints (changes based on active pane/mode)
- Center: active filter indicator when `:f` is applied (e.g. `[filter: todo, urgent]`)

## Project Structure

```
markdown-viewer/
├── main.go                  # Entry point, CLI arg parsing
├── go.mod
├── go.sum
├── ARCHITECTURE.md
├── internal/
│   ├── app/
│   │   └── app.go           # Root Bubble Tea model, orchestrates sub-models
│   ├── browser/
│   │   ├── browser.go       # File browser pane model
│   │   └── entry.go         # File entry type (name, heading, modified)
│   ├── viewer/
│   │   └── viewer.go        # Markdown viewer pane model
│   ├── search/
│   │   ├── search.go        # Fuzzy search overlay model
│   │   └── fuzzy.go         # Fuzzy matching logic
│   ├── command/
│   │   └── command.go       # Command input bar model (: commands)
│   ├── filter/
│   │   └── filter.go        # Checkbox status filter logic
│   ├── markdown/
│   │   ├── renderer.go      # Glamour wrapper + checkbox preprocessor
│   │   └── checkbox.go      # Checkbox parsing and state definitions
│   ├── theme/
│   │   └── theme.go         # Theme definitions (emoji, symbols, custom)
│   ├── config/
│   │   └── config.go        # YAML config loading from ~/.mdview.yaml
│   └── watcher/
│       └── watcher.go       # fsnotify file watcher integration
└── testdata/
    ├── sample.md             # Test markdown with all checkbox states
    └── nested/
        └── sub.md            # Test nested directory support
```

## Implementation Order

### Phase 1: Foundation
1. `go mod init`, install dependencies (bubbletea, glamour, lipgloss, fsnotify, yaml)
2. Config loading (`internal/config/`)
3. Checkbox parsing and theme rendering (`internal/markdown/`, `internal/theme/`)

### Phase 2: Core UI
4. Root app model with two-pane layout (`internal/app/`)
5. File browser model — list `.md` files, show name/heading/date (`internal/browser/`)
6. Markdown viewer model — render with Glamour + checkbox preprocessor (`internal/viewer/`)
7. Pane focus switching (Tab), resizing (+/-), status bar

### Phase 3: Navigation
8. File browser navigation (j/k/enter/backspace)
9. Auto-preview on cursor movement
10. Viewer scrolling (j/k/g/G)
11. Directory traversal (enter into subfolder, backspace to parent)

### Phase 4: Commands & Search
12. Command input bar (`:` prefix) (`internal/command/`)
13. Fuzzy search (`:s`) with full-screen overlay (`internal/search/`)
14. Status filter (`:f`) with combinable statuses (`internal/filter/`)

### Phase 5: File Watching & Polish
15. fsnotify integration for auto-reload (`internal/watcher/`)
16. Arrow key alternative keybindings (config-driven)
17. Edge cases: terminal resize, very large files
18. Empty directory handling: show "No markdown files found" in viewer pane, still allow navigating to subdirectories via the file browser

## Verification

1. **Build:** `go build -o mdview .`
2. **Basic launch:** `./mdview testdata/` — should show two-pane layout with sample files
3. **Navigation:** j/k moves cursor, enter opens folders, backspace goes up, Tab switches panes
4. **Checkbox rendering:** Open sample.md with all 5 checkbox states, verify correct icons per theme
5. **Search:** `:s <term>` shows overlay with fuzzy results, enter jumps to file
6. **Filter:** `:f todo urgent` filters to only those states, `:f clear` resets
7. **Auto-reload:** Edit a file in Vim in another terminal, viewer updates automatically
8. **Resize:** +/- adjusts pane widths, terminal resize reflows layout
9. **Config:** Modify `~/.mdview.yaml` theme setting, restart, verify theme changes
