# Gosite API Structure & Design

## 1. Abstract

The goal of `gosite` is to provide a programmatic Go framework for building a wide range of websites, from backend-rendered Server-Side Rendering (SSR) or Multi-Page Applications (MPA) to frontend-rendered WebAssembly (WASM) Single-Page Applications (SPA).

### 1.1. Core Philosophy: Write Once, Run Anywhere

**The fundamental principle:** You write your site structure **once** in shared code (no build tags). The framework adapts its behavior based on the compilation target:

- **Backend (`!wasm`)**: Generates static HTML/CSS/JS files to disk
- **Frontend (`wasm`)**: Renders pages dynamically in the browser DOM

```go
// ✅ CORRECT: Single source of truth
// app/site.go (NO build tags)
func BuildSite(site *gosite.Site) {
    page := site.NewPage("Home", "index.html")
    section := page.NewSection("Welcome")
    section.Add(components...)
}

// Backend: calls BuildSite(site) then site.Generate()
// Frontend: calls BuildSite(site) then mountPage(page)
```

This approach eliminates code duplication and ensures your site structure is consistent across environments.

### 1.3. Architecture Overview

```
┌─────────────────────────────────────────────────────────┐
│                   Your Application                      │
│  ┌───────────────────────────────────────────────────┐ │
│  │        app/site.go (NO BUILD TAGS)                │ │
│  │                                                    │ │
│  │  func BuildSite(site *gosite.Site) {              │ │
│  │      page := site.NewPage(...)                    │ │
│  │      section := page.NewSection(...)              │ │
│  │      section.Add(components...)                   │ │
│  │  }                                                 │ │
│  │                                                    │ │
│  │  ← Written ONCE, used by BOTH environments        │ │
│  └────────────────┬───────────────┬──────────────────┘ │
└──────────────────│───────────────│─────────────────────┘
                   │               │
       ┌───────────▼──────┐   ┌───▼──────────────┐
       │   Backend Build  │   │  Frontend Build  │
       │  //go:build !wasm│   │ //go:build wasm  │
       └───────────┬──────┘   └───┬──────────────┘
                   │               │
       ┌───────────▼──────┐   ┌───▼──────────────┐
       │  gosite.New()    │   │  gosite.New()    │
       │  (Backend Mode)  │   │ (Frontend Mode)  │
       │                  │   │                  │
       │  • Has WriteFile │   │  • Has Event     │
       │  • Has Generate()│   │    Binder        │
       │  • Writes to disk│   │  • Renders to DOM│
       └───────────┬──────┘   └───┬──────────────┘
                   │               │
       ┌───────────▼──────┐   ┌───▼──────────────┐
       │     Output       │   │     Browser      │
       │                  │   │                  │
       │  dist/           │   │  [Live SPA]      │
       │  ├── index.html  │   │  • Dynamic       │
       │  ├── style.css   │   │  • Interactive   │
       │  └── script.js   │   │  • Client-side   │
       └──────────────────┘   └──────────────────┘
```

### 1.2. Core Principles

- **Unified Go Codebase**: Write frontend and backend logic in Go.
- **Single Site Definition**: Define your site structure once, use it everywhere.
- **Performance**: Leverage build tags (`wasm` and `!wasm`) to separate environment-specific logic only.
- **TinyGo Compatibility**: Strictly use `tinystring` for all string manipulations to minimize binary size.
- **Fluent API**: Offer a clean, chained API for building pages and components.
- **Decoupled DOM Interaction**: Abstract all browser/DOM APIs through interfaces.

## 2. Core Concepts

The structure of a website is hierarchical: a `Site` contains `Pages`, each `Page` contains `Sections`, and each `Section` contains `Components`.

### 2.1. Site
The root object for a website. It holds global configuration, manages pages, and orchestrates the final build process. It is the main entry point.

**Responsibilities:**
- Global configuration management
- Page registration and lifecycle
- CSS/JS asset aggregation and deduplication
- File generation (backend only)
- Navigation building

### 2.2. Page
Represents a single HTML page (e.g., `index.html`, `about.html`). It manages its own content, composed of sections.

**Responsibilities:**
- Section management
- HTML structure generation
- Head content management (meta tags, title, etc.)

