package core

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

// WithPage loads state, connects to the active browser, and returns the active page.
// The caller should NOT close the browser — it persists across invocations.
func WithPage(ctx *Context) (*State, *rod.Browser, *rod.Page, error) {
	s, err := LoadState(ctx)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to load state: %w", err)
	}
	if s == nil || s.Active == "" {
		return nil, nil, nil, fmt.Errorf("no browser running (use 'browser start' or 'browser connect')")
	}

	browser := rod.New().ControlURL(s.Active).NoDefaultDevice()
	if err := browser.Connect(); err != nil {
		return nil, nil, nil, fmt.Errorf("failed to connect to browser: %w", err)
	}

	pages, err := browser.Pages()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to list pages: %w", err)
	}
	if len(pages) == 0 {
		// Create a blank page if none exist (e.g., Chrome started with no-startup-window)
		page := browser.MustPage()
		pages = rod.Pages{page}
	}

	idx := s.ActivePage
	if idx < 0 || idx >= len(pages) {
		idx = 0
		s.ActivePage = 0
	}

	page := pages[idx].Timeout(ctx.Timeout)

	// Auto-dismiss JS dialogs (alert/confirm/prompt/beforeunload) so commands
	// don't hang waiting for user interaction that can't happen in a CLI.
	go page.EachEvent(func(e *proto.PageJavascriptDialogOpening) {
		shouldAccept := e.Type != proto.PageDialogTypePrompt
		_ = proto.PageHandleJavaScriptDialog{
			Accept:     shouldAccept,
			PromptText: e.DefaultPrompt,
		}.Call(page)
		fmt.Fprintf(os.Stderr, "[%s] %s\n", e.Type, e.Message)
	})()

	return s, browser, page, nil
}

// WithBrowser loads state and connects to the active browser without selecting a page.
func WithBrowser(ctx *Context) (*State, *rod.Browser, error) {
	s, err := LoadState(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load state: %w", err)
	}
	if s == nil || s.Active == "" {
		return nil, nil, fmt.Errorf("no browser running (use 'browser start' or 'browser connect')")
	}

	browser := rod.New().ControlURL(s.Active).NoDefaultDevice()
	if err := browser.Connect(); err != nil {
		return nil, nil, fmt.Errorf("failed to connect to browser: %w", err)
	}

	return s, browser, nil
}

// IsAlive checks if a browser at the given debug URL is alive by hitting /json/version.
func IsAlive(debugURL string) bool {
	base := debugURLToHTTP(debugURL)
	if base == "" {
		return false
	}

	client := &http.Client{
		Timeout: 2 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := client.Get(base + "/json/version")
	if err != nil {
		return false
	}
	resp.Body.Close()
	return resp.StatusCode == 200
}

// debugURLToHTTP converts a ws(s):// debug URL to an http(s):// base URL.
func debugURLToHTTP(wsURL string) string {
	if strings.HasPrefix(wsURL, "wss://") {
		host := strings.TrimPrefix(wsURL, "wss://")
		if idx := strings.Index(host, "/"); idx != -1 {
			host = host[:idx]
		}
		return "https://" + host
	}
	if strings.HasPrefix(wsURL, "ws://") {
		host := strings.TrimPrefix(wsURL, "ws://")
		if idx := strings.Index(host, "/"); idx != -1 {
			host = host[:idx]
		}
		return "http://" + host
	}
	return ""
}
