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
	//*css
	return `/* Main nav styles */
.main-nav {
	height: 60px;
	background: linear-gradient(135deg, var(--color-primary), #2c6aa0);
	box-shadow: 0 2px 10px rgba(0,0,0,0.1);
	display: flex;
	justify-content: flex-end;
	align-items: center;
	position: sticky;
	top: 0;
	z-index: 100;
}

/* Links container */
.links-container {
	height: 100%;
	width: 100%;
	display: flex;
	flex-direction: row;
	align-items: center;
}

/* Nav links */
.main-nav a {
	height: 100%;
	padding: 0 20px;
	display: flex;
	align-items: center;
	text-decoration: none;
	color: white;
	transition: background 0.2s;
	font-weight: 500;
}

.main-nav a:hover {
	background: rgba(255,255,255,0.2);
}

.main-nav .home-link {
	margin-right: auto;
	font-weight: 600;
}

/* SVG styles */
.main-nav svg {
	fill: white;
}

/* Hide checkbox and buttons by default */
#sidebar-active {
	display: none;
}

.open-sidebar-button,
.close-sidebar-button {
	display: none;
}

/* Mobile responsive styles */
@media (max-width: 768px) {
	.links-container {
		flex-direction: column;
		align-items: flex-start;
		position: fixed;
		top: 0;
		right: -100%;
		z-index: 10;
		width: 300px;
		height: 100vh;
		background: linear-gradient(180deg, var(--color-primary), #2c6aa0);
		box-shadow: -5px 0 15px rgba(0, 0, 0, 0.3);
		transition: right 0.3s ease-out;
	}

	.main-nav a {
		box-sizing: border-box;
		height: auto;
		width: 100%;
		padding: 20px 30px;
		justify-content: flex-start;
		border-bottom: 1px solid rgba(255,255,255,0.1);
	}

	.main-nav .home-link {
		margin-right: 0;
	}

	.open-sidebar-button,
	.close-sidebar-button {
		padding: 20px;
		display: block;
		cursor: pointer;
	}

	#sidebar-active:checked ~ .links-container {
		right: 0;
	}

	#sidebar-active:checked ~ #overlay {
		height: 100%;
		width: 100%;
		position: fixed;
		top: 0;
		left: 0;
		z-index: 9;
		background: rgba(0,0,0,0.5);
	}
}

/* View Transition API */
@view-transition {
	navigation: auto;
}

::view-transition-group(*) {
	animation-duration: 0.3s;
}

/* Smooth fade transition for page content */
@keyframes fade-in {
	from {
		opacity: 0;
	}
}

@keyframes fade-out {
	to {
		opacity: 0;
	}
}

::view-transition-old(root) {
	animation: 150ms cubic-bezier(0.4, 0, 1, 1) both fade-out;
}

::view-transition-new(root) {
	animation: 300ms cubic-bezier(0, 0, 0.2, 1) both fade-in;
}
`
}

// RenderJS generates the JavaScript for view transitions
func (n *NavbarBuilder) RenderJS() string {
	return `// View Transition API for smooth page navigation
(function() {
	// Check if View Transition API is supported
	if (!document.startViewTransition) {
		return; // Fallback to normal navigation
	}

	// Intercept navigation link clicks
	document.addEventListener('click', function(e) {
		const link = e.target.closest('a');

		// Only handle internal navigation links
		if (!link || !link.href || link.target === '_blank') return;

		const url = new URL(link.href);

		// Check if it's a same-origin navigation
		if (url.origin !== location.origin) return;

		// Prevent default navigation
		e.preventDefault();

		// Start view transition
		document.startViewTransition(async () => {
			// Fetch the new page
			const response = await fetch(url.href);
			const html = await response.text();

			// Parse the HTML
			const parser = new DOMParser();
			const doc = parser.parseFromString(html, 'text/html');

			// Update the document title
			document.title = doc.title;

			// Replace the main content
			const newMain = doc.querySelector('main');
			const currentMain = document.querySelector('main');

			if (newMain && currentMain) {
				currentMain.replaceWith(newMain);
			} else {
				// Fallback: replace entire body content
				document.body.innerHTML = doc.body.innerHTML;
			}

			// Update the URL
			history.pushState({}, '', url.href);

			// Re-attach event listeners after DOM update
			initializeEventListeners();
		});
	});

	// Handle browser back/forward buttons
	window.addEventListener('popstate', function() {
		document.startViewTransition(async () => {
			const response = await fetch(location.href);
			const html = await response.text();
			const parser = new DOMParser();
			const doc = parser.parseFromString(html, 'text/html');

			document.title = doc.title;

			const newMain = doc.querySelector('main');
			const currentMain = document.querySelector('main');

			if (newMain && currentMain) {
				currentMain.replaceWith(newMain);
			} else {
				document.body.innerHTML = doc.body.innerHTML;
			}

			initializeEventListeners();
		});
	});

	// Function to reinitialize event listeners after DOM updates
	function initializeEventListeners() {
		// Re-attach any component-specific event listeners here
		// This ensures interactive elements work after page transitions
	}
})();
`
}
