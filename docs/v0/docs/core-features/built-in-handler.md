
# Built-in handler


High-level function, returns `http.HandlerFunc` that can be used directly by `net/http` or a compatible framework.  
Takes 1 paramter - page pointer.  
Under the hood writes 2 context variables, that you can use with `GetContext`:

- `internal:rw` - `http.ResponseWriter`
- `internal:r` - `*http.Request`

Usage:

```go
...
func main() {
    ...
    mux.HandleFunc("/", kyoto.PageHandler(&PageIndex{}))
    ...
}
```
