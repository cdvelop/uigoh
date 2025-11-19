package banner

import (
	. "github.com/cdvelop/tinystring"
)

// BannerType defines the banner style
type BannerType string

const (
	BannerTypeQuote  BannerType = "quote"
	BannerTypeAction BannerType = "action"
)

// Button represents a button in the banner
type Button struct {
	Label    string
	Href     string
	CSSClass string
}

// Banner implements HTMLRenderer and CSSRenderer interfaces.
// It provides a banner section with quote or call-to-action content.
type Banner struct {
	Type       BannerType
	Quote      string
	Author     string
	Text       string
	ImageSrc   string
	Buttons    []Button
	BgImageSrc string
	CSSClass   string
}

// RenderHTML generates the HTML for the banner.
func (b *Banner) RenderHTML() string {
	class := "banner"
	if b.Type == BannerTypeQuote {
		class = "banner-one text-center"
	} else {
		class = "banner-two text-center"
	}
	if b.CSSClass != "" {
		class += " " + b.CSSClass
	}
	classEsc := Convert(class).EscapeAttr()

	var content string
	if b.Type == BannerTypeQuote {
		quoteEsc := Convert(b.Quote).EscapeHTML()
		authorEsc := Convert(b.Author).EscapeHTML()
		content = Fmt(`        <div class="container text-white">
            <blockquote class="lead"><i class="fas fa-quote-left"></i> %s <i class="fas fa-quote-right"></i></blockquote>
            <small class="text text-sm">- %s</small>
        </div>
`, quoteEsc, authorEsc)
	} else {
		textEsc := Convert(b.Text).EscapeHTML()
		imgSrcEsc := Convert(b.ImageSrc).EscapeAttr()

		buttonsHTML := ""
		for _, btn := range b.Buttons {
			labelEsc := Convert(btn.Label).EscapeHTML()
			hrefEsc := Convert(btn.Href).EscapeAttr()
			btnClassEsc := Convert("btn " + btn.CSSClass).EscapeAttr()
			buttonsHTML += Fmt(`                    <a href="%s" class="%s">%s</a>
`, hrefEsc, btnClassEsc, labelEsc)
		}

		content = Fmt(`        <div class="container grid">
            <div class="banner-two-left">
                <img src="%s">
            </div>
            <div class="banner-two-right">
                <p class="lead text-white">%s</p>
                <div class="btn-group">
%s                </div>
            </div>
        </div>
`, imgSrcEsc, textEsc, buttonsHTML)
	}

	tpl := `    <section class="%s">
%s    </section>
`

	return Fmt(tpl, classEsc, content)
}
