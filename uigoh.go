
package uigoh

import "github.com/cdvelop/uigoh/components"

// Site is the main entry point for the UI generation system.
// It manages all pages, components, and assets (CSS, JS).
type Site = components.Site

// NewSite creates a new Site manager.
// title: The overall title for the site.
// outputDir: The directory where the generated files will be saved.
func NewSite(title, outputDir string) *Site {
    return components.NewSite(title, outputDir)
}

// Page represents a single HTML page.
type Page = components.Page

// SectionBuilder helps in constructing a section within a page.
type SectionBuilder = components.SectionBuilder

// CarouselImage defines the structure for an image in a carousel.
type CarouselImage = components.CarouselImage

// FormConfig holds the configuration for a form.
type FormConfig = components.FormConfig

// FormField defines a single field within a form.
type FormField = components.FormField

// NavItem represents a single item in the navigation.
type NavItem = components.NavItem

// Note on Usage:
//
// To correctly build a UI, you must use the methods provided by the Site object.
// Standalone functions for components are not provided because the Site object
// is responsible for collecting and managing the CSS and JavaScript assets
// required by each component.
//
// Example:
//
// site := uigoh.NewSite("My Awesome Site", "dist")
// homePage := site.NewPage(nil, "Home")
// section := homePage.Section("Welcome")
// section.AddCard("Title", "Description", "icon-name")
// site.GenerateSite()
