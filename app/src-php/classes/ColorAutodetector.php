<?php
require('../vendor/autoload.php');

use League\ColorExtractor\Palette;
use League\ColorExtractor\Color;
use League\ColorExtractor\ColorExtractor;

class ColorAutodetector {

  public function extractColors($filename) {

    try {
      if (empty($filename)) {
        throw new \InvalidArgumentException('Filename is empty');
      }

      $fullPath = '/data/backgrounds/' . basename($filename);
      if (!file_exists($fullPath)) {
        throw new \InvalidArgumentException("$fullPath does not exist");
      }

      $palette = Palette::fromFilename($fullPath);
      $extractor = new ColorExtractor($palette);
      $colors = $extractor->extract(5);

      $htmlColors = [];
      foreach($colors as $color) {
        $htmlColors[] = Color::fromIntToHex($color);
      }

      return $htmlColors;
    }
    catch (\Exception $e) {
      error_log('ERROR: ' . $e->getMessage());
      return array('black', 'black', 'black'); 
    }
  }

}
