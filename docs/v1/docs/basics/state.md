
# State

The most obvious part of any component-based framework - state management.
Kyoto is far away from reactive approach, but also have own state system.

Let's take a look.

=== "main.go"

	```go

	...

	func ComponentFoo(core *kyoto.Core) {
		lifecycle.Init(core, func() {
			core.State.Set("Foo", "Bar")
		})
	}

	func PageBar(core *kyoto.Core) {
		lifecycle.Init(core, func() {
			core.Component("Foo", ComponentFoo)
		})
		lifecycle.Async(core, func() error {
			// You will get a state of ComponentFoo instance (as a map[string]interface{})
			fmt.Println(core.State.Get("Foo"))
			return nil
		})
	}

	...
	```

So, what is state for kyoto?
It's a way to pass data structure for template rendering.
If we're talking about dynamic components, it's also a way to save component data for dynamic actions after page rendering.

In the basement of state we have `kyoto.Store`, a simple atomic map wrapper.
To interact with a state we are using it's instance inside of a core (`core.State`).
There are 3 methods for state modification: `Set`, `Get` and `Del`.

Store instance is unique for each page or component and represents only current scope (unlike Context and Scheduler).
Under the hood, store's underlying map becomes injected into the parent's store.
