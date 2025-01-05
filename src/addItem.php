<?php
require('db.php');
require('Icons.php');

$db = new DB();
$listId = $_GET['listId'];
$title = $_POST['title'];
$href = $_POST['href'];
$icon = $_POST['icon'];

if ($db->AddItem($listId, $title, $href, $icon) === false) {
  error_log("ERROR: Could not add Item!");
}

$icons = new Icons();
if (!$icons->get($icon)) {
  error_log("ERROR: Could not download icon!");
}

header("Location: /index.php");
