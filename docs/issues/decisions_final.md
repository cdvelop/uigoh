# Final Architecture Decisions

**Parent**: [Project Structure](./project_structure.md)  
**Status**: **APPROVED - Ready for Implementation**  
**Created**: 2025-10-23  
**Updated**: 2025-10-23  

---

## 🎯 Architecture Overview

### Site/Page Model
- **Site**: Singleton managing entire UI generation system
- **SPA**: Single `index.html` with embedded sections
- **Separate Pages**: Optional standalone pages (services.html, contact.html, etc.)
- **Shared Assets**: All pages use same `style.css` and `main.js`

---

## ✅ All Decisions Finalized

### 1️⃣ Site Architecture: **Hybrid SPA + Multi-Page**
- Default: Modules add sections to SPA `index.html`
- Optional: Modules create separate pages
- Ordered pages using **slice** (`[]*Page`), not map
- Auto-generated navigation combines both types

**Structure**:
```go
type Site struct {
    indexPage *Page           // SPA index.html
    pages     []*Page         // Ordered separate pages
    cssBlocks map[string]string // Shared CSS
    jsBlocks  map[string]string // Shared JS
}
```

---

### 2️⃣ RenderUI Signature: **Variadic Context**
```go
func (m *Module) RenderUI(context ...any) string
```

**Rationale**:
- Maximum flexibility
- Each module handles its own requirements internally
- No coupling to specific context types
- Module decides what it needs from context

**Return Values**:
- Returns section HTML → Added to `index.html` (SPA)
- Returns `""` → Module created separate page

**Usage Examples**:

```go
// Example 1: SPA Section (returns HTML)
func (h *Homepage) RenderUI(context ...any) string {
    section := pkg.UI.Section("Inicio")
    section.AddCard("Title", "Description", "icon")
    return section.Render() // ✅ Added to index.html
}

// Example 2: Separate Page (returns "")
func (s *Services) RenderUI(context ...any) string {
    page := pkg.UI.NewPage(s, "Servicios")
    section := page.Section("Servicios") // ✅ Auto-added!
    section.AddCard("Title", "Description", "icon")
    return "" // ✅ Creates services.html
}

// Example 3: With Context Extraction
func (p *PatientCare) RenderUI(context ...any) string {
    var user UserInfo
    if len(context) > 0 {
        if u, ok := context[0].(UserInfo); ok {
            user = u
        }
    }
    
    section := pkg.UI.Section("Patient Care")
    section.AddCard(user.Name, user.Status, "icon-user")
    return section.Render()
}
```

---

### 3️⃣ AddSection Interface: **Reflection-Based**

**Signature**:
```go
func (s *Site) AddSection(module any) *Site
```

**How it works**:
1. Accepts any type
2. Uses reflection to find `RenderUI` method
3. Calls `RenderUI()` with empty context
4. If returns HTML → adds to index.html
5. If returns "" → module created separate page (skip)

**Implementation**:
```go
func (s *Site) AddSection(module any) *Site {
    renderMethod := reflect.ValueOf(module).MethodByName("RenderUI")
    if !renderMethod.IsValid() {
        return s // No RenderUI, skip
    }
    
    result := renderMethod.Call([]reflect.Value{})
    html := result[0].String()
    
    if html == "" {
        return s // Separate page, already registered
    }
    
    s.indexPage.AddRawSection(html)
    return s
}
```

---

### 4️⃣ Auto-Generated IDs and Navigation: **Convention-Based**

**Core Principle**: Everything auto-generated - no manual IDs or navigation.

**Section ID Detection**:
```go
// Using runtime.Callers to detect module name
func (s *Site) detectModuleID() string {
    pc := make([]uintptr, 15)
    n := runtime.Callers(3, pc)
    frames := runtime.CallersFrames(pc[:n])
    
    for {
        frame, more := frames.Next()
        if strings.Contains(frame.Function, "/internal/") {
            // Extract module name: "internal/services" → "services"
            return extractModuleName(frame.Function)
        }
        if !more {
            break
        }
    }
}
```

