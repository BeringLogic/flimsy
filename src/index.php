<?php session_start(); ?>

<!DOCTYPE html>
<head>
  <title>Flimsy</title>
  <link rel="stylesheet" href="style.css">
  <link rel="stylesheet" href="https://code.jquery.com/ui/1.14.1/themes/base/jquery-ui.css">
  <script
			  src="https://code.jquery.com/jquery-3.7.1.min.js"
			  integrity="sha256-/JqT3SQfawRcv/BIHPThkBvs0OEvtFFmqPF/lYI/Cxo="
			  crossorigin="anonymous"></script>
  <script
			  src="https://code.jquery.com/ui/1.14.1/jquery-ui.min.js"
			  integrity="sha256-AlTido85uXPlSyyaZNsjJXeCs07eSv3r43kyCVc8ChI="
			  crossorigin="anonymous"></script>
</head>

<body>
  <header>
    <div class="weather">
      <img>
      <div class="temp"></div>
      <div class="description"></div>
    </div>
    <?php if (empty($_SESSION['loggedIn'])) { ?>
      <a class="login" href="#">Login</a>
    <?php } else { ?>
      <a class="logout" href="logout.php">Logout</a>
    <?php } ?>
    <h1>Flimsy Home Page</h1>
  </header>

  <button id="addList">➕ Add List</button>

  <div id="loginDialog" class="dialog">
    <form action="login.php" method="post">
      <div class="dialog-field">
        <label for="username">Username</label>
        <input id="username" type="text" name="username" required>
      </div>
      <div class="dialog-field">
        <label for="password">Password</label>
        <input id="password" type="password" name="password" required>
      </div>
    </form>
  </div>

  <div id="editListDialog" class="dialog">
    <form action="editList.php" method="post">
      <div class="dialog-field">
        <label for="listTitle">Title</label>
        <input id="listTitle" type="text" name="title" required>
      </div>
    </form>
  </div>

  <div id="editItemDialog" class="dialog">
    <form action="addItem.php" method="post">
      <div class="dialog-field">
        <label for="itemTitle">Title</label>
        <input id="itemTitle" type="text" name="title" required>
      </div>
      <div class="dialog-field">
        <label for="itemHref">URL</label>
        <input id="itemHref" type="text" name="href" required>
      </div>
      <div class="dialog-field">
        <label for="itemIcon">Icon</label>
        <input id="itemIcon" type="text" name="icon" required>
      </div>
    </form>
  </div>

  <script type="text/javascript">
    function render(data) {
      data.forEach(b => {
        var list = $('<div class="list"></div>');
        list.append($('<h2>' + b.title + ' <button id="editList_' + b.id + '" class="editList">✍️</button><button id="removeList_' + b.id + '" class="removeList">❌</button></h2>'));
        items = $('<div class="items"></div>');

        b.items.forEach(i => {
          var item = $('<div class="item"></div>');
          item.append($('<button id="removeItem_' + b.id + '_' + i.id + '" class="removeItem">❌</button><button id="editItem_' + b.id + '_' + i.id + '" class="editItem">✍️</button>'));

          var link = $('<a class="href" href="' + i.href + '" target="_blank"></a>');
          link.append($('<img class="icon" src="dashboard-icons/' + i.icon + '.png" />'));

          var details = $('<div class="details"></div>');
          details.append($('<div class="title">' + i.title + '</div>'));
          details.append($('<div class="href">' + i.href + '</div>'));

          link.append(details);
          item.append(link);

          items.append(item);
        });

        items.append($('<div class="addItem"><button id="addItem_' + b.id + '">➕ Add Item</button></div>'));

        list.append(items);
        list.insertBefore($('#addList'));
      });

      if (<?php echo !empty($_SESSION['loggedIn']) && $_SESSION['loggedIn'] == true ? "true" : "false"; ?>) {
        $('button').css('display', 'inline-block');
      }
    }

    function updateWeather() {
      if (<?php echo empty($_ENV["FLIMSY_WEATHER_API_KEY"]) ? "true" : "false"; ?>) {
        $('.weather').css('display', 'none');
        return;
      }
      $.ajax({
        url: 'https://api.openweathermap.org/data/2.5/weather',
        data: {
          lon : <?php echo $_ENV["FLIMSY_WEATHER_LONGITUDE"] ?: 0; ?>,
          lat : <?php echo $_ENV["FLIMSY_WEATHER_LATITUDE"] ?: 0; ?>,
          units : '<?php echo $_ENV["FLIMSY_WEATHER_UNITS"] ?: "metric"; ?>',
          lang : '<?php echo $_ENV["FLIMSY_WEATHER_LANGUAGE"] ?: "en"; ?>',
          appid : '<?php echo $_ENV["FLIMSY_WEATHER_API_KEY"] ?: ""; ?>',
        },
        success: (data) => {
          $('.weather img').attr('src', 'https://openweathermap.org/img/wn/' + data.weather[0].icon + '@2x.png');
          $('.weather .description').html(data.weather[0].description);
          $('.weather .temp').html(data.main.temp + ' °C');
        }
      });
    }

    window.onload = () => {
      $.ajax({
        url: 'getAllData.php',
        success: (data) => {
          updateWeather();

          render(data);

          $('button.editList').click((e) => {
            const id = e.target.id.replace('editList_', '');
            $('#editListDialog form').attr('action', 'editList.php?id=' + id);
            $('#editListDialog form #listTitle').val(data.find(l => l.id == id).title);
            $('#editListDialog').dialog('open');
          });
          $('#addList').click(() => {
            $('#editListDialog form').attr('action', 'addList.php');
            $('#editListDialog form #listTitle').val('');
            $('#editListDialog').dialog('open');
          });
          $('button.removeList').click((e) => {
            const id = e.target.id.replace('removeList_', '');
            if (!confirm('Are you sure you want to remove this list, and all its items?')) {
              return;
            }
            location.href = 'removeList.php?id=' + id;
          });

          $('button.editItem').click((e) => {
            const idParts = e.target.id.split('_');
            const listId = idParts[1];
            const itemId = idParts[2];
            $('#editItemDialog form').attr('action', 'editItem.php?listId=' + listId + '&itemId=' + itemId);
            $('#editItemDialog form #itemTitle').val(data.find(l => l.id == listId).items.find(i => i.id == itemId).title);
            $('#editItemDialog form #itemHref').val(data.find(l => l.id == listId).items.find(i => i.id == itemId).href);
            $('#editItemDialog form #itemIcon').val(data.find(l => l.id == listId).items.find(i => i.id == itemId).icon);
            $('#editItemDialog').dialog('open');
          });
          $('div.addItem button').click((e) => {
            const listId = e.target.id.replace('addItem_', '');
            $('#editItemDialog form').attr('action', 'addItem.php?listId=' + listId);
            $('#editItemDialog form #itemTitle').val('');
            $('#editItemDialog form #itemHref').val('');
            $('#editItemDialog form #itemIcon').val('');
            $('#editItemDialog').dialog('open');
          });
          $('button.removeItem').click((e) => {
            const idParts = e.target.id.split('_');
            const listId = idParts[1];
            const itemId = idParts[2];
            if (!confirm('Are you sure you want to remove this item?')) {
              return;
            }
            location.href = 'removeItem.php?listId=' + listId + '&itemId=' + itemId;
          });

          <?php if (!empty($_SESSION['message'])) {
            echo 'alert("' . $_SESSION['message'] . '");';
            unset($_SESSION['message']);
          } ?>
        },
        error: (error, status, xhr) => {
          alert("An error occured while loading the data! Make sure /data is writable by the www-data user (UID 33, GID 33).");
          console.log(error);
          console.log(status);
          console.log(xhr);
        }

      })

      $('#loginDialog').dialog({
        title: 'Login',
        autoOpen: false,
        modal: true,
        buttons: {
          "Login": () => {
            if ($('#loginDialog :invalid').length > 0) {
              return;
            }
            $('#loginDialog form').submit();
          },
          "Cancel": () => {
            $('#loginDialog').dialog('close');
          }
        }
      });
      $('a.login').click(() => {
        $('#loginDialog').dialog('open');
      });

      $('#editListDialog').dialog({
        title: 'Edit List',
        autoOpen: false,
        modal: true,
        buttons: {
          "Save": () => {
            if ($('#editListDialog :invalid').length > 0) {
              return;
            }
            $('#editListDialog form').submit();
          },
          "Cancel": () => {
            $('#editListDialog').dialog('close');
          }
        }
      });

      $('#editItemDialog').dialog({
        title: 'Edit Item',
        autoOpen: false,
        modal: true,
        buttons: {
          "Save": () => {
            if ($('#editItemDialog :invalid').length > 0) {
              return;
            }
            $('#editItemDialog form').submit();
          },
          "Cancel": () => {
            $('#editItemDialog').dialog('close');
          }
        }
      });
    }
  </script>

</body>
</html>
