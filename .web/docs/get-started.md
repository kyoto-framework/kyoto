# Get Started

Overall project setup may seem complicated and unusual. But you will get used to it and realize all the benefits over time.  
If you want to implement own project from zero, or extend already existing, it's highly recommended to rely on demo projects or examples.  

## Quick Start

The quickest way to get started with `kyoto` is to clone ready demo project, like [`kyoto-hn`](https://github.com/yuriizinets/kyoto-hn), and try to expore/modify things inside.

## Installation

To install this library, just use `go get github.com/yuriizinets/kyoto`

## Integration

It's easy to integrate into existing projects. There are 2 ways to integrate: define own generic page handler, or use low-level page rendering method. In the most cases, second method is the easiest (especially when no context is needed).  

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

Check [Page rendering](/core-features/#page-rendering) section for details.
