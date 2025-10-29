//go:build !wasm

package gosite

import (
	. "github.com/cdvelop/tinystring"
)

// Site manages the global state of the website for the backend.
// It includes fields for asset management and file generation.
type Site struct {
	Cfg       *Config
	pages     []*Page
	cssBlocks []assetBlock
	jsBlocks  []assetBlock
	buff      *Conv
}

// New creates a new site manager for the backend.
func New(cfg *Config) *Site {
	if cfg.ColorScheme == nil {
		cfg.ColorScheme = DefaultColorScheme()
	}
	return &Site{
		Cfg:       cfg,
		pages:     make([]*Page, 0),
		cssBlocks: make([]assetBlock, 0),
		jsBlocks:  make([]assetBlock, 0),
		buff:      Convert(),
	}
}

// AddCSS accumulates CSS with deduplication at the site level.
func (s *Site) AddCSS(css string) {
	if css == "" {
		return
	}
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
		pagePath := PathJoin(s.Cfg.OutputDir, page.filename).String()
		if err := s.Cfg.WriteFile(pagePath, page.RenderHTML()); err != nil {
			// In Go, it's conventional to return errors rather than panic.
			// The caller can decide how to handle the error.
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

// generateBaseCSS generates the base CSS with variables and reset styles.
func (s *Site) generateBaseCSS() string {
	cs := s.Cfg.ColorScheme
	tpl := `:root {
	--color-primary: %s;
	--color-secondary: %s;
	--color-text: %s;
	--color-background: %s;
	--color-border: %s;
	--color-heading: %s;
	--color-card-bg: %s;
}
*, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }
body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background: var(--color-background); color: var(--color-text); line-height: 1.6; }
section { padding: 2rem; max-width: 1200px; margin: 0 auto; }
h1 { color: var(--color-heading); font-size: 2.5rem; margin-bottom: 1.5rem; text-align: center; }
h2 { color: var(--color-heading); font-size: 2rem; margin-bottom: 1rem; }
.card-container { display: grid; grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); gap: 1.5rem; margin-top: 2rem; }
`
	return Fmt(tpl, cs.Primary, cs.Secondary, cs.Text, cs.Background, cs.Border, cs.Primary, cs.Background)
}

// writeCSSFile writes the combined CSS to a file.
func (s *Site) writeCSSFile() error {
	if len(s.cssBlocks) == 0 {
		return nil // No CSS to write
	}
	s.buff.Reset()
	s.buff.Write(s.generateBaseCSS())
	for _, b := range s.cssBlocks {
		s.buff.Write(b.Content)
		s.buff.Write("\n")
	}
	cssPath := PathJoin(s.Cfg.OutputDir, "style.css").String()
	return s.Cfg.WriteFile(cssPath, s.buff.String())
}

// writeJSFile writes the combined JS to a file.
func (s *Site) writeJSFile() error {
	if len(s.jsBlocks) == 0 {
		return nil // No JS to write
	}
	s.buff.Reset()
	for _, b := range s.jsBlocks {
		s.buff.Write(b.Content)
		s.buff.Write("\n")
	}
	jsPath := PathJoin(s.Cfg.OutputDir, "script.js").String()
	return s.Cfg.WriteFile(jsPath, s.buff.String())
}
