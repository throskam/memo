package lib

import (
	"bytes"

	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

var (
	markdownParser goldmark.Markdown
	sanitizer      *bluemonday.Policy
)

// RenderMarkdown renders markdown content to HTML
func RenderMarkdown(content string) string {
	var buf bytes.Buffer

	err := markdownParser.Convert([]byte(content), &buf)
	if err != nil {
		return ""
	}

	return sanitizer.Sanitize(buf.String())
}

func init() {
	markdownParser = goldmark.New(
		goldmark.WithExtensions(
			extension.Table,
			extension.Strikethrough,
			extension.TaskList,
			extension.Linkify,
			extension.GFM,
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
		),
	)

	sanitizer = bluemonday.UGCPolicy()
}
