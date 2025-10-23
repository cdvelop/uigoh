
package components

import (
	"github.com/cdvelop/tinystrings"
)

// Page represents a single HTML page and is the root of the site structure.
type Page struct {
	Site     SiteLink
	sections []*SectionBuilder
	Title    string
	Filename string
	Head     []string
}

// NewSection adds a new section to the page and returns a builder for it.
func (p *Page) NewSection(title string) *SectionBuilder {
	section := &SectionBuilder{
		page:  p,
		site:  p.Site,
		Title: title,
	}
	p.sections = append(p.sections, section)
	return section
}

// AddHead adds content to the <head> section of the page.
func (p *Page) AddHead(content string) *Page {
	p.Head = append(p.Head, content)
	return p
}

// RenderHTML generates the complete HTML for the page.
func (p *Page) RenderHTML() string {
	var b tinystrings.Builder

	b.WriteString("<!DOCTYPE html>\n")
	b.WriteString("<html lang=\"es\">\n")
	b.WriteString("<head>\n")
	b.WriteString("  <meta charset=\"UTF-8\">\n")
	b.WriteString("  <meta name=\"viewport\" content=\"width=device-width, initial-scale-1.0\">\n")
	b.WriteString("  <title>")
	b.WriteString(tinystrings.EscapeHTML(p.Title))
	b.WriteString("</title>\n")
	b.WriteString("  <link rel=\"stylesheet\" href=\"style.css\">\n")

	for _, h := range p.Head {
		b.WriteString("  ")
		b.WriteString(h)
		b.WriteString("\n")
	}

	b.WriteString("</head>\n")
	b.WriteString("<body>\n")

	// Render navigation if any pages are registered
	if p.Site.PageCount() > 1 {
		b.WriteString(p.Site.BuildNav())
	}

	b.WriteString("  <main class=\"content\">\n")
	for _, section := range p.sections {
		b.WriteString(section.Render())
	}
	b.WriteString("  </main>\n")

	b.WriteString("  <script src=\"main.js\"></script>\n")
	b.WriteString("</body>\n")
	b.WriteString("</html>\n")

	return b.String()
}
