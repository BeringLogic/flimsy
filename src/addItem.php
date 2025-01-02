<?php

$id = $_GET['id'];
$title = $_POST['title'];
$desc = $_POST['desc'];
$icon = $_POST['icon'];

$newId = bin2hex(random_bytes(32));

function findList($data, $id) {
  foreach ($data as $index=>$item) {
    if ($item['id'] == $id) {
      return $index;
    }
  }
  return null;
}

if (file_exists("/data/data.json")) {
  $data = json_decode(file_get_contents("/data/data.json"), true);
  $listIndex = findList($data, $id);
  if ($listIndex === null) {
    error_log("List with id $id not found");
  }
  else {
    $data[$listIndex]['items'][] = array("id" => $newId, "title" => $title, "desc" => $desc, "icon" => $icon);
    file_put_contents("/data/data.json", json_encode($data, JSON_PRETTY_PRINT));
  }
}

header("Location: /index.php");
