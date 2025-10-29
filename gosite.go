package gosite

// ColorScheme holds the basic color configuration for the site.
type ColorScheme struct {
	Primary    string
	Secondary  string
	Text       string
	Background string
	Border     string
}

// DefaultColorScheme returns the default color scheme.
func DefaultColorScheme() *ColorScheme {
	return &ColorScheme{
		Primary:    "#3f88bf",
		Secondary:  "#ff9300",
		Text:       "#000000",
		Background: "#ffffff",
		Border:     "#e9e9e9",
	}
}

// Config holds the configuration for the site.
// It uses build tags to include environment-specific fields.
type Config struct {
	Title       string
	OutputDir   string
	ColorScheme *ColorScheme
	EventBinder EventBinder // Frontend only
	WriteFile   func(path string, content string) error // Backend only
}

// NewPage creates a new page and registers it with the site.
func (s *Site) NewPage(title, filename string) *Page {
	p := &Page{
		site:     s,
		title:    title,
		filename: filename,
		sections: make([]*Section, 0),
		head:     make([]string, 0),
	}
	s.pages = append(s.pages, p)
	return p
}

// PageCount returns the number of pages in the site.
func (s *Site) PageCount() int {
	return len(s.pages)
}

// BuildNav creates the navigation menu.
// This is a shared method, as nav structure is the same in both environments.
func (s *Site) BuildNav() string {
	nav := &NavbarBuilder{site: s}
	// In the backend, this will add CSS/JS. In frontend, it's a no-op.
	s.AddCSS(nav.RenderCSS())
	s.AddJS(nav.RenderJS())
	return nav.Render()
}
