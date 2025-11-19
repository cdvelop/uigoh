package contactform

import (
	. "github.com/cdvelop/tinystring"
)

// ContactForm implements HTMLRenderer and CSSRenderer interfaces.
// It provides a contact form with optional embedded map.
type ContactForm struct {
	Title       string
	Description string
	MapEmbedURL string
	ShowMap     bool
	BgColor     string
	CSSClass    string
}

// RenderHTML generates the HTML for the contact form.
func (c *ContactForm) RenderHTML() string {
	class := "contact py"
	if c.CSSClass != "" {
		class += " " + c.CSSClass
	}
	classEsc := Convert(class).EscapeAttr()

	rightBg := "contact-right text-white text-center bg-blue"
	if c.BgColor != "" {
		rightBg = "contact-right text-white text-center " + c.BgColor
	}
	rightBgEsc := Convert(rightBg).EscapeAttr()

	titleEsc := Convert(c.Title).EscapeHTML()
	descEsc := Convert(c.Description).EscapeHTML()

	mapHTML := ""
	if c.ShowMap && c.MapEmbedURL != "" {
		mapURLEsc := Convert(c.MapEmbedURL).EscapeAttr()
		mapHTML = Fmt(`            <div class="contact-left">
                <iframe src="%s" width="600" height="450" style="border:0;" allowfullscreen="" loading="lazy"></iframe>
            </div>
`, mapURLEsc)
	}

	tpl := `    <section class="%s">
        <div class="container grid">
%s            <div class="%s">
                <div class="contact-head">
                    <h3 class="lead">%s</h3>
                    <p class="text text-md">%s</p>
                </div>
                <form>
                    <div class="form-element">
                        <input type="text" class="form-control" placeholder="Your name">
                    </div>
                    <div class="form-element">
                        <input type="email" class="form-control" placeholder="Your email">
                    </div>
                    <div class="form-element">
                        <textarea rows="5" placeholder="Your Message" class="form-control"></textarea>
                    </div>
                    <button type="submit" class="btn btn-white btn-submit">
                        <i class="fas fa-arrow-right"></i> Send Message
                    </button>
                </form>
            </div>
        </div>
    </section>
`

	return Fmt(tpl, classEsc, mapHTML, rightBgEsc, titleEsc, descEsc)
}