**Page Filename Detection**:
```go
// Using reflection on module struct
func (s *Site) NewPage(module any, title string) *Page {
    // Detect: &services.Services{} → "services.html"
    t := reflect.TypeOf(module)
    if t.Kind() == reflect.Ptr {
        t = t.Elem()
    }
    filename := strings.ToLower(t.Name()) + ".html"
    
    page := NewPage(title)
    page.filename = filename
    s.pages = append(s.pages, page)
    
    return page
}
```

**Auto-Navigation Generation**:
```go
// Combined navigation: SPA sections + separate pages
func (s *Site) buildCombinedNav() string {
    var nav strings.Builder
    
    // SPA sections: #section-id (anchor links)
    for _, section := range s.indexPage.sections {
        nav.WriteString("<a href=\"#" + section.moduleID + "\">")
        nav.WriteString(section.title)
        nav.WriteString("</a>")
    }
    
    // Separate pages: page.html (file links)
    for _, page := range s.pages {
        nav.WriteString("<a href=\"" + page.filename + "\">")
        nav.WriteString(page.title)
        nav.WriteString("</a>")
    }
    
    return nav.String()
}
```

**Rationale**: 
- ✅ Zero manual configuration
- ✅ Impossible ID/filename conflicts
- ✅ Navigation always synchronized
- ✅ Single source of truth

---

### 5️⃣ Page.Section Auto-Add: **No Redundancy**

**Key Decision**: `page.Section()` automatically adds section to page.

```go
// ❌ OLD WAY - Redundant
page := pkg.UI.NewPage(m, "Services")
section := pkg.UI.Section("Services")
page.AddSection(section) // ❌ Manual step

// ✅ NEW WAY - Automatic
page := pkg.UI.NewPage(m, "Services")
section := page.Section("Services") // ✅ Auto-added!
```

**Implementation**:
```go
func (p *Page) Section(title string) *SectionBuilder {
    section := p.site.Section(title)
    section.page = p
    
    // Auto-add to page
    p.sections = append(p.sections, section)
    
    // Auto-accumulate CSS
    p.site.AddCSS(section.RenderCSS())
    
    return section
}
```

---

### 6️⃣ CSS/JS Deduplication: **Site-Level with Hash**

**Implementation**:
```go
type Site struct {
    cssBlocks map[string]string // hash -> css content
    cssOrder  []string           // maintain insertion order
    jsBlocks  map[string]string // hash -> js content
    jsOrder   []string           // maintain insertion order
}

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

func (s *Site) RenderCSS() string {
    var b strings.Builder
    for i, hash := range s.cssOrder {
        if i > 0 {
            b.WriteString("\n/* Separator */\n\n")
        }
        b.WriteString(s.cssBlocks[hash])
    }
    return b.String()
}
```

**Key Points**:
- ✅ Shared across all pages (index.html, services.html, etc.)
- ✅ Single `style.css` file
- ✅ Single `main.js` file
- ✅ Hash-based deduplication
- ✅ Insertion order maintained

---

### 7️⃣ Slice for Pages: **Ordered, Not Random**

**Decision**: Use `[]*Page` instead of `map[string]*Page`

```go
type Site struct {
    pages []*Page // ✅ Maintains insertion order
}
```

**Rationale**:
- ❌ Maps have random iteration order
- ✅ Slices maintain insertion order
- ✅ Navigation links appear in predictable order
- ✅ Better UX

---

### 8️⃣ GenerateSite Frequency: **Every Startup**

**Implementation**:
```go
func main() {
    // Initialize site
    pkg.UI = uigoh.NewSite("Title", "output/")
    
    // Register modules
    for _, mod := range pkg.Modules {
        pkg.UI.AddSection(mod)
    }
    
    // Generate all files
    pkg.UI.GenerateSite()
    
    // Start server
    // ...
}
```

**Workflow**:
1. golite detects Go file changes
2. Server restarts
3. `GenerateSite()` runs at startup
4. Writes all files (index.html + pages)
5. golite detects HTML/CSS/JS changes
6. Minifies to public/
7. Browser auto-reloads

---

### 9️⃣ No Manual HTML in Modules: **Strict Prohibition**

**Module Responsibility**: Only provide content and title

**uigoh Responsibility**: Generate IDs, navigation, structure

