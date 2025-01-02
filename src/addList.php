<?php
require('db.php');

$db = new DB();
$title = $_POST['title'];

if ($db->AddList($title) === false) {
  error_log("ERROR: Could not add List!");
}

header("Location: /index.php");
