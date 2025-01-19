<?php
require('classes/DB.php');

$db = new DB();

if ($db->init() === false) {
  error_log("ERROR: Could not init DB! You must bind a volume to /data and it must be writable by the selected user in compose.yaml.");
  http_response_code(500);
  exit;
}