### 2.3. Section
A logical division within a page (e.g., a hero, a gallery, a contact form). It acts as a container for components and helps structure the page layout.

**Responsibilities:**
- Component registration
- Section HTML rendering
- Component CSS/JS registration with Site

### 2.4. Component
A reusable UI element (e.g., a card, a button, a form). Components are the smallest building blocks and are added to sections.

**Characteristics:**
- Must implement `HTMLRenderer` interface at minimum
- Can optionally implement `CSSRenderer` and `JSRenderer`
- Self-contained and reusable
- Can work in both backend and frontend environments (with proper build tags)

## 3. API Design

To enforce a clear and logical workflow, `gosite` will use public structs with private fields, where instantiation is controlled via a fluent, chained API.

### 3.1. Entry Point & Chaining

The user starts by creating a `Site` instance via the `New()` function. From there, all other objects are created through chained method calls. Direct instantiation of `Page` or `Section` is not intended.

```go
// The only public constructor.
// Returns a public struct `Site` with private fields.
site := gosite.New(&gosite.Config{
    Title:     "My Website",
    OutputDir: "output",
    WriteFile: os.WriteFile, // For backend only
})

// `NewPage` is a method on `Site` and returns a `*Page`.
page1 := site.NewPage("Home", "index.html")

// `NewSection` is a method on `Page` and returns a `*Section`.
// The title parameter helps structure the page.
section1 := page1.NewSection("Welcome")

// `Add` is a method on `Section` that accepts any component.
// It returns `*Section` to allow further chaining.
section1.Add(cardComponent1).
    Add(cardComponent2).
    Add(customComponent)

// Alternative: create another section
page1.NewSection("Features").
    Add(featureCard1).
    Add(featureCard2)

// Finally, generate the site (backend only)
err := site.Generate()
```

### 3.2. Struct Visibility

#### Public Structs (Exported)
- **`gosite.Site`**: Public struct, but fields like `pages`, `cssBlocks`, `jsBlocks`, `buff` are private.
- **`gosite.Page`**: Public struct, but fields like `site`, `sections`, `title`, `filename`, `head` are private.
- **`gosite.Section`**: Public struct, with private fields `page`, `site`, `components`.

#### Public Fields (When Necessary)
- **`Section.Title`**: Public string - allows customization of section title.
- **`Section.ModuleID`**: Public string - allows custom ID for the section element.

#### Component Interfaces (Exported)
All components must implement at least one of these interfaces:
- **`HTMLRenderer`**: Required - renders component HTML.
- **`CSSRenderer`**: Optional - provides component-specific CSS.
- **`JSRenderer`**: Optional - provides component-specific JavaScript.

### 3.3. Method Signatures

```go
// Site methods
func New(cfg *Config) *Site
func (s *Site) NewPage(title, filename string) *Page
func (s *Site) Generate() error              // Backend only
func (s *Site) PageCount() int
func (s *Site) BuildNav() string
func (s *Site) AddCSS(css string)            // Backend only (no-op in frontend)
func (s *Site) AddJS(js string)              // Backend only (no-op in frontend)

// Page methods
func (p *Page) NewSection(title string) *Section
func (p *Page) AddHead(content string) *Page
func (p *Page) RenderHTML() string

// Section methods
func (s *Section) Add(component any) *Section
func (s *Section) Render() string
```

## 4. Build Architecture (Backend vs. Frontend)

The codebase will be split to isolate backend-only logic from the frontend WASM binary. This is achieved using Go's build tags.

### 4.1. File Naming Convention

We will adopt a clear file naming convention to separate concerns:
- **`env.backend.go`**: Contains code that should only be compiled for the backend. It will be guarded by `//go:build !wasm`. This includes file writing (`Generate`), CSS/JS aggregation, and other server-side logic.
- **`env.frontend.go`**: Contains code specific to the WASM environment. It will be guarded by `//go:build wasm`. This includes the logic that interacts with the `EventBinder`.
- **Files without `env` prefix**: These files contain shared logic, such as struct definitions (`Site`, `Page`), interfaces, and component logic that is common to both environments.

### 4.2. What Gets Build Tags vs What Doesn't

**Rule of Thumb:** Only entry points and environment-specific APIs need build tags.

#### ❌ NO Build Tags Needed (Shared Code):

