# Server Side Actions

## Overview

Server Side Actions (SSA) are very similar to component methods in traditional front-end frameworks.
The main difference - all actions are executed on the server side, code placed on the server and client has just a thin communication layer.
The front-end then only receives ready-to-use HTML mark-up.

## Installation

To use SSA, you will need to define a template builder, register a SSA handler, and include a communication layer on target page.

Here is how you can create the template builder:

```go
func ssatemplate(p kyoto.Page) *template.Template {
    return template.Must(template.New("SSA").Funcs(kyoto.TFuncMap()).ParseGlob("*.html"))
}
```

After the creation of the template builder, you need to register an SSA handler under an `/SSA/` route.

```go
...

mux.HandleFunc("/SSA/", kyoto.SSAHandler(ssatemplate))

// In case of default http mux, use this
http.HandleFunc("/SSA/", kyoto.SSAHandler(ssatemplate))

...
```

Now we need to include a thin communication layer, implemented with JS, into target page.  
This can be done with a `dynamics` template function, provided by `kyoto.Funcs()` (check [Page Rendering](/core-features/#page-rendering) for details)

```html
<html>
  <head>
    ...
  </head>
  <body>
    ... {{ dynamics }}
  </body>
</html>
```

## Usage

### Actions Definition

Now you can implement an `Actions` method to define your own component's methods.  
This method must return `kyoto.ActionMap`, a map which holds your methods. Each method must accept dynamic arguments with `...interface{}`.
In these methods you can modify the component's state e.g. dynamically create and initialize others components, etc.

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

In the situation where you need the instance of the page, i.e for getting context, this method has an overloaded option with a page argument

```go
func (c *ComponentExample) Actions(p kyoto.Page) kyoto.ActionMap {
    ...
}
```

### Attributes Injection

You need to include component attributes into your top-level node with the `componentattrs` template function. This function accepts components as an argument.  
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

The library provides multiple ways of action triggering. One of them: the `action` template function. This function accepts multiple arguments: the first being the action name, all arguments after that will be passed as `...interface{}` as action arguments.

> Please note that you can use `action` template function only in event handlers, like `onclick="..."`.

Usage:

```html
{{ define "ComponentExample" }}
<div {{ componentattrs . }}>
  <button onclick="{{ action 'ExampleAction' }}">Click Me</button>
</div>
{{ end }}
```

#### Form Handling

`action` is not the only way to trigger an action. `formsubmit` allows handling form submission. Upon trigger it calls the `Submit` action, defined in your `kyoto.ActionMap`.
Instead of passing form values as arguments, library unpacks that data directly into the component by the name attribute.

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

#### Input Binding

Not all operations needs to be done on server side. Some actions like input binding are better implemented on the client side to avoid delays and unnecessary server calls.
For input binding, the Kyoto library provides the `bind` template function. This function accepts one argument - the target component field name.

Usage:

```html
{{ define "ComponentExample" }}
<div {{ componentattrs . }}>
  <input value="{{ .InputData }}" oninput="{{ bind 'InputData' }}" />
  <button onclick="{{ action 'ExampleAction' }}">Click Me</button>
</div>
{{ end }}
```

#### On Load Trigger

The library provides a way to trigger an action on page load. This may be useful for components' lazy loading. This feature is implemented as `ssa:onload` HTML attribute and accepts just an action name.

Usage:

```html
{{ define "ComponentExample" }}
<div {{ componentattrs . }} ssa:onload="Load">...</div>
{{ end }}
```

#### Intersection Trigger

You can use the `ssa:onintersect` HTML attribute to trigger an action on element intersection.
This functionality was built on top of the browser's built-in `IntersectionObserver`.

Usage:

```html
{{ define "ComponentExample" }}
<div {{ componentattrs . }} ssa:onintersect="OnIntersect">...</div>
{{ end }}
```

#### Poll Trigger

For components that must to be updated over time (f.e. charts, stats, etc), the library has a polling mechanism, based on browser's built-in `setInternal`.  
You can use this trigger with `ssa:poll` and `ssa:poll.interval` HTML attributes.

Usage:

```html
{{ define "ComponentExample" }}
<div {{ componentattrs . }} ssa:poll="Poll" ssa:poll.interval="1000">...</div>
{{ end }}
```

#### Control Display on Action Call

Because `kyoto` makes a roundtrip to the server every time an action is triggered on the page,
there are cases where the page may not react immediately to a user event (like a click).
That's why the library provides a way to easily control display attributes on action call.
You can use `ssa:oncall.display` HTML attribute to control display during action call.
At the end of an action the layout will be restored.

!!! note
    Don't forget to set a default display for loading elements like spinners and loaders.

Usage:

```html
{{ define "ComponentExample" }}
<div {{ componentattrs . }}>
  <div ssa:oncall.display="block" style="display: none">Loading ...</div>
  <button onclick="{{ action 'Load' }}">Load</button>
</div>
{{ end }}
```

### Control Rendering Mode

There are cases when `morphdom` may fail. The library then falls back to `replace` mode instead, which just replaces element's `outerHTML`.
To force the library to use replace mode, you can use the `ssa:render.mode` HTML attribute.

Usage:

```html
{{ define "ComponentExample" }}
<div {{ componentattrs . }} ssa:render.mode="replace">...</div>
{{ end }}
```

## Lifecycle

SSA has its own lifecycle, which is a bit different in comparison with page rendering

- Creating a request on the client-side with the communication layer
- Extracting action data from requests on the server-side
- Finding registered component types
- Creating component instances
- Triggering component's initialization method (if implemented)
- Populating component's state
- Calling actions
- If any new components were registered during the action execution, do asynchronous operations for them (overall async process is the same as for page rendering)
- Rendering components and returning HTML to client-side
- Morphing received HTML with component, or replacing in case of morph failure or explicit `ssa:render.mode="replace"` attribute

## Limitations and Advice

- Your component must to be JSON serializable. This is related to how the Actions feature works under the hood. Component state is stored directly in DOM and can be used by client-size operations like binding.
- Avoid using of interface types in components with Actions. That way your component will become incompatible with JSON serialization process.
- Avoid huge states. This will increase the total page size and will slow down Actions operations. You can avoid field JSON serialization with `json:"-"` if it's not needed. i.e in case of list/table data. We're trying to avoid storing rows of data in the state and using database directly instead.
