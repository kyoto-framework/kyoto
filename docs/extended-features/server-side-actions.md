
# Server Side Actions

## Overview

Server Side Actions (SSA) very similar to component methods in traditional frontend frameworks.
Main difference - all actions are executed on server side, code placed only on server and client has only thin communication layer.
Frontend only recieves ready for use HTML markup.

## Installation

For using SSA, you need to define template builder, register SSA handler, and include communication layer on target page.

Here is how you can create template builder:

```go
func ssatemplate(p kyoto.Page) *template.Template {
    return template.Must(template.New("SSA").Funcs(kyoto.TFuncMap()).ParseGlob("*.html"))
}
```

After creation of template builder, you need to register SSA handler under `/SSA/` route.

```go
...

mux.HandleFunc("/SSA/", kyoto.SSAHandler(ssatemplate))

// In case of default http mux, use this
http.HandleFunc("/SSA/", kyoto.SSAHandler(ssatemplate))

...
```

And now, we need to include thin communication layer, implemented with JS, into target page.  
This can be done with `dynamics` template function, provided by `kyoto.Funcs()` function (check [Page rendering](/core-features/#page-rendering) for details)

```html
<html>
    <head>
        ...
    </head>
    <body>
        ...
        {{ dynamics }}
    </body>
</html>
```

## Usage

### Actions definition

Now you can implement `Actions` method to define own component's methods.  
This method must return `kyoto.ActionMap`, map which holds your methods. Each method accepts dynamic arguments amount with `...interface{}`.
In the method you can modify component's state, dynamicaly create and initialize another components, etc.

Usage:

```go
...

func (c *ComponentExample) Actions() kyoto.ActionMap {
    return kyoto.ActionMap{
        "ExampleAction": func(args ...interface{}) {
            // Do what you want here
        },
        "Submit": func(args ...interface{}) {
            // Do what you want here
        },
    }
}

...
```

In case when you need page instance, f.e. for getting context, this method have overload option with page argument

```go
func (c *ComponentExample) Actions(p kyoto.Page) kyoto.ActionMap {
    ...
}
```

### Attributes injection

You need to include component attributes into your top-level node with `componentattrs` template function. This function accepts component as argument.  
This includes different internal library data and component state.

Usage:

```html
{{ define "ComponentExample" }}
<div {{ componentattrs . }}>
    ...
</div>
{{ end }}
```

### Features

#### Actions

Library provides multiple ways of action triggering. One of them - `action` template function. This function accepts multiple arguments: first argument is always action name, all arguments after that will be passed as `...interface{}` to action arguments.

> Please note, that you can use `action` template function only in event handlers, like `onclick="..."`.

Usage:

```html
{{ define "ComponentExample" }}
<div {{ componentattrs . }}>
    <button onclick="{{ action 'ExampleAction' }}">Click Me</button>
</div>
{{ end }}
```

#### Form handling

`action` is not the only way to trigger an action. `formsubmit` allows to handle form submition. On trigger, it calls `Submit` action, defined in your `kyoto.ActionMap`.
Instead of passing form values as arguments, library unpacks that data directly into component by name attribute.

Usage:

```html
<form
    {{ componentattrs . }}
    action="/" 
    method="POST"
    onsubmit="{{ formsubmit }}"
>
    <input name="Email" value="{{ .Email }}" type="email" />
    <button type="submit">Submit</button>
</form>
```

#### Input binding

Not all operations needs to be done on server side. Some actions like inputs binding better to implement on client side to avoid delays and unnecessary server calls.
For input binding, library provides `bind` template function. This function accepts one argument - target component field name.

Usage:

```html
{{ define "ComponentExample" }}
<div {{ componentattrs . }}>
    <input value="{{ .InputData }}" oninput="{{ bind 'InputData' }}" />
    <button onclick="{{ action 'ExampleAction' }}">Click Me</button>
</div>
{{ end }}
```

#### On load trigger

Library provides a way to trigger an action on page load. May be useful for components lazy loading. This feature is implemented as `ssa:onload` HTML attribute and accepts just an action name.

Usage:  

```html
{{ define "ComponentExample" }}
<div {{ componentattrs . }} ssa:onload="Load">
    ...
</div>
{{ end }}
```

#### Intersection trigger

You can use `ssa:onintersect` HTML attribute to trigger an action on element intersection.
This functionality was built on top of the browser's built-in `IntersectionObserver`.  

Usage:

```html
{{ define "ComponentExample" }}
<div {{ componentattrs . }} ssa:onintersect="OnIntersect">
    ...
</div>
{{ end }}
```

#### Poll trigger

For components that must to be updated over the time (f.e. charts, stats, etc), library have polling mechanism, based on browser's built-in `setInternal`.  
You can use this trigger with `ssa:poll` and `ssa:poll.interval` HTML attributes.

Usage:

```html
{{ define "ComponentExample" }}
<div {{ componentattrs . }} ssa:poll="Poll" ssa:poll.interval=1000>
    ...
</div>
{{ end }}
```

#### Control display on action call

Because `kyoto` makes a roundtrip to the server every time an action is triggered on the page,
there are cases when the page may not react immediately to a user event (like a click).
That's why library provides a way to easily control display attribute on action call.
You can use `ssa:oncall.display` HTML attribute to control display during action call.
In the end of an action, layout will be restored.  

!!! note
    Don't forget to set default display for loading elements like spinners and loaders

Usage:

```html
{{ define "ComponentExample" }}
<div {{ componentattrs . }}>
    <div ssa:oncall.display="block" style="display: none">
        Loading ...
    </div>
    <button onclick="{{ action 'Load' }}">Load</button>
</div>
{{ end }}
```

### Control rendering mode

There are cases when `morphdom` may fail. Then library falls back to `replace` mode instead, which just replaces element's `outerHTML`.
To force library use replace mode, you can use `ssa:render.mode` HTML attribute.  

Usage:

```html
{{ define "ComponentExample" }}
<div {{ componentattrs . }} ssa:render.mode="replace">
    ...
</div>
{{ end }}
```

## Lifecycle

SSA has own lifecycle, which is a bit different in comparison with page rendering

- Creating request on client side with communication layer
- Extracting action data from request on server side
- Finding registered component type
- Creating component instance
- Triggering component's initialization method (if implemented)
- Populating component's state
- Calling action
- If new components where registed while action execution, do asynchronous operations for them (overall async process is the same as for page rendering)
- Rendering component and returning HTML to client side
- Morphing recieved HTML with component, or replacing in case of morph failure or explicit `ssa:render.mode="replace"` attribute

## Limitations

- Your component must to be JSON serializable. This is related to how Actions feature works under the hood. Component state is stored directly in DOM and can be used by client-size operations like binding.
- Avoid using of interface types in components with Actions. In that way your component will become incompatible with JSON serialization process.
- Avoid huge states. This will increase total page size and will slow down Actions operations. You can always avoid field JSON serialization with `json:"-"` if it's not needed to be stored. F.e. in case of list/table data, we're trying to avoid storing rows in state and using database directly instead.
