<?php
require('db.php');

$db = new DB();
$id = $_GET["id"];
$title = $_POST['title'];
$numberOfRows = $_POST['number_of_rows'];

if ($db->EditList($id, $title, $numberOfRows) === false) {
  error_log("ERROR: Could not edit List!");
}

header("Location: /index.php");
