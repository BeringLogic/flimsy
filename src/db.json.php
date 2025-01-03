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
      error_log("ERROR: /data/data.json could not be created! Check permissions. Owner and Group must be www-data (UID 33 and GID 33).");
      return false;
    }

    if ($this->saveData('[]') === false) {
      error_log("ERROR: /data/data.json could not be written to! Check permissions. Owner and Group must be www-data (UID 33 and GID 33).");
      return false;
    }

    return true;
  }

  public function getRawData() {
    if (!file_exists($this->file)) {
      error_log("ERROR: File does not exist!");
      return null;
    }

    $fileContents = file_get_contents($this->file);
    if ($fileContents === false) {
      error_log("ERROR: Could not get file contents!");
      return null;
    }

    return $fileContents;
  }

  public function AddList($title) {
    $data = $this->loadData();
    if ($data === null) {
      error_log("ERROR: Could not load data");
      return false;
    }
  
    $data[] = array("id" => $this->getNewId(), "title" => $title, "items" => array());
  
    return $this->saveData($data);
  }

  public function EditList($id, $title) {
    $data = $this->loadData();
    if ($data === null) {
      error_log("ERROR: Could not load data");
      return false;
    }
  
    $listIndex = $this->getIndex($data, $id);
    if ($listIndex === null) {
      error_log("ERROR: List with id $id not found");
      return false;
    }

    $data[$listIndex]['title'] = $title;
  
    return $this->saveData($data);
  }

  public function RemoveList($id) {
    $data = $this->loadData();
    if ($data === null) {
      error_log("ERROR: Could not load data");
      return false;
    }

    $listIndex = $this->getIndex($data, $id);
    if ($listIndex === null) {
      error_log("ERROR: List with id $id not found");
      return false;
    }

    unset($data[$listIndex]);

    return $this->saveData(array_values($data));
  }

  public function AddItem($listId, $itemId, $title, $href, $icon) {
    $data = $this->loadData();
    if ($data === null) {
      error_log("ERROR: Could not load data");
      return false;
    }

    $listIndex = $this->getIndex($data, $listId);
    if ($listIndex === null) {
      error_log("ERROR: List with id $listId not found");
      return false;
    }

    $data[$listIndex]['items'][] = array("id" => $itemId, "title" => $title, "href" => $href, "icon" => $icon);

    return $this->saveData($data);
  }

  public function EditItem($listId, $itemId, $title, $href, $icon) {
    $data = $this->loadData();
    if ($data === null) {
      error_log("ERROR: Could not load data");
      return false;
    }

    $listIndex = $this->getIndex($data, $listId);
    if ($listIndex === null) {
      error_log("ERROR: List with id $listId not found");
      return false;
    }

    $itemIndex = $this->getIndex($data[$listIndex]['items'], $itemId);
    if ($itemIndex === null) {
      error_log("ERROR: Item with id $itemId not found");
      return false;
    }

    $data[$listIndex]['items'][$itemIndex]['title'] = $title;
    $data[$listIndex]['items'][$itemIndex]['href'] = $href;
    $data[$listIndex]['items'][$itemIndex]['icon'] = $icon;

    return $this->saveData($data);
  }

  public function RemoveItem($listId, $itemId) {
    $data = $this->loadData();
    if ($data === null) {
      error_log("ERROR: Could not load data");
      return false;
    }

    $listIndex = $this->getIndex($data, $listId);
    if ($listIndex === null) {
      error_log("ERROR: List with id $listId not found");
      return false;
    }

    $itemIndex = $this->getIndex($data[$listIndex]['items'], $itemId);
    if ($itemIndex === null) {
      error_log("ERROR: Item with id $itemId not found");
      return false;
    }

    unset($data[$listIndex]['items'][$itemIndex]);
    $data[$listIndex]['items'] = array_values($data[$listIndex]['items']);

    return $this->saveData($data);
  }




  public function getNewId() {
    return bin2hex(random_bytes(32));
  }


  private function getIndex($data, $id) {
    foreach ($data as $index=>$item) {
      if ($item['id'] == $id) {
        return $index;
      }
    }
    return null;
  }

  private function loadData() {
    $fileContents = $this->getRawData();
    if ($fileContents === null) {
      return null;
    }

    $data = json_decode($fileContents, true);
    if ($data === null) {
      error_log("ERROR: Could not decode file contents!");
      error_log(print_r($fileContents, true));
      return null;
    }

    return $data;
  }

  private function saveData($data) {
    if (!file_exists($this->file)) {
      error_log("ERROR: File does not exist!");
      return false;
    }

    if (@file_put_contents($this->file, json_encode($data, JSON_PRETTY_PRINT)) === false) {
      error_log("ERROR: Could not write to file!");
      return false;
    }

    return true;
  }

}
