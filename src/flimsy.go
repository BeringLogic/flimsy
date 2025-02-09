package main

import (
  "os"
  "fmt"
  "net/http"
  "github.com/gin-gonic/gin"

  "github.com/BeringLogic/flimsy/db"
)

type listAndItems struct {
  List *db.List
  Items []*db.Item
}

var config *db.Config
var lists *[]db.List
var items *[]db.Item
var listsAndItems []listAndItems

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
    "IsLoggedIn" : false,
    "session_message" : "", 
    "FLIMSY_WEATHER_API_KEY" : getEnv("FLIMSY_WEATHER_API_KEY", ""),
    "FLIMSY_WEATHER_LOCATION" : getEnv("FLIMSY_WEATHER_LOCATION", "New York"),
    "FLIMSY_WEATHER_UNITS" : getEnv("FLIMSY_WEATHER_UNITS", "standard"),
    "FLIMSY_WEATHER_LANGUAGE" : getEnv("FLIMSY_WEATHER_LANGUAGE", "en"),
    "config" : config,
    "listsAndItems" : listsAndItems,
  })
}

func GET_config(c *gin.Context) {
  c.JSON(http.StatusOK, config)
}

func GET_onlineStatus(c *gin.Context) {
  url := c.Query("href")

  resp, err := http.Get(url); if err != nil {
    c.JSON(http.StatusOK, gin.H{
      "online" : false,
      "error" : err.Error(),
    })
  } else {
    c.JSON(http.StatusOK, gin.H{
      "online" : true,
    })
  }
  resp.Body.Close()
}

func InitDB() error {
  err := db.Open(); if err != nil {
    return fmt.Errorf("Failed to open DB: %s", err)
  }

  config, err = db.LoadConfig(); if err != nil {
    err = db.Seed(); if err != nil {
      return fmt.Errorf("Failed to seed DB: %s", err)
    }
    config, err = db.LoadConfig(); if err != nil {
      return fmt.Errorf("Failed to load config: %s", err)
    }
  }

  lists, err = db.LoadLists(); if err != nil {
    return fmt.Errorf("Failed to load lists: %s", err)
  }

  items, err = db.LoadItems(); if err != nil {
    return fmt.Errorf("Failed to load items: %s", err)
  }

  for _, list := range *lists {
    var lai listAndItems
    lai.List = &list
    for _, item := range *items {
      if item.List_id == list.Id {
        lai.Items = append(lai.Items, &item)
      }
    }
    listsAndItems = append(listsAndItems, lai)
  }

  return nil
}

func InitServer() {
  gin.ForceConsoleColor()

  r := gin.Default()

  r.LoadHTMLGlob("/var/lib/flimsy/templates/*.tmpl")

  r.Static("/static", "/var/lib/flimsy/static")
  r.Static("/data/icons", "/data/icons")
  r.Static("/data/backgrounds", "/data/backgrounds")
  r.GET("/", func(c *gin.Context) { GET_root(c) })
  r.GET("/config", func(c *gin.Context) { GET_config(c) })
  r.GET("/onlineStatus", func(c *gin.Context) { GET_onlineStatus(c) })

  r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func main() {
  err := InitDB(); if err != nil {
    fmt.Fprintln(os.Stderr, fmt.Sprintf("Failed to init DB: %s", err))
    return 
  }
  defer db.Close()

  InitServer()
}
