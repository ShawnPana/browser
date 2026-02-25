package tools

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/proto"
)

// IsCoordinate checks if the first two args look like x,y coordinates.
func IsCoordinate(args []string) bool {
	if len(args) < 2 {
		return false
	}
	_, err1 := strconv.ParseFloat(args[0], 64)
	_, err2 := strconv.ParseFloat(args[1], 64)
	return err1 == nil && err2 == nil
}

// ParseCoords parses x,y from two string args.
func ParseCoords(xStr, yStr string) (float64, float64, error) {
	x, err := strconv.ParseFloat(xStr, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid x coordinate: %s", xStr)
	}
	y, err := strconv.ParseFloat(yStr, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid y coordinate: %s", yStr)
	}
	return x, y, nil
}

func ClickSelector(page *rod.Page, selector string) error {
	el, err := page.Element(selector)
	if err != nil {
		return err
	}
	if err := el.Click(proto.InputMouseButtonLeft, 1); err != nil {
		return err
	}
	time.Sleep(100 * time.Millisecond)
	return nil
}

func ClickCoords(page *rod.Page, x, y float64) error {
	if err := page.Mouse.MoveTo(proto.Point{X: x, Y: y}); err != nil {
		return err
	}
	if err := page.Mouse.Click(proto.InputMouseButtonLeft, 1); err != nil {
		return err
	}
	time.Sleep(100 * time.Millisecond)
	return nil
}

func InputSelector(page *rod.Page, selector, text string) error {
	el, err := page.Element(selector)
	if err != nil {
		return err
	}
	el.MustSelectAllText().MustInput(text)
	return nil
}

func InputCoordsImpl(page *rod.Page, x, y float64, text string) error {
	if err := ClickCoords(page, x, y); err != nil {
		return err
	}
	// Select all (Ctrl+A) then type
	_ = page.KeyActions().Press(input.ControlLeft).Type(input.KeyA).Do()
	return page.InsertText(text)
}

func Clear(page *rod.Page, selector string) error {
	el, err := page.Element(selector)
	if err != nil {
		return err
	}
	el.MustSelectAllText().MustInput("")
	return nil
}

func Select(page *rod.Page, selector, value string) error {
	js := fmt.Sprintf(`() => {
		const el = document.querySelector(%q);
		if (!el) throw new Error('element not found: %s');
		el.value = %q;
		el.dispatchEvent(new Event('change', {bubbles: true}));
	}`, selector, selector, value)
	_, err := page.Eval(js)
	return err
}

func Submit(page *rod.Page, selector string) error {
	js := fmt.Sprintf(`() => {
		const el = document.querySelector(%q);
		if (!el) throw new Error('element not found: %s');
		if (el.form) {
			el.form.submit();
		} else if (el.submit) {
			el.submit();
		} else {
			throw new Error('element has no submit method');
		}
	}`, selector, selector)
	_, err := page.Eval(js)
	return err
}

func HoverSelector(page *rod.Page, selector string) error {
	el, err := page.Element(selector)
	if err != nil {
		return err
	}
	return el.Hover()
}

func HoverCoords(page *rod.Page, x, y float64) error {
	return page.Mouse.MoveTo(proto.Point{X: x, Y: y})
}

func Focus(page *rod.Page, selector string) error {
	el, err := page.Element(selector)
	if err != nil {
		return err
	}
	return el.Focus()
}

func ScrollSelector(page *rod.Page, selector string, deltaY float64) error {
	// Use JS to get element position — avoids WaitInteractable timeout on complex pages
	js := fmt.Sprintf(`() => {
		const el = document.querySelector(%q);
		if (!el) throw new Error('element not found: %s');
		const r = el.getBoundingClientRect();
		return {x: r.x + r.width/2, y: r.y + r.height/2};
	}`, selector, selector)
	res, err := page.Eval(js)
	if err != nil {
		return err
	}
	x := res.Value.Get("x").Num()
	y := res.Value.Get("y").Num()
	page.Mouse.MoveTo(proto.Point{X: x, Y: y})
	return page.Mouse.Scroll(0, deltaY, 1)
}

func ScrollCoords(page *rod.Page, x, y, deltaY float64) error {
	if err := page.Mouse.MoveTo(proto.Point{X: x, Y: y}); err != nil {
		return err
	}
	return page.Mouse.Scroll(0, deltaY, 1)
}

func Drag(page *rod.Page, x1, y1, x2, y2 float64, steps int) error {
	if steps <= 0 {
		steps = 10
	}

	if err := page.Mouse.MoveTo(proto.Point{X: x1, Y: y1}); err != nil {
		return err
	}
	if err := page.Mouse.Down(proto.InputMouseButtonLeft, 1); err != nil {
		return err
	}

	for i := 1; i <= steps; i++ {
		t := float64(i) / float64(steps)
		px := x1 + (x2-x1)*t
		py := y1 + (y2-y1)*t
		if err := page.Mouse.MoveTo(proto.Point{X: px, Y: py}); err != nil {
			return err
		}
		time.Sleep(10 * time.Millisecond)
	}

	return page.Mouse.Up(proto.InputMouseButtonLeft, 1)
}

func ElementAt(page *rod.Page, x, y int) (*proto.DOMDescribeNodeResult, error) {
	nodeRes, err := proto.DOMGetNodeForLocation{X: x, Y: y}.Call(page)
	if err != nil {
		return nil, err
	}

	desc, err := proto.DOMDescribeNode{
		BackendNodeID: nodeRes.BackendNodeID,
	}.Call(page)
	if err != nil {
		return nil, err
	}

	return desc, nil
}
