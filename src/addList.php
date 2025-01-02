<?php
require('db.php');

$db = new DB();
$id = $db->getNewId();
$title = $_POST['title'];
$items = array();

if ($db->SaveList($id, $title, $items) === false) {
  error_log("ERROR: Could not save List!");
}

header("Location: /index.php");
