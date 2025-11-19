//go:build !wasm
// +build !wasm

package postcard

import (
	_ "embed"
)

//go:embed style.css
var styleCss string

// RenderCSS returns the CSS for the postcard.
func (c *PostCard) RenderCSS() string {
	return styleCss
}
