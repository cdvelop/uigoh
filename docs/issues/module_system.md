# Module System & Reflection-Based UI Rendering

**Parent**: [Project Structure](./project_structure.md)  
**Status**: Proposal  
**Created**: 2025-10-23  

---

## 🎯 Overview

The module system provides:
1. **Centralized module registry** (`pkg/modules.go`)
2. **Automatic module discovery** via reflection
3. **Convention-based routing** (struct name → URL path)
4. **Auto-rendering pipeline** for modules with UI

---

## 📦 Module Registry

### `src/pkg/modules.go`

```go
package pkg

import (
    "github.com/cdvelop/monjitaschillan.cl/src/internal/patientCare"
    "github.com/cdvelop/monjitaschillan.cl/src/internal/servicesPage"
    "github.com/cdvelop/monjitaschillan.cl/src/internal/homePage"
)

// Modules contains all registered internal modules
// Add new modules here to make them available system-wide
var Modules = []any{
    &homePage.Module{},
    &servicesPage.Module{},
    &patientCare.Module{},
    // Add more modules as needed
}
```

**Purpose**:
- Single source of truth for all modules
- Easy to add/remove modules
- Used by different subsystems (UI rendering, routing, etc.)

---

## 🔍 Reflection-Based Detection

### Auto-Discovery Process

```go
// Example reflection logic in gosite/discovery.go
package gosite

import (
    "reflect"
    "strings"
)

// ModuleInfo holds discovered module metadata
type ModuleInfo struct {
    Name       string      // e.g., "patientcare"
    StructName string      // e.g., "Module"
    Package    string      // e.g., "patientCare"
    Instance   any
    HasRenderUI bool
}

// DiscoverModules inspects the module registry
func DiscoverModules(modules []any) []ModuleInfo {
    discovered := make([]ModuleInfo, 0, len(modules))
    
    for _, module := range modules {
        info := inspectModule(module)
        discovered = append(discovered, info)
    }
    
    return discovered
}

// inspectModule extracts metadata from a module instance
func inspectModule(module any) ModuleInfo {
    t := reflect.TypeOf(module)
    v := reflect.ValueOf(module)
    
    // Handle pointers
    if t.Kind() == reflect.Ptr {
        t = t.Elem()
    }
    
    // Extract package name from full type path
    // e.g., "github.com/.../internal/patientCare.Module" → "patientCare"
    fullPath := t.PkgPath()
    parts := strings.Split(fullPath, "/")
    packageName := parts[len(parts)-1]
    
    // Convert to lowercase for URL routing
    moduleName := strings.ToLower(packageName)
    
    // Check if RenderUI method exists
    hasRenderUI := false
    if t.Kind() == reflect.Struct {
        method := v.MethodByName("RenderUI")
        hasRenderUI = method.IsValid()
    }
    
    return ModuleInfo{
        Name:       moduleName,
        StructName: t.Name(),
        Package:    packageName,
        Instance:   module,
        HasRenderUI: hasRenderUI,
    }
}
```

---

## 🎨 Auto-Rendering Pipeline

### How It Works

```go
// Example in main.go
func main() {
    // ... setup code ...
    
    // Discover all modules
    modules := gosite.DiscoverModules(pkg.Modules)
    
    // Create main page
    page := ui.UI.NewPage("Monjitas Chillán")
    
    // Render each module that has RenderUI
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
    
    // Generate site
    if err := ui.UI.GenerateSite(); err != nil {
        log.Fatal(err)
    }
}

// callRenderUI uses reflection to invoke RenderUI method
func callRenderUI(module any) (string, error) {
    v := reflect.ValueOf(module)
    method := v.MethodByName("RenderUI")
    
    if !method.IsValid() {
        return "", fmt.Errorf("RenderUI method not found")
    }
    
    // Call with no params (can extend to pass params)
    results := method.Call([]reflect.Value{})
    
    // Expect (string, error) return
    if len(results) != 2 {
        return "", fmt.Errorf("RenderUI should return (string, error)")
    }
    
    html := results[0].Interface().(string)
    
    if !results[1].IsNil() {
        err := results[1].Interface().(error)
        return "", err
    }
    
    return html, nil
}
```

---

## 🗺️ Convention-Based Routing

### Module Name → URL Mapping

