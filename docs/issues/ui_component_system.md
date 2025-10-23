# UI Component System Design

**Parent**: [Project Structure](./project_structure.md)  
**Status**: Proposal  
**Created**: 2025-10-23  

---

## 🎨 Overview

This document details the technical implementation of the UI component system using **pure Go with String Builders** and **Builder Pattern**.

---

## 🧩 Component Architecture

### Core Principle
Every UI component is a **Go function** that returns **strings** (HTML, CSS, JS).

### Component Structure

```go
// Example: card.go
package gosite

import "strings"

type CardConfig struct {
    Title       string
    Description string
    Icon        string
    CSSClass    string
}

// RenderHTML returns the card HTML structure
func (c *CardConfig) RenderHTML() string {
    var b strings.Builder
    
    b.WriteString("<div class=\"card")
    if c.CSSClass != "" {
        b.WriteString(" ")
        b.WriteString(c.CSSClass)
    }
    b.WriteString("\">\n")
    
    if c.Icon != "" {
        b.WriteString("  <svg class=\"icon\"><use href=\"icons.svg#")
        b.WriteString(c.Icon)
        b.WriteString("\"></use></svg>\n")
    }
    
    b.WriteString("  <h3>")
    b.WriteString(escapeHTML(c.Title))
    b.WriteString("</h3>\n")
    
    b.WriteString("  <p>")
    b.WriteString(escapeHTML(c.Description))
    b.WriteString("</p>\n")
    
    b.WriteString("</div>\n")
    
    return b.String()
}

// RenderCSS returns the card CSS styles
func (c *CardConfig) RenderCSS() string {
    return `.card {
  border: 1px solid var(--color-border);
  border-radius: 8px;
  padding: 1.5rem;
  background: var(--color-card-bg);
  transition: transform 0.2s;
}

.card:hover {
  transform: translateY(-4px);
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
}

.card .icon {
  width: 48px;
  height: 48px;
  margin-bottom: 1rem;
}

.card h3 {
  margin: 0 0 0.5rem 0;
  color: var(--color-heading);
}

.card p {
  margin: 0;
  color: var(--color-text);
}
`
}

// RenderJS returns the card JavaScript (if needed)
func (c *CardConfig) RenderJS() string {
    // Most components won't need JS
    return ""
}
```

---

## 🏗️ Builder Pattern Implementation

### Page Builder

```go
// gosite.go
package gosite

import "strings"

```go
// page.go
package gosite

import "strings"

type Page struct {
    title      string
    filename   string            // e.g., "services.html", "contact.html"
    sections   []*SectionBuilder // Sections within this page
    navigation string            // Pre-rendered navigation (set by Site)
    site       *Site             // Back-reference to site
    head       []string          // Additional <head> content
}

func NewPage(title string) *Page {
    return &Page{
        title:    title,
        sections: make([]*SectionBuilder, 0),
        head:     make([]string, 0),
    }
}

// Section creates a section within THIS page
// Automatically adds section to page - NO need for AddSection()
func (p *Page) Section(title string) *SectionBuilder {
    section := p.site.Section(title)
    section.page = p
    
    // Auto-add section to page
    p.sections = append(p.sections, section)
    
    // Auto-accumulate section CSS to site
    p.site.AddCSS(section.RenderCSS())
    
    return section
}

// AddRawSection adds raw HTML section content to page
// Used by Site.AddSection() for SPA sections
func (p *Page) AddRawSection(html string) {
    // For SPA, just collect raw HTML
    // Navigation is handled by Site
}

// AddHead adds content to <head> section
func (p *Page) AddHead(content string) *Page {
    p.head = append(p.head, content)
    return p
}

// RenderHTML generates the complete HTML page
func (p *Page) RenderHTML() string {
    var b strings.Builder
    
    b.WriteString("<!DOCTYPE html>\n")
    b.WriteString("<html lang=\"es\">\n")
    b.WriteString("<head>\n")
    b.WriteString("  <meta charset=\"UTF-8\">\n")
    b.WriteString("  <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n")
    b.WriteString("  <title>")
    b.WriteString(escapeHTML(p.title))
    b.WriteString("</title>\n")
    b.WriteString("  <link rel=\"stylesheet\" href=\"style.css\">\n")
    
    // Additional head content
    for _, h := range p.head {
        b.WriteString("  ")
        b.WriteString(h)
        b.WriteString("\n")
    }
    
    b.WriteString("</head>\n")
    b.WriteString("<body>\n")
    
    // Navigation (set by Site.buildCombinedNav)
    if p.navigation != "" {
        b.WriteString(p.navigation)
    }
    
    // Render main content
    b.WriteString("  <main class=\"content\">\n")
    
    // Render all sections
    for _, section := range p.sections {
        b.WriteString(section.Render())
    }
    
    b.WriteString("  </main>\n")
    
    b.WriteString("  <script src=\"main.js\"></script>\n")
    b.WriteString("</body>\n")
    b.WriteString("</html>\n")
    
    return b.String()
}
```
```

