<?php
require('classes/DB.php');

$db = new DB();

$result = array();
if ($db->init() === false) {
  $result['success'] = false;
  $message = "ERROR: Could not init DB! You must bind a volume to /data and it must be writable by the 1000:1000 user.";
  $result['error'] = $message;
  error_log($message);
}
else {
  $result['success'] = true;
}

header('Content-Type: application/json');
echo json_encode($result);
