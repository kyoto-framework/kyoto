# Quick Start

> Please note, currently we recommend to use demo project as a start. Documentation needs a lot of work.

## Simple

Just let's render a page  

```go
package main

import(
    "html/template"

    "github.com/gin-gonic/gin"
    "github.com/yuriizinets/ssceng"
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

## Advanced

So, where is a profit?
Let’s create more realistic example with data component and use it multiple times.
Component instances will fetch dynamic data concurrently with goroutines under the hood. Total page load time will be lower than using more traditional approach and you don’t have to think about concurrency by yourself.  

index.go

```go
package main

import(
    "html/template"

    "github.com/yuriizinets/ssceng"
)

type PageIndex struct{
    FirstInstance ssc.Component
    SecondInstance ssc.Component
    ThirdInstance ssc.Component
}

func (*PageIndex) Template() *template.Template {
    tmpl, _ := template.New("index.html").ParseGlob("*.html")
    return tmpl
}

func (p *PageIndex) Init() {
    p.FirstInstance = ssc.RegC(p, &ComponentLorem{})
    p.SecondInstance = ssc.RegC(p, &ComponentLorem{})
    p.ThirdInstance = ssc.RegC(p, &ComponentLorem{})
}
```

index.html

```html
<html>
    <body>
        {{ template "ComponentLorem" .FirstInstance }}
        {{ template "ComponentLorem" .SecondInstance }}
        {{ template "ComponentLorem" .ThirdInstance }}
    </body>
</html>
```

component.lorem.go

```go
package main

type ComponentLorem struct {
    Data template.HTML
}

func (c *ComponentLorem) Async() error {
    resp, _ := http.Get("https://loripsum.net/api")
    bcontent, err := ioutil.ReadAll(resp.Body)
    c.Data = template.HTML(bcontent)
    return err
}
```

component.lorem.html

```html
<div>
    {{ .Data }}
</div>
```
