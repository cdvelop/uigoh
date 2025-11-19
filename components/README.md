# GoSite Components

## Component Structure

Each component follows this **exact** file structure:

```
components/
└── category/
    └── componentname/
        ├── componentname.go    # REQUIRED: Component logic
        ├── style.css           # OPTIONAL: Component styles
        ├── script.js           # OPTIONAL: Component interactivity
        └── env.back.go         # REQUIRED if CSS/JS exist
```

## File Requirements

### 1. Main Component File (`componentname.go`)

**Contains ONLY:**
- Package declaration
- Import of `tinystring`
- Struct definition with public fields
- `RenderHTML()` method

**Must NOT contain:**
- `RenderCSS()` method (goes in `env.back.go`)
- `RenderJS()` method (goes in `env.back.go`)
- Inline CSS or JS strings

**Example:**
```go
package componentname

import (
	. "github.com/cdvelop/tinystring"
)

// ComponentName brief description.
type ComponentName struct {
	Title    string
	CSSClass string
}

// RenderHTML generates the HTML.
func (c *ComponentName) RenderHTML() string {
	classEsc := Convert("component " + c.CSSClass).EscapeAttr()
	titleEsc := Convert(c.Title).EscapeHTML()
	
	return Fmt(`<div class="%s"><h3>%s</h3></div>`, classEsc, titleEsc)
}
```

### 2. Style File (`style.css`) - OPTIONAL

Create only if component needs CSS.

```css
/* Component: ComponentName */

.component {
  padding: 2rem;
}

@media (min-width: 768px) {
  .component {
    padding: 3rem;
  }
}
```

### 3. Script File (`script.js`) - OPTIONAL

Create only if component needs JavaScript. Always wrap in IIFE.

```javascript
// Component: ComponentName
(function() {
  const elements = document.querySelectorAll('.component');
  elements.forEach(el => {
    el.addEventListener('click', function() {
      // Handle interaction
    });
  });
})();
```

### 4. Backend Embed File (`env.back.go`) - REQUIRED if CSS/JS exist

```go
//go:build !wasm
// +build !wasm

package componentname

import (
	_ "embed"
)

//go:embed style.css
var styleCss string

func (c *ComponentName) RenderCSS() string {
	return styleCss
}

//go:embed script.js
var scriptJs string

func (c *ComponentName) RenderJS() string {
	return scriptJs
}
```

**Note:** Only include embed directives for files that exist.

## Component Categories

- **navigation/** - Navigation bars, menus, breadcrumbs
- **layout/** - Headers, footers, heroes, banners, sections
- **content/** - Cards, panels, lists, media displays
- **forms/** - Form components, inputs, validation

## Naming Conventions

- **Struct:** PascalCase (e.g., `ServiceCard`)
- **Package:** lowercase (e.g., `servicecard`)
- **File:** lowercase + `.go` (e.g., `servicecard.go`)
- **CSS Class:** kebab-case (e.g., `.service-card`)

## String Manipulation Rules

**ONLY use `github.com/cdvelop/tinystring`:**
- `Fmt()` - String formatting
- `Convert()` - String conversion
- `.EscapeHTML()` - Escape HTML entities
- `.EscapeAttr()` - Escape HTML attributes
- `Error` handling with `tinystring` methods

**NEVER use:** `fmt`, `strings`, `strconv`
