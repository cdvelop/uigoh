package gosite

import (
	"fmt"

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
	// Build head entries
	headHTML := ""
	for _, h := range p.head {
		headHTML += "  " + h + "\n"
	}

	// Build sections
	sectionsHTML := ""
	for _, section := range p.sections {
		sectionsHTML += section.Render()
	}

	title := Convert(p.title).EscapeHTML()

	tpl := `<!DOCTYPE html>
<html lang="es">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>%s</title>
  <link rel="stylesheet" href="style.css">
%s</head>
<body>
%s  <main class="content">
%s  </main>

  <script src="script.js"></script>
</body>
</html>
`

	// Optionally include nav if multiple pages exist
	navHTML := ""
	if p.site.PageCount() > 1 {
		navHTML = p.site.BuildNav()
	}

	return fmt.Sprintf(tpl, title, headHTML, navHTML, sectionsHTML)
}
