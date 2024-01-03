/*
Kyoto is a library for creating fast, server side frontend avoiding vanilla templating downsides.

It tries to address complexities in frontend domain like
responsibility separation, components structure, asynchronous load
and hassle-free dynamic layout updates.
These issues are common for frontends written with Go.

The library provides you with primitives for pages and components creation,
state and rendering management, dynamic layout updates (with client configuration),
utility functions and asynchronous components out of the box.
Still, it bundles with minimal dependencies
and tries to utilize built-ins as much as possible.

You would probably want to opt out from this library in few cases, like,
if you're not ready for drastic API changes between major version,
you want to develop SPA/PWA and/or complex client-side logic,
or you're just feeling OK with your current setup.
Please, don't compare kyoto with a popular JS libraries like React, Vue or Svelte.
I know you will have such a desire, but most likely you will be wrong.
Use cases and underlying principles are just too different.

If you want to get an idea of what a typical static component would look like, here's some sample code.
It's very ascetic and simplistic, as we don't want to overload you with implementation details.
Markup is also not included here (it's just a well-known `html/template`).

	// State is declared separately from component itself
	type ComponentState struct {
		// We're providing component with some abilities here
		component.Universal // This component uses universal state, that can be (un)marshalled with both server and client

		// Actual component state is just struct fields
		Foo string
		Bar string
	}

	// Component is a function, that returns configured state.
	// To be able to provide additional arguments to the component on initialization,
	// you have to wrap component with additional function that will handle args and return actual component.
	// Until then, you may keep component declaration as-is.
	func Component(ctx *component.Context) component.State {
		// Initialize state here.
		// As far as component.Universal provided in declaration,
		// we're implementing component.State interface.
		state := &ComponentState{}
		// Do whatever you want with a state
		state.Foo = "foo"
		state.Bar = "bar"
		// Done
		return state
	}

For details, please check project's website on https://kyoto.codes.
Also, you may check the library index to explore available sub-packages
and https://pkg.go.dev for Go'ish documentation style.

# Quick start

We don't want you to deal with boilerplate code on your own,
so you can proceed with our simple starter project.

	$ git clone https://kyoto.codes/new <your-new-project>
	$ rm -r <your-new-project>/.git

Feel free to use it as an example for your own setup.

# Components

Components is a common approach for modern libraries to manage frontend parts.
Kyoto components are trying to be mostly independent (but configurable) part of the project.

To create component, it would be enough to implement component.Component.
It's a function, a context receiver which returns a component state.
State is an implementation of component.State,
which is easy to implement with nesting one of the state implementations (options will be described later).

	package main

	type ComponentState struct {
		component.Disposable // You're providing component with some abilities here
	}

	func Component(ctx *component.Context) component.State {
		state := &ComponentState{}
		return state
	}

	...

Each component becomes a part of the page or top-level component,
which executes component function asynchronously and gets a state future object.
In that way your components are executing in a non-blocking way.

Pages are just top-level components, where you can configure rendering and page related stuff.

# Components with state

Stateful components are pretty similar to stateless ones,
but they are actualy implementing marshal/unmarshal interface instead of mocking it.

You have multiple state options to choose from: universal or server.

Universal state is a state, that can be marshalled and unmarshalled both on server and client.
It's a common state option without functionality limitations.
On the other hand, the whole state must to be sent and received,
which applies some limitations on the state size.

	package main

	type ComponentState struct {
		component.Universal // This state allows you to operate with data on both server and client
	}

	func Component(ctx *component.Context) component.State {
		state := &ComponentState{}
		return state
	}

Server state can be marshalled and unmarshalled only on server.
It's a good option for components, that are not supposed to be updated on client side (f.e. no inputs).
Also, it's a good option for components with lots of state data.

	package main

	type ComponentState struct {
		component.Server // This state allows you to operate with data on server only
	}

	func Component(ctx *component.Context) component.State {
		state := &ComponentState{}
		return state
	}

# Components with arguments

Sometimes you may want to pass some arguments to the component.
It's easy to do with wrapping component with additional function.

	package main

	type ComponentState struct {
		component.Universal

		Data string
	}

	func Component(data string) component.Component {
		return func(ctx *component.Context) component.State {
			state := &ComponentState{}
			state.Data = data // We are passing arg to the component state, but it's not a requirement.
			return state
		}
	}

# Rendering

...

# Routing

...

# Actions

...
*/
package kyoto
