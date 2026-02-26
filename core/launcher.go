package core

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/go-rod/rod/lib/launcher"
)

type LaunchOptions struct {
	Headless bool
	Insecure bool
}

// Launch starts a new Chrome instance and returns the debug WebSocket URL and data dir.
func Launch(ctx *Context, opts LaunchOptions) (debugURL string, dataDir string, err error) {
	dataDir = ctx.StateDir + "/chrome-data"

	l := launcher.New().
		Set("no-sandbox").
		Set("disable-gpu").
		Leakless(false).
		UserDataDir(dataDir).
		Headless(opts.Headless)

	if !opts.Headless {
		l = l.Delete("no-startup-window")
	}

	if opts.Insecure {
		l = l.Set("ignore-certificate-errors")
	}

	bin, err := findChrome()
	if err != nil {
		return "", "", err
	}
	l = l.Bin(bin)

	u, err := l.Launch()
	if err != nil {
		return "", "", err
	}
	return u, dataDir, nil
}

// findChrome resolves the Chrome/Chromium binary path.
// Resolution order:
// 1. BROWSER_CHROME_BIN env var
// 2. exec.LookPath for common binary names
// 3. Error with instructions
func findChrome() (string, error) {
	if bin := os.Getenv("BROWSER_CHROME_BIN"); bin != "" {
		return bin, nil
	}

	for _, name := range []string{"chromium", "chromium-browser", "google-chrome", "google-chrome-stable"} {
		if path, err := exec.LookPath(name); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("Chrome not found. Install Chrome or Chromium, or set BROWSER_CHROME_BIN.")
}
