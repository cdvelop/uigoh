
package components

import "strings"

type CardConfig struct {
    Title       string
    Description string
    Icon        string
    CSSClass    string
}

// RenderHTML returns the card HTML structure
func (c *CardConfig) RenderHTML() string {
    var b strings.Builder

    b.WriteString("<div class=\"card")
    if c.CSSClass != "" {
        b.WriteString(" ")
        b.WriteString(c.CSSClass)
    }
    b.WriteString("\">\n")

    if c.Icon != "" {
        b.WriteString("  <svg class=\"icon\"><use href=\"icons.svg#")
        b.WriteString(c.Icon)
        b.WriteString("\"></use></svg>\n")
    }

    b.WriteString("  <h3>")
    b.WriteString(escapeHTML(c.Title))
    b.WriteString("</h3>\n")

    b.WriteString("  <p>")
    b.WriteString(escapeHTML(c.Description))
    b.WriteString("</p>\n")

    b.WriteString("</div>\n")

    return b.String()
}

// RenderCSS returns the card CSS styles
func (c *CardConfig) RenderCSS() string {
    return `.card {
  border: 1px solid var(--color-border);
  border-radius: 8px;
  padding: 1.5rem;
  background: var(--color-card-bg);
  transition: transform 0.2s;
}

.card:hover {
  transform: translateY(-4px);
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
}

.card .icon {
  width: 48px;
  height: 48px;
  margin-bottom: 1rem;
}

.card h3 {
  margin: 0 0 0.5rem 0;
  color: var(--color-heading);
}

.card p {
  margin: 0;
  color: var(--color-text);
}
`
}

// RenderJS returns the card JavaScript (if needed)
func (c *CardConfig) RenderJS() string {
    // Most components won't need JS
    return ""
}
