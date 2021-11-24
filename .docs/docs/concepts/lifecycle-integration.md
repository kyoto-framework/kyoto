
# Lifecycle Integration

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

You can find detailed usage examples in the `Core Features` section.  
This overview does not include features from the `Extended Features` section, because most of them not related to the lifecycle directly.
