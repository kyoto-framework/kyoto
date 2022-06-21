package kyoto

import (
	"html/template"
)

// ****************
// Template configuration
// ****************

// TemplateConfiguration holds template building configuration.
type TemplateConfiguration struct {
	ParseGlob string
	FuncMap   template.FuncMap
}

// FuncMap holds a library predefined template functions.
// You have to include it into TemplateConf.FuncMap (or your raw templates) to use kyoto properly.
var FuncMap = template.FuncMap{
	"await":  Await,
	"state":  actionFuncState,
	"client": actionFuncClient,
}

// TemplateConf is a global configuration that will be used during template building.
// Feel free to modify it as needed.
var TemplateConf = TemplateConfiguration{
	ParseGlob: "*.html",
	FuncMap:   FuncMap,
}

// ComposeFuncMap is a function for composing multiple FuncMap instances into one.
//
// Example:
//
//		func MyFuncMap() template.FuncMap {
//			return kyoto.ComposeFuncMap(
//				funcmap1,
//				funcmap2,
//				...
//			)
//		}
//
func ComposeFuncMap(fmaps ...template.FuncMap) template.FuncMap {
	var result = template.FuncMap{}
	for _, fmap := range fmaps {
		for k, v := range fmap {
			result[k] = v
		}
	}
	return result
}

// ****************
// Template building functions
// ****************

// Template creates a new template with a given name,
// using global parameters stored in kyoto.TemplateConf.
// Stores template in the context.
//
// Example:
//
//		func PageFoo(ctx *kyoto.Context) (state PageFooState) {
//			// By default uses kyoto.FuncMap
//			// and parses everything in the current directory with a .ParseGlob("*.html")
//			kyoto.Template(ctx, "page.foo.html")
//			...
//		}
//
func Template(c *Context, name string) {
	c.Template = template.Must(template.New(name).Funcs(TemplateConf.FuncMap).ParseGlob(TemplateConf.ParseGlob))
}

// TemplateInline creates a new template with a given template source,
// using global parameters stored in kyoto.TemplateConf.
// Stores template in the context.
//
// Example:
//
//		func PageFoo(ctx *kyoto.Context) (state PageFooState) {
//			// By default uses kyoto.FuncMap
//			// and parses everything in the current directory with a .ParseGlob("*.html")
//			kyoto.TemplateInline(ctx, `<html>...</html>`)
//			...
//		}
//
func TemplateInline(c *Context, tmpl string) {
	c.Template = template.Must(template.Must(template.New("inline").Funcs(TemplateConf.FuncMap).ParseGlob(TemplateConf.ParseGlob)).Parse(tmpl))
}

// TemplateRaw handles a raw template.
// Stores template in the context.
//
// Example:
//
//		func PageFoo(ctx *kyoto.Context) (state PageFooState) {
//			tmpl := MyTemplateBuilder("page.foo.html") // *template.Template
//			kyoto.TemplateRaw(ctx, tmpl)
//			...
//		}
//
func TemplateRaw(c *Context, tmpl *template.Template) {
	c.Template = tmpl
}
