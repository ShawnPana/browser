package cli

import (
	"browser/core"
	"browser/core/tools"
	"fmt"
	"strconv"
)

func cmdTabs(args []string) {
	ctx := core.NewContext()
	s, bro, err := core.WithBrowser(ctx)
	if err != nil {
		Fatal("%v", err)
	}

	pages, err := tools.ListPages(bro)
	if err != nil {
		Fatal("%v", err)
	}

	var tabs []TabInfo
	for _, p := range pages {
		tabs = append(tabs, TabInfo{Title: p.Title, URL: p.URL})
	}
	fmt.Println(FormatTabList(tabs, s.ActivePage))
}

func cmdSwitch(args []string) {
	if len(args) < 1 {
		Fatal("usage: browser switch <index>")
	}
	idx, err := strconv.Atoi(args[0])
	if err != nil {
		Fatal("invalid index: %s", args[0])
	}

	ctx := core.NewContext()
	s, bro, err := core.WithBrowser(ctx)
	if err != nil {
		Fatal("%v", err)
	}

	if err := tools.SwitchPage(bro, idx); err != nil {
		Fatal("%v", err)
	}

	s.ActivePage = idx
	if err := core.SaveState(ctx, s); err != nil {
		Fatal("failed to save state: %v", err)
	}

	// Print info about the new active page
	pages, err := tools.ListPages(bro)
	if err != nil {
		Fatal("%v", err)
	}
	if idx < len(pages) {
		fmt.Printf("[%d] %s - %s\n", idx, pages[idx].Title, pages[idx].URL)
	}
}

func cmdNewTab(args []string) {
	url := ""
	if len(args) > 0 {
		url = args[0]
	}

	ctx := core.NewContext()
	s, bro, err := core.WithBrowser(ctx)
	if err != nil {
		Fatal("%v", err)
	}

	_, idx, err := tools.NewPage(bro, url)
	if err != nil {
		Fatal("%v", err)
	}

	s.ActivePage = idx
	if err := core.SaveState(ctx, s); err != nil {
		Fatal("failed to save state: %v", err)
	}

	fmt.Printf("[%d]\n", idx)
}

func cmdCloseTab(args []string) {
	ctx := core.NewContext()
	s, bro, err := core.WithBrowser(ctx)
	if err != nil {
		Fatal("%v", err)
	}

	idx := s.ActivePage
	if len(args) > 0 {
		var err error
		idx, err = strconv.Atoi(args[0])
		if err != nil {
			Fatal("invalid index: %s", args[0])
		}
	}

	if err := tools.ClosePage(bro, idx); err != nil {
		Fatal("%v", err)
	}

	// Adjust active page if needed
	pages, _ := tools.ListPages(bro)
	if s.ActivePage >= len(pages) {
		s.ActivePage = len(pages) - 1
	}
	if s.ActivePage < 0 {
		s.ActivePage = 0
	}

	if err := core.SaveState(ctx, s); err != nil {
		Fatal("failed to save state: %v", err)
	}

	fmt.Printf("closed [%d]\n", idx)
}
