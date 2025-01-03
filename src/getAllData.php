<?php
require("db.php");
$db = new DB();

if ($db->init() === false) {
  error_log("ERROR: Could not init DB!");
  http_response_code(500);
  exit;
}

$lists = $db->getAllLists();
$items = $db->getAllItems();

foreach($lists as &$list) {
  $list['items'] = array();
  foreach($items as &$item) {
    if ($item['list_id'] == $list['id']) {
      $list['items'][] = $item;
    }
  }
}

header('Content-Type: application/json');
echo json_encode($lists);
