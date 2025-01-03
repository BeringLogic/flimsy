<?php

class DB {
  private $dbh; 
  
  public function __construct()
  {
    $this->dbh = new SQLite3('/data/flimsy.db', SQLITE3_OPEN_CREATE | SQLITE3_OPEN_READWRITE);
  }

  /**
   * Try to get the config and seed the new DB if needed
   * @return false if there was an error while seeding
   **/
  public function init() {
    if (@$this->dbh->prepare('SELECT * FROM config WHERE id = 1') === false) {
      error_log("Seeding /data/flimsy.db...");

      $seed = array(
        "CREATE TABLE config (
          id INTEGER PRIMARY KEY AUTOINCREMENT,
          icon TEXT NULL,
          title TEXT NULL,
          backround_image TEXT NULL,
          number_of_rows INTEGER NOT NULL
        );",
        "INSERT INTO config (icon, title, backround_image, number_of_rows) VALUES (NULL, 'Flimsy Home Page', NULL, 4);",
        "CREATE TABLE list (
          id INTEGER PRIMARY KEY AUTOINCREMENT,
          title TEXT NOT NULL
        );",
        "CREATE TABLE item ( 
          id INTEGER PRIMARY KEY AUTOINCREMENT,
          list_id INTEGER NOT NULL,
          title TEXT NOT NULL,
          href TEXT NOT NULL,
          icon TEXT NOT NULL,
          FOREIGN KEY(list_id) REFERENCES lists(id)
        );"
      );

      foreach ($seed as $query) {
        #error_log($query);
        if ($this->dbh->query($query) === false) {
          error_log("ERROR: Could not seed /data/flimsy.db! Check permissions. Owner and Group must be www-data (UID 33 and GID 33).");
          return false;
        }
      }
    }

    return true;
  }

  public function getAllLists() {
    $stmt = $this->dbh->prepare('SELECT * FROM list');
    $result = $stmt->execute();
    $data = array();
    while ($row = $result->fetchArray(SQLITE3_ASSOC)) {
      $data[] = $row;
    }
    $result->finalize();
    return $data;
  }

  public function getAllItems() {
    $stmt = $this->dbh->prepare('SELECT * FROM item');
    $result = $stmt->execute();
    $data = array();
    while ($row = $result->fetchArray(SQLITE3_ASSOC)) {
      $data[] = $row;
    }
    $result->finalize();
    return $data;
  }

  public function getRawData() {
    $config = $this->dbh->querySingle('SELECT * FROM config WHERE id = 1');
    echo "config = "; print_r($config);

    $list = $this->dbh->querySingle('SELECT * FROM list');
    echo "list = "; print_r($list);

    $item = $this->dbh->querySingle('SELECT * FROM item');
    echo "item = "; print_r($item);


    $stmt = $this->dbh->prepare('SELECT * FROM config');
    $result = $stmt->execute();
    $data = $result->fetchArray(SQLITE3_ASSOC);
    echo "data = "; print_r($data);
    $result->finalize();
    return false;
  }

  public function AddList($title) {
    $stmt = $this->dbh->prepare('INSERT INTO list (title) VALUES (:title)');
    $stmt->bindValue(':title', $title);
    return $stmt->execute() !== false;
  }

  public function EditList($id, $title) {
    $stmt = $this->dbh->prepare('UPDATE list SET title = :title WHERE id = :id');
    $stmt->bindValue(':title', $title);
    $stmt->bindValue(':id', $id);
    return $stmt->execute() !== false;
  }

  public function RemoveList($id) {
    $stmt = $this->dbh->prepare('DELETE FROM list WHERE id = :id');
    $stmt->bindValue(':id', $id);
    return $stmt->execute() !== false;
  }

  public function AddItem($listId, $title, $href, $icon) {
    $stmt = $this->dbh->prepare('INSERT INTO item (list_id, title, href, icon) VALUES (:list_id, :title, :href, :icon)');
    $stmt->bindValue(':list_id', $listId);
    $stmt->bindValue(':title', $title);
    $stmt->bindValue(':href', $href);
    $stmt->bindValue(':icon', $icon);
    return $stmt->execute() !== false;
  }

  public function EditItem($itemId, $title, $href, $icon) {
    $stmt = $this->dbh->prepare('UPDATE item SET title = :title, href = :href, icon = :icon WHERE id = :id');
    $stmt->bindValue(':title', $title);
    $stmt->bindValue(':href', $href);
    $stmt->bindValue(':icon', $icon);
    $stmt->bindValue(':id', $itemId);
    return $stmt->execute() !== false;
  }

  public function RemoveItem($id) {
    $stmt = $this->dbh->prepare('DELETE FROM item WHERE id = :id');
    $stmt->bindValue(':id', $id);
    return $stmt->execute() !== false;
  }

}