---

## 🎛️ Site Manager (Singleton)

### Main Site Manager

```go
// site.go
package gosite

import (
    "fmt"
    "os"
    "path/filepath"
    "reflect"
    "strings"
)

// Site manages the entire UI generation system
// Supports both SPA (index.html) and separate pages
type Site struct {
    title     string          // Site title
    outputDir string          // Output directory for generated files
    indexPage *Page           // Main index.html (SPA)
    pages     []*Page         // Separate pages (services.html, contact.html, etc.)
    cssBlocks map[string]string // Shared CSS (deduplicated)
    cssOrder  []string          // CSS insertion order
    jsBlocks  map[string]string // Shared JS (deduplicated)
    jsOrder   []string          // JS insertion order
}

// NewSite creates a new Site manager
// This is the entry point for UI generation
func NewSite(title, outputDir string) *Site {
    return &Site{
        title:     title,
        outputDir: outputDir,
        indexPage: NewPage(title),
        pages:     make([]*Page, 0),
        cssBlocks: make(map[string]string),
        cssOrder:  make([]string, 0),
        jsBlocks:  make(map[string]string),
        jsOrder:   make([]string, 0),
    }
}

// AddSection adds a module's UI to the SPA index.html
// Accepts any module that implements RenderUI(context ...any) string
func (s *Site) AddSection(module any) *Site {
    // Use reflection to find RenderUI method
    renderMethod := reflect.ValueOf(module).MethodByName("RenderUI")
    if !renderMethod.IsValid() {
        return s // No RenderUI method, skip
    }
    
    // Call RenderUI with empty context
    result := renderMethod.Call([]reflect.Value{})
    if len(result) == 0 {
        return s
    }
    
    html := result[0].String()
    
    // Empty string means module created its own page
    // (via NewPage) - nothing to add to index
    if html == "" {
        return s
    }
    
    // Add section HTML to index page
    s.indexPage.AddRawSection(html)
    
    return s
}

// NewPage creates a separate HTML page for a module
// Automatically detects module name via reflection for filename
// Registers the page and creates nav link in index.html
func (s *Site) NewPage(module any, title string) *Page {
    // Detect module name via reflection
    moduleName := s.getModuleName(module)
    
    // Create new page
    page := NewPage(title)
    page.filename = moduleName + ".html" // Set filename
    page.site = s                        // Link back to site
    
    // Register page in site
    s.pages = append(s.pages, page)
    
    return page
}

// getModuleName extracts module name from struct using reflection
func (s *Site) getModuleName(module any) string {
    t := reflect.TypeOf(module)
    
    // Handle pointer types
    if t.Kind() == reflect.Ptr {
        t = t.Elem()
    }
    
    // Get struct name and convert to lowercase
    name := t.Name()
    return strings.ToLower(name)
}

// Section creates a section wrapper component with auto-generated ID
// Used by Page for creating sections within a page
func (s *Site) Section(title string) *SectionBuilder {
    // Auto-detect module ID from caller using reflection
    moduleID := s.detectModuleID()
    
    section := &SectionBuilder{
        title:    title,
        moduleID: moduleID,
        content:  make([]string, 0),
        site:     s,
    }
    
    return section
}

// detectModuleID uses reflection to find the calling module name
func (s *Site) detectModuleID() string {
    // Walk up the call stack to find the module name
    pc := make([]uintptr, 15)
    n := runtime.Callers(3, pc)
    frames := runtime.CallersFrames(pc[:n])
    
    for {
        frame, more := frames.Next()
        
        // Look for internal/ package
        if strings.Contains(frame.Function, "/internal/") {
            parts := strings.Split(frame.Function, "/internal/")
            if len(parts) > 1 {
                moduleParts := strings.Split(parts[1], ".")
                if len(moduleParts) > 0 {
                    moduleName := moduleParts[0]
                    return strings.ToLower(moduleName)
                }
            }
        }
        
        if !more {
            break
        }
    }
    
    // Fallback to timestamp-based ID
    return fmt.Sprintf("section-%d", time.Now().UnixNano())
}

// Card creates a card component
func (s *Site) Card(title, description, icon string) string {
    card := &CardConfig{
        Title:       title,
        Description: description,
        Icon:        icon,
    }
    
    // Accumulate CSS in site-wide collection
    s.AddCSS(card.RenderCSS())
    
    return card.RenderHTML()
}

// Carousel creates a carousel component
func (s *Site) Carousel(images []CarouselImage) string {
    carousel := &CarouselConfig{Images: images}
    
    s.AddCSS(carousel.RenderCSS())
    s.AddJS(carousel.RenderJS())
    
    return carousel.RenderHTML()
}

// Form creates a form component
func (s *Site) Form(config FormConfig) string {
    form := &FormBuilder{Config: config}
    
    s.AddCSS(form.RenderCSS())
    s.AddJS(form.RenderJS())
    
    return form.RenderHTML()
}

// AddCSS accumulates CSS with deduplication at site level
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

// AddJS accumulates JavaScript with deduplication at site level
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

// GenerateSite renders all files to disk
// Creates index.html (SPA) + separate pages + shared CSS/JS
func (s *Site) GenerateSite() error {
    // Ensure output directory exists
    if err := os.MkdirAll(s.outputDir, 0755); err != nil {
        return fmt.Errorf("failed to create output directory: %w", err)
    }
    
    // Generate combined navigation for index
    s.indexPage.navigation = s.buildCombinedNav()
    
    // Write index.html (SPA)
    indexPath := filepath.Join(s.outputDir, "index.html")
    if err := os.WriteFile(indexPath, []byte(s.indexPage.RenderHTML()), 0644); err != nil {
        return fmt.Errorf("failed to write index.html: %w", err)
    }
    
    // Write separate pages
    for _, page := range s.pages {
        page.navigation = s.buildCombinedNav() // Same nav for all pages
        pagePath := filepath.Join(s.outputDir, page.filename)
        if err := os.WriteFile(pagePath, []byte(page.RenderHTML()), 0644); err != nil {
            return fmt.Errorf("failed to write %s: %w", page.filename, err)
        }
    }
    
    // Write shared CSS (used by all pages)
    cssPath := filepath.Join(s.outputDir, "style.css")
    if err := os.WriteFile(cssPath, []byte(s.RenderCSS()), 0644); err != nil {
        return fmt.Errorf("failed to write CSS: %w", err)
    }
    
    // Write shared JS (used by all pages)
    jsPath := filepath.Join(s.outputDir, "main.js")
    if err := os.WriteFile(jsPath, []byte(s.RenderJS()), 0644); err != nil {
        return fmt.Errorf("failed to write JS: %w", err)
    }
    
    return nil
}

// buildCombinedNav creates navigation mixing SPA sections + separate pages
func (s *Site) buildCombinedNav() string {
    var b strings.Builder
    
    b.WriteString("<nav class=\"main-nav\">\n")
    
    // SPA sections (anchor links #section-id)
    for _, section := range s.indexPage.sections {
        navItem := section.GetNavItem()
        
        b.WriteString("  <a href=\"")
        b.WriteString(escapeAttr(navItem.Href))
        b.WriteString("\" class=\"nav-link\">\n")
        
        if navItem.Icon != "" {
            b.WriteString("    <svg class=\"icon\"><use href=\"icons.svg#")
            b.WriteString(navItem.Icon)
            b.WriteString("\"></use></svg>\n")
        }
        
        b.WriteString("    <span>")
        b.WriteString(escapeHTML(navItem.Label))
        b.WriteString("</span>\n")
        b.WriteString("  </a>\n")
    }
    
    // Separate pages (file links page.html)
    for _, page := range s.pages {
        b.WriteString("  <a href=\"")
        b.WriteString(escapeAttr(page.filename))
        b.WriteString("\" class=\"nav-link\">\n")
        
        b.WriteString("    <span>")
        b.WriteString(escapeHTML(page.title))
        b.WriteString("</span>\n")
        b.WriteString("  </a>\n")
    }
    
    b.WriteString("</nav>\n")
    
    return b.String()
}

// RenderCSS generates shared CSS for all pages
func (s *Site) RenderCSS() string {
    var b strings.Builder
    
    b.WriteString("/* Generated CSS from UI Components */\n")
    b.WriteString("/* Shared across all pages */\n")
    b.WriteString("/* Deduplicated - each block appears only once */\n\n")
    
    for i, hash := range s.cssOrder {
        if i > 0 {
            b.WriteString("\n/* Component Separator */\n\n")
        }
        b.WriteString(s.cssBlocks[hash])
    }
    
    return b.String()
}

// RenderJS generates shared JavaScript for all pages
func (s *Site) RenderJS() string {
    var b strings.Builder
    
    b.WriteString("// Generated JavaScript from UI Components\n")
    b.WriteString("// Shared across all pages\n")
    b.WriteString("// Deduplicated - each block appears only once\n\n")
    
    for i, hash := range s.jsOrder {
        if i > 0 {
            b.WriteString("\n// Component Separator\n\n")
        }
        b.WriteString(s.jsBlocks[hash])
    }
    
    return b.String()
}
```

