package tools

import (
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
	el, err := page.Element(selector)
	if err != nil {
		return "", err
	}
	return el.HTML()
}

func Text(page *rod.Page, selector string) (string, error) {
	el, err := page.Element(selector)
	if err != nil {
		return "", err
	}
	return el.Text()
}

func Attr(page *rod.Page, selector, name string) (string, error) {
	el, err := page.Element(selector)
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
