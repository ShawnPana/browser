package core

import (
	"os"

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

	if bin := os.Getenv("ROD_CHROME_BIN"); bin != "" {
		l = l.Bin(bin)
	}

	u, err := l.Launch()
	if err != nil {
		return "", "", err
	}
	return u, dataDir, nil
}
