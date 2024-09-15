package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()

	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*")

	r.POST("/", promptHandler)
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "protected.html", gin.H{})
	})

	r.Run(":8088")
}