```go
// Site structure definition
func BuildSite(site *gosite.Site) { ... }

// Component definitions
type Card struct { ... }
func (c *Card) RenderHTML() string { ... }

// Business logic
func ProcessData(input string) string { ... }

// Data structures
type User struct { ... }

// Helper functions using tinystring
func FormatTitle(title string) string { ... }
```

#### ✅ Build Tags Required:

```go
// Entry point functions
//go:build !wasm
func main() { ... }

//go:build wasm  
func main() { ... }

// Environment-specific configuration
//go:build !wasm
func getConfig() *gosite.Config {
    return &gosite.Config{
        WriteFile: os.WriteFile,  // Backend only
    }
}

//go:build wasm
func getConfig() *gosite.Config {
    return &gosite.Config{
        EventBinder: eventBinder,  // Frontend only
    }
}

// DOM manipulation
//go:build wasm
func mountToDOM(html string) { ... }

// File system operations
//go:build !wasm
func saveToFile(path, content string) error { ... }
```

### 4.3. Example: Site Struct

The `Site` struct will have different fields depending on the build target.

```go
// env.backend.go
//go:build !wasm

package gosite

// Site struct for the backend, includes build-related fields.
type Site struct {
    Cfg       *Config
    pages     []*Page
    cssBlocks []assetBlock
    jsBlocks  []assetBlock
    buff      *Conv
}
// ... methods for file generation ...
```

```go
// env.frontend.go
//go:build wasm

package gosite

// Site struct for the frontend, lightweight.
type Site struct {
    Cfg   *Config
    pages []*Page
}
// ... methods for WASM rendering ...
```

## 5. Event Handling (WASM)

To keep `gosite` independent of the browser's DOM API, all event interactions will be handled through an `EventBinder` interface provided in the `Config`.

### 5.1. The `EventBinder` Interface

The user must provide an implementation for this interface when creating a WASM application.

```go
// in interfaces.go
package gosite

type EventBinder interface {
    // EventListener adds or removes an event listener from a DOM element.
    // The implementation will use syscall/js to interact with the DOM.
    // 
    // Parameters:
    //   add: true to add listener, false to remove it
    //   elementID: the HTML element's id attribute
    //   eventType: event name (e.g., "click", "change", "submit")
    //   callback: Go function to execute when event fires
    EventListener(add bool, elementID, eventType string, callback func())
}
```

### [5.2. Implementation Example](examples/5.2-implementation.md)



### [5.3. Usage in Components](examples/5.3-usage-in-components.md)


## 6. Dependencies

- **`tinystring`**: This is a hard requirement. All internal string operations, conversions, and concatenations must use the `tinystring` package to ensure compatibility with TinyGo and to produce the smallest possible WASM binaries. The use of `fmt`, `strings`, `strconv`, `errors`, `path`, or `bytes` is disallowed in the core framework.

### [6.1. TinyString API Coverage](examples/6.1-tinystring-api-coverage.md)


### 6.2. Error Handling Without `errors` Package

Since we cannot use the standard `errors` package, error handling must be done through:

1. **Return boolean flags** for simple success/failure
2. **Return error strings** when more context is needed
3. **Use panic** for truly exceptional, unrecoverable errors (sparingly)

```go
// Example: WriteFile function signature in Config
type Config struct {
    // Returns error as string, not error interface
    WriteFile func(path string, content string) error
}

// Implementation can use standard library since it's user-provided
cfg := &Config{
    WriteFile: func(path, content string) error {
        return os.WriteFile(path, []byte(content), 0644)
    },
}
```

---

## 7. Complete Examples

### 7.1. Shared Application Logic (No Build Tags)

The key principle: **Write your site structure once, gosite adapts to the environment.**

