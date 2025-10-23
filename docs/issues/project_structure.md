# Project Structure Refactoring

**Status**: Proposal  
**Created**: 2025-10-23  
**Author**: System Architecture  

---

## 🎯 Executive Summary

This document outlines the architectural refactoring of `monjitaschillan.cl` to separate business logic from UI presentation, enabling scalability for a complex online medical system while maintaining clean, maintainable Go code.

### Core Problem
Currently, UI is generated using Go's `html/template` with external HTML files (`template.html`). This approach:
- ❌ Couples business logic with presentation
- ❌ Lacks component reusability
- ❌ Difficult to test UI rendering independently
- ❌ Not scalable for complex dynamic systems

### Solution
Implement a **Component-Based UI System** written entirely in Go using **String Builders** and **Builder Pattern**, with complete separation between:
- **Business Logic** (`src/internal/*`)
- **UI Components** (`src/pkg/gosite/*`)
- **Module Registry** (`src/pkg/modules.go`)

---

## 📚 Related Documentation

- **[⭐ FINAL DECISIONS](./decisions_final.md)** - All approved decisions consolidated
- [UI Component System Design](./ui_component_system.md) - Detailed component architecture
- [Separation of Concerns](./separation_of_concerns.md) - Layer boundaries and data flow
- [Module System & Reflection](./module_system.md) - Module registry and auto-rendering
- [Migration Strategy](./migration_strategy.md) - Step-by-step refactoring plan

---

## 🏗️ Architecture Decisions

### ✅ Decision 1: Pure Go Components (String Builders)
**Rationale**: 
- Type-safe compilation
- No external template parsing overhead
- Better IDE support and refactoring
- Easier to compose and test

**Trade-offs**:
- ➕ Compile-time safety
- ➕ Better performance (no template parsing)
- ➖ More verbose than HTML templates
- ➖ Visual changes require recompilation

---

### ✅ Decision 2: Builder Pattern for Page Construction
```go
// Each internal handler constructs its UI
page := UI.NewPage("Patient Care")
page.AddNav(navItems)
page.AddSection(content)
html := page.Render()
```

**Rationale**:
- Fluent, readable API
- Accumulates components progressively
- Single render point in `main.go`

---

### ✅ Decision 3: Component-Level Asset Management
Each component can render:
- **HTML** (structure)
- **CSS** (styles)
- **JS** (behavior)

The main UI manager accumulates all assets and renders **once**:
- `src/web/ui/style.css` (single file, non-minified)
- `src/web/ui/script.js` (single file, non-minified)

**Note**: `golite` handles minification and moves to `src/web/public/` dynamically.

---

### ✅ Decision 4: Internal Module Ownership of UI Construction
**Responsibility**:
- `src/internal/patientCare` → Builds its own UI sections
- `src/internal/servicesPage` → Builds service listings
- `src/internal/homePage` → Builds homepage layout

**Constraint**: Can **only** use methods from `src/pkg/gosite/ui.go` (public API).

Optional method signature for each module:
```go
func (m *Module) RenderUI(params ...any) (string, error)
```

---

### ✅ Decision 5: Module Registry System
**Location**: `src/pkg/modules.go`

```go
var Modules = []any{
    &patientCare.Module{},
    &servicesPage.Module{},
    &homePage.Module{},
}
```

**Auto-detection**: Using reflection, the `HtmlUI` manager:
1. Detects struct name (e.g., `patientCare`)
2. Converts to lowercase
3. Stores for rendering association

---

### ✅ Decision 6: Package Location
**Path**: `src/pkg/gosite/gosite.go`

**Rationale**: Will be extracted as external package later, but initialized here during development.

**Structure**:
```
src/pkg/gosite/
├── gosite.go          # Main UI manager (HtmlUI struct, New())
├── page.go           # Page builder
├── nav.go            # Navigation component
├── card.go           # Card component
├── carousel.go       # Carousel component
├── form.go           # Form component
├── section.go        # Section component
├── css.go            # CSS accumulator
└── js.go             # JS accumulator
```

**No subdirectories** - all `.go` files at root level.

---

### ✅ Decision 7: No Backward Compatibility
**Rationale**: Current `builder` package is not scalable.

**Action**: Complete replacement, no migration path needed.

---

## 🔄 GenerateSite API

### New Signature
```go
func GenerateSite(destinations ...string) error
```

**Behavior**:
- **Variadic parameter**: Optional destination paths
- **Default**: `src/web/ui/` if no destinations provided
- **Outputs**:
  - `style.css` (accumulated CSS from all components)
  - `script.js` (accumulated JS from all components)
  - `index.html` (rendered page HTML)

**Called from**: `src/cmd/appserver/main.go` (line ~51)

---

## 📊 Data Flow

