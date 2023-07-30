
<p align="center">
    <img width="200" src="logo.svg" />
</p>

<h1 align="center">kyoto</h1>

<p align="center">
    Go server side frontend framework
</p>

```go
import "go.kyoto.codes/v3"
```

Kyoto is a library for creating fast, server side frontend avoiding vanilla templating downsides.

It tries to address complexities in frontend domain like responsibility separation, components structure, asynchronous load and hassle\-free dynamic layout updates. These issues are common for frontends written with Go.

The library provides you with primitives for pages and components creation, state and rendering management, dynamic layout updates \(with client configuration\), utility functions and asynchronous components out of the box. Still, it bundles with minimal dependencies and tries to utilize built\-ins as much as possible.

You would probably want to opt out from this library in few cases, like, if you're not ready for drastic API changes between major version, you want to develop SPA/PWA and/or complex client\-side logic, or you're just feeling OK with your current setup. Please, don't compare kyoto with a popular JS libraries like React, Vue or Svelte. I know you will have such a desire, but most likely you will be wrong. Use cases and underlying principles are just too different.

If you want to get an idea of what a typical static component would look like, here's some sample code. It's very ascetic and simplistic, as we don't want to overload you with implementation details. Markup is also not included here \(it's just a well\-known \`html/template\`\).

```go
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
```

For details, please check project's website on https://kyoto.codes. Also, you may check the library index to explore available sub\-packages and https://pkg.go.dev for Go'ish documentation style.
