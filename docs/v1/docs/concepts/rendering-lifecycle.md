
# Rendering lifecycle

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
