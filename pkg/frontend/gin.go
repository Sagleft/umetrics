package frontend

import (
	"bot/pkg/memory"
	"encoding/json"
	"fmt"
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
	f.renderErrorPage(c, http.StatusNotFound, "page not found")
}

func (f *ginFront) renderErrorPage(c *gin.Context, statusCode int, errInfo string) {
	f.renderTemplate(
		c,
		statusCode,
		"error.html",
		gin.H{
			"code":       statusCode,
			"errorTitle": getErrorTitle(statusCode),
			"errorInfo":  errInfo,
		},
	)
}

func (f *ginFront) renderError(c *gin.Context, err error) {
	f.renderErrorPage(c, http.StatusInternalServerError, err.Error())
}

func (f *ginFront) renderHomePage(c *gin.Context) {
	peers, err := f.Memory.GetPeers()
	if err != nil {
		f.renderError(c, err)
		return
	}

	channelsCount, err := f.Memory.GetChannelsCount()
	if err != nil {
		f.renderError(c, err)
		return
	}

	peersBytes, err := json.Marshal(peers)
	if err != nil {
		f.renderError(c, fmt.Errorf("encode peers data: %w", err))
		return
	}

	f.renderTemplate(
		c,
		http.StatusOK,
		"home.html",
		gin.H{
			"peersData":     string(peersBytes),
			"channelsCount": channelsCount,
		},
	)
}