---

## 🏗️ Section Builder

### Section Component with Auto-ID Generation

```go
// section.go
package gosite

import (
    "strings"
    "reflect"
    "runtime"
)

type SectionBuilder struct {
    title    string
    moduleID string // Auto-generated from caller module
    content  []string
    site     *Site  // Reference to site for component methods
    page     *Page  // Reference to parent page
}

// Render generates the section HTML with auto-generated ID
func (s *SectionBuilder) Render() string {
    var b strings.Builder
    
    b.WriteString("<section id=\"")
    b.WriteString(s.moduleID) // Auto-generated ID
    b.WriteString("\" class=\"page\">\n")
    
    if s.title != "" {
        b.WriteString("  <h1>")
        b.WriteString(escapeHTML(s.title))
        b.WriteString("</h1>\n")
    }
    
    b.WriteString("  <div class=\"card-container\">\n")
    for _, item := range s.content {
        b.WriteString("    ")
        b.WriteString(item)
        b.WriteString("\n")
    }
    b.WriteString("  </div>\n")
    
    b.WriteString("</section>\n")
    
    return b.String()
}

// GetNavItem returns navigation item data for auto-nav generation
func (s *SectionBuilder) GetNavItem() NavItem {
    return NavItem{
        Label: s.title,
        Href:  "#" + s.moduleID,
        Icon:  s.detectIcon(), // Optional: auto-detect icon
    }
}

// detectIcon attempts to auto-detect icon from module name
func (s *SectionBuilder) detectIcon() string {
    // Map common module names to icons
    iconMap := map[string]string{
        "homepage":     "icon-home",
        "servicespage": "icon-service",
        "staffpage":    "icon-staff",
        "contactpage":  "icon-contact",
    }
    
    if icon, ok := iconMap[s.moduleID]; ok {
        return icon
    }
    
    return "" // No icon
}

// AddCard adds a card to the section
func (s *SectionBuilder) AddCard(title, description, icon string) *SectionBuilder {
    card := s.site.Card(title, description, icon)
    s.content = append(s.content, card)
    return s
}

// AddCarousel adds a carousel to the section
func (s *SectionBuilder) AddCarousel(images []CarouselImage) *SectionBuilder {
    carousel := s.site.Carousel(images)
    s.content = append(s.content, carousel)
    return s
}

// AddForm adds a form to the section
func (s *SectionBuilder) AddForm(config FormConfig) *SectionBuilder {
    form := s.site.Form(config)
    s.content = append(s.content, form)
    return s
}

// RenderCSS returns section CSS
func (s *SectionBuilder) RenderCSS() string {
    return `.page {
  min-height: 100vh;
  padding: 2rem;
}

