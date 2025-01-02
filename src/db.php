<?php

class DB {
  private $file = '/data/data.json';

  /**
   * Check if there is a DB and create one if there are none
   * @return false if there is no DB and it failed creating a new one
   **/
  public function init() {
    if (file_exists('/data/data.json')) {
      return true;
    }

    error_log("Creating /data/data.json...");
    if (!@touch('/data/data.json')) {
      error_log("/data/data.json could not be created! Check permissions. Owner and Group must be www-data (UID 33 and GID 33).");
      return false;
    }

    if (!$this->saveData('[]')) {
      error_log("/data/data.json could not be written to! Check permissions. Owner and Group must be www-data (UID 33 and GID 33).");
      return false;
    }

    return true;
  }

  public function SaveList($id, $title, $items) {
    $data = $this->loadData();
    if ($data === null) {
      return false;
    }

    $data[] = array("id" => $id, "title" => $title, "items" => $items);
  
    return $this->saveData($data);
  }

  public function SaveItem($listId, $itemId, $title, $href, $icon) {
    $data = $this->loadData();
    if ($data === null) {
      return false;
    }

    $listIndex = $this->getListIndex($data, $listId);
    if ($listIndex === null) {
      error_log("List with id $listId not found");
      return false;
    }

    $data[$listIndex]['items'][] = array("id" => $itemId, "title" => $title, "href" => $href, "icon" => $icon);

    return $this->saveData($data);
  }

  public function getNewId() {
    return bin2hex(random_bytes(32));
  }


  private function getListIndex($data, $id) {
    foreach ($data as $index=>$item) {
      if ($item['id'] == $id) {
        return $index;
      }
    }
    return null;
  }

  private function loadData() {
    if (!file_exists($this->file)) {
      return null;
    }

    return json_decode(file_get_contents($this->file), true);
  }

  private function saveData($data) {
    if (!file_exists($this->file)) {
      return false;
    }

    if (@file_put_contents($this->file, json_encode($data, JSON_PRETTY_PRINT)) === false) {
      return false;
    }

    return true;
  }

}
