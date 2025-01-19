<?php
require('../classes/DB.php');

$db = new db();

$result = $db->GetConfig();

$result['backgrounds'] = scandir('/data/backgrounds');

exec("sensors -jA", $output);
$result['sensors'] = json_decode(implode('', $output), true);

header('Content-Type: application/json');
echo json_encode($result);
