# GoSite - Architecture Summary

**Updated**: 2025-10-23  
**Status**: Ready for Implementation  

---

## 🎯 Core Concept

**Hybrid Architecture**: Combines SPA (Single Page Application) with optional separate pages.

- **Default**: Modules add sections to `index.html` (SPA)
- **Optional**: Modules create separate pages (services.html, contact.html, etc.)
- **Shared Assets**: All pages use same CSS/JS
- **Auto-Everything**: IDs, filenames, navigation all auto-generated

---

## 🏗️ Architecture

```
Site (Singleton)
├── indexPage → index.html
│   └── sections[] (SPA sections)
├── pages[] → [services.html, contact.html, ...]
├── cssBlocks{} → style.css (shared, deduplicated)
└── jsBlocks{} → main.js (shared, deduplicated)
```

---

## 🔑 Key Components

### 1. Site Manager
```go
type Site struct {
    title     string
    outputDir string
    indexPage *Page           // SPA index.html
    pages     []*Page         // Ordered separate pages
    cssBlocks map[string]string
    jsBlocks  map[string]string
}

// Entry point
site := gosite.NewSite("Site Title", "output/dir")
```

### 2. Module Interface
```go
type Module interface {
    RenderUI(context ...any) string
}

// Returns HTML → SPA section
// Returns "" → Separate page
```

### 3. Main Flow
```go
func main() {
    pkg.UI = gosite.NewSite("Monjitas Chillán", "src/web/ui/")
    
    for _, mod := range pkg.Modules {
        pkg.UI.AddSection(mod) // Auto-detects type
    }
    
    pkg.UI.GenerateSite() // Writes all files
}
```

---

## 📋 Module Patterns

### Pattern A: SPA Section
```go
func (h *Homepage) RenderUI(context ...any) string {
    section := pkg.UI.Section("Inicio")
    section.AddCard("Title", "Desc", "icon")
    return section.Render() // ✅ Added to index.html
}
```

### Pattern B: Separate Page
```go
func (s *Services) RenderUI(context ...any) string {
    page := pkg.UI.NewPage(s, "Servicios")
    section := page.Section("Servicios") // Auto-added!
    section.AddCard("Title", "Desc", "icon")
    return "" // ✅ Creates services.html
}
```

---

## ✨ Auto-Generation Features

### 1. Section IDs
- Uses `runtime.Callers()` to detect module name
- `internal/services` → section ID: `"services"`

### 2. Page Filenames
- Uses `reflect.TypeOf()` to detect struct name
- `&Services{}` → filename: `"services.html"`

### 3. Navigation
- Auto-combines SPA sections + separate pages
- SPA sections: `<a href="#section-id">`
- Separate pages: `<a href="page.html">`

### 4. CSS/JS Deduplication
- SHA256 hash-based
- Shared across all pages
- Single `style.css` and `main.js`

---

## 🚫 Anti-Patterns

### ❌ Manual AddSection
```go
page := pkg.UI.NewPage(m, "Title")
section := pkg.UI.Section("Title")
page.AddSection(section) // ❌ Redundant!
```

### ✅ Correct - Auto-Add
```go
page := pkg.UI.NewPage(m, "Title")
section := page.Section("Title") // ✅ Auto-added!
```

### ❌ Manual HTML
```go
func (m *Module) RenderUI(context ...any) string {
    return "<div>HTML</div>" // ❌ Never!
}
```

### ✅ Correct - UI API
```go
func (m *Module) RenderUI(context ...any) string {
    section := pkg.UI.Section("Title")
    section.AddCard(...) // ✅ Use components!
    return section.Render()
}
```

---

## 📂 File Output

```
src/web/ui/
├── index.html       ← SPA with embedded sections
├── services.html    ← Separate page
├── contact.html     ← Separate page
├── style.css        ← Shared (deduplicated)
└── main.js          ← Shared (deduplicated)
```

---

## 🔄 Development Workflow

1. **Edit Go Module** → golite detects change
2. **Server Restarts** → `GenerateSite()` runs
3. **HTML/CSS/JS Generated** → golite detects change
4. **Minified to public/** → Browser auto-reloads

---

## 🎨 Base Styles

Initial CSS framework based on `monjitaschillan.cl`:
- Variables (colors, fonts, spacing)
- Reset
- Navigation
- Cards
- Forms
- Carousel
- Utilities

See: [`docs/base-styles.css`](./base-styles.css)

---

## 📝 Implementation Checklist

- [ ] Implement `Site` struct
- [ ] Implement `NewSite()` constructor
- [ ] Implement `AddSection()` with reflection
- [ ] Implement `NewPage()` with struct name detection
- [ ] Implement `Page.Section()` with auto-add
- [ ] Implement `detectModuleID()` with runtime.Callers
- [ ] Implement `buildCombinedNav()`
- [ ] Implement hash-based CSS/JS deduplication
- [ ] Implement `GenerateSite()` multi-file output
- [ ] Copy base styles to package
- [ ] Write unit tests
- [ ] Benchmark performance
- [ ] Document all public APIs

---

## 📚 Related Documents

- [UI Component System](./issues/ui_component_system.md) - Detailed implementation
- [Final Decisions](./issues/decisions_final.md) - All approved decisions
- [Module System](./issues/module_system.md) - Module registry & reflection
- [Migration Strategy](./issues/migration_strategy.md) - Step-by-step migration plan

---

## 🎯 Success Criteria

- ✅ Zero manual configuration (IDs, filenames, nav)
- ✅ Module decides SPA vs separate page
- ✅ All pages share CSS/JS (deduplicated)
- ✅ Navigation auto-synchronized
- ✅ Pure Go (no templates)
- ✅ Type-safe with builder pattern
- ✅ Follows Go best practices
