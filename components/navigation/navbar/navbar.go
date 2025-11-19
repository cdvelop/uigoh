package navbar

import (
	. "github.com/cdvelop/tinystring"
)

// NavItem represents a navigation link item
type NavItem struct {
	Label string
	Href  string
}

// Navbar implements HTMLRenderer, CSSRenderer, and JSRenderer interfaces.
// It provides a responsive navigation bar with logo, menu items, and search bar.
type Navbar struct {
	LogoSrc    string
	LogoAlt    string
	LogoHref   string
	NavItems   []NavItem
	ShowSearch bool
	BgColor    string // CSS class for background color
	CSSClass   string
}

// RenderHTML generates the HTML for the navbar.
func (n *Navbar) RenderHTML() string {
	// Build background class
	bgClass := "navbar bg-blue"
	if n.BgColor != "" {
		bgClass = "navbar " + n.BgColor
	}
	if n.CSSClass != "" {
		bgClass += " " + n.CSSClass
	}
	bgClassEsc := Convert(bgClass).EscapeAttr()

	// Build logo
	logoSrcEsc := Convert(n.LogoSrc).EscapeAttr()
	logoAltEsc := Convert(n.LogoAlt).EscapeAttr()
	logoHrefEsc := Convert(n.LogoHref).EscapeAttr()

	logoHTML := Fmt(`        <a href="%s" class="navbar-brand">
            <img src="%s" alt="%s">
        </a>`, logoHrefEsc, logoSrcEsc, logoAltEsc)

	// Build nav items
	navItemsHTML := ""
	for _, item := range n.NavItems {
		labelEsc := Convert(item.Label).EscapeHTML()
		hrefEsc := Convert(item.Href).EscapeAttr()
		navItemsHTML += Fmt(`                        <li class="nav-item">
                            <a href="%s" class="nav-link">%s</a>
                        </li>
`, hrefEsc, labelEsc)
	}

	// Build search bar
	searchHTML := ""
	if n.ShowSearch {
		searchHTML = `                    <div class="search-bar">
                        <form>
                            <div class="search-bar-box flex">
                                <span class="search-icon flex">
                                    <img src="images/search-icon-dark.png">
                                </span>
                                <input type="search" class="search-control" placeholder="Search here">
                            </div>
                        </form>
                    </div>
`
	}

	tpl := `    <nav class="%s">
        <div class="container flex">
%s
        <button type="button" class="navbar-show-btn">
            <img src="images/ham-menu-icon.png">
        </button>

        <div class="navbar-collapse bg-white">
            <button type="button" class="navbar-hide-btn">
                <img src="images/close-icon.png">
            </button>
            <ul class="navbar-nav">
%s            </ul>
%s        </div> 
        </div>
    </nav>
`

	return Fmt(tpl, bgClassEsc, logoHTML, navItemsHTML, searchHTML)
}
