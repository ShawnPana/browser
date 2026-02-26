package cli

import (
	"browser/core"
	"browser/core/tools"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func cmdPDF(args []string) {
	if len(args) < 1 {
		Fatal("usage: browser pdf <path>")
	}
	ctx := core.NewContext()
	_, _, page, err := core.WithPage(ctx)
	if err != nil {
		Fatal("%v", err)
	}

	reader, err := tools.PDF(page)
	if err != nil {
		Fatal("%v", err)
	}

	outPath := ctx.ResolveOutputPath(args[0])
	os.MkdirAll(filepath.Dir(outPath), 0755)

	f, err := os.Create(outPath)
	if err != nil {
		Fatal("failed to create file: %v", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, reader); err != nil {
		Fatal("failed to write PDF: %v", err)
	}
	fmt.Println(outPath)
}
