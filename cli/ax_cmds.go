package cli

import (
	"browser/browser"
	"browser/browser/tools"
	"encoding/json"
	"fmt"
	"strconv"
)

func cmdAXTree(args []string) {
	depth := 0
	asJSON := false

	i := 0
	for i < len(args) {
		switch args[i] {
		case "--depth":
			i++
			if i >= len(args) {
				Fatal("--depth requires a value")
			}
			d, err := strconv.Atoi(args[i])
			if err != nil {
				Fatal("invalid depth: %s", args[i])
			}
			depth = d
		case "--json":
			asJSON = true
		}
		i++
	}

	ctx := browser.NewContext()
	_, _, page, err := browser.WithPage(ctx)
	if err != nil {
		Fatal("%v", err)
	}

	nodes, err := tools.AXTree(page, depth)
	if err != nil {
		Fatal("%v", err)
	}

	if asJSON {
		data, err := json.MarshalIndent(nodes, "", "  ")
		if err != nil {
			Fatal("%v", err)
		}
		fmt.Println(string(data))
	} else {
		fmt.Print(tools.FormatAXTree(nodes, 0))
	}
}

func cmdAXFind(args []string) {
	name := ""
	role := ""
	asJSON := false

	i := 0
	for i < len(args) {
		switch args[i] {
		case "--name":
			i++
			if i >= len(args) {
				Fatal("--name requires a value")
			}
			name = args[i]
		case "--role":
			i++
			if i >= len(args) {
				Fatal("--role requires a value")
			}
			role = args[i]
		case "--json":
			asJSON = true
		}
		i++
	}

	ctx := browser.NewContext()
	_, _, page, err := browser.WithPage(ctx)
	if err != nil {
		Fatal("%v", err)
	}

	nodes, err := tools.AXFind(page, name, role)
	if err != nil {
		Fatal("%v", err)
	}

	if len(nodes) == 0 {
		CheckFail("no matching nodes")
	}

	if asJSON {
		data, err := json.MarshalIndent(nodes, "", "  ")
		if err != nil {
			Fatal("%v", err)
		}
		fmt.Println(string(data))
	} else {
		fmt.Print(tools.FormatAXNodes(nodes))
	}
}

func cmdAXNode(args []string) {
	if len(args) < 1 {
		Fatal("usage: browser ax-node <selector> [--json]")
	}

	selector := args[0]
	asJSON := false
	for _, a := range args[1:] {
		if a == "--json" {
			asJSON = true
		}
	}

	ctx := browser.NewContext()
	_, _, page, err := browser.WithPage(ctx)
	if err != nil {
		Fatal("%v", err)
	}

	node, err := tools.AXNodeInfo(page, selector)
	if err != nil {
		Fatal("%v", err)
	}

	if asJSON {
		data, err := json.MarshalIndent(node, "", "  ")
		if err != nil {
			Fatal("%v", err)
		}
		fmt.Println(string(data))
	} else {
		props := ""
		if len(node.Properties) > 0 {
			props = " (" + joinStrings(node.Properties, ", ") + ")"
		}
		fmt.Printf("[%s] %q%s\n", node.Role, node.Name, props)
	}
}

func joinStrings(ss []string, sep string) string {
	result := ""
	for i, s := range ss {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}
