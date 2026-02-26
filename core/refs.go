package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func refsPath(ctx *Context) string {
	return filepath.Join(ctx.StateDir, "refs.json")
}

// SaveRefMap persists a ref→selector map to refs.json.
func SaveRefMap(ctx *Context, refs map[int]string) error {
	if err := os.MkdirAll(ctx.StateDir, 0755); err != nil {
		return err
	}
	// Convert int keys to strings for JSON
	m := make(map[string]string, len(refs))
	for k, v := range refs {
		m[strconv.Itoa(k)] = v
	}
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return os.WriteFile(refsPath(ctx), data, 0644)
}

// LoadRefMap reads the ref→selector map from refs.json.
func LoadRefMap(ctx *Context) (map[int]string, error) {
	data, err := os.ReadFile(refsPath(ctx))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var m map[string]string
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	refs := make(map[int]string, len(m))
	for k, v := range m {
		n, err := strconv.Atoi(k)
		if err != nil {
			continue
		}
		refs[n] = v
	}
	return refs, nil
}

// ResolveSelector resolves a selector that may be a @ref (e.g. "@1") into
// a CSS selector by looking up refs.json. Non-ref selectors are returned as-is.
func ResolveSelector(ctx *Context, selector string) (string, error) {
	if !strings.HasPrefix(selector, "@") {
		return selector, nil
	}
	numStr := selector[1:]
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return "", fmt.Errorf("invalid ref %q: must be @<number>", selector)
	}
	refs, err := LoadRefMap(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to load refs: %w", err)
	}
	if refs == nil {
		return "", fmt.Errorf("ref %s not found (run browser ax-tree to refresh refs)", selector)
	}
	sel, ok := refs[num]
	if !ok {
		return "", fmt.Errorf("ref %s not found (run browser ax-tree to refresh refs)", selector)
	}
	return sel, nil
}
