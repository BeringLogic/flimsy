<?php
session_start();

if (isset($_POST['username']) && $_POST['username'] == $_ENV['FLIMSY_USERNAME'] && isset($_POST['password']) && $_POST['password'] == $_ENV["FLIMSY_PASSWORD"]) {
  $_SESSION['loggedIn'] = true;
  $_SESSION['message'] = "Successfully logged in!";
}
else {
  $_SESSION['loggedIn'] = false;
  $_SESSION['message'] = "Incorrect username or password!";
}

header('Location: /index.php');
