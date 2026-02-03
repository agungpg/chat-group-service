package utils

import (
	"mime"
	"strings"
)

func ExtFromContentType(contentType string) string {
	// contentType may include charset: "text/plain; charset=utf-8"
	ct := strings.TrimSpace(strings.Split(contentType, ";")[0])

	exts, err := mime.ExtensionsByType(ct)
	if err != nil || len(exts) == 0 {
		return "" // unknown
	}

	// exts usually includes leading dot, e.g. ".png"
	return exts[0]
}
