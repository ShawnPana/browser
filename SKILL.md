---
name: browser-use
description: Browser automation CLI for AI agents using Chrome DevTools Protocol. Use when the user needs to interact with websites, including navigating pages, filling forms, clicking buttons, taking screenshots, extracting data, testing web apps, or automating any browser task. Triggers include requests to "open a website", "fill out a form", "click a button", "take a screenshot", "scrape data from a page", "test this web app", "login to a site", "automate browser actions", or any task requiring programmatic web interaction. Also use for Browser Use Cloud API operations.
allowed-tools: Bash(browser:*)
---

# Browser Automation with `browser`

## Core Workflow

Every browser automation follows this pattern:

1. **Start**: `browser start` (launch Chrome) or `browser connect <url>` (attach to existing)
2. **Navigate**: `browser open <url>`
3. **Inspect**: `browser ax-tree` or `browser screenshot` to understand the page
4. **Interact**: Use CSS selectors or x,y coordinates to click, type, etc.
5. **Verify**: `browser exists <sel>`, `browser assert <expr>`, or `browser screenshot`

```bash
browser start
browser open https://example.com/form
browser ax-tree --depth 3
# Identify elements from accessibility tree

browser click 'input[name="email"]'
browser input 'input[name="email"]' 'user@example.com'
browser input 'input[name="password"]' 'secret123'
browser click 'button[type="submit"]'
browser wait-load
browser get title
```

## Command Chaining

Commands can be chained with `&&` in a single shell invocation. The browser persists between commands (Chrome runs as a separate process), so chaining is safe and efficient.

```bash
# Navigate and verify
browser open https://example.com && browser wait-load && browser get title

# Fill a form
browser input '#email' 'user@example.com' && browser input '#pass' 'secret' && browser click '#submit'

# Navigate and capture
browser open https://example.com && browser wait-load && browser screenshot page.png  # → ~/.browser/tmp/page.png
```

**When to chain:** Use `&&` when you don't need intermediate output. Run commands separately when you need to parse output first (e.g., `ax-tree` to discover elements, then interact).

## Essential Commands