.card-container {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 1.5rem;
  max-width: 1200px;
  margin: 0 auto;
}
`
}
```

---

## 🧭 Navigation Component

```go
// nav.go
package gosite

import "strings"

type NavItem struct {
    Label string
    Href  string
    Icon  string
}

type NavConfig struct {
    Items []NavItem
}

func (n *NavConfig) RenderHTML() string {
    var b strings.Builder
    
    b.WriteString("<nav class=\"main-nav\">\n")
    
    for _, item := range n.Items {
        b.WriteString("  <a href=\"")
        b.WriteString(escapeAttr(item.Href))
        b.WriteString("\" class=\"nav-link\">\n")
        
        if item.Icon != "" {
            b.WriteString("    <svg class=\"icon\"><use href=\"icons.svg#")
            b.WriteString(item.Icon)
            b.WriteString("\"></use></svg>\n")
        }
        
        b.WriteString("    <span>")
        b.WriteString(escapeHTML(item.Label))
        b.WriteString("</span>\n")
        b.WriteString("  </a>\n")
    }
    
    b.WriteString("</nav>\n")
    
    return b.String()
}

func (n *NavConfig) RenderCSS() string {
    return `.main-nav {
  display: flex;
  gap: 1rem;
  padding: 1rem;
  background: var(--color-nav-bg);
  border-bottom: 1px solid var(--color-border);
}

