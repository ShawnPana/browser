package tools

import (
	"encoding/json"
	"fmt"

	"github.com/go-rod/rod"
)

func URL(page *rod.Page) (string, error) {
	info, err := page.Info()
	if err != nil {
		return "", err
	}
	return info.URL, nil
}

func Title(page *rod.Page) (string, error) {
	info, err := page.Info()
	if err != nil {
		return "", err
	}
	return info.Title, nil
}

func HTML(page *rod.Page, selector string) (string, error) {
	if selector == "" {
		res, err := page.Eval(`() => document.documentElement.outerHTML`)
		if err != nil {
			return "", err
		}
		return res.Value.Str(), nil
	}
	el, err := findElement(page, selector)
	if err != nil {
		return "", err
	}
	return el.HTML()
}

func Text(page *rod.Page, selector string) (string, error) {
	el, err := findElement(page, selector)
	if err != nil {
		return "", err
	}
	return el.Text()
}

func Attr(page *rod.Page, selector, name string) (string, error) {
	el, err := findElement(page, selector)
	if err != nil {
		return "", err
	}
	val, err := el.Attribute(name)
	if err != nil {
		return "", err
	}
	if val == nil {
		return "", nil
	}
	return *val, nil
}

func Value(page *rod.Page, selector string) (string, error) {
	el, err := findElement(page, selector)
	if err != nil {
		return "", err
	}
	val, err := el.Property("value")
	if err != nil {
		return "", err
	}
	return val.Str(), nil
}

func Box(page *rod.Page, selector string) (map[string]float64, error) {
	el, err := findElement(page, selector)
	if err != nil {
		return nil, err
	}
	shape, err := el.Shape()
	if err != nil {
		return nil, err
	}
	box := shape.Box()
	return map[string]float64{
		"x":      box.X,
		"y":      box.Y,
		"width":  box.Width,
		"height": box.Height,
	}, nil
}

func Styles(page *rod.Page, selector string, props []string) (interface{}, error) {
	if len(props) > 0 {
		propsJSON, _ := json.Marshal(props)
		js := fmt.Sprintf(`() => {
			const el = document.querySelector(%q);
			if (!el) throw new Error('element not found: %s');
			const cs = getComputedStyle(el);
			const props = %s;
			const result = {};
			for (const p of props) { result[p] = cs.getPropertyValue(p); }
			return result;
		}`, selector, selector, string(propsJSON))
		res, err := page.Eval(js)
		if err != nil {
			return nil, err
		}
		return res.Value.Val(), nil
	}

	js := fmt.Sprintf(`() => {
		const el = document.querySelector(%q);
		if (!el) throw new Error('element not found: %s');
		const cs = getComputedStyle(el);
		const result = {};
		for (let i = 0; i < cs.length; i++) {
			const p = cs[i];
			result[p] = cs.getPropertyValue(p);
		}
		return result;
	}`, selector, selector)
	res, err := page.Eval(js)
	if err != nil {
		return nil, err
	}
	return res.Value.Val(), nil
}
