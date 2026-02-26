package cli

import (
	"fmt"

	"github.com/ShawnPana/browser/core"
	"github.com/ShawnPana/browser/core/tools"
)

func cmdExists(args []string) {
	if len(args) < 1 {
		Fatal("usage: browser exists <selector>")
	}
	ctx := core.NewContext()
	sel := resolveSelector(ctx, args[0])
	_, _, page, err := core.WithPage(ctx)
	if err != nil {
		Fatal("%v", err)
	}
	exists, err := tools.Exists(page, sel)
	if err != nil {
		Fatal("%v", err)
	}
	fmt.Println(exists)
	if !exists {
		CheckFail("element not found: %s", args[0])
	}
}

func cmdCount(args []string) {
	if len(args) < 1 {
		Fatal("usage: browser count <selector>")
	}
	ctx := core.NewContext()
	sel := resolveSelector(ctx, args[0])
	_, _, page, err := core.WithPage(ctx)
	if err != nil {
		Fatal("%v", err)
	}
	count, err := tools.Count(page, sel)
	if err != nil {
		Fatal("%v", err)
	}
	fmt.Println(count)
}

func cmdVisible(args []string) {
	if len(args) < 1 {
		Fatal("usage: browser visible <selector>")
	}
	ctx := core.NewContext()
	sel := resolveSelector(ctx, args[0])
	_, _, page, err := core.WithPage(ctx)
	if err != nil {
		Fatal("%v", err)
	}
	vis, err := tools.Visible(page, sel)
	if err != nil {
		Fatal("%v", err)
	}
	fmt.Println(vis)
	if !vis {
		CheckFail("element not visible: %s", args[0])
	}
}

func cmdAssert(args []string) {
	if len(args) < 1 {
		Fatal("usage: browser assert <expression> [expected] [-m message]")
	}

	expr := args[0]
	expected := ""
	message := ""

	i := 1
	for i < len(args) {
		if args[i] == "-m" && i+1 < len(args) {
			message = args[i+1]
			i += 2
		} else if expected == "" {
			expected = args[i]
			i++
		} else {
			i++
		}
	}

	ctx := core.NewContext()
	_, _, page, err := core.WithPage(ctx)
	if err != nil {
		Fatal("%v", err)
	}

	pass, result, err := tools.Assert(page, expr, expected)
	if err != nil {
		Fatal("%v", err)
	}

	if pass {
		fmt.Println("pass")
	} else {
		if message != "" {
			CheckFail("%s (got: %s)", message, result)
		} else if expected != "" {
			CheckFail("expected %q, got %q", expected, result)
		} else {
			CheckFail("expression is falsy (got: %s)", result)
		}
	}
}