.nav-link {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 1rem;
  text-decoration: none;
  color: var(--color-text);
  border-radius: 4px;
  transition: background 0.2s;
}

.nav-link:hover {
  background: var(--color-nav-hover);
}

.nav-link .icon {
  width: 20px;
  height: 20px;
}
`
}

func (n *NavConfig) RenderJS() string {
    return `// Navigation mobile toggle (if needed)
document.addEventListener('DOMContentLoaded', function() {
    // Add mobile menu logic here
});
`
}
```

---

## 🎠 Carousel Component

```go
// carousel.go
package gosite

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
    
    items[current].classList.add('active');
    
    setInterval(function() {
        items[current].classList.remove('active');
        current = (current + 1) % items.length;
        items[current].classList.add('active');
    }, 3000);
})();
`
}
```

---

## 📋 Form Component

```go
// form.go
package gosite

import "strings"

type FormField struct {
    Type        string // text, email, textarea
    Name        string
    Placeholder string
    Required    bool
}

type FormConfig struct {
    Action string
    Method string
    Fields []FormField
}

type FormBuilder struct {
    Config FormConfig
}

func (f *FormBuilder) RenderHTML() string {
    var b strings.Builder
    
    b.WriteString("<form class=\"contact-form\" action=\"")
    b.WriteString(escapeAttr(f.Config.Action))
    b.WriteString("\" method=\"")
    b.WriteString(f.Config.Method)
    b.WriteString("\">\n")
    
    for _, field := range f.Config.Fields {
        if field.Type == "textarea" {
            b.WriteString("  <textarea name=\"")
            b.WriteString(escapeAttr(field.Name))
            b.WriteString("\" placeholder=\"")
            b.WriteString(escapeAttr(field.Placeholder))
            b.WriteString("\"")
            if field.Required {
                b.WriteString(" required")
            }
            b.WriteString("></textarea>\n")
        } else {
            b.WriteString("  <input type=\"")
            b.WriteString(field.Type)
            b.WriteString("\" name=\"")
            b.WriteString(escapeAttr(field.Name))
            b.WriteString("\" placeholder=\"")
            b.WriteString(escapeAttr(field.Placeholder))
            b.WriteString("\"")
            if field.Required {
                b.WriteString(" required")
            }
            b.WriteString(">\n")
        }
    }
    
    b.WriteString("  <button type=\"submit\">Enviar Mensaje</button>\n")
    b.WriteString("</form>\n")
    
    return b.String()
}

func (f *FormBuilder) RenderCSS() string {
    return `.contact-form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  max-width: 500px;
}

.contact-form input,
.contact-form textarea {
  padding: 0.75rem;
  border: 1px solid var(--color-border);
  border-radius: 4px;
  font-family: inherit;
}

.contact-form button {
  padding: 0.75rem 1.5rem;
  background: var(--color-primary);
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.contact-form button:hover {
  opacity: 0.9;
}
`
}

func (f *FormBuilder) RenderJS() string {
    return `// Form validation
document.addEventListener('DOMContentLoaded', function() {
    const forms = document.querySelectorAll('.contact-form');
    forms.forEach(form => {
        form.addEventListener('submit', function(e) {
            // Add validation logic here
        });
    });
});
`
}
```

---

## 🔐 Security & Utilities

### HTML Escaping

```go
// utils.go
package gosite

import (
    "crypto/sha256"
    "encoding/hex"
    "html"
)

func escapeHTML(s string) string {
    return html.EscapeString(s)
}

func escapeAttr(s string) string {
    // More strict escaping for attributes
    return html.EscapeString(s)
}

