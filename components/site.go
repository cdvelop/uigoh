
package components

import (
    "fmt"
    "os"
    "path/filepath"
    "reflect"
    "runtime"
    "strings"
    "time"
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
    site := &Site{
        title:     title,
        outputDir: outputDir,
        pages:     make([]*Page, 0),
        cssBlocks: make(map[string]string),
        cssOrder:  make([]string, 0),
        jsBlocks:  make(map[string]string),
        jsOrder:   make([]string, 0),
    }
    site.indexPage = NewPage(title)
    site.indexPage.site = site
    return site
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
