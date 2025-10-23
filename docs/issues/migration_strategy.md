# Migration Strategy

**Parent**: [Project Structure](./project_structure.md)  
**Status**: Proposal  
**Created**: 2025-10-23  

---

## 🎯 Migration Goal

Replace current `builder` package (template-based) with new `gosite` system (component-based, pure Go) **without breaking production**.

---

## 📋 Pre-Migration Checklist

### Current State Analysis

- [x] Current system uses `html/template`
- [x] Template file: `src/web/ui/template.html`
- [x] Data models in `builder/models.go`
- [x] Single `GenerateSite()` call in `main.go`
- [x] Output to `src/web/ui/index.html`
- [x] `golite` watches `src/web/ui/` and minifies to `public/`

### Dependencies to Review

```bash
# Check what depends on builder package
grep -r "github.com/cdvelop/monjitaschillan.cl/builder" .

# Expected results:
# - src/cmd/appserver/main.go (line 9)
# - go.mod (replace directive if exists)
```

---

## 🗺️ Migration Phases

### Phase 1: Foundation Setup ✅
**Goal**: Create new package structure without breaking existing system

**Tasks**:
1. Create `src/pkg/gosite/` directory
2. Create `src/pkg/ui.go` (public API)
3. Create `src/pkg/modules.go` (empty for now)
4. Implement core gosite files:
   - `gosite.go` (HtmlUI manager)
   - `page.go` (Page builder)
   - `utils.go` (HTML escaping)

**Test**: Build succeeds, old system still works

```bash
go build ./src/cmd/appserver
```

---

### Phase 2: Component Implementation 🔨
**Goal**: Build all necessary UI components

**Order** (dependencies first):
1. ✅ `utils.go` - Escaping functions
2. ✅ `css.go` - CSS utilities
3. ✅ `js.go` - JS utilities
4. ✅ `page.go` - Page builder
5. ✅ `section.go` - Section wrapper
6. ✅ `nav.go` - Navigation component
7. ✅ `card.go` - Card component
8. ✅ `carousel.go` - Carousel component
9. ✅ `form.go` - Form component

**For each component**:
- Write implementation
- Write unit test
- Test HTML escaping
- Verify CSS/JS accumulation

**Validation**:
```bash
cd src/pkg/gosite
go test -v
```

---

### Phase 3: Create First Internal Module 🏗️
**Goal**: Prove the module system works end-to-end

**Create**: `src/internal/homePage/`

**Files**:
```
src/internal/homePage/
├── module.go      # Module struct
├── content.go     # Business logic (get carousel images, etc.)
└── ui.go          # RenderUI() method
```

**Implementation**:
```go
// module.go
package homePage

type Module struct{}

// content.go
package homePage

type CarouselImage struct {
    Src string
    Alt string
}

func (m *Module) GetCarouselImages() []CarouselImage {
    return []CarouselImage{
        {Src: "img/med-img-01.jpg", Alt: "Imagen de Bienvenida 1"},
        {Src: "img/med-img-02.jpg", Alt: "Imagen de Bienvenida 2"},
        {Src: "img/med-img-03.jpg", Alt: "Imagen de Bienvenida 3"},
    }
}

// ui.go
package homePage

import (
    "github.com/cdvelop/monjitaschillan.cl/pkg"
    "strings"
)

func (m *Module) RenderUI(params ...any) (string, error) {
    page := pkg.UI.NewPage("Monjitas Chillán")
    
    // Build carousel
    var carousel strings.Builder
    carousel.WriteString("<div class=\"carousel\">\n")
    
    for _, img := range m.GetCarouselImages() {
        carousel.WriteString("  <div class=\"carousel-item\">\n")
        carousel.WriteString("    <img src=\"")
        carousel.WriteString(img.Src)
        carousel.WriteString("\" alt=\"")
        carousel.WriteString(img.Alt)
        carousel.WriteString("\">\n")
        carousel.WriteString("  </div>\n")
    }
    
    carousel.WriteString("</div>\n")
    
    page.AddSection(carousel.String())
    
    return page.RenderHTML(), nil
}
```

**Register**:
```go
// pkg/modules.go
package pkg

import "github.com/cdvelop/monjitaschillan.cl/src/internal/homePage"

var Modules = []any{
    &homePage.Module{},
}
```

**Test standalone**:
```go
// Test file
package main

import (
    "log"
    "github.com/cdvelop/monjitaschillan.cl/pkg"
)

func main() {
    module := &homePage.Module{}
    html, err := module.RenderUI()
    if err != nil {
        log.Fatal(err)
    }
    log.Println(html)
}
```

---

### Phase 4: Parallel System Testing 🔄
**Goal**: Run new system alongside old one for comparison

**Approach**: Dual rendering

