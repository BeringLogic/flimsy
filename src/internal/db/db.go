package db

import (
  "database/sql"
  _ "github.com/mattn/go-sqlite3"
)

var sqlDb *sql.DB

func Open() (error) {
  var err error
  sqlDb, err = sql.Open("sqlite3", "/data/flimsy.db?_busy_timeout=5000&_foreign_keys=ON&_journal_mode=WAL"); if err != nil {
    return err
  }
  return nil
}

func Seed() error {
  queries := []string {
    `CREATE TABLE config (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      icon TEXT NULL,
      title TEXT NULL,
      background_image TEXT NULL,
      color_background TEXT NOT NULL,
      color_foreground TEXT NOT NULL,
      color_items TEXT NOT NULL,
      color_borders TEXT NOT NULL,
      cpu_temp_sensor TEXT,
      show_free_ram INTEGER,
      show_free_swap INTEGER,
      show_public_ip INTEGER,
      show_free_space INTEGER
    );`,
    `INSERT INTO config (
      icon,
      title,
      color_background,
      color_foreground,
      color_items,
      color_borders,
      cpu_temp_sensor,
      show_free_ram,
      show_free_swap,
      show_public_ip, 
      show_free_space
    )
    VALUES
      NULL,
      'Flimsy Home Page',
      '#1e1e2e',
      '#cdd6f4',
      '#11111b',
      '#6c7086',
      '$cpuTempSensor',
      true,
      true,
      true,
      true
    );`,
    `CREATE TABLE list (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      title TEXT NOT NULL,
      number_of_rows INTEGER NOT NULL,
      position INTEGER NOT NULL
    );`,
    `CREATE TABLE item (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      list_id INTEGER NOT NULL,
      title TEXT NOT NULL,
      href TEXT NOT NULL,
      icon TEXT NOT NULL,
      position INTEGER NOT  NULL,
      FOREIGN KEY(list_id) REFERENCES lists(id)
    );`,
  }

  for _, query := range queries {
    _, err := sqlDb.Exec(query); if err != nil {
      return err
    }
  }
 
  return nil 
}

func Close() {
  sqlDb.Close()
}
