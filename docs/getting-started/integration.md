# Integration

It's easy to integrate `kyoto` into existing projects. There are 2 ways to integrate: define your own generic page handler or use a low-level page rendering method. In most cases, the second method is the easiest (especially when no context is needed).

Example with `net/http`:

```go
func IndexPageHandler(rw http.ResponseWriter, r *http.Request) {
    RenderPage(rw, &IndexPage{})
}
```

Example with `gin`:

```go
func IndexPageHandler(g *gin.Context) {
    kyoto.RenderPage(g.Writer, &IndexPage{})
}
```

Check the [Page Rendering](#) section for details.
