package main

import (
	"html/template"
	"os"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/yuriizinets/go-ssc"
)

func funcmap() template.FuncMap {
	return ssc.Funcs()
}

func main() {
	g := gin.Default()

	g.GET("/", func(c *gin.Context) {
		ssc.RenderPage(c.Writer, &PageIndex{})
	})

	g.Use(static.Serve("/static/", static.LocalFile("./static", true)))

	g.POST("/SSA/:Component/:Action", func(c *gin.Context) {
		ssc.HandleSSA(
			c.Writer,
			template.Must(template.New(c.Param("Component")).Funcs(funcmap()).ParseGlob("*.html")),
			c.Param("Component"),
			c.PostForm("State"),
			c.Param("Action"),
			c.PostForm("Args"),
			[]ssc.Component{
				&ComponentCounter{},
				&ComponentSampleBinding{},
				&ComponentSampleParent{},
				&ComponentSampleChild{},
			},
		)
	})

	addr := "localhost:25025"
	if os.Getenv("PORT") != "" {
		addr = ":" + os.Getenv("PORT")
	}
	g.Run(addr)
}
