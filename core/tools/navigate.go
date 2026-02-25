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
	res, err := proto.PageGetNavigationHistory{}.Call(page)
	if err != nil || res.CurrentIndex <= 0 {
		_, err = page.Eval(`() => history.back()`)
		return err
	}
	return proto.PageNavigateToHistoryEntry{
		EntryID: res.Entries[res.CurrentIndex-1].ID,
	}.Call(page)
}

func Forward(page *rod.Page) error {
	res, err := proto.PageGetNavigationHistory{}.Call(page)
	if err != nil || res.CurrentIndex >= len(res.Entries)-1 {
		_, err = page.Eval(`() => history.forward()`)
		return err
	}
	return proto.PageNavigateToHistoryEntry{
		EntryID: res.Entries[res.CurrentIndex+1].ID,
	}.Call(page)
}

func Reload(page *rod.Page, hard bool) error {
	if err := (proto.PageReload{IgnoreCache: hard}).Call(page); err != nil {
		return err
	}
	return page.WaitLoad()
}
