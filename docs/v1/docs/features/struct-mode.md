
# Struct Mode

Core receivers are not the only way to define a component or page.
Kyoto allows you to use structure with methods instead of core receiver with adapters.
It might be useful for highly customizable components.
Thanks to struct mode, you can use struct fields as arguments during instance creation.

## Usage

### Pages

Let's start from page definition, the most basic version.
We need to define a structure and `Template` method for that structure.

=== "page.index.go"

	```go
	package main

	import (
		"html/template"

		"github.com/kyoto-framework/kyoto/render"
	)


	type PageIndex struct {}

	func (p *PageIndex) Template() *template.Template {
		template.Must(template.New("page.index.html").Funcs(render.FuncMap()).ParseGlob("*.html"))
	}
	
	```

It would be nice to open [Basics • Pages](basics/pages.md) to compare approach.
You will definitely notice similar parts.
Struct mode tries to replicate existing adapters as much as possible.

Next,let's attach this page to our router.

=== "main.go"

	```go
	...

	mux.HandleFunc("/", render.PageHandler(smode.Adapt(&PageIndex{})))

	...
	```

I agree, it looks a bit explicit.
But it will make a sense if I will tell you that `smode.Adapt` function translates our structure to core receiver.
In that way we can use adapted structure in existing functional code.

!!! note ""
	To simplify struct pages registration, you can use own small wrapper.

	```go
	func myhandler(page smode.Page) http.HandlerFunc {
		return render.PageHandler(smode.Adapt(page))
	}
	```

You can define `Init` method, that will represent `lifecycle.Init` in this structure.

```go
...

func (p *PageIndex) Init() {
	// Do what you want here
}

...
```

### Components

As well as pages, you can define structure components.
Let's start right from an example.

=== "component.uuid.go"

	```go
	package main

	type ComponentUUID struct {
		UUID string
	}

	func (c *ComponentUUID) Init() {
		c.UUID = "None"
	}

	func (c *ComponentUUID) Async() error {
		resp, _ := http.Get("http://httpbin.org/uuid")
        data := map[string]string{}
        json.NewDecoder(resp.Body).Decode(&data)
        c.UUID = data["uuid"]
        return nil
	}
	```

You may notice how we defined lifecycle methods here.
Feel free to open [Basics • Components](basics/components.md) to compare approach.
`smode.Adapt` function takes care about registration of our methods.
Also, in case of struct definitions we are using struct fields as a state.

And now let's attach this component to the page multiple times.

=== "page.index.go"

	```go
	...

	type PageIndex struct {
		UUID1 smode.Component
		UUID2 smode.Component
	}

	func (p *PageIndex) Init() {
		p.UUID1 = smode.RegC(p, &ComponentUUID{})
		p.UUID2 = smode.RegC(p, &ComponentUUID{})
	}

	...
	```

As you can see, we are using `smode.RegC` for attaching components as an alternative to `core.Component`.
You can pass parameters to a component on initialization.

### Context

Instead of `core.Context` we can use `smode.SetContext`, `smode.GetContext` and `smode.DelContext` for context management.
Context uses page instance as namespace for correct concurrency handling on requests level (page instance is creating for each new request).

=== "page.index.go"

	```go
	...

	type PageIndex struct{}

	func (p *PageIndex) Init() {
		smode.SetContext(p, "key", "value")
	}

	func (p *PageIndex) Async() error {
		println(smode.GetContext(p, "key"))
		return nil
	}

	...
	```

### Actions

To define actions for a component you can implement `Actions` method.
This method must return `smode.ActionMap`, a map which holds your methods.

=== "component.foo.go"

	```go
	...

	type ComponentFoo struct {
		Status string
	}

	func (c *ComponentFoo) Actions() smode.ActionMap {
		return smode.ActionMap{
			`Bar`: func(args ...interface{}) {
				c.Status = "Baz"
			}
		}
	}
	```

## Limitations

- You can't use custom kyoto modules, only built-in `smode` functions.
