package main

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestWRK(t *testing.T) {
	server := gin.Default()

	server.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, "ok")
		return
	})

	server.Run("localhost:8082")
}