| Module Package | Detected Name | URL Path | Section ID |
|----------------|---------------|----------|------------|
| `homePage` | `homepage` | `/` | `#inicio` |
| `servicesPage` | `servicespage` | `/servicios` | `#servicios` |
| `patientCare` | `patientcare` | `/pacientes` | `#pacientes` |
| `patientAppointments` | `patientappointments` | `/citas` | `#citas` |

**Auto-detection Example**:
```go
// internal/patientCare/ui.go
package patientCare

type Module struct{}

func (m *Module) RenderUI(params ...any) (string, error) {
    // Detected name: "patientcare" (lowercase)
    // Auto-generated section ID: "pacientes" (can customize)
    
    page := ui.UI.NewPage("Patient Care")
    
    // Module name available for routing
    moduleName := getModuleName(m) // "patientcare"
    
    content := fmt.Sprintf(`
        <section id="%s" class="page">
            <h1>Atención de Pacientes</h1>
            <!-- content -->
        </section>
    `, moduleName)
    
    page.AddSection(content)
    return page.RenderHTML(), nil
}
```

---

## 📐 Module Interface (Optional)

### Recommended Interface

```go
// pkg/interfaces.go
package pkg

// UIRenderer defines the optional interface for modules with UI
type UIRenderer interface {
    RenderUI(params ...any) (string, error)
}

// NameProvider allows modules to override auto-detected names
type NameProvider interface {
    GetModuleName() string
}

// RouteProvider allows modules to define custom routes
type RouteProvider interface {
    GetRoutes() []Route
}

type Route struct {
    Path    string
    Handler func(w http.ResponseWriter, r *http.Request)
}
```

### Usage Example

```go
// internal/patientCare/module.go
package patientCare

type Module struct {
    service *Service
}

// RenderUI implements UIRenderer
func (m *Module) RenderUI(params ...any) (string, error) {
    // Implementation
}

// GetModuleName implements NameProvider (optional override)
func (m *Module) GetModuleName() string {
    return "pacientes" // Custom name instead of auto-detected "patientcare"
}

// GetRoutes implements RouteProvider (optional)
func (m *Module) GetRoutes() []pkg.Route {
    return []pkg.Route{
        {Path: "/api/patients", Handler: m.handleGetPatients},
        {Path: "/api/patients/:id", Handler: m.handleGetPatient},
    }
}
```

---

## 🔄 Initialization Lifecycle

### Startup Sequence

```
1. main.go starts
   ↓
2. Load pkg.Modules registry
   ↓
3. For each module:
   a. Reflect on struct type
   b. Extract package name → lowercase
   c. Check for RenderUI method
   d. Check for custom name (NameProvider)
   e. Store in ModuleInfo
   ↓
4. Initialize UI manager (ui.UI)
   ↓
5. Create main page
   ↓
6. For each module with RenderUI:
   a. Call RenderUI() via reflection
   b. Collect HTML output
   c. Add to page sections
   d. Auto-accumulate CSS/JS
   ↓
7. Call ui.UI.GenerateSite()
   a. Write index.html
   b. Write style.css
   c. Write script.js
   ↓
8. golite watches and minifies
   ↓
9. Server starts, serves public/
```

---

## 🎛️ Module Configuration

### Per-Module Config (Future)

```go
// pkg/modules.go
type ModuleConfig struct {
    Module  any
    Enabled bool
    Order   int  // Render order
    Name    string // Override auto-detected name
}

var Modules = []ModuleConfig{
    {Module: &homePage.Module{}, Enabled: true, Order: 1},
    {Module: &servicesPage.Module{}, Enabled: true, Order: 2},
    {Module: &patientCare.Module{}, Enabled: false, Order: 3}, // Disabled
}
```

---

## 🧪 Testing Module Discovery

### Unit Test Example

```go
// pkg/modules_test.go
package pkg

import (
    "testing"
    "github.com/cdvelop/monjitaschillan.cl/src/internal/homePage"
)

func TestModuleDiscovery(t *testing.T) {
    modules := []any{
        &homePage.Module{},
    }
    
    discovered := gosite.DiscoverModules(modules)
    
    if len(discovered) != 1 {
        t.Fatalf("Expected 1 module, got %d", len(discovered))
    }
    
    mod := discovered[0]
    
    if mod.Name != "homepage" {
        t.Errorf("Expected name 'homepage', got '%s'", mod.Name)
    }
    
    if mod.Package != "homePage" {
        t.Errorf("Expected package 'homePage', got '%s'", mod.Package)
    }
    
    if !mod.HasRenderUI {
        t.Error("Expected module to have RenderUI method")
    }
}
```

