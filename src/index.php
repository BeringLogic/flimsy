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
    <img id="icon">
    <h1><span id="title"></span><button id="config">⚙️</button></h1>
  </header>

  <button id="addList">➕ Add List</button>

  <div id="configDialog" class="dialog">
    <form action="setConfig.php" method="post">
      <div class="dialog-field">
        <label for="configIcon">Icon</label>
        <input id="configIcon" type="text" name="icon">
      </div>
      <div class="dialog-field">
        <label for="configTitle">Title</label>
        <input id="configTitle" type="text" name="title">
      </div>
      <div class="dialog-field">
        <label for="configBackroundImage">Backround Image</label>
        <input id="configBackroundImage" type="text" name="backround_image">
      </div>
    </form>
  </div>

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
      <div class="dialog-field">
        <label for="listRows">Number of Rows</label>
        <input id="listRows" type="number" min="1" name="number_of_rows" required>
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
    function loadConfig() {
      $.ajax({
        url: 'getConfig.php',
        success: (config) => {
          if (config.icon) {
            $('#icon').attr('src', 'dashboard-icons/' + config.icon + '.png')
                      .css('display', 'inline-block');
          }
          if (config.title) {
            $('#title').html(config.title);
          }
          if (config.backround_image) {
            $('body').css('background-image', 'url(backgrounds/' + config.backround_image + ')');
          }

        }
      })
    }
    function loadData() {
      $.ajax({
        url: 'getAllData.php',
        success: (data) => {
          render(data);

          $('button.editList').click((e) => {
            const id = e.target.id.replace('editList_', '');
            $('#editListDialog form').attr('action', 'editList.php?id=' + id);
            $('#listTitle').val(data.find(l => l.id == id).title);
            $('#listRows').val(data.find(l => l.id == id).number_of_rows);
            $('#editListDialog').dialog('open');
          });
          $('#addList').click(() => {
            $('#editListDialog form').attr('action', 'addList.php');
            $('#listTitle').val('');
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
            const itemId = e.target.id.replace('editItem_', '');
            $.ajax({
              url: 'getItem.php?id=' + itemId,
              success: (item) => {
                $('#editItemDialog form').attr('action', 'editItem.php?itemId=' + item.id);
                $('#itemTitle').val(item.title);
                $('#itemHref').val(item.href);
                $('#itemIcon').val(item.icon);
                $('#editItemDialog').dialog('open');
              }
            });
          });
          $('div.addItem button').click((e) => {
            const listId = e.target.id.replace('addItem_', '');
            $('#editItemDialog form').attr('action', 'addItem.php?listId=' + listId);
            $('#itemTitle').val('');
            $('#itemHref').val('');
            $('#itemIcon').val('');
            $('#editItemDialog').dialog('open');
          });
          $('button.removeItem').click((e) => {
            const itemId = e.target.id.replace('removeItem_', '');
            if (!confirm('Are you sure you want to remove this item?')) {
              return;
            }
            location.href = 'removeItem.php?itemId=' + itemId;
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
      });
    }

    function render(data) {
      data.forEach(l => {
        var list = $('<div class="list"></div>');
        list.append($('<h2>' + l.title + ' <button id="editList_' + l.id + '" class="editList">✍️</button><button id="removeList_' + l.id + '" class="removeList">❌</button></h2>'));
        items = $('<div class="items"></div>');
        items.css('grid-template-columns', 'repeat(' + l.number_of_rows + ', 1fr)');

        l.items.forEach(i => {
          var item = $('<div class="item"></div>');
          item.append($('<button id="removeItem_' + i.id + '" class="removeItem">❌</button><button id="editItem_' + i.id + '" class="editItem">✍️</button>'));

          var link = $('<a class="href" href="' + i.href + '" target="_blank"></a>');
          link.append($('<img class="icon" src="dashboard-icons/' + i.icon + '.png" />'));

          var details = $('<div class="details"></div>');
          details.css('max-width', 'calc(100vw / ' + l.number_of_rows + ' - 135px)');
          details.append($('<div class="title">' + i.title + '</div>'));
          details.append($('<div class="href">' + i.href + '</div>'));

          link.append(details);
          item.append(link);

          items.append(item);
        });

        items.append($('<div class="addItem"><button id="addItem_' + l.id + '">➕ Add Item</button></div>'));

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
        url: 'initDB.php',
        success: (data) => {
          loadConfig();
          loadData();
        }
      });

      updateWeather();

      $('#config').click(() => {
        $.ajax({
          url: 'getConfig.php',
          success: (config) => {
            $('#configTitle').val(config.title);
            $('#configIcon').val(config.icon);
            $('#configBackroundImage').val(config.backround_image);
            $('#configRows').val(config.number_of_rows);
            $('#configDialog').dialog('open');
          }
        })
      });

      $('#configDialog').dialog({
        title: 'Config',
        autoOpen: false,
        modal: true,
        buttons: {
          "Save": () => {
            if ($('#configDialog :invalid').length > 0) {
              return;
            }
            $('#configDialog form').submit();
          },
          "Cancel": () => {
            $('#configDialog').dialog('close');
          }
        }
      });

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