```go
// app/site.go
// No build tags - this code works in both backend and frontend
package app

import (
    "github.com/cdvelop/gosite"
    "github.com/cdvelop/gosite/components/card"
)

// BuildSite creates the complete site structure
// This is the ONLY place where you define your site
func BuildSite(site *gosite.Site) {
    // Create home page
    homePage := site.NewPage("Home", "index.html")
    
    // Add sections and components
    heroSection := homePage.NewSection("Welcome to My Portfolio")
    heroSection.Add(&card.Card{
        Title:       "About Me",
        Description: "Full-stack developer with 5 years of experience",
        Image:       "/assets/avatar.jpg",
    })

    // Create about page
    aboutPage := site.NewPage("About", "about.html")
    aboutSection := aboutPage.NewSection("My Skills")
    aboutSection.Add(&card.Card{
        Title:       "Go",
        Description: "Backend development",
    }).Add(&card.Card{
        Title:       "JavaScript",
        Description: "Frontend development",
    })
    
    // Create contact page
    contactPage := site.NewPage("Contact", "contact.html")
    contactSection := contactPage.NewSection("Get in Touch")
    contactSection.Add(&card.Card{
        Title:       "Email",
        Description: "contact@example.com",
    })
}
```

### 7.2. Backend Entry Point (SSR/MPA Generator)

```go
//go:build !wasm

// cmd/generator/main.go
package main

import (
    "os"
    "github.com/cdvelop/gosite"
    "myproject/app"
)

func main() {
    // Backend-specific configuration
    cfg := &gosite.Config{
        Title:     "My Portfolio",
        OutputDir: "dist",
        WriteFile: func(path, content string) error {
            return os.WriteFile(path, []byte(content), 0644)
        },
        ColorScheme: gosite.DefaultColorScheme(),
    }

    // Create site
    site := gosite.New(cfg)

    // Use the shared build logic
    app.BuildSite(site)

    // Generate files (only works in backend)
    if err := site.Generate(); err != nil {
        panic(err)
    }
    
    println("Site generated successfully in", cfg.OutputDir)
}
```

**Output structure:**
```
dist/
├── index.html          # Home page
├── about.html          # About page  
├── contact.html        # Contact page
├── style.css           # Combined CSS (base + components)
└── script.js           # Combined JS (if any components use it)
```

### 7.3. Frontend Entry Point (WASM/SPA)

```go
//go:build wasm

// cmd/webapp/main.go
package main

import (
    "syscall/js"
    "github.com/cdvelop/gosite"
    "myproject/app"
)

func main() {
    // Frontend-specific configuration
    cfg := &gosite.Config{
        Title:       "My Portfolio",
        EventBinder: NewDOMEventBinder(), // Frontend needs event handling
        ColorScheme: gosite.DefaultColorScheme(),
    }

    // Create site
    site := gosite.New(cfg)

    // Use the SAME shared build logic
    app.BuildSite(site)

    // Frontend-specific: mount to DOM
    mountInitialPage(site)
    
    // Setup client-side routing
    setupRouting(site)

    // Keep WASM running
    <-make(chan struct{})
}

func mountInitialPage(site *gosite.Site) {
    // Get the initial page based on URL or default to first page
    // In frontend, Generate() is not called - pages are rendered on demand
    page := site.Pages[0] // Assuming we expose Pages or have a GetPage method
    mountPage(page)
}

func mountPage(page *gosite.Page) {
    html := page.RenderHTML()
    root := js.Global().Get("document").Call("getElementById", "app")
    root.Set("innerHTML", html)
}

func setupRouting(site *gosite.Site) {
    // Setup browser history navigation
    js.Global().Call("addEventListener", "popstate", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
        // Handle back/forward navigation
        path := js.Global().Get("location").Get("pathname").String()
        // Find and mount corresponding page
        return nil
    }))
}
```

### 7.4. Project Structure

```
myproject/
├── app/
│   └── site.go              # Shared site structure (NO build tags)
├── cmd/
│   ├── generator/
│   │   └── main.go          # Backend entry (//go:build !wasm)
│   └── webapp/
│       └── main.go          # Frontend entry (//go:build wasm)
├── components/
│   ├── card/
│   │   └── card.go          # Shared component (NO build tags)
│   └── form/
│       └── form.go          # Shared component (NO build tags)
└── go.mod
```

### 7.5. Advanced: Environment-Specific Components

Sometimes you need components that work differently per environment:

```go
// components/interactive/button.go
// NO build tags - shared interface
package interactive

type Button struct {
    ID      string
    Text    string
    OnClick func()
}

func (b *Button) RenderHTML() string {
    return `<button id="` + b.ID + `">` + b.Text + `</button>`
}
```

