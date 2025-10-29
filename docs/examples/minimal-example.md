## 14. Appendix

### A. Complete Minimal Examples

**Shared Site Logic (NO build tags):**
```go
// app/site.go
package app

import "github.com/cdvelop/gosite"

func BuildSite(site *gosite.Site) {
    page := site.NewPage("Home", "index.html")
    section := page.NewSection("Welcome")
    section.Add(&MyComponent{Text: "Hello World"})
}

type MyComponent struct{ Text string }
func (c *MyComponent) RenderHTML() string {
    return "<p>" + c.Text + "</p>"
}
```

**Backend Entry (12 lines):**
```go
//go:build !wasm
// cmd/generator/main.go
package main
import "myproject/app"

func main() {
    site := gosite.New(&gosite.Config{
        OutputDir: "dist",
        WriteFile: os.WriteFile,
    })
    app.BuildSite(site)
    site.Generate()
}
```

**Frontend Entry (15 lines):**
```go
//go:build wasm
// cmd/webapp/main.go
package main
import "myproject/app"

func main() {
    site := gosite.New(&gosite.Config{
        EventBinder: NewDOMEventBinder(),
    })
    app.BuildSite(site)
    html := site.Pages[0].RenderHTML()
    js.Global().Get("document").
        Call("getElementById", "app").
        Set("innerHTML", html)
    <-make(chan struct{})
}
```

**Key Point:** The site structure (`BuildSite`) is written ONCE and used by both environments.
