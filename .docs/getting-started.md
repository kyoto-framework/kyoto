# Getting Started

The overall project setup may seem complicated and unusual but you will soon get used to it and realize all the benefits over time.  
If you want to implement your own project from scratch, or extend an already existing project, it's highly recommended to rely on demo projects and/or examples.

## Quick Start

The quickest way to get started with `kyoto` is to use a starter project.  
[https://github.com/yuriizinets/kyoto-starter](https://github.com/yuriizinets/kyoto-starter)

## Installation

To install this library, just use `go get github.com/yuriizinets/kyoto`

## Integration

It's easy to integrate into existing projects. There are 2 ways to integrate: define own generic page handler, or use low-level page rendering method. In most cases, the second method is the easiest (especially when no context is needed).

Example with `net/http`:

```go
func IndexPageHandler(rw http.ResponseWriter, r *http.Request) {
    RenderPage(rw, &IndexPage{})
}
```

Example with `gin`:

```go
func IndexPageHandler(g *gin.Context) {
    kyoto.RenderPage(g.Writer, &IndexPage{})
}
```

Check the [Page rendering](/core-features/#page-rendering) section for details.
