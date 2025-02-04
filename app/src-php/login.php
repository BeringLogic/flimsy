<?php
session_start();

if (empty($_SERVER['FLIMSY_USERNAME']) || empty($_SERVER['FLIMSY_PASSWORD']) || (isset($_POST['username']) && $_POST['username'] == $_SERVER['FLIMSY_USERNAME'] && isset($_POST['password']) && $_POST['password'] == $_SERVER["FLIMSY_PASSWORD"])) {
  $_SESSION['loggedIn'] = true;
  $_SESSION['message'] = "Successfully logged in!\\n- Click on the gear button to customize the appearance\\n- Click on items and lists to edit them\\n- Drag & drop to reorder.";
}
else {
  $_SESSION['loggedIn'] = false;
  $_SESSION['message'] = "Incorrect username or password!";
}

header('Location: /index.php');
