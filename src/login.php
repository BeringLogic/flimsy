<?php
session_start();

if (isset($_POST['username']) && $_POST['username'] == $_ENV['FLIMSY_USERNAME'] && isset($_POST['password']) && $_POST['password'] == $_ENV["FLIMSY_PASSWORD"]) {
  $_SESSION['loggedIn'] = true;
}
else {
  session_destroy();
}

header('Location: /index.php');