```go
//go:build !wasm

// components/interactive/button_backend.go
package interactive

// In backend, events are just markers in HTML
func (b *Button) AttachEvents(site *gosite.Site) {
    // No-op in backend, or add data attributes for progressive enhancement
}
```

```go
//go:build wasm

// components/interactive/button_frontend.go
package interactive

import "github.com/cdvelop/gosite"

// In frontend, actually attach DOM events
func (b *Button) AttachEvents(site *gosite.Site) {
    if site.Cfg.EventBinder != nil && b.OnClick != nil {
        site.Cfg.EventBinder.EventListener(true, b.ID, "click", b.OnClick)
    }
}
```

**Usage in shared code:**
```go
// app/site.go (still no build tags!)
func BuildSite(site *gosite.Site) {
    btn := &interactive.Button{
        ID:   "submit-btn",
        Text: "Submit",
        OnClick: func() {
            // This callback works in WASM, ignored in backend
            println("Button clicked!")
        },
    }
    
    section.Add(btn)
    btn.AttachEvents(site) // Different behavior per environment
}
```

### 7.6. Complete Minimal Example (Single File)

For simple cases, everything can be in one file:

```go
// main.go
package main

import (
    "github.com/cdvelop/gosite"
)

func buildSite(site *gosite.Site) {
    page := site.NewPage("Home", "index.html")
    section := page.NewSection("Welcome")
    section.Add(&SimpleCard{Title: "Hello", Text: "World"})
}

type SimpleCard struct {
    Title string
    Text  string
}

func (c *SimpleCard) RenderHTML() string {
    return `<div><h3>` + c.Title + `</h3><p>` + c.Text + `</p></div>`
}

// Backend entry point
//go:build !wasm
func main() {
    site := gosite.New(&gosite.Config{
        OutputDir: "dist",
        WriteFile: os.WriteFile,
    })
    buildSite(site)
    site.Generate()
}

// Frontend entry point
//go:build wasm
func main() {
    site := gosite.New(&gosite.Config{
        EventBinder: NewDOMEventBinder(),
    })
    buildSite(site)
    mountPage(site.Pages[0])
    <-make(chan struct{})
}
```

**Note:** The above won't compile as-is (can't have two `main` functions), but illustrates the concept. In practice, use separate files with build tags for the entry points.

```go
// components/card/card.go
// No build tags - works in both backend and frontend

package card

import (
    . "github.com/cdvelop/tinystring"
)

type Card struct {
    Title       string
    Description string
    Image       string
}

func (c *Card) RenderHTML() string {
    b := Convert()
    b.Write(`<div class="card">`)
    
    if c.Image != "" {
        b.Write(`<img src="`)
        b.Write(Convert(c.Image).EscapeAttr())
        b.Write(`" alt="`)
        b.Write(Convert(c.Title).EscapeAttr())
        b.Write(`">`)
    }
    
    b.Write(`<h3>`)
    b.Write(Convert(c.Title).EscapeHTML())
    b.Write(`</h3>`)
    
    b.Write(`<p>`)
    b.Write(Convert(c.Description).EscapeHTML())
    b.Write(`</p>`)
    
    b.Write(`</div>`)
    return b.String()
}

func (c *Card) RenderCSS() string {
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

.card img {
    width: 100%;
    border-radius: 4px;
    margin-bottom: 1rem;
}

.card h3 {
    color: var(--color-primary);
    margin-bottom: 0.5rem;
}`
}

// No RenderJS needed for this component
```

---

## 8. Component Development

### 8.1. Component Interfaces

All components must implement at least the `HTMLRenderer` interface:

```go
type HTMLRenderer interface {
    RenderHTML() string
}

type CSSRenderer interface {
    RenderCSS() string
}

type JSRenderer interface {
    RenderJS() string
}
```

### 8.2. Creating a Custom Component

```go
// components/gallery/gallery.go
package mycomponent

import (
    . "github.com/cdvelop/tinystring"
)

type Gallery struct {
    Images []string
    Columns int
}

// shared - no build tags
func (g *Gallery) RenderHTML() string {
    b := Convert()
    b.Write(`<div class="gallery">`)
    
    for _, img := range g.Images {
        b.Write(`<div class="gallery-item"><img src="`)
        b.Write(Convert(img).EscapeAttr())
        b.Write(`"></div>`)
    }
    
    b.Write(`</div>`)
    return b.String()
}

