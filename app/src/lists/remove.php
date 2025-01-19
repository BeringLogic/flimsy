<?php
require('../classes/DB.php');

$db = new DB();
$id = $_GET['id'];

if ($db->RemoveList($id) === false) {
  error_log("ERROR: Could not remove List!");
}

header("Location: /index.php");
