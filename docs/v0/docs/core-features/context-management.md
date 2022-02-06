
# Context management

You can use `kyoto.SetContext`, `kyoto.GetContext` and `kyoto.DelContext` for managing your context.

Context uses page instance as namespace for correct concurrency handling on requests level (page instance is creating for each new request).
Context can be used for passing additional state (f.e.`http.Request`, `gin.Context`) which can be accessed inside of lifecycle methods, like `Init` or `Async`.  
**It's important to cleanup context with `kyoto.DelContext(p, "")` after page processing to avoid memory leaks! Built-in handler already doing it.**

Example of usage:

```go
func IndexPageHandler(rw http.ResponseWriter, r *http.Request) {
    p := &PageIndex{}
    kyoto.SetContext(p, "internal:r", r)
    kyoto.SetContext(p, "internal:rw", rw)
    kyoto.RenderPage(rw, )
    kyoto.DelContext(p, "")
}
...
func (p *PageIndex) Init() {
    r := kyoto.GetContext(p, "internal:r").(*http.Request)
    rw := kyoto.GetContext(p, "internal:rw").(http.ResponseWriter)
    ...
}
...
```

Most of the component methods have an overload option with a page argument. This way you don't need to store the page pointer in the component itself. Check full [interfaces specification](/concepts/#interfaces) in the [Concepts](/concepts) section.  
Example of overloaded asynchronous method:

```go
...

func (*ComponentExample) Async(p kyoto.Page) error {
    r := kyoto.GetContext(p, "internal:r").(*http.Request)
}

...
```