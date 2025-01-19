<?php

class Icons {
  
  public function get($iconName) {
    $extension = pathinfo($iconName, PATHINFO_EXTENSION);
    $source = "https://cdn.jsdelivr.net/gh/walkxcode/dashboard-icons/{$extension}/{$iconName}";
    $destination = "/data/icons/{$iconName}";

    if (!is_dir("/data/icons")) {
      if (!mkdir("/data/icons", 0777, true)) {
        return false;
      };
    }

    if (!file_exists($destination)) {
      return copy($source, $destination);
    }

    return true;
  }

}
