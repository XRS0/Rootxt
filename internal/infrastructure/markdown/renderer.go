package markdown

import (
	"bytes"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

type Renderer struct {
	engine goldmark.Markdown
}

func NewRenderer() Renderer {
	engine := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithRendererOptions(html.WithHardWraps(), html.WithUnsafe()),
	)
	return Renderer{engine: engine}
}

func (r Renderer) Render(markdown string) (string, error) {
	var buf bytes.Buffer
	if err := r.engine.Convert([]byte(markdown), &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}
