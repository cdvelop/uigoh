
package components

import "strings"

type NavItem struct {
    Label string
    Href  string
    Icon  string
}

type NavConfig struct {
    Items []NavItem
}

func (n *NavConfig) RenderHTML() string {
    var b strings.Builder

    b.WriteString("<nav class=\"main-nav\">\n")

    // Add a button for mobile view
    b.WriteString("  <button class=\"nav-toggle\" aria-label=\"Toggle navigation\">&#9776;</button>\n")

    b.WriteString("  <div class=\"nav-links\">\n")
    for _, item := range n.Items {
        b.WriteString("    <a href=\"")
        b.WriteString(escapeAttr(item.Href))
        b.WriteString("\" class=\"nav-link\">\n")

        if item.Icon != "" {
            b.WriteString("      <svg class=\"icon\"><use href=\"icons.svg#")
            b.WriteString(item.Icon)
            b.WriteString("\"></use></svg>\n")
        }

        b.WriteString("      <span>")
        b.WriteString(escapeHTML(item.Label))
        b.WriteString("</span>\n")
        b.WriteString("    </a>\n")
    }
    b.WriteString("  </div>\n")

    b.WriteString("</nav>\n")

    return b.String()
}

func (n *NavConfig) RenderCSS() string {
    return `.main-nav {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 1rem;
  padding: 1rem;
  background: var(--color-nav-bg);
  border-bottom: 1px solid var(--color-border);
}

.nav-links {
  display: flex;
  gap: 1rem;
}

.nav-link {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 1rem;
  text-decoration: none;
  color: var(--color-text);
  border-radius: 4px;
  transition: background 0.2s;
}

.nav-link:hover {
  background: var(--color-nav-hover);
}

.nav-link .icon {
  width: 20px;
  height: 20px;
}

.nav-toggle {
  display: none; // Hidden by default on larger screens
  background: none;
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
}

// Media query for mobile devices
@media (max-width: 768px) {
  .nav-links {
    display: none; // Hide links by default on mobile
    flex-direction: column;
    width: 100%;
    position: absolute;
    top: 60px; // Adjust based on nav height
    left: 0;
    background: var(--color-nav-bg);
  }

  .nav-links.active {
    display: flex; // Show links when active
  }

  .nav-toggle {
    display: block; // Show the toggle button on mobile
  }
}
`
}

func (n *NavConfig) RenderJS() string {
    return `// Navigation mobile toggle
document.addEventListener('DOMContentLoaded', function() {
    const navToggle = document.querySelector('.nav-toggle');
    const navLinks = document.querySelector('.nav-links');

    if (navToggle && navLinks) {
        navToggle.addEventListener('click', function() {
            navLinks.classList.toggle('active');
        });
    }
});
`
}
