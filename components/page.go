
package components

import "strings"

type Page struct {
    title      string
    filename   string            // e.g., "services.html", "contact.html"
    sections   []*SectionBuilder // Sections within this page
    rawContent []string          // For SPA sections
    navigation string            // Pre-rendered navigation (set by Site)
    site       *Site             // Back-reference to site
    head       []string          // Additional <head> content
}

func NewPage(title string) *Page {
    return &Page{
        title:      title,
        sections:   make([]*SectionBuilder, 0),
        rawContent: make([]string, 0),
        head:       make([]string, 0),
    }
}

// Section creates a section within THIS page
// Automatically adds section to page - NO need for AddSection()
func (p *Page) Section(title string) *SectionBuilder {
    section := p.site.Section(title)
    section.page = p

    // Auto-add section to page
    p.sections = append(p.sections, section)

    // Auto-accumulate section CSS to site
    p.site.AddCSS(section.RenderCSS())

    return section
}

// AddRawSection adds raw HTML section content to page
// Used by Site.AddSection() for SPA sections
func (p *Page) AddRawSection(html string) {
    p.rawContent = append(p.rawContent, html)
}

// AddHead adds content to <head> section
func (p *Page) AddHead(content string) *Page {
    p.head = append(p.head, content)
    return p
}

// RenderHTML generates the complete HTML page
func (p *Page) RenderHTML() string {
    var b strings.Builder

    b.WriteString("<!DOCTYPE html>\n")
    b.WriteString("<html lang=\"es\">\n")
    b.WriteString("<head>\n")
    b.WriteString("  <meta charset=\"UTF-8\">\n")
    b.WriteString("  <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n")
    b.WriteString("  <title>")
    b.WriteString(escapeHTML(p.title))
    b.WriteString("</title>\n")
    b.WriteString("  <link rel=\"stylesheet\" href=\"style.css\">\n")

    // Additional head content
    for _, h := range p.head {
        b.WriteString("  ")
        b.WriteString(h)
        b.WriteString("\n")
    }

    b.WriteString("</head>\n")
    b.WriteString("<body>\n")

    // Navigation (set by Site.buildCombinedNav)
    if p.navigation != "" {
        b.WriteString(p.navigation)
    }

    b.WriteString("  <main class=\"content\">\n")

    // Render all sections from builders
    for _, section := range p.sections {
        b.WriteString(section.Render())
    }

    // Render raw HTML sections for SPA
    for _, raw := range p.rawContent {
        b.WriteString(raw)
    }

    b.WriteString("  </main>\n")

    b.WriteString("  <script src=\"main.js\"></script>\n")
    b.WriteString("</body>\n")
    b.WriteString("</html>\n")

    return b.String()
}
