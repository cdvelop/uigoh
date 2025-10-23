package gosite

import (
	. "github.com/cdvelop/tinystring"
)

// SectionBuilder handles the construction of a page section.
type SectionBuilder struct {
	page     *page
	site     SiteLink
	Title    string
	ModuleID string
	content  []Component
}

// Component is an interface for all UI components.
type Component interface {
	RenderHTML() string
	RenderCSS() string
	RenderJS() string
}

// Add appends a new component to the section.
func (s *SectionBuilder) Add(component Component) *SectionBuilder {
	s.content = append(s.content, component)
	s.site.AddCSS(component.RenderCSS())
	s.site.AddJS(component.RenderJS())
	return s
}

// Render generates the section's HTML.
func (s *SectionBuilder) Render() string {
	var b = Convert()
	id := s.ModuleID
	if id == "" {
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
		b.Write("    ")
		b.Write(item.RenderHTML())
		b.Write("\n")
	}
	b.Write("  </div>\n")
	b.Write("</section>\n")
	return b.String()
}
