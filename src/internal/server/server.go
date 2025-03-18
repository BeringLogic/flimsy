package server

import (
	"bytes"
  "crypto/tls"
	"encoding/json"
	"fmt"
	"html/template"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"

	"github.com/dustin/go-humanize"
	"github.com/BeringLogic/palette-extractor"

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
  flimsyServer.router.HandleFunc("GET /onlineStatus/{id}", flimsyServer.GET_onlineStatus)
  flimsyServer.router.HandleFunc("GET /systemInfo", flimsyServer.GET_systemInfo)
  flimsyServer.router.HandleFunc("GET /weather", flimsyServer.GET_weather)
  flimsyServer.router.HandleFunc("GET /login", flimsyServer.GET_login)
  flimsyServer.router.HandleFunc("POST /login", flimsyServer.POST_login)
  flimsyServer.router.HandleFunc("GET /logout", flimsyServer.GET_logout)

  adminRouter := http.NewServeMux()
  adminRouter.HandleFunc("GET /config", flimsyServer.GET_config)
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

  adminMiddlewareStack := middleware.CreateStack(
    middleware.MustBeAuthenticated,
  )
  flimsyServer.router.Handle("/", adminMiddlewareStack(adminRouter))

  flimsyServer.middlewareStack = middleware.CreateStack(
    middleware.Logging(flimsyServer.log, false),
    middleware.IsAuthenticated(flimsyServer.storage),
  )

  return flimsyServer
}

func (flimsyServer *FlimsyServer) Start(host string, port string) error {
  flimsyServer.log.Printf("Starting server on %s:%s", host, port)
  return http.ListenAndServe(host + ":" + port, flimsyServer.middlewareStack(flimsyServer.router))
}

func (flimsyServer *FlimsyServer) error(w http.ResponseWriter, statusCode int, errorMessage string) {
  flimsyServer.log.Print(errorMessage)
  flimsyServer.executeTemplate("error.tmpl", &w, map[string]string{
    "header": http.StatusText(statusCode),
    "error": errorMessage,
  })
  w.WriteHeader(statusCode)
}

func (flimsyServer *FlimsyServer) executeTemplate(templateName string, w *http.ResponseWriter, data any) error {
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
  data := map[string]any{
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
    flimsyServer.error(w, http.StatusNotFound, "Item not found")
    return
  }

  if item.Check_url == "" {
    return
  }

  transport := &http.Transport{
		TLSClientConfig: &tls.Config{
    },
	}
  if item.Skip_certificate_verification == 1 {
    transport.TLSClientConfig.InsecureSkipVerify = true
  }
  client := &http.Client{
    Transport: transport,
    Timeout: 5 * time.Second,
  }

  resp, err := client.Get(item.Check_url); if err != nil {
    flimsyServer.executeTemplate("onlineStatus.tmpl", &w, map[string]string{
      "class" : "offline",
      "color" : "red",
      "title" : err.Error(),
      "Id" : idString,
    })
    return
  }

  if resp.StatusCode > 399 && resp.StatusCode != http.StatusUnauthorized && resp.StatusCode != http.StatusForbidden {
    flimsyServer.executeTemplate("onlineStatus.tmpl", &w, map[string]string{
      "class" : "offline",
      "color" : "red",
      "title" : fmt.Sprintf("%d %s", resp.StatusCode, http.StatusText(resp.StatusCode)),
      "Id" : idString,
    })
    return
  }

  resp.Body.Close()

  flimsyServer.executeTemplate("onlineStatus.tmpl", &w, map[string]string{
    "class" : "online",
    "color" : "green",
    "title" : "Online",
    "Id" : idString,
  })
}

func (flimsyServer *FlimsyServer) GET_systemInfo(w http.ResponseWriter, r *http.Request) {
  flimsyServer.executeTemplate("systemInfo.tmpl", &w, systemInfo.GetSystemInfo(flimsyServer.storage.Config))
}

