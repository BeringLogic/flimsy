package main

import (
  "io"
  "os"
  "fmt"
  "bytes"
  "net/http"
  "html/template"
  "encoding/json"

  "github.com/BeringLogic/flimsy/internal/db"
  "github.com/BeringLogic/flimsy/internal/logger"
  "github.com/BeringLogic/flimsy/internal/systemInfo"
  "github.com/BeringLogic/flimsy/internal/middleware"
)

type listAndItems struct {
  List *db.List
  Items []*db.Item
}

var config *db.Config
var lists *[]db.List
var items *[]db.Item
var listsAndItems []listAndItems
var templates *template.Template

func getEnv(key, def string) string {
  if value, exists := os.LookupEnv(key); exists {
    return value
  } else {
    return def
  }
}

func ExecuteTemplate(templateName string, w *http.ResponseWriter, data interface{}) error {
  buffer := &bytes.Buffer{};

  err := templates.ExecuteTemplate(buffer, templateName, data); if err != nil {
    logger.Print(err.Error())
    templates.ExecuteTemplate(*w, "500.tmpl", nil)
  } else {
    buffer.WriteTo(*w)
  }

  return err
}

func GET_root(w http.ResponseWriter, r *http.Request) {
  data := map[string]interface{}{
    "auth_disabled" : true,
    "IsLoggedIn" : false,
    "session_message" : "", 
    "FLIMSY_WEATHER_API_KEY" : getEnv("FLIMSY_WEATHER_API_KEY", ""),
    "FLIMSY_WEATHER_LOCATION" : getEnv("FLIMSY_WEATHER_LOCATION", "New York"),
    "FLIMSY_WEATHER_UNITS" : getEnv("FLIMSY_WEATHER_UNITS", "standard"),
    "FLIMSY_WEATHER_LANGUAGE" : getEnv("FLIMSY_WEATHER_LANGUAGE", "en"),
    "config" : config,
    "listsAndItems" : listsAndItems,
  }

  ExecuteTemplate("index.tmpl", &w, data)
}

func GET_config(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(config);
}

func GET_onlineStatus(w http.ResponseWriter, r *http.Request) {
  href := r.URL.Query().Get("href")
  w.Header().Set("Content-Type", "application/json")
  
  resp, err := http.Get(href); if err != nil {
    json.NewEncoder(w).Encode(map[string]interface{}{
      "online" : false,
      "error" : err.Error(),
    })
    return
  }

  _, _ = io.ReadAll(resp.Body);
  resp.Body.Close()

  json.NewEncoder(w).Encode(map[string]interface{}{
    "online" : true,
  })
}

func GET_systemInfo(w http.ResponseWriter, r *http.Request) {
  ExecuteTemplate("systemInfo.tmpl", &w, systemInfo.GetSystemInfo(config))
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

func InitServer() error {
  templates = template.Must(template.ParseGlob("/var/lib/flimsy/templates/*.tmpl"))

  router := http.NewServeMux()

  router.Handle("GET /static/", http.StripPrefix("/static", http.FileServer(http.Dir("/var/lib/flimsy/static"))))
  router.Handle("GET /data/icons/", http.StripPrefix("/data/icons", http.FileServer(http.Dir("/data/icons"))))
  router.Handle("GET /data/backgrounds/", http.StripPrefix("/data/backgrounds", http.FileServer(http.Dir("/data/backgrounds"))))
  router.HandleFunc("GET /{$}", GET_root)
  router.HandleFunc("GET /config", GET_config)
  router.HandleFunc("GET /onlineStatus", GET_onlineStatus)
  router.HandleFunc("GET /systemInfo", GET_systemInfo)

  middlewareStack := middleware.CreateStack(
    middleware.Logging,
  )

  return http.ListenAndServe(":8080", middlewareStack(router))
}

func main() {
  logger.Init()

  if err := InitDB(); err != nil {
    logger.Fatalf("Failed to init DB: %s", err.Error())
  }
  defer db.Close()

  logger.Print("Starting server on port 0.0.0.0:8080")
  if err := InitServer(); err != nil {
    logger.Fatalf("Failed to init server: %s", err.Error())
  }
}
