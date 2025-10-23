
package components

import "strings"

type CarouselImage struct {
    Src string
    Alt string
}

type CarouselConfig struct {
    Images []CarouselImage
}

func (c *CarouselConfig) RenderHTML() string {
    var b strings.Builder

    b.WriteString("<div class=\"carousel\">\n")

    for _, img := range c.Images {
        b.WriteString("  <div class=\"carousel-item\">\n")
        b.WriteString("    <img src=\"")
        b.WriteString(escapeAttr(img.Src))
        b.WriteString("\" alt=\"")
        b.WriteString(escapeAttr(img.Alt))
        b.WriteString("\">\n")
        b.WriteString("  </div>\n")
    }

    b.WriteString("</div>\n")

    return b.String()
}

func (c *CarouselConfig) RenderCSS() string {
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

func (c *CarouselConfig) RenderJS() string {
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
