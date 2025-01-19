<?php session_start(); ?>

<!DOCTYPE html>
<head>
  <title>Flimsy</title>
  <link rel="stylesheet" href="style.css">
  <link rel="stylesheet" href="https://code.jquery.com/ui/1.14.1/themes/base/jquery-ui.css">
  <link rel="stylesheet" href="https://www.nerdfonts.com/assets/css/webfont.css">
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
    <div id="system-info">
      <div><i class="nf nf-oct-cpu"></i><span id="cpu-temp">?</span></div>
      <div><i class="nf nf-fa-memory"></i><span id="free-memory">?</span> free</div>
      <div><i class="nf nf-md-swap_horizontal"></i><span id="free-swap">?</span> free</div>
    </div>
    <?php if (empty($_SESSION['loggedIn'])) { ?>
      <a class="login" href="#">Login</a>
    <?php } else { ?>
      <a class="logout" href="logout.php">Logout</a>
    <?php } ?>
    <button id="configButton">⚙️</button>
    <img id="icon">
    <h1 id="title"></h1>
  </header>

  <div id="lists"></div>

  <button id="addList">➕ Add List</button>

  <div id="configDialog" class="dialog">
    <form action="setConfig.php" method="post" enctype="multipart/form-data">
      <div class="dialog-column">
        <fieldset>
          <legend>Header</legend>
          <div class="dialog-field">
            <label for="configIcon">Icon</label>
            <input id="configIcon" type="text" name="icon">
          </div>
          <div class="dialog-field">
            <label for="configTitle">Title</label>
            <input id="configTitle" type="text" name="title">
          </div>
        </fieldset>
        <fieldset>
          <legend>Background</legend>
          <div class="dialog-field">
            <input id="configBackgroundTypeUpload" type="radio" style="width:auto;" name="background_type" value="upload">
            <input id="configBackgroundUpload" type="file" style="width:85%; display:inline-block;" name="background_image">
          </div>
          <div class="dialog-field">
            <input id="configBackgroundTypeKeep" type="radio" style="width:auto;" name="background_type" value="keep" checked>
            <select id="configBackgroundImage" style="width:85%; display:inline-block;" name="background_image"></select>
          </div>
          <div class="dialog-field">
            <input id="configBackgroundTypeNone" type="radio" style="width:auto;" name="background_type" value="none">
            <label for="configBackgroundTypeNone" style="display:inline-block;">None</label>
          </div>
        </fieldset>
      </div>
      <div class="dialog-column">
        <fieldset>
          <legend>Colors</legend>
          <div class="dialog-field">
            <input type="radio" id="autodetect_colors" name="color_type" style="width:auto;" value="autodetect">
            <label for="autodetect_colors" style="display:inline-block;">Autodetect Colors</label>
          </div>
          <div class="dialog-field">
            <input type="radio" id="catppuccin_latte_colors" name="color_type" style="width:auto;" value="catppuccin_latte">
            <label for="catppuccin_latte_colors" style="display:inline-block;">Catppuccin Latte</label>
          </div>
          <div class="dialog-field">
            <input type="radio" id="catppuccin_mocha_colors" name="color_type" style="width:auto;" value="catppuccin_mocha">
            <label for="catppuccin_mocha_colors" style="display:inline-block;">Catppuccin Mocha</label>
          </div>
          <div class="dialog-field">
            <input type="radio" id="manual_colors" name="color_type" style="width:auto;" value="manual" checked>
            <label for="manual_colors" style="display:inline-block;">Manual Colors</label>
          </div>
          <div class="dialog-field">
            <label for="configColorBackground">Background</label>
            <input id="configColorBackground" type="color" name="color_background" required>
          </div>
          <div class="dialog-field">
            <label for="configColorForeground">Foreground</label>
            <input id="configColorForeground" type="color" name="color_foreground" required>
          </div>
          <div class="dialog-field">
            <label for="configColorItems">Items</label>
            <input id="configColorItems" type="color" name="color_items" required>
          </div>
          <div class="dialog-field">
            <label for="configColorBorders">Borders</label>
            <input id="configColorBorders" type="color" name="color_borders" required>
          </div>
        </fieldset>
      </div>
      <div class="dialog-column">
        <fieldset>
          <legend>System Info</legend>
          <div class="dialog-field">
            <label for="configCpuTempSensor">CPU Temp Sensor</label>
            <select id="configCpuTempSensor" name="cpu_temp_sensor">
              <option value="">Don't show</option>
            </select>
          </div>
          <div class="dialog-field">
            <label for="configMountPoints">Mount Points</label>
            <input id="configMountPoints" name="mount_points">
          </div>
        </fieldset>
      </div>
    </form>
  </div>

  <div id="uploadDialog" class="dialog">
    <form action="uploadBackground.php" method="post" enctype="multipart/form-data">
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
    function loadData() {
      $.ajax({
        url: 'getAllData.php',
        success: (data) => {
          render(data);

          <?php if (!empty($_SESSION['loggedIn']) && $_SESSION['loggedIn'] == true) { ?>
          $(function() {
            $('#lists').sortable({
              stop: function(event, ui) {
                const listIds = $(ui.item).parent().children("div.list").map((l, el) => el.id.replace('list_', '')).toArray();
                $.ajax({
                  url: 'reorderLists.php',
                  method: 'POST',
                  data: {
                    listIds: listIds
                  }
                });
              } 
            })
            $("div.list").sortable({
              items: "div.item",
              connectWith: ".list",
              stop: function(event, ui) {
                const listId = $(ui.item).parent().parent().attr('id').replace('list_', '');
                const itemIds = $(ui.item).parent().children("div.item").map((i, el) => el.id.replace('item_', '')).toArray();
                $.ajax({
                  url: 'reorderItems.php',
                  method: 'POST',
                  data: {
                    listId: listId,
                    itemIds: itemIds
                  }
                });
              }
            });
          });
          <?php } ?>

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

          if (<?php echo !empty($_SESSION['loggedIn']) && $_SESSION['loggedIn'] == true ? "true" : "false"; ?>) {
            $('button').css('display', 'inline-block');
          }

          loadConfig();
          updateWeather();
          updateSystemInfo();
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
        var list = $('<div id="list_' + l.id + '" class="list"></div>');
        list.append($('<h2>' + l.title + ' <button id="editList_' + l.id + '" class="editList">✍️</button><button id="removeList_' + l.id + '" class="removeList">❌</button></h2>'));
        items = $('<div class="items"></div>');
        items.css('grid-template-columns', 'repeat(' + l.number_of_rows + ', 1fr)');

        l.items.forEach(i => {
          var item = $('<div id="item_' + i.id + '" class="item"></div>');
          item.append($('<button id="removeItem_' + i.id + '" class="removeItem">❌</button><button id="editItem_' + i.id + '" class="editItem">✍️</button>'));

          var link = $('<a class="href" href="' + i.href + '" target="_blank"></a>');
          link.append($('<img class="icon" src="/data/icons/' + i.icon + '" />'));

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

        $('#lists').append(list);
      });

    }

    function loadConfig() {
      $.ajax({
        url: 'getConfig.php',
        success: (config) => {
          if (config.icon) {
            $('#icon').attr('src', '/data/icons/' + config.icon)
                      .css('display', 'inline-block');
          }
          if (config.title) {
            $('#title').html(config.title);
          }
          if (config.background_image) {
            $('body').css('background-image', 'url(/data/backgrounds/' + config.background_image + ')');
          }
          $('body').css('background-color', config.color_background);
          $('body, a.login, a.logout').css('color', config.color_foreground);
          $('#title, .list h2').css('text-shadow', '1px 1px ' + config.color_background);
          $('.item a').css('color', config.color_foreground);
          $('.item').css('background-color', config.color_items);
          $('.item').css('border-color', config.color_borders);

          if (!config.cpu_temp_sensor) {
            $('#system-info div:first-child()').css('display', 'none');
          }
        }
      });
    }

    function updateWeather() {
      if (<?php echo empty($_SERVER["FLIMSY_WEATHER_API_KEY"]) ? "true" : "false"; ?>) {
        $('.weather').css('display', 'none');
        return;
      }
      $.ajax({
        url: 'https://api.openweathermap.org/data/2.5/weather',
        data: {
          lat : <?php echo empty($_SERVER["FLIMSY_WEATHER_LATITUDE"]) ? 0 : $_SERVER["FLIMSY_WEATHER_LATITUDE"]; ?>,
          lon : <?php echo empty($_SERVER["FLIMSY_WEATHER_LONGITUDE"]) ? 0 : $_SERVER["FLIMSY_WEATHER_LONGITUDE"]; ?>,
          units : '<?php echo empty($_SERVER["FLIMSY_WEATHER_UNITS"]) ? "standard" : $_SERVER["FLIMSY_WEATHER_UNITS"]; ?>',
          lang : '<?php echo empty($_SERVER["FLIMSY_WEATHER_LANGUAGE"]) ? "en" : $_SERVER["FLIMSY_WEATHER_LANGUAGE"]; ?>',
          appid : '<?php echo empty($_SERVER["FLIMSY_WEATHER_API_KEY"]) ? "" : $_SERVER["FLIMSY_WEATHER_API_KEY"]; ?>',
        },
        success: (data) => {
          var units;
          switch ('<?php echo empty($_SERVER["FLIMSY_WEATHER_UNITS"]) ? "standard" : $_SERVER["FLIMSY_WEATHER_UNITS"]; ?>') {
            default:
            case 'standard': units = 'K'; break;
            case 'metric': units = 'C'; break;
            case 'imperial': units = 'F'; break;
          };
          $('.weather img').attr('src', 'https://openweathermap.org/img/wn/' + data.weather[0].icon + '@2x.png');
          $('.weather .description').html(data.weather[0].description);
          $('.weather .temp').html(data.main.temp + ' °' + units);
        }
      });
    }

    function updateSystemInfo()
    {
      $.ajax({
        url: 'getSystemInfo.php',
        success: (data) => {
          if (data.cpu_temp.endsWith('C')) {
            $('#cpu-temp').html(data.cpu_temp);
          }
          else {
            $('#cpu-temp').html(data.cpu_temp + ' °C');
          }
          $('#free-memory').html(data.free_memory);
          $('#free-swap').html(data.free_swap);

          data.disks.forEach((d) => {
            var div = $('<div></div>');
            $('<i class="nf nf-md-harddisk"></i>').appendTo(div);
            $('<div class="free-disk-space">' + d.free_disk_space + ' free</div>').appendTo(div);
            $('<div class="mount-point">' + d.mount_point + '</div>').appendTo(div);
            div.appendTo($('#system-info'));
          })
        }
      });
    }

    $('document').ready(() => {

      $.ajax({
        url: 'initDB.php',
        success: (data) => {
          loadData();

          <?php if (!empty($_SESSION['message'])) {
            echo 'alert("' . $_SESSION['message'] . '");';
            unset($_SESSION['message']);
          } ?>
        }
      });

      $('#configButton').click(() => {
        $.ajax({
          url: 'getConfig.php',
          success: (config) => {
            $('#configTitle').val(config.title);
            $('#configIcon').val(config.icon);

            config.backgrounds.forEach((b) => {
              if (b == '.' || b == '..') return;
              $('#configBackgroundImage').append($('<option value="' + b + '">' + b + '</option>'));
            })
            $('#configBackgroundImage').val(config.background_image);
            $('#configBackgroundImage').trigger('change');
            $('#configColorBackground').val(config.color_background);
            $('#configColorForeground').val(config.color_foreground);
            $('#configColorItems').val(config.color_items);
            $('#configColorBorders').val(config.color_borders);

            for (const [chip, sensors] of Object.entries(config.sensors)) {
              $('#configCpuTempSensor').append($('<optgroup label="' + chip + '"></optgroup>'));
              for (const [sensorName, values] of Object.entries(sensors)) {
                for (const [valueName, value] of Object.entries(values)) {
                  if (valueName.endsWith('_input')) {
                    $('#configCpuTempSensor').append($('<option value="' + sensorName + '">&nbsp;' + sensorName + ': ' + value + '</option>'));
                  }
                }
              }
            }
            $('#configCpuTempSensor').val(config.cpu_temp_sensor);
            $('#configMountPoints').val(config.mount_points);
            $('#configDialog').dialog('open');
          }
        })
      });

      $('#configDialog').dialog({
        title: 'Config',
        autoOpen: false,
        modal: true,
        width: 1010,
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
      $('#configBackgroundUpload').change(() => {
        $('#autodetect_colors').removeAttr('disabled');
        $('#configBackgroundTypeUpload').prop('checked', true);
      });
      $('#configBackgroundImage').change((e) => {
        if (e.target.value != '') { 
          $('#autodetect_colors').removeAttr('disabled');
          $('#configBackgroundTypeKeep').prop('checked', true);
        }
        else {
          $('#autodetect_colors').attr('disabled', 'disabled');
          $('#manual_colors').prop('checked', true);
        }
      });
      $('#autodetect_colors').change(() => {
        if ($('#autodetect_colors').is(':checked')) {
          $('#configDialog input[type=color]').attr('disabled', 'disabled');
        }
      });
      $('#manual_colors').change(() => {
        if ($('#manual_colors').is(':checked')) {
          $('#configDialog input[type=color]').removeAttr('disabled');
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
        if (<?php if (empty($_SERVER) || empty($_SERVER["FLIMSY_PASSWORD"])) { echo 'false'; } else { echo 'true'; } ?>) {
          $('#loginDialog').dialog('open');
        }
        else {
          window.location = 'login.php';
          return;
        }
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
    });
  </script>

</body>
</html>
