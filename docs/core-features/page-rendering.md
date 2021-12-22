
# Page rendering

The most important part of the library.  
A low-level function, responsible for rendering the page directly into `io.Writer`.  
Under the hood, executes full rendering lifecycle. Library has more high-level wrappers with context setters and another features, but all of them rely on this function. Accepts 2 parameters - page pointer and `io.Writer`.

## Usage

First of all, let's create a page structure

```go title="page.index.go"
package main

import (
    "github.com/kyoto-framework/kyoto"
)

type PageIndex struct {}
```

As a requirement, each page must have an html template builder method.  
Please note that providing `kyoto.Funcs()` is not required, but highly recommended as far as it provides some library features.

```go title="page.index.go"
...

func (*PageIndex) Template() *template.Template {
    return template.Must(template.New("page.index.html").Funcs(kyoto.Funcs()).ParseGlob("*.html"))
}
```

After creating the page structure, it's time to create template

```html title="page.index.html"
<html>
  <head>
    <title>kyoto page</title>
  </head>
  <body>
    Hello World!
  </body>
</html>
```

Now you can use the rendering function

```go
func ExampleHandler(rw http.ResponseWriter, r *http.Request) {
    RenderPage(rw, &PageIndex{})
}
```

For example:

```go title="main.go"
...

func main() {
    http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
        kyoto.RenderPage(rw, &PageIndex{})
    })

    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

??? note "Complete example with component"
    Let's add an example component to make this example more complete. This component will generate random number and will hold that number as content.
    Please note that component template definition must to much with actual component structure name.

    ```go title="component.rand.go"
    package main

    import (
        "math/rand"
        "strconv"
    )

    type ComponentRand struct {
        Content string
    }

    func (c *ComponentRand) Init() {
        c.Content = strconv.Itoa(rand.Intn(1000))
    }
    ```

    ```html title="component.rand.html"
    {{ define "ComponentRand" }}
    <div>Random number: {{ .Content }}</div>
    {{ end }}
    ```

    After component creation, let's register and include it into page.  
    Check `lifecycle integration` section for detailed documentation.

    ```go title="page.index.go"
    package main

    import (
        "html/template"
        "github.com/kyoto-framework/kyoto"
    )

    type PageIndex struct {
        Rand kyoto.Component
    }

    func (*PageIndex) Template() *template.Template {
        return template.Must(template.New("page.index.html").Funcs(kyoto.TFuncMap()).ParseGlob("*.html"))
    }

    func (p *PageIndex) Init() {
        p.Rand = kyoto.RegC(p, &ComponentRand{})
    }
    ```

    ```html title="page.index.html"
    <html>
    <head>
        <title>kyoto page</title>
    </head>
    <body>
        {{ template "ComponentRand" .Rand }}
    </body>
    </html>
    ```

    That's it! Now you have component instance, included into lifecycle and rendered on the page.

## Alternative rendering

Library gives an option to use own rendering function.
It's enough to implement `ImplementsRender` interface and use `render` template function in case of components.  

Let's assume we want to create a component without relying on the `html/template`.
As a first step, let's create a component structure with a state. On initialization step we will generate a random number.  

```go title="component.rand.go"
package main

type ComponentRand struct {
    Content string
}

func (c *ComponentRand) Init() {
    c.Content = strconv.Itoa(rand.Intn(1000))
}
```

Second step is to implement `ImplementsRender` interface.  

```go title="component.rand.go"
...

func (c *ComponentRand) Render() string {
    return fmt.Sprintf(`
        <div>Random number: %s</div>
    `, c.Content)
}
```

To use the component, just include it into page with `render` function. Please note that we don't need to provide any template definitions now.  

```html title="page.index.html"
<html>
  <head>
    <title>kyoto page</title>
  </head>
  <body>
    {{ render .ComponentRand }}
  </body>
```

Also, you can implement this interface in page structure in the same way.

??? note "Custom template engine"
    Right! With this approach, you can use any template engine.  
    All template functions are available with `T` preffix, like `kyoto.TRender`.

??? note "Components library"
    This functionality is very suitable for components library.  
    In this way you can avoid templates packaging problem.
