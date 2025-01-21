<?php
require('classes/DB.php');
$db = new db();
$config = $db->getConfig();

$result = array();

$result['cpu_temp'] = exec('sensors | grep "'.$config["cpu_temp_sensor"].'" | cut -d ":" -f 2 | awk \'{ print $1 }\'');
$result['free_memory'] = exec('free -h | grep Mem | awk \'{ print $4 }\'');
$result['free_swap'] = exec('free -h | grep Swap | awk \'{ print $4 }\'');
$result['storage'] = array();

if (!empty($config['mount_points'])) {
  foreach (explode(',', $config['mount_points']) as $mount_point) {
    if (!is_dir($mount_point)) {
      error_log("ERROR: Invalid mount point: " . $mount_point . ". It doesn't exist or user www-data doesn't have permission to reach it.");
      continue;
    }
    $result['storage'][] = array(
      'mount_point' => basename($mount_point) ?: '/',
      'free_space' => exec("df -h " . $mount_point . " | tail -n 1 | awk '{ print $4 }'")
    );
  }
}

header('Content-Type: application/json');
echo json_encode($result);
