<?php

$title = $_POST['title'];

if (file_exists("data/data.json")) {
  $data = json_decode(file_get_contents("data/data.json"), true);
  $data[] = array("title" => $title, "items" => array());
  file_put_contents("data/data.json", json_encode($data));
}

