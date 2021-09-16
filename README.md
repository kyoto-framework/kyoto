
# SSC Engine

An HTML render engine concept that brings frontend-like components experience to the server side with native `html/template` on steroids. Supports any serving basis (`net/http`/`gin`/etc), that provides io.Writer in response.

> **Disclaimer**  
> This project is an experimental concept. Code quality may not meet your expectations. **Don’t use in production.** In case of any issues/proposals, feel free to open an issue

## Quick Start

Let's just render a page

```go
package main

import(
    "html/template"

    "github.com/gin-gonic/gin"
    "github.com/yuriizinets/go-ssc"
)

// PageIndex is an implementation of ssc.Page interface
type PageIndex struct{}

// Template is a required page method. It tells about template configuration
func (*PageIndex) Template() *template.Template {
    // Template body is located in index.html
    // <html>
    //   <body>The most basic example</body>
    // </html>
    tmpl, _ := template.New("index.html").ParseGlob("*.html")
    return tmpl
}

func main() {
    g := gin.Default()

    g.GET("/", func(c *gin.Context) {
        ssc.RenderPage(c.Writer, &PageIndex{})
    })

    g.Run("localhost:25025")
}
```

## Documentation & overview

For more details, check project's page
