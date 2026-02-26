package cli

import (
	"browser/core"
	"browser/core/tools"
	"fmt"
)

func cmdScreenshot(args []string) {
	fullPage := false
	file := ""

	for _, arg := range args {
		switch arg {
		case "--full":
			fullPage = true
		default:
			if file == "" {
				file = arg
			}
		}
	}

	ctx := core.NewContext()

	if file == "" {
		file = tools.AutoScreenshotName(ctx.OutputDir())
	} else {
		file = ctx.ResolveOutputPath(file)
	}

	_, _, page, err := core.WithPage(ctx)
	if err != nil {
		Fatal("%v", err)
	}

	filename, err := tools.Screenshot(page, tools.ScreenshotOptions{
		FullPage: fullPage,
		File:     file,
	})
	if err != nil {
		Fatal("%v", err)
	}

	fmt.Println(filename)
}
