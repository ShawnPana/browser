package tools

import (
	"fmt"

	"github.com/go-rod/rod"
)

type PageInfo struct {
	Title string
	URL   string
}

func ListPages(browser *rod.Browser) ([]PageInfo, error) {
	pages, err := browser.Pages()
	if err != nil {
		return nil, err
	}
	var result []PageInfo
	for _, p := range pages {
		info, err := p.Info()
		if err != nil {
			result = append(result, PageInfo{Title: "(error)", URL: "(error)"})
			continue
		}
		result = append(result, PageInfo{Title: info.Title, URL: info.URL})
	}
	return result, nil
}

func SwitchPage(browser *rod.Browser, index int) error {
	pages, err := browser.Pages()
	if err != nil {
		return err
	}
	if index < 0 || index >= len(pages) {
		return fmt.Errorf("page index %d out of range (have %d pages)", index, len(pages))
	}
	return nil
}

func NewPage(browser *rod.Browser, url string) (*rod.Page, int, error) {
	var page *rod.Page
	var err error

	if url != "" {
		page = browser.MustPage(url)
	} else {
		page = browser.MustPage()
	}

	// Find the new page's index
	pages, err := browser.Pages()
	if err != nil {
		return page, 0, err
	}
	for i, p := range pages {
		if p.TargetID == page.TargetID {
			return page, i, nil
		}
	}
	return page, len(pages) - 1, nil
}

func ClosePage(browser *rod.Browser, index int) error {
	pages, err := browser.Pages()
	if err != nil {
		return err
	}
	if len(pages) <= 1 {
		return fmt.Errorf("cannot close the last page")
	}
	if index < 0 || index >= len(pages) {
		return fmt.Errorf("page index %d out of range (have %d pages)", index, len(pages))
	}
	return pages[index].Close()
}
