package cli

import (
	"browser/browser"
	"browser/browser/tools"
)

func cmdURL(args []string) {
	ctx := browser.NewContext()
	_, _, page, err := browser.WithPage(ctx)
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
	ctx := browser.NewContext()
	_, _, page, err := browser.WithPage(ctx)
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
	ctx := browser.NewContext()
	_, _, page, err := browser.WithPage(ctx)
	if err != nil {
		Fatal("%v", err)
	}
	selector := ""
	if len(args) > 0 {
		selector = args[0]
	}
	h, err := tools.HTML(page, selector)
	if err != nil {
		Fatal("%v", err)
	}
	PrintResult(h)
}

func cmdText(args []string) {
	if len(args) < 1 {
		Fatal("usage: browser text <selector>")
	}
	ctx := browser.NewContext()
	_, _, page, err := browser.WithPage(ctx)
	if err != nil {
		Fatal("%v", err)
	}
	t, err := tools.Text(page, args[0])
	if err != nil {
		Fatal("%v", err)
	}
	PrintResult(t)
}

func cmdAttr(args []string) {
	if len(args) < 2 {
		Fatal("usage: browser attr <selector> <name>")
	}
	ctx := browser.NewContext()
	_, _, page, err := browser.WithPage(ctx)
	if err != nil {
		Fatal("%v", err)
	}
	v, err := tools.Attr(page, args[0], args[1])
	if err != nil {
		Fatal("%v", err)
	}
	PrintResult(v)
}
