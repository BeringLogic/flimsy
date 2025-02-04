<?php
$host = parse_url($_GET["href"], PHP_URL_HOST);
if ($host === false) {
  echo 500;
  return;
}

$devnull = fopen('/dev/null', 'w');

$ch = curl_init($host);
curl_setopt($ch, CURLOPT_FILE, $devnull);
curl_exec($ch);
echo curl_strerror(curl_errno($ch));
curl_close($ch);

fclose($devnull);
