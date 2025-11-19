package servicecard

import (
	. "github.com/cdvelop/tinystring"
)

// ServiceCard implements HTMLRenderer and CSSRenderer interfaces.
// It provides a service card with icon, title, and description.
type ServiceCard struct {
	Title       string
	Description string
	IconSrc     string
	CSSClass    string
}

// RenderHTML generates the HTML for the service card.
func (s *ServiceCard) RenderHTML() string {
	class := "service-item"
	if s.CSSClass != "" {
		class += " " + s.CSSClass
	}
	classEsc := Convert(class).EscapeAttr()

	imageSrcEsc := Convert(s.IconSrc).EscapeAttr()
	titleEsc := Convert(s.Title).EscapeHTML()
	descriptionEsc := Convert(s.Description).EscapeHTML()

	tpl := `    <article class="%s">
        <div class="icon">
            <img src="%s">
        </div>
        <h3>%s</h3>
        <p class="text text-sm">%s</p>
    </article>
`

	return Fmt(tpl, classEsc, imageSrcEsc, titleEsc, descriptionEsc)
}
