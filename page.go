package gosite

import (
	. "github.com/cdvelop/tinystring"
)

// Page represents a single HTML page. Its fields are unexported to maintain
// a controlled, fluent API.
type Page struct {
	site     SiteLink
	sections []*Section
	title    string
	filename string
	head     []string
}

// NewSection adds a new section to the page and returns it for chaining.
func (p *Page) NewSection(title string) *Section {
	section := &Section{
		page:    p,
		site:    p.site,
		Title:   title,
		content: make([]any, 0),
	}
	p.sections = append(p.sections, section)
	return section
}

// AddHead adds content to the <head> section of the page.
func (p *Page) AddHead(content string) *Page {
	p.head = append(p.head, content)
	return p
}

// RenderHTML generates the complete HTML for the page.
func (p *Page) RenderHTML() string {
	b := Convert()

	// Build head entries
	for _, h := range p.head {
		b.Write("  ")
		b.Write(h)
		b.Write("\n")
	}
	headHTML := b.String()

	// Build sections
	b.Reset()
	for _, section := range p.sections {
		b.Write(section.Render())
	}
	sectionsHTML := b.String()

	title := Convert(p.title).EscapeHTML()

	// Optionally include nav if multiple pages exist
	navHTML := ""
	if p.site.PageCount() > 1 {
		navHTML = p.site.BuildNav()
	}

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
	return Fmt(tpl, title, headHTML, navHTML, sectionsHTML)
}
