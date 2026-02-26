package cli

import (
	"github.com/ShawnPana/browser/core"
	"github.com/ShawnPana/browser/core/tools"
	"strings"
)

func cmdKeyboard(args []string) {
	if len(args) < 1 {
		Fatal("usage: browser keyboard <type|inserttext> <text>")
	}
	sub := args[0]
	subArgs := args[1:]
	switch sub {
	case "type":
		cmdKeyboardType(subArgs)
	case "inserttext":
		cmdKeyboardInsertText(subArgs)
	default:
		Fatal("unknown keyboard subcommand: %s", sub)
	}
}

func cmdKeyboardType(args []string) {
	if len(args) < 1 {
		Fatal("usage: browser keyboard type <text>")
	}
	ctx := core.NewContext()
	_, _, page, err := core.WithPage(ctx)
	if err != nil {
		Fatal("%v", err)
	}
	text := strings.Join(args, " ")
	if err := tools.KeyboardType(page, text); err != nil {
		Fatal("%v", err)
	}
}

func cmdKeyboardInsertText(args []string) {
	if len(args) < 1 {
		Fatal("usage: browser keyboard inserttext <text>")
	}
	ctx := core.NewContext()
	_, _, page, err := core.WithPage(ctx)
	if err != nil {
		Fatal("%v", err)
	}
	text := strings.Join(args, " ")
	if err := tools.KeyboardInsertText(page, text); err != nil {
		Fatal("%v", err)
	}
}
