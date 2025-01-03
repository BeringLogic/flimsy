CREATE TABLE config (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  icon TEXT NULL,
  title TEXT NULL,
  backround_image TEXT NULL,
  number_of_rows INTEGER NOT NULL
);
INSERT INTO config (icon, title, backround_image, number_of_rows) VALUES (NULL, 'Flimsy Home Page', NULL, 4);

CREATE TABLE list (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  title TEXT NOT NULL,
);

CREATE TABLE item (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  list_id INTEGER NOT NULL,
  title TEXT NOT NULL,
  href TEXT NOT NULL,
  icon TEXT NOT NULL,
  FOREIGNGN KEY(list_id) REFERENCES lists(id)
);

