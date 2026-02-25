package cli

import (
	"browser/browser"
	"browser/browser/tools"
	"fmt"
	"strconv"
)

func cmdScreenshot(args []string) {
	width := 1280
	height := 720
	fullPage := true
	file := ""

	i := 0
	for i < len(args) {
		switch args[i] {
		case "-w":
			i++
			if i >= len(args) {
				Fatal("-w requires a value")
			}
			w, err := strconv.Atoi(args[i])
			if err != nil {
				Fatal("invalid width: %s", args[i])
			}
			width = w
		case "-h":
			i++
			if i >= len(args) {
				Fatal("-h requires a value")
			}
			h, err := strconv.Atoi(args[i])
			if err != nil {
				Fatal("invalid height: %s", args[i])
			}
			height = h
			fullPage = false
		default:
			if file == "" {
				file = args[i]
			}
		}
		i++
	}

	ctx := browser.NewContext()
	_, _, page, err := browser.WithPage(ctx)
	if err != nil {
		Fatal("%v", err)
	}

	filename, err := tools.Screenshot(page, tools.ScreenshotOptions{
		Width:    width,
		Height:   height,
		FullPage: fullPage,
		File:     file,
	})
	if err != nil {
		Fatal("%v", err)
	}

	fmt.Printf("%s (%dx%d viewport)\n", filename, width, height)
}
