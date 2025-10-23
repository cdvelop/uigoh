
package gosite

import (
	"os"
	"path/filepath"
	"crypto/sha256"
	"encoding/hex"

	"github.com/cdvelop/gosite/components"
	"github.com/cdvelop/tinystrings"
)

// Config holds the configuration for the site.
type Config struct {
	Title     string
	OutputDir string
}

// Site manages the global state of the website, including pages, assets, etc.
type Site struct {
	Cfg       *Config
	pages     []*components.Page
	cssBlocks map[string]string
	cssOrder  []string
	jsBlocks  map[string]string
	jsOrder   []string
}

// NewSite creates a new site manager.
func NewSite(cfg *Config) *Site {
	return &Site{
		Cfg:       cfg,
		pages:     make([]*components.Page, 0),
		cssBlocks: make(map[string]string),
		cssOrder:  make([]string, 0),
		jsBlocks:  make(map[string]string),
		jsOrder:   make([]string, 0),
	}
}

// NewPage creates a new page and registers it with the site.
func (s *Site) NewPage(title, filename string) *components.Page {
	p := &components.Page{
		Site:     s,
		Title:    title,
		Filename: filename,
	}
	s.pages = append(s.pages, p)
	return p
}

// PageCount returns the number of pages in the site.
func (s *Site) PageCount() int {
	return len(s.pages)
}

// BuildNav creates the navigation menu.
func (s *Site) BuildNav() string {
	var b tinystrings.Builder
	b.WriteString("<nav class=\"main-nav\">\n")
	for _, page := range s.pages {
		b.WriteString("  <a href=\"")
		b.WriteString(tinystrings.EscapeAttr(page.Filename))
		b.WriteString("\" class=\"nav-link\">")
		b.WriteString(tinystrings.EscapeHTML(page.Title))
		b.WriteString("</a>\n")
	}
	b.WriteString("</nav>\n")
	return b.String()
}

// AddCSS accumulates CSS with deduplication at the site level.
func (s *Site) AddCSS(css string) {
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
func (s *Site) AddJS(js string) {
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
func (s *Site) Generate() error {
	if err := os.MkdirAll(s.Cfg.OutputDir, 0755); err != nil {
		return err
	}

	for _, page := range s.pages {
		pagePath := filepath.Join(s.Cfg.OutputDir, page.Filename)
		if err := os.WriteFile(pagePath, []byte(page.RenderHTML()), 0644); err != nil {
			return err
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
func (s *Site) writeCSS() error {
	var b tinystrings.Builder
	for _, hash := range s.cssOrder {
		b.WriteString(s.cssBlocks[hash])
		b.WriteString("\n")
	}
	cssPath := filepath.Join(s.Cfg.OutputDir, "style.css")
	return os.WriteFile(cssPath, []byte(b.String()), 0644)
}

// writeJS writes the combined and deduplicated JS to a file.
func (s *Site) writeJS() error {
	var b tinystrings.Builder
	for _, hash := range s.jsOrder {
		b.WriteString(s.jsBlocks[hash])
		b.WriteString("\n")
	}
	jsPath := filepath.Join(s.Cfg.OutputDir, "main.js")
	return os.WriteFile(jsPath, []byte(b.String()), 0644)
}

// hashString creates a SHA-256 hash of a string.
func hashString(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
