package db

import (
  "database/sql"
  _ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

type Config struct {
  ID int
  icon string
  title string
  background_image string
  color_background string
  color_foreground string
  color_items string
  color_borders string
  cpu_temp_sensor string
  show_free_ram int
  show_free_swap int
  show_public_ip int
  show_free_space int
}

func Open() (*sql.DB, error) {
  db, err := sql.Open("sqlite3", "/data/flimsy.db?_busy_timeout=5000&_foreign_keys=ON&_journal_mode=WAL"); if err != nil {
    return nil, err
  }
  return db, nil
}

// Try to get the config and seed the DB if it doesn't exist
func GetConfig() (*Config, error) {
  rows, err := db.Query("SELECT * FROM config WHERE id = 1"); if err != nil {
    // If there is an error it's probably because the table doesn't exist so seed the DB
    err = seed(); if err != nil {
      return nil, err
    }
    // Try again
    rows, err = db.Query("SELECT * FROM config WHERE id = 1"); if err != nil {
      return nil, err
    }
  }

  config := Config{}
  err = rows.Scan(
      &config.ID,
      &config.icon,
      &config.title,
      &config.background_image,
      &config.color_background,
      &config.color_foreground,
      &config.color_items,
      &config.color_borders,
      &config.cpu_temp_sensor,
      &config.show_free_ram,
      &config.show_free_swap,
      &config.show_public_ip,
      &config.show_free_space); if err != nil {
    return nil, err
  }

  return &config, nil
}

func seed() error {
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
    _, err := db.Exec(query); if err != nil {
      return err
    }
  }
  
  return nil 
}

func Close() {
  db.Close()
}
