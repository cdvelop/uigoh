package hero

import (
	. "github.com/cdvelop/tinystring"
)

// Button represents a button in the hero section
type Button struct {
	Label    string
	Href     string
	CSSClass string // e.g., "btn-white", "btn-light-blue"
}

// Hero implements HTMLRenderer, CSSRenderer interfaces.
// It provides a hero/header section with title, description, image, and call-to-action buttons.
type Hero struct {
	Title       string   // Main title
	TitleSpan   string   // Highlighted part of title (optional)
	Lead        string   // Lead text below title
	Description string   // Main description text
	ImageSrc    string   // Hero image source
	ImageAlt    string   // Hero image alt text
	Buttons     []Button // Call-to-action buttons
	BgColor     string   // CSS class for background color
	CSSClass    string
}

// RenderHTML generates the HTML for the hero section.
func (h *Hero) RenderHTML() string {
	// Build background class
	bgClass := "header bg-blue"
	if h.BgColor != "" {
		bgClass = "header " + h.BgColor
	}
	if h.CSSClass != "" {
		bgClass += " " + h.CSSClass
	}
	bgClassEsc := Convert(bgClass).EscapeAttr()

	// Build title with optional span
	titleHTML := Convert(h.Title).EscapeHTML()
	if h.TitleSpan != "" {
		spanEsc := Convert(h.TitleSpan).EscapeHTML()
		titleHTML += Fmt("<br> <span>%s</span>", spanEsc)
	}

	// Build lead and description
	leadEsc := Convert(h.Lead).EscapeHTML()
	descEsc := Convert(h.Description).EscapeHTML()

	// Build buttons
	buttonsHTML := ""
	for _, btn := range h.Buttons {
		labelEsc := Convert(btn.Label).EscapeHTML()
		hrefEsc := Convert(btn.Href).EscapeAttr()
		btnClassEsc := Convert("btn " + btn.CSSClass).EscapeAttr()
		buttonsHTML += Fmt(`                        <a href="%s" class="%s">%s</a>
`, hrefEsc, btnClassEsc, labelEsc)
	}

	// Build image
	imgSrcEsc := Convert(h.ImageSrc).EscapeAttr()
	imgAltEsc := Convert(h.ImageAlt).EscapeAttr()

	tpl := `    <header class="%s">
        <div class="header-inner text-white text-center">
            <div class="container grid">
                <div class="header-inner-left">
                    <h1>%s</h1>
                    <p class="lead">%s</p>
                    <p class="text text-md">%s</p>
                    <div class="btn-group">
%s                    </div>
                </div>
                <div class="header-inner-right">
                    <img src="%s" alt="%s">
                </div>
            </div>
        </div>
    </header>
`

	return Fmt(tpl, bgClassEsc, titleHTML, leadEsc, descEsc, buttonsHTML, imgSrcEsc, imgAltEsc)
}
