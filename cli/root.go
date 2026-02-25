package cli

import (
	"fmt"
	"os"
)

const Version = "0.1.0"

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
	case "url":
		cmdURL(rest)
	case "title":
		cmdTitle(rest)
	case "html":
		cmdHTML(rest)
	case "text":
		cmdText(rest)
	case "attr":
		cmdAttr(rest)

	// Interaction
	case "js":
		cmdJS(rest)
	case "click":
		cmdClick(rest)
	case "input":
		cmdInput(rest)
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
	case "file":
		cmdFile(rest)
	case "download":
		cmdDownload(rest)
	case "scroll":
		cmdScroll(rest)
	case "drag":
		cmdDrag(rest)
	case "element-at":
		cmdElementAt(rest)

	// Waiting
	case "wait":
		cmdWait(rest)
	case "waitload":
		cmdWaitLoad(rest)
	case "waitstable":
		cmdWaitStable(rest)
	case "waitidle":
		cmdWaitIdle(rest)
	case "sleep":
		cmdSleep(rest)

	// Screenshot
	case "screenshot":
		cmdScreenshot(rest)

	// Tabs
	case "pages":
		cmdPages(rest)
	case "page":
		cmdPage(rest)
	case "newpage":
		cmdNewPage(rest)
	case "closepage":
		cmdClosePage(rest)

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
  url                        Print current URL
  title                      Print page title
  html [selector]            Print HTML (full page or element)
  text <selector>            Print element text
  attr <selector> <name>     Print element attribute

Interaction:
  js <expr>                  Evaluate JavaScript
  click <sel|x y>            Click element or coordinates
  input <sel|x y> <text>     Type text into element or at coordinates
  clear <sel>                Clear input field
  select <sel> <value>       Select option value
  submit <sel>               Submit form
  hover <sel|x y>            Hover over element or coordinates
  focus <sel>                Focus element
  file <sel> <path|->        Upload file to input
  download <sel> [file|-]    Download linked resource
  scroll <sel|x y> <delta>   Scroll element or at coordinates
  drag <x1 y1> <x2 y2>      Drag from point to point
  element-at <x y>           Describe element at coordinates

Waiting:
  wait <sel>                 Wait for element to be visible
  waitload                   Wait for page load
  waitstable                 Wait for DOM stability
  waitidle                   Wait for idle
  sleep <seconds>            Sleep for duration

Screenshots:
  screenshot [-w N] [-h N] [file]  Capture screenshot

Tabs:
  pages                      List open pages
  page <index>               Switch to page by index
  newpage [url]              Open new page
  closepage [index]          Close page

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
  BROWSER_HOME        State directory (default: ~/.browser)
  ROD_TIMEOUT         Command timeout in seconds (default: 30)
  ROD_CHROME_BIN      Chrome binary path
  BROWSER_USE_API_KEY Cloud API key

Version: ` + Version + "\n")
}
