package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/ShawnPana/browser/core"
	"github.com/ShawnPana/browser/core/tools"
)

// Fatal prints an error to stderr and exits with code 2.
func Fatal(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "error: "+format+"\n", args...)
	os.Exit(2)
}

// CheckFail prints a failure message to stderr and exits with code 1.
func CheckFail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "fail: "+format+"\n", args...)
	os.Exit(1)
}

// PrintResult prints a value to stdout. Strings are printed unquoted,
// objects/arrays are pretty-printed as JSON, booleans/numbers are raw.
func PrintResult(v interface{}) {
	s := tools.FormatResult(v)
	// Pretty-print if it looks like JSON object/array
	if len(s) > 0 && (s[0] == '{' || s[0] == '[') {
		var raw json.RawMessage
		if json.Unmarshal([]byte(s), &raw) == nil {
			if out, err := json.MarshalIndent(raw, "", "  "); err == nil {
				fmt.Println(string(out))
				return
			}
		}
	}
	fmt.Println(s)
}

// PrintJSON pretty-prints a JSON-encoded byte slice.
func PrintJSON(data []byte) {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		fmt.Println(string(data))
		return
	}
	out, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Println(string(data))
		return
	}
	fmt.Println(string(out))
}

// FormatTabList formats a list of pages for the `pages` command.
func FormatTabList(pages []TabInfo, activeIndex int) string {
	var b strings.Builder
	for i, p := range pages {
		marker := "  "
		if i == activeIndex {
			marker = "* "
		}
		fmt.Fprintf(&b, "%s[%d] %s - %s\n", marker, i, p.Title, p.URL)
	}
	return strings.TrimRight(b.String(), "\n")
}

// TabInfo holds basic page info for listing.
type TabInfo struct {
	Title string
	URL   string
}

// resolveSelector resolves a @ref selector to a CSS selector, or returns
// the selector unchanged if it's not a ref.
func resolveSelector(ctx *core.Context, sel string) string {
	resolved, err := core.ResolveSelector(ctx, sel)
	if err != nil {
		Fatal("%v", err)
	}
	return resolved
}
