
package main

import (
	"log"

	"github.com/cdvelop/uigoh"
	"github.com/cdvelop/uigoh/components"
)

func main() {
	// 1. Configure the site
	cfg := &components.Config{
		Title:     "My Awesome Site",
		OutputDir: "dist",
	}

	// 2. Create a new site; NewSite returns the root page (index.html)
	page := uigoh.NewSite(cfg)
	page.Title = "Home Page"

	// 3. Add a head tag to the page
	page.AddHead(`<meta name="description" content="Welcome to the awesome site!">`)

	// 4. Create a new section on the root page
	introSection := page.NewSection("Welcome to Our Website")

	// 5. Add components to the section using the builder API
	introSection.Add(&components.Card{
		Title:       "Feature One",
		Description: "This is a description for the first amazing feature.",
		Icon:        "icon-star",
	}).Add(&components.Card{
		Title:       "Feature Two",
		Description: "Discover the second feature that will change your life.",
		Icon:        "icon-heart",
	})

	// 6. Add another section with a form
	contactSection := page.NewSection("Contact Us")
	contactSection.Add(&components.Form{
		Config: components.FormConfig{
			Action: "/submit-form",
			Method: "POST",
			Fields: []components.FormField{
				{Type: "text", Name: "username", Placeholder: "Your Name", Required: true},
				{Type: "email", Name: "email", Placeholder: "Your Email", Required: true},
				{Type: "textarea", Name: "message", Placeholder: "Your Message"},
			},
		},
	})

	// 7. Generate the site
	if err := page.Generate(); err != nil {
		log.Fatalf("Error generating site: %v", err)
	}

	log.Println("Site generated successfully in 'dist' directory!")
}