```go
// main.go (temporary dual mode)
func main() {
    // ... setup ...
    
    // OLD SYSTEM (keep for now)
    webUi := filepath.Dir(absPublicDir)
    templatePath := filepath.Join(webUi, "ui", "template.html")
    indexPath := filepath.Join(webUi, "ui", "index.html")
    
    err = builder.GenerateSite(templatePath, indexPath)
    if err != nil {
        log.Fatalf("Error generating site (OLD): %v", err)
    }
    
    // NEW SYSTEM (test output)
    testOutput := filepath.Join(webUi, "ui-new")
    os.MkdirAll(testOutput, 0755)
    
    newUI := pkg.UI
    page := newUI.NewPage("Monjitas Chillán")
    
    for _, mod := range pkg.Modules {
        if renderer, ok := mod.(interface{ RenderUI(...any) (string, error) }); ok {
            html, err := renderer.RenderUI()
            if err != nil {
                log.Printf("Error rendering module: %v", err)
                continue
            }
            page.AddSection(html)
        }
    }
    
    if err := newUI.GenerateSite(testOutput); err != nil {
        log.Printf("Error generating new site: %v", err)
    }
    
    log.Println("✅ Old system output: src/web/ui/")
    log.Println("✅ New system output: src/web/ui-new/")
    log.Println("Compare files to verify correctness")
    
    // Server still uses old output
    log.Printf("Serving from: %s", absPublicDir)
    // ... rest of server code ...
}
```

**Validation**:
```bash
# Compare outputs
diff src/web/ui/index.html src/web/ui-new/index.html
diff src/web/ui/style.css src/web/ui-new/style.css
```

---

### Phase 5: Migrate Remaining Modules 📦
**Goal**: Port all content from `template.html` to modules

**Modules to create**:
1. ✅ `homePage` (carousel, hero)
2. ⏳ `servicesPage` (service cards)
3. ⏳ `staffPage` (staff cards)
4. ⏳ `contactPage` (contact form, map)
5. ⏳ `navigation` (global nav)

**For each module**:
1. Extract data from `builder/models.go`
2. Create `internal/<module>/content.go`
3. Create `internal/<module>/ui.go`
4. Add to `pkg/modules.go`
5. Test rendering
6. Compare with old output

---

### Phase 6: CSS Migration 🎨
**Goal**: Port CSS from `src/web/ui/css/*.css` to component methods

**Current CSS structure**:
```
src/web/ui/css/
├── 010-variables.css
├── 020-reset.css
├── 030-base.css
└── ... (component styles)
```

**Strategy**:
- **Global styles** (variables, reset, base) → `gosite/global.go`
- **Component styles** → Each component's `RenderCSS()`

**Example**:
```go
// gosite/global.go
package gosite

func GlobalCSS() string {
    return `/* CSS Variables */
:root {
    --color-primary: #2563eb;
    --color-border: #e5e7eb;
    /* ... */
}

/* CSS Reset */
* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}

/* Base Styles */
body {
    font-family: system-ui, sans-serif;
    line-height: 1.6;
}
`
}
```

**Integration**:
```go
// page.go
func (p *Page) RenderCSS() string {
    var b strings.Builder
    
    // Always include global CSS first
    b.WriteString(GlobalCSS())
    b.WriteString("\n\n")
    
    // Then component-specific CSS
    for _, css := range p.cssBlocks {
        b.WriteString(css)
        b.WriteString("\n")
    }
    
    return b.String()
}
```

---

### Phase 7: Switch Cutover 🔄
**Goal**: Replace old system with new in main.go

**Changes**:
```go
// main.go - BEFORE
import (
    "github.com/cdvelop/monjitaschillan.cl/builder"
)

func main() {
    // ...
    err = builder.GenerateSite(templatePath, indexPath)
    // ...
}
```

```go
// main.go - AFTER
import (
    "github.com/cdvelop/monjitaschillan.cl/pkg"
    "github.com/cdvelop/monjitaschillan.cl/pkg/gosite"
)

func main() {
    // ... setup code ...
    
    // Discover and render modules
    modules := gosite.DiscoverModules(pkg.Modules)
    
    page := pkg.UI.NewPage("Monjitas Chillán")
    
    for _, mod := range modules {
        if mod.HasRenderUI {
            html, err := callRenderUI(mod.Instance)
            if err != nil {
                log.Printf("Error rendering %s: %v", mod.Name, err)
                continue
            }
            page.AddSection(html)
        }
    }
    
    // Generate to src/web/ui/ (same output path as before)
    outputDir := filepath.Join(filepath.Dir(absPublicDir), "ui")
    if err := pkg.UI.GenerateSite(outputDir); err != nil {
        log.Fatalf("Error generating site: %v", err)
    }
    
    log.Printf("✅ Site generated at: %s", outputDir)
    
    // golite picks up changes automatically (no changes needed)
    
    // ... rest of server code (unchanged) ...
}
```

