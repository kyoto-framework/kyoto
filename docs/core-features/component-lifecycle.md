
# Component lifecycle

This section extends `Lifecycle integration` documentation with examples.

## Init

`Init` method is triggering on initialization lifecycle step.

Usage:

```go
...

func (*ComponentExample) Init() {
    // Do what you want here
}

...
```

This method have overload option with page argument

```go
...

func (*ComponentExample) Init(p kyoto.Page) {
    // Do what you want here
}

...
```

## Async

`Async` method is triggering an asynchronous operations lifecycle step. You need to return an error in case of async operation failure.

> Very useful in case of time-consuming operations, when you need to fetch data from external API or database.

Usage:

```go
...

func (*ComponentExample) Async() error {
    // Do what you want here
    return nil
}

...
```

This method have overload option with page argument

```go
...

func (*ComponentExample) Async(p kyoto.Page) error {
    // Do what you want here
    return nil
}

...
```

## AfterAsync

`AfterAsync` method is triggering after asynchronous operations lifecycle step.

Usage:

```go
...

func (*ComponentExample) AfterAsync() {
    // Do what you want here
}

...
```

This method have overload option with page argument

```go
...

func (*ComponentExample) AfterAsync(p kyoto.Page) {
    // Do what you want here
}

...
```
