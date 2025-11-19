//go:build !wasm
// +build !wasm

package navbar

import (
_ "embed"
)

//go:embed style.css
var styleCss string

// RenderCSS returns the CSS for the navbar.
func (n *Navbar) RenderCSS() string {
return styleCss
}

//go:embed script.js
var scriptJs string

// RenderJS returns the JavaScript for the navbar.
func (n *Navbar) RenderJS() string {
return scriptJs
}
