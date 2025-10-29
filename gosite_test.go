package gosite_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cdvelop/gosite"
	"github.com/cdvelop/gosite/components/card"
	"github.com/cdvelop/gosite/components/form"
)

func TestGenerateExample(t *testing.T) {
	// Persisted output directory so generated files remain for manual inspection
	outDir := filepath.Join(".", "output")
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		t.Fatalf("creating outDir: %v", err)
	}

	cfg := &gosite.Config{
		Title:     "My Awesome Site",
		OutputDir: "output",
		ColorScheme: &gosite.ColorScheme{
			Primary:    "#3f88bf",
			Secondary:  "#ff9300",
			Text:       "#333333",
			Background: "#ffffff",
			Border:     "#e0e0e0",
		},
		WriteFile: func(path, content string) error {
			return os.WriteFile(path, []byte(content), 0o644)
		},
	}

	site := gosite.New(cfg)

	// Home page
	homePage := site.NewPage("Home Page", "index.html")
	homePage.AddHead(`<meta name="description" content="Welcome to the awesome site!">`)
	introSection := homePage.NewSection("Welcome to Our Website")
	introSection.Add(&card.Card{
		Title:       "Feature One",
		Description: "This is a description for the first amazing feature.",
		Icon:        "icon-star",
	}).Add(&card.Card{
		Title:       "Feature Two",
		Description: "Discover the second feature that will change your life.",
		Icon:        "icon-heart",
	})

	// About page
	aboutPage := site.NewPage("About Us", "about.html")
	aboutPage.NewSection("About Our Company").Add(&card.Card{
		Title:       "Our Mission",
		Description: "To build the best and most awesome things.",
	})

	// Contact page with form
	contactPage := site.NewPage("Contact Us", "contact.html")
	contactSection := contactPage.NewSection("Contact Us")
	contactSection.Add(&form.Form{
		Config: form.FormConfig{
			Action: "/submit-form",
			Method: "POST",
			Fields: []form.FormField{
				{Type: "text", Name: "username", Placeholder: "Your Name", Required: true},
				{Type: "email", Name: "email", Placeholder: "Your Email", Required: true},
				{Type: "textarea", Name: "message", Placeholder: "Your Message"},
			},
		},
	})

	if err := site.Generate(); err != nil {
		t.Fatalf("site.Generate() returned error: %v", err)
	}

	// Verify HTML files were written
	var htmlFiles []string
	err := filepath.WalkDir(outDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if strings.HasSuffix(d.Name(), ".html") {
			htmlFiles = append(htmlFiles, path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("error walking out dir: %v", err)
	}
	if len(htmlFiles) == 0 {
		t.Fatalf("no .html files were generated in out dir %s", outDir)
	}

	// Verify CSS file was generated with base styles
	cssPath := filepath.Join(outDir, "style.css")
	cssContent, err := os.ReadFile(cssPath)
	if err != nil {
		t.Fatalf("CSS file not generated: %v", err)
	}

	cssStr := string(cssContent)
	// Check for CSS variables
	if !strings.Contains(cssStr, "--color-primary") {
		t.Error("CSS missing --color-primary variable")
	}
	if !strings.Contains(cssStr, "#3f88bf") {
		t.Error("CSS missing custom primary color")
	}
	// Check for base styles
	if !strings.Contains(cssStr, ".main-nav") {
		t.Error("CSS missing .main-nav styles")
	}
	if !strings.Contains(cssStr, ".card") {
		t.Error("CSS missing .card component styles")
	}

	// Ensure at least one HTML contains expected text
	found := false
	for _, p := range htmlFiles {
		b, err := os.ReadFile(p)
		if err != nil {
			t.Fatalf("reading generated file %s: %v", p, err)
		}
		s := string(b)
		if strings.Contains(s, "My Awesome Site") || strings.Contains(s, "Home Page") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("no generated HTML files contain expected site title or page title")
	}

	t.Logf("✓ Generated site in: %s", outDir)
	t.Logf("✓ HTML files: %d", len(htmlFiles))
	t.Logf("✓ CSS file size: %d bytes", len(cssContent))
}
