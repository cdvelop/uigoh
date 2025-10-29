# Gosite API Structure & Design

## 1. Abstract

The goal of `gosite` is to provide a programmatic Go framework for building a wide range of websites, from backend-rendered Server-Side Rendering (SSR) or Multi-Page Applications (MPA) to frontend-rendered WebAssembly (WASM) Single-Page Applications (SPA).

The core principles are:
- **Unified Go Codebase**: Write frontend and backend logic in Go.
- **Performance**: Leverage build tags (`wasm` and `!wasm`) to separate logic, ensuring that frontend binaries (WASM) are minimal and do not include backend-specific code (like file system operations).
- **TinyGo Compatibility**: Strictly use `tinystring` for all string manipulations to minimize binary size, avoiding standard libraries like `fmt`, `strings`, and `strconv`.
- **Fluent API**: Offer a clean, chained API for building pages and components, enforcing a logical structure.
- **Decoupled DOM Interaction**: Abstract all browser/DOM APIs through interfaces, keeping the core `gosite` logic platform-agnostic.

## 2. Core Concepts

The structure of a website is hierarchical: a `Site` contains `Pages`, each `Page` contains `Sections`, and each `Section` contains `Components`.

- **`Site`**: The root object for a website. It holds global configuration, manages pages, and orchestrates the final build process. It is the main entry point.
- **`Page`**: Represents a single HTML page (e.g., `index.html`, `about.html`). It manages its own content, composed of sections.
- **`Section`**: A logical division within a page (e.g., a hero, a gallery, a contact form). It acts as a container for components and helps structure the page layout.
- **`Component`**: A reusable UI element (e.g., a card, a button, a form). Components are the smallest building blocks and are added to sections.

## 3. API Design

To enforce a clear and logical workflow, `gosite` will use public structs with private fields, where instantiation is controlled via a fluent, chained API.

### 3.1. Entry Point & Chaining

The user starts by creating a `Site` instance. From there, all other objects are created through chained method calls. Direct instantiation of `Page` or `Section` is not intended.

```go
// The only public constructor.
// Returns a public struct `Site` with private fields.
site := gosite.New(&gosite.Config{...})

// `NewPage` is a method on `Site` and returns a `*Page`.
page1 := site.NewPage("Home", "index.html")

// `AddSection` is a method on `Page` and returns a `*Section`.
// Methods on `Section` return the `*Section` pointer to allow further chaining.
page1.AddSection().
    AddCard(cardData1).
    AddCard(cardData2)

page1.AddSection().
    AddCustomComponent(...)
```

### 3.2. Struct Visibility

- **`gosite.Site`**: Public struct, but fields like `pages`, `cssBlocks`, etc., are private.
- **`gosite.Page`**: Public struct, but fields like `site`, `sections`, etc., are private.
- **`gosite.Section`**: Public struct, with private fields.
- **`gosite.Component`**: Components are defined in their own packages (e.g., `components/card`). They are added to sections via methods.

## 4. Build Architecture (Backend vs. Frontend)

The codebase will be split to isolate backend-only logic from the frontend WASM binary. This is achieved using Go's build tags.

### 4.1. File Naming Convention

We will adopt a clear file naming convention to separate concerns:
- **`env.backend.go`**: Contains code that should only be compiled for the backend. It will be guarded by `//go:build !wasm`. This includes file writing (`Generate`), CSS/JS aggregation, and other server-side logic.
- **`env.frontend.go`**: Contains code specific to the WASM environment. It will be guarded by `//go:build wasm`. This includes the logic that interacts with the `EventBinder`.
- **Files without `env` prefix**: These files contain shared logic, such as struct definitions (`Site`, `Page`), interfaces, and component logic that is common to both environments.

### 4.2. Example: `Site` Struct

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
    EventListener(add bool, elementID, eventType string, callback func())
}
```

- `add bool`: If `true`, the listener is added. If `false`, it is removed.
- `elementID string`: The `id` of the target HTML element.
- `eventType string`: The event to listen for (e.g., "click", "change").
- `callback func()`: The Go function to execute when the event fires.

### 5.2. Usage in Components

A component that needs to be interactive will use the `EventBinder` from the site's configuration.

```go
// Example of a button component
func (s *Section) AddButton(text, id string, onClick func()) {
    // ... logic to generate button HTML with the given id ...

    // Register the event
    if s.page.site.Cfg.EventBinder != nil {
        s.page.site.Cfg.EventBinder.EventListener(true, id, "click", onClick)
    }
}
```

## 6. Dependencies

- **`tinystring`**: This is a hard requirement. All internal string operations, conversions, and concatenations must use the `tinystring` package to ensure compatibility with TinyGo and to produce the smallest possible WASM binaries. The use of `fmt`, `strings`, `strconv`, `errors`, `path`, or `bytes` is disallowed in the core framework.