**Complete Example**:
```go
// Module: internal/servicesPage/ui.go
func (m *ServicesPage) RenderUI(context ...any) string {
    page := pkg.UI.NewPage("Monjitas Chillán")
    
    // Create section - NO ID needed, NO nav needed
    section := pkg.UI.Section("Nuestros Servicios")
    
    // Add content using API methods only
    services := m.GetServices() // Business logic
    for _, svc := range services {
        section.AddCard(svc.Title, svc.Description, svc.Icon)
    }
    
    // Add section to page
    page.AddSection(section)
    
    // Navigation auto-generated from all sections
    return page.RenderHTML()
}
```

**Behind the Scenes (in uigoh)**:
```go
// section.go
type SectionBuilder struct {
    title    string
    moduleID string // Auto-detected via reflection
    content  []string
    ui       *HtmlUI
}

func (ui *HtmlUI) Section(title string) *SectionBuilder {
    // Auto-detect module name using reflection
    moduleID := ui.detectModuleName() // e.g., "servicespage"
    
    section := &SectionBuilder{
        title:    title,
        moduleID: moduleID,
        content:  make([]string, 0),
        ui:       ui,
    }
    
    return section
}

// page.go
type Page struct {
    title    string
    sections []*SectionBuilder // Store section builders
}

func (p *Page) AddSection(section *SectionBuilder) *Page {
    p.sections = append(p.sections, section)
    return p
}

func (p *Page) RenderHTML() string {
    var b strings.Builder
    
    // Auto-generate navigation from sections
    b.WriteString(p.generateNavigation())
    
    // Render all sections
    for _, section := range p.sections {
        b.WriteString(section.RenderWithID()) // ID auto-injected
    }
    
    return b.String()
}

func (p *Page) generateNavigation() string {
    var nav strings.Builder
    
    nav.WriteString("<nav class=\"main-nav\">\n")
    
    for _, section := range p.sections {
        // Auto-generate nav item from section
        nav.WriteString("  <a href=\"#")
        nav.WriteString(section.moduleID) // Auto-generated ID
        nav.WriteString("\" class=\"nav-link\">\n")
        nav.WriteString("    <span>")
        nav.WriteString(escapeHTML(section.title))
        nav.WriteString("</span>\n")
        nav.WriteString("  </a>\n")
    }
    
    nav.WriteString("</nav>\n")
    
    return nav.String()
}
```

---

### 7️⃣ Module HTML Generation: **FORBIDDEN**

**❌ CRITICAL RULE**: Modules MUST NOT create HTML directly

**✅ CORRECT**:
```go
func (m *Module) RenderUI(context ...any) string {
    page := pkg.UI.NewPage("Services")
    
    // Just title, no ID
    section := pkg.UI.Section("Nuestros Servicios")
    
    services := m.GetServices() // Business logic
    for _, svc := range services {
        section.AddCard(svc.Title, svc.Description, svc.Icon)
    }
    
    page.AddSection(section)
    
    // Navigation auto-generated, no manual nav code
    return page.RenderHTML()
}
```

**❌ WRONG**:
```go
func (m *Module) RenderUI(context ...any) string {
    html := "<section id='servicios'><h1>Services</h1>" // ❌ NO!
    return html
}

// Also wrong - manual ID specification
section := pkg.UI.Section("servicios", "Title") // ❌ NO ID param!

// Also wrong - manual nav construction
nav := pkg.UI.Nav([]pkg.NavItem{...}) // ❌ Auto-generated!
```

---

## 📐 Updated Component API

### HtmlUI Public Methods

```go
// Page management
func (ui *HtmlUI) NewPage(title string) *Page

// Section creation - NO ID parameter, auto-generated
func (ui *HtmlUI) Section(title string) *SectionBuilder

// Individual components (used by Section)
func (ui *HtmlUI) Card(title, description, icon string) string
func (ui *HtmlUI) Carousel(images []CarouselImage) string
func (ui *HtmlUI) Form(config FormConfig) string
func (ui *HtmlUI) Header(text string) string
func (ui *HtmlUI) Button(label, action string) string

// NO Nav() method - navigation is auto-generated

// Rendering
func (ui *HtmlUI) GenerateSite(destinations ...string) error
```

### Page Methods

```go
// Add sections - accepts SectionBuilder
func (p *Page) AddSection(section *SectionBuilder) *Page

// RenderHTML auto-generates navigation from sections
func (p *Page) RenderHTML() string
```

