<?php
require('db.php');

$db = new DB();
$listId = $_GET['listId'];
$itemId = $_GET['itemId'];

if ($db->RemoveItem($listId, $itemId) === false) {
  error_log('ERROR: Could not remove the item');
}

header('Location: /index.php');
