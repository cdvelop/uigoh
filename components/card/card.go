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
	// Build escaped values first
	class := "card"
	if c.CSSClass != "" {
		class += " " + c.CSSClass
	}

	classEsc := Convert(class).EscapeAttr()
	titleEsc := Convert(c.Title).EscapeHTML()
	descEsc := Convert(c.Description).EscapeHTML()

	iconHTML := ""
	if c.Icon != "" {
		iconEsc := Convert(c.Icon).EscapeAttr()
		iconHTML = Fmt("  <svg class=\"icon\"><use href=\"icons.svg#%s\"></use></svg>\n", iconEsc)
	}

	tpl := `<div class="%s">
%s  <h3>%s</h3>
  <p>%s</p>
</div>
`

	return Fmt(tpl, classEsc, iconHTML, titleEsc, descEsc)
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
