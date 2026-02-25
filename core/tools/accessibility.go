package tools

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

type AXBox struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

type AXNode struct {
	Role       string   `json:"role"`
	Name       string   `json:"name"`
	Value      string   `json:"value,omitempty"`
	Properties []string `json:"properties,omitempty"`
	Box        *AXBox   `json:"box,omitempty"`
	Selector   string   `json:"selector,omitempty"`
	Children   []AXNode `json:"children,omitempty"`
}

func AXTree(page *rod.Page, maxDepth int, withCoords bool, withSelectors bool) ([]AXNode, error) {
	return axTreeForFrame(page, "", maxDepth, 0, withCoords, withSelectors)
}

func axTreeForFrame(page *rod.Page, frameID proto.PageFrameID, maxDepth, startDepth int, withCoords, withSelectors bool) ([]AXNode, error) {
	req := proto.AccessibilityGetFullAXTree{}
	if frameID != "" {
		req.FrameID = frameID
	}
	res, err := req.Call(page)
	if err != nil {
		return nil, err
	}

	nodeMap := make(map[proto.AccessibilityAXNodeID]*proto.AccessibilityAXNode)
	for i := range res.Nodes {
		nodeMap[res.Nodes[i].NodeID] = res.Nodes[i]
	}

	var boxMap map[proto.AccessibilityAXNodeID]*AXBox
	if withCoords {
		boxMap = resolveBoxes(page, res.Nodes)
	}

	var selMap map[proto.AccessibilityAXNodeID]string
	if withSelectors {
		selMap = resolveSelectors(page, res.Nodes)
	}

	var roots []AXNode
	for _, n := range res.Nodes {
		if len(n.ParentID) == 0 || n.ParentID == "" {
			roots = append(roots, buildAXNode(page, n, nodeMap, boxMap, selMap, startDepth, maxDepth, withCoords, withSelectors))
		}
	}
	return roots, nil
}

func resolveBoxes(page *rod.Page, nodes []*proto.AccessibilityAXNode) map[proto.AccessibilityAXNodeID]*AXBox {
	boxMap := make(map[proto.AccessibilityAXNodeID]*AXBox)
	for _, n := range nodes {
		if n.Ignored || n.BackendDOMNodeID == 0 {
			continue
		}
		res, err := proto.DOMGetBoxModel{BackendNodeID: n.BackendDOMNodeID}.Call(page)
		if err != nil {
			continue // hidden or off-screen elements
		}
		q := res.Model.Content
		if len(q) < 8 {
			continue
		}
		// Content quad: [x1,y1, x2,y2, x3,y3, x4,y4]
		x := q[0]
		y := q[1]
		w := q[2] - q[0]
		h := q[5] - q[1]
		boxMap[n.NodeID] = &AXBox{X: x, Y: y, Width: w, Height: h}
	}
	return boxMap
}

func buildAXNode(page *rod.Page, n *proto.AccessibilityAXNode, nodeMap map[proto.AccessibilityAXNodeID]*proto.AccessibilityAXNode, boxMap map[proto.AccessibilityAXNodeID]*AXBox, selMap map[proto.AccessibilityAXNodeID]string, depth, maxDepth int, withCoords, withSelectors bool) AXNode {
	node := convertAXNode(n)
	if boxMap != nil {
		node.Box = boxMap[n.NodeID]
	}
	if selMap != nil {
		node.Selector = selMap[n.NodeID]
	}

	if maxDepth > 0 && depth >= maxDepth {
		return node
	}

	// Traverse into iframe child frames
	role := axValueStr(n.Role)
	if role == "Iframe" && n.BackendDOMNodeID != 0 {
		desc, err := proto.DOMDescribeNode{BackendNodeID: n.BackendDOMNodeID}.Call(page)
		if err == nil && desc.Node.FrameID != "" {
			children, err := axTreeForFrame(page, desc.Node.FrameID, maxDepth, depth+1, withCoords, withSelectors)
			if err == nil {
				node.Children = children
			}
		}
		return node
	}

	for _, childID := range n.ChildIDs {
		if child, ok := nodeMap[childID]; ok {
			if isIgnored(child) {
				for _, grandchildID := range child.ChildIDs {
					if gc, ok := nodeMap[grandchildID]; ok {
						node.Children = append(node.Children, buildAXNode(page, gc, nodeMap, boxMap, selMap, depth+1, maxDepth, withCoords, withSelectors))
					}
				}
			} else {
				node.Children = append(node.Children, buildAXNode(page, child, nodeMap, boxMap, selMap, depth+1, maxDepth, withCoords, withSelectors))
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

func resolveSelectors(page *rod.Page, nodes []*proto.AccessibilityAXNode) map[proto.AccessibilityAXNodeID]string {
	selMap := make(map[proto.AccessibilityAXNodeID]string)
	for _, n := range nodes {
		if n.Ignored || n.BackendDOMNodeID == 0 {
			continue
		}
		desc, err := proto.DOMDescribeNode{BackendNodeID: n.BackendDOMNodeID}.Call(page)
		if err != nil {
			continue
		}
		sel := buildSelector(desc.Node)
		if sel != "" {
			selMap[n.NodeID] = sel
		}
	}
	return selMap
}

func buildSelector(node *proto.DOMNode) string {
	attrs := parseAttributes(node.Attributes)
	tag := node.LocalName

	// #id
	if id, ok := attrs["id"]; ok && id != "" {
		return "#" + id
	}

	// tag[name="value"]
	if name, ok := attrs["name"]; ok && name != "" {
		return tag + `[name="` + name + `"]`
	}

	// tag.class1.class2
	if class, ok := attrs["class"]; ok && class != "" {
		classes := strings.Fields(class)
		return tag + "." + strings.Join(classes, ".")
	}

	// input[type="value"]
	if typ, ok := attrs["type"]; ok && typ != "" && tag == "input" {
		return tag + `[type="` + typ + `"]`
	}

	return ""
}

func parseAttributes(attrs []string) map[string]string {
	m := make(map[string]string)
	for i := 0; i+1 < len(attrs); i += 2 {
		m[attrs[i]] = attrs[i+1]
	}
	return m
}

func FormatAXTree(nodes []AXNode, indent int) string {
	var b strings.Builder
	for _, n := range nodes {
		prefix := strings.Repeat("  ", indent)
		props := ""
		if len(n.Properties) > 0 {
			props = " (" + strings.Join(n.Properties, ", ") + ")"
		}
		coords := ""
		if n.Box != nil {
			cx := int(math.Round(n.Box.X + n.Box.Width/2))
			cy := int(math.Round(n.Box.Y + n.Box.Height/2))
			coords = fmt.Sprintf(" @%d,%d", cx, cy)
		}
		sel := ""
		if n.Selector != "" {
			sel = fmt.Sprintf(" {%s}", n.Selector)
		}
		fmt.Fprintf(&b, "%s[%s] %q%s%s%s\n", prefix, n.Role, n.Name, props, coords, sel)
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