func (flimsyServer *FlimsyServer) GET_weather(w http.ResponseWriter, r *http.Request) {
  u := "https://api.openweathermap.org/data/2.5/weather"
  u += "?q=" + url.QueryEscape(utils.GetEnv("FLIMSY_WEATHER_LOCATION", "New York"))
  u += "&units=" + url.QueryEscape(utils.GetEnv("FLIMSY_WEATHER_UNITS", "standard"))
  u += "&lang=" + url.QueryEscape(utils.GetEnv("FLIMSY_WEATHER_LANGUAGE", "en"))
  u += "&appid=" + url.QueryEscape(utils.GetEnv("FLIMSY_WEATHER_API_KEY", ""))

  response, err := http.Get(u); if err != nil {
    flimsyServer.error(w, http.StatusInternalServerError, err.Error())
    return
  }

  if response.StatusCode != 200 {
    flimsyServer.error(w, http.StatusInternalServerError, "Failed to fetch weather")
    return
  }

  defer response.Body.Close()

  bytes, err := io.ReadAll(response.Body); if err != nil {
    flimsyServer.error(w, http.StatusInternalServerError, err.Error())
    return
  }

  var data map[string]any
  err = json.Unmarshal(bytes, &data); if err != nil {
    flimsyServer.error(w, http.StatusInternalServerError, err.Error())
    return
  }

  icon := data["weather"].([]any)[0].(map[string]any)["icon"].(string)
  iconString := fmt.Sprintf("https://openweathermap.org/img/wn/%s@2x.png", icon)

  temp := data["main"].(map[string]any)["temp"].(float64)
  tempString := ""

  switch utils.GetEnv("FLIMSY_WEATHER_UNITS", "standard") {
    case "standard":
      tempString = fmt.Sprintf("%.0f °K", temp)
    case "metric":
      tempString = fmt.Sprintf("%.0f °C", temp)
    case "imperial":
      tempString = fmt.Sprintf("%.0f °F", temp)
  }

  flimsyServer.executeTemplate("weather.tmpl", &w, map[string]any{
    "icon" : iconString,
    "description" : data["weather"].([]any)[0].(map[string]any)["description"].(string),
    "location" : data["name"].(string) + ", " + data["sys"].(map[string]any)["country"].(string),
    "temp" : tempString,
  })
}

func (flimsyServer *FlimsyServer) logUserIn(w http.ResponseWriter, r *http.Request) {
  session, err := flimsyServer.storage.CreateNewSession(); if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  sessionCookie := http.Cookie{
    Name: "session_token",
    Value: session.Token,
    HttpOnly: true,
    Expires: session.ExpiresAt,
    SameSite: http.SameSiteStrictMode,
  }

  if r.URL.Scheme == "https" {
    sessionCookie.Secure = true
  }

  http.SetCookie(w, &sessionCookie)

  flimsyServer.log.Print("User logged in")
}

