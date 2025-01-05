<?php

class Icons {
  
  public function get($iconName) {
    $source = "https://cdn.jsdelivr.net/gh/walkxcode/dashboard-icons/png/{$iconName}.png";
    $destination = "/data/icons/{$iconName}.png";

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