---

## ⚡ Performance Considerations

### Reflection Overhead

**When is reflection used?**
- ✅ Once at startup (module discovery)
- ✅ Once per render cycle (calling RenderUI)
- ❌ NOT in hot paths

**Optimization**:
```go
// Cache method lookups
type ModuleCache struct {
    renderMethod reflect.Value
    nameMethod   reflect.Value
}

var methodCache = make(map[string]*ModuleCache)

func getCachedMethod(module any, methodName string) reflect.Value {
    key := fmt.Sprintf("%T.%s", module, methodName)
    
    if cached, ok := methodCache[key]; ok {
        return cached.renderMethod
    }
    
    v := reflect.ValueOf(module)
    method := v.MethodByName(methodName)
    
    methodCache[key] = &ModuleCache{renderMethod: method}
    
    return method
}
```

---

## 🚧 Edge Cases & Error Handling

### Missing RenderUI Method
```go
if !mod.HasRenderUI {
    log.Printf("⚠️  Module %s has no RenderUI, skipping", mod.Name)
    continue
}
```

### RenderUI Returns Error
```go
html, err := callRenderUI(mod.Instance)
if err != nil {
    log.Printf("❌ Error rendering %s: %v", mod.Name, err)
    // Option A: Skip module
    continue
    // Option B: Render error placeholder
    html = fmt.Sprintf("<div class=\"error\">Failed to load %s</div>", mod.Name)
}
```

### Duplicate Module Names
```go
func validateModuleNames(modules []ModuleInfo) error {
    seen := make(map[string]bool)
    
    for _, mod := range modules {
        if seen[mod.Name] {
            return fmt.Errorf("duplicate module name: %s", mod.Name)
        }
        seen[mod.Name] = true
    }
    
    return nil
}
```

---

## 📝 Adding a New Module

### Step-by-Step

1. **Create module package**
   ```bash
   mkdir -p src/internal/appointments
   ```

2. **Define module struct**
   ```go
   // src/internal/appointments/module.go
   package appointments
   
   type Module struct {
       service *Service
   }
   ```

3. **Implement RenderUI (optional)**
   ```go
   // src/internal/appointments/ui.go
   package appointments
   
   func (m *Module) RenderUI(params ...any) (string, error) {
       page := ui.UI.NewPage("Appointments")
       // Build UI
       return page.RenderHTML(), nil
   }
   ```

4. **Register in pkg/modules.go**
   ```go
   import "github.com/.../internal/appointments"
   
   var Modules = []any{
       &homePage.Module{},
       &appointments.Module{}, // ← Add here
   }
   ```

5. **Run application**
   ```bash
   go run ./src/cmd/appserver
   ```

✅ Module automatically discovered and rendered!

---

## 🎯 Benefits of This System

1. **Zero boilerplate**: Just add to `Modules` slice
2. **Convention over configuration**: Names auto-detected
3. **Type-safe**: Compile-time checks for module structs
4. **Flexible**: Override conventions when needed
5. **Testable**: Each module independently testable
6. **Scalable**: Add 100 modules without changing core code

---

## ❓ Open Questions

### Question 1: Variadic Parameters in RenderUI
```go
func (m *Module) RenderUI(params ...any) (string, error)
```

**What should be passed in `params`?**

**Option A - Context object**:
```go
type RenderContext struct {
    User    *User
    Request *http.Request
}
RenderUI(ctx RenderContext)
```

**Option B - Keep variadic, document conventions**:
```go
// params[0] = *User (optional)
// params[1] = *http.Request (optional)
RenderUI(params ...any)
```

**Recommendation**: Start with **no params** (simpler), add context later if needed.

---

### Question 2: Render Order
How should modules be ordered in the final page?

**Option A - Registry order** (current):
```go
var Modules = []any{
    &homePage.Module{},      // First
    &servicesPage.Module{},  // Second
}
```

**Option B - Explicit priority**:
```go
type Module interface {
    RenderUI() (string, error)
    Priority() int // Lower = earlier
}
```

**Recommendation**: Start with **registry order**, add priority if needed.

---

**Next**: [Migration Strategy](./migration_strategy.md) - Implementation plan
