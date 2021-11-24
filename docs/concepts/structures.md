
# Structures

Each page or component must to be represented by its own defintion and structure.  
Just think of this as a component, where the page is the highest level component with its own set of functionality.

Examples:

```go title="page.go"
package main

import (
    "html/template"
)

type PageIndex struct{}

func (p *PageIndex) Template() *template.Template {
    return template.Must(template.New("page.index.html").Funcs(kyoto.Funcs()).ParseGlob("*.html"))
}
```

```go title="component.go"
package main

type ComponentExample struct{}
```
