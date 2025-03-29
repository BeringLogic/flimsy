CREATE TABLE config (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  icon TEXT NOT NULL,
  title TEXT NOT NULL,
  background_image TEXT NOT NULL,
  color_background TEXT NOT NULL,
  color_foreground TEXT NOT NULL,
  color_items TEXT NOT NULL,
  color_borders TEXT NOT NULL,
  cpu_temp_sensor TEXT,
  show_free_ram INTEGER,
  show_free_swap INTEGER,
  show_public_ip INTEGER,
  show_free_space INTEGER
);
CREATE TABLE list (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  title TEXT NOT NULL,
  number_of_cols INTEGER NOT NULL,
  position INTEGER NOT NULL
);
CREATE TABLE item (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  list_id INTEGER NOT NULL,
  title TEXT NOT NULL,
  url TEXT NOT NULL,
  icon TEXT NOT NULL,
  position INTEGER NOT NULL,
  FOREIGN KEY(list_id) REFERENCES list(id)
);
CREATE TABLE session (
  token TEXT NOT NULL,
  expires_at TEXT NOT NULL
);