```
┌──────────────────────┐
│   main.go            │
│   - Initialize UI    │
│   - Load Modules     │
│   - GenerateSite()   │
└──────────┬───────────┘
           │
           ▼
┌──────────────────────┐
│  pkg/modules.go      │
│  var Modules = []any │
└──────────┬───────────┘
           │
           ▼
┌──────────────────────┐     ┌──────────────────────┐
│ internal/patientCare │────▶│   pkg/ui.go          │
│ RenderUI()           │     │   var UI (singleton) │
└──────────────────────┘     └──────────┬───────────┘
                                        │
           ┌────────────────────────────┘
           ▼
┌──────────────────────┐
│  pkg/gosite/*         │
│  - Component funcs   │
│  - CSS/JS builders   │
│  - Page rendering    │
└──────────┬───────────┘
           │
           ▼
┌──────────────────────┐
│  src/web/ui/         │
│  - index.html        │
│  - style.css         │
│  - script.js         │
└──────────────────────┘
           │
           ▼ (golite watches)
┌──────────────────────┐
│  src/web/public/     │
│  - index.html (min)  │
│  - style.css (min)   │
│  - script.js (min)   │
└──────────────────────┘
```

---

## ✅ Resolved Decisions

### ✅ Decision 8: Singleton State Management
**Decision**: **Stateless** - Each call is independent

```go
// Each call is independent
html1 := UI.NewPage("Home").Render()
html2 := UI.NewPage("Services").Render()
// No shared state between calls
```

**Rationale**: Simpler, more predictable, easier to test.

---

### ✅ Decision 9: CSS Deduplication
**Decision**: **Prevent CSS duplication** - Each CSS block should be unique

**Implementation**: Use a hash-based registry to track rendered CSS:
```go
type Page struct {
    cssBlocks map[string]string // hash -> css
    cssHashes []string           // maintain order
}

func (p *Page) AddCSS(css string) *Page {
    if css == "" {
        return p
    }
    hash := hashString(css)
    if _, exists := p.cssBlocks[hash]; !exists {
        p.cssBlocks[hash] = css
        p.cssHashes = append(p.cssHashes, hash)
    }
    return p
}
```

**Rationale**: Avoid redundant CSS, optimize output size.

---

### ✅ Decision 10: Error Handling in RenderUI
**Decision**: **No error return** - Simplified signature

```go
func (m *Module) RenderUI(params ...any) string
```

**Rationale**: 
- Simpler API
- Components should be defensive
- Errors logged internally, never stop render
- Empty string returned on failure

---

### ✅ Decision 11: GenerateSite Frequency
**Decision**: **Generate only on changes**

**Implementation**:
- `golite` detects Go file changes
- Server restarts automatically
- `GenerateSite()` runs once at startup
- Compare output with existing files
- Write only if content changed

**Rationale**: Avoid unnecessary file system operations.

---

### ❓ Question 2: Component Naming Convention
You want reflection to detect module names (e.g., `patientCare` → `"patientcare"`).

**Suggestion**: Should we support multiple naming strategies?
- `PatientCare` → `"patientcare"` (lowercase)
- `PatientCare` → `"patient-care"` (kebab-case)
- `PatientCare` → `"patient_care"` (snake_case)

**My recommendation**: Start with **lowercase** only, add others if needed.

---

### ⚠️ Question 3: RenderUI Parameters - User & Permissions (NEEDS GUIDANCE)

**Context**: Some sections need user information (name, role, permissions) but must remain decoupled.

**Current Signature**:
```go
func (m *Module) RenderUI(params ...any) string
```

**Challenge**: How to pass user context without coupling to specific structs?

**Option A - Interface-Based** (Recommended):
```go
// pkg/interfaces.go
type UserContext interface {
    GetUsername() string
    GetRole() string
    HasPermission(perm string) bool
    IsAuthenticated() bool
}

// Module uses interface
func (m *Module) RenderUI(ctx UserContext) string {
    if ctx.HasPermission("view_patients") {
        // Render sensitive section
    }
}
```

**Pros**:
- ✅ Decoupled from concrete user struct
- ✅ Type-safe
- ✅ Self-documenting
- ✅ Easy to mock for testing
- ✅ Can evolve without breaking modules

**Cons**:
- ⚠️ Less flexible than `...any`
- ⚠️ Need to define interface upfront

---

**Option B - Context Struct**:
```go
type RenderContext struct {
    User interface {
        GetUsername() string
        GetRole() string
        HasPermission(string) bool
    }
    Language string
    Theme    string
}

func (m *Module) RenderUI(ctx RenderContext) string
```

**Pros**:
- ✅ Groups related data
- ✅ Can add more context later
- ✅ Still decoupled via embedded interface

**Cons**:
- ⚠️ More verbose

---

**Option C - Variadic with Convention** (Current):
```go
func (m *Module) RenderUI(params ...any) string {
    // params[0] = UserContext interface (optional)
    // params[1] = Additional context (optional)
}
```

**Pros**:
- ✅ Maximum flexibility
- ✅ Optional parameters

