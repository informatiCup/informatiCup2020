package main

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Delims("{{{", "}}}")
	r.Static("/static", "./static")
	r.LoadHTMLGlob("./gohtml/*.gohtml")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.gohtml", nil)
	})
	r.GET("/game/:name", func(c *gin.Context) {
		c.File(fmt.Sprintf("games/%s.json", c.Param("name")))
	})
	r.Run(":50124")
}
