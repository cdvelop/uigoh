package doctorcard

import (
	. "github.com/cdvelop/tinystring"
)

// DoctorCard implements HTMLRenderer and CSSRenderer interfaces.
// It provides a doctor card with image and overlay info (name and specialty).
type DoctorCard struct {
	Name      string
	Specialty string
	ImageSrc  string
	ImageAlt  string
	BgColor   string
	CSSClass  string
}

// RenderHTML generates the HTML for the doctor card.
func (d *DoctorCard) RenderHTML() string {
	class := "doc-panel-item"
	if d.CSSClass != "" {
		class += " " + d.CSSClass
	}
	classEsc := Convert(class).EscapeAttr()

	bgClass := "info text-center bg-blue text-white flex"
	if d.BgColor != "" {
		bgClass = "info text-center " + d.BgColor + " text-white flex"
	}
	bgClassEsc := Convert(bgClass).EscapeAttr()

	imageSrcEsc := Convert(d.ImageSrc).EscapeAttr()
	imageAltEsc := Convert(d.ImageAlt).EscapeAttr()
	nameEsc := Convert(d.Name).EscapeHTML()
	specialtyEsc := Convert(d.Specialty).EscapeHTML()

	tpl := `    <div class="%s">
        <div class="img flex">
            <img src="%s" alt="%s">
            <div class="%s">
                <p class="lead">%s</p>
                <p class="text-lg">%s</p>
            </div>
        </div>
    </div>
`

	return Fmt(tpl, classEsc, imageSrcEsc, imageAltEsc, bgClassEsc, nameEsc, specialtyEsc)
}
