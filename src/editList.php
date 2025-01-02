<?php
require('db.php');

$db = new DB();
$id = $_GET["id"];
$title = $_POST['title'];

if ($db->EditList($id, $title) === false) {
  error_log("ERROR: Could not edit List!");
}

header("Location: /index.php");
