package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// gin
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.GET("/index", func(c *gin.Context) {
		data := struct {
			Title string
			Items []string
		}{
			Title: "My page",
			Items: []string{
				"My photos",
				"My blog",
			},
		}
		c.HTML(http.StatusOK, "index.tmpl", data)
	})

	router.Run() // listens on 0.0.0.0:8080 by default
}
