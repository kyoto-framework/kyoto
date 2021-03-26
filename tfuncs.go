package gofr

import "html/template"

func Funcs() template.FuncMap {
	return template.FuncMap{
		"meta": func(p Page) template.HTML {
			builder := ""
			meta := p.Meta()
			if meta.Title != "" {
				builder += "<title>" + meta.Title + "</title>\n"
			}
			if meta.Canonical != "" {
				builder += `<link rel="canonical" href="` + meta.Canonical + `">` + "\n"
			}
			if len(meta.Hreflangs) != 0 {
				for _, hreflang := range meta.Hreflangs {
					builder += `<link rel="alternate" hreflang="` + hreflang.Lang + `" href="` + hreflang.Href + `">` + "\n"
				}
			}
			return template.HTML(builder)
		},
	}
}
