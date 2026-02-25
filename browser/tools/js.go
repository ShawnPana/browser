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