**Cons**:
- ❌ Not type-safe
- ❌ Requires runtime type assertions
- ❌ Not self-documenting
- ❌ Easy to misuse

---

**Recommendation**: **Option A or B**

**Suggested Approach**:
1. Define minimal `UserContext` interface in `pkg/interfaces.go`
2. Actual user struct implements interface
3. Modules depend only on interface
4. Start with minimal methods, extend as needed

**Example Implementation**:
```go
// pkg/interfaces.go
package pkg

type UserContext interface {
    GetUsername() string
    GetRole() string
    HasPermission(permission string) bool
    IsAuthenticated() bool
}

// Default implementation for anonymous users
type AnonymousUser struct{}

func (a *AnonymousUser) GetUsername() string        { return "anonymous" }
func (a *AnonymousUser) GetRole() string            { return "guest" }
func (a *AnonymousUser) HasPermission(string) bool  { return false }
func (a *AnonymousUser) IsAuthenticated() bool      { return false }

// internal/users/user.go (example implementation)
type User struct {
    Username    string
    Role        string
    Permissions []string
}

func (u *User) GetUsername() string { return u.Username }
func (u *User) GetRole() string     { return u.Role }
func (u *User) HasPermission(perm string) bool {
    for _, p := range u.Permissions {
        if p == perm {
            return true
        }
    }
    return false
}
func (u *User) IsAuthenticated() bool { return u.Username != "" }

// Usage in module
func (m *PatientCare) RenderUI(ctx pkg.UserContext) string {
    page := pkg.UI.NewPage("Patient Care")
    
    if !ctx.HasPermission("view_patients") {
        return pkg.UI.AccessDenied()
    }
    
    // Build UI with user context
    header := pkg.UI.Header("Welcome, " + ctx.GetUsername())
    page.AddSection(header)
    
    // ...
}
```

**Question**: Which option do you prefer? Need more details on your permission system?

---

### ❓ Question 4: CSS/JS Deduplication
If multiple components use the same CSS classes or JS functions:

**Question**: Should the system auto-deduplicate?

**Example**:
```go
// card.go renders: .card { border: 1px solid #ccc; }
// product-card.go also renders: .card { border: 1px solid #ccc; }
```

**Options**:
- **A**: Keep duplicates (simpler, CSS gzip handles it well)
- **B**: Hash-based deduplication (complex, optimal output)

**My recommendation**: Start with **A**, optimize later if needed.

---

### ❓ Question 5: Error Handling Strategy
When a component fails to render:

**Question**: Should it:
- **A**: Return error, stop entire page render?
- **B**: Log error, render placeholder, continue?
- **C**: Panic (fail fast)?

**My recommendation**: **A** for development, **B** for production (with logging).

---

### ⚠️ Potential Issue: Hot Reload with golite
You mentioned `golite` watches `src/web/ui/` and minifies to `public/`.

**Concern**: If `GenerateSite()` runs on every request during development:
- Constant file writes to `src/web/ui/`
- golite triggers on every change
- Could cause file system thrashing

**Solutions**:
1. **GenerateSite only on startup** (static mode)
2. **Add debouncing** in golite
3. **Dev mode**: Direct render to `public/` (skip ui/ intermediate)

**Question**: How often should `GenerateSite()` run?
- On server startup only?
- On every HTTP request (for development)?
- Triggered by file changes?

---

### 💡 Suggestion 1: Component Interface
Consider defining a standard component interface:

```go
type Component interface {
    RenderHTML() string
    RenderCSS() string
    RenderJS() string
}
```

**Pros**: Uniform component structure, easier testing  
**Cons**: More boilerplate

**Alternative**: Keep free functions (more flexible).

---

### 💡 Suggestion 2: Template Literals Helper
Since using string builders, consider a helper for multiline strings:

```go
// Instead of:
html := "<div class=\"card\">\n"
html += "  <h2>" + title + "</h2>\n"
html += "</div>\n"

// Use:
html := Multiline(`
    <div class="card">
        <h2>%s</h2>
    </div>
`, title)
```

---

### 💡 Suggestion 3: Development vs Production Mode
Add a mode flag:

```go
UI.SetMode(gosite.Development) // Pretty-printed, comments
UI.SetMode(gosite.Production)  // Compact, no comments
```

**Use case**: Easier debugging during development.

---

## ✅ Next Steps

1. **Review this document** and answer open questions
2. **Approve architecture** or request changes
3. **Read detailed docs** (ui_component_system.md, etc.)
4. **Begin implementation** following migration_strategy.md

---

## 🏁 Success Criteria

- ✅ Zero `html/template` usage
- ✅ All UI code in pure Go
- ✅ `internal/*` has no HTML knowledge
- ✅ Single `UI` entry point from `pkg/ui.go`
- ✅ Component reusability across modules
- ✅ CSS/JS bundled in single files
- ✅ Compatible with `golite` pipeline

---

**Awaiting approval before proceeding to implementation.**
