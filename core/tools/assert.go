package tools

import (
	"fmt"
	"strings"

	"github.com/go-rod/rod"
)

func Exists(page *rod.Page, selector string) (bool, error) {
	has, _, err := page.Has(selector)
	if err != nil {
		return false, err
	}
	return has, nil
}

func Count(page *rod.Page, selector string) (int, error) {
	els, err := page.Elements(selector)
	if err != nil {
		return 0, err
	}
	return len(els), nil
}

func Visible(page *rod.Page, selector string) (bool, error) {
	has, el, err := page.Has(selector)
	if err != nil {
		return false, err
	}
	if !has {
		return false, nil
	}
	vis, err := el.Visible()
	if err != nil {
		return false, err
	}
	return vis, nil
}

// Assert evaluates a JS expression.
// In truthy mode (expected==""), checks if result is truthy.
// In equality mode, compares string representation to expected.
func Assert(page *rod.Page, expr, expected string) (bool, string, error) {
	result, err := JSEval(page, expr)
	if err != nil {
		return false, "", err
	}

	resultStr := fmt.Sprintf("%v", result)

	if expected == "" {
		// Truthy mode
		if isFalsy(result) {
			return false, resultStr, nil
		}
		return true, resultStr, nil
	}

	// Equality mode
	if strings.TrimSpace(resultStr) == strings.TrimSpace(expected) {
		return true, resultStr, nil
	}
	return false, resultStr, nil
}

func isFalsy(v interface{}) bool {
	if v == nil {
		return true
	}
	switch val := v.(type) {
	case bool:
		return !val
	case string:
		return val == ""
	case float64:
		return val == 0
	case int:
		return val == 0
	}
	return false
}
