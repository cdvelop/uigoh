
package components

import (
    "strings"
)

type SectionBuilder struct {
    title    string
    moduleID string // Auto-generated from caller module
    content  []string
    site     *Site  // Reference to site for component methods
    page     *Page  // Reference to parent page
}

// Render generates the section HTML with auto-generated ID
func (s *SectionBuilder) Render() string {
    var b strings.Builder

    b.WriteString("<section id=\"")
    b.WriteString(s.moduleID) // Auto-generated ID
    b.WriteString("\" class=\"page\">\n")

    if s.title != "" {
        b.WriteString("  <h1>")
        b.WriteString(escapeHTML(s.title))
        b.WriteString("</h1>\n")
    }

    b.WriteString("  <div class=\"card-container\">\n")
    for _, item := range s.content {
        b.WriteString("    ")
        b.WriteString(item)
        b.WriteString("\n")
    }
    b.WriteString("  </div>\n")

    b.WriteString("</section>\n")

    return b.String()
}

// GetNavItem returns navigation item data for auto-nav generation
func (s *SectionBuilder) GetNavItem() NavItem {
    return NavItem{
        Label: s.title,
        Href:  "#" + s.moduleID,
        Icon:  s.detectIcon(), // Optional: auto-detect icon
    }
}

// detectIcon attempts to auto-detect icon from module name
func (s *SectionBuilder) detectIcon() string {
    // Map common module names to icons
    iconMap := map[string]string{
        "homepage":     "icon-home",
        "servicespage": "icon-service",
        "staffpage":    "icon-staff",
        "contactpage":  "icon-contact",
    }

    if icon, ok := iconMap[s.moduleID]; ok {
        return icon
    }

    return "" // No icon
}

// AddCard adds a card to the section
func (s *SectionBuilder) AddCard(title, description, icon string) *SectionBuilder {
    card := s.site.Card(title, description, icon)
    s.content = append(s.content, card)
    return s
}

// AddCarousel adds a carousel to the section
func (s *SectionBuilder) AddCarousel(images []CarouselImage) *SectionBuilder {
    carousel := s.site.Carousel(images)
    s.content = append(s.content, carousel)
    return s
}

// AddForm adds a form to the section
func (s *SectionBuilder) AddForm(config FormConfig) *SectionBuilder {
    form := s.site.Form(config)
    s.content = append(s.content, form)
    return s
}

// RenderCSS returns section CSS
func (s *SectionBuilder) RenderCSS() string {
    return `.page {
  min-height: 100vh;
  padding: 2rem;
}

.card-container {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 1.5rem;
  max-width: 1200px;
  margin: 0 auto;
}
`
}
