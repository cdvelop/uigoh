package gosite

import (
	. "github.com/cdvelop/tinystring"
)

// Section handles the construction of a page section.
// Its fields are unexported to maintain a controlled, fluent API.
type Section struct {
	page     *Page
	site     SiteLink
	Title    string
	ModuleID string
	content  []any
}

// Add appends a new component to the section and returns the section for chaining.
func (s *Section) Add(component any) *Section {
	s.content = append(s.content, component)

	// Cast and handle CSS if the component implements CSSRenderer.
	if cssRenderer, ok := component.(CSSRenderer); ok {
		s.site.AddCSS(cssRenderer.RenderCSS())
	}

	// Cast and handle JS if the component implements JSRenderer.
	if jsRenderer, ok := component.(JSRenderer); ok {
		s.site.AddJS(jsRenderer.RenderJS())
	}

	return s
}

// Render generates the section's HTML.
func (s *Section) Render() string {
	b := Convert()
	id := s.ModuleID
	if id == "" {
		// Generate a default ID from the title if none is provided.
		id = Convert(s.Title).ToLower().Replace(" ", "-").String()
	}
	b.Write("<section id=\"")
	b.Write(Convert(id).EscapeAttr())
	b.Write("\" class=\"page\">\n")

	if s.Title != "" {
		b.Write("  <h1>")
		b.Write(Convert(s.Title).EscapeHTML())
		b.Write("</h1>\n")
	}
	b.Write("  <div class=\"card-container\">\n")
	for _, item := range s.content {
		// Only render HTML if the component implements HTMLRenderer.
		if htmlRenderer, ok := item.(HTMLRenderer); ok {
			b.Write("    ")
			b.Write(htmlRenderer.RenderHTML())
			b.Write("\n")
		}
	}
	b.Write("  </div>\n")
	b.Write("</section>\n")
	return b.String()
}
