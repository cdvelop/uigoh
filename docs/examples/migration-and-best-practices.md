## 13. Migration & Best Practices

### 13.1. From Static HTML

**Before:**
```html
<!-- index.html -->
<html>
<body>
    <section>
        <h1>Welcome</h1>
        <div class="card">
            <h3>Card Title</h3>
            <p>Description</p>
        </div>
    </section>
</body>
</html>
```

**After:**
```go
site := gosite.New(cfg)
page := site.NewPage("Welcome", "index.html")
section := page.NewSection("Welcome")
section.Add(&card.Card{
    Title:       "Card Title",
    Description: "Description",
})
site.Generate()
```

### 13.2. Best Practices

**1. Single Source of Truth - Site Structure:**
```go
// ✅ GOOD: One function defines the entire site
// app/site.go (NO build tags)
func BuildSite(site *gosite.Site) {
    homePage := site.NewPage("Home", "index.html")
    aboutPage := site.NewPage("About", "about.html")
    // ... build all pages
}

// ❌ BAD: Duplicating site structure
// Don't write the same page creation twice in different files
```

**2. Environment-Specific Entry Points Only:**
```go
// ✅ GOOD: Only entry points have build tags
// cmd/generator/main.go
//go:build !wasm
func main() {
    site := gosite.New(backendConfig)
    app.BuildSite(site)  // Call shared logic
    site.Generate()
}

// cmd/webapp/main.go
//go:build wasm
func main() {
    site := gosite.New(frontendConfig)
    app.BuildSite(site)  // Same shared logic
    mountPage(site.Pages[0])
}
```

**3. Component Organization:**
```
components/
├── card/
│   └── card.go              # NO build tags - works everywhere
├── form/
│   ├── form.go              # NO build tags - base component
│   ├── env.backend.go      # //go:build !wasm - backend specific
│   └── env.frontend.go     # //go:build wasm - frontend specific
└── navbar/
    └── navbar.go            # NO build tags
```

**4. Configuration Pattern:**
```go
// ✅ GOOD: Separate configs, shared logic
type AppConfig struct {
    Title string
    Theme *gosite.ColorScheme
}

func BuildSiteWithConfig(site *gosite.Site, appCfg AppConfig) {
    // Use appCfg to customize site building
    // This function has NO build tags
}

// Backend
//go:build !wasm
cfg := &gosite.Config{
    Title: appCfg.Title,
    OutputDir: "dist",
    WriteFile: os.WriteFile,
}

// Frontend  
//go:build wasm
cfg := &gosite.Config{
    Title: appCfg.Title,
    EventBinder: eventBinder,
}
```

**5. Avoid Site References in Components:**
```go
// ✅ GOOD: Self-contained component
type Card struct {
    Title       string
    Description string
}

func (c *Card) RenderHTML() string { ... }
func (c *Card) RenderCSS() string { ... }

// ❌ BAD: Component depends on Site
type Card struct {
    site *gosite.Site  // Don't do this
    Title string
}
```

**6. Progressive Enhancement Pattern:**
```go
// Shared component that works everywhere
type Form struct {
    ID     string
    Fields []Field
}

// Always renders functional HTML
func (f *Form) RenderHTML() string {
    // Renders a standard HTML form
    // Works without JavaScript (backend)
    // Enhanced with JavaScript (frontend)
}

// Frontend-specific enhancement
//go:build wasm
func (f *Form) AttachValidation(site *gosite.Site) {
    // Add client-side validation in WASM
}
```

**7. Build Tags File Organization:**
```go
// ✅ GOOD: Clear separation
main.go              // Shared types and functions
main_backend.go      // //go:build !wasm - backend main()
main_frontend.go     // //go:build wasm - frontend main()

// ❌ BAD: Mixing concerns
main.go              // Everything in one file with complex build tag logic
```

---