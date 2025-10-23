
package components

import (
	"github.com/cdvelop/tinystrings"
)

// SectionBuilder handles the construction of a page section.
type SectionBuilder struct {
	page     *Page
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
	var b tinystrings.Builder
	id := s.ModuleID
	if id == "" {
		id = tinystrings.ToLower(tinystrings.ReplaceAll(s.Title, " ", "-"))
	}
	b.WriteString("<section id=\"")
	b.WriteString(tinystrings.EscapeAttr(id))
	b.WriteString("\" class=\"page\">\n")

	if s.Title != "" {
		b.WriteString("  <h1>")
		b.WriteString(tinystrings.EscapeHTML(s.Title))
		b.WriteString("</h1>\n")
	}
	b.WriteString("  <div class=\"card-container\">\n")
	for _, item := range s.content {
		b.WriteString("    ")
		b.WriteString(item.RenderHTML())
		b.WriteString("\n")
	}
	b.WriteString("  </div>\n")
	b.WriteString("</section>\n")
	return b.String()
}
