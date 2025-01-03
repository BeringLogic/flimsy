<?php
require('db.php');
require('icons.php');

$db = new DB();
$listId = $_GET['listId'];

$itemId = $db->getNewId();
$title = $_POST['title'];
$href = $_POST['href'];
$icon = $_POST['icon'];

if ($db->AddItem($listId, $itemId, $title, $href, $icon) === false) {
  error_log("ERROR: Could not add Item!");
}

$icons = new Icons();
if (!$icons->get($icon)) {
  error_log("ERROR: Could not download icon!");
}

header("Location: /index.php");
