package cli

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/ShawnPana/browser/core"
	"github.com/ShawnPana/browser/core/tools"
)

func cmdAXTree(args []string) {
	depth := 0
	asJSON := false
	withCoords := false
	withSelectors := false
	noRefs := false

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
		case "--with-coords":
			withCoords = true
		case "--selectors":
			withSelectors = true
		case "--no-refs":
			noRefs = true
		}
		i++
	}

	ctx := core.NewContext()
	_, _, page, err := core.WithPage(ctx)
	if err != nil {
		Fatal("%v", err)
	}

	// Always resolve selectors internally (needed for refs).
	// The --selectors flag controls whether selectors are displayed.
	nodes, err := tools.AXTree(page, depth, withCoords, true)
	if err != nil {
		Fatal("%v", err)
	}

	// Assign refs and persist them (unless --no-refs)
	if !noRefs {
		refs := tools.AssignRefs(nodes)
		if err := core.SaveRefMap(ctx, refs); err != nil {
			Fatal("failed to save refs: %v", err)
		}
	}

	// Strip selectors from output if --selectors was not requested
	if !withSelectors {
		stripSelectors(nodes)
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

// stripSelectors removes the Selector field from all nodes so it doesn't
// appear in output (text or JSON) when --selectors is not requested.
func stripSelectors(nodes []AXNode) {
	for i := range nodes {
		nodes[i].Selector = ""
		if len(nodes[i].Children) > 0 {
			stripSelectors(nodes[i].Children)
		}
	}
}

type AXNode = tools.AXNode

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

	ctx := core.NewContext()
	_, _, page, err := core.WithPage(ctx)
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

	ctx := core.NewContext()
	selector := resolveSelector(ctx, args[0])
	asJSON := false
	for _, a := range args[1:] {
		if a == "--json" {
			asJSON = true
		}
	}

	_, _, page, err := core.WithPage(ctx)
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
