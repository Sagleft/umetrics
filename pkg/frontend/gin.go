package frontend

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ginFront struct {
	Gin *gin.Engine
}

func NewGINFrontend() (Frontend, error) {
	f := &ginFront{
		Gin: gin.Default(),
	}
	f.setupRoutes()

	return f, nil
}

func (f *ginFront) Run() error {
	return f.Gin.Run(":8080")
}

func (f *ginFront) loadTemplates(pattern string) {
	f.Gin.LoadHTMLGlob(pattern)
}

func (f *ginFront) setupRoutes() error {
	f.loadTemplates("./templates/*")
	f.Gin.Static("/assets", "./assets")

	f.Gin.GET("/", func(c *gin.Context) {
		f.renderTemplate(
			c,
			http.StatusOK,
			"home.html",
			gin.H{},
		)
	})
	f.Gin.NoRoute(func(c *gin.Context) {
		f.renderTemplate(
			c,
			http.StatusNotFound,
			"404.html",
			gin.H{},
		)
	})

	return nil
}

func (f *ginFront) renderTemplate(c *gin.Context, code int, name string, obj interface{}) {
	c.HTML(code, name, obj)
}
