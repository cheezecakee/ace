package ui

import (
	"strings"
)

type Render struct {
	Styles Styles

	Head   string
	Body   string
	Footer string
}

func NewRender(styles Styles) *Render {
	return &Render{
		Styles: styles,
	}
}

func (r *Render) Build() string {
	var s strings.Builder

	if r.Head != "" {
		s.WriteString(r.Styles.Header.Render(r.Head))
	}
	if r.Body != "" {
		s.WriteString(r.Styles.Body.Render(r.Body))
	}

	if r.Footer != "" {
		s.WriteString(r.Styles.Footer.Render(r.Footer))
	}
	return s.String()
}
