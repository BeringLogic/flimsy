<?php
require('db.php');

$db = new db();
$icon = $_POST['icon'];
$title = $_POST['title'];
$backround_image = $_POST['backround_image'];

$db->SetConfig($icon, $title, $backround_image);

header("Location: /index.php");