// hashString creates a SHA256 hash of the input string
// Used for CSS/JS deduplication
func hashString(s string) string {
    h := sha256.New()
    h.Write([]byte(s))
    return hex.EncodeToString(h.Sum(nil))
}
```

---

## 📦 Component Catalog

### Planned Components

| Component | File | HTML | CSS | JS |
|-----------|------|------|-----|-----|
| Navigation | `nav.go` | ✅ | ✅ | ✅ (mobile toggle) |
| Card | `card.go` | ✅ | ✅ | ❌ |
| Carousel | `carousel.go` | ✅ | ✅ | ✅ (auto-slide) |
| Form | `form.go` | ✅ | ✅ | ✅ (validation) |
| Section | `section.go` | ✅ | ✅ | ❌ |
| Button | `button.go` | ✅ | ✅ | ❌ |
| Modal | `modal.go` | ✅ | ✅ | ✅ (open/close) |
| Table | `table.go` | ✅ | ✅ | ✅ (sorting) |

---

## 🔄 System Flow

### Architecture Overview

```
Site (Singleton)
├── indexPage (*Page)           → index.html (SPA)
│   └── sections []Section      → Embedded sections
├── pages []*Page               → Separate pages (ordered)
│   ├── services.html
│   ├── contact.html
│   └── about.html
├── cssBlocks (deduplicated)    → style.css (shared)
└── jsBlocks (deduplicated)     → main.js (shared)
```

### Generation Flow

```
1. main.go creates Site
   pkg.UI = gosite.NewSite("Site Title", "output/dir")
   ↓
2. main.go iterates modules
   for _, mod := range pkg.Modules {
       pkg.UI.AddSection(mod)
   }
   ↓
3. AddSection calls mod.RenderUI()
   - If returns HTML → add to index.html (SPA section)
   - If returns "" → module created separate page
   ↓
4. Module creates components
   section.AddCard(title, desc, icon)
   ↓
5. Components accumulate CSS/JS to Site
   site.AddCSS(card.RenderCSS())
   site.AddJS(carousel.RenderJS())
   ↓
6. main.go calls GenerateSite()
   pkg.UI.GenerateSite()
   ↓
7. Site generates combined navigation
   - SPA sections: #section-id
   - Separate pages: page.html
   ↓
8. Files written to disk
   - index.html (SPA with sections)
   - services.html (separate page)
   - contact.html (separate page)
   - style.css (shared, deduplicated)
   - main.js (shared, deduplicated)
   ↓
9. golite watches src/web/ui/
   - Detects changes
   - Minifies to src/web/public/
   ↓
10. appserver serves from public/
    - index.html → SPA with client routing
    - services.html → Standalone page
    - Shared CSS/JS cached
```

### Decision Tree: SPA vs Separate Page

```
Module.RenderUI() called
│
├─ Wants SPA section?
│  │
│  ├─ section := pkg.UI.Section("Title")
│  ├─ section.AddCard(...)
│  └─ return section.Render() ✅ Added to index.html
│
└─ Wants separate page?
   │
   ├─ page := pkg.UI.NewPage(module, "Title")
   ├─ section := page.Section("Title") ← Auto-added!
   ├─ section.AddCard(...)
   └─ return "" ✅ Creates page.html + nav link
```

---

## 🎯 Usage Examples

### ✅ Example 1: SPA Section (index.html)

```go
// src/internal/homepage/homepage.go
package homepage

import "github.com/yourorg/project/pkg"

type Homepage struct{}

type Feature struct {
    Title       string
    Description string
    Icon        string
}

func (h *Homepage) GetFeatures() []Feature {
    // Business logic only - returns data
    return []Feature{
        {Title: "Atención Personalizada", Description: "Cuidado dedicado", Icon: "icon-care"},
        {Title: "Equipo Profesional", Description: "Personal capacitado", Icon: "icon-staff"},
        {Title: "Instalaciones Modernas", Description: "Espacios cómodos", Icon: "icon-building"},
    }
}

// RenderUI returns section HTML for SPA
func (h *Homepage) RenderUI(context ...any) string {
    section := pkg.UI.Section("Inicio")
    
    // Add feature cards using UI API
    features := h.GetFeatures()
    for _, feat := range features {
        section.AddCard(feat.Title, feat.Description, feat.Icon)
    }
    
    // Return section HTML - will be added to index.html
    return section.Render()
}
```

### ✅ Example 2: Separate Page (services.html)

```go
// src/internal/services/services.go
package services

import "github.com/yourorg/project/pkg"

type Services struct{}

type Service struct {
    Title       string
    Description string
    Icon        string
}