// components/gallery/env.backend.go
// !wasm only - CSS generation for backend
func (g *Gallery) RenderCSS() string {
    cols := "3"
    if g.Columns > 0 {
        cols = Convert().WriteInt(g.Columns).String()
    }
    
    return `.gallery {
    display: grid;
    grid-template-columns: repeat(` + cols + `, 1fr);
    gap: 1rem;
}

.gallery-item img {
    width: 100%;
    height: auto;
    border-radius: 4px;
}`
}

// Usage:
gallery := &Gallery{
    Images: []string{"/img1.jpg", "/img2.jpg", "/img3.jpg"},
    Columns: 3,
}
section.Add(gallery)
```

### 8.3. Interactive Component (Frontend Only)

```go
//go:build wasm

package counter

import (
    . "github.com/cdvelop/tinystring"
    "github.com/cdvelop/gosite"
)

type Counter struct {
    ID    string
    count int
    site  *gosite.Site
}

func New(id string, site *gosite.Site) *Counter {
    c := &Counter{
        ID:   id,
        site: site,
    }
    c.attachEvents()
    return c
}

func (c *Counter) RenderHTML() string {
    b := Convert()
    b.Write(`<div class="counter">`)
    b.Write(`<button id="`)
    b.Write(c.ID)
    b.Write(`-dec">-</button>`)
    b.Write(`<span id="`)
    b.Write(c.ID)
    b.Write(`-value">0</span>`)
    b.Write(`<button id="`)
    b.Write(c.ID)
    b.Write(`-inc">+</button>`)
    b.Write(`</div>`)
    return b.String()
}

func (c *Counter) attachEvents() {
    if c.site.Cfg.EventBinder == nil {
        return
    }
    
    binder := c.site.Cfg.EventBinder
    
    // Increment button
    binder.EventListener(true, c.ID+"-inc", "click", func() {
        c.count++
        c.updateDisplay()
    })
    
    // Decrement button
    binder.EventListener(true, c.ID+"-dec", "click", func() {
        c.count--
        c.updateDisplay()
    })
}

func (c *Counter) updateDisplay() {
    // Update DOM directly (only in WASM)
    js.Global().Get("document").
        Call("getElementById", c.ID+"-value").
        Set("textContent", Convert().WriteInt(c.count).String())
}
```

---

## 9. Build Process

### 9.1. Backend Build

```bash
# Standard Go build
go build -o bin/generator ./cmd/generator

# Run the generator
./bin/generator

# Output will be in the configured OutputDir (e.g., "dist/")
```

### 9.2. Frontend WASM Build

```bash
# Using TinyGo for smaller binaries
tinygo build -o public/app.wasm -target wasm ./cmd/webapp

# Copy the WASM exec helper
cp $(tinygo env TINYGOROOT)/targets/wasm_exec.js public/

# Serve the application
# public/index.html should load app.wasm
```

**Example index.html for WASM:**
```html
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>My WASM App</title>
    <link rel="stylesheet" href="style.css">
</head>
<body>
    <div id="app"></div>
    <script src="wasm_exec.js"></script>
    <script>
        const go = new Go();
        WebAssembly.instantiateStreaming(
            fetch("app.wasm"),
            go.importObject
        ).then(result => {
            go.run(result.instance);
        });
    </script>
</body>
</html>
```

### 9.3. Hybrid Approach

You can combine both:
1. Use backend to generate initial HTML structure
2. Enhance with WASM for interactivity

```bash
# 1. Generate static structure
go run ./cmd/generator

# 2. Build WASM components
tinygo build -o dist/components.wasm -target wasm ./cmd/components

