<?php

class Icons {
  
  public function get($iconName) {
    $extension = pathinfo($iconName, PATHINFO_EXTENSION);
    $source = "https://cdn.jsdelivr.net/gh/homarr-labs/dashboard-icons/{$extension}/{$iconName}";
    $destination = "/data/icons/{$iconName}";

    if (!is_dir("/data/icons")) {
      if (!mkdir("/data/icons", 0777, true)) {
        error_log('ERROR: Could not create /data/icons directory!');
        return false;
      };
    }

    if (!file_exists($destination)) {
      if (copy($source, $destination) === false) {
        error_log("ERROR: Could not download icon: {$iconName}. Available icons: https://github.com/homarr-labs/dashboard-icons/blob/main/ICONS.md");
        return false;
      }
      return true;
    }

    return true;
  }

}
