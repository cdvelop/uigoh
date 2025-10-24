package gosite

import (
	. "github.com/cdvelop/tinystring"
)

// NavbarBuilder handles the construction of the navigation bar
type NavbarBuilder struct {
	site *Site
}

// Render generates the navbar HTML with mobile-responsive structure
func (n *NavbarBuilder) Render() string {
	var b = Convert()

	b.Write("<nav class=\"main-nav\">\n")
	b.Write("  <input type=\"checkbox\" id=\"sidebar-active\">\n")
	b.Write("  <label for=\"sidebar-active\" class=\"open-sidebar-button\">\n")
	b.Write("    <svg xmlns=\"http://www.w3.org/2000/svg\" height=\"32\" viewBox=\"0 -960 960 960\" width=\"32\">\n")
	b.Write("      <path d=\"M120-240v-80h720v80H120Zm0-200v-80h720v80H120Zm0-200v-80h720v80H120Z\"/>\n")
	b.Write("    </svg>\n")
	b.Write("  </label>\n")
	b.Write("  <label id=\"overlay\" for=\"sidebar-active\"></label>\n")
	b.Write("  <div class=\"links-container\">\n")
	b.Write("    <label for=\"sidebar-active\" class=\"close-sidebar-button\">\n")
	b.Write("      <svg xmlns=\"http://www.w3.org/2000/svg\" height=\"32\" viewBox=\"0 -960 960 960\" width=\"32\">\n")
	b.Write("        <path d=\"m256-200-56-56 224-224-224-224 56-56 224 224 224-224 56 56-224 224 224 224-56 56-224-224-224 224Z\"/>\n")
	b.Write("      </svg>\n")
	b.Write("    </label>\n")

	// Add all page links
	for i, page := range n.site.pages {
		if i == 0 {
			b.Write("    <a class=\"home-link\" href=\"")
		} else {
			b.Write("    <a href=\"")
		}
		b.Write(Convert(page.filename).EscapeAttr())
		b.Write("\">")
		b.Write(Convert(page.title).EscapeHTML())
		b.Write("</a>\n")
	}

	b.Write("  </div>\n")
	b.Write("</nav>\n")

	return b.String()
}

// RenderCSS generates the navbar CSS with responsive styles
func (n *NavbarBuilder) RenderCSS() string {
	var b = Convert()

	// Main nav styles
	b.Write(".main-nav {\n")
	b.Write("  height: 60px;\n")
	b.Write("  background: linear-gradient(135deg, var(--color-primary), #2c6aa0);\n")
	b.Write("  box-shadow: 0 2px 10px rgba(0,0,0,0.1);\n")
	b.Write("  display: flex;\n")
	b.Write("  justify-content: flex-end;\n")
	b.Write("  align-items: center;\n")
	b.Write("  position: sticky;\n")
	b.Write("  top: 0;\n")
	b.Write("  z-index: 100;\n")
	b.Write("}\n\n")

	// Links container
	b.Write(".links-container {\n")
	b.Write("  height: 100%;\n")
	b.Write("  width: 100%;\n")
	b.Write("  display: flex;\n")
	b.Write("  flex-direction: row;\n")
	b.Write("  align-items: center;\n")
	b.Write("}\n\n")

	// Nav links
	b.Write(".main-nav a {\n")
	b.Write("  height: 100%;\n")
	b.Write("  padding: 0 20px;\n")
	b.Write("  display: flex;\n")
	b.Write("  align-items: center;\n")
	b.Write("  text-decoration: none;\n")
	b.Write("  color: white;\n")
	b.Write("  transition: background 0.2s;\n")
	b.Write("  font-weight: 500;\n")
	b.Write("}\n\n")

	b.Write(".main-nav a:hover {\n")
	b.Write("  background: rgba(255,255,255,0.2);\n")
	b.Write("}\n\n")

	b.Write(".main-nav .home-link {\n")
	b.Write("  margin-right: auto;\n")
	b.Write("  font-weight: 600;\n")
	b.Write("}\n\n")

	// SVG styles
	b.Write(".main-nav svg {\n")
	b.Write("  fill: white;\n")
	b.Write("}\n\n")

	// Hide checkbox and buttons by default
	b.Write("#sidebar-active {\n")
	b.Write("  display: none;\n")
	b.Write("}\n\n")

	b.Write(".open-sidebar-button,\n")
	b.Write(".close-sidebar-button {\n")
	b.Write("  display: none;\n")
	b.Write("}\n\n")

	// Mobile responsive styles
	b.Write("@media (max-width: 768px) {\n")
	b.Write("  .links-container {\n")
	b.Write("    flex-direction: column;\n")
	b.Write("    align-items: flex-start;\n")
	b.Write("    position: fixed;\n")
	b.Write("    top: 0;\n")
	b.Write("    right: -100%;\n")
	b.Write("    z-index: 10;\n")
	b.Write("    width: 300px;\n")
	b.Write("    height: 100vh;\n")
	b.Write("    background: linear-gradient(180deg, var(--color-primary), #2c6aa0);\n")
	b.Write("    box-shadow: -5px 0 15px rgba(0, 0, 0, 0.3);\n")
	b.Write("    transition: right 0.3s ease-out;\n")
	b.Write("  }\n\n")

	b.Write("  .main-nav a {\n")
	b.Write("    box-sizing: border-box;\n")
	b.Write("    height: auto;\n")
	b.Write("    width: 100%;\n")
	b.Write("    padding: 20px 30px;\n")
	b.Write("    justify-content: flex-start;\n")
	b.Write("    border-bottom: 1px solid rgba(255,255,255,0.1);\n")
	b.Write("  }\n\n")

	b.Write("  .main-nav .home-link {\n")
	b.Write("    margin-right: 0;\n")
	b.Write("  }\n\n")

	b.Write("  .open-sidebar-button,\n")
	b.Write("  .close-sidebar-button {\n")
	b.Write("    padding: 20px;\n")
	b.Write("    display: block;\n")
	b.Write("    cursor: pointer;\n")
	b.Write("  }\n\n")

	b.Write("  #sidebar-active:checked ~ .links-container {\n")
	b.Write("    right: 0;\n")
	b.Write("  }\n\n")

	b.Write("  #sidebar-active:checked ~ #overlay {\n")
	b.Write("    height: 100%;\n")
	b.Write("    width: 100%;\n")
	b.Write("    position: fixed;\n")
	b.Write("    top: 0;\n")
	b.Write("    left: 0;\n")
	b.Write("    z-index: 9;\n")
	b.Write("    background: rgba(0,0,0,0.5);\n")
	b.Write("  }\n")
	b.Write("}\n\n")

	return b.String()
}
