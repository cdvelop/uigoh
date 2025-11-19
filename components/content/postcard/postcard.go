package postcard

import (
	. "github.com/cdvelop/tinystring"
)

// PostCard implements HTMLRenderer and CSSRenderer interfaces.
// It provides a blog post card with image, title, content, date, and comments count.
type PostCard struct {
	Title         string
	ImageSrc      string
	ImageAlt      string
	Content       string
	ContentExtra  string
	Date          string
	CommentsCount string
	CSSClass      string
}

// RenderHTML generates the HTML for the post card.
func (p *PostCard) RenderHTML() string {
	class := "post-item bg-white"
	if p.CSSClass != "" {
		class += " " + p.CSSClass
	}
	classEsc := Convert(class).EscapeAttr()

	imageSrcEsc := Convert(p.ImageSrc).EscapeAttr()
	imageAltEsc := Convert(p.ImageAlt).EscapeAttr()
	titleEsc := Convert(p.Title).EscapeHTML()
	contentEsc := Convert(p.Content).EscapeHTML()
	contentExtraEsc := Convert(p.ContentExtra).EscapeHTML()
	dateEsc := Convert(p.Date).EscapeHTML()
	commentsCountEsc := Convert(p.CommentsCount).EscapeHTML()

	tpl := `    <article class="%s">
        <div class="img">
            <img src="%s" alt="%s">
        </div>
        <div class="content">
            <h4>%s</h4>
            <p class="text text-sm">%s</p>
            <p class="text text-sm">%s</p>
            <div class="info flex">
                <small class="text text-sm"><i class="fas fa-clock"></i> %s</small>
                <small class="text text-sm"><i class="fas fa-comment"></i> %s</small>
            </div>
        </div>
    </article>
`

	return Fmt(tpl, classEsc, imageSrcEsc, imageAltEsc, titleEsc, contentEsc, contentExtraEsc, dateEsc, commentsCountEsc)
}
