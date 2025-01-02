<?php

$id = bin2hex(random_bytes(32));
$title = $_POST['title'];

if (file_exists("/data/data.json")) {
  $data = json_decode(file_get_contents("/data/data.json"), true);
  $data[] = array("id" => $id, "title" => $title, "items" => array());
  file_put_contents("/data/data.json", json_encode($data, JSON_PRETTY_PRINT));
}

header("Location: /index.php");
