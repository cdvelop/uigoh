package gosite

import (
	. "github.com/cdvelop/tinystring"
)

// page represents a single HTML page and is the root of the site structure.
// It's unexported to keep the package API encapsulated.
type page struct {
	site     SiteLink
	sections []*SectionBuilder
	title    string
	filename string
	head     []string
}

// NewSection adds a new section to the page and returns a builder for it.
func (p *page) NewSection(title string) *SectionBuilder {
	section := &SectionBuilder{
		page:  p,
		site:  p.site,
		Title: title,
	}
	p.sections = append(p.sections, section)
	return section
}

// AddHead adds content to the <head> section of the page.
func (p *page) AddHead(content string) *page {
	p.head = append(p.head, content)
	return p
}

// RenderHTML generates the complete HTML for the page.
func (p *page) RenderHTML() string {
	var b = Convert()

	b.Write("<!DOCTYPE html>\n")
	b.Write("<html lang=\"es\">\n")
	b.Write("<head>\n")
	b.Write("  <meta charset=\"UTF-8\">\n")
	b.Write("  <meta name=\"viewport\" content=\"width=device-width, initial-scale-1.0\">\n")
	b.Write("  <title>")
	b.Write(Convert(p.title).EscapeHTML())
	b.Write("</title>\n")
	b.Write("  <link rel=\"stylesheet\" href=\"style.css\">\n")

	for _, h := range p.head {
		b.Write("  ")
		b.Write(h)
		b.Write("\n")
	}

	b.Write("</head>\n")
	b.Write("<body>\n")

	// Render navigation if any pages are registered
	if p.site.PageCount() > 1 {
		b.Write(p.site.BuildNav())
	}

	b.Write("  <main class=\"content\">\n")
	for _, section := range p.sections {
		b.Write(section.Render())
	}
	b.Write("  </main>\n")

	b.Write("  <script src=\"script.js\"></script>\n")
	b.Write("</body>\n")
	b.Write("</html>\n")

	return b.String()
}
