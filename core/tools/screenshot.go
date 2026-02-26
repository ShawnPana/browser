package tools

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

type ScreenshotOptions struct {
	FullPage bool
	File     string
}

func Screenshot(page *rod.Page, opts ScreenshotOptions) (string, error) {
	metrics, err := proto.PageGetLayoutMetrics{}.Call(page)
	if err != nil {
		return "", err
	}

	var w, h float64
	if opts.FullPage {
		w = metrics.CSSContentSize.Width
		h = metrics.CSSContentSize.Height
	} else {
		w = float64(metrics.CSSLayoutViewport.ClientWidth)
		h = float64(metrics.CSSLayoutViewport.ClientHeight)
	}

	req := &proto.PageCaptureScreenshot{
		Format: proto.PageCaptureScreenshotFormatPng,
		Clip: &proto.PageViewport{
			X:      0,
			Y:      0,
			Width:  w,
			Height: h,
			Scale:  1,
		},
	}
	shot, err := req.Call(page)
	if err != nil {
		return "", err
	}
	data := shot.Data

	filename := opts.File
	if filename == "" {
		filename = AutoScreenshotName(".")
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return "", err
	}

	return filename, nil
}

func AutoScreenshotName(dir string) string {
	os.MkdirAll(dir, 0755)
	name := filepath.Join(dir, "screenshot.png")
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return name
	}
	for i := 2; ; i++ {
		name = filepath.Join(dir, fmt.Sprintf("screenshot-%d.png", i))
		if _, err := os.Stat(name); os.IsNotExist(err) {
			return name
		}
	}
}
