package cli

import (
	"github.com/ShawnPana/browser/core"
	"github.com/ShawnPana/browser/core/tools"
)

func cmdOpen(args []string) {
	if len(args) < 1 {
		Fatal("usage: browser open <url>")
	}
	ctx := core.NewContext()
	_, _, page, err := core.WithPage(ctx)
	if err != nil {
		Fatal("%v", err)
	}
	if err := tools.Open(page, args[0]); err != nil {
		Fatal("%v", err)
	}
}

func cmdBack(args []string) {
	ctx := core.NewContext()
	_, _, page, err := core.WithPage(ctx)
	if err != nil {
		Fatal("%v", err)
	}
	if err := tools.Back(page); err != nil {
		Fatal("%v", err)
	}
}

func cmdForward(args []string) {
	ctx := core.NewContext()
	_, _, page, err := core.WithPage(ctx)
	if err != nil {
		Fatal("%v", err)
	}
	if err := tools.Forward(page); err != nil {
		Fatal("%v", err)
	}
}

func cmdReload(args []string) {
	hard := false
	for _, a := range args {
		if a == "--hard" {
			hard = true
		}
	}
	ctx := core.NewContext()
	_, _, page, err := core.WithPage(ctx)
	if err != nil {
		Fatal("%v", err)
	}
	if err := tools.Reload(page, hard); err != nil {
		Fatal("%v", err)
	}
}
