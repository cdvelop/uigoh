package gosite

// SiteLink defines the interface for communication between components and the site.
type SiteLink interface {
	PageCount() int
	BuildNav() string
	AddCSS(css string)
	AddJS(js string)
}

// EventBinder adds or removes an event listener from a DOM element.
// The implementation will use syscall/js to interact with the DOM.
type EventBinder interface {
	EventListener(add bool, elementID, eventType string, callback func())
}

// HTMLRenderer is an interface for components that render HTML.
type HTMLRenderer interface {
	RenderHTML() string
}

// CSSRenderer is an interface for components that render CSS.
type CSSRenderer interface {
	RenderCSS() string
}

// JSRenderer is an interface for components that render JavaScript.
type JSRenderer interface {
	RenderJS() string
}
