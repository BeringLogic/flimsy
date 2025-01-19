<?php
require('db.php');

$db = new DB();

if ($db->init() === false) {
  error_log("ERROR: Could not init DB!");
  http_response_code(500);
  exit;
}
