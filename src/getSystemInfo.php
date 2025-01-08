<?php
require('db.php');
$db = new db();
$config = $db->getConfig();

$result = array();

$result['cpu_temp'] = exec('sensors | grep "'.$config["cpu_temp_sensor"].'" | cut -d ":" -f 2 | awk \'{ print $1 }\'');
$result['free_memory'] = exec('free -h | grep Mem | awk \'{ print $4 }\'');
$result['free_swap'] = exec('free -h | tail -n 1 | awk \'{ print $4 }\'');
$result['disks'] = array();

if (!empty($config['mount_points'])) {
  foreach (explode(',', $config['mount_points']) as $mount_point) {
    $result['disks'][] = array(
      'mount_point' => basename($mount_point) ?: '/',
      'free_disk_space' => exec("df -h " . $mount_point . " | tail -n 1 | awk '{ print $4 }'")
    );
  }
}

header('Content-Type: application/json');
echo json_encode($result);
