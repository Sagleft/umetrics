package frontend

import (
	"bot/pkg/memory"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ginFront struct {
	Gin    *gin.Engine
	Memory memory.Memory
}

func NewGINFrontend(db memory.Memory) (Frontend, error) {
	f := &ginFront{
		Gin:    gin.Default(),
		Memory: db,
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

	f.Gin.GET("/", f.renderHomePage)
	f.Gin.NoRoute(f.renderNotFoundPage)
	return nil
}

func (f *ginFront) renderTemplate(c *gin.Context, code int, name string, obj interface{}) {
	c.HTML(code, name, obj)
}

func (f *ginFront) renderNotFoundPage(c *gin.Context) {
	f.renderTemplate(
		c,
		http.StatusNotFound,
		"404.html",
		gin.H{},
	)
}

func (f *ginFront) renderHomePage(c *gin.Context) {
	f.renderTemplate(
		c,
		http.StatusOK,
		"home.html",
		gin.H{},
	)
}
