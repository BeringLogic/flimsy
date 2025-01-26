<?php
require('../classes/DB.php');

$db = new DB();
$itemId = $_GET['itemId'];

if ($db->deleteItem($itemId) === false) {
  error_log('ERROR: Could not delete the item');
}

header('Location: /index.php');
