
package main

import (
	"log"

	"github.com/cdvelop/gosite"
	"github.com/cdvelop/gosite/components"
)

func main() {
	// 1. Configure the site
	cfg := &gosite.Config{
		Title:     "My Awesome Site",
		OutputDir: "dist",
	}

	// 2. Create a new site
	site := gosite.NewSite(cfg)

	// 3. Create the home page
	homePage := site.NewPage("Home Page", "index.html")
	homePage.AddHead(`<meta name="description" content="Welcome to the awesome site!">`)

	// 4. Add content to the home page
	introSection := homePage.NewSection("Welcome to Our Website")
	introSection.Add(&components.Card{
		Title:       "Feature One",
		Description: "This is a description for the first amazing feature.",
		Icon:        "icon-star",
	}).Add(&components.Card{
		Title:       "Feature Two",
		Description: "Discover the second feature that will change your life.",
		Icon:        "icon-heart",
	})

	// 5. Create the about page
	aboutPage := site.NewPage("About Us", "about.html")
	aboutPage.NewSection("About Our Company").Add(&components.Card{
		Title:       "Our Mission",
		Description: "To build the best and most awesome things.",
	})

	// 6. Create the contact page
	contactPage := site.NewPage("Contact Us", "contact.html")
	contactSection := contactPage.NewSection("Contact Us")
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
	if err := site.Generate(); err != nil {
		log.Fatalf("Error generating site: %v", err)
	}

	log.Println("Site generated successfully in 'dist' directory!")
}
