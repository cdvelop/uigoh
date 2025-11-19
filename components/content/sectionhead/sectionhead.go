package sectionhead

import (
	. "github.com/cdvelop/tinystring"
)

// SectionHead implements HTMLRenderer and CSSRenderer interfaces.
// It provides a section heading with title, optional subtitle, and decorative line.
type SectionHead struct {
	Title        string
	Subtitle     string // Optional subtitle below title
	ShowBorder   bool   // Show border-line decoration
	ShowLineArt  bool   // Show line-art decoration with dots image
	DotsImageSrc string // Source for dots decoration image (only if ShowLineArt is true)
	TextCenter   bool   // Center align text
	CSSClass     string
}

// RenderHTML generates the HTML for the section head.
func (s *SectionHead) RenderHTML() string {
	// Build class
	class := "section-head"
	if s.TextCenter {
		class += " text-center"
	}
	if s.CSSClass != "" {
		class += " " + s.CSSClass
	}
	classEsc := Convert(class).EscapeAttr()

	// Build title
	titleEsc := Convert(s.Title).EscapeHTML()
	titleHTML := Fmt("    <h2>%s</h2>\n", titleEsc)

	// Build subtitle
	subtitleHTML := ""
	if s.Subtitle != "" {
		subtitleEsc := Convert(s.Subtitle).EscapeHTML()
		subtitleHTML = Fmt("    <p class=\"text text-lg\">%s</p>\n", subtitleEsc)
	}

	// Build border line
	borderHTML := ""
	if s.ShowBorder {
		borderHTML = "    <div class=\"border-line\"></div>\n"
	}

	// Build line art
	lineArtHTML := ""
	if s.ShowLineArt {
		dotsImgEsc := Convert(s.DotsImageSrc).EscapeAttr()
		lineArtHTML = Fmt(`    <div class="line-art flex">
        <div></div>
        <img src="%s">
        <div></div>
    </div>
`, dotsImgEsc)
	}

	tpl := `<div class="%s">
%s%s%s%s</div>
`

	return Fmt(tpl, classEsc, titleHTML, subtitleHTML, borderHTML, lineArtHTML)
}
