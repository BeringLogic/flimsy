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
      <img src="https://openweathermap.org/img/wn/10d@2x.png">
      <div class="temp">5°C</div>
    </div>
    <h1>Flimsy Home Page Dev</h1>
  </header>

  <div class="list">
    <h2>List 1 </h2>
    <div class="items">
      <div class="item">
        <img class="icon" src="https://cdn.jsdelivr.net/gh/walkxcode/dashboard-icons/png/dockge-light.png">
        <div class="details">
          <div class="title">Item 1</div>
          <div class="href">https://item1.example.com:8080</div>
        </div>
      </div>
      <div class="item">Item 2</div>
      <div class="item">Item 3</div>
      <div class="item">Item 4</div>
      <div class="item">Item 5</div>
      <div class="item">Item 6</div>
      <div class="item">Item 7</div>
      <div class="item">Item 8</div>
      <div class="item">Item 9</div>
      <div class="addItem"><button id="addItem_1">➕ Add Item</button></div>
    </div>
  </div>

  <button id="addList">➕ Add List</button>

  <div id="addListDialog" class="dialog">
    <form action="addList.php" method="post">
      <div class="dialog-field">
        <label for="listTitle">Title</label>
        <input id="listTitle" type="text" name="title" required>
      </div>
    </form>
  </div>

  <div id="addItemDialog" class="dialog">
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
    function initData() {
      <?php if (!file_exists('/data/data.json')) {
        error_log("Creating /data/data.json...");
        if (!@touch('/data/data.json')) {
          error_log("/data/data.json could not be created! Check permissions. Owner must be www-data.");
          return false;
        }
        file_put_contents('/data/data.json', '[]');
        return true;
      }?>
    }

    function readData() {
      return <?php echo rtrim(file_get_contents("/data/data.json")); ?>;
    }

    function render(data) {
      data.forEach(b => {
        var list = $('<div class="list"></div>');
        list.append($('<h2>' + b.title + '</h2>'));
        items = $('<div class="items"></div>');

        b.items.forEach(i => {
          var item = $('<div class="item"></div>');
          item.append($('<img class="icon" src="' + i.icon + '" />'));

          var details = $('<div class="details"></div>');
          details.append($('<div class="title">' + i.title + '</div>'));
          details.append($('<div class="href">' + i.href + '</div>'));

          item.append(details);
          items.append(item);
        });

        items.append($('<div class="addItem"><button id="addItem_' + b.id + '">➕ Add Item</button></div>'));

        list.append(items);
        list.insertBefore($('#addList'));
      });
    }

    window.onload = () => {
      if (initData() === false) {
        alert("/data is not writable by the www-data user!");
        return;
      }

      $.ajax({
        url: 'getAllData.php',
        success: (data) => {
          render(data);

          $('div.addItem button').click((e) => {
            const id = e.target.id.replace('addItem_', '');
            $('#addItemDialog form').attr('action', 'addItem.php?id=' + id);
            $('#addItemDialog').dialog('open');
          });

          $('#addList').click(() => {
            $('#addListDialog').dialog('open');
          });
        },
        error: (error, status, xhr) => {
          alert("An error occured while loading the data!");
          console.log(error);
          console.log(status);
          console.log(xhr);
        }

      })

      $('#addListDialog').dialog({
        title: 'Add List',
        autoOpen: false,
        modal: true,
        buttons: {
          "Add": () => {
            if ($('#addListDialog :invalid').length > 0) {
              return;
            }
            $('#addListDialog form').submit();
          },
          "Cancel": () => {
            $('#addListDialog').dialog('close');
          }
        }
      });

      $('#addItemDialog').dialog({
        title: 'Add Item',
        autoOpen: false,
        modal: true,
        buttons: {
          "Add": () => {
            if ($('#addItemDialog :invalid').length > 0) {
              return;
            }
            $('#addItemDialog form').submit();
          },
          "Cancel": () => {
            $('#addItemDialog').dialog('close');
          }
        }
      });
    }
  </script>

</body>
</html>
