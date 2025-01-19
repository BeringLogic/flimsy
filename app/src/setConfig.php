<?php
require('db.php');
require('Icons.php');
require('ColorAutodetector.php');

$db = new db();
$icon = $_POST['icon'];
$title = $_POST['title'];
$backround_image = $_POST['backround_image'];
$color_type = $_POST['color_type'];
$cpu_temp_sensor = $_POST['cpu_temp_sensor'];
$mount_points = $_POST['mount_points'];

switch ($color_type) {
  case 'autodetect':
    $colorDetector = new ColorAutodetector();
    $colors = $colorDetector->extractColors($backround_image);
    $color_background = "black";
    $color_foreground = $colors[1];
    $color_items = $colors[2];
    $color_borders = $colors[0];
    break;
  case 'catppuccin_latte':         # Colors from https://github.com/catppuccin/catppuccin/blob/main/docs/style-guide.md
    $backround_image = '';
    $color_background = "#eff1f5"; # Background Pane
    $color_foreground = "#4c4f69"; # Cursor Line
    $color_items = "#dce0e8";      # Secondary Panes, Crust
    $color_borders = "#9ca0b0";    # Inactive Border
    break;
  case 'catppuccin_mocha':         # Colors from https://github.com/catppuccin/catppuccin/blob/main/docs/style-guide.md
    $backround_image = '';
    $color_background = "#1e1e2e"; # Background Pane
    $color_foreground = "#cdd6f4"; # Cursor Line
    $color_items = "#11111b";      # Secondary Panes, Crust
    $color_borders = "#6c7086";    # Inactive Border
    break;
  default:
  case 'manual':
    $color_background = $_POST['color_background'];
    $color_foreground = $_POST['color_foreground'];
    $color_items = $_POST['color_items'];
    $color_borders = $_POST['color_borders'];
    break;
}

$db->SetConfig($icon, $title, $backround_image, $color_background, $color_foreground, $color_items, $color_borders, $cpu_temp_sensor, $mount_points);

if (!empty($icon)) {
  $icons = new Icons();
  if (!$icons->get($icon)) {
    error_log("ERROR: Could not download icon!");
  }
}

#header("Location: /index.php");
