package footer

import (
	. "github.com/cdvelop/tinystring"
)

// FooterColumn represents a footer column section
type FooterColumn struct {
	Title   string
	Content FooterContent
}

// FooterContent can be different types of content
type FooterContent struct {
	Type            string
	Text            string
	LogoSrc         string
	Address         string
	Tags            []string
	Links           []Link
	AppointmentInfo []string
}

// Link represents a footer link
type Link struct {
	Label string
	Href  string
}

// SocialLink represents a social media link
type SocialLink struct {
	IconClass string
	Href      string
}

// Footer implements HTMLRenderer and CSSRenderer interfaces.
// It provides a multi-column footer with various content types.
type Footer struct {
	Columns     []FooterColumn
	SocialLinks []SocialLink
	CSSClass    string
}

// RenderHTML generates the HTML for the footer.
func (f *Footer) RenderHTML() string {
	class := "footer text-center"
	if f.CSSClass != "" {
		class += " " + f.CSSClass
	}
	classEsc := Convert(class).EscapeAttr()

	columnsHTML := ""
	for _, col := range f.Columns {
		titleEsc := Convert(col.Title).EscapeHTML()
		contentHTML := f.renderColumnContent(col.Content)
		columnsHTML += Fmt(`            <div class="footer-item">
                <h3 class="footer-head">%s</h3>
%s            </div>

`, titleEsc, contentHTML)
	}

	socialHTML := ""
	if len(f.SocialLinks) > 0 {
		for _, social := range f.SocialLinks {
			iconEsc := Convert(social.IconClass).EscapeAttr()
			hrefEsc := Convert(social.Href).EscapeAttr()
			socialHTML += Fmt(`                    <li><a href="%s" class="text-white flex"> <i class="%s"></i></a></li>
`, hrefEsc, iconEsc)
		}
	}

	tpl := `    <footer class="%s">
        <div class="container">
            <div class="footer-inner text-white py grid">
%s            </div>

            <div class="footer-links">
                <ul class="flex">
%s                </ul>
            </div>
        </div>
    </footer>
`

	return Fmt(tpl, classEsc, columnsHTML, socialHTML)
}

func (f *Footer) renderColumnContent(content FooterContent) string {
	switch content.Type {
	case "about":
		logoEsc := Convert(content.LogoSrc).EscapeAttr()
		textEsc := Convert(content.Text).EscapeHTML()
		addressEsc := Convert(content.Address).EscapeHTML()
		return Fmt(`                <div class="icon">
                    <img src="%s">
                </div>
                <p class="text text-md">%s</p>
                <address>%s</address>
`, logoEsc, textEsc, addressEsc)

	case "tags":
		tagsHTML := ""
		for _, tag := range content.Tags {
			tagEsc := Convert(tag).EscapeHTML()
			tagsHTML += Fmt("                        <li>%s</li>\n", tagEsc)
		}
		return Fmt(`                <ul class="tags-list flex">
%s                </ul>
`, tagsHTML)

	case "links":
		linksHTML := ""
		for _, link := range content.Links {
			labelEsc := Convert(link.Label).EscapeHTML()
			hrefEsc := Convert(link.Href).EscapeAttr()
			linksHTML += Fmt(`                        <li><a href="%s" class="text-white">%s</a></li>
`, hrefEsc, labelEsc)
		}
		return Fmt(`                <ul>
%s                </ul>
`, linksHTML)

	case "appointment":
		textEsc := Convert(content.Text).EscapeHTML()
		infoHTML := ""
		for _, info := range content.AppointmentInfo {
			infoEsc := Convert(info).EscapeHTML()
			infoHTML += Fmt("                        <li>%s</li>\n", infoEsc)
		}
		return Fmt(`                <p class="text text-md">%s</p>
                <ul class="appointment-info">
%s                </ul>
`, textEsc, infoHTML)

	default:
		return ""
	}
}
