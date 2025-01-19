<?php
require('db.php');

$db = new DB();

$listId = $_POST['listId'];
$itemIds = $_POST['itemIds'];

if (!$db->reorderItems($listId, $itemIds)) {
  error_log('ERROR: Could not reorder items');
  http_response_code(500);
  exit();
}
