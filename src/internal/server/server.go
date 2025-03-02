package server

import (
	"bytes"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/BeringLogic/flimsy/internal/auth"
	"github.com/BeringLogic/flimsy/internal/icons"
	"github.com/BeringLogic/flimsy/internal/logger"
	"github.com/BeringLogic/flimsy/internal/middleware"
	"github.com/BeringLogic/flimsy/internal/storage"
	"github.com/BeringLogic/flimsy/internal/systemInfo"
	"github.com/BeringLogic/flimsy/internal/utils"
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
  flimsyServer.router.HandleFunc("GET /style.css", flimsyServer.GET_style)
  flimsyServer.router.HandleFunc("GET /onlineStatus", flimsyServer.GET_onlineStatus)
  flimsyServer.router.HandleFunc("GET /systemInfo", flimsyServer.GET_systemInfo)
  flimsyServer.router.HandleFunc("GET /login", flimsyServer.GET_login)
  flimsyServer.router.HandleFunc("POST /login", flimsyServer.POST_login)
  flimsyServer.router.HandleFunc("GET /logout", flimsyServer.GET_logout)
  flimsyServer.router.HandleFunc("GET /config", flimsyServer.GET_config)

  adminRouter := http.NewServeMux()
  adminRouter.HandleFunc("POST /config", flimsyServer.POST_config)
  adminRouter.HandleFunc("PUT /list", flimsyServer.PUT_list)
  adminRouter.HandleFunc("GET /list/{id}", flimsyServer.GET_list)
  adminRouter.HandleFunc("PATCH /list/{id}", flimsyServer.PATCH_list)
  adminRouter.HandleFunc("DELETE /list/{id}", flimsyServer.DELETE_list)
  adminRouter.HandleFunc("POST /reorderLists", flimsyServer.POST_reorderLists)
  adminRouter.HandleFunc("PUT /item/{list_id}", flimsyServer.PUT_item)
  adminRouter.HandleFunc("GET /item/{id}", flimsyServer.GET_item)
  adminRouter.HandleFunc("PATCH /item/{id}", flimsyServer.PATCH_item)
  adminRouter.HandleFunc("DELETE /item/{id}", flimsyServer.DELETE_item)
  adminRouter.HandleFunc("POST /reorderItems", flimsyServer.POST_reorderItems)
  flimsyServer.router.Handle("/", middleware.MustBeAuthenticated(adminRouter))

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

func (flimsyServer *FlimsyServer) error(w http.ResponseWriter, statusCode int, errorMessage string) {
  flimsyServer.log.Print(errorMessage)
  flimsyServer.executeTemplate("error.tmpl", &w, map[string]string{
    "header": http.StatusText(statusCode),
    "error": errorMessage,
  })
  w.WriteHeader(statusCode)
}

func (flimsyServer *FlimsyServer) executeTemplate(templateName string, w *http.ResponseWriter, data interface{}) error {
  buffer := &bytes.Buffer{};

  if err := flimsyServer.templates.ExecuteTemplate(buffer, templateName, data); err != nil {
    flimsyServer.error(*w, http.StatusInternalServerError, err.Error())
    return err
  } else {
    buffer.WriteTo(*w)
    return nil
  }
}

func (flimsyServer *FlimsyServer) GET_root(w http.ResponseWriter, r *http.Request) {
  data := map[string]interface{}{
    "IsAuthDisabled" : utils.GetEnv("FLIMSY_USERNAME", "") == "" && utils.GetEnv("FLIMSY_PASSWORD", "") == "",
    "IsLoggedIn" : r.Context().Value(middleware.IsAuthenticatedContextKey).(bool),
    "session_message" : session_message, 
    "FLIMSY_WEATHER_API_KEY" : utils.GetEnv("FLIMSY_WEATHER_API_KEY", ""),
    "FLIMSY_WEATHER_LOCATION" : utils.GetEnv("FLIMSY_WEATHER_LOCATION", "New York"),
    "FLIMSY_WEATHER_UNITS" : utils.GetEnv("FLIMSY_WEATHER_UNITS", "standard"),
    "FLIMSY_WEATHER_LANGUAGE" : utils.GetEnv("FLIMSY_WEATHER_LANGUAGE", "en"),
    "config" : flimsyServer.storage.Config,
    "listsAndItems" : flimsyServer.storage.AllListsAndItems,
  }

  session_message = ""

  flimsyServer.executeTemplate("index.tmpl", &w, data)
}

func (flimsyServer *FlimsyServer) GET_style(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "text/css")
  flimsyServer.executeTemplate("style.tmpl", &w, flimsyServer.storage.Config)
}

