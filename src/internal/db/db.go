package db

import (
	"embed"
	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/mattn/go-sqlite3"

	"github.com/BeringLogic/flimsy/internal/utils"
)


type FlimsyDB struct {
  sqlDb *sql.DB
}


//go:embed migrations
var migrations embed.FS


func CreateNew() *FlimsyDB {
  return &FlimsyDB {}
}

func (flimsyDB *FlimsyDB) Open() (error) {
  var openError error
  if flimsyDB.sqlDb, openError = sql.Open("sqlite3", "/data/flimsy.db?_busy_timeout=5000&_foreign_keys=ON&_journal_mode=WAL"); openError != nil {
    return openError
  }

  // check if this is a first init
  var tableCount int
  if err := flimsyDB.sqlDb.QueryRow("SELECT count(*) FROM sqlite_master WHERE type = 'table' AND name = 'config'").Scan(&tableCount); err != nil {
    return err
  }
  firstInit := tableCount == 0

  // apply migrations
	if err := flimsyDB.Migrate(); err != nil {
    return err
	}

  // If this is the first init, set the cpu temp sensor
  if firstInit {
    _, err := flimsyDB.sqlDb.Exec("UPDATE config SET cpu_temp_sensor = ? WHERE id = 1", utils.GetEnv("FLIMSY_CPU_TEMP_SENSOR", "")); if err != nil {
      return err
    }
  }

  return nil 
}

func (flimsyDB *FlimsyDB) Migrate() error {
  source, err := iofs.New(migrations, "migrations"); if err != nil {
    return err
  }

  dbDriver, err := sqlite3.WithInstance(flimsyDB.sqlDb, &sqlite3.Config{}); if err != nil {
    return err
  }

  m, err := migrate.NewWithInstance("iofs", source, "sqlite3", dbDriver); if err != nil {
    return err
  }

  if err := m.Up(); err != nil && err != migrate.ErrNoChange {
    return err
  }

  return nil
}

func (flimsyDB *FlimsyDB) Close() {
  flimsyDB.sqlDb.Close()
}