func (flimsyServer *FlimsyServer) GET_login(w http.ResponseWriter, r *http.Request) {
  if utils.GetEnv("FLIMSY_USERNAME", "") == "" && utils.GetEnv("FLIMSY_PASSWORD", "") == "" {
    flimsyServer.log.Print("Authentication is disabled. You can enable it by setting the environment variables FLIMSY_USERNAME and FLIMSY_PASSWORD.")
    flimsyServer.logUserIn(w, r)
    session_message = "Authentication is disabled. You can enable it by setting the environment variables FLIMSY_USERNAME and FLIMSY_PASSWORD.\n\nYou are now logged in!\n\n- Click on the gear button to customize the appearance\n- Click on items and lists to edit them\n- Drag & drop to reorder."
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

  flimsyServer.logUserIn(w, r)
  session_message = "You are now logged in!\n- Click on the gear button to customize the appearance\n- Click on items and lists to edit them\n- Drag & drop to reorder."
  w.Header().Set("HX-Location", "/")
}

func (flimsyServer *FlimsyServer) GET_logout(w http.ResponseWriter, r *http.Request) {
  currentSessionCookie, err := r.Cookie("session_token"); if err != nil {
    http.Redirect(w, r, "/", http.StatusSeeOther)
    return
  }

  if err := flimsyServer.storage.DeleteSession(currentSessionCookie.Value); err != nil {
    flimsyServer.error(w, http.StatusInternalServerError, err.Error())
    return
  }

  newSessionCookie := http.Cookie{
    Name: "session_token",
    Value: "",
    HttpOnly: true,
    Expires: time.Now(),
    SameSite: http.SameSiteStrictMode,
  }

  if r.URL.Scheme == "https" {
    newSessionCookie.Secure = true
  }

  http.SetCookie(w, &newSessionCookie)

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

  flimsyServer.executeTemplate("configDialog.tmpl", &w, map[string]any{
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
    var uploadError error
    var uploadedBackgroundSize uint64
    flimsyServer.storage.Config.Background_image, uploadedBackgroundSize, uploadError = saveUploadedBackground(r); if uploadError != nil {
      flimsyServer.error(w, http.StatusInternalServerError, uploadError.Error())
      return
    }
    flimsyServer.log.Printf("Background image uploaded: %s (%s)", flimsyServer.storage.Config.Background_image, humanize.Bytes(uploadedBackgroundSize))
  case "keep":
    flimsyServer.storage.Config.Background_image = r.FormValue("background_image")
  case "none":
    flimsyServer.storage.Config.Background_image = "";
  }

  switch Color_type {
  case "autodetect":
    var getColorsError error
    flimsyServer.storage.Config.Color_background = "black"
    flimsyServer.storage.Config.Color_foreground,
    flimsyServer.storage.Config.Color_items,
    flimsyServer.storage.Config.Color_borders,
    getColorsError = getColorsFromBackground(flimsyServer.storage.Config.Background_image); if getColorsError != nil {
      flimsyServer.error(w, http.StatusInternalServerError, getColorsError.Error())
      return
    }
    flimsyServer.log.Printf("Autodetected colors from background image: %s, %s, %s", flimsyServer.storage.Config.Color_foreground, flimsyServer.storage.Config.Color_items, flimsyServer.storage.Config.Color_borders)
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

  flimsyServer.log.Print("Config saved")

  http.Redirect(w, r, "/", http.StatusSeeOther)
}

func saveUploadedBackground(r *http.Request) (string, uint64, error) {
  err := r.ParseMultipartForm(100 << 20); if err != nil {
    return "", 0, fmt.Errorf("ParseMultipartForm error: %w", err)
  }

  file, header, err := r.FormFile("background_file"); if err != nil {
    return "", 0, fmt.Errorf("FormFile error: %w", err)
  }
  defer file.Close()

  bytes, err := io.ReadAll(file); if err != nil {
    return "", 0, fmt.Errorf("ReadAll error: %w", err)
  }

  outputFileName := "/data/backgrounds/" + strings.ReplaceAll(header.Filename, "/", "_")

  err = os.WriteFile(outputFileName, bytes, 0644); if err != nil {
    return "", 0, fmt.Errorf("WriteFile error: %w", err)
  }

  return path.Base(outputFileName), uint64(header.Size), nil
}

func getColorsFromBackground(backgroundImageName string) (string, string, string, error) {
  backgroundImage, err := os.Open("/data/backgrounds/" + backgroundImageName); if err != nil {
    return "", "", "", fmt.Errorf("Open error: %w", err)
  }
  defer backgroundImage.Close()

  img, _, err := image.Decode(backgroundImage); if err != nil {
    return "", "", "", fmt.Errorf("Decode error: %w", err)
  }

  extractor, err := extractor.NewExtractor(img, 1); if err != nil {
    return "", "", "", fmt.Errorf("NewExtractor error: %w", err)
  }

  palette := extractor.GetPalette(3)

  foregroundColor := fmt.Sprintf("#%02x%02x%02x", palette[2][0], palette[2][1], palette[2][2])
  itemsColor := fmt.Sprintf("#%02x%02x%02x", palette[0][0], palette[0][1], palette[0][2])
  bordersColor := fmt.Sprintf("#%02x%02x%02x", palette[1][0], palette[1][1], palette[1][2])

  return foregroundColor, itemsColor, bordersColor, nil
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

  numberOfColsString := r.FormValue("number_of_cols")
  if numberOfColsString == "" {
    flimsyServer.error(w, http.StatusBadRequest, "Missing number of cols")
    return
  }

  numberOfCols, err := strconv.Atoi(numberOfColsString); if err != nil {
    flimsyServer.error(w, http.StatusBadRequest, err.Error())
    return
  }

  lai, err := flimsyServer.storage.AddList(title, numberOfCols); if err != nil {
    flimsyServer.error(w, http.StatusInternalServerError, err.Error())
    return
  }

  flimsyServer.log.Print("List added: ", title)

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

  flimsyServer.log.Print("List saved: ", list.Title)

  list.Title = r.FormValue("title")

  Number_of_cols_string := r.FormValue("number_of_cols")
  list.Number_of_cols, err = strconv.Atoi(Number_of_cols_string); if err != nil {
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

  flimsyServer.log.Print("List deleted")
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

  flimsyServer.log.Print("Lists reordered")
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

  if err := icons.DownloadIcon(icon); err != nil {
    flimsyServer.error(w, http.StatusInternalServerError, err.Error())
    return
  }
  
  var skipCertificateVerification int
  skipCertificateVerificationString := r.FormValue("skip_certificate_verification")
  if skipCertificateVerificationString == "1" {
    skipCertificateVerification = 1
  }

  check_url := r.FormValue("check_url")

  item, err := flimsyServer.storage.AddItem(listId, title, url, icon, skipCertificateVerification, check_url); if err != nil {
    flimsyServer.error(w, http.StatusInternalServerError, err.Error())
    return
  }

  flimsyServer.log.Print("Item added: ", title)

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
  
  skipCertificateVerificationString := r.FormValue("skip_certificate_verification")
  if skipCertificateVerificationString == "1" {
    item.Skip_certificate_verification = 1
  }

  item.Check_url = r.FormValue("check_url")

  if err := icons.DownloadIcon(item.Icon); err != nil {
    flimsyServer.error(w, http.StatusInternalServerError, err.Error())
    return
  }

  if err := flimsyServer.storage.SaveItem(item); err != nil {
    flimsyServer.error(w, http.StatusInternalServerError, err.Error())
    return
  }

  flimsyServer.log.Print("Item updated: ", item.Title)

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

  flimsyServer.log.Print("Item deleted")
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

  flimsyServer.log.Print("Items reordered")
}
