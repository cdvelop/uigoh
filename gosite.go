package gosite

import (
	. "github.com/cdvelop/tinystring"
)

type writerFile interface {
	WriteFile(name string, data []byte, perm uint32) error
}

// Config holds the configuration for the site.
type Config struct {
	Title     string
	OutputDir string // eg: "dist"
	// 	type writerFile interface {
	//     WriteFile(name string, data []byte, perm uint32) error
	// }
	WriteFile writerFile
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

// BuildNav creates the navigation menu.
func (s *Site) BuildNav() string {
	s.buff.Reset()
	s.buff.Write("<nav class=\"main-nav\">\n")
	for _, page := range s.pages {
		s.buff.Write("  <a href=\"")
		s.buff.Write(Convert(page.filename).EscapeAttr())
		s.buff.Write("\" class=\"nav-link\">")
		s.buff.Write(Convert(page.title).EscapeHTML())
		s.buff.Write("</a>\n")
	}
	s.buff.Write("</nav>\n")
	return s.buff.String()
}

// AddCSS accumulates CSS with deduplication at the site level.
func (s *Site) AddCSS(css string) {
	if css == "" {
		return
	}

	s.cssBlocks = append(s.cssBlocks, assetBlock{Content: css})
}

// AddJS accumulates JavaScript with deduplication at the site level.
func (s *Site) AddJS(js string) {
	if js == "" {
		return
	}

	s.jsBlocks = append(s.jsBlocks, assetBlock{Content: js})
}

// Generate renders all site files to disk.
func (s *Site) Generate() error {

	for _, page := range s.pages {
		pagePath := PathJoin(s.Cfg.OutputDir, page.filename)
		if err := s.Cfg.WriteFile.WriteFile(pagePath, []byte(page.RenderHTML()), 0644); err != nil {
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
	s.buff.Reset()
	for _, b := range s.cssBlocks {
		s.buff.Write(b.Content)
		s.buff.Write("\n")
	}
	cssPath := PathJoin(s.Cfg.OutputDir, "style.css")
	return s.Cfg.WriteFile.WriteFile(cssPath, []byte(s.buff.String()), 0644)
}

// writeJS writes the combined and deduplicated JS to a file.
func (s *Site) writeJS() error {
	s.buff.Reset()
	for _, b := range s.jsBlocks {
		s.buff.Write(b.Content)
		s.buff.Write("\n")
	}
	jsPath := PathJoin(s.Cfg.OutputDir, "script.js")
	return s.Cfg.WriteFile.WriteFile(jsPath, []byte(s.buff.String()), 0644)
}
