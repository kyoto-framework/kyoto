# Core

This part of documentation holds core library features.

## Page

Entrypoint for page rendering.  
Check [page interfaces](/concepts.html#page-interfaces) concept for more details.  

Example of usage:

```go
type PageIndex struct {
    Demo ssc.Component
}

func (p *PageIndex) Init() {
    p.DemoComponent = ssc.RegC(p, &ComponentDemo{})
}

func (p *PageIndex) Template() *template.Template {
    return template.Must(template.New("page.index.html").Funcs(ssc.Funcs()).ParseGlob("*.html"))
}
```

```html
<html>
    <head></head>
    <body>
        {{ template "ComponentDemo" .Demo }}
    </body>
</html>
```

## Component

Represents independent, reusable part of frontend functionality.  
Check [component interfaces](/concepts.html#component-interfaces) concept for more details.  

Example of usage:

```go
type ComponentDemoUUID struct {
    UUID string
}

func (c *ComponentDemoUUID) Async() error {
    resp, err := http.Get("http://httpbin.org/uuid")
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    data := map[string]string{}
    json.NewDecoder(resp.Body).Decode(&data)
    c.UUID = data["uuid"]
    return nil
}
```

```html
{{ define "ComponentDemoUUID" }}
<div class="bg-gray-100 rounded p-4">
    <div class="text-2xl font-semibold">UUID Generator Demo</div>
    <div class="mt-2 text-lg">{{ .UUID }}</div>
</div>
{{ end }}
```

## Render page

Low-level method, responsible for page rendering directly in `io.Writer`.  
Under the hood, executes full rendering [lifecycle](/concepts.html#lifecycle).
Rarely used directly by developers, library has more high-level wrappers with context setters, etc.
Accepts 2 parameters - `Page` pointer and `io.Writer`.  

Example of usage:  

```go
func IndexPageHandler(rw http.ResponseWriter, r *http.Request) {
    ssc.RenderPage(rw, &PageIndex{})
}
```

## Context

You can use
[`ssc.SetContext`](https://github.com/yuriizinets/ssceng/blob/master/context.go#L9),
[`ssc.GetContext`](https://github.com/yuriizinets/ssceng/blob/master/context.go#L20) and
[`ssc.DelContext`](https://github.com/yuriizinets/ssceng/blob/master/context.go#L26) for managing your context.
Context uses `Page` instance as namespace for correct concurrency handling on requests level (`Page` instance is creating for each new request).
Context can be used for passing additional state (f.e.`Request`, `gin.Context`) which can be accessed inside of lifecycle methods, like `Init` or `Async`.

Example of usage:

```go
func IndexPageHandler(rw http.ResponseWriter, r *http.Request) {
    p := &PageIndex{}
    ssc.SetContext(p, "internal:r", r)
    ssc.SetContext(p, "internal:rw", rw)
    ssc.RenderPage(rw, )
}
...
func (p *PageIndex) Init() {
    r := ssc.GetContext(p, "internal:r").(*http.Request)
    rw := ssc.GetContext(p, "internal:rw").(http.ResponseWriter)
    ...
}
...
```

## Handler factory

High-level method for creating own page handler with context setting.  
Takes 2 parameters: `Page` pointer and `map[string]interface{}` containing context items (`SetContext` will be called for each item in map)

Example of usage:

```go
func pagehandler(p ssc.Page) http.HandlerFunc {
    return func(rw http.ResponseWriter, r *http.Request) {
        ssc.PageHandlerFactory(p, map[string]interface{}{
            "internal:rw": rw,
            "internal:r":  r,
        })(rw, r)
    }
}

...

func main() {
    ...
    mux.HandleFunc("/", pagehandler(&PageIndex{}))
    ...
}
```

## Async components

In case when you need to fetch data from external API or database,
you can implement it as [`ImplementsAsync`](https://github.com/yuriizinets/ssceng/blob/master/types.go#L69) interface.
All needed async methods will will be called concurrently in [lifecycle](/concepts.html#lifecycle).

Example of usage: check [component](/core.html#component) section

## Flags

You can set execution flags as part of configuration.  
All flags are collected in [`flags.go`](https://github.com/yuriizinets/ssceng/blob/master/flags.go).

> This flag names are temporary and will be changed in nearest future

- `BENCH_LOWLEVEL` - responsible for logging lifecycle timings
- `BENCH_HANDLERS` - responsible for logging more highlevel timings
