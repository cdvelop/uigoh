package gosite

import (
	. "github.com/cdvelop/tinystring"
)

// ColorScheme holds the basic color configuration for the site.
type ColorScheme struct {
	Primary    string // Main brand color
	Secondary  string // Secondary/accent color
	Text       string // Main text color
	Background string // Background color
	Border     string // Border color
}

// DefaultColorScheme returns the default color scheme based on monjitaschillan.cl
func DefaultColorScheme() *ColorScheme {
	return &ColorScheme{
		Primary:    "#3f88bf", // Blue
		Secondary:  "#ff9300", // Orange
		Text:       "#000000", // Black
		Background: "#ffffff", // White
		Border:     "#e9e9e9", // Light gray
	}
}

// Config holds the configuration for the site.
type Config struct {
	Title       string
	OutputDir   string // eg: "dist"
	ColorScheme *ColorScheme

	WriteFile func(path string, content string) error
}

// Site manages the global state of the website, including pages, assets, etc.
type Site struct {
	Cfg   *Config
	pages []*page
	// ordered slices of assetBlock preserve insertion order and allow
	// deduplication by hash without using maps (maps don't preserve order)
	cssBlocks []assetBlock
	jsBlocks  []assetBlock

	buff *Conv // reusable buffer

}

// NewSite creates a new site manager.
func NewSite(cfg *Config) *Site {
	// Set default color scheme if not provided
	if cfg.ColorScheme == nil {
		cfg.ColorScheme = DefaultColorScheme()
	}

	return &Site{
		Cfg:       cfg,
		pages:     make([]*page, 0),
		cssBlocks: make([]assetBlock, 0),
		jsBlocks:  make([]assetBlock, 0),
		buff:      Convert(),
	}
}

// NewPage creates a new page and registers it with the site.
func (s *Site) NewPage(title, filename string) *page {
	p := &page{
		site:     s,
		title:    title,
		filename: filename,
	}
	s.pages = append(s.pages, p)
	return p
}

// PageCount returns the number of pages in the site.
func (s *Site) PageCount() int {
	return len(s.pages)
}

// BuildNav creates the navigation menu using NavbarBuilder.
func (s *Site) BuildNav() string {
	nav := &NavbarBuilder{site: s}
	// Register navbar CSS and JS once
	s.AddCSS(nav.RenderCSS())
	s.AddJS(nav.RenderJS())
	return nav.Render()
}

// AddCSS accumulates CSS with deduplication at the site level.
func (s *Site) AddCSS(css string) {
	if css == "" {
		return
	}

	// Check if this CSS block already exists
	for _, existing := range s.cssBlocks {
		if existing.Content == css {
			return // Already added, skip duplicate
		}
	}

	s.cssBlocks = append(s.cssBlocks, assetBlock{Content: css})
}

// AddJS accumulates JavaScript with deduplication at the site level.
func (s *Site) AddJS(js string) {
	if js == "" {
		return
	}

	// Check if this JS block already exists
	for _, existing := range s.jsBlocks {
		if existing.Content == js {
			return // Already added, skip duplicate
		}
	}

	s.jsBlocks = append(s.jsBlocks, assetBlock{Content: js})
}

// Generate renders all site files to disk.
func (s *Site) Generate() error {

	for _, page := range s.pages {
		pagePath := PathJoin(s.Cfg.OutputDir, page.filename)
		if err := s.Cfg.WriteFile(pagePath, page.RenderHTML()); err != nil {
			return err
		}
	}

	if err := s.writeCSSFile(); err != nil {
		return err
	}
	if err := s.writeJSFile(); err != nil {
		return err
	}
	return nil
}

// generateBaseCSS generates the base CSS with variables and reset styles
func (s *Site) generateBaseCSS() string {
	cs := s.Cfg.ColorScheme
	s.buff.Reset()

	// CSS Variables
	s.buff.Write(":root {\n")
	s.buff.Write("  --color-primary: ")
	s.buff.Write(cs.Primary)
	s.buff.Write(";\n")
	s.buff.Write("  --color-secondary: ")
	s.buff.Write(cs.Secondary)
	s.buff.Write(";\n")
	s.buff.Write("  --color-text: ")
	s.buff.Write(cs.Text)
	s.buff.Write(";\n")
	s.buff.Write("  --color-background: ")
	s.buff.Write(cs.Background)
	s.buff.Write(";\n")
	s.buff.Write("  --color-border: ")
	s.buff.Write(cs.Border)
	s.buff.Write(";\n")
	s.buff.Write("  --color-heading: ")
	s.buff.Write(cs.Primary)
	s.buff.Write(";\n")
	s.buff.Write("  --color-card-bg: ")
	s.buff.Write(cs.Background)
	s.buff.Write(";\n")
	s.buff.Write("}\n\n")

	// CSS Reset
	s.buff.Write("*, *::before, *::after {\n")
	s.buff.Write("  box-sizing: border-box;\n")
	s.buff.Write("  margin: 0;\n")
	s.buff.Write("  padding: 0;\n")
	s.buff.Write("}\n\n")

	// Base styles
	s.buff.Write("body {\n")
	s.buff.Write("  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;\n")
	s.buff.Write("  background: var(--color-background);\n")
	s.buff.Write("  color: var(--color-text);\n")
	s.buff.Write("  line-height: 1.6;\n")
	s.buff.Write("  padding: 0;\n")
	s.buff.Write("  margin: 0;\n")
	s.buff.Write("}\n\n")

	// Section styles
	s.buff.Write("section {\n")
	s.buff.Write("  padding: 2rem;\n")
	s.buff.Write("  max-width: 1200px;\n")
	s.buff.Write("  margin: 0 auto;\n")
	s.buff.Write("}\n\n")

	s.buff.Write("h1 {\n")
	s.buff.Write("  color: var(--color-heading);\n")
	s.buff.Write("  font-size: 2.5rem;\n")
	s.buff.Write("  margin-bottom: 1.5rem;\n")
	s.buff.Write("  text-align: center;\n")
	s.buff.Write("}\n\n")

	s.buff.Write("h2 {\n")
	s.buff.Write("  color: var(--color-heading);\n")
	s.buff.Write("  font-size: 2rem;\n")
	s.buff.Write("  margin-bottom: 1rem;\n")
	s.buff.Write("}\n\n")

	// Card container (grid layout)
	s.buff.Write(".card-container {\n")
	s.buff.Write("  display: grid;\n")
	s.buff.Write("  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));\n")
	s.buff.Write("  gap: 1.5rem;\n")
	s.buff.Write("  margin-top: 2rem;\n")
	s.buff.Write("}\n\n")

	return s.buff.String()
}

// writeCSSFile writes separate CSS files (base + components)
func (s *Site) writeCSSFile() error {
	s.buff.Reset()

	// Write base CSS (variables + reset + layout)
	s.buff.Write(s.generateBaseCSS())

	// Add component CSS
	for _, b := range s.cssBlocks {
		s.buff.Write(b.Content)
		s.buff.Write("\n")
	}

	cssPath := PathJoin(s.Cfg.OutputDir, "style.css")
	return s.Cfg.WriteFile(cssPath, s.buff.String())
}

// writeJSFile writes the combined and deduplicated JS to a file.
func (s *Site) writeJSFile() error {
	s.buff.Reset()
	for _, b := range s.jsBlocks {
		s.buff.Write(b.Content)
		s.buff.Write("\n")
	}
	jsPath := PathJoin(s.Cfg.OutputDir, "script.js")
	return s.Cfg.WriteFile(jsPath, s.buff.String())
}
