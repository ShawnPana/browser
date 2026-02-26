package cli

import (
	"encoding/json"
	"fmt"

	"github.com/ShawnPana/browser/core"
	"github.com/ShawnPana/browser/core/tools"
)

func cmdGet(args []string) {
	if len(args) < 1 {
		Fatal("usage: browser get <url|title|html|text|attr|value|box|styles> ...")
	}
	sub := args[0]
	subArgs := args[1:]
	switch sub {
	case "url":
		cmdURL(subArgs)
	case "title":
		cmdTitle(subArgs)
	case "html":
		cmdHTML(subArgs)
	case "text":
		cmdText(subArgs)
	case "attr":
		cmdAttr(subArgs)
	case "value":
		cmdValue(subArgs)
	case "box":
		cmdBox(subArgs)
	case "styles":
		cmdStyles(subArgs)
	default:
		Fatal("unknown get subcommand: %s", sub)
	}
}

func cmdURL(args []string) {
	ctx := core.NewContext()
	_, _, page, err := core.WithPage(ctx)
	if err != nil {
		Fatal("%v", err)
	}
	u, err := tools.URL(page)
	if err != nil {
		Fatal("%v", err)
	}
	PrintResult(u)
}

func cmdTitle(args []string) {
	ctx := core.NewContext()
	_, _, page, err := core.WithPage(ctx)
	if err != nil {
		Fatal("%v", err)
	}
	t, err := tools.Title(page)
	if err != nil {
		Fatal("%v", err)
	}
	PrintResult(t)
}

func cmdHTML(args []string) {
	ctx := core.NewContext()
	_, _, page, err := core.WithPage(ctx)
	if err != nil {
		Fatal("%v", err)
	}
	selector := ""
	if len(args) > 0 {
		selector = resolveSelector(ctx, args[0])
	}
	h, err := tools.HTML(page, selector)
	if err != nil {
		Fatal("%v", err)
	}
	PrintResult(h)
}

func cmdText(args []string) {
	if len(args) < 1 {
		Fatal("usage: browser get text <selector>")
	}
	ctx := core.NewContext()
	_, _, page, err := core.WithPage(ctx)
	if err != nil {
		Fatal("%v", err)
	}
	sel := resolveSelector(ctx, args[0])
	t, err := tools.Text(page, sel)
	if err != nil {
		Fatal("%v", err)
	}
	PrintResult(t)
}

func cmdAttr(args []string) {
	if len(args) < 2 {
		Fatal("usage: browser get attr <selector> <name>")
	}
	ctx := core.NewContext()
	_, _, page, err := core.WithPage(ctx)
	if err != nil {
		Fatal("%v", err)
	}
	sel := resolveSelector(ctx, args[0])
	v, err := tools.Attr(page, sel, args[1])
	if err != nil {
		Fatal("%v", err)
	}
	PrintResult(v)
}

func cmdValue(args []string) {
	if len(args) < 1 {
		Fatal("usage: browser get value <selector>")
	}
	ctx := core.NewContext()
	_, _, page, err := core.WithPage(ctx)
	if err != nil {
		Fatal("%v", err)
	}
	sel := resolveSelector(ctx, args[0])
	v, err := tools.Value(page, sel)
	if err != nil {
		Fatal("%v", err)
	}
	PrintResult(v)
}

func cmdBox(args []string) {
	if len(args) < 1 {
		Fatal("usage: browser get box <selector>")
	}
	ctx := core.NewContext()
	_, _, page, err := core.WithPage(ctx)
	if err != nil {
		Fatal("%v", err)
	}
	sel := resolveSelector(ctx, args[0])
	box, err := tools.Box(page, sel)
	if err != nil {
		Fatal("%v", err)
	}
	out, _ := json.MarshalIndent(box, "", "  ")
	fmt.Println(string(out))
}

func cmdStyles(args []string) {
	if len(args) < 1 {
		Fatal("usage: browser get styles <selector> [prop...]")
	}
	ctx := core.NewContext()
	_, _, page, err := core.WithPage(ctx)
	if err != nil {
		Fatal("%v", err)
	}
	sel := resolveSelector(ctx, args[0])
	props := args[1:]
	result, err := tools.Styles(page, sel, props)
	if err != nil {
		Fatal("%v", err)
	}
	out, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(out))
}
