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
browser title
browser ax-tree --depth 3

# Interact
browser click 'a[href="/more"]'
browser input '#search' 'hello world'
browser screenshot result.png

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

The browser persists as a background process across CLI invocations. State is stored in `~/.browser/state.json`. Multiple browsers can be registered simultaneously — one is active at a time.

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
| `browser url` | Print current URL |
| `browser title` | Print page title |
| `browser html` | Print full page HTML |
| `browser html <selector>` | Print element's outer HTML |
| `browser text <selector>` | Print element's visible text |
| `browser attr <selector> <name>` | Print element attribute value |

### Interaction

All interaction commands accept either a **CSS selector** or **x,y coordinates**.

**With CSS selectors:**

| Command | Description |
|---------|-------------|
| `browser click <selector>` | Click element |
| `browser input <selector> <text>` | Clear field and type text |
| `browser clear <selector>` | Clear input field |
| `browser select <selector> <value>` | Select dropdown option by value |
| `browser submit <selector>` | Submit form |
| `browser hover <selector>` | Hover over element |
| `browser focus <selector>` | Focus element |

**With coordinates:**

| Command | Description |
|---------|-------------|
| `browser click <x> <y>` | Click at coordinates |
| `browser input <x> <y> <text>` | Click at coordinates, then type |
| `browser hover <x> <y>` | Hover at coordinates |
| `browser scroll <x> <y> <delta>` | Scroll at coordinates (delta in px) |
| `browser scroll <selector> <delta>` | Scroll inside an element |
| `browser drag <x1> <y1> <x2> <y2>` | Drag from one point to another |
| `browser element-at <x> <y>` | Describe the DOM element at coordinates |

### JavaScript

```bash
browser js <expression>
```

Evaluates a JavaScript expression in the page context. The expression is auto-wrapped in `() => { return (expr); }`.

```bash
browser js 'document.title'                                    # → Example Domain
browser js '2 + 2'                                             # → 4
browser js 'document.querySelectorAll("a").length'             # → 12
browser js 'JSON.stringify([...document.querySelectorAll("h2")].map(e => e.textContent))'
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
| `browser waitload` | Wait for page load event |
| `browser waitstable` | Wait for DOM stability (300ms) |
| `browser waitidle` | Wait for requestIdleCallback (5s) |
| `browser sleep <seconds>` | Sleep for a duration |

### Screenshots

```bash
browser screenshot                          # → screenshot.png (auto-named)
browser screenshot page.png                 # → page.png
browser screenshot -w 1920 file.png         # Custom width, full-page height
browser screenshot -w 1920 -h 1080 file.png # Fixed viewport clip
```

When `-h` is specified, the screenshot clips to the viewport height. Without it, the full scrollable page is captured.

### Tabs

| Command | Description |
|---------|-------------|
| `browser pages` | List all open pages |
| `browser page <index>` | Switch active page |
| `browser newpage [url]` | Open a new tab |
| `browser closepage [index]` | Close a tab (defaults to active) |

```bash
browser newpage https://site-a.com
browser newpage https://site-b.com
browser pages
# * [0] about:blank - about:blank
#   [1] Site A - https://site-a.com
#   [2] Site B - https://site-b.com

browser page 1      # switch to Site A
browser closepage 0  # close the blank tab
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
browser screenshot page.png

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
| `BROWSER_HOME` | State directory | `~/.browser` |
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
│   ├── page_cmds.go     # url, title, html, text, attr
│   ├── interact_cmds.go # click, input, hover, scroll, drag, ...
│   ├── wait_cmds.go     # wait, waitload, waitstable, waitidle, sleep
│   ├── screenshot_cmds.go
│   ├── tab_cmds.go      # pages, page, newpage, closepage
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
│       ├── js.go
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
