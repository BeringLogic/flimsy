<?php

class DB {
  private $dbh; 
  
  public function __construct()
  {
    try {
      $this->dbh = new SQLite3('/data/flimsy.db', SQLITE3_OPEN_CREATE | SQLITE3_OPEN_READWRITE);
    }
    catch (Exception) {
      $this->dbh = null;
    }
  }

  /**
   * Try to get the config and seed the new DB if needed
   * @return false if there was an error while seeding
   **/
  public function init() {
    $cpuTempSensor = empty($_SERVER['FLIMSY_CPU_TEMP_SENSOR']) ? null : $_SERVER['FLIMSY_CPU_TEMP_SENSOR'];
    $mountPoints = empty($_SERVER['FLIMSY_MOUNT_POINTS']) ? null : $_SERVER['FLIMSY_MOUNT_POINTS'];
    
    if ($this->dbh === null) {
      error_log("Unable to open database!");
      return false;
    }

    if (@$this->dbh->prepare('SELECT * FROM config WHERE id = 1') === false) {
      error_log("Seeding /data/flimsy.db...");

      $seed = array(
        "CREATE TABLE config (
          id INTEGER PRIMARY KEY AUTOINCREMENT,
          icon TEXT NULL,
          title TEXT NULL,
          backround_image TEXT NULL,
          color_background TEXT NOT NULL,
          color_foreground TEXT NOT NULL,
          color_items TEXT NOT NULL,
          color_borders TEXT NOT NULL,
          cpu_temp_sensor TEXT,
          mount_points TEXT
        );",
        "INSERT INTO config (icon, title, color_background, color_foreground, color_items, color_borders, cpu_temp_sensor, mount_points) VALUES (NULL, 'Flimsy Home Page', '#1e1e2e', '#cdd6f4', '#11111b', '#6c7086', '$cpuTempSensor', '$mountPoints');",
        "CREATE TABLE list (
          id INTEGER PRIMARY KEY AUTOINCREMENT,
          title TEXT NOT NULL,
          number_of_rows INTEGER NOT NULL,
          position INTEGER NOT NULL
        );",
        "CREATE TABLE item ( 
          id INTEGER PRIMARY KEY AUTOINCREMENT,
          list_id INTEGER NOT NULL,
          title TEXT NOT NULL,
          href TEXT NOT NULL,
          icon TEXT NOT NULL,
          position INTEGER NOT NULL,
          FOREIGN KEY(list_id) REFERENCES lists(id)
        );"
      );

      foreach ($seed as $query) {
        #error_log($query);
        if ($this->dbh->query($query) === false) {
          error_log("ERROR: Could not seed /data/flimsy.db!");
          return false;
        }
      }
    }

    return true;
  }

  public function GetConfig() {
    $stmt = $this->dbh->prepare('SELECT * FROM config WHERE id = 1');
    $result = $stmt->execute();
    $data = $result->fetchArray(SQLITE3_ASSOC);
    $result->finalize();
    return $data;
  }
  public function SetConfig($icon, $title, $backround_image, $color_background, $color_foreground, $color_items, $color_borders, $cpu_temp_sensor, $mount_points) {
    $stmt = $this->dbh->prepare('UPDATE config SET icon = :icon, title = :title, backround_image = :backround_image, color_background = :color_background, color_foreground = :color_foreground, color_items = :color_items, color_borders = :color_borders, cpu_temp_sensor = :cpu_temp_sensor, mount_points = :mount_points WHERE id = 1');
    $stmt->bindValue(':icon', $icon);
    $stmt->bindValue(':title', $title);
    $stmt->bindValue(':backround_image', $backround_image);
    $stmt->bindValue(':color_background', $color_background);
    $stmt->bindValue(':color_foreground', $color_foreground);
    $stmt->bindValue(':color_items', $color_items);
    $stmt->bindValue(':color_borders', $color_borders);
    $stmt->bindValue(':cpu_temp_sensor', $cpu_temp_sensor);
    $stmt->bindValue(':mount_points', $mount_points);
    return $stmt->execute() !== false;
  }

  public function getAllLists() {
    $stmt = $this->dbh->prepare('SELECT * FROM list ORDER BY position');
    $result = $stmt->execute();
    $data = array();
    while ($row = $result->fetchArray(SQLITE3_ASSOC)) {
      $data[] = $row;
    }
    $result->finalize();
    return $data;
  }

  public function getAllItems() {
    $stmt = $this->dbh->prepare('SELECT * FROM item ORDER BY list_id, position');
    $result = $stmt->execute();
    $data = array();
    while ($row = $result->fetchArray(SQLITE3_ASSOC)) {
      $data[] = $row;
    }
    $result->finalize();
    return $data;
  }

  public function AddList($title, $numberOfRows) {
    $stmt = $this->dbh->prepare('INSERT INTO list (title, number_of_rows, position) VALUES (:title, :number_of_rows, (SELECT IFNULL(max(position), 0) + 1 FROM list))');
    $stmt->bindValue(':title', $title);
    $stmt->bindValue(':number_of_rows', $numberOfRows);
    return $stmt->execute() !== false;
  }

  public function EditList($id, $title, $numberOfRows) {
    $stmt = $this->dbh->prepare('UPDATE list SET title = :title, number_of_rows = :number_of_rows WHERE id = :id');
    $stmt->bindValue(':title', $title);
    $stmt->bindValue(':number_of_rows', $numberOfRows);
    $stmt->bindValue(':id', $id);
    return $stmt->execute() !== false;
  }

  public function RemoveList($id) {
    $stmt = $this->dbh->prepare('DELETE FROM list WHERE id = :id');
    $stmt->bindValue(':id', $id);
    return $stmt->execute() !== false;
  }

  public function reorderLists($listIds) {
    $ret = true;

    $this->dbh->exec('BEGIN TRANSACTION'); 
    $stmt = $this->dbh->prepare('UPDATE list SET position = :position WHERE id = :id');
    for ($i = 0; $i < count($listIds); $i++) {
      $stmt->bindValue(':position', $i);
      $stmt->bindValue(':id', $listIds[$i]);
      if ($stmt->execute() === false) {
        $ret = false;
        break;
      }
    }
    $this->dbh->exec('COMMIT');

    return $ret;
  }
  public function AddItem($listId, $title, $href, $icon) {
    $stmt = $this->dbh->prepare('INSERT INTO item (list_id, title, href, icon, position) VALUES (:list_id, :title, :href, :icon, (select IFNULL(max(position), 0) + 1 from item where list_id = :list_id))');
    $stmt->bindValue(':list_id', $listId);
    $stmt->bindValue(':title', $title);
    $stmt->bindValue(':href', $href);
    $stmt->bindValue(':icon', $icon);
    return $stmt->execute() !== false;
  }

  public function GetItem($itemId) {
    $stmt = $this->dbh->prepare('SELECT * FROM item WHERE id = :id');
    $stmt->bindValue(':id', $itemId);
    $result = $stmt->execute();
    $data = $result->fetchArray(SQLITE3_ASSOC);
    $result->finalize();
    return $data;
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

  public function reorderItems($listId, $itemIds) {
    $ret = true;

    $this->dbh->exec('BEGIN TRANSACTION'); 
    $stmt = $this->dbh->prepare('UPDATE item SET list_id = :list_id, position = :position WHERE id = :id');
    for ($i = 0; $i < count($itemIds); $i++) {
      $stmt->bindValue(':list_id', $listId);
      $stmt->bindValue(':position', $i);
      $stmt->bindValue(':id', $itemIds[$i]);
      if ($stmt->execute() === false) {
        $ret = false;
        break;
      }
    }
    $this->dbh->exec('COMMIT');
    return $ret;
  }
}

