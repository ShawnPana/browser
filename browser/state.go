package browser

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type BrowserEntry struct {
	URL    string `json:"url"`    // ws:// or wss:// debug URL
	Source string `json:"source"` // "local" or "cloud"
}

type State struct {
	Active     string         `json:"active"`      // debug URL of active browser
	Browsers   []BrowserEntry `json:"browsers"`    // all known browsers
	ActivePage int            `json:"active_page"` // active tab index on active browser
	DataDir    string         `json:"data_dir"`    // Chrome user data dir (for local)
}

func statePath(ctx *Context) string {
	return filepath.Join(ctx.StateDir, "state.json")
}

func LoadState(ctx *Context) (*State, error) {
	data, err := os.ReadFile(statePath(ctx))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var s State
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

func SaveState(ctx *Context, s *State) error {
	if err := os.MkdirAll(ctx.StateDir, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(statePath(ctx), data, 0644)
}

func RemoveState(ctx *Context) error {
	return os.Remove(statePath(ctx))
}

// AddBrowser adds a browser to the registry and sets it as active.
func (s *State) AddBrowser(url, source string) {
	for _, b := range s.Browsers {
		if b.URL == url {
			s.Active = url
			return
		}
	}
	s.Browsers = append(s.Browsers, BrowserEntry{URL: url, Source: source})
	s.Active = url
}

// RemoveBrowser removes a browser from the registry.
// If it was active, switches to the first remaining browser or clears active.
func (s *State) RemoveBrowser(url string) {
	for i, b := range s.Browsers {
		if b.URL == url {
			s.Browsers = append(s.Browsers[:i], s.Browsers[i+1:]...)
			break
		}
	}
	if s.Active == url {
		if len(s.Browsers) > 0 {
			s.Active = s.Browsers[0].URL
		} else {
			s.Active = ""
		}
	}
}

// ActiveEntry returns the BrowserEntry for the active browser, or nil.
func (s *State) ActiveEntry() *BrowserEntry {
	for i := range s.Browsers {
		if s.Browsers[i].URL == s.Active {
			return &s.Browsers[i]
		}
	}
	return nil
}
