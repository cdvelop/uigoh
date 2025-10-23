
package components

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// site manages the global state of the website, including pages, assets, etc.
type site struct {
	cfg       *Config
	pages     []*Page
	cssBlocks map[string]string
	cssOrder  []string
	jsBlocks  map[string]string
	jsOrder   []string
}

// newSite creates a new site manager.
func newSite(cfg *Config) *site {
	return &site{
		cfg:       cfg,
		pages:     make([]*Page, 0),
		cssBlocks: make(map[string]string),
		cssOrder:  make([]string, 0),
		jsBlocks:  make(map[string]string),
		jsOrder:   make([]string, 0),
	}
}

// AddPage registers a new page with the site.
func (s *site) AddPage(p *Page) {
	s.pages = append(s.pages, p)
}

// AddCSS accumulates CSS with deduplication at the site level.
func (s *site) AddCSS(css string) {
	if css == "" {
		return
	}
	hash := hashString(css)
	if _, exists := s.cssBlocks[hash]; !exists {
		s.cssBlocks[hash] = css
		s.cssOrder = append(s.cssOrder, hash)
	}
}

// AddJS accumulates JavaScript with deduplication at the site level.
func (s *site) AddJS(js string) {
	if js == "" {
		return
	}
	hash := hashString(js)
	if _, exists := s.jsBlocks[hash]; !exists {
		s.jsBlocks[hash] = js
		s.jsOrder = append(s.jsOrder, hash)
	}
}

// Generate renders all site files to disk.
func (s *site) Generate() error {
	if err := os.MkdirAll(s.cfg.OutputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	for _, page := range s.pages {
		pagePath := filepath.Join(s.cfg.OutputDir, page.Filename)
		if err := os.WriteFile(pagePath, []byte(page.RenderHTML()), 0644); err != nil {
			return fmt.Errorf("failed to write page %s: %w", page.Filename, err)
		}
	}

	if err := s.writeCSS(); err != nil {
		return err
	}
	if err := s.writeJS(); err != nil {
		return err
	}
	return nil
}

// writeCSS writes the combined and deduplicated CSS to a file.
func (s *site) writeCSS() error {
	var b strings.Builder
	for _, hash := range s.cssOrder {
		b.WriteString(s.cssBlocks[hash])
		b.WriteString("\n")
	}
	cssPath := filepath.Join(s.cfg.OutputDir, "style.css")
	return os.WriteFile(cssPath, []byte(b.String()), 0644)
}

// writeJS writes the combined and deduplicated JS to a file.
func (s *site) writeJS() error {
	var b strings.Builder
	for _, hash := range s.jsOrder {
		b.WriteString(s.jsBlocks[hash])
		b.WriteString("\n")
	}
	jsPath := filepath.Join(s.cfg.OutputDir, "main.js")
	return os.WriteFile(jsPath, []byte(b.String()), 0644)
}

// buildCombinedNav creates the navigation menu.
func (s *site) buildCombinedNav() string {
	var b strings.Builder
	b.WriteString("<nav class=\"main-nav\">\n")
	for _, page := range s.pages {
		fmt.Fprintf(&b, "  <a href=\"%s\" class=\"nav-link\">%s</a>\n", escapeAttr(page.Filename), escapeHTML(page.Title))
	}
	b.WriteString("</nav>\n")
	return b.String()
}
