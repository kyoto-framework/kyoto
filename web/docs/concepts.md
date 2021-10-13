# Concepts

## Interfaces

Each page or component is represented by its own structure.
For implementing specific functionality, you need to implement one of predefined interfaces. You need to follow declaration rules to match the required interface (you can find all interfaces in [`types.go`](https://github.com/yuriizinets/ssceng/blob/master/types.go#L59)).  
In some cases, you can use methods overloading with extended interfaces. This allows to simplify setup and avoid unnecessary code.

### Page interfaces

- [`Page`](https://github.com/yuriizinets/ssceng/blob/master/types.go#L51) - main page interface with minimal requirements
- [`ImplementsInitWithoutPage`](https://github.com/yuriizinets/ssceng/blob/master/types.go#L61) - page initialization method
- [`ImplementsMeta`](https://github.com/yuriizinets/ssceng/blob/master/types.go#L81) - page meta builder, you can find more [here](/extended.html#meta-builder)

### Component interfaces

- [`Component`](https://github.com/yuriizinets/ssceng/blob/master/types.go#L57) - main component interface with minimal requirements
- [`ImplementsInit`](https://github.com/yuriizinets/ssceng/blob/master/types.go#L61) - component initialization method, for initializing default values or registering nested components
- [`ImplementsAsync`](https://github.com/yuriizinets/ssceng/blob/master/types.go#L69) - async method will be called concurrently with another async methods
- [`ImplementsAfterAsync`](https://github.com/yuriizinets/ssceng/blob/master/types.go#L77) - method is called when all async method finished execution
- [`ImplementsActions`](https://github.com/yuriizinets/ssceng/blob/master/types.go#L85) - method, returning [`ssc.ActionMap`](https://github.com/yuriizinets/ssceng/blob/master/types.go#L11) with component [SSA](/extended#server-side-actions-ssa) methods

Overloads

- [`ImplementsInitWithoutPage`](https://github.com/yuriizinets/ssceng/blob/master/types.go#L65) - same as [`ImplementsInit`](https://github.com/yuriizinets/ssceng/blob/master/types.go#L61), but without page argument
- [`ImplementsAsyncWithoutPage`](https://github.com/yuriizinets/ssceng/blob/master/types.go#L73) - same as [`ImplementsAsync`](https://github.com/yuriizinets/ssceng/blob/master/types.go#L69), but without page argument
- [`ImplementsAfterAsyncWithoutPage`](https://github.com/yuriizinets/ssceng/blob/master/types.go#L81) - same as [`ImplementsAfterAsync`](https://github.com/yuriizinets/ssceng/blob/master/types.go#L77), but without page argument
- [`ImplementsActionsWithoutPage`](https://github.com/yuriizinets/ssceng/blob/master/types.go#L89) - same as [`ImplementsActions`](https://github.com/yuriizinets/ssceng/blob/master/types.go#L85), but without page argument

## Lifecycle

Before implementing any method, you need to understand the rendering lifecycle.  

### Page rendering

Each page's lifecycle is hidden under the render function and follows this steps:

- Defining shared variables (waitgroup, errors channel)
- Triggering the page's `Init()` to initialize and register components
- Running all component's `Async()` functions in separate goroutines
- Waiting untill all asynchronous operations are completed
- If new components were registered while `Async` execution, repeat `Async` stage for newly created components
- Calling `AfterAsync()` for each component
- Cleaning up registered components (not needed more for internal usage)
- Getting page's template and render

### SSA Rendering

If you want to use SSA in your project, it's better to know how it works first. SSA has own, shorten lifecycle.  

- Creating request with JS on client side (icluding component name, state, action name, args)
- Extracting action data from request on server side
- Finding registered component type
- Creating component struct
- Triggering the component's `Init()`
- Populating component's state
- Calling action
- If new components where registed while action execution, do `Async` stage in the same way as for page
- Rendering component and returning HTML to client side
- Replacing component's HTML with recieved version
