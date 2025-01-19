<?php

class UploadBackground
{
  public function upload()
  {
    try {
      if (
        empty($_FILES['background_image']['name'])
        || empty($_FILES['background_image']['tmp_name'])
      ) {
        throw new Exception("No file uploaded");
      }

      if (!is_uploaded_file($_FILES["background_image"]["tmp_name"])) {
        throw new Exception("is_uploaded_file() failed.");
      }

      // Check if image file is a actual image or fake image
      $check = getimagesize($_FILES["background_image"]["tmp_name"]);
      if ($check === false) {
        throw new Exception("File is not an image - " . $check["mime"] . ".");
      }

      $target_dir = "/data/backgrounds/";
      if (!is_dir($target_dir)) {
        mkdir($target_dir, 0755, true);
      }

      $target_file = $target_dir . basename($_FILES["background_image"]["name"]);
      if (!move_uploaded_file($_FILES["background_image"]["tmp_name"], $target_file)) {
        throw new Exception("Could not save file: " . $target_file . ".");
      }

      return basename($target_file);
    }
    catch (Exception $e) {
      error_log('ERROR: ' . $e->getMessage());
      return null;
    }
  }
}
