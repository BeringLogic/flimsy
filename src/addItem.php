<?php
require('db.php');

$db = new DB();
$listId = $_GET['id'];

$itemId = $db->getNewId();
$title = $_POST['title'];
$href = $_POST['href'];
$icon = $_POST['icon'];

if (!$db->SaveItem($listId, $itemId, $title, $href, $icon)) {
  error_log("ERROR: Could not save Item!");
}

header("Location: /index.php");
