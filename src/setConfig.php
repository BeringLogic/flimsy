<?php
require('db.php');

$db = new db();
$icon = $_POST['icon'];
$title = $_POST['title'];
$backround_image = $_POST['backround_image'];
$color_background = $_POST['color_background'];
$color_foreground = $_POST['color_foreground'];
$color_items = $_POST['color_items'];
$color_borders = $_POST['color_borders'];

$db->SetConfig($icon, $title, $backround_image, $color_background, $color_foreground, $color_items, $color_borders);

header("Location: /index.php");
