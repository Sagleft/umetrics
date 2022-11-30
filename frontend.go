package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (b *bot) setupFrontend() error {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	return r.Run(":8080")
}
