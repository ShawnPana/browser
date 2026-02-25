package cli

import (
	"browser/core"
	"browser/core/tools"
	"strconv"
	"time"
)

func cmdWait(args []string) {
	if len(args) < 1 {
		Fatal("usage: browser wait <selector>")
	}
	ctx := core.NewContext()
	_, _, page, err := core.WithPage(ctx)
	if err != nil {
		Fatal("%v", err)
	}
	if err := tools.Wait(page, args[0]); err != nil {
		Fatal("%v", err)
	}
}

func cmdWaitLoad(args []string) {
	ctx := core.NewContext()
	_, _, page, err := core.WithPage(ctx)
	if err != nil {
		Fatal("%v", err)
	}
	if err := tools.WaitLoad(page); err != nil {
		Fatal("%v", err)
	}
}

func cmdWaitStable(args []string) {
	ctx := core.NewContext()
	_, _, page, err := core.WithPage(ctx)
	if err != nil {
		Fatal("%v", err)
	}
	if err := tools.WaitStable(page); err != nil {
		Fatal("%v", err)
	}
}

func cmdWaitIdle(args []string) {
	ctx := core.NewContext()
	_, _, page, err := core.WithPage(ctx)
	if err != nil {
		Fatal("%v", err)
	}
	if err := tools.WaitIdle(page); err != nil {
		Fatal("%v", err)
	}
}

func cmdSleep(args []string) {
	if len(args) < 1 {
		Fatal("usage: browser sleep <seconds>")
	}
	secs, err := strconv.ParseFloat(args[0], 64)
	if err != nil {
		Fatal("invalid duration: %s", args[0])
	}
	tools.Sleep(time.Duration(secs * float64(time.Second)))
}
