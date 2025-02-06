package main

import (
  "net/http"
  "github.com/gin-gonic/gin"
)

func GET_root(c *gin.Context) {
  c.HTML(http.StatusOK, "index.tmpl", nil)
}

func main() {
  gin.ForceConsoleColor()

  r := gin.Default()

  r.LoadHTMLGlob("/var/lib/flimsy/templates/*.tmpl")

  r.Static("/static", "/var/lib/flimsy/static")
  r.GET("/", func(c *gin.Context) { GET_root(c) })

  r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
