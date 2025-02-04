<?php
require('classes/DB.php');
$db = new DB();

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
