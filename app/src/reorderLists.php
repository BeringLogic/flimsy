<?php
require('db.php');

$db = new DB();

$listIds = $_POST['listIds'];

if (!$db->reorderLists($listIds)) {
  error_log('ERROR: Could not reorder lists');
  http_response_code(500);
  exit();
}
