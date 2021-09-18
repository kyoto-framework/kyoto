# Extended

This part of documentation holds extended library features.

## Meta builder

::: warning
Not stable. In active development.
:::

To simplify work with additional head tags, library includes meta builder.  
Useful for setting title, description, canonical url, etc.
Check [`ssc.Meta`](https://github.com/yuriizinets/ssceng/blob/master/types.go#L31) for details.
For using meta builder, you need to implement page's [`ImplementsMeta`](https://github.com/yuriizinets/ssceng/blob/master/types.go#L81) interface.

Example of usage:

```go
...
func (p *PageIndex) Meta() ssc.Meta {
    return ssc.Meta{
        Title: "Example Page",
    }
}
```

## Server Side Actions (SSA)

::: warning
Not stable. In active development.
:::

Server Side Actions very similar to component methods in traditional frontend frameworks.
Main difference - all actions are executed on server side, code placed only on server and client has only thin communication layer.
Frontend only recieves ready for use HTML markup.  
Check [SSA lifecycle](/concepts/#ssa-rendering) for more details.

### Basic

First, you need to include thin JS layer for client-server communication

```html
...
<body>
    ...
    {{ dynamics }}
</body>
...
```

For implementing component actions, you need to implement [`ImplementsActions`](https://github.com/yuriizinets/ssceng/blob/master/types.go#L77) interface

```go
...
func (c *ComponentDemoUUID) Actions() ssc.ActionMap {
    return ssc.ActionMap{
        "Update": func(args ...interface{}) {
            resp, err := http.Get("http://httpbin.org/uuid")
            if err != nil {
                return
            }
            defer resp.Body.Close()
            data := map[string]string{}
            json.NewDecoder(resp.Body).Decode(&data)
            c.UUID = data["uuid"]
        },
    }
}
...
```

As one of requirements for Actions, you need to include additional attributes into your component with `componentattrs` template function.
This function injects internal data (like component's state).

```html
{{ define "ComponentDemoUUID" }}
<div {{ componentattrs . }}>
    ...
</div>
{{ end }}
```

After that, you can trigger an action with `action` template function.
First argument is the action name, all farther arguments are passed as `args ..interface{}` to action.
After execution, component's HTML will be updated with HTML markup from Action server response.

```html
...
<button onclick="{{ action `Update` }}">update</button>
...
```

### Cross-component calls

Documentation not ready yet

## Server Side State

::: danger
Not implemented yet. Check [issue](https://github.com/yuriizinets/ssceng/issues/28) state
:::

This feature is useful in case of large state payloads.
Instead of saving state inline as html tag, store state on server side and inject state hash as html tag.
Using this, you will decrease amount of data sent with SSA request and total HTML document size.