func (flimsyServer *FlimsyServer) GET_onlineStatus(w http.ResponseWriter, r *http.Request) {
  url := r.URL.Query().Get("url")
  
  resp, err := http.Get(url); if err != nil {
    flimsyServer.executeTemplate("onlineStatus.tmpl", &w, map[string]string{
      "class" : "offline",
      "color" : "red",
      "title" : err.Error(),
    })
    return
  }

  if resp.StatusCode != 200 {
    flimsyServer.executeTemplate("onlineStatus.tmpl", &w, map[string]string{
      "class" : "offline",
      "color" : "red",
      "title" : err.Error(),
    })
    return
  }

  _, _ = io.ReadAll(resp.Body);
  resp.Body.Close()

  flimsyServer.executeTemplate("onlineStatus.tmpl", &w, map[string]string{
    "class" : "online",
    "color" : "green",
    "title" : "Online",
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
    flimsyServer.executeTemplate("loginDialog.tmpl", &w, nil)
  }
}

func (flimsyServer *FlimsyServer) POST_login(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")

  username := r.FormValue("username")
  password := r.FormValue("password")
  if username != utils.GetEnv("FLIMSY_USERNAME", "") || password != utils.GetEnv("FLIMSY_PASSWORD", "") {
    flimsyServer.log.Print("Invalid username or password")
    flimsyServer.executeTemplate("loginDialog.tmpl", &w, map[string]string{
      "error" : "Invalid username or password",
    })
    return
  }

  flimsyServer.logUserIn(w)
  session_message = "You are now logged in!\n- Click on the gear button to customize the appearance\n- Click on items and lists to edit them\n- Drag & drop to reorder."
  w.Header().Set("HX-Location", "/")
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

func (flimsyServer *FlimsyServer) GET_config(w http.ResponseWriter, r *http.Request) {
  backgrounds, err := utils.GetBackgrounds(); if err != nil {
    flimsyServer.log.Print(err.Error())
  }
  sensors, err := systemInfo.GetSensors(); if err != nil {
    flimsyServer.log.Print(err.Error())
  }

  flimsyServer.executeTemplate("configDialog.tmpl", &w, map[string]interface{}{
    "config" : flimsyServer.storage.Config,
    "backgrounds" : backgrounds,
    "sensors" : sensors,
  })
}

func (flimsyServer *FlimsyServer) POST_config(w http.ResponseWriter, r *http.Request) {
  flimsyServer.storage.Config.Icon = r.FormValue("icon")
  flimsyServer.storage.Config.Title = r.FormValue("title")
  Background_type := r.FormValue("background_type")
  Color_type := r.FormValue("color_type")
  flimsyServer.storage.Config.Color_background = r.FormValue("color_background")
  flimsyServer.storage.Config.Color_foreground = r.FormValue("color_foreground")
  flimsyServer.storage.Config.Color_items = r.FormValue("color_items")
  flimsyServer.storage.Config.Color_borders = r.FormValue("color_borders")
  flimsyServer.storage.Config.Cpu_temp_sensor = r.FormValue("cpu_temp_sensor")
  Show_free_ram := r.FormValue("show_free_ram")
  Show_free_swap := r.FormValue("show_free_swap")
  Show_public_ip := r.FormValue("show_public_ip")
  Show_free_space := r.FormValue("show_free_space")

  switch Background_type {
  case "upload":
    // flimsyServer.storage.Config.Background_image = $upBg->upload(); 
    flimsyServer.log.Print("TODO: Background image upload")
  case "keep":
    flimsyServer.storage.Config.Background_image = r.FormValue("background_image")
  case "none":
    flimsyServer.storage.Config.Background_image = "";
  }

  switch Color_type {
  case "autodetect":
    flimsyServer.log.Print("TODO: autodetect background colors")
  case "catppuccin_latte":
    flimsyServer.storage.Config.Color_background = "#eff1f5"
    flimsyServer.storage.Config.Color_foreground = "#4c4f69"
    flimsyServer.storage.Config.Color_items = "#dce0e8"
    flimsyServer.storage.Config.Color_borders = "#9ca0b0"
  case "catppuccin_mocha":
    flimsyServer.storage.Config.Color_background = "#1e1e2e"
    flimsyServer.storage.Config.Color_foreground = "#cdd6f4"
    flimsyServer.storage.Config.Color_items = "#11111b"
    flimsyServer.storage.Config.Color_borders = "#6c7086"
  case "manual":
    flimsyServer.storage.Config.Color_background = r.FormValue("color_background")
    flimsyServer.storage.Config.Color_foreground = r.FormValue("color_foreground")
    flimsyServer.storage.Config.Color_items = r.FormValue("color_items")
    flimsyServer.storage.Config.Color_borders = r.FormValue("color_borders")
  }

  if Show_free_ram == "" {
    flimsyServer.storage.Config.Show_free_ram = 0
  } else {
    flimsyServer.storage.Config.Show_free_ram = 1
  }
  if Show_free_swap == "" {
    flimsyServer.storage.Config.Show_free_swap = 0
  } else {
    flimsyServer.storage.Config.Show_free_swap = 1
  }
  if Show_public_ip == "" {
    flimsyServer.storage.Config.Show_public_ip = 0
  } else {
    flimsyServer.storage.Config.Show_public_ip = 1
  }
  if Show_free_space == "" {
    flimsyServer.storage.Config.Show_free_space = 0
  } else {
    flimsyServer.storage.Config.Show_free_space = 1
  }

  if err := flimsyServer.storage.SaveConfig(); err != nil {
    flimsyServer.error(w, http.StatusInternalServerError, err.Error())
    return
  }

  if flimsyServer.storage.Config.Icon != "" {
    err := icons.DownloadIcon(flimsyServer.storage.Config.Icon); if err != nil {
      flimsyServer.log.Print(err.Error())
    }
  }

  http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (flimsyServer *FlimsyServer) PUT_list(w http.ResponseWriter, r *http.Request) {
  confirm := r.FormValue("confirm")
  if confirm != "true" {
    flimsyServer.executeTemplate("addListDialog.tmpl", &w, nil)
    return
  }

  title := r.FormValue("title")
  if title == "" {
    flimsyServer.error(w, http.StatusBadRequest, "Missing title")
    return
  }

  numberOfRowsString := r.FormValue("number_of_rows")
  if numberOfRowsString == "" {
    flimsyServer.error(w, http.StatusBadRequest, "Missing number of rows")
    return
  }

  numberOfRows, err := strconv.Atoi(numberOfRowsString); if err != nil {
    flimsyServer.error(w, http.StatusBadRequest, err.Error())
    return
  }

  lai, err := flimsyServer.storage.AddList(title, numberOfRows); if err != nil {
    flimsyServer.error(w, http.StatusInternalServerError, err.Error())
    return
  }

  flimsyServer.executeTemplate("list.loggedin.tmpl", &w, lai)
}

func (flimsyServer *FlimsyServer) GET_list(w http.ResponseWriter, r *http.Request) {
  idString := r.PathValue("id")
  if idString == "" {
    flimsyServer.log.Print("Missing id")
    http.Error(w, "Missing id", http.StatusBadRequest)
    return
  }

  id, err := strconv.ParseInt(idString, 10, 64); if err != nil {
    flimsyServer.log.Print(err.Error())
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  list, err := flimsyServer.storage.GetList(id); if err != nil {
    flimsyServer.error(w, http.StatusNotFound, err.Error())
    return
  }

  flimsyServer.executeTemplate("editListDialog.tmpl", &w, list)
}

func (flimsyServer *FlimsyServer) PATCH_list(w http.ResponseWriter, r *http.Request) {
  idString := r.PathValue("id")
  if idString == "" {
    flimsyServer.error(w, http.StatusBadRequest, "Missing id")
    return
  }

  id, err := strconv.ParseInt(idString, 10, 64); if err != nil {
    flimsyServer.error(w, http.StatusBadRequest, err.Error())
    return
  }

  list, err := flimsyServer.storage.GetList(id); if err != nil {
    flimsyServer.error(w, http.StatusNotFound, err.Error())
    return
  }

  list.Title = r.FormValue("title")

  Number_of_rows_string := r.FormValue("number_of_rows")
  list.Number_of_rows, err = strconv.Atoi(Number_of_rows_string); if err != nil {
    flimsyServer.error(w, http.StatusBadRequest, err.Error())
    return
  }

  listAndItems, err := flimsyServer.storage.SaveList(list); if err != nil { 
    flimsyServer.error(w, http.StatusInternalServerError, err.Error())
    return
  }

  flimsyServer.executeTemplate("list.loggedin.tmpl", &w, listAndItems)
}

func (flimsyServer *FlimsyServer) DELETE_list(w http.ResponseWriter, r *http.Request) {
  idString := r.PathValue("id")
  if idString == "" {
    flimsyServer.error(w, http.StatusBadRequest, "Missing id")
    return
  }

  confirm := r.FormValue("confirm")
  if confirm != "true" {
    flimsyServer.executeTemplate("deleteListDialog.tmpl", &w, idString)
    return
  }

  id, err := strconv.ParseInt(idString, 10, 64); if err != nil {
    flimsyServer.error(w, http.StatusBadRequest, err.Error())
    return
  }

  if err := flimsyServer.storage.DeleteList(id); err != nil {
    flimsyServer.error(w, http.StatusInternalServerError, err.Error())
    return
  }
}

func (flimsyServer *FlimsyServer) POST_reorderLists(w http.ResponseWriter, r *http.Request) {
  if err := r.ParseForm(); err != nil {
    flimsyServer.error(w, http.StatusBadRequest, err.Error())
    return
  }

  list_id_strings, exists := r.Form["ids"]
  if !exists || list_id_strings == nil {
    flimsyServer.error(w, http.StatusBadRequest, "Missing ids")
    return
  }

  list_ids := make([]int64, len(list_id_strings))
  for i, list_id_string := range list_id_strings {
    list_id, err := strconv.ParseInt(list_id_string, 10, 64); if err != nil {
      flimsyServer.error(w, http.StatusBadRequest, err.Error())
      return
    }
    list_ids[i] = list_id
  }

  if err := flimsyServer.storage.ReorderLists(list_ids); err != nil {
    flimsyServer.error(w, http.StatusInternalServerError, err.Error())
    return
  }
}

func (flimsyServer *FlimsyServer) PUT_item(w http.ResponseWriter, r *http.Request) {
  listIdString := r.PathValue("list_id")
  if listIdString == "" {
    flimsyServer.error(w, http.StatusBadRequest, "Missing list id")
    return
  }

  confirm := r.FormValue("confirm")
  if confirm != "true" {
    flimsyServer.executeTemplate("addItemDialog.tmpl", &w, listIdString)
    return
  }

  listId, err := strconv.ParseInt(listIdString, 10, 64); if err != nil {
    flimsyServer.error(w, http.StatusBadRequest, err.Error())
    return
  }

  title := r.FormValue("title")
  if title == "" {
    flimsyServer.error(w, http.StatusBadRequest, "Missing title")
    return
  }

  url := r.FormValue("url")
  if url == "" {
    flimsyServer.error(w, http.StatusBadRequest, "Missing url")
    return
  }

  icon := r.FormValue("icon")
  if icon == "" {
    flimsyServer.error(w, http.StatusBadRequest, "Missing icon")
    return
  }

  item, err := flimsyServer.storage.AddItem(listId, title, url, icon); if err != nil {
    flimsyServer.error(w, http.StatusInternalServerError, err.Error())
    return
  }

  flimsyServer.executeTemplate("item.loggedin.tmpl", &w, item)
}

func (flimsyServer *FlimsyServer) GET_item(w http.ResponseWriter, r *http.Request) {
  idString := r.PathValue("id")
  if idString == "" {
    flimsyServer.error(w, http.StatusBadRequest, "Missing id")
    return
  }

  id, err := strconv.ParseInt(idString, 10, 64); if err != nil {
    flimsyServer.error(w, http.StatusBadRequest, err.Error())
    return
  }

  item, err := flimsyServer.storage.GetItem(id); if err != nil {
    flimsyServer.error(w, http.StatusNotFound, err.Error())
    return
  }

  flimsyServer.executeTemplate("editItemDialog.tmpl", &w, item)
}

func (flimsyServer *FlimsyServer) PATCH_item(w http.ResponseWriter, r *http.Request) {
  idString := r.PathValue("id")
  if idString == "" {
    flimsyServer.error(w, http.StatusBadRequest, "Missing id")
    return
  }

  id, err := strconv.ParseInt(idString, 10, 64); if err != nil {
    flimsyServer.error(w, http.StatusBadRequest, err.Error())
    return
  }

  item, err := flimsyServer.storage.GetItem(id); if err != nil {
    flimsyServer.error(w, http.StatusNotFound, err.Error())
    return
  }

  item.Title = r.FormValue("title")
  item.Url = r.FormValue("url")
  item.Icon = r.FormValue("icon")

  if err := flimsyServer.storage.SaveItem(item); err != nil {
    flimsyServer.error(w, http.StatusInternalServerError, err.Error())
    return
  }

  flimsyServer.executeTemplate("item.loggedin.tmpl", &w, item)
}

func (flimsyServer *FlimsyServer) DELETE_item(w http.ResponseWriter, r *http.Request) {
  idString := r.PathValue("id")
  if idString == "" {
    flimsyServer.error(w, http.StatusBadRequest, "Missing id")
    return
  }

  confirm := r.FormValue("confirm")

  if confirm != "true" {
    flimsyServer.executeTemplate("deleteItemDialog.tmpl", &w, idString)
    return
  }

  id, err := strconv.ParseInt(idString, 10, 64); if err != nil {
    flimsyServer.error(w, http.StatusBadRequest, err.Error())
    return
  }

  if err := flimsyServer.storage.DeleteItem(id); err != nil {
    flimsyServer.error(w, http.StatusInternalServerError, err.Error())
    return
  }
}

func (flimsyServer *FlimsyServer) POST_reorderItems(w http.ResponseWriter, r *http.Request) {
  if err := r.ParseForm(); err != nil {
    flimsyServer.error(w, http.StatusBadRequest, err.Error())
    return
  }

  list_id_string, exists := r.Form["list_id"]
  if !exists || len(list_id_string) < 1 || list_id_string[0] == "" {
    flimsyServer.error(w, http.StatusBadRequest, "Missing list_id")
    return
  }

  list_id, err := strconv.ParseInt(list_id_string[0], 10, 64); if err != nil {
    flimsyServer.error(w, http.StatusBadRequest, err.Error())
    return
  }

  item_id_strings, exists := r.Form["ids"]
  if !exists || item_id_strings == nil {
    flimsyServer.error(w, http.StatusBadRequest, "Missing ids")
    return
  }

  item_ids := make([]int64, len(item_id_strings))
  for i, item_id_string := range item_id_strings {
    item_id, err := strconv.ParseInt(item_id_string, 10, 64); if err != nil {
      flimsyServer.error(w, http.StatusBadRequest, err.Error())
      return
    }
    item_ids[i] = item_id
  }

  if err := flimsyServer.storage.ReorderItems(list_id, item_ids); err != nil {
    flimsyServer.error(w, http.StatusInternalServerError, err.Error())
    return
  }
}
