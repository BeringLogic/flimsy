package main

import (
  "net/http"
  "github.com/gin-gonic/gin"
)

func GET_root(c *gin.Context) {
  c.HTML(http.StatusOK, "index.tmpl", nil)
}

func GET_homepage_png(c *gin.Context) {
  c.File("/var/lib/flimsy/public/homepage.png")
}

func main() {
  gin.ForceConsoleColor()

  r := gin.Default()

  r.LoadHTMLGlob("/var/lib/flimsy/templates/*.tmpl")

  r.GET("/", func(c *gin.Context) { GET_root(c) })
  r.GET("/homepage.png", func(c *gin.Context) { GET_homepage_png(c) })

  r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
