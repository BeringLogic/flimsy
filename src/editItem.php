<?php
require('db.php');
require('Icons.php');

$db = new DB();
$itemId = $_GET['itemId'];
$title = $_POST['title'];
$href = $_POST['href'];
$icon = $_POST['icon'];

if ($db->EditItem($itemId, $title, $href, $icon) === false) {
  error_log("ERROR: Could not edit Item!");
}

$icons = new Icons();
if (!$icons->get($icon)) {
  error_log("ERROR: Could not download icon!");
}

header("Location: /index.php");

