# Component Extraction Task

## Context

The `gosite` library is a Go-based framework for building websites that compile to both:
- **Backend** (`//go:build !wasm`): Generates static HTML, CSS, and JS files
- **Frontend** (`//go:build wasm`): Renders dynamically in the browser's DOM

For complete API design details, see [API_DESIGN.md](../API_DESIGN.md).

## Current State

### Existing Components (Limited)
Currently, we have only a few components in the `components/` directory:
- `card/` - Card component with title, description, and icon
- `carousel/` - Carousel/slider component
- `form/` - Form component

### Available Templates
We have two complete website templates in `templates/`:

1. **hospital-website/** - A modern hospital website template
   - ✅ Excellent mobile and desktop responsive design
   - Contains multiple reusable UI patterns (navbar, hero section, service cards, doctor panels, package cards, contact forms, footer, etc.)

2. **platform/** - A platform/dashboard template
   - ⚠️ Desktop-focused (mobile version needs refactoring)
   - Used in previous systems
   - Contains different UI patterns suitable for admin/platform interfaces

## Task Requirements

**Extract and structure reusable components from these two templates following the gosite API architecture.**

Components must:
1. Implement the required interfaces (`HTMLRenderer`, `CSSRenderer`, `JSRenderer`)
2. Follow the gosite API design principles
3. Be easy to use through the fluent, chained API
4. Work in both backend (static generation) and frontend (WASM) environments
5. Use **only** `github.com/cdvelop/tinystring` for string manipulation (NO `fmt`, `strings`, `strconv`, etc.)

## Component Interface Requirements

Every component must implement at least `HTMLRenderer`:

```go
type HTMLRenderer interface {
    RenderHTML() string
}
```

Optionally, components can implement:

```go
type CSSRenderer interface {
    RenderCSS() string
}

type JSRenderer interface {
    RenderJS() string
}
```

## Example Component Structure

**CRITICAL: File Organization Pattern**

Each component MUST follow this exact structure:

```
components/
└── category/
    └── componentname/
        ├── componentname.go      # Main component file (REQUIRED)
        ├── style.css             # CSS styles (if component has CSS)
        ├── script.js             # JavaScript code (if component is interactive)
        └── env.back.go           # Backend embed file (if CSS/JS exist)
```

### Required Files

#### 1. `componentname.go` - Main Component File
Contains the component struct and RenderHTML() method ONLY:

```go
package componentname

import (
    . "github.com/cdvelop/tinystring"
)

// ComponentName implements HTMLRenderer, CSSRenderer, and JSRenderer interfaces.
// Brief description of what the component does.
type ComponentName struct {
    // Public fields for configuration
    Title       string
    Description string
    CSSClass    string
}

// RenderHTML generates the HTML for the component.
func (c *ComponentName) RenderHTML() string {
    // Use tinystring methods: Fmt(), Convert(), EscapeHTML(), EscapeAttr()
    // Build and return HTML string
    
    class := "component-name"
    if c.CSSClass != "" {
        class += " " + c.CSSClass
    }
    classEsc := Convert(class).EscapeAttr()
    
    titleEsc := Convert(c.Title).EscapeHTML()
    descEsc := Convert(c.Description).EscapeHTML()
    
    tpl := `<div class="%s">
    <h3>%s</h3>
    <p>%s</p>
</div>
`
    
    return Fmt(tpl, classEsc, titleEsc, descEsc)
}
```

#### 2. `style.css` - Component Styles (OPTIONAL)
Create ONLY if component needs CSS:

```css
/* Component: ComponentName */

.component-name {
  padding: 2rem;
  background: var(--bg-color);
}

.component-name h3 {
  font-size: 2rem;
  margin-bottom: 1rem;
}

/* Responsive */
@media (min-width: 768px) {
  .component-name {
    padding: 3rem;
  }
}
```

#### 3. `script.js` - Component JavaScript (OPTIONAL)
Create ONLY if component needs interactivity:

```javascript
// Component: ComponentName
(function() {
  const elements = document.querySelectorAll('.component-name');
  
  elements.forEach(el => {
    el.addEventListener('click', function() {
      // Handle interaction
    });
  });
})();
```

#### 4. `env.back.go` - Backend Embed File (REQUIRED if CSS or JS exist)
Embeds CSS and JS files for backend compilation:

```go
//go:build !wasm
// +build !wasm

package componentname

import (
	_ "embed"
)

//go:embed style.css
var styleCss string

// RenderCSS returns the CSS for the component.
func (c *ComponentName) RenderCSS() string {
	return styleCss
}

//go:embed script.js
var scriptJs string

