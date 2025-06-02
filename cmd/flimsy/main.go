package main


import (
  "github.com/BeringLogic/flimsy/internal/logger"
  "github.com/BeringLogic/flimsy/internal/server"
  "github.com/BeringLogic/flimsy/internal/storage"
  "github.com/BeringLogic/flimsy/internal/utils"
)


type Application struct {
  log *logger.FlimsyLogger
  storage *storage.FlimsyStorage
  server *server.FlimsyServer
}


var app *Application


func main() {
  app = &Application{}
  app.log = logger.CreateNew()
  app.storage = storage.CreateNew()
  app.server = server.CreateNew(app.log, app.storage)

  if err := app.storage.Init(); err != nil {
    app.log.Fatalf("Failed to init DB: %s", err.Error())
  }
  defer app.storage.Close()

  host := utils.GetEnv("FLIMSY_HOST", "0.0.0.0")
  port := utils.GetEnv("FLIMSY_PORT", "8080")
  if err := app.server.Start(host, port); err != nil {
    app.log.Fatalf("Failed to start server: %s", err.Error())
  }
}
