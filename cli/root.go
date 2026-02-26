package cli

import (
	"fmt"
	"os"
)

var Version = "dev"

func Execute() {
	args := os.Args[1:]
	if len(args) == 0 {
		printHelp()
		return
	}

	cmd := args[0]
	rest := args[1:]

	switch cmd {
	// Browser lifecycle
	case "start":
		cmdStart(rest)
	case "stop":
		cmdStop(rest)
	case "connect":
		cmdConnect(rest)
	case "status":
		cmdStatus(rest)

	// Navigation
	case "open":
		cmdOpen(rest)
	case "back":
		cmdBack(rest)
	case "forward":
		cmdForward(rest)
	case "reload":
		cmdReload(rest)

	// Page info
	case "get":
		cmdGet(rest)

	// Interaction
	case "eval":
		cmdEval(rest)
	case "click":
		cmdClick(rest)
	case "dblclick":
		cmdDblClick(rest)
	case "rightclick":
		cmdRightClick(rest)
	case "input":
		cmdInput(rest)
	case "type":
		cmdType(rest)
	case "press":
		cmdPress(rest)
	case "clear":
		cmdClear(rest)
	case "select":
		cmdSelect(rest)
	case "submit":
		cmdSubmit(rest)
	case "hover":
		cmdHover(rest)
	case "focus":
		cmdFocus(rest)
	case "check":
		cmdCheck(rest)
	case "uncheck":
		cmdUncheck(rest)
	case "file":
		cmdFile(rest)
	case "download":
		cmdDownload(rest)
	case "scroll":
		cmdScroll(rest)
	case "scrollintoview":
		cmdScrollIntoView(rest)
	case "drag":
		cmdDrag(rest)
	case "element-at":
		cmdElementAt(rest)
	case "keyboard":
		cmdKeyboard(rest)

	// Waiting
	case "wait":
		cmdWait(rest)
	case "wait-load":
		cmdWaitLoad(rest)
	case "wait-stable":
		cmdWaitStable(rest)
	case "wait-idle":
		cmdWaitIdle(rest)
	case "sleep":
		cmdSleep(rest)

	// Output
	case "screenshot":
		cmdScreenshot(rest)
	case "pdf":
		cmdPDF(rest)

	// Tabs
	case "tabs":
		cmdTabs(rest)
	case "switch":
		cmdSwitch(rest)
	case "new-tab":
		cmdNewTab(rest)
	case "close-tab":
		cmdCloseTab(rest)

	// Checks
	case "exists":
		cmdExists(rest)
	case "count":
		cmdCount(rest)
	case "visible":
		cmdVisible(rest)
	case "assert":
		cmdAssert(rest)

	// Accessibility
	case "ax-tree":
		cmdAXTree(rest)
	case "ax-find":
		cmdAXFind(rest)
	case "ax-node":
		cmdAXNode(rest)

	// Cloud
	case "cloud":
		cmdCloud(rest)

	// Meta
	case "help", "--help", "-h":
		printHelp()
	case "version", "--version", "-v":
		fmt.Println("browser " + Version)

	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", cmd)
		fmt.Fprintf(os.Stderr, "Run 'browser help' for usage.\n")
		os.Exit(2)
	}
}

func printHelp() {
	fmt.Print(`browser - CLI for browser automation via Chrome DevTools Protocol

Usage: browser <command> [arguments]

Browser Lifecycle:
  start [--show] [-k]        Launch a new Chrome instance
  stop                       Close the active browser
  connect <url|index>        Connect to a browser or switch active
  status                     Show browser status

Navigation:
  open <url>                 Navigate to URL
  back                       Go back
  forward                    Go forward
  reload [--hard]            Reload page

Page Info:
  get url                    Print current URL
  get title                  Print page title
  get html [selector]        Print HTML (full page or element)
  get text <selector>        Print element text
  get attr <selector> <name> Print element attribute
  get value <selector>       Print input element value
  get box <selector>         Print bounding box (x, y, width, height)
  get styles <sel> [prop...] Print computed styles

Interaction:
  eval <expr>                Evaluate JavaScript
  click <sel|x y>            Click element or coordinates
  dblclick <sel|x y>         Double-click element or coordinates
  rightclick <sel|x y>       Right-click element or coordinates
  input <sel|x y> <text>     Clear and type text (replace contents)
  type <sel|x y> <text>      Type text without clearing (append)
  press <key>                Press key combo (e.g. Enter, Control+a)
  clear <sel>                Clear input field
  select <sel> <value>       Select option value
  submit <sel>               Submit form
  hover <sel|x y>            Hover over element or coordinates
  focus <sel>                Focus element
  check <sel>                Check checkbox
  uncheck <sel>              Uncheck checkbox
  file <sel> <path|->        Upload file to input
  download <sel> [file|-]    Download linked resource
  scroll <sel|x y> <delta>   Scroll element or at coordinates
  scrollintoview <sel>       Scroll element into viewport
  drag <x1 y1> <x2 y2>      Drag from point to point
  element-at <x y>           Describe element at coordinates

Keyboard:
  keyboard type <text>       Type text with real keystrokes (no selector)
  keyboard inserttext <text> Insert text without key events (no selector)

Waiting:
  wait <sel>                 Wait for element to be visible
  wait-load                  Wait for page load
  wait-stable                Wait for DOM stability
  wait-idle                  Wait for idle
  sleep <seconds>            Sleep for duration

Output:
  screenshot [-w N] [-h N] [file]  Capture screenshot (default: ~/.browser/tmp/)
  pdf <path>                       Save page as PDF (default: ~/.browser/tmp/)

Tabs:
  tabs                       List open tabs
  switch <index>             Switch to tab by index
  new-tab [url]              Open new tab
  close-tab [index]          Close tab

Checks:
  exists <sel>               Check if element exists (exit 0/1)
  count <sel>                Count matching elements
  visible <sel>              Check if element is visible (exit 0/1)
  assert <expr> [expected]   Assert JS expression (exit 0/1)

Accessibility:
  ax-tree [--depth N] [--json]       Print accessibility tree
  ax-find [--name N] [--role R]      Find accessible nodes
  ax-node <sel> [--json]             Get node accessibility info

Cloud:
  cloud login <key>          Save API key
  cloud logout               Remove API key
  cloud <METHOD> <path>      REST passthrough to Browser Use API
  cloud poll <task-id>       Poll task status
  cloud --help               Show API endpoints

Environment:
  BROWSER_HOME        State directory (default: ~/.browser); also controls output dir ($BROWSER_HOME/tmp/)
  ROD_TIMEOUT         Command timeout in seconds (default: 30)
  ROD_CHROME_BIN      Chrome binary path
  BROWSER_USE_API_KEY Cloud API key

Version: ` + Version + "\n")
}