### SectionBuilder Methods

```go
// Content addition - NO HTML strings
func (s *SectionBuilder) AddCard(title, desc, icon string) *SectionBuilder
func (s *SectionBuilder) AddCarousel(images []CarouselImage) *SectionBuilder
func (s *SectionBuilder) AddForm(config FormConfig) *SectionBuilder
func (s *SectionBuilder) AddCustom(componentHTML string) *SectionBuilder // Escape hatch

// NO SetID() - ID is auto-generated
// NO Render() - handled by Page
```

---

## 🗂️ File Structure (Final)

```
src/
├── pkg/
│   ├── ui.go            # var UI = uigoh.New()
│   ├── modules.go       # var Modules = []any{...}
│   ├── interfaces.go    # UserContext interface
│   │
│   └── uigoh/
│       ├── uigoh.go     # HtmlUI manager
│       ├── page.go      # Page builder (with dedup)
│       ├── utils.go     # escapeHTML, hashString
│       ├── section.go   # Section component
│       ├── nav.go       # Navigation component
│       ├── card.go      # Card component
│       ├── carousel.go  # Carousel component
│       ├── form.go      # Form component
│       ├── button.go    # Button component
│       ├── header.go    # Header component
│       └── css.go       # Global CSS
│
├── internal/
│   ├── homePage/
│   │   ├── module.go    # struct, business logic
│   │   └── ui.go        # RenderUI - ONLY calls UI API
│   │
│   ├── servicesPage/
│   │   ├── module.go
│   │   ├── service.go   # Business models
│   │   └── ui.go        # RenderUI - ONLY calls UI API
│   │
│   └── patientCare/
│       ├── module.go
│       ├── patient.go
│       ├── repository.go
│       ├── service.go
│       └── ui.go        # RenderUI - ONLY calls UI API
│
└── web/
    ├── ui/              # Generated (watched by golite)
    │   ├── index.html
    │   ├── style.css
    │   └── script.js
    │
    └── public/          # Served (minified by golite)
        ├── index.html
        ├── style.css
        ├── script.js
        ├── icons.svg
        └── img/
```

---

## 🔄 Data Flow (Corrected)

```
┌─────────────────────────────────────────────────────────┐
│ main.go                                                  │
│  - Initializes UI (pkg.UI)                              │
│  - Loads modules (pkg.Modules)                          │
│  - Calls each module's RenderUI(context...)             │
└─────────────────┬───────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────────────┐
│ internal/servicesPage/ui.go                             │
│                                                          │
│  func (m *Module) RenderUI(context ...any) string {     │
│      // Extract context internally if needed            │
│      var user UserInfo                                  │
│      if len(context) > 0 {                              │
│          if u, ok := context[0].(UserInfo); ok {        │
│              user = u                                   │
│          }                                              │
│      }                                                  │
│                                                          │
│      page := pkg.UI.NewPage("Services")                 │
│                                                          │
│      // NO ID specified - auto-generated                │
│      section := pkg.UI.Section("Nuestros Servicios")    │
│                                                          │
│      services := m.GetServices() // ✅ Business logic   │
│                                                          │
│      for _, svc := range services {                     │
│          section.AddCard(                               │
│              svc.Title,      // ✅ Data only            │
│              svc.Description,                           │
│              svc.Icon                                   │
│          )                                              │
│      }                                                  │
│                                                          │
│      page.AddSection(section) // Pass builder directly  │
│                                                          │
│      return page.RenderHTML() // Auto-generates nav     │
│  }                                                       │
│                                                          │
│  ✅ NO HTML strings                                     │
│  ✅ NO CSS                                              │
│  ✅ NO JavaScript                                       │
│  ✅ NO IDs                                              │
│  ✅ NO manual navigation                                │
│  ✅ ONLY calls pkg.UI methods                           │
└─────────────────┬───────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────────────┐
│ pkg/ui.go → pkg/uigoh/*                                 │
│  - Section() auto-generates ID from module name         │
│  - section.AddCard() calls uigoh.Card()                 │
│  - Card() generates HTML                                │
│  - Card() accumulates CSS (deduplicated)                │
│  - Page.RenderHTML() auto-generates navigation          │
└─────────────────┬───────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────────────┐
│ main.go calls pkg.UI.GenerateSite()                     │
│  - Compares with existing files                         │
│  - Writes only if changed:                              │
│    * src/web/ui/index.html (with auto-nav)              │
│    * src/web/ui/style.css (deduplicated)                │
│    * src/web/ui/script.js (deduplicated)                │
└─────────────────┬───────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────────────┐
│ golite watches src/web/ui/                              │
│  - Detects file changes                                 │
│  - Minifies files                                       │
│  - Copies to src/web/public/                            │
└─────────────────────────────────────────────────────────┘
```

