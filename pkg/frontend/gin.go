package frontend

import (
	"bot/pkg/memory"
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type ginFront struct {
	Gin    *gin.Engine
	Memory memory.Memory

	filehashes map[string]string
}

func NewGINFrontend(db memory.Memory) (Frontend, error) {
	f := &ginFront{
		Gin:        gin.Default(),
		Memory:     db,
		filehashes: make(map[string]string),
	}

	if err := f.hashFiles(); err != nil {
		return nil, err
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
	f.Gin.GET("geoData.json", f.renderGeoData)
	f.Gin.NoRoute(f.renderNotFoundPage)
	return nil
}

func (f *ginFront) hashFiles() error {
	for _, filePath := range frontFilesToHash {
		file, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("read file to get hash: %w", err)
		}

		defer file.Close()

		hash := md5.New()
		_, err = io.Copy(hash, file)
		if err != nil {
			return fmt.Errorf("get file hash: %w", err)
		}

		f.filehashes[getFileHashName(file)] = string(hash.Sum(nil))
	}

	return nil
}

func getFileHashName(file *os.File) string {
	return strings.ReplaceAll(file.Name(), ".", "_")
}

func (f *ginFront) renderGeoData(c *gin.Context) {
	peers, err := f.Memory.GetPeers()
	if err != nil {
		f.renderError(c, err)
		return
	}

	c.JSON(http.StatusOK, peers)
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
			"version":    f.filehashes,
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
	channelsCount, err := f.Memory.GetChannelsCount()
	if err != nil {
		f.renderError(c, err)
		return
	}

	usersCount, err := f.Memory.GetUsersCount()
	if err != nil {
		f.renderError(c, err)
		return
	}

	topChannels, err := f.Memory.GetTopChannels(maxTopChannels)
	if err != nil {
		f.renderError(c, err)
		return
	}

	topUsers, err := f.Memory.GetTopUsers(maxTopUsers)
	if err != nil {
		f.renderError(c, err)
		return
	}

	f.renderTemplate(
		c,
		http.StatusOK,
		"home.html",
		gin.H{
			"version":       f.filehashes,
			"channelsCount": channelsCount,
			"usersCount":    usersCount,
			"topChannels":   topChannels,
			"topUsers":      topUsers,
		},
	)
}
