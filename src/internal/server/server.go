package server


import (
  "io"
  "time"
  "bytes"
  "net/http"
  "html/template"
  "encoding/json"

  "github.com/BeringLogic/flimsy/internal/auth"
  "github.com/BeringLogic/flimsy/internal/utils"
  "github.com/BeringLogic/flimsy/internal/logger"
  "github.com/BeringLogic/flimsy/internal/storage"
  "github.com/BeringLogic/flimsy/internal/systemInfo"
  "github.com/BeringLogic/flimsy/internal/middleware"
)


type FlimsyServer struct {
  log *logger.FlimsyLogger
  storage *storage.FlimsyStorage
  templates *template.Template
  router *http.ServeMux
  middlewareStack middleware.Middleware
}


var session_message string


func CreateNew(log *logger.FlimsyLogger, storage *storage.FlimsyStorage) *FlimsyServer {
  flimsyServer := &FlimsyServer{} 

  flimsyServer.log = log
  flimsyServer.storage = storage

  flimsyServer.templates = template.Must(template.ParseGlob("/var/lib/flimsy/templates/*.tmpl"))

  flimsyServer.router = http.NewServeMux()
  flimsyServer.router.Handle("GET /static/", http.StripPrefix("/static", http.FileServer(http.Dir("/var/lib/flimsy/static"))))
  flimsyServer.router.Handle("GET /data/icons/", http.StripPrefix("/data/icons", http.FileServer(http.Dir("/data/icons"))))
  flimsyServer.router.Handle("GET /data/backgrounds/", http.StripPrefix("/data/backgrounds", http.FileServer(http.Dir("/data/backgrounds"))))
  flimsyServer.router.HandleFunc("GET /{$}", flimsyServer.GET_root)
  flimsyServer.router.HandleFunc("GET /config", flimsyServer.GET_config)
  flimsyServer.router.HandleFunc("GET /onlineStatus", flimsyServer.GET_onlineStatus)
  flimsyServer.router.HandleFunc("GET /systemInfo", flimsyServer.GET_systemInfo)
  flimsyServer.router.HandleFunc("GET /login", flimsyServer.GET_login)
  flimsyServer.router.HandleFunc("POST /login", flimsyServer.POST_login)
  flimsyServer.router.HandleFunc("GET /logout", flimsyServer.GET_logout)

  wrappedLogger := middleware.Logging(flimsyServer.log)

  flimsyServer.middlewareStack = middleware.CreateStack(
    wrappedLogger,
    middleware.IsAuthenticated,
  )

  return flimsyServer
}

func (flimsyServer *FlimsyServer) Start() error {
  return http.ListenAndServe(":8080", flimsyServer.middlewareStack(flimsyServer.router))
}

func (flimsyServer *FlimsyServer) executeTemplate(templateName string, w *http.ResponseWriter, data interface{}) error {
  buffer := &bytes.Buffer{};

  if err := flimsyServer.templates.ExecuteTemplate(buffer, templateName, data); err != nil {
    flimsyServer.log.Print(err.Error())
    flimsyServer.templates.ExecuteTemplate(*w, "500.tmpl", nil)
    return err
  } else {
    buffer.WriteTo(*w)
    return nil
  }
}

func (flimsyServer *FlimsyServer) GET_root(w http.ResponseWriter, r *http.Request) {
  data := map[string]interface{}{
    "auth_disabled" : true,
    "IsLoggedIn" : r.Context().Value(middleware.IsAuthenticatedContextKey).(bool),
    "session_message" : session_message, 
    "FLIMSY_WEATHER_API_KEY" : utils.GetEnv("FLIMSY_WEATHER_API_KEY", ""),
    "FLIMSY_WEATHER_LOCATION" : utils.GetEnv("FLIMSY_WEATHER_LOCATION", "New York"),
    "FLIMSY_WEATHER_UNITS" : utils.GetEnv("FLIMSY_WEATHER_UNITS", "standard"),
    "FLIMSY_WEATHER_LANGUAGE" : utils.GetEnv("FLIMSY_WEATHER_LANGUAGE", "en"),
    "config" : flimsyServer.storage.Config,
    "listsAndItems" : flimsyServer.storage.ListsAndItems,
  }

  session_message = ""

  flimsyServer.executeTemplate("index.tmpl", &w, data)
}

func (flimsyServer *FlimsyServer) GET_config(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(flimsyServer.storage.Config);
}

func (flimsyServer *FlimsyServer) GET_onlineStatus(w http.ResponseWriter, r *http.Request) {
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

func (flimsyServer *FlimsyServer) GET_systemInfo(w http.ResponseWriter, r *http.Request) {
  flimsyServer.executeTemplate("systemInfo.tmpl", &w, systemInfo.GetSystemInfo(flimsyServer.storage.Config))
}

func (flimsyServer *FlimsyServer) logUserIn(w http.ResponseWriter) {
  session_token, csrf_token, err := auth.GenerateTokens(); if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  http.SetCookie(w, &http.Cookie{
    Name: "session_token",
    Value: session_token,
    HttpOnly: true,
    Expires: time.Now().Add(time.Hour * 3),
    SameSite: http.SameSiteLaxMode,
  })

  http.SetCookie(w, &http.Cookie{
    Name: "csrf_token",
    Value: csrf_token,
    HttpOnly: false,
    Expires: time.Now().Add(time.Hour * 3),
    SameSite: http.SameSiteLaxMode,
  })

  flimsyServer.log.Print("User logged in")
}

func (flimsyServer *FlimsyServer) GET_login(w http.ResponseWriter, r *http.Request) {
  if utils.GetEnv("FLIMSY_USERNAME", "") == "" && utils.GetEnv("FLIMSY_PASSWORD", "") == "" {
    flimsyServer.logUserIn(w)
    session_message = "Authentication is disabled. You can enable it by setting the environment variables FLIMSY_USERNAME and FLIMSY_PASSWORD.\n\n- Click on the gear button to customize the appearance\n- Click on items and lists to edit them\n- Drag & drop to reorder."
    http.Redirect(w, r, "/", http.StatusSeeOther)
  } else {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
      "success" : false,
      "error" : "Authentication is enabled. You must POST credentials to authenticate.",
    })
    w.WriteHeader(http.StatusForbidden)
  }
}

func (flimsyServer *FlimsyServer) POST_login(w http.ResponseWriter, r *http.Request) {
  username := r.FormValue("username")
  password := r.FormValue("password")
  if username != utils.GetEnv("FLIMSY_USERNAME", "") || password != utils.GetEnv("FLIMSY_PASSWORD", "") {
    http.Error(w, "Forbidden", http.StatusForbidden)
    return
  }

  flimsyServer.logUserIn(w)
  session_message = "You are now logged in!\n- Click on the gear button to customize the appearance\n- Click on items and lists to edit them\n- Drag & drop to reorder."
}

func (flimsyServer *FlimsyServer) GET_logout(w http.ResponseWriter, r *http.Request) {
  http.SetCookie(w, &http.Cookie{
    Name: "session_token",
    Value: "",
    HttpOnly: true,
    Expires: time.Now(),
    SameSite: http.SameSiteLaxMode,
  })
  http.SetCookie(w, &http.Cookie{
    Name: "csrf_token",
    Value: "",
    HttpOnly: false,
    Expires: time.Now(),
    SameSite: http.SameSiteLaxMode,
  })
  http.Redirect(w, r, "/", http.StatusSeeOther)

  flimsyServer.log.Print("User logged out")
}