// RenderJS returns the JavaScript for the component.
func (c *ComponentName) RenderJS() string {
	return scriptJs
}
```

**IMPORTANT RULES:**
1. Main `.go` file contains ONLY struct definition and RenderHTML()
2. CSS goes in separate `style.css` file (if needed)
3. JavaScript goes in separate `script.js` file (if needed)
4. `env.back.go` embeds CSS/JS files (create only if CSS or JS exist)
5. DO NOT inline CSS or JS strings in the main `.go` file

## Required Analysis and Deliverables

### 1. Component Inventory
Analyze both templates (`hospital-website/` and `platform/`) and provide a comprehensive table with:

| Component Name | Template Source | Description | Properties | Responsive | Interactive |
|---------------|-----------------|-------------|------------|------------|-------------|
| ... | ... | ... | ... | ✅/⚠️ | Yes/No |

Include:
- Component name (descriptive, clear)
- Template source (hospital-website, platform, or both)
- Visual/functional description (1-2 sentences)
- Key configurable properties (list main fields)
- Responsive behavior (✅ mobile+desktop, ⚠️ desktop-only, ⚠️ needs-work)
- Interactive elements (Yes if requires JS/WASM)

### 2. Priority Matrix
Rank ALL identified components in three priority categories:

**HIGH PRIORITY** (common, simple, highly reusable):
- Component Name - Justification (1 sentence why it's high priority)
- ...

**MEDIUM PRIORITY** (useful but more specific):
- Component Name - Justification
- ...

**LOW PRIORITY** (complex, template-specific, or requires refactoring):
- Component Name - Justification
- ...

### 3. Directory Structure Proposal
Provide the exact directory structure for `components/` showing how ALL components will be organized:

```
components/
├── category1/
│   ├── component1/
│   │   └── component1.go
│   └── component2/
│       └── component2.go
└── category2/
    └── component3/
        └── component3.go
```

Justify your organization strategy in 2-3 sentences.

### 4. Component Dependency Graph
List any composite components and their dependencies:

```
CompositeComponentName
├── DependsOn: BaseComponent1
├── DependsOn: BaseComponent2
└── Strategy: [direct import / registry / other]
```

If no composite components exist, state: "No composite components identified."

### 5. Responsive Design Strategy
Choose ONE strategy and justify:

**Strategy:** [Embedded media queries / Component variants / CSS class overrides / Configurable breakpoints]

**Justification:** (2-3 sentences explaining why this strategy is best for gosite)

**Implementation Pattern:**
```go
// Show code example of how this strategy will be implemented
```

### 6. Styling Architecture
Choose ONE approach and provide specifications:

**Approach:** [CSS variables / BEM / Scoped / Utility-first]

**CSS Variable Naming Convention:** (if using CSS variables)
```css
/* Example */
--component-name-property: value;
```

**Class Naming Pattern:**
```css
/* Example */
.component-name { }
.component-name__element { }
.component-name--modifier { }
```

### 7. Interactive Components List
Provide a table of components that need JavaScript/WASM:

| Component | Interaction Type | EventBinder Usage | State Management |
|-----------|------------------|-------------------|------------------|
| Navbar | Toggle menu | Yes | Local boolean |
| ... | ... | ... | ... |

### 8. Conversion Workflow
Step-by-step checklist for converting ONE component from template to gosite:

**Phase 1: Analysis**
- [ ] Identify component HTML structure in template
- [ ] Extract all CSS rules related to the component
- [ ] Identify any JavaScript functionality needed
- [ ] List all configurable properties (fields)

**Phase 2: Directory Setup**
- [ ] Create component directory: `components/category/componentname/`
- [ ] Create main file: `componentname.go`
- [ ] Create `style.css` (if component has CSS)
- [ ] Create `script.js` (if component is interactive)
- [ ] Create `env.back.go` (if CSS or JS exist)

**Phase 3: Implementation**
- [ ] Write component struct with all properties in `componentname.go`
- [ ] Implement `RenderHTML()` method using tinystring
- [ ] Copy and clean CSS to `style.css`
- [ ] Copy and wrap JavaScript in IIFE in `script.js`
- [ ] Set up embed directives in `env.back.go`
- [ ] Add component documentation comments

**Phase 4: Validation**
- [ ] Verify no compilation errors
- [ ] Test HTML output format
- [ ] Verify CSS is properly embedded
- [ ] Verify JS is properly embedded (if exists)
- [ ] Check responsive behavior in CSS

### 9. Documentation Standard
Choose the documentation approach for EACH component:

**Required:**
- [ ] Inline code comments for public methods
- [ ] Package-level documentation
- [ ] README.md in component folder


State which are REQUIRED vs OPTIONAL and provide a template.



### 11. Naming Conventions (SPECIFY EXACT RULES)
Provide exact naming rules:

**Component Struct:** `[Exact pattern with example]`
**File Name:** `[Exact pattern with example]`
**Package Name:** `[Exact pattern with example]`
**CSS Classes:** `[Exact pattern with example]`
**Properties:** `[Exact pattern with example]`

### 12. Implementation Roadmap
Provide a phased implementation plan:

**Phase 1: Foundation (Components to extract first)**
1. Component Name (Reason: ...)
2. Component Name (Reason: ...)
...

**Phase 2: Core Features**
1. Component Name (Reason: ...)
...

**Phase 3: Advanced Components**
1. Component Name (Reason: ...)
...

### 13. Platform Template Decision
Choose ONE approach for the desktop-only `platform/` template:

**Decision:** [Fix mobile first / Extract as desktop-only / Skip for now]

**Action Plan:**
- Step 1: ...
- Step 2: ...

## Submission Format

Provide all answers in a single, well-structured markdown document with clear headings matching the sections above. Be specific, concrete, and actionable. No open-ended questions or suggestions—make definitive decisions and justify them briefly.
