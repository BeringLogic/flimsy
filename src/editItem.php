<?php
require('db.php');

$db = new DB();
$listId = $_GET['listId'];
$itemId = $_GET['itemId'];
$title = $_POST['title'];
$href = $_POST['href'];
$icon = $_POST['icon'];

if ($db->EditItem($listId, $itemId, $title, $href, $icon) === false) {
  error_log("ERROR: Could not edit Item!");
}

header("Location: /index.php");

