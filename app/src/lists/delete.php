<?php
require('../classes/DB.php');

$db = new DB();
$id = $_GET['id'];

if ($db->deleteList($id) === false) {
  error_log("ERROR: Could not delete List!");
}

header("Location: /index.php");
