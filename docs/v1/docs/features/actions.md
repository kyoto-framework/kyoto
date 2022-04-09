
# Actions

Kyoto has a way to simplify building dynamic UIs.
For this purpose it has a module named "actions".
Principle is very similar to component methods in traditional front-end frameworks.
The main difference - all actions are executed on the server side, code placed on the server and client has just a thin communication layer.
Even template rendering remains on the server side to avoid bringing render functionality to the client.
Module uses approach similar to Laravel Livewire or Hotwired which sending HTML instead of JSON over the wire.

## Installation

To use actions, you will need to prepare your project:

- include a communication layer on resulting page with `dynamics` function
- register an actions handler (provided by `actions.Handler`) with a specific route
- register every dynamic component with `actions.Register`

=== "main.go"

	```go
	...

	// Register actions handler
	mux.HandleFunc("/internal/actions/", actions.Handler(func(c *kyoto.Core) *template.Template {
		return template.Must(template.New("Actions").Funcs(render.FuncMap(c)).ParseGlob("*.html"))
	}))
	// Register dynamic components
	actions.Register(
		ComponentFoo,
		ComponentBar,
		ComponentBaz(""),  // In case of wrapped core receiver, we need to call a wrapper
		...
	)
	
	...
	```

=== "page.index.html"

	```html
	<html>
		<head>
		...
		</head>
		<body>
		...
		<!-- Include actions communication layer -->
		{{ dynamics }}
		</body>
	</html>
	```

## Usage

### Definition

To define an action for a component you can use `actions.Define` adapter.
Also, you will need to include dynamic component attributes with `componentattrs` function into top-level HTML tag.

=== "component.foo.go"

	```go
	...

	func ComponentFoo(core *kyoto.Core) {
		...
		actions.Define(core, "Bar", func(args ...interface{}) {
			// Your action logic
		})
		actions.Define(core, "Baz", func(args ...interface{}) {
			// Your action logic
		})
	}

	...
	```

=== "component.foo.html"

	```html
	{{ define "ComponentFoo" }}
	<div {{ componentattrs . }} >
		...
	</div>
	{{ end }}
	```

### Action call

Library provides many ways to make an action call.
Most of them have shortcurts, implemented as template functions to simplify usage.

#### Direct

The most basic component action call.
There are 2 ways to trigger a component action: with JS function `Action(this, ...)` or template function `{{ action ... }}`.
In both cases additional arguments will be passed to action handler as `...interface{}`.

=== "component.foo.html"

	```html
	{{ define "ComponentFoo" }}
	<div {{ componentattrs . }}>
		<div> Current state: {{ .Content }} </div>
		<!-- Example of call with shortcut -->
		<button onclick="{{ action `Bar` }}">Bar</button>
		<!-- Example of call with JS function -->
		<button onclick="Action(this, 'Baz')">Baz</button>
	</div>
	{{ end }}
	```

=== "component.foo.go"

	```go
	...

	func ComponentFoo(core *kyoto.Core) {
		lifecycle.Init(core, func() {
			core.State.Set("Content", "Foo")
		})
		actions.Define(core, "Bar", func(args ...interface{}) {
			core.State.Set("Content", "Bar")
		})
		actions.Define(core, "Baz", func(args ...interface{}) {
			core.State.Set("Content", "Baz")
		})
	}

	...
	```


#### Cross-component

You can call other components' methods with different preffixes.
`$` preffix allows you to call methods of a parent component.
To call a components' method by it's id, you can use `:` delimiter with such syntax `<component-id>:<method-name>`.

=== "component.foo.html"

	```html
	{{ define "ComponentFoo" }}
	<div {{ componentattrs . }}>
		<button onclick="{{ action `$Trigger` }}">Trigger parent components' ("Bar") action named "Trigger"</button>
		<button onclick="{{ action `Baz:Trigger` }}">Trigger action named "Trigger" on component with id "Baz"</button>
	</div>
	{{ end }}
	```

=== "component.foo.go"

	```go
	...
	
	func ComponentFoo(core *kyoto.Core) {
		// Nothing to do here
	}

	...
	```

=== "component.bar.go"

	```go
	...

	func ComponentBar(core *kyoto.Core) {
		lifecycle.Init(core, func() {
			core.Component("Foo", ComponentFoo)
		})
		actions.Define(core, "Trigger", func(args ...interface{}) {
			...
		})
	}

	...
	```

=== "component.baz.go"

	```go
	...

	func ComponentBaz(core *kyoto.Core) {
		actions.Define(core, "Trigger", func(args ...interface{}) {
			...
		})
	}

	...
	```

#### Form submit

Thanks to actions, kyoto have a way to submit a form without page reloading.
Form submission will be received in the component as an action.
Instead of passing form values as arguments, library unpacks that data directly into the component by the name attribute.
There are 2 ways to use this feature: with JS function `FormSubmit(this, event)` or template function `{{ formsubmit }}`.

=== "component.form.html"

	```html
	{{ define "ComponentForm" }}
	<form
	  {{ componentattrs . }}
	  action="#"
	  method="POST"
	  onsubmit="{{ formsubmit }}"
	>
	  <input name="Email" value="{{ .Email }}" type="email" />
	  <button type="submit">Submit</button>
	</form>
	{{ end }}
	```

=== "component.form.go" 

	```go
	...

	func ComponentForm(core *kyoto.Core) {
		actions.Define(core, "Submit", func() {
			...
		})
	}

	...
	```


#### Trigger with onload

