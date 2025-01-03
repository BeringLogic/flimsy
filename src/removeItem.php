<?php
require('db.php');

$db = new DB();
$itemId = $_GET['itemId'];

if ($db->RemoveItem($itemId) === false) {
  error_log('ERROR: Could not remove the item');
}

header('Location: /index.php');