```bash
# Browser Lifecycle
browser start                      # Launch headless Chrome
browser start --show               # Launch with UI visible
browser start -k                   # Launch with --ignore-certificate-errors
browser stop                       # Close active browser
browser connect <host:port>        # Connect to existing Chrome
browser connect <https://url>      # Connect to cloud browser
browser connect <index>            # Switch active browser by registry index
browser status                     # Show all browsers with liveness

# Navigation
browser open <url>                 # Navigate to URL and wait for load
browser back                       # Go back in history
browser forward                    # Go forward in history
browser reload                     # Reload page
browser reload --hard              # Hard reload (bypass cache)

# Page Info
browser get url                    # Print current URL
browser get title                  # Print page title
browser get html                   # Print full page HTML
browser get html <selector>        # Print element's outer HTML
browser get text <selector>        # Print element's visible text
browser get attr <selector> <name> # Print element attribute value
browser get value <selector>       # Print input element value
browser get box <selector>         # Print bounding box {x, y, width, height}
browser get styles <sel> [prop...] # Print computed styles (all or specific)

# Interaction (CSS selectors)
browser click <selector>           # Click element
browser dblclick <selector>        # Double-click element
browser rightclick <selector>      # Right-click element
browser input <selector> <text>    # Clear and type text (replace contents)
browser type <selector> <text>     # Type text without clearing (append)
browser press <key>                # Press key combo (Enter, Control+a, Shift+Tab)
browser clear <selector>           # Clear input field
browser select <selector> <value>  # Select dropdown option by value
browser submit <selector>          # Submit form
browser hover <selector>           # Hover over element
browser focus <selector>           # Focus element
browser check <selector>           # Check checkbox
browser uncheck <selector>         # Uncheck checkbox
browser scrollintoview <selector>  # Scroll element into viewport

# Interaction (x,y coordinates)
browser click <x> <y>              # Click at coordinates
browser dblclick <x> <y>           # Double-click at coordinates
browser rightclick <x> <y>         # Right-click at coordinates
browser input <x> <y> <text>       # Click at coords, clear, then type
browser type <x> <y> <text>        # Click at coords then type (no clear)
browser hover <x> <y>              # Hover at coordinates
browser scroll <x> <y> <delta>     # Scroll at coordinates
browser drag <x1> <y1> <x2> <y2>  # Drag from point to point
browser element-at <x> <y>         # Describe element at coordinates

# Keyboard (no selector — acts on focused element)
browser keyboard type <text>       # Type with real keystrokes
browser keyboard inserttext <text> # Insert text without key events

# File Operations
browser file <selector> <path|->   # Upload file (- for stdin)
browser download <selector> [file] # Download linked resource (default: ~/.browser/tmp/)
browser download <selector> -      # Download to stdout

# Waiting
browser wait <selector>            # Wait for element to be visible
browser wait-load                  # Wait for page load event
browser wait-stable                # Wait for DOM stability
browser wait-idle                  # Wait for idle callback
browser sleep <seconds>            # Sleep for duration

# Output (bare filenames default to ~/.browser/tmp/; use ./ or / for explicit paths)
browser screenshot                 # → ~/.browser/tmp/screenshot.png
browser screenshot file.png        # → ~/.browser/tmp/file.png
browser screenshot ./file.png      # → ./file.png (explicit path)
browser screenshot --full file.png # Full page → ~/.browser/tmp/file.png
browser pdf page.pdf               # → ~/.browser/tmp/page.pdf

# Tabs
browser tabs                       # List all open tabs
browser switch <index>             # Switch to tab by index
browser new-tab [url]              # Open new tab
browser close-tab [index]          # Close tab

# Checks (exit 0 = pass, exit 1 = fail)
browser exists <selector>          # Check element exists
browser count <selector>           # Count matching elements
browser visible <selector>         # Check element visibility
browser assert <js-expr>           # Assert JS expression is truthy
browser assert <js-expr> <expected>  # Assert JS result equals expected

# Accessibility
browser ax-tree                    # Print accessibility tree with refs on interactive elements
browser ax-tree --depth 3          # Limit tree depth
browser ax-tree --json             # Output as JSON (includes ref field)
browser ax-tree --with-coords      # Include CSS pixel coordinates per element
browser ax-tree --selectors        # Include CSS selectors per element
browser ax-tree --no-refs          # Suppress ref annotations
browser ax-find --name "Submit"    # Find nodes by accessible name
browser ax-find --role button      # Find nodes by ARIA role
browser ax-node <selector>         # Get accessibility info for element
browser ax-node <selector> --json  # Output as JSON

# Cloud (Browser Use API)
browser cloud login <api-key>      # Save API key
browser cloud logout               # Remove API key
browser cloud GET /browsers        # List cloud browsers
browser cloud POST /browsers '{}' # Create cloud browser
browser cloud PATCH /browsers/<id> '{"action":"stop"}'  # Stop cloud browser
browser cloud poll <task-id>       # Poll task until completion
browser cloud --help               # Show all API endpoints (from OpenAPI)
```

## Common Patterns

### Form Submission

```bash
browser start
browser open https://example.com/signup
browser ax-tree --depth 3

browser input 'input[name="name"]' 'Jane Doe'
browser input 'input[name="email"]' 'jane@example.com'
browser select '#country' 'US'
browser click 'button[type="submit"]'
browser wait-load
browser get title
```

### Coordinate-Based Interaction

When CSS selectors are unreliable (e.g., canvas apps, complex SPAs), use coordinates. Take a screenshot first to identify positions:

```bash
browser screenshot current.png    # → ~/.browser/tmp/current.png
# View the screenshot to identify coordinates

browser click 450 320
browser element-at 450 320        # Verify what's at those coords
browser input 300 200 'hello'     # Click at coords, then type
browser scroll 640 360 500        # Scroll down 500px at center
browser drag 100 100 400 400      # Drag from (100,100) to (400,400)
```

### Data Extraction

```bash
browser open https://example.com/products
browser get text 'h1'              # Get heading text
browser get attr 'a.product' 'href'  # Get link URL
browser eval 'document.querySelectorAll(".price").length'
browser eval 'JSON.stringify(Array.from(document.querySelectorAll(".item")).map(e => e.textContent))'
```

### Ref-Based Interaction (Recommended)

Use `ax-tree` to inspect the page — interactive elements get short refs you can use directly in commands:

```bash
browser open https://example.com
browser ax-tree --depth 4
# Output:
# [RootWebArea] "Example" (url=https://example.com)
#   [navigation] "Main"
#     [i1] [link] "Home" (focusable=true)
#     [i2] [link] "About" (focusable=true)
#   [main] ""
#     [i3] [heading] "Welcome" (level=1)

# Click using the ref instead of a CSS selector
browser click i2                    # Clicks the "About" link
browser get text i3                 # Gets text of the heading
```

