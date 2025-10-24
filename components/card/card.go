package card

import (
	. "github.com/cdvelop/tinystring"
)

// Card implements HTMLRenderer, CSSRenderer, and JSRenderer interfaces.
type Card struct {
	Title       string
	Description string
	Icon        string
	CSSClass    string
}

// RenderHTML generates the HTML for the card.
func (c *Card) RenderHTML() string {
	var b = Convert()
	class := "card"
	if c.CSSClass != "" {
		class += " " + c.CSSClass
	}
	b.Write("<div class=\"")
	b.Write(Convert(class).EscapeAttr())
	b.Write("\">\n")

	if c.Icon != "" {
		b.Write("  <svg class=\"icon\"><use href=\"icons.svg#")
		b.Write(Convert(c.Icon).EscapeAttr())
		b.Write("\"></use></svg>\n")
	}
	b.Write("  <h3>")
	b.Write(Convert(c.Title).EscapeHTML())
	b.Write("</h3>\n")

	b.Write("  <p>")
	b.Write(Convert(c.Description).EscapeHTML())
	b.Write("</p>\n")

	b.Write("</div>\n")
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
