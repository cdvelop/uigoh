# API Design Principles

## 1. Core Philosophy: Write Once, Run Anywhere

The fundamental principle of `gosite` is that you define your site structure and components a single time in shared Go code. The framework then adapts its behavior based on the compilation target:

- **Backend (`//go:build !wasm`)**: Generates static HTML, CSS, and JS files.
- **Frontend (`//go:build wasm`)**: Renders the site dynamically in the browser's DOM.

This approach eliminates code duplication and ensures a consistent structure across both server-side and client-side environments.

## 2. Fluent, Chained API

The API is designed to be fluent and intuitive. All instantiation is controlled through a chained API, starting from a single entry point: `gosite.New()`.

```go
// 1. Create a Site (the only public constructor)
site := gosite.New(&gosite.Config{...})

// 2. Chain methods to build pages and sections
page := site.NewPage("Home", "index.html")
section := page.NewSection("Welcome")

// 3. Add components to sections
section.Add(component1).Add(component2)
```

Direct instantiation of `Page` or `Section` is disallowed by keeping their struct fields private.

## 3. Environment-Specific Logic with Build Tags

To separate backend and frontend concerns, the codebase uses a clear file-naming convention with Go build tags:

- **`env.backend.go` (`//go:build !wasm`)**: Contains backend-only logic, such as file generation (`site.Generate()`) and asset bundling.
- **`env.frontend.go` (`//go:build wasm`)**: Contains frontend-only logic, such as DOM rendering and event binding.
- **Other files (no build tags)**: Contain shared logic, including struct definitions (`Site`, `Page`, `Section`), component interfaces, and the core building methods (`NewPage`, `NewSection`, `Add`).

## 4. Component-Based Architecture

The smallest unit of the UI is the **component**. Every component must implement at least the `HTMLRenderer` interface.

- **`HTMLRenderer`**: `RenderHTML() string` (Required)
- **`CSSRenderer`**: `RenderCSS() string` (Optional)
- **`JSRenderer`**: `RenderJS() string` (Optional)

When a component is added to a section, the `gosite` framework automatically collects and deduplicates its CSS and JS for final bundling (in the backend).

## 5. Decoupling with Interfaces

To maintain a clean architecture and prevent circular dependencies, `gosite` uses two key interfaces to mediate communication between different parts of the framework.

### `EventBinder`: Abstracting DOM Interaction

For frontend (`wasm`) builds, components often need to respond to user interactions like clicks or input changes. To keep the core component logic independent of the browser's `syscall/js` API, `gosite` uses the `EventBinder` interface.

- **Purpose**: To provide a formal, testable way for components to request DOM event listeners without directly depending on browser-specific code.
- **Implementation**: The user provides a concrete implementation of this interface in their `gosite.Config` when setting up a frontend application.

```go
// The EventBinder interface, defined in interfaces.go
type EventBinder interface {
    EventListener(add bool, elementID, eventType string, callback func())
}
```

### `SiteLink`: Preventing Circular Dependencies

Components, pages, and sections often need to communicate "up" to the main `Site` object—for example, to add CSS/JS assets or get site-wide information. A direct dependency (`*gosite.Site`) would create a circular import cycle.

The `SiteLink` interface solves this by exposing only a minimal, stable set of methods from the `Site` object.

- **Purpose**: To allow child elements (like components and pages) to safely interact with the site manager without creating import cycles.
- **Implementation**: The `*gosite.Site` object implements this interface, and a `SiteLink` is passed down to its children.

```go
// The SiteLink interface, defined in interfaces.go
type SiteLink interface {
    PageCount() int
    BuildNav() string
    AddCSS(css string)
    AddJS(js string)
}
```

## 6. Dependency Constraints for TinyGo/WASM

To ensure minimal binary sizes and compatibility with TinyGo, the core `gosite` library adheres to a strict dependency policy:

- **`tinystring` is mandatory**: All string manipulation, formatting, and path operations **must** use the `github.com/cdvelop/tinystring` package.
- **Forbidden Packages**: The use of standard packages like `fmt`, `strings`, `strconv`, `errors`, `path`, and `bytes` is **disallowed**.
- **Error Handling**: Errors are handled by returning `(bool, string)` tuples or simple error strings, avoiding the standard `error` interface.