func (s *Services) GetServices() []Service {
    // Business logic only - returns data
    return []Service{
        {Title: "Medicina General", Description: "Atención primaria", Icon: "icon-service"},
        {Title: "Curaciones", Description: "Manejo de heridas", Icon: "icon-service"},
        {Title: "Control de Signos", Description: "Monitoreo vital", Icon: "icon-service"},
    }
}

// RenderUI creates separate page - returns empty string
func (s *Services) RenderUI(context ...any) string {
    // Create separate page - filename auto-detected as "services.html"
    page := pkg.UI.NewPage(s, "Nuestros Servicios")
    
    // Create section within page - NO need for AddSection!
    // page.Section() automatically adds section to page
    section := page.Section("Servicios Disponibles")
    
    // Add service cards
    services := s.GetServices()
    for _, svc := range services {
        section.AddCard(svc.Title, svc.Description, svc.Icon)
    }
    
    // Return empty string - page already registered
    // Link auto-added to index.html navigation
    return ""
}
```

### ✅ Example 3: main.go Integration

```go
// src/cmd/appserver/main.go
package main

import (
    "log"
    "github.com/yourorg/project/pkg"
    "github.com/yourorg/project/internal/homepage"
    "github.com/yourorg/project/internal/services"
    "github.com/yourorg/project/internal/contact"
)

func main() {
    // Initialize UI system
    pkg.UI = gosite.NewSite("Monjitas Chillán", "src/web/ui/")
    
    // Register modules
    pkg.Modules = []any{
        &homepage.Homepage{},    // Will be section in index.html
        &services.Services{},    // Will create services.html
        &contact.Contact{},      // Will create contact.html
    }
    
    // Add each module to site
    for _, mod := range pkg.Modules {
        pkg.UI.AddSection(mod) // Auto-detects if section or page
    }
    
    // Generate all files
    // - index.html (with homepage section)
    // - services.html
    // - contact.html
    // - style.css (shared)
    // - main.js (shared)
    if err := pkg.UI.GenerateSite(); err != nil {
        log.Fatal(err)
    }
    
    // Start server
    // ...
}
```

---

### ❌ WRONG Examples

#### ❌ Wrong 1: Module Creating HTML Directly
```go
// ❌ NEVER DO THIS - Module should NOT know about HTML
func (m *Module) RenderUI(context ...any) string {
    var html strings.Builder
    html.WriteString("<section id=\"servicios\" class=\"page\">\n")  // ❌
    html.WriteString("  <h1>Nuestros Servicios</h1>\n")            // ❌
    html.WriteString("  <div class=\"card-container\">\n")          // ❌
    
    return html.String()
}
```

#### ❌ Wrong 2: Calling AddSection on Page
```go
// ❌ REDUNDANT - page.Section() already adds section
func (s *Services) RenderUI(context ...any) string {
    page := pkg.UI.NewPage(s, "Servicios")
    
    section := pkg.UI.Section("Servicios")
    section.AddCard(...)
    
    page.AddSection(section) // ❌ NOT NEEDED!
    
    return ""
}

// ✅ CORRECT - page.Section() auto-adds
func (s *Services) RenderUI(context ...any) string {
    page := pkg.UI.NewPage(s, "Servicios")
    
    section := page.Section("Servicios") // ✅ Auto-added!
    section.AddCard(...)
    
    return ""
}
```

#### ❌ Wrong 3: Manual ID Specification
```go
// ❌ NO - ID is auto-generated via reflection
section := pkg.UI.Section("servicios", "Title") // ❌ No ID param!

// ✅ CORRECT - ID auto-detected from module name
section := pkg.UI.Section("Title") // ✅ ID = "services" (from struct name)
```

#### ❌ Wrong 4: Manual Navigation Creation
```go
// ❌ NO - Navigation is auto-generated by Site
nav := pkg.UI.Nav([]pkg.NavItem{...}) // ❌ Never call this!

// ✅ CORRECT - Site.buildCombinedNav() handles this automatically
```

#### ❌ Wrong 5: Returning HTML from Separate Page
```go
// ❌ NO - Should return empty string
func (s *Services) RenderUI(context ...any) string {
    page := pkg.UI.NewPage(s, "Servicios")
    section := page.Section("Servicios")
    
    return page.RenderHTML() // ❌ Return "" instead!
}