**Testing**:
```bash
# Build and run
go run ./src/cmd/appserver

# Check output
ls -la src/web/ui/
# Should see: index.html, style.css, script.js

# golite should auto-minify to public/
ls -la src/web/public/
```

---

### Phase 8: Cleanup 🧹
**Goal**: Remove old system entirely

**Tasks**:
1. Delete `builder/` directory
   ```bash
   rm -rf builder/
   ```

2. Delete old templates
   ```bash
   rm src/web/ui/template.html
   ```

3. Remove unused imports from main.go

4. Update `go.mod` (remove builder references)

5. Update documentation

6. Run tests
   ```bash
   go test ./...
   ```

---

## 🧪 Testing Strategy

### Unit Tests
```bash
# Test individual components
go test ./src/pkg/gosite/... -v

# Test modules
go test ./src/internal/homePage/... -v
```

### Integration Tests
```bash
# Test full rendering pipeline
go test ./src/cmd/appserver/... -v
```

### Visual Regression
```bash
# Manual comparison
firefox src/web/public/index.html
# Compare with screenshots from old system
```

---

## 📊 Migration Progress Tracker

| Phase | Status | Completion | Blockers |
|-------|--------|------------|----------|
| 1. Foundation | ⏳ Not started | 0% | - |
| 2. Components | ⏳ Not started | 0% | Phase 1 |
| 3. First Module | ⏳ Not started | 0% | Phase 2 |
| 4. Parallel Test | ⏳ Not started | 0% | Phase 3 |
| 5. All Modules | ⏳ Not started | 0% | Phase 4 |
| 6. CSS Migration | ⏳ Not started | 0% | Phase 5 |
| 7. Cutover | ⏳ Not started | 0% | Phase 6 |
| 8. Cleanup | ⏳ Not started | 0% | Phase 7 |

---

## ⚠️ Risk Mitigation

### Risk 1: Output Doesn't Match
**Likelihood**: High  
**Impact**: Medium  

**Mitigation**:
- Use diff tools extensively
- Visual QA on every change
- Keep old system running in parallel during Phase 4

---

### Risk 2: golite Integration Issues
**Likelihood**: Low  
**Impact**: High  

**Mitigation**:
- Test golite watching early (Phase 4)
- Verify file paths are unchanged
- Check file write timing (avoid race conditions)

---

### Risk 3: Performance Regression
**Likelihood**: Low  
**Impact**: Medium  

**Mitigation**:
- Benchmark before/after
- Profile `GenerateSite()` execution
- Optimize string builder usage

---

## 🎯 Success Criteria

- [ ] All visual output matches old system (pixel-perfect)
- [ ] No broken links or images
- [ ] CSS/JS behaves identically
- [ ] Build time < old system
- [ ] golite integration works unchanged
- [ ] All tests pass
- [ ] No `html/template` imports remain
- [ ] Documentation updated
- [ ] Team can add new modules without guidance

---

## 📅 Estimated Timeline

| Phase | Estimated Time | Dependencies |
|-------|----------------|--------------|
| 1. Foundation | 2 hours | None |
| 2. Components | 8 hours | Phase 1 |
| 3. First Module | 2 hours | Phase 2 |
| 4. Parallel Test | 2 hours | Phase 3 |
| 5. All Modules | 6 hours | Phase 4 |
| 6. CSS Migration | 4 hours | Phase 5 |
| 7. Cutover | 1 hour | Phase 6 |
| 8. Cleanup | 1 hour | Phase 7 |
| **Total** | **~26 hours** | Sequential |

*Note: Times are estimates, actual may vary based on discoveries*

---

## 🚀 Quick Start Commands

### Phase 1
```bash
mkdir -p src/pkg/gosite
touch src/pkg/ui.go src/pkg/modules.go
cd src/pkg/gosite
touch gosite.go page.go utils.go
```

### Phase 2
```bash
cd src/pkg/gosite
touch nav.go card.go carousel.go form.go section.go css.go js.go
```

### Phase 3
```bash
mkdir -p src/internal/homePage
touch src/internal/homePage/module.go
touch src/internal/homePage/content.go
touch src/internal/homePage/ui.go
```

### Test
```bash
go test ./src/pkg/gosite/... -v
go build ./src/cmd/appserver
```

---

## 📚 References

- [Project Structure](./project_structure.md) - Architecture overview
- [UI Component System](./ui_component_system.md) - Component details
- [Module System](./module_system.md) - Module registry
- [Separation of Concerns](./separation_of_concerns.md) - Layer boundaries

---

**Ready to begin? Await approval of architecture documents before starting Phase 1.**
