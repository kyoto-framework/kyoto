
<p align="center">
    <img src="https://i.imgur.com/iHqOPQN.jpg">
</p>

An HTML render engine concept that brings frontend-like components experience to the server side with native `html/template` on steroids. Supports any serving basis (nethttp/Gin/etc), that provides `io.Writer` in response.  

## Disclaimer

> Under heavy development, not stable **(!!!)**
> **I'm not Golang "rockstar"**, and code may be not so good quality as you may expect. If you see any problems in the project - feel free to open new Issue.

## Quick Start

Just let's render a page

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

## More realistic example

So, where is a profit?  
Let's create more realistic example with data component an use it multiple times.  
Component instances will fetch dynamic data concurrently with goroutines under the hood. So, total page load time will be lower than using more traditional approach and you don't have to think about concurrency by yourself.

index.go

```go
package main

import(
    "html/template"

    "github.com/gin-gonic/gin"
    "github.com/yuriizinets/go-ssc"
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
<!-- index.html -->
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

## Features

- Component approach in mix with `html/template`
- Asynchronous operations
- Component methods that can be called from the client side (Server Side Actions, SSA)
- Different types of component communication (parent, cross, etc)

If you want to find more usage information, please, check `demo` folder. Usage documentation is not ready yet.

## Why?

I am trying to minimize the usage of popular SPA/PWA frameworks where it's not needed because it adds a lot of complexity and overhead. I don't want to bring significant runtime, VirtualDOM, and Webpack into the project with minimal dynamic frontend behavior.

This project proves the possibility of keeping most of the logic on the server's side.

## What problems does it solve? Why not using plain GoKit?

While developing the website's frontend with Go, I discovered some of the downsides of this approach:

- With plain html/template you're starting to repeat yourself. It's harder to define reusable parts.
- You must repeat DTO calls for each page, where you're using reusable parts.
- With Go's routines approach it's hard to make async-like DTO calls in the handlers.
- For dynamic things, you still need to use JS and client-side DOM modification.

Complexity is much higher when all of them get combined.

This engine tries to bring components and async experience to the traditional server-side rendering.

## Zen

- Don't replace Go features that exist already
- Don't do work that's already done
- Don't force developers to use a specific solution (Gin/Chi/GORM/sqlx/etc). Let them choose
- Rely on the server to do the rendering, minimum JS specifics or client-side only behavior

## Basic concepts

Each page or component is represented by its own structure. For implementing specific functionality, you can use structure's methods with a predefined declaration (f.e. `Init(p ssc.Page)`). You need to follow declaration rules to match the interfaces required (you can find all interfaces in `types.go`).  
Before implementing any method, you need to understand the rendering lifecycle.

### Lifecycle

Each page's lifecycle is hidden under the render function and follows this steps:

- Defining shared variables (waitgroup, errors channel)
- Triggering the page's `Init()` to initialize and register components
- Running all component's `Async()` functions in separate goroutines
- Waiting untill all asynchronous operations are completed
- Calling `AfterAsync()` for each component
- Cleaning up registered components (not needed more for internal usage)
- Getting page's template and render

### SSA Lifecycle

If you want to use SSA in your project, it's better to know how it works first. SSA has own, shorten lifecycle.  

- Creating request with JS on client side
- Extracting action data from request (component name, component state, action name, action args)
- Finding registered component type
- Creating component struct
- Triggering the component's `Init()`
- Populating component's state
- Calling action
- Rendering component and returning HTML to client side
- Replacing component's HTML with recieved version

## Sponsors

Broker One  
[https://brokerone.io](https://brokerone.io)  
[https://mybrokerone.com](https://mybrokerone.com)

![Broker One](https://i.imgur.com/jLBC4jV.jpg)