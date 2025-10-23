
package components

import (
	"fmt"
	"strings"
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
	var b strings.Builder
	class := "card"
	if c.CSSClass != "" {
		class += " " + c.CSSClass
	}
	fmt.Fprintf(&b, "<div class=\"%s\">\n", escapeAttr(class))
	if c.Icon != "" {
		fmt.Fprintf(&b, "  <svg class=\"icon\"><use href=\"icons.svg#%s\"></use></svg>\n", escapeAttr(c.Icon))
	}
	fmt.Fprintf(&b, "  <h3>%s</h3>\n", escapeHTML(c.Title))
	fmt.Fprintf(&b, "  <p>%s</p>\n", escapeHTML(c.Description))
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
