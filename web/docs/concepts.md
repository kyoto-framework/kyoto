# Concepts

Each page or component is represented by its own structure. For implementing specific functionality, you can use structure's methods with a predefined declaration (f.e. `Init(p ssc.Page)`). You need to follow declaration rules to match the interfaces required (you can find all interfaces in `types.go`).  
Before implementing any method, you need to understand the rendering lifecycle.

## Component structure

## Lifecycle

Each page's lifecycle is hidden under the render function and follows this steps:

- Defining shared variables (waitgroup, errors channel)
- Triggering the page's `Init()` to initialize and register components
- Running all component's `Async()` functions in separate goroutines
- Waiting untill all asynchronous operations are completed
- Calling `AfterAsync()` for each component
- Cleaning up registered components (not needed more for internal usage)
- Getting page's template and render

## SSA Lifecycle

If you want to use SSA in your project, it's better to know how it works first. SSA has own, shorten lifecycle.  

- Creating request with JS on client side
- Extracting action data from request (component name, component state, action name, action args)
- Finding registered component type
- Creating component struct
- Triggering the component's `Init()`
- Populating component's state
- Calling action
- Rendering component and returning HTML to client side
- Replacing component's HTML with recieved version
