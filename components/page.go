
package components

import (
	"strings"
)

// Page represents a single HTML page and is the root of the site structure.
type Page struct {
	site     *site
	sections []*SectionBuilder
	Title    string
	Filename string
	Head     []string
}

// NewPage creates a new site and returns the root page.
func NewPage(cfg *Config) *Page {
	s := newSite(cfg)
	p := &Page{
		site:     s,
		Title:    cfg.Title,
		Filename: "index.html", // Default filename for the root page
	}
	s.AddPage(p) // Register the root page
	return p
}

// NewSection adds a new section to the page and returns a builder for it.
func (p *Page) NewSection(title string) *SectionBuilder {
	section := &SectionBuilder{
		page:  p,
		site:  p.site,
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

// Generate builds the entire site, including all pages, CSS, and JS.
func (p *Page) Generate() error {
	return p.site.Generate()
}

// RenderHTML generates the complete HTML for the page.
func (p *Page) RenderHTML() string {
	var b strings.Builder

	b.WriteString("<!DOCTYPE html>\n")
	b.WriteString("<html lang=\"es\">\n")
	b.WriteString("<head>\n")
	b.WriteString("  <meta charset=\"UTF-8\">\n")
	b.WriteString("  <meta name=\"viewport\" content=\"width=device-width, initial-scale-1.0\">\n")
	b.WriteString("  <title>")
	b.WriteString(escapeHTML(p.Title))
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
	if len(p.site.pages) > 1 {
		b.WriteString(p.site.buildCombinedNav())
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
