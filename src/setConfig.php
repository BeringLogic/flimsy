<?php
require('db.php');
require('ColorAutodetector.php');

$db = new db();
$icon = $_POST['icon'];
$title = $_POST['title'];
$backround_image = $_POST['backround_image'];
$autodetect_colors = $_POST['autodetect_colors'];

if ($autodetect_colors == 'autodetect') {
  $colorDetector = new ColorAutodetector();
  $colors = $colorDetector->extractColors($backround_image);
  $color_background = "black";
  $color_foreground = $colors[1];
  $color_items = $colors[2];
  $color_borders = $colors[0];
}
else {
  $color_background = $_POST['color_background'];
  $color_foreground = $_POST['color_foreground'];
  $color_items = $_POST['color_items'];
  $color_borders = $_POST['color_borders'];
}

$db->SetConfig($icon, $title, $backround_image, $color_background, $color_foreground, $color_items, $color_borders);

header("Location: /index.php");
