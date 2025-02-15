package main


import (
  "github.com/BeringLogic/flimsy/internal/logger"
  "github.com/BeringLogic/flimsy/internal/server"
  "github.com/BeringLogic/flimsy/internal/storage"
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

  app.log.Print("Starting server on port 0.0.0.0:8080")
  if err := app.server.Start(); err != nil {
    app.log.Fatalf("Failed to start server: %s", err.Error())
  }
}
