package packagecard

import (
	. "github.com/cdvelop/tinystring"
)

// PackageCard implements HTMLRenderer and CSSRenderer interfaces.
// It provides a package service card with icon, title, description, and button.
type PackageCard struct {
	Title       string
	Description string
	IconClass   string
	ButtonLabel string
	ButtonHref  string
	CSSClass    string
}

// RenderHTML generates the HTML for the package card.
func (p *PackageCard) RenderHTML() string {
	class := "package-service-item bg-white"
	if p.CSSClass != "" {
		class += " " + p.CSSClass
	}
	classEsc := Convert(class).EscapeAttr()

	iconClassEsc := Convert(p.IconClass).EscapeAttr()
	titleEsc := Convert(p.Title).EscapeHTML()
	descriptionEsc := Convert(p.Description).EscapeHTML()
	buttonLabelEsc := Convert(p.ButtonLabel).EscapeHTML()
	buttonHrefEsc := Convert(p.ButtonHref).EscapeAttr()

	tpl := `    <div class="%s">
        <div class="icon flex">
            <i class="%s"></i>
        </div>
        <h3>%s</h3>
        <p class="text text-sm">%s</p>
        <a href="%s" class="btn btn-blue">%s</a>
    </div>
`

	return Fmt(tpl, classEsc, iconClassEsc, titleEsc, descriptionEsc, buttonHrefEsc, buttonLabelEsc)
}
