package core

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const DefaultTimeout = 30 * time.Second

type Context struct {
	StateDir string
	Timeout  time.Duration
}

func NewContext() *Context {
	return &Context{
		StateDir: stateDir(),
		Timeout:  timeout(),
	}
}

func (c *Context) OutputDir() string {
	return filepath.Join(c.StateDir, "tmp")
}

func (c *Context) ResolveOutputPath(name string) string {
	if strings.ContainsRune(name, filepath.Separator) || strings.ContainsRune(name, '/') {
		return name
	}
	return filepath.Join(c.OutputDir(), name)
}

func stateDir() string {
	if v := os.Getenv("BROWSER_HOME"); v != "" {
		return v
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ".browser"
	}
	return filepath.Join(home, ".browser")
}

func timeout() time.Duration {
	if v := os.Getenv("BROWSER_TIMEOUT"); v != "" {
		if secs, err := strconv.ParseFloat(v, 64); err == nil {
			return time.Duration(secs * float64(time.Second))
		}
	}
	return DefaultTimeout
}
