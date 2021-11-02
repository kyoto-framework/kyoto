# Additional notes

## Dealing with Go packages

### General

There are few rules and limitations that you need to follow when working with packages:

- Component template definition name needs to be the same as component struct name.
- There are no namespaces or packages for templates.
- Component name needs to be unique (due to lack of namespaces in templates)

### Shared component library

For now, it's impossible to define component, that can be uploaded/downloaded with go packages.
Go packaging system doesn't support external resources like templates, styles, etc.  

Now, there are 2 ways of installing shared library:

- Adding library as git submodule with `git submodule add <package> <folder>` + adding `replace` statement into `go.mod` file
- Installing as usual Go package with `go get <package>` with manual copying of needed templates into local templates folder

There are plans to extend library with `Render` interface, where developer can define own way to render component/page, 
without relying on `html/template` (but not excluding it). In that way developer will be able to include all needed markup inside of component's code.

## Downsides of the library

- Pretty complex and hard to understand setup
- There is no option to provide "Quick Start" documentation as far as library created for complex setup
- Actions feature sometimes is hard to use due to [downsides](/extended-features/#ssa-limitations)
- As far as Actions feature relies on server side code, network connection is quite important
- Not so feature-rich as modern frontend frameworks, as far as library tries to solve another issues
