<?php session_start(); ?>

<!DOCTYPE html>
<head>
  <title>Flimsy</title>
  <link rel="stylesheet" href="style.css">
  <link rel="stylesheet" href="https://code.jquery.com/ui/1.14.1/themes/base/jquery-ui.css">
  <link rel="stylesheet" href="https://www.nerdfonts.com/assets/css/webfont.css">
  <link rel="icon" type="image/png" href="/homepage.png">
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
    <div class="right">
      <button id="configButton">⚙️</button>
      <?php if (empty($_SESSION['loggedIn'])) { ?>
        <a class="login" href="#">Login</a>
      <?php } else { ?>
        <a class="logout" href="logout.php">Logout</a>
      <?php } ?>
      <div id="system-info">
      </div>
      <div class="weather">
        <img>
        <div class="location"></div>
        <div class="temp"></div>
        <div class="description"></div>
      </div>
    </div>
    <div>
      <img id="icon">
      <h1 id="title"></h1>
    </div>
  </header>

  <div id="lists" style="clear:both;"></div>

  <button id="addList">➕ Add List</button>

  <div id="configDialog" class="dialog">
    <form action="config/set.php" method="post" enctype="multipart/form-data">
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
            <input id="configShowFreeRam" type="checkbox" style="width:auto;" name="show_free_ram">
            <label for="configShowFreeRam" style="display:inline-block">Show Free RAM</label>
          </div>
          <div class="dialog-field">
            <input id="configShowFreeSwap" type="checkbox" style="width:auto;" name="show_free_swap">
            <label for="configShowFreeSwap" style="display:inline-block">Show Free Swap</label>
          </div>
          <div class="dialog-field">
            <input id="configShowPublicIp" type="checkbox" style="width:auto;" name="show_public_ip">
            <label for="configShowPublicIp" style="display:inline-block">Show Public IP</label>
          </div>
          <div class="dialog-field">
            <input id="configShowFreeSpace" type="checkbox" style="width:auto;" name="show_free_space">
            <label for="configShowFreeSpace" style="display:inline-block">Show Free Space</label>
          </div>
        </fieldset>
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
    <form method="post">
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
    <form method="post">
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
    function deleteList() {
      if (!confirm("Are you sure that you want to delete this list and all its items?")) return;
      
      const id = $('#editListDialog form').attr('action').replace('lists/edit.php?id=', '');
      $.ajax({
        url: 'lists/delete.php?id=' + id,
        success: () => {
          $('#list_' + id).remove();
          $('#editListDialog').dialog('close');
        }
      });
    }
    function deleteItem() {
      if (!confirm("Are you sure that you want to delete this item?")) return;
      
      const id = $('#editItemDialog form').attr('action').replace('items/edit.php?itemId=', '');
      $.ajax({
        url: 'items/delete.php?itemId=' + id,
        success: () => {
          $('#item_' + id).remove();
          $('#editItemDialog').dialog('close');
        }
      });
    }

    function loadData() {
      $.ajax({
        url: 'getAllData.php',
        success: (data) => {
          render(data);

          $('div.list h2').click((e) => {
            <?php if (!empty($_SESSION['loggedIn']) && $_SESSION['loggedIn'] == true) { ?>
              const id = e.currentTarget.parentNode.id.replace('list_', '');
              $('#listTitle').val(data.find(l => l.id == id).title);
              $('#listRows').val(data.find(l => l.id == id).number_of_rows);

              var buttons = $('#editListDialog').dialog('option', 'buttons');
              if (buttons.hasOwnProperty('Delete') == false) {
                buttons = { 'Delete': deleteList, ...buttons };
              }
              $('#editListDialog').dialog('option', 'buttons', buttons);
              $('#editListDialog').dialog('option', 'title', 'Edit List');
              $('#editListDialog form').attr('action', 'lists/edit.php?id=' + id);
              $('#editListDialog').dialog('open');
            <?php } ?>
          })

          $('div.item').click((e) => {
            <?php if (!empty($_SESSION['loggedIn']) && $_SESSION['loggedIn'] == true) { ?>
              if (e.target.id.startsWith('addItem_')) return;

              const itemId = e.currentTarget.id.replace('item_', '');
              $.ajax({
                url: 'items/get.php?id=' + itemId,
                success: (item) => {
                  $('#itemTitle').val(item.title);
                  $('#itemHref').val(item.href);
                  $('#itemIcon').val(item.icon);

                  var buttons = $('#editItemDialog').dialog('option', 'buttons');
                  if (!buttons.hasOwnProperty('Delete')) {
                    buttons = { 'Delete': deleteItem, ...buttons };
                  }
                  $('#editItemDialog').dialog('option', 'buttons', buttons);
                  $('#editItemDialog').dialog('option', 'title', 'Edit Item');
                  $('#editItemDialog form').attr('action', 'items/edit.php?itemId=' + item.id);
                  $('#editItemDialog').dialog('open');
                }
              });
            <?php } else { ?>
              const href = $('#' + e.currentTarget.id + ' div.details div.href').text();
              window.open(href);
            <?php } ?>
          });

          $('#addList').click(() => {
            $('#listTitle').val('');

            var buttons = $('#editListDialog').dialog('option', 'buttons');
            if (buttons.hasOwnProperty('Delete')) {
              delete buttons.Delete;
            }
            $('#editListDialog').dialog('option', 'buttons', buttons);
            $('#editListDialog').dialog('option', 'title', 'Add New List');
            $('#editListDialog form').attr('action', 'lists/add.php');
            $('#editListDialog').dialog('open');
          });

          $('div.addItem button').click((e) => {
            const listId = e.target.id.replace('addItem_', '');
            $('#itemTitle').val('');
            $('#itemHref').val('');
            $('#itemIcon').val('');

            var buttons = $('#editItemDialog').dialog('option', 'buttons');
            if (buttons.hasOwnProperty('Delete')) {
              delete buttons.Delete;
            }
            $('#editItemDialog').dialog('option', 'buttons', buttons);
            $('#editItemDialog').dialog('option', 'title', 'Add New Item');
            $('#editItemDialog form').attr('action', 'items/add.php?listId=' + listId);
            $('#editItemDialog').dialog('open');
          });

          <?php if (!empty($_SESSION['loggedIn']) && $_SESSION['loggedIn'] == true) { ?>
          $(function() {
            $('#lists').sortable({
              helper : 'clone',
              stop: function(event, ui) {
                const listIds = $(ui.item).parent().children("div.list").map((l, el) => el.id.replace('list_', '')).toArray();
                $.ajax({
                  url: 'lists/reorder.php',
                  method: 'POST',
                  data: {
                    listIds: listIds
                  }
                });
              } 
            })
            $("div.list").sortable({
              helper: 'clone',
              items: "div.item",
              connectWith: ".list",
              stop: function(event, ui) {
                const listId = $(ui.item).parent().parent().attr('id').replace('list_', '');
                const itemIds = $(ui.item).parent().children("div.item").map((i, el) => el.id.replace('item_', '')).toArray();

                // Move the addItem button to the end of the list, in case the item was dropped after it
                $('#addItem_' + listId).parent().appendTo($('#addItem_' + listId).parent().parent());

                $.ajax({
                  url: 'items/reorder.php',
                  method: 'POST',
                  data: {
                    listId: listId,
                    itemIds: itemIds
                  }
                });
              }
            });
          });
          $('button').css('display', 'inline-block');
          <?php } else { ?>
          $('.addItem').css('display', 'none');
          <?php } ?>

          loadConfig();
          updateWeather();
          updateSystemInfo();
          updateStatusIcons(data);
        },
        error: (error, status, xhr) => {
          alert("An error occured while loading the data! Make sure /data is writable by the user UID 1000, GID 1000.");
          console.log(error);
          console.log(status);
          console.log(xhr);
        }
      });
    }

    function render(data) {
      data.forEach(l => {
        var list = $('<div id="list_' + l.id + '" class="list"></div>');
        list.append($('<h2>' + l.title + '</h2>'));

        items = $('<div class="items"></div>');
        items.css('grid-template-columns', 'repeat(' + l.number_of_rows + ', 1fr)');

        l.items.forEach(i => {
          var item = $('<div id="item_' + i.id + '" class="item"></div>');
          var iconDiv = $('<div class="iconDiv"></div>');
          iconDiv.append($('<img class="icon" src="/data/icons/' + i.icon + '" />'));
          item.append(iconDiv);

          var details = $('<div class="details"></div>');
          details.append($('<i class="status nf nf-oct-dot_fill"></i>'));
          details.append($('<div class="title">' + i.title + '</div>'));
          details.append($('<div class="href">' + i.href + '</div>'));

          item.append(details);
          items.append(item);
        });

        items.append($('<div class="item addItem"><button id="addItem_' + l.id + '">➕ Add Item</button></div>'));

        list.append(items);

        $('#lists').append(list);
      });

    }

    function loadConfig() {
      $.ajax({
        url: 'config/get.php',
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

          $('.ui-resizable-n, .ui-resizable-e, .ui-resizable-s, .ui-resizable-w, .ui-widget-content').css('background-color', config.color_items);
          $('body').css('background-color', config.color_background);
          $('body, .ui-dialog-title, .ui-dialog-content, .ui-dialog-buttonpane, a.login, a.logout, .item a').css('color', config.color_foreground);
          $('#title, .list h2, #system-info, div.weather, a.login, a.logout, div.details').css('text-shadow', '2px 2px ' + config.color_background);
          $('.ui-dialog-title, .ui-dialog-titlebar, .ui-dialog-content, .ui-dialog-buttonpane').css('background-color', config.color_borders);
          $('.item').each((i, el) => { el.style.backgroundColor = 'rgba(from ' + config.color_items + ' r g b / 0.75)'; });
          $('.ui-dialog, .ui-draggable-handle, .ui-dialog-content fieldset, .ui-widget-content, .item').css('border-color', config.color_items);
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
          q: '<?php echo empty($_SERVER["FLIMSY_WEATHER_LOCATION"]) ? "new york" : $_SERVER["FLIMSY_WEATHER_LOCATION"]; ?>',
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
          $('.weather .location').html(data.name + ', ' + data.sys.country);
          $('.weather .description').html(data.weather[0].description);
          $('.weather .temp').html(Math.round(data.main.temp) + ' °' + units);
        }
      });
    }

    function updateSystemInfo()
    {
      $.ajax({
        url: 'getSystemInfo.php',
        success: (data) => {
          if (data.cpu_temp) {
            $('#system-info').append($('<div><i class="nf nf-oct-cpu"></i><span id="cpu-temp">?</span></div>'));
            if (data.cpu_temp.endsWith('C')) {
              $('#cpu-temp').html(data.cpu_temp);
            }
            else {
              $('#cpu-temp').html(data.cpu_temp + ' °C');
            }
          }
          if (data.free_memory) {
            $('#system-info').append($('<div><i class="nf nf-fa-memory"></i><span id="free-memory">?</span> free</div>'));
            $('#free-memory').html(data.free_memory);
          }
          if (data.free_swap) {
            $('#system-info').append($('<div><i class="nf nf-md-file_swap"></i><span id="free-swap">?</span> free</div>'));
            $('#free-swap').html(data.free_swap);
          }
          if (data.public_ip) {
            $('#system-info').append($('<div><i class="nf nf-md-ethernet"></i><span id="public-ip">?</span></div>'));
            $('#public-ip').html(data.public_ip);
          }
          if (data.storage) {
            data.storage.forEach((d) => {
              var div = $('<div></div>');
              $('<i class="nf nf-md-harddisk"></i>').appendTo(div);
              $('<div class="free-space">' + d.free_space + ' free</div>').appendTo(div);
              $('<div class="mount-point">' + d.mount_point + '</div>').appendTo(div);
              div.appendTo($('#system-info'));
            })
          }
        }
      });
    }

    function updateStatusIcons(data) {
      data.forEach((list) => {
        list.items.forEach((item) => {
          $.ajax({
            url: 'getStatus.php?href=' + item.href,
            success: (status) => {
              const i = $('#item_' + item.id + ' i.status');
              i.attr('title', status)
              if (status == 'No error') {
                i.addClass('online').removeClass('offline');
              }
              else {
                i.addClass('offline').removeClass('online');
              }
            }
          })
        })
      });
    }

    $('document').ready(() => {

      $.ajax({
        url: 'initDB.php',
        success: (data) => {
          if (data.success) {
            loadData();
          }
          else {
            alert(data.error);
          }

          <?php if (!empty($_SESSION['message'])) {
            echo 'alert("' . $_SESSION['message'] . '");';
            unset($_SESSION['message']);
          } ?>
        }
      });

      $('#configButton').click(() => {
        $.ajax({
          url: 'config/get.php',
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
            $('#configShowFreeRam').prop('checked', config.show_free_ram);
            $('#configShowFreeSwap').prop('checked', config.show_free_swap);
            $('#configShowPublicIp').prop('checked', config.show_public_ip);
            $('#configShowFreeSpace').prop('checked', config.show_free_space);
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
          "Cancel": () => {
            $('#configDialog').dialog('close');
          },
          "Save": () => {
            if ($('#configDialog :invalid').length > 0) {
              return;
            }
            $('#configDialog form').submit();
          },
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
          "Cancel": () => {
            $('#loginDialog').dialog('close');
          },
          "Login": () => {
            if ($('#loginDialog :invalid').length > 0) {
              return;
            }
            $('#loginDialog form').submit();
          },
        }
      });
      $('a.login').click(() => {
        if (<?php if (empty($_SERVER["FLIMSY_PASSWORD"])) { echo 'false'; } else { echo 'true'; } ?>) {
          $('#loginDialog').dialog('open');
        }
        else {
          window.location = 'login.php';
          return;
        }
      });

      $('#editListDialog').dialog({
        autoOpen: false,
        modal: true,
        buttons: {
          "Cancel": () => {
            $('#editListDialog').dialog('close');
          },
          "Save": () => {
            if ($('#editListDialog :invalid').length > 0) {
              return;
            }
            $('#editListDialog form').submit();
          },
        }
      });

      $('#editItemDialog').dialog({
        autoOpen: false,
        modal: true,
        buttons: {
          "Cancel": () => {
            $('#editItemDialog').dialog('close');
          },
          "Save": () => {
            if ($('#editItemDialog :invalid').length > 0) {
              return;
            }
            $('#editItemDialog form').submit();
          },
        }
      });

    });
  </script>

</body>
</html>
