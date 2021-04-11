package main

import (
	"html/template"

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

	g.Run("localhost:25025")
}
