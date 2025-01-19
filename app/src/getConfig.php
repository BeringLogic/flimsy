<?php
require('db.php');

$db = new db();

header('Content-Type: application/json');
echo json_encode($db->GetConfig());
