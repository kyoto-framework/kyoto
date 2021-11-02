# Concepts

In this page you'll learn the basic library concepts such as the page rendering lifecycle and how to use most of library features.

## Structures

Each page or component must to be represented by its own defintion and structure.  
Just think of this as a component, where the page is the highest level component with its own set of functionality.

## Rendering Lifecycle

![arch-rendering](https://i.imgur.com/72xIkzx.png)

> Diagram might become outdated with time, as far as library in active development

In the plain `html/template` there is no such thing as a lifecycle. The lifecycle concept was taken from popular JS frameworks.

Let's explore how the lifecycle works in `kyoto`:

- Definition of shared variables (waitgroup, errors channel)
- Initializing the page and registering components
- Triggering asynchronous components in separate goroutines
- Waiting until all asynchronous operations are completed
- If new components are registered during asychronous operations, repeat async steps for new components
- Triggering after async operations
- Clean-up registered components
- Rendering the final page to `io.Writer`

## Lifecycle Integration

To integrate your component into the lifecycle, you must do two things:

- Register the component within the page initialization
- Implement the lifecycle method to match one of existing interfaces (will be described further later)

Component registration must to be done within the page initialization stage. It can be done with implementing the `ImplementsInitWithoutPage` interface and registering components with the`RegC` method. `RegC` method will return an initialized component pointer, so you can use it to assign the component.

For example:

```go
type PageIndex struct {
    Example kyoto.Component
}

...

func (p *PageIndex) Init() {
    p.Example = kyoto.RegC(p, &ComponentExample{})
}
```

After component registration, you can use all the lifecycle features. Let's explore them:

- Initialization step: Triggered immediately on component registration. You can integrate your component into this step with the `Init` method
- Asynchronous operations step: Triggered concurrently after all initializations are done. You can integrate your component into this step with the `Aync` method
- After asynchronous operations step: Triggered after all concurrent operations are done. You can integrate your component into this step with the `AfterAsync` method

You can find detailed usage examples in the [Core Features](/core-features) section.  
This overview does not include features from the [Extended Features](/extended-features) section, because most of them not related to the lifecycle directly.

## Method Overloading

This library can handle different method signatures. Under the hood it tries to cast your component to different interfaces, depending on the current lifecycle step. This approach allows you to reduce code complexity and extend library features with minimal pain.  
For example, you can use `Async(p kyoto.Page)` or `Async()` signature depending on your needs.

You can find detailed usage examples in the [Core Features](/core-features) section.

## Interfaces

This library strongly relies on interfaces. For method overloading, lifecycle integrations, etc.  
List of all (or almost all) available interfaces:

- `ImplementsInit` - interface for checking implementation of initialization method.
  - Part of [lifecycle](/concepts/#rendering-lifecycle).
  - Appliable for components.
- `ImplementsInitWithoutPage` - same as `ImplementsInit`, but without `kyoto.Page` argument.
  - Part of [lifecycle](/concepts/#rendering-lifecycle).
  - Appliable for both pages and components.
- `ImplementsMeta` - interface for checking implementation of meta builder.
  - Appliable for pages.
- `ImplementsAsync` - interface for checking implementation of asynchronous method.
  - Part of [lifecycle](/concepts/#rendering-lifecycle).
  - Appliable for components.
- `ImplementsAsyncWithoutPage` - same as `ImplementsAsync`, but without `kyoto.Page` argument.
  - Part of [lifecycle](/concepts/#rendering-lifecycle).
  - Appliable for components.
- `ImplementsAfterAsync` - interface for checking implementation of after async method.
  - Part of [lifecycle](/concepts/#rendering-lifecycle).
  - Appliable for components.
- `ImplementsAfterAsyncWithoutPage` - same as `ImplementsAfterAsync`, but without `kyoto.Page` argument.
  - Part of [lifecycle](/concepts/#rendering-lifecycle).
  - Appliable for components.
- `ImplementsAtions` - interface for checking implementation of [Server Side Actions](/extended-features/#server-side-actions).
  - Appliable for components.

You can find the complete list of interfaces here:
[https://github.com/yuriizinets/kyoto/blob/master/types.go](https://github.com/yuriizinets/kyoto/blob/master/types.go)
