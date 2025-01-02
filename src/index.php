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
          item.append($('<button id="removeItem_' + b.id + '" class="removeItem">❌</button><button id="editItem_' + b.id + '_' + i.id + '" class="editItem">✍️</button>'));
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
      $.ajax({
        url: 'getAllData.php',
        success: (data) => {
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

        },
        error: (error, status, xhr) => {
          alert("An error occured while loading the data! Make sure /data is writable by the www-data user (UID 33, GID 33).");
          console.log(error);
          console.log(status);
          console.log(xhr);
        }

      })

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
