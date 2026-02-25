package browser

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

const CloudBaseURL = "https://api.browser-use.com/api/v2"

type CloudConfig struct {
	APIKey string `json:"api_key"`
}

type CachedBrowser struct {
	ID     string `json:"id"`
	CdpURL string `json:"cdp_url"`
	Status string `json:"status"`
}

type CloudCache struct {
	Browsers  []CachedBrowser `json:"browsers"`
	FetchedAt time.Time       `json:"fetched_at"`
}

func cloudConfigPath(ctx *Context) string {
	return filepath.Join(ctx.StateDir, "cloud.json")
}

func cloudCachePath(ctx *Context) string {
	return filepath.Join(ctx.StateDir, "cloud_cache.json")
}

func LoadCloudConfig(ctx *Context) CloudConfig {
	// Env var takes priority
	if key := os.Getenv("BROWSER_USE_API_KEY"); key != "" {
		return CloudConfig{APIKey: key}
	}

	data, err := os.ReadFile(cloudConfigPath(ctx))
	if err != nil {
		return CloudConfig{}
	}
	var cfg CloudConfig
	json.Unmarshal(data, &cfg)
	return cfg
}

func SaveCloudConfig(ctx *Context, cfg CloudConfig) error {
	if err := os.MkdirAll(ctx.StateDir, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(cloudConfigPath(ctx), data, 0644)
}

func RemoveCloudConfig(ctx *Context) error {
	return os.Remove(cloudConfigPath(ctx))
}

func LoadCloudCache(ctx *Context) *CloudCache {
	data, err := os.ReadFile(cloudCachePath(ctx))
	if err != nil {
		return nil
	}
	var cache CloudCache
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil
	}
	// Check TTL (30 seconds)
	if time.Since(cache.FetchedAt) > 30*time.Second {
		return nil
	}
	return &cache
}

func SaveCloudCache(ctx *Context, cache *CloudCache) error {
	if err := os.MkdirAll(ctx.StateDir, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(cloudCachePath(ctx), data, 0644)
}