# 3. Include WASM in generated HTML pages
```

---

## 10. Lifecycle & State Management

### 10.1. Backend Lifecycle

```
1. New(cfg)           → Create Site
2. site.NewPage()     → Create Pages
3. page.NewSection()  → Create Sections
4. section.Add()      → Add Components
5. site.Generate()    → Write files to disk
```

**Important Notes:**
- CSS/JS are automatically collected when components are added
- Deduplication happens automatically
- Files are only written when `Generate()` is called
- After `Generate()`, modifications won't affect output

### 10.2. Frontend Lifecycle

```
1. New(cfg)              → Create Site (with EventBinder)
2. site.NewPage()        → Create Pages
3. page.NewSection()     → Create Sections
4. section.Add()         → Add Components
5. component.AttachEvents() → Register DOM events
6. mountPage()           → Render to DOM
```

**Important Notes:**
- No file generation occurs
- Pages are rendered to strings and mounted to DOM
- Navigation is handled client-side
- Events must be attached after mounting

### 10.3. State Management Patterns

**Backend (Stateless):**
```go
// Each generation is independent
func generateSite(data SomeData) {
    site := gosite.New(cfg)
    // Build pages from data
    site.Generate()
}
```

**Frontend (Stateful):**
```go
type AppState struct {
    currentPage *gosite.Page
    userData    UserData
}

func (a *AppState) Navigate(page *gosite.Page) {
    a.currentPage = page
    mountPage(page)
}
```

---

## 11. Limitations & Constraints

### 11.1. TinyGo/WASM Constraints

**Package Restrictions:**
- ❌ Cannot use `fmt`, `strings`, `strconv`, `errors`, `path`, `bytes`
- ❌ Cannot use reflection extensively
- ❌ Some standard library packages are not available
- ✅ Must use `tinystring` for all string operations
- ✅ `syscall/js` is available for DOM manipulation

**Binary Size Considerations:**
- Each imported package increases WASM size
- Use build tags to exclude unnecessary code
- Prefer simple data structures over complex ones

### 11.2. Architecture Limitations

**The Golden Rule: Write Once, Build Tags Only for Entry Points**

```go
// ✅ CORRECT PATTERN
// app/site.go (NO build tags)
func BuildSite(site *gosite.Site) {
    // All site structure here
    page := site.NewPage(...)
    // This works in BOTH backend and frontend
}

// main_backend.go
//go:build !wasm
func main() {
    site := gosite.New(backendConfig)
    BuildSite(site)  // Reuse shared logic
    site.Generate()  // Backend-specific
}

// main_frontend.go
//go:build wasm
func main() {
    site := gosite.New(frontendConfig)
    BuildSite(site)  // Same shared logic
    mountPage(...)   // Frontend-specific
}

// ❌ WRONG PATTERN - Don't do this!
// Duplicating site structure in both environments
```

**Component Sharing:**
```go
// ✅ Shared component (works everywhere, NO build tags)
type Card struct {
    Title string
}
func (c *Card) RenderHTML() string { ... }

// ✅ Component with environment-specific enhancements
// card.go (NO build tags)
type InteractiveCard struct {
    Title string
    OnClick func()
}
func (c *InteractiveCard) RenderHTML() string { ... }

// card_backend.go
//go:build !wasm
func (c *InteractiveCard) Enhance(site *gosite.Site) {
    // No-op or progressive enhancement
}

// card_frontend.go
//go:build wasm
func (c *InteractiveCard) Enhance(site *gosite.Site) {
    // Attach actual DOM events
    site.Cfg.EventBinder.EventListener(...)
}
```

**What Can Be Shared:**
- ✅ Site structure (`NewPage`, `NewSection`, `Add`)
- ✅ Component definitions and HTML rendering
- ✅ Data structures
- ✅ Business logic (using tinystring only)
- ✅ Validation logic
- ✅ Component CSS/JS definitions

**What Cannot Be Shared (Needs Build Tags):**
- ❌ `main()` functions
- ❌ File system operations (`os.WriteFile`)
- ❌ DOM manipulation (`syscall/js`)
- ❌ `site.Generate()` calls
- ❌ Browser-specific APIs
- ❌ Server-specific APIs

**CSS/JS Aggregation:**
- Only works in backend mode
- Frontend mode ignores `AddCSS()` and `AddJS()` calls
- WASM apps must handle styling differently (inline or external)

### 11.3. Known Issues & Workarounds

**Issue: Error Handling**
```go
// ❌ Cannot use errors.New()
// ✅ Workaround: Return error strings or use boolean flags
func validate(data string) (bool, string) {
    if data == "" {
        return false, "data cannot be empty"
    }
    return true, ""
}
```

**Issue: String Formatting**
```go
// ❌ Cannot use fmt.Sprintf()
// ✅ Workaround: Use tinystring.Fmt()
msg := Fmt("User %s has %d points", username, points)
```

**Issue: Path Operations**
```go
// ❌ Cannot use path.Join()
// ✅ Workaround: Use tinystring.PathJoin()
fullPath := PathJoin("output", "pages", "index.html").String()
```

### 11.4. Browser Compatibility

**WASM Requirements:**
- Modern browsers with WebAssembly support (2017+)
- Chrome 57+, Firefox 52+, Safari 11+, Edge 16+
- No IE11 support

**Recommended Polyfills:** None required for core functionality

---

## 12. API Reference

### 12.1. Configuration

```go
type Config struct {
    // Site title (used in <title> tags)
    Title string

    // Output directory for generated files (backend only)
    // Example: "dist", "public", "output"
    OutputDir string

    // Color scheme for the site
    // If nil, DefaultColorScheme() is used
    ColorScheme *ColorScheme

    // Event binder for DOM interactions (frontend only)
    // Required for interactive WASM applications
    EventBinder EventBinder

    // File writing function (backend only)
    // Typically os.WriteFile
    WriteFile func(path string, content string) error
}

