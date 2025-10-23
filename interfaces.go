package gosite

// SiteLink defines the interface for communication between components and the site.
type SiteLink interface {
	PageCount() int
	BuildNav() string
	AddCSS(css string)
	AddJS(js string)
}
