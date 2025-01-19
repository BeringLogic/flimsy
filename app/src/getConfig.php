<?php
require('db.php');

$db = new db();

$result = $db->GetConfig();
// $result['sensors'] = array();
exec("sensors -jA", $output);
$result['sensors'] = json_decode(implode('', $output), true);
// foreach($output as $line) {
  // $result['sensors'][] = array(
  //   'text' => $line,
  //   'value' => substr($line, 0, strpos($line, ':'))
  // );
// }

header('Content-Type: application/json');
echo json_encode($result);
