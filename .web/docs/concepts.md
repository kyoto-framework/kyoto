# Concepts

In this part you'll understand basic library concepts, page rendering lifecycle and how to use most of library features.

## Structures

Each page or component must to be represented by its own structure.  
Just think of this structure as a component, where page - toplevel component with own set of functionality.

## Rendering lifecycle

![arch-rendering](https://i.imgur.com/72xIkzx.png)

> Diagram might become outdated with time, as far as library in active development

In the plain `html/template` there is no such thing as lifecycle. Lifecycle concept was taken from popular JS frameworks.  
Let's explore how lifecycle works in `kyoto`:

- Definition of shared variables (waitgroup, errors channel)
- Initializing page and registering components
- Triggering asynchronous components in separate goroutines
- Waiting untill all asynchronous operations are completed
- If new components was registered during asychronous operations, repeat async steps for new components
- Triggering after async operations
- Cleanup registered components
- Rendering final page to `io.Writer`

## Lifecycle integration

To integrate your component into lifecycle, you need at least 2 things to be done:  

- Register component on page initialization
- Implement lifecycle method to match one of existing interfaces (will be described further)

Components registration must to be done on page initialization stage. It can be done with implementing `ImplementsInitWithoutPage` interface and registering components with `RegC` method. `RegC` method immediately returns initialized component pointer, so you can use it on component assign.  
Example:

```go
type PageIndex struct {
    Example kyoto.Component
}

...

func (p *PageIndex) Init() {
    p.Example = kyoto.RegC(p, &ComponentExample{})
}
```

After component registration, you can use all lifecycle features. Let's explore them:

- Initialization step. Triggered immediately on component registration. You can integrate your component into this step with `Init` method
- Asynchronous operations step. Triggered concurrently after all initializations are done. You can integrate your component into this step with `Aync` method
- After asynchronous operations step. Triggered after all concurrent operations are done. You can integrate your component into this step with `AfterAsync` method

You can find detailed usage examples in [Core Features](/docs/core-features) section.  
This overview not includes features from [Extended Features](/docs/extended-features) section, because most of them not related to lifecycle directly.

## Methods overloading

This library can handle different method signatures. Under the hood it tries to cast your component to different interfaces, depending on current lifecycle step. This approach allows to reduce code complexity and extend library features with minimal pain.  
For example, you can use `Async(p kyoto.Page)` or `Async()` signature depending on your needs.  

You can find detailed usage examples in [Core Features](/docs/core-features) section.  

## Interfaces

This library strongly relies on interfaces. For methods overloading, lifecycle integrations, etc.  
List of all (or almost all) available interfaces:

- `ImplementsInit` - interface for checking implementation of initialization method.  
  - Part of [lifecycle](/docs/concepts/#rendering-lifecycle).  
  - Appliable for components.
- `ImplementsInitWithoutPage` - same as `ImplementsInit`, but without `kyoto.Page` argument.  
  - Part of [lifecycle](/docs/concepts/#rendering-lifecycle).  
  - Appliable for both pages and components.
- `ImplementsMeta` - interface for checking implementation of meta builder.  
  - Appliable for pages.
- `ImplementsAsync` - interface for checking implementation of asynchronous method.  
  - Part of [lifecycle](/docs/concepts/#rendering-lifecycle).  
  - Appliable for components.
- `ImplementsAsyncWithoutPage` - same as `ImplementsAsync`, but without `kyoto.Page` argument.  
  - Part of [lifecycle](/docs/concepts/#rendering-lifecycle).  
  - Appliable for components.
- `ImplementsAfterAsync` - interface for checking implementation of after async method.  
  - Part of [lifecycle](/docs/concepts/#rendering-lifecycle).  
  - Appliable for components.
- `ImplementsAfterAsyncWithoutPage` - same as `ImplementsAfterAsync`, but without `kyoto.Page` argument.  
  - Part of [lifecycle](/docs/concepts/#rendering-lifecycle).  
  - Appliable for components.
- `ImplementsAtions` - interface for checking implementation of [Server Side Actions](/docs/extended-features/#server-side-actions).  
  - Appliable for components.

You can find complete list of interfaces here:
[https://github.com/yuriizinets/kyoto/blob/master/types.go](https://github.com/yuriizinets/kyoto/blob/master/types.go)
