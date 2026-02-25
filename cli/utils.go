package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
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
	switch val := v.(type) {
	case nil:
		fmt.Println("null")
	case string:
		fmt.Println(val)
	case bool:
		fmt.Println(val)
	case json.Number:
		fmt.Println(val)
	case float64:
		// Print as integer if it's a whole number
		if val == float64(int64(val)) {
			fmt.Println(int64(val))
		} else {
			fmt.Println(val)
		}
	case int:
		fmt.Println(val)
	case int64:
		fmt.Println(val)
	default:
		data, err := json.MarshalIndent(val, "", "  ")
		if err != nil {
			fmt.Println(val)
		} else {
			fmt.Println(string(data))
		}
	}
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
