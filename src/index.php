<!DOCTYPE html>
<head>
  <title>Flimsy</title>
  <link rel="stylesheet" href="style.css">
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
    <h2>List 1</h2>
    <div class="items">
      <div class="item">
        <img class="icon" src="https://cdn.jsdelivr.net/gh/walkxcode/dashboard-icons/png/dockge-light.png">
        <div class="details">
          <div class="title">Item 1</div>
          <div class="desc">https://item1.example.com:8080</div>
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
    </div>
  </div>

  <div id="addList" class="list">➕</div>

  <div id="addListDialog" class="dialog">
    <h2>Add List</h2>
    <div class="dialog-field">
      <label for="title">Title</label>
      <input type="text" id="title">
    </div>
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
      return <?php echo file_get_contents('/data/data.json'); ?>;
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
          details.append($('<div class="desc">' + i.desc + '</div>'));

          item.append(details);
          items.append(item);
        });

        list.append(items);
        list.insertBefore($('#addList'));
      });
    }

    window.onload = () => {
      if (initData() === false) {
        alert("/data is not writable by the www-data user!");
        return;
      }
      const data = readData();
      console.log(data);
      render(data);
    }
  </script>

</body>
</html>
