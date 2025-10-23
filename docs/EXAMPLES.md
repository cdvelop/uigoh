# GoSite - Complete Examples

**Updated**: 2025-10-23  

---

## 📦 Package Setup

```go
// src/pkg/ui.go
package pkg

import "github.com/yourorg/gosite"

// Global UI instance
var UI *gosite.Site
```

---

## 🏠 Example 1: Homepage (SPA Section)

```go
// src/internal/homepage/homepage.go
package homepage

import "github.com/yourorg/project/pkg"

type Homepage struct{}

type Feature struct {
    Title       string
    Description string
    Icon        string
}

func (h *Homepage) GetFeatures() []Feature {
    return []Feature{
        {
            Title:       "Atención Personalizada",
            Description: "Cuidado dedicado a cada paciente",
            Icon:        "icon-care",
        },
        {
            Title:       "Equipo Profesional",
            Description: "Personal altamente capacitado",
            Icon:        "icon-staff",
        },
        {
            Title:       "Instalaciones Modernas",
            Description: "Espacios cómodos y equipados",
            Icon:        "icon-building",
        },
    }
}

func (h *Homepage) RenderUI(context ...any) string {
    // Create section (ID auto-detected as "homepage")
    section := pkg.UI.Section("Inicio")
    
    // Add feature cards
    features := h.GetFeatures()
    for _, feat := range features {
        section.AddCard(feat.Title, feat.Description, feat.Icon)
    }
    
    // Add carousel
    section.AddCarousel([]pkg.CarouselImage{
        {Src: "img/facility1.jpg", Alt: "Instalación 1"},
        {Src: "img/facility2.jpg", Alt: "Instalación 2"},
        {Src: "img/facility3.jpg", Alt: "Instalación 3"},
    })
    
    // Return HTML → added to index.html as section
    return section.Render()
}
```

---

## 🛠️ Example 2: Services (Separate Page)

```go
// src/internal/services/services.go
package services

import "github.com/yourorg/project/pkg"

type Services struct{}

type Service struct {
    Title       string
    Description string
    Icon        string
    Price       string
}

func (s *Services) GetServices() []Service {
    return []Service{
        {
            Title:       "Medicina General",
            Description: "Atención médica primaria y consultas generales",
            Icon:        "icon-medicine",
            Price:       "$15.000",
        },
        {
            Title:       "Curaciones",
            Description: "Manejo profesional de heridas y curaciones",
            Icon:        "icon-bandage",
            Price:       "$8.000",
        },
        {
            Title:       "Control de Signos Vitales",
            Description: "Monitoreo de presión, temperatura, y más",
            Icon:        "icon-heartbeat",
            Price:       "$5.000",
        },
    }
}

func (s *Services) RenderUI(context ...any) string {
    // Create separate page (filename auto-detected as "services.html")
    page := pkg.UI.NewPage(s, "Nuestros Servicios")
    
    // Create section within page (auto-added, no AddSection needed!)
    section := page.Section("Servicios Disponibles")
    
    // Add service cards
    services := s.GetServices()
    for _, svc := range services {
        section.AddCard(
            svc.Title,
            svc.Description + " - " + svc.Price,
            svc.Icon,
        )
    }
    
    // Return empty string → page already registered
    // Link automatically added to index.html navigation
    return ""
}
```

---

## 📞 Example 3: Contact (Separate Page with Form)

```go
// src/internal/contact/contact.go
package contact

import "github.com/yourorg/project/pkg"

type Contact struct{}

type ContactInfo struct {
    Phone   string
    Email   string
    Address string
}

func (c *Contact) GetInfo() ContactInfo {
    return ContactInfo{
        Phone:   "+56 9 1234 5678",
        Email:   "contacto@monjitas.cl",
        Address: "Avenida Principal 123, Chillán",
    }
}

func (c *Contact) RenderUI(context ...any) string {
    // Create separate page (filename: "contact.html")
    page := pkg.UI.NewPage(c, "Contáctanos")
    
    // Contact info section
    infoSection := page.Section("Información de Contacto")
    info := c.GetInfo()
    infoSection.AddCard("Teléfono", info.Phone, "icon-phone")
    infoSection.AddCard("Email", info.Email, "icon-email")
    infoSection.AddCard("Dirección", info.Address, "icon-location")
    
    // Contact form section
    formSection := page.Section("Envíanos un Mensaje")
    formSection.AddForm(pkg.FormConfig{
        Action: "/api/contact",
        Method: "POST",
        Fields: []pkg.FormField{
            {Type: "text", Name: "name", Placeholder: "Tu nombre", Required: true},
            {Type: "email", Name: "email", Placeholder: "Tu email", Required: true},
            {Type: "text", Name: "phone", Placeholder: "Tu teléfono", Required: false},
            {Type: "textarea", Name: "message", Placeholder: "Tu mensaje", Required: true},
        },
    })
    
    // Return empty → separate page created
    return ""
}
```

---

## 👨‍⚕️ Example 4: Staff (SPA Section with Context)