---

## ⚠️ KEY ARCHITECTURE CHANGES

### 🔄 From Previous Version:

| Aspect | Before | **After** |
|--------|--------|-----------|
| RenderUI params | `ctx pkg.UserContext` | `context ...any` |
| Section creation | `Section(id, title)` | `Section(title)` only |
| Section ID | Manual | **Auto-generated** |
| Navigation | `UI.Nav(items)` manual | **Auto-generated** from sections |
| AddSection param | `section.Render()` | `section` (builder) |

### ✅ Zero Configuration Principle

**Modules only provide**:
- Section title
- Content data
- Business logic

**uigoh handles**:
- ID generation (from module name)
- Navigation generation (from sections)
- HTML structure
- CSS/JS accumulation

---

## ⚠️ NO LONGER PENDING - APPROVED

All architectural decisions are now finalized and approved:

- ✅ RenderUI uses `...any` for maximum flexibility
- ✅ Section IDs are auto-generated
- ✅ Navigation is auto-generated
- ✅ No manual configuration needed
- ✅ CSS deduplication required
- ✅ Generate only on file changes

---

## 🎯 Implementation Ready

**Status**: ✅ **ALL DECISIONS APPROVED - READY FOR IMPLEMENTATION**

Next steps:

1. ✅ Start Phase 1 (Foundation)
2. ✅ Create `src/pkg/uigoh/` structure  
3. ✅ Implement auto-ID generation via reflection
4. ✅ Implement auto-navigation generation
5. ✅ Implement components with deduplication
6. ✅ Create first module (homePage)
7. ✅ Test end-to-end

---

## 📋 Summary of Final Architecture

| Component | Responsibility |
|-----------|----------------|
| **Module** | Data + Business logic + Call UI API |
| **uigoh** | HTML/CSS/JS generation + Auto-IDs + Auto-nav |
| **Section** | Container builder, no ID knowledge |
| **Page** | Section aggregator + Nav generator |
| **GenerateSite** | Write files only on changes |

### Complete Minimal Example:

```go
// internal/homePage/ui.go
package homePage

import "github.com/cdvelop/monjitaschillan.cl/pkg"

type Module struct{}

func (m *Module) RenderUI(context ...any) string {
    section := pkg.UI.Section("Salud sin barreras")
    

    // Add carousel
    section.AddCarousel([]pkg.CarouselImage{
        {Src: "img/med-img-01.jpg", Alt: "Bienvenida 1"},
        {Src: "img/med-img-02.jpg", Alt: "Bienvenida 2"},
    })
    
    
    // Navigation auto-generated from all sections
    // HTML/CSS/JS auto-accumulated
    return section.RenderHTML()
}
```

**That's it.** No IDs, no nav, no HTML strings.

---

## 📋 Summary of Changes from Original

| Aspect | Original | **Final** |
|--------|----------|-----------|
| RenderUI params | `ctx UserContext` | `context ...any` ✅ |
| RenderUI return | `(string, error)` | `string` ✅ |
| Section creation | `Section(id, title)` | `Section(title)` ✅ |
| Section ID | Manual | **Auto-generated** ✅ |
| Navigation | Manual `Nav()` call | **Auto-generated** ✅ |
| AddSection | `section.Render()` | `section` ✅ |
| CSS Dedup | Optional | **Required** ✅ |
| Module HTML | ❌ Was allowed | **Forbidden** ✅ |
| Singleton State | Question | **Stateless** ✅ |
| GenerateSite | Question | **On changes only** ✅ |

---

**Status**: ✅ **APPROVED - Ready for implementation**
