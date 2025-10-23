Prompt:

> I have the following Go code that re-exports structures from the components package:

package gosite

import "github.com/cdvelop/gosite/components"

type Site = components.Site

func NewSite(title, outputDir string) *Site {
    return components.NewSite(title, outputDir)
}

type Page = components.Page
type SectionBuilder = components.SectionBuilder
type CarouselImage = components.CarouselImage
type FormConfig = components.FormConfig
type FormField = components.FormField
type NavItem = components.NavItem

I need this code to be restructured and redesigned as follows:

NewSite should now receive a pointer to a Config struct (e.g. *Config) instead of title and outputDir parameters.

The Config struct must have all its fields public so it can be modified easily from outside.

NewSite should return a pointer to Page, which will act as the root page of the website.

The returned Page must have methods to add sections (for example: NewSection()).

Each Section should have its own builder-like API to add nested elements such as forms, carousels, and other common UI components — following a hierarchical and chained structure, similar to HTML nesting (e.g. page.Section().AddForm().AddField()).

The design must support method chaining and be easy to extend with new component types in the future.

Keep the naming and structure idiomatic to Go (avoid unnecessary abstractions, prefer explicitness, etc.).


Please rewrite the code according to these requirements, keeping it clean, modular, and following Go best practices.




---

Would you like me to include an example of the expected chained API syntax (like how the user should call it, e.g. site.Page().Section().Form()…) inside the prompt too? That can help the model understand the structure better.
