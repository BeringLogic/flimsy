<?php
require('db.php');

$db = new db();

$result = $db->GetConfig();

exec("sensors -jA", $output);
$result['sensors'] = json_decode(implode('', $output), true);

header('Content-Type: application/json');
echo json_encode($result);
