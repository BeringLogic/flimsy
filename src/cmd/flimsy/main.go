package main

import (
  "io"
  "os"
  "fmt"
  "net/http"
  "html/template"
  "encoding/json"

  "github.com/BeringLogic/flimsy/internal/db"
  "github.com/BeringLogic/flimsy/internal/systemInfo"
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
    value, exists := os.LookupEnv(key)
    if !exists {
        value = def
    }
    return value
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

  templates.ExecuteTemplate(w, "index.tmpl", data)
}

func GET_config(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  je := json.NewEncoder(w)
  je.Encode(config);
}

func GET_onlineStatus(w http.ResponseWriter, r *http.Request) {
  href := r.URL.Query().Get("href")
  je := json.NewEncoder(w)

  w.Header().Set("Content-Type", "application/json")
  
  resp, err := http.Get(href); if err != nil {
    je.Encode(map[string]interface{}{
      "online" : false,
      "error" : err.Error(),
    })
    return
  }

  _, _ = io.ReadAll(resp.Body);
  resp.Body.Close()

  je.Encode(map[string]interface{}{
    "online" : true,
  })
}

func GET_systemInfo(w http.ResponseWriter, r *http.Request) {
  templates.ExecuteTemplate(w, "systemInfo.tmpl", systemInfo.GetSystemInfo(config))
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

  return http.ListenAndServe(":8080", router)
}

func main() {
  if err := InitDB(); err != nil {
    fmt.Fprintln(os.Stderr, fmt.Sprintf("Failed to init DB: %s", err))
    return 
  }
  defer db.Close()

  fmt.Println("Starting server on port 0.0.0.0:8080")
  if err := InitServer(); err != nil {
    fmt.Fprintln(os.Stderr, fmt.Sprintf("Failed to init server: %s", err))
    return 
  }
}
