
package uigoh

import (
	"github.com/cdvelop/uigoh/components"
)

// NewSite creates a new site and returns the root page.
func NewSite(cfg *components.Config) *components.Page {
	return components.NewPage(cfg)
}
