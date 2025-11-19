//go:build !wasm
// +build !wasm

package sectionhead

import (
	_ "embed"
)

//go:embed style.css
var styleCss string

// RenderCSS returns the CSS for the sectionhead.
func (c *SectionHead) RenderCSS() string {
	return styleCss
}
