<?php
require('../classes/DB.php');

$db = new DB();
$title = $_POST['title'];
$numberOfRows = $_POST['number_of_rows'];

if ($db->AddList($title, $numberOfRows) === false) {
  error_log("ERROR: Could not add List!");
}

header("Location: /index.php");
