package tools

import (
	"encoding/json"
	"fmt"

	"github.com/go-rod/rod"
)

// JSEval evaluates a JavaScript expression and returns the result.
// The expression is auto-wrapped in an arrow function: () => { return (expr); }
func JSEval(page *rod.Page, expr string) (interface{}, error) {
	js := fmt.Sprintf(`() => { return (%s); }`, expr)
	res, err := page.Eval(js)
	if err != nil {
		return nil, err
	}
	return parseGsonResult(res.Value), nil
}

// FormatResult formats a JS evaluation result as a human-readable string.
// Strings are returned as-is, numbers use whole-number detection,
// and objects/arrays are JSON-encoded.
func FormatResult(v interface{}) string {
	switch val := v.(type) {
	case nil:
		return "null"
	case string:
		return val
	case bool:
		return fmt.Sprintf("%t", val)
	case float64:
		if val == float64(int64(val)) {
			return fmt.Sprintf("%d", int64(val))
		}
		return fmt.Sprintf("%g", val)
	case json.Number:
		return val.String()
	case int:
		return fmt.Sprintf("%d", val)
	case int64:
		return fmt.Sprintf("%d", val)
	default:
		data, err := json.Marshal(val)
		if err != nil {
			return fmt.Sprintf("%v", val)
		}
		return string(data)
	}
}

// parseGsonResult extracts a native Go type from a gson.JSON value.
func parseGsonResult(v interface{ Val() interface{} }) interface{} {
	raw := v.Val()
	if raw == nil {
		return nil
	}

	switch val := raw.(type) {
	case string:
		return val
	case bool:
		return val
	case json.Number:
		return val
	case float64:
		return val
	default:
		return raw
	}
}
