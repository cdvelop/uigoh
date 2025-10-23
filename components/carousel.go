
package components

import (
	"fmt"
	"strings"
)

// CarouselImage defines the structure for an image in a carousel.
type CarouselImage struct {
	Src string
	Alt string
}

// Carousel implements the Component interface for a carousel.
type Carousel struct {
	Images []CarouselImage
}

// RenderHTML generates the HTML for the carousel.
func (c *Carousel) RenderHTML() string {
	var b strings.Builder
	b.WriteString("<div class=\"carousel\">\n")
	for _, img := range c.Images {
		fmt.Fprintf(&b, "  <div class=\"carousel-item\"><img src=\"%s\" alt=\"%s\"></div>\n", escapeAttr(img.Src), escapeAttr(img.Alt))
	}
	b.WriteString("</div>\n")
	return b.String()
}

// RenderCSS returns the CSS for the carousel.
func (c *Carousel) RenderCSS() string {
    return `.carousel {
  position: relative;
  width: 100%;
  overflow: hidden;
}

.carousel-item {
  display: none;
}

.carousel-item.active {
  display: block;
}

.carousel-item img {
  width: 100%;
  height: auto;
}
`
}

// RenderJS returns the JavaScript for the carousel.
func (c *Carousel) RenderJS() string {
    return `// Carousel auto-slide
(function() {
    const carousel = document.querySelector('.carousel');
    if (!carousel) return;

    const items = carousel.querySelectorAll('.carousel-item');
    let current = 0;

    if (items.length > 0) {
        items[current].classList.add('active');
    }

    setInterval(function() {
        if (items.length > 0) {
            items[current].classList.remove('active');
            current = (current + 1) % items.length;
            items[current].classList.add('active');
        }
    }, 3000);
})();
`
}
