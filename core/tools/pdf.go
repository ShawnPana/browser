package tools

import (
	"io"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

// PDF generates a PDF of the page with sensible defaults.
func PDF(page *rod.Page) (io.Reader, error) {
	return page.PDF(&proto.PagePrintToPDF{
		PrintBackground: true,
	})
}
