package cli

import (
	"github.com/ShawnPana/browser/core"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func cmdStart(args []string) {
	headless := true
	insecure := false
	for _, a := range args {
		switch a {
		case "--show":
			headless = false
		case "-k", "--insecure":
			insecure = true
		}
	}

	ctx := core.NewContext()
	debugURL, dataDir, err := core.Launch(ctx, core.LaunchOptions{
		Headless: headless,
		Insecure: insecure,
	})
	if err != nil {
		Fatal("failed to start browser: %v", err)
	}

	s, _ := core.LoadState(ctx)
	if s == nil {
		s = &core.State{}
	}
	s.DataDir = dataDir
	s.AddBrowser(debugURL, "local")
	s.ActivePage = 0

	if err := core.SaveState(ctx, s); err != nil {
		Fatal("failed to save state: %v", err)
	}

	fmt.Println(debugURL)
}

func cmdStop(args []string) {
	ctx := core.NewContext()
	s, bro, err := core.WithBrowser(ctx)
	if err != nil {
		Fatal("%v", err)
	}

	activeURL := s.Active

	// Check if this is a cloud browser and stop via API
	if isCloudURL(activeURL) {
		stopCloudBrowser(ctx, activeURL)
	}

	// Close the browser via CDP
	if err := bro.Close(); err != nil {
		// Ignore close errors — browser may already be gone
	}

	s.RemoveBrowser(activeURL)
	s.ActivePage = 0

	if len(s.Browsers) == 0 {
		core.RemoveState(ctx)
	} else {
		core.SaveState(ctx, s)
	}

	fmt.Println("stopped")
}

var cloudURLPattern = regexp.MustCompile(`wss://[a-f0-9-]+\.cdp\d*\.browser-use\.com`)

func isCloudURL(url string) bool {
	return cloudURLPattern.MatchString(url)
}

func stopCloudBrowser(ctx *core.Context, debugURL string) {
	cfg := core.LoadCloudConfig(ctx)
	if cfg.APIKey == "" {
		return
	}

	// Extract UUID from the debug URL
	// Pattern: wss://UUID.cdpN.browser-use.com/...
	parts := strings.SplitN(strings.TrimPrefix(debugURL, "wss://"), ".", 2)
	if len(parts) < 2 {
		return
	}
	browserID := parts[0]

	body := strings.NewReader(`{"action":"stop"}`)
	req, err := http.NewRequest("PATCH", core.CloudBaseURL+"/browsers/"+browserID, body)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Browser-Use-API-Key", cfg.APIKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	resp.Body.Close()
}

func cmdConnect(args []string) {
	if len(args) < 1 {
		Fatal("usage: browser connect <host:port | https://... | index>")
	}

	ctx := core.NewContext()
	s, _ := core.LoadState(ctx)
	if s == nil {
		s = &core.State{}
	}

	target := args[0]

	// Check if it's a numeric index into the registry
	if idx, ok := parseIndex(target); ok {
		if idx < 0 || idx >= len(s.Browsers) {
			Fatal("browser index %d out of range (have %d)", idx, len(s.Browsers))
		}
		s.Active = s.Browsers[idx].URL
		s.ActivePage = 0
		if err := core.SaveState(ctx, s); err != nil {
			Fatal("failed to save state: %v", err)
		}
		fmt.Println(s.Active)
		return
	}

	// Resolve debug URL from target
	debugURL, source, err := resolveDebugURL(target)
	if err != nil {
		Fatal("failed to connect: %v", err)
	}

	s.AddBrowser(debugURL, source)
	s.ActivePage = 0
	if err := core.SaveState(ctx, s); err != nil {
		Fatal("failed to save state: %v", err)
	}

	fmt.Println(debugURL)
}

func resolveDebugURL(target string) (string, string, error) {
	var versionURL string
	var source string

	if strings.HasPrefix(target, "https://") {
		versionURL = target + "/json/version"
		source = "cloud"
	} else if strings.HasPrefix(target, "http://") {
		versionURL = target + "/json/version"
		source = "local"
	} else {
		// Bare host:port
		normalized := strings.Replace(target, "localhost", "127.0.0.1", 1)
		versionURL = "http://" + normalized + "/json/version"
		source = "local"
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Get(versionURL)
	if err != nil {
		return "", "", fmt.Errorf("cannot reach %s: %w", versionURL, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	var info struct {
		WebSocketDebuggerURL string `json:"webSocketDebuggerUrl"`
	}
	if err := json.Unmarshal(body, &info); err != nil {
		return "", "", fmt.Errorf("invalid /json/version response: %w", err)
	}
	if info.WebSocketDebuggerURL == "" {
		return "", "", fmt.Errorf("no webSocketDebuggerUrl in response")
	}

	debugURL := info.WebSocketDebuggerURL
	// For cloud URLs, upgrade ws:// to wss:// if the source was HTTPS
	if source == "cloud" && strings.HasPrefix(debugURL, "ws://") {
		debugURL = "wss://" + strings.TrimPrefix(debugURL, "ws://")
	}
	// Normalize localhost
	debugURL = strings.Replace(debugURL, "localhost", "127.0.0.1", 1)

	return debugURL, source, nil
}

func parseIndex(s string) (int, bool) {
	if len(s) == 0 {
		return 0, false
	}
	n := 0
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, false
		}
		n = n*10 + int(c-'0')
	}
	return n, true
}

func cmdStatus(args []string) {
	ctx := core.NewContext()
	s, err := core.LoadState(ctx)
	if err != nil {
		Fatal("failed to load state: %v", err)
	}
	if s == nil || len(s.Browsers) == 0 {
		fmt.Println("no browsers registered")
		return
	}

	for i, b := range s.Browsers {
		marker := "  "
		if b.URL == s.Active {
			marker = "* "
		}
		alive := "dead"
		if core.IsAlive(b.URL) {
			alive = "alive"
		}
		fmt.Printf("%s[%d] %s (%s, %s)\n", marker, i, b.URL, b.Source, alive)
	}
}
