<?php
require('db.php');

$db = new db();
$id = $_GET['id'];

header('Content-Type: application/json');
echo json_encode($db->GetItem($id));
