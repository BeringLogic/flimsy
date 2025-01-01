<!DOCTYPE html>
<head>
  <title>Flimsy</title>
  <link rel="stylesheet" href="style.css">
</head>

<body>
  <header>
    <div class="weather">
      <img src="https://openweathermap.org/img/wn/10d@2x.png">
      <div class="temp">5°C</div>
    </div>
    <h1>Flimsy Home Page Dev</h1>
  </header>

  <div class="block">
    <h2>Block 1</h2>
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

  <div class="block">
    <h2>Block 2</h2>
    <div class="items">
      <div class="item">Item 1</div>
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

  <script type="text/javascript">
    function initData() {
      <?php if (!file_exists('/data/data.json')) {
        error_log("Creating /data/data.json...");
        if (!@touch('/data/data.json')) {
          error_log("/data/data.json could not be created! Check permissions. Owner must be www-data.");
          return false;
        }
        file_put_contents('/data/data.json', '{}');
        return true;
      }?>
    }
    function readData() {
      return <?php echo json_encode(file_get_contents('/data/data.json')); ?>;
    }
  
    window.onload = () => {
      if (initData() === false) {
        alert("/data is not writable by the www-data user!");
        return;
      }
      const data = readData();
      console.log(data);
    }
  </script>

</body>
</html>
