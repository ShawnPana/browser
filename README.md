# browser

A CLI for browser automation via Chrome DevTools Protocol, built on [go-rod](https://github.com/go-rod/rod).

Launch a headless Chrome, then drive it entirely from the command line — navigate pages, fill forms, click buttons, take screenshots, run JavaScript, inspect the accessibility tree, and more. The browser runs as a persistent background process, so each CLI invocation is fast and stateless.

## Install

```bash
go install browser@latest
```

Or build from source:

```bash
git clone https://github.com/user/browser.git
cd browser
go build -o browser .
```

Requires Chrome or Chromium installed on the system (or set `ROD_CHROME_BIN`).

## Quick Start

```bash
# Launch a headless Chrome
browser start

# Navigate to a page
browser open https://example.com

# Inspect the page
browser get title
browser ax-tree --depth 3

# Interact
browser click 'a[href="/more"]'
browser input '#search' 'hello world'
browser screenshot result.png   # saves to ~/.browser/tmp/result.png; use ./result.png for CWD

# Clean up
browser stop
```

## Commands

### Browser Lifecycle

| Command | Description |
|---------|-------------|
| `browser start` | Launch headless Chrome |
| `browser start --show` | Launch with UI visible |
| `browser start -k` | Launch with `--ignore-certificate-errors` |
| `browser stop` | Close the active browser |
| `browser connect <host:port>` | Connect to an existing Chrome instance |
| `browser connect <https://url>` | Connect to a cloud browser |
| `browser connect <index>` | Switch active browser by registry index |
| `browser status` | Show all registered browsers with liveness |

The browser persists as a background process across CLI invocations. State is stored in `~/.browser/state.json`, and output files (screenshots, PDFs, downloads) default to `~/.browser/tmp/`. Multiple browsers can be registered simultaneously — one is active at a time.

### Navigation

| Command | Description |
|---------|-------------|
| `browser open <url>` | Navigate to URL (waits for load) |
| `browser back` | Go back in history |
| `browser forward` | Go forward in history |
| `browser reload` | Reload page |
| `browser reload --hard` | Hard reload (bypass cache) |

### Page Info

| Command | Description |
|---------|-------------|
| `browser get url` | Print current URL |
| `browser get title` | Print page title |
| `browser get html` | Print full page HTML |
| `browser get html <selector>` | Print element's outer HTML |
| `browser get text <selector>` | Print element's visible text |
| `browser get attr <selector> <name>` | Print element attribute value |
| `browser get value <selector>` | Print input element value |
| `browser get box <selector>` | Print bounding box (x, y, width, height) |
| `browser get styles <sel> [prop...]` | Print computed styles |

### Interaction

All interaction commands accept either a **CSS selector** or **x,y coordinates**.

**With CSS selectors:**

| Command | Description |
|---------|-------------|
| `browser click <selector>` | Click element |
| `browser dblclick <selector>` | Double-click element |
| `browser rightclick <selector>` | Right-click element |
| `browser input <selector> <text>` | Clear field and type text (replace) |
| `browser type <selector> <text>` | Type text without clearing (append) |
| `browser press <key>` | Press key combo (e.g. `Enter`, `Control+a`) |
| `browser clear <selector>` | Clear input field |
| `browser select <selector> <value>` | Select dropdown option by value |
| `browser submit <selector>` | Submit form |
| `browser hover <selector>` | Hover over element |
| `browser focus <selector>` | Focus element |
| `browser check <selector>` | Check checkbox |
| `browser uncheck <selector>` | Uncheck checkbox |
| `browser scrollintoview <selector>` | Scroll element into viewport |

**With coordinates:**

| Command | Description |
|---------|-------------|
| `browser click <x> <y>` | Click at coordinates |
| `browser dblclick <x> <y>` | Double-click at coordinates |
| `browser rightclick <x> <y>` | Right-click at coordinates |
| `browser input <x> <y> <text>` | Click at coordinates, clear, then type |
| `browser type <x> <y> <text>` | Click at coordinates, then type (no clear) |
| `browser hover <x> <y>` | Hover at coordinates |
| `browser scroll <x> <y> <delta>` | Scroll at coordinates (delta in px) |
| `browser scroll <selector> <delta>` | Scroll inside an element |
| `browser drag <x1> <y1> <x2> <y2>` | Drag from one point to another |
| `browser element-at <x> <y>` | Describe the DOM element at coordinates |

### Keyboard

Send keystrokes without a selector (acts on the currently focused element):

| Command | Description |
|---------|-------------|
| `browser keyboard type <text>` | Type text with real keystrokes |
| `browser keyboard inserttext <text>` | Insert text without key events |

### JavaScript

```bash
browser eval <expression>
```

Evaluates a JavaScript expression in the page context. The expression is auto-wrapped in `() => { return (expr); }`.

```bash
browser eval 'document.title'                                    # → Example Domain
browser eval '2 + 2'                                             # → 4
browser eval 'document.querySelectorAll("a").length'             # → 12
browser eval 'JSON.stringify([...document.querySelectorAll("h2")].map(e => e.textContent))'
```

Output formatting: strings are printed unquoted, numbers and booleans are printed raw, objects and arrays are pretty-printed as JSON.

### File Operations

| Command | Description |
|---------|-------------|
| `browser file <selector> <path>` | Upload a file to a file input |
| `browser file <selector> -` | Upload from stdin |
| `browser download <selector> [file]` | Download resource from element's `href`/`src` |
| `browser download <selector> -` | Download to stdout |

```bash
browser file 'input[type="file"]' ./document.pdf
cat image.png | browser file 'input[type="file"]' -
browser download 'a.download-link' ./output.pdf
```

### Waiting

| Command | Description |
|---------|-------------|
| `browser wait <selector>` | Wait for element to become visible |
| `browser wait-load` | Wait for page load event |
| `browser wait-stable` | Wait for DOM stability (300ms) |
| `browser wait-idle` | Wait for requestIdleCallback (5s) |
| `browser sleep <seconds>` | Sleep for a duration |

### Output

Bare filenames default to `~/.browser/tmp/`. Use an explicit path (`./`, `../`, `/`) to write elsewhere.

```bash
browser screenshot                          # → ~/.browser/tmp/screenshot.png (auto-named)
browser screenshot page.png                 # → ~/.browser/tmp/page.png
browser screenshot ./page.png               # → ./page.png (explicit path)
browser screenshot -w 1920 file.png         # Custom width, full-page height
browser screenshot -w 1920 -h 1080 file.png # Fixed viewport clip
browser pdf page.pdf                        # → ~/.browser/tmp/page.pdf
browser pdf ./page.pdf                      # → ./page.pdf (explicit path)
```

When `-h` is specified, the screenshot clips to the viewport height. Without it, the full scrollable page is captured.

### Tabs

| Command | Description |
|---------|-------------|
| `browser tabs` | List all open tabs |
| `browser switch <index>` | Switch active tab |
| `browser new-tab [url]` | Open a new tab |
| `browser close-tab [index]` | Close a tab (defaults to active) |

```bash
browser new-tab https://site-a.com
browser new-tab https://site-b.com
browser tabs
# * [0] about:blank - about:blank
#   [1] Site A - https://site-a.com
#   [2] Site B - https://site-b.com

browser switch 1     # switch to Site A
browser close-tab 0  # close the blank tab
```

### Checks and Assertions

These commands use exit codes for scripting: **0** = pass, **1** = fail, **2** = error.

| Command | Description |
|---------|-------------|
| `browser exists <selector>` | Check if element exists |
| `browser visible <selector>` | Check if element is visible |
| `browser count <selector>` | Count matching elements |
| `browser assert <expr>` | Assert JS expression is truthy |
| `browser assert <expr> <expected>` | Assert JS result equals expected string |
| `browser assert <expr> <expected> -m <msg>` | Assert with custom failure message |

```bash
browser exists 'h1' && echo "found"
browser assert 'document.title' 'My Page'
browser assert 'location.pathname' '/dashboard' -m 'should be on dashboard'
```

### Accessibility

Inspect the accessibility tree to understand page structure without relying on CSS selectors or visual layout.

| Command | Description |
|---------|-------------|
| `browser ax-tree` | Print full accessibility tree |
| `browser ax-tree --depth N` | Limit tree depth |
| `browser ax-tree --json` | Output as JSON |
| `browser ax-find --name <name>` | Find nodes by accessible name |
| `browser ax-find --role <role>` | Find nodes by ARIA role |
| `browser ax-node <selector>` | Get accessibility info for a specific element |
| `browser ax-node <selector> --json` | Output as JSON |

```bash
browser ax-tree --depth 3
# [RootWebArea] "Example" (url=https://example.com)
#   [navigation] "Main"
#     [link] "Home" (focusable=true)
#     [link] "About" (focusable=true)
#   [main] ""
#     [heading] "Welcome" (level=1)

browser ax-find --role button
browser ax-find --name "Submit"
browser ax-node '#save-btn' --json
```

### Cloud (Browser Use API)

Manage cloud browsers via the [Browser Use](https://browser-use.com) API.

| Command | Description |
|---------|-------------|
| `browser cloud login <api-key>` | Save API key |
| `browser cloud logout` | Remove API key |
| `browser cloud <METHOD> <path> [body]` | REST passthrough |
| `browser cloud poll <task-id>` | Poll a task until completion |
| `browser cloud --help` | Show API endpoints (fetches OpenAPI spec) |

```bash
browser cloud login "$BROWSER_USE_API_KEY"
browser cloud POST /browsers '{"headless": true}'
browser connect https://<uuid>.cdp0.browser-use.com

# All commands work identically on cloud browsers
browser open https://example.com
browser screenshot page.png   # → ~/.browser/tmp/page.png

# Stop auto-detects cloud URLs and calls the API
browser stop

# Run a cloud task
browser cloud POST /tasks '{"url":"https://example.com","task":"Extract the main heading"}'
browser cloud poll <task-id>
```

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success, or check passed (`exists`, `visible`, `assert`) |
| 1 | Check failed (`exists`, `visible`, `assert` returned false) |
| 2 | Error (bad arguments, no browser, timeout, network failure) |

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `BROWSER_HOME` | State directory (also controls output dir: `$BROWSER_HOME/tmp/`) | `~/.browser` |
| `ROD_TIMEOUT` | Command timeout in seconds | `30` |
| `ROD_CHROME_BIN` | Chrome binary path | auto-detect |
| `BROWSER_USE_API_KEY` | Cloud API key (alternative to `cloud login`) | none |

## Architecture

```
browser/
├── main.go              # Entry point
├── cli/                 # Command routing and arg parsing
│   ├── root.go          # Main switch/case dispatcher, help text
│   ├── utils.go         # Output formatting helpers
│   ├── browser_cmds.go  # start, stop, connect, status
│   ├── nav_cmds.go      # open, back, forward, reload
│   ├── page_cmds.go     # get url/title/html/text/attr/value/box/styles
│   ├── interact_cmds.go # click, dblclick, rightclick, input, type, press, ...
│   ├── keyboard_cmds.go # keyboard type, keyboard inserttext
│   ├── wait_cmds.go     # wait, wait-load, wait-stable, wait-idle, sleep
│   ├── screenshot_cmds.go
│   ├── pdf_cmds.go      # pdf
│   ├── tab_cmds.go      # tabs, switch, new-tab, close-tab
│   ├── assert_cmds.go   # exists, count, visible, assert
│   ├── ax_cmds.go       # ax-tree, ax-find, ax-node
│   └── cloud_cmds.go    # cloud login, REST passthrough, poll
├── core/                # State management and browser connectivity
│   ├── context.go       # Context (state dir, timeout)
│   ├── state.go         # JSON state persistence (~/.browser/state.json)
│   ├── service.go       # WithPage / WithBrowser — connect to active browser
│   ├── launcher.go      # Chrome launch via rod/launcher
│   ├── cloud.go         # Cloud config and cache
│   └── tools/           # Stateless tool functions
│       ├── navigate.go
│       ├── pageinfo.go
│       ├── interact.go
│       ├── keyboard.go
│       ├── js.go
│       ├── pdf.go
│       ├── screenshot.go
│       ├── tabs.go
│       ├── wait.go
│       ├── assert.go
│       ├── accessibility.go
│       └── file.go
```

Three layers:

- **`core/tools/`** — Pure functions that take a `*rod.Page` or `*rod.Browser` and do one thing. No state, no CLI concerns.
- **`core/`** — Browser lifecycle and persistent state. Manages a JSON registry of browser connections so the process survives across CLI calls.
- **`cli/`** — Thin command layer. Parses args, calls into core, prints results.

## License

MIT
