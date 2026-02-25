package tools

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

type AXNode struct {
	Role       string   `json:"role"`
	Name       string   `json:"name"`
	Value      string   `json:"value,omitempty"`
	Properties []string `json:"properties,omitempty"`
	Children   []AXNode `json:"children,omitempty"`
}

func AXTree(page *rod.Page, maxDepth int) ([]AXNode, error) {
	res, err := proto.AccessibilityGetFullAXTree{}.Call(page)
	if err != nil {
		return nil, err
	}

	// Build a map of nodeID -> node
	nodeMap := make(map[proto.AccessibilityAXNodeID]*proto.AccessibilityAXNode)
	for i := range res.Nodes {
		nodeMap[res.Nodes[i].NodeID] = res.Nodes[i]
	}

	// Build tree from root
	var roots []AXNode
	for _, n := range res.Nodes {
		if len(n.ParentID) == 0 || n.ParentID == "" {
			roots = append(roots, buildAXTree(n, nodeMap, 0, maxDepth))
		}
	}
	return roots, nil
}

func buildAXTree(n *proto.AccessibilityAXNode, nodeMap map[proto.AccessibilityAXNodeID]*proto.AccessibilityAXNode, depth, maxDepth int) AXNode {
	node := convertAXNode(n)

	if maxDepth > 0 && depth >= maxDepth {
		return node
	}

	for _, childID := range n.ChildIDs {
		if child, ok := nodeMap[childID]; ok {
			if isIgnored(child) {
				// Recurse into children of ignored nodes
				for _, grandchildID := range child.ChildIDs {
					if gc, ok := nodeMap[grandchildID]; ok {
						node.Children = append(node.Children, buildAXTree(gc, nodeMap, depth+1, maxDepth))
					}
				}
			} else {
				node.Children = append(node.Children, buildAXTree(child, nodeMap, depth+1, maxDepth))
			}
		}
	}

	return node
}

func isIgnored(n *proto.AccessibilityAXNode) bool {
	return n.Ignored
}

func convertAXNode(n *proto.AccessibilityAXNode) AXNode {
	node := AXNode{
		Role: axValueStr(n.Role),
		Name: axValueStr(n.Name),
	}
	if n.Value != nil {
		node.Value = axValueStr(n.Value)
	}
	for _, p := range n.Properties {
		if p.Value != nil {
			val := axValueStr(p.Value)
			if val != "" && val != "false" {
				node.Properties = append(node.Properties, fmt.Sprintf("%s=%s", p.Name, val))
			}
		}
	}
	return node
}

func axValueStr(v *proto.AccessibilityAXValue) string {
	if v == nil {
		return ""
	}
	// gson.JSON has a Val() method that returns the underlying value
	val := v.Value.Val()
	if val == nil {
		return ""
	}
	switch s := val.(type) {
	case string:
		return s
	default:
		data, err := json.Marshal(val)
		if err != nil {
			return fmt.Sprintf("%v", val)
		}
		return string(data)
	}
}

func AXFind(page *rod.Page, name, role string) ([]AXNode, error) {
	// Get document root node ID
	doc, err := proto.DOMGetDocument{}.Call(page)
	if err != nil {
		return nil, err
	}

	params := &proto.AccessibilityQueryAXTree{
		BackendNodeID: doc.Root.BackendNodeID,
	}
	if name != "" {
		params.AccessibleName = name
	}
	if role != "" {
		params.Role = role
	}

	res, err := params.Call(page)
	if err != nil {
		return nil, err
	}

	var nodes []AXNode
	for _, n := range res.Nodes {
		nodes = append(nodes, convertAXNode(n))
	}
	return nodes, nil
}

func AXNodeInfo(page *rod.Page, selector string) (*AXNode, error) {
	el, err := findElement(page, selector)
	if err != nil {
		return nil, err
	}

	// Get backend node ID
	desc, err := el.Describe(0, false)
	if err != nil {
		return nil, err
	}

	res, err := proto.AccessibilityGetPartialAXTree{
		BackendNodeID: desc.BackendNodeID,
	}.Call(page)
	if err != nil {
		return nil, err
	}

	// Find first non-ignored node
	for _, n := range res.Nodes {
		if !n.Ignored {
			node := convertAXNode(n)
			return &node, nil
		}
	}

	if len(res.Nodes) > 0 {
		node := convertAXNode(res.Nodes[0])
		return &node, nil
	}

	return nil, fmt.Errorf("no accessibility info for selector: %s", selector)
}

func FormatAXTree(nodes []AXNode, indent int) string {
	var b strings.Builder
	for _, n := range nodes {
		prefix := strings.Repeat("  ", indent)
		props := ""
		if len(n.Properties) > 0 {
			props = " (" + strings.Join(n.Properties, ", ") + ")"
		}
		fmt.Fprintf(&b, "%s[%s] %q%s\n", prefix, n.Role, n.Name, props)
		if len(n.Children) > 0 {
			b.WriteString(FormatAXTree(n.Children, indent+1))
		}
	}
	return b.String()
}

func FormatAXNodes(nodes []AXNode) string {
	var b strings.Builder
	for _, n := range nodes {
		props := ""
		if len(n.Properties) > 0 {
			props = " (" + strings.Join(n.Properties, ", ") + ")"
		}
		fmt.Fprintf(&b, "[%s] %q%s\n", n.Role, n.Name, props)
	}
	return b.String()
}
