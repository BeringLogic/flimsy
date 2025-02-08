package main

import (
  "os"
  "net/http"
  "database/sql"
  _ "github.com/mattn/go-sqlite3"
  "github.com/gin-gonic/gin"
)

func getEnv(key, def string) string {
    value, exists := os.LookupEnv(key)
    if !exists {
        value = def
    }
    return value
}

func GET_root(c *gin.Context) {
  c.HTML(http.StatusOK, "index.tmpl", gin.H{
    "auth_disabled" : true,
    "loggedIn" : false,
    "session_message" : "", 
    "FLIMSY_WEATHER_API_KEY" : getEnv("FLIMSY_WEATHER_API_KEY", ""),
    "FLIMSY_WEATHER_LOCATION" : getEnv("FLIMSY_WEATHER_LOCATION", "New York"),
    "FLIMSY_WEATHER_UNITS" : getEnv("FLIMSY_WEATHER_UNITS", "standard"),
    "FLIMSY_WEATHER_LANGUAGE" : getEnv("FLIMSY_WEATHER_LANGUAGE", "en"),
  })
}

func main() {
  db, err := sql.Open("sqlite3", ":memory:?_busy_timeout=5000&_foreign_keys=ON&_journal_mode=WAL")
  if err != nil {
    panic(err)
  }
  defer db.Close()

  gin.ForceConsoleColor()

  r := gin.Default()

  r.LoadHTMLGlob("/var/lib/flimsy/templates/*.tmpl")

  r.Static("/static", "/var/lib/flimsy/static")
  r.GET("/", func(c *gin.Context) { GET_root(c) })

  r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
