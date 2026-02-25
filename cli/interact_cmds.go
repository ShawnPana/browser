package cli

import (
	"browser/core"
	"browser/core/tools"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func cmdEval(args []string) {
	if len(args) < 1 {
		Fatal("usage: browser eval <expr>")
	}
	ctx := core.NewContext()
	_, _, page, err := core.WithPage(ctx)
	if err != nil {
		Fatal("%v", err)
	}
	expr := strings.Join(args, " ")
	result, err := tools.JSEval(page, expr)
	if err != nil {
		Fatal("%v", err)
	}
	PrintResult(result)
}

func cmdClick(args []string) {
	if len(args) < 1 {
		Fatal("usage: browser click <selector|x y>")
	}
	ctx := core.NewContext()
	_, _, page, err := core.WithPage(ctx)
	if err != nil {
		Fatal("%v", err)
	}
	if tools.IsCoordinate(args) {
		x, y, err := tools.ParseCoords(args[0], args[1])
		if err != nil {
			Fatal("%v", err)
		}
		if err := tools.ClickCoords(page, x, y); err != nil {
			Fatal("%v", err)
		}
	} else {
		if err := tools.ClickSelector(page, args[0]); err != nil {
			Fatal("%v", err)
		}
	}
}

func cmdInput(args []string) {
	if len(args) < 2 {
		Fatal("usage: browser input <selector|x y> <text>")
	}
	ctx := core.NewContext()
	_, _, page, err := core.WithPage(ctx)
	if err != nil {
		Fatal("%v", err)
	}
	if tools.IsCoordinate(args) {
		if len(args) < 3 {
			Fatal("usage: browser input <x> <y> <text>")
		}
		x, y, err := tools.ParseCoords(args[0], args[1])
		if err != nil {
			Fatal("%v", err)
		}
		text := strings.Join(args[2:], " ")
		if err := tools.InputCoordsImpl(page, x, y, text); err != nil {
			Fatal("%v", err)
		}
	} else {
		text := strings.Join(args[1:], " ")
		if err := tools.InputSelector(page, args[0], text); err != nil {
			Fatal("%v", err)
		}
	}
}

func cmdClear(args []string) {
	if len(args) < 1 {
		Fatal("usage: browser clear <selector>")
	}
	ctx := core.NewContext()
	_, _, page, err := core.WithPage(ctx)
	if err != nil {
		Fatal("%v", err)
	}
	if err := tools.Clear(page, args[0]); err != nil {
		Fatal("%v", err)
	}
}

func cmdSelect(args []string) {
	if len(args) < 2 {
		Fatal("usage: browser select <selector> <value>")
	}
	ctx := core.NewContext()
	_, _, page, err := core.WithPage(ctx)
	if err != nil {
		Fatal("%v", err)
	}
	if err := tools.Select(page, args[0], args[1]); err != nil {
		Fatal("%v", err)
	}
}

func cmdSubmit(args []string) {
	if len(args) < 1 {
		Fatal("usage: browser submit <selector>")
	}
	ctx := core.NewContext()
	_, _, page, err := core.WithPage(ctx)
	if err != nil {
		Fatal("%v", err)
	}
	if err := tools.Submit(page, args[0]); err != nil {
		Fatal("%v", err)
	}
}

func cmdHover(args []string) {
	if len(args) < 1 {
		Fatal("usage: browser hover <selector|x y>")
	}
	ctx := core.NewContext()
	_, _, page, err := core.WithPage(ctx)
	if err != nil {
		Fatal("%v", err)
	}
	if tools.IsCoordinate(args) {
		x, y, err := tools.ParseCoords(args[0], args[1])
		if err != nil {
			Fatal("%v", err)
		}
		if err := tools.HoverCoords(page, x, y); err != nil {
			Fatal("%v", err)
		}
	} else {
		if err := tools.HoverSelector(page, args[0]); err != nil {
			Fatal("%v", err)
		}
	}
}

func cmdFocus(args []string) {
	if len(args) < 1 {
		Fatal("usage: browser focus <selector>")
	}
	ctx := core.NewContext()
	_, _, page, err := core.WithPage(ctx)
	if err != nil {
		Fatal("%v", err)
	}
	if err := tools.Focus(page, args[0]); err != nil {
		Fatal("%v", err)
	}
}

func cmdFile(args []string) {
	if len(args) < 2 {
		Fatal("usage: browser file <selector> <path|->")
	}
	ctx := core.NewContext()
	_, _, page, err := core.WithPage(ctx)
	if err != nil {
		Fatal("%v", err)
	}
	if err := tools.FileUpload(page, args[0], args[1]); err != nil {
		Fatal("%v", err)
	}
}

func cmdDownload(args []string) {
	if len(args) < 1 {
		Fatal("usage: browser download <selector> [file|-]")
	}
	ctx := core.NewContext()
	_, _, page, err := core.WithPage(ctx)
	if err != nil {
		Fatal("%v", err)
	}

	data, filename, err := tools.Download(page, args[0])
	if err != nil {
		Fatal("%v", err)
	}

	outFile := ""
	if len(args) > 1 {
		outFile = args[1]
	}

	if outFile == "-" {
		os.Stdout.Write(data)
		return
	}

	if outFile == "" {
		outFile = filename
	}

	if err := os.WriteFile(outFile, data, 0644); err != nil {
		Fatal("failed to write file: %v", err)
	}
	fmt.Println(outFile)
}

func cmdScroll(args []string) {
	if len(args) < 2 {
		Fatal("usage: browser scroll <selector|x y> <delta>")
	}
	ctx := core.NewContext()
	_, _, page, err := core.WithPage(ctx)
	if err != nil {
		Fatal("%v", err)
	}

	if tools.IsCoordinate(args) {
		if len(args) < 3 {
			Fatal("usage: browser scroll <x> <y> <delta>")
		}
		x, y, err := tools.ParseCoords(args[0], args[1])
		if err != nil {
			Fatal("%v", err)
		}
		delta, err := strconv.ParseFloat(args[2], 64)
		if err != nil {
			Fatal("invalid delta: %s", args[2])
		}
		if err := tools.ScrollCoords(page, x, y, delta); err != nil {
			Fatal("%v", err)
		}
	} else {
		delta, err := strconv.ParseFloat(args[len(args)-1], 64)
		if err != nil {
			Fatal("invalid delta: %s", args[len(args)-1])
		}
		if err := tools.ScrollSelector(page, args[0], delta); err != nil {
			Fatal("%v", err)
		}
	}
}

func cmdDrag(args []string) {
	if len(args) < 4 {
		Fatal("usage: browser drag <x1> <y1> <x2> <y2>")
	}
	ctx := core.NewContext()
	_, _, page, err := core.WithPage(ctx)
	if err != nil {
		Fatal("%v", err)
	}

	x1, y1, err := tools.ParseCoords(args[0], args[1])
	if err != nil {
		Fatal("%v", err)
	}
	x2, y2, err := tools.ParseCoords(args[2], args[3])
	if err != nil {
		Fatal("%v", err)
	}

	steps := 10
	if err := tools.Drag(page, x1, y1, x2, y2, steps); err != nil {
		Fatal("%v", err)
	}
}

func cmdType(args []string) {
	if len(args) < 2 {
		Fatal("usage: browser type <selector|x y> <text>")
	}
	ctx := core.NewContext()
	_, _, page, err := core.WithPage(ctx)
	if err != nil {
		Fatal("%v", err)
	}
	if tools.IsCoordinate(args) {
		if len(args) < 3 {
			Fatal("usage: browser type <x> <y> <text>")
		}
		x, y, err := tools.ParseCoords(args[0], args[1])
		if err != nil {
			Fatal("%v", err)
		}
		text := strings.Join(args[2:], " ")
		if err := tools.TypeCoordsImpl(page, x, y, text); err != nil {
			Fatal("%v", err)
		}
	} else {
		text := strings.Join(args[1:], " ")
		if err := tools.TypeSelector(page, args[0], text); err != nil {
			Fatal("%v", err)
		}
	}
}

func cmdDblClick(args []string) {
	if len(args) < 1 {
		Fatal("usage: browser dblclick <selector|x y>")
	}
	ctx := core.NewContext()
	_, _, page, err := core.WithPage(ctx)
	if err != nil {
		Fatal("%v", err)
	}
	if tools.IsCoordinate(args) {
		x, y, err := tools.ParseCoords(args[0], args[1])
		if err != nil {
			Fatal("%v", err)
		}
		if err := tools.DblClickCoords(page, x, y); err != nil {
			Fatal("%v", err)
		}
	} else {
		if err := tools.DblClickSelector(page, args[0]); err != nil {
			Fatal("%v", err)
		}
	}
}

func cmdRightClick(args []string) {
	if len(args) < 1 {
		Fatal("usage: browser rightclick <selector|x y>")
	}
	ctx := core.NewContext()
	_, _, page, err := core.WithPage(ctx)
	if err != nil {
		Fatal("%v", err)
	}
	if tools.IsCoordinate(args) {
		x, y, err := tools.ParseCoords(args[0], args[1])
		if err != nil {
			Fatal("%v", err)
		}
		if err := tools.RightClickCoords(page, x, y); err != nil {
			Fatal("%v", err)
		}
	} else {
		if err := tools.RightClickSelector(page, args[0]); err != nil {
			Fatal("%v", err)
		}
	}
}

func cmdPress(args []string) {
	if len(args) < 1 {
		Fatal("usage: browser press <key>")
	}
	ctx := core.NewContext()
	_, _, page, err := core.WithPage(ctx)
	if err != nil {
		Fatal("%v", err)
	}
	combo := strings.Join(args, "+")
	if err := tools.Press(page, combo); err != nil {
		Fatal("%v", err)
	}
}

func cmdCheck(args []string) {
	if len(args) < 1 {
		Fatal("usage: browser check <selector>")
	}
	ctx := core.NewContext()
	_, _, page, err := core.WithPage(ctx)
	if err != nil {
		Fatal("%v", err)
	}
	if err := tools.Check(page, args[0]); err != nil {
		Fatal("%v", err)
	}
}

func cmdUncheck(args []string) {
	if len(args) < 1 {
		Fatal("usage: browser uncheck <selector>")
	}
	ctx := core.NewContext()
	_, _, page, err := core.WithPage(ctx)
	if err != nil {
		Fatal("%v", err)
	}
	if err := tools.Uncheck(page, args[0]); err != nil {
		Fatal("%v", err)
	}
}

func cmdScrollIntoView(args []string) {
	if len(args) < 1 {
		Fatal("usage: browser scrollintoview <selector>")
	}
	ctx := core.NewContext()
	_, _, page, err := core.WithPage(ctx)
	if err != nil {
		Fatal("%v", err)
	}
	if err := tools.ScrollIntoView(page, args[0]); err != nil {
		Fatal("%v", err)
	}
}

func cmdElementAt(args []string) {
	if len(args) < 2 {
		Fatal("usage: browser element-at <x> <y>")
	}
	ctx := core.NewContext()
	_, _, page, err := core.WithPage(ctx)
	if err != nil {
		Fatal("%v", err)
	}

	x, err := strconv.Atoi(args[0])
	if err != nil {
		Fatal("invalid x coordinate: %s", args[0])
	}
	y, err := strconv.Atoi(args[1])
	if err != nil {
		Fatal("invalid y coordinate: %s", args[1])
	}

	desc, err := tools.ElementAt(page, x, y)
	if err != nil {
		Fatal("%v", err)
	}

	node := desc.Node
	// Format: <tag attr1="val1" attr2="val2"> Backend Node ID: N
	var attrs []string
	for i := 0; i+1 < len(node.Attributes); i += 2 {
		attrs = append(attrs, fmt.Sprintf(`%s="%s"`, node.Attributes[i], node.Attributes[i+1]))
	}
	tag := node.LocalName
	if len(attrs) > 0 {
		fmt.Printf("<%s %s>\n", tag, strings.Join(attrs, " "))
	} else {
		fmt.Printf("<%s>\n", tag)
	}
	fmt.Printf("Backend Node ID: %d\n", node.BackendNodeID)
}
