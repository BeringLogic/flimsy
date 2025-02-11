<?php
require('classes/DB.php');
$db = new db();
$config = $db->getConfig();

$result = array();

if ($config['cpu_temp_sensor'] != "") {
  $result['cpu_temp'] = exec('sensors | grep "'.$config["cpu_temp_sensor"].'" | cut -d ":" -f 2 | awk \'{ print $1 }\'');
}
if ($config['show_free_ram']) {
  $result['free_memory'] = exec('free -h | grep Mem | awk \'{ print $4 }\'');
}
if ($config['show_free_swap']) {
  $result['free_swap'] = exec('free -h | grep Swap | awk \'{ print $4 }\'');
}
if ($config['show_public_ip']) {
  $result['public_ip'] = file_get_contents('https://api.ipify.org');
}
if ($config['show_free_space']) {
  $result['storage'] = array();
  $mountPoints = array_merge(array('/'), glob('/mnt/*', GLOB_ONLYDIR));
  foreach ($mountPoints as $mountPoint) {
    $result['storage'][] = array(
      'mount_point' => basename($mountPoint) ?: '/',
      'free_space' => exec("df -h " . $mountPoint . " | tail -n 1 | awk '{ print $4 }'")
    );
  }
}

header('Content-Type: application/json');
echo json_encode($result);
