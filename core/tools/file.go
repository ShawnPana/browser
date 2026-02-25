package tools

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-rod/rod"
)

func FileUpload(page *rod.Page, selector string, path string) error {
	filePath := path

	// If path is "-", read stdin to a temp file
	if path == "-" {
		tmp, err := os.CreateTemp("", "browser-upload-*")
		if err != nil {
			return err
		}
		if _, err := io.Copy(tmp, os.Stdin); err != nil {
			tmp.Close()
			return err
		}
		tmp.Close()
		filePath = tmp.Name()
	}

	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return err
	}

	el, err := findElement(page, selector)
	if err != nil {
		return err
	}
	return el.SetFiles([]string{absPath})
}

// Download downloads the resource referenced by an element's href or src attribute.
// Returns the data, suggested filename, and any error.
func Download(page *rod.Page, selector string) ([]byte, string, error) {
	el, err := findElement(page, selector)
	if err != nil {
		return nil, "", err
	}

	// Get URL from href or src
	url := ""
	if href, err := el.Attribute("href"); err == nil && href != nil {
		url = *href
	}
	if url == "" {
		if src, err := el.Attribute("src"); err == nil && src != nil {
			url = *src
		}
	}
	if url == "" {
		return nil, "", fmt.Errorf("element has no href or src attribute")
	}

	// Handle data: URLs
	if strings.HasPrefix(url, "data:") {
		return decodeDataURL(url)
	}

	// Use fetch() in page context to preserve cookies/session
	js := fmt.Sprintf(`() => {
		return fetch(%q).then(r => {
			const cd = r.headers.get('content-disposition') || '';
			return r.arrayBuffer().then(buf => {
				const arr = new Uint8Array(buf);
				let binary = '';
				for (let i = 0; i < arr.length; i++) binary += String.fromCharCode(arr[i]);
				return {data: btoa(binary), filename: cd};
			});
		});
	}`, url)

	res, err := page.Eval(js)
	if err != nil {
		return nil, "", err
	}

	dataB64 := res.Value.Get("data").Str()
	filename := inferFilename(url, res.Value.Get("filename").Str())

	data, err := base64.StdEncoding.DecodeString(dataB64)
	if err != nil {
		return nil, "", err
	}

	return data, filename, nil
}

func decodeDataURL(dataURL string) ([]byte, string, error) {
	// data:[<mediatype>][;base64],<data>
	parts := strings.SplitN(dataURL, ",", 2)
	if len(parts) != 2 {
		return nil, "", fmt.Errorf("invalid data URL")
	}
	header := parts[0]
	payload := parts[1]

	if strings.Contains(header, ";base64") {
		data, err := base64.StdEncoding.DecodeString(payload)
		return data, "download", err
	}

	return []byte(payload), "download", nil
}

func inferFilename(url, contentDisposition string) string {
	// Try content-disposition
	if contentDisposition != "" {
		if idx := strings.Index(contentDisposition, "filename="); idx != -1 {
			name := contentDisposition[idx+9:]
			name = strings.Trim(name, `"' `)
			if name != "" {
				return name
			}
		}
	}

	// Extract from URL path
	path := url
	if idx := strings.Index(path, "?"); idx != -1 {
		path = path[:idx]
	}
	if idx := strings.LastIndex(path, "/"); idx != -1 {
		name := path[idx+1:]
		if name != "" {
			return name
		}
	}

	return "download"
}
