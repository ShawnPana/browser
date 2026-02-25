package tools

import (
	"time"

	"github.com/go-rod/rod"
)

func Wait(page *rod.Page, selector string) error {
	el, err := page.Element(selector)
	if err != nil {
		return err
	}
	return el.WaitVisible()
}

func WaitLoad(page *rod.Page) error {
	return page.WaitLoad()
}

func WaitStable(page *rod.Page) error {
	return page.WaitStable(300 * time.Millisecond)
}

func WaitIdle(page *rod.Page) error {
	return page.WaitIdle(5 * time.Second)
}

func Sleep(d time.Duration) {
	time.Sleep(d)
}
