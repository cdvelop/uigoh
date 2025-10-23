
package components

import (
	"github.com/cdvelop/tinystrings"
)

// Card implements the Component interface for a card.
type Card struct {
	Title       string
	Description string
	Icon        string
	CSSClass    string
}

// RenderHTML generates the HTML for the card.
func (c *Card) RenderHTML() string {
	var b tinystrings.Builder
	class := "card"
	if c.CSSClass != "" {
		class += " " + c.CSSClass
	}
	b.WriteString("<div class=\"")
	b.WriteString(tinystrings.EscapeAttr(class))
	b.WriteString("\">\n")

	if c.Icon != "" {
		b.WriteString("  <svg class=\"icon\"><use href=\"icons.svg#")
		b.WriteString(tinystrings.EscapeAttr(c.Icon))
		b.WriteString("\"></use></svg>\n")
	}
	b.WriteString("  <h3>")
	b.WriteString(tinystrings.EscapeHTML(c.Title))
	b.WriteString("</h3>\n")

	b.WriteString("  <p>")
	b.WriteString(tinystrings.EscapeHTML(c.Description))
	b.WriteString("</p>\n")

	b.WriteString("</div>\n")
	return b.String()
}

// RenderCSS returns the CSS for the card.
func (c *Card) RenderCSS() string {
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

// RenderJS returns the JavaScript for the card (empty in this case).
func (c *Card) RenderJS() string {
	return ""
}