type ColorScheme struct {
    Primary    string // Main brand color (hex)
    Secondary  string // Accent color (hex)
    Text       string // Text color (hex)
    Background string // Background color (hex)
    Border     string // Border color (hex)
}
```

### 12.2. Site Methods

```go
// New creates a new site instance
// This is the only way to create a Site
func New(cfg *Config) *Site

// NewPage creates and registers a new page
// title: Page title for <title> tag and navigation
// filename: Output filename (e.g., "index.html") or route name for SPA
func (s *Site) NewPage(title, filename string) *Page

// Generate writes all pages and assets to disk (backend only)
// Returns error if file writing fails
func (s *Site) Generate() error

// PageCount returns the number of registered pages
func (s *Site) PageCount() int

// BuildNav generates navigation HTML for all pages
// Automatically called when rendering pages with multiple pages
func (s *Site) BuildNav() string

// AddCSS adds CSS to the site's stylesheet (backend only, no-op in frontend)
// Automatically deduplicates identical CSS blocks
func (s *Site) AddCSS(css string)

// AddJS adds JavaScript to the site's script (backend only, no-op in frontend)
// Automatically deduplicates identical JS blocks
func (s *Site) AddJS(js string)
```

### 12.3. Page Methods

```go
// NewSection creates a new section in the page
// title: Section title (rendered as <h1> in the section)
// Returns *Section for method chaining
func (p *Page) NewSection(title string) *Section

// AddHead adds content to the page's <head> section
// content: Raw HTML to insert in <head> (meta tags, links, etc.)
// Returns *Page for method chaining
func (p *Page) AddHead(content string) *Page

// RenderHTML generates the complete HTML for the page
// Returns the full HTML document as a string
func (p *Page) RenderHTML() string
```

### 12.4. Section Methods

```go
// Add adds a component to the section
// component: Any type that implements HTMLRenderer
// Optionally implements CSSRenderer and/or JSRenderer
// Returns *Section for method chaining
func (s *Section) Add(component any) *Section

// Render generates the section's HTML
// Called automatically by Page.RenderHTML()
func (s *Section) Render() string

// Public fields
Title    string // Section title
ModuleID string // Custom HTML id for the section element
```

### 12.5. Component Interfaces

```go
// HTMLRenderer must be implemented by all components
type HTMLRenderer interface {
    RenderHTML() string
}

// CSSRenderer is optional for components that need custom styles
type CSSRenderer interface {
    RenderCSS() string
}

// JSRenderer is optional for components that need JavaScript
type JSRenderer interface {
    RenderJS() string
}

// EventBinder is implemented by the user for WASM applications
type EventBinder interface {
    EventListener(add bool, elementID, eventType string, callback func())
}
```



### 12.7. Helper Functions

```go
// DefaultColorScheme returns a predefined color scheme
func DefaultColorScheme() *ColorScheme
```

---
## [Migration & Best Practices](examples/migration-and-best-practices.md)


### [Minimal Example](examples/minimal-example.md)

### [Glossary](glossary.md)

### [Resources](resources.md)