// ✅ CORRECT
func (s *Services) RenderUI(context ...any) string {
    page := pkg.UI.NewPage(s, "Servicios")
    section := page.Section("Servicios")
    
    return "" // ✅ Page already registered
}
```

**Why are these wrong?**
- ❌ Violates separation of concerns
- ❌ Duplicates framework functionality
- ❌ Creates maintenance burden
- ❌ Breaks auto-generation features
- ❌ CSS/JS not properly accumulated
- ❌ No HTML escaping (security risk)
- ❌ Tightly couples business logic to UI structure

---

## ⚡ Performance Considerations

### String Builder Efficiency
- ✅ Preallocate capacity when size is known: `strings.Builder.Grow(n)`
- ✅ Reuse builders with `.Reset()`
- ✅ Avoid concatenation with `+`, always use `.WriteString()`

### CSS/JS Deduplication
- **Required**: Hash-based deduplication to avoid duplicates
- Use map with content hash as key
- Maintain insertion order with separate slice

### Caching Strategy
- **Development**: Regenerate on every request
- **Production**: Generate once at startup

---

## 🧪 Testing Strategy

### Unit Tests per Component

```go
// card_test.go
package gosite

import "testing"

func TestCardRenderHTML(t *testing.T) {
    card := &CardConfig{
        Title:       "Test Card",
        Description: "Test description",
        Icon:        "icon-test",
    }
    
    html := card.RenderHTML()
    
    if !strings.Contains(html, "Test Card") {
        t.Error("Card HTML should contain title")
    }
    
    if !strings.Contains(html, "icon-test") {
        t.Error("Card HTML should contain icon reference")
    }
}

func TestCardEscaping(t *testing.T) {
    card := &CardConfig{
        Title:       "<script>alert('xss')</script>",
        Description: "Safe",
    }
    
    html := card.RenderHTML()
    
    if strings.Contains(html, "<script>") {
        t.Error("Card should escape HTML in title")
    }
}
```

---

## 📝 Component Development Guidelines

### 1. Always Escape User Input
```go
b.WriteString(escapeHTML(userInput))
```

### 2. Provide Sensible Defaults
```go
type ButtonConfig struct {
    Label   string
    Type    string // default: "button"
    CSSClass string // default: ""
}

func (b *ButtonConfig) SetDefaults() {
    if b.Type == "" {
        b.Type = "button"
    }
}
```

### 3. Document Component Config
```go
// CardConfig defines the configuration for a card component.
//
// Fields:
//   - Title: Main heading (required, escaped)
//   - Description: Body text (required, escaped)
//   - Icon: SVG icon ID from icons.svg (optional)
//   - CSSClass: Additional CSS classes (optional)
type CardConfig struct { ... }
```

### 4. Keep Components Focused
- One component = one UI pattern
- Avoid "god components"
- Compose complex UIs from simple components

---

## � Key Architecture Decisions

### ✅ Final Confirmed Decisions

1. **Site/Page Architecture**
   - `Site` manages entire system (singleton)
   - `indexPage` for SPA sections
   - `pages []*Page` for separate pages (slice maintains order)
   - Shared CSS/JS across all pages

2. **Module Interface**
   - `RenderUI(context ...any) string`
   - Returns section HTML → added to index.html
   - Returns "" → separate page created

3. **Auto-Detection**
   - Section IDs via reflection from module name
   - Page filenames via reflection from struct name
   - Navigation auto-generated mixing sections + pages

4. **No Redundancy**
   - `page.Section()` auto-adds section to page
   - NO manual `AddSection()` calls in modules
   - NO manual navigation construction

5. **CSS/JS Shared**
   - Single `style.css` used by all pages
   - Single `main.js` used by all pages
   - Hash-based deduplication at Site level

6. **File Structure**
   ```
   src/web/ui/
   ├── index.html      (SPA with sections)
   ├── services.html   (separate page)
   ├── contact.html    (separate page)
   ├── style.css       (shared)
   └── main.js         (shared)
   ```

---

## 🚀 Next Steps

1. ✅ Architecture documented
2. ✅ Base CSS created (`docs/base-styles.css`)
3. ⏳ Implement `Site` struct
4. ⏳ Implement `Page` struct with `Section()` method
5. ⏳ Implement reflection-based filename detection
6. ⏳ Implement `buildCombinedNav()`
7. ⏳ Test with sample modules
8. ⏳ Add remaining components (Modal, Table, Button)

---

**See also**:
- [Module System](./module_system.md) - Integration with internal modules
- [Migration Strategy](./migration_strategy.md) - Implementation steps
- [Final Decisions](./decisions_final.md) - All approved decisions
- [Base Styles](../base-styles.css) - Initial CSS framework
