package tools

import (
	"fmt"
	"os"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

type ScreenshotOptions struct {
	Width    int
	Height   int
	FullPage bool
	File     string
}

func Screenshot(page *rod.Page, opts ScreenshotOptions) (string, error) {
	// Set viewport
	if err := setViewport(page, opts.Width, opts.Height); err != nil {
		return "", err
	}

	var data []byte
	var err error

	if opts.FullPage {
		data, err = page.Screenshot(true, &proto.PageCaptureScreenshot{
			Format: proto.PageCaptureScreenshotFormatPng,
		})
	} else {
		data, err = page.Screenshot(false, &proto.PageCaptureScreenshot{
			Format: proto.PageCaptureScreenshotFormatPng,
			Clip: &proto.PageViewport{
				X:      0,
				Y:      0,
				Width:  float64(opts.Width),
				Height: float64(opts.Height),
				Scale:  1,
			},
		})
	}
	if err != nil {
		return "", err
	}

	filename := opts.File
	if filename == "" {
		filename = autoScreenshotName()
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return "", err
	}

	return filename, nil
}

func setViewport(page *rod.Page, width, height int) error {
	return (proto.EmulationSetDeviceMetricsOverride{
		Width:             width,
		Height:            height,
		DeviceScaleFactor: 1,
		Mobile:            false,
	}).Call(page)
}

func autoScreenshotName() string {
	name := "screenshot.png"
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return name
	}
	for i := 2; ; i++ {
		name = fmt.Sprintf("screenshot-%d.png", i)
		if _, err := os.Stat(name); os.IsNotExist(err) {
			return name
		}
	}
}