```go
// src/internal/staffpage/staff.go
package staffpage

import "github.com/yourorg/project/pkg"

type StaffPage struct{}

type Staff struct {
    Name        string
    Role        string
    Photo       string
    Description string
}

// UserContext example - shows how modules can use context
type UserContext struct {
    IsAdmin bool
}

func (s *StaffPage) GetStaff() []Staff {
    return []Staff{
        {
            Name:        "Dra. María González",
            Role:        "Directora Médica",
            Photo:       "img/staff/maria.jpg",
            Description: "Especialista en medicina familiar",
        },
        {
            Name:        "Enf. Juan Pérez",
            Role:        "Enfermero Jefe",
            Photo:       "img/staff/juan.jpg",
            Description: "10 años de experiencia",
        },
    }
}

func (s *StaffPage) RenderUI(context ...any) string {
    // Extract user context if provided
    var user UserContext
    if len(context) > 0 {
        if u, ok := context[0].(UserContext); ok {
            user = u
        }
    }
    
    // Create section (ID: "staffpage")
    section := pkg.UI.Section("Nuestro Equipo")
    
    // Add staff cards
    staff := s.GetStaff()
    for _, member := range staff {
        cardDesc := member.Description
        
        // Admins see additional info
        if user.IsAdmin {
            cardDesc += " [Admin: Contact info available]"
        }
        
        section.AddCard(
            member.Name + " - " + member.Role,
            cardDesc,
            "icon-user",
        )
    }
    
    // Return HTML → SPA section
    return section.Render()
}
```

---

## 🚀 Example 5: Main Application

```go
// src/cmd/appserver/main.go
package main

import (
    "log"
    "net/http"
    
    "github.com/yourorg/project/pkg"
    "github.com/yourorg/project/internal/homepage"
    "github.com/yourorg/project/internal/services"
    "github.com/yourorg/project/internal/contact"
    "github.com/yourorg/project/internal/staffpage"
    "github.com/yourorg/gosite"
)

func main() {
    // Initialize UI system
    pkg.UI = gosite.NewSite("Monjitas Chillán", "src/web/ui/")
    
    // Register modules
    pkg.Modules = []any{
        &homepage.Homepage{},    // Will be section in index.html
        &services.Services{},    // Will create services.html
        &staffpage.StaffPage{},  // Will be section in index.html
        &contact.Contact{},      // Will create contact.html
    }
    
    // Generate all UI files
    for _, mod := range pkg.Modules {
        pkg.UI.AddSection(mod) // Auto-detects if section or page
    }
    
    // Write files to disk
    if err := pkg.UI.GenerateSite(); err != nil {
        log.Fatal("Failed to generate site:", err)
    }
    
    log.Println("✅ Site generated successfully:")
    log.Println("   - index.html (with homepage & staff sections)")
    log.Println("   - services.html (separate page)")
    log.Println("   - contact.html (separate page)")
    log.Println("   - style.css (shared)")
    log.Println("   - main.js (shared)")
    
    // Start HTTP server
    fs := http.FileServer(http.Dir("src/web/public"))
    http.Handle("/", noCacheMiddleware(fs))
    
    log.Println("🚀 Server running at http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

// noCacheMiddleware prevents browser caching during development
func noCacheMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
        w.Header().Set("Pragma", "no-cache")
        w.Header().Set("Expires", "0")
        next.ServeHTTP(w, r)
    })
}
```

---

## 📊 Expected Output

### File Structure
```
src/web/ui/
├── index.html       (267 KB - SPA with 2 sections)
├── services.html    (134 KB - Separate page)
├── contact.html     (156 KB - Separate page)
├── style.css        (45 KB - Shared, deduplicated)
└── main.js          (12 KB - Shared, deduplicated)
```

### index.html Navigation
```html
<nav class="main-nav">
  <!-- SPA Sections (anchor links) -->
  <a href="#homepage" class="nav-link">
    <span>Inicio</span>
  </a>
  <a href="#staffpage" class="nav-link">
    <span>Nuestro Equipo</span>
  </a>
  
  <!-- Separate Pages (file links) -->
  <a href="services.html" class="nav-link">
    <span>Nuestros Servicios</span>
  </a>
  <a href="contact.html" class="nav-link">
    <span>Contáctanos</span>
  </a>
</nav>
```

---

## 🧪 Testing Example

```go
// src/internal/homepage/homepage_test.go
package homepage

import (
    "strings"
    "testing"
    
    "github.com/yourorg/project/pkg"
    "github.com/yourorg/gosite"
)

func TestHomepageRenderUI(t *testing.T) {
    // Setup
    pkg.UI = gosite.NewSite("Test Site", "test/output")
    h := &Homepage{}
    
    // Execute
    html := h.RenderUI()
    
    // Verify
    if html == "" {
        t.Error("RenderUI should return HTML for SPA section")
    }
    
    if !strings.Contains(html, "Atención Personalizada") {
        t.Error("HTML should contain feature titles")
    }
    
    if !strings.Contains(html, "icon-care") {
        t.Error("HTML should contain feature icons")
    }
    
    if !strings.Contains(html, `id="homepage"`) {
        t.Error("Section should have auto-generated ID")
    }
}
```

---

## 📝 Notes

### Module Return Values
- **Return HTML** → Module is SPA section
- **Return ""** → Module created separate page

### Auto-Generation
- Section IDs: From module package name
- Page filenames: From struct name
- Navigation: Combined from all sources

### CSS/JS Sharing
- All pages reference same `style.css` and `main.js`
- Deduplicated at site level
- Only one copy of each component's styles/scripts

---

## 🔗 Related Documents

- [Architecture Summary](../ARCHITECTURE.md)
- [UI Component System](./issues/ui_component_system.md)
- [Final Decisions](./issues/decisions_final.md)