Refs are assigned to interactive elements (links, buttons, inputs, etc.) and named content elements (headings, list items). They're refreshed every time `ax-tree` runs.

### Accessibility-Driven Interaction

Use the accessibility tree to understand page structure without relying on implementation details:

```bash
browser ax-find --role link         # Find all links
browser ax-find --name "Submit"     # Find elements named "Submit"
browser ax-node 'button#save'       # Get a11y info for specific element
```

### Tab Management

```bash
browser new-tab https://site-a.com
browser new-tab https://site-b.com
browser tabs
# * [0] about:blank - about:blank
#   [1] Site A - https://site-a.com
#   [2] Site B - https://site-b.com

browser switch 1                   # Switch to Site A
browser get title                  # "Site A"
browser close-tab 0                # Close the blank tab
```

### Assertions and Testing

```bash
# Check elements exist
browser exists 'h1'                # exit 0 if found, exit 1 if not
browser visible '.modal'           # exit 0 if visible, exit 1 if not
browser count '.list-item'         # prints count

# Assert JS expressions
browser assert 'document.title === "Expected Title"'
browser assert 'document.title' 'Expected Title'
browser assert 'document.querySelectorAll(".item").length > 0'
browser assert 'location.pathname' '/dashboard' -m 'should be on dashboard'
```

### Cloud Browser Automation

```bash
# Login with API key
browser cloud login "$BROWSER_USE_API_KEY"

# Create a cloud browser
browser cloud POST /browsers '{"headless": true}'

# Connect to a cloud browser
browser connect https://<uuid>.cdp0.browser-use.com

# Use normal commands — they work identically on local and cloud browsers
browser open https://example.com
browser screenshot page.png   # → ~/.browser/tmp/page.png
browser get title

# Stop cloud browser (browser stop auto-detects cloud URLs)
browser stop

# Run a cloud task and poll
browser cloud POST /tasks '{"url":"https://example.com","task":"Extract the main heading"}'
browser cloud poll <task-id>
```

### File Upload and Download

```bash
# Upload a file to an input
browser file 'input[type="file"]' ./document.pdf

# Upload from stdin
cat image.png | browser file 'input[type="file"]' -

# Download a linked resource
browser download 'a.download-link' ./output.pdf
browser download 'img.logo' -     # Download to stdout (pipe to file or tool)
```

### Multiple Browsers

```bash
# Start a local browser
browser start
browser open https://site-a.com
browser status
# * [0] ws://127.0.0.1:9222/... (local, alive)

# Connect to another
browser connect localhost:9333
browser status
# * [1] ws://127.0.0.1:9333/... (local, alive)
#   [0] ws://127.0.0.1:9222/... (local, alive)

# Switch between them
browser connect 0                  # Switch to first browser
browser connect 1                  # Switch to second browser
```

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success (also: `exists`/`visible`/`assert` passed) |
| 1 | Check failed (`exists`/`visible`/`assert` returned false, `ax-find` no matches) |
| 2 | Error (bad arguments, no browser, timeout, network failure) |

## JavaScript Evaluation

The `eval` command auto-wraps expressions in `() => { return (expr); }`, so you can write concise expressions:

```bash
browser eval 'document.title'                    # String printed unquoted
browser eval '2 + 2'                             # Number printed raw
browser eval 'document.querySelectorAll("a").length'
browser eval 'location.href'
browser eval 'JSON.stringify({url: location.href, title: document.title})'
```

Output formatting: strings are printed unquoted, numbers/booleans are raw, objects/arrays are pretty-printed as JSON.

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `BROWSER_HOME` | State directory (also controls output dir: `$BROWSER_HOME/tmp/`) | `~/.browser` |
| `BROWSER_TIMEOUT` | Command timeout in seconds | `30` |
| `BROWSER_CHROME_BIN` | Chrome binary path | auto-detect |
| `BROWSER_USE_API_KEY` | Cloud API key | none |

## Persistent State

The browser runs as a separate process that survives CLI exit. State is stored in `~/.browser/state.json`:

- **Browser registry**: Multiple browsers can be tracked, one is active
- **Active page**: Which tab is currently targeted
- **Data dir**: Chrome user data directory for local browsers
- **Output dir**: Screenshots, PDFs, and downloads default to `~/.browser/tmp/`

Always `browser stop` when done to clean up. Use `browser status` to check what's running.