Library provides alternative ways to trigger an action.
One of them is triggering on page load.
This may be useful for components' lazy loading.
This feature is implemented as `ssa:onload` HTML attribute and accepts just an action name.

=== "component.foo.html"

	```html
	{{ define "ComponentFoo" }}
	<div {{ componentattrs . }} ssa:onload="Bar">
		...
	</div>
	{{ end }}
	```

=== "component.foo.go"

	```go
	...

	func ComponentFoo(core *kyoto.Core) {
		...
		actions.Define(core, "Bar", func(args ...interface{}) {
			...
		})
	}
	
	...
	```

#### Trigger with poll

Another way to trigger an action is triggering with interval.
Useful for components that must to be updated over time (f.e. charts, stats, etc).
You can use this trigger with `ssa:poll` and `ssa:poll.interval` HTML attributes.

=== "component.foo.html"

	```html
	{{ define "ComponentFoo" }}
	<div {{ componentattrs . }} ssa:poll="Bar" ssa:poll.interval="1000">
		...
	</div>
	{{ end }}
	```

=== "component.foo.go"

	```go
	...

	func ComponentFoo(core *kyoto.Core) {
		...
		actions.Define(core, "Bar", func(args ...interface{}) {
			...
		})
	}

	...
	```

#### Trigger with intersection

You can use the `ssa:onintersect` HTML attribute to trigger an action on element intersection.
This functionality was built on top of the browser's built-in `IntersectionObserver`.

=== "component.foo.html"

	```html
	{{ define "ComponentFoo" }}
	<div {{ componentattrs . }} ssa:onintersect="Bar">...</div>
	{{ end }}
	```

=== "component.foo.go"

	```go
	...
	
	func ComponentFoo(core *kyoto.Core) {
		...
		actions.Define(core, "Bar", func(args ...interface{}) {
			...
		})
	}

	...
	```

### Binding

Not all operations needs to be done on server side.
Some actions like input binding are better implemented on the client side to avoid delays and unnecessary server calls.
That's why the library have a way to bind controls to component state.

#### Input binding

For input binding, the Kyoto library provides the `bind` template function.
This function accepts one argument - the target component field name.
Also, you can use `Bind(this, 'FieldName')` JS function if you would like to avoid shortcuts.

=== "component.foo.html"

	```html
	{{ define "ComponentFoo" }}
	<div {{ componentattrs . }}>
		<input oninput="{{ bind `FieldName` }}">
		...
	</div>
	{{ end }}
	```

=== "component.foo.go"

	```go
	...
		
	func ComponentFoo(core *kyoto.Core) {
		lifecycle.Init(core, func() {
			core.State.Set("FieldName", "")
		})
	}

	...
	```

### Flow control

#### Control display on action

Because kyoto makes a roundtrip to the server every time an action is triggered on the page,
there are cases where the page may not react immediately to a user event (like a click).
That's why the library provides a way to easily control display attributes on action call.
You can use `ssa:oncall.display` HTML attribute to control display during action call.
At the end of an action the layout will be restored.

!!! note ""
	Don't forget to set a default display for loading elements like spinners and loaders.

=== "component.foo.html"

	```html
	{{ define "ComponentFoo" }}
	<div {{ componentattrs . }}>
	  <div ssa:oncall.display="block" style="display: none">Loading ...</div>
	  <button onclick="{{ action 'Bar' }}">Load</button>
	</div>
	{{ end }}
	```

=== "component.foo.go"

	```go
	...

	func ComponentFoo(core *kyoto.Core) {
		...
		actions.Define(core, "Bar", func(args ...interface{}) {
			...
		})
	}

	...
	```

#### Multi-stage UI update

You can push multiple component UI updates during single action call.
Just call `actions.Flush(core)` to initiate an update.

=== "component.foo.go"

	```go
	...

	func ComponentFoo(core *kyoto.Core) {
		lifecycle.Init(core, func() {
			core.State.Set("Status", "Not loaded")
		})
		actions.Define(core, "Bar", func(args ...interface{}) {
			// Update status UI
            c.Status = "Loading ..."
            actions.Flush(core)
            // Do some actions
            // ...
            // Update status UI again
            c.Status = "Loaded"
            actions.Flush(core)
		})
	}

	...
	```

=== "component.foo.html"

	```html
	{{ define "ComponentFoo" }}
	<div {{ componentattrs . }}>
	  <div> Status: {{ .Status }} </div>
	  <button onclick="{{ action 'Bar' }}">Bar</button>
	</div>
	{{ end }}
	```

### Rendering options

#### Mode

There are cases when morphdom may fail.
The library then falls back to replace mode instead, which just replaces element's outerHTML.
To force the library to use replace mode, you can use the `ssa:render.mode` HTML attribute.

=== "component.foo.html"

	```go
	{{ define "ComponentExample" }}
	<div {{ componentattrs . }} ssa:render.mode="replace">
		...
	</div>
	{{ end }}
	```

## Limitations and advices

* Your component must to be JSON serializable.
  This is related to how the Actions feature works under the hood.
  Component state is stored directly in DOM and can be used by client-size operations like binding.
* Avoid using of interface types in components with Actions.
  That way your component will become incompatible with JSON serialization process.
* Avoid huge states.
  This will increase the total page size and will slow down Actions operations.
  You can avoid field JSON serialization with `json:"-"` if it's not needed. i.e in case of list/table data.
  We're trying to avoid storing rows of data in the state and using database directly instead.
* When you're modifying state in the child component, state of the parent component is not updating.
  Be careful with dynamic components nesting.
