package tools

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

func Open(page *rod.Page, url string) error {
	if err := page.Navigate(url); err != nil {
		return err
	}
	return page.WaitLoad()
}

func Back(page *rod.Page) error {
	_, err := page.Eval(`() => history.back()`)
	if err != nil {
		return err
	}
	return page.WaitLoad()
}

func Forward(page *rod.Page) error {
	_, err := page.Eval(`() => history.forward()`)
	if err != nil {
		return err
	}
	return page.WaitLoad()
}

func Reload(page *rod.Page, hard bool) error {
	if err := (proto.PageReload{IgnoreCache: hard}).Call(page); err != nil {
		return err
	}
	return page.WaitLoad()
}
