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

  <div id="addBlock" class="block">➕</div>

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
      data.forEach(block => {
        const blockEl = document.createElement('div');
        blockEl.classList.add('block');

        // h2 title
        h2El = document.createElement('h2');
        h2El.innerText = block.title;
        blockEl.appendChild(h2El);

        // div.items
        const itemsEl = document.createElement('div');
        itemsEl.classList.add('items');
        block.items.forEach(item => {
          const itemEl = document.createElement('div');
          itemEl.classList.add('item');

          // img
          const itemImgEl = document.createElement('img');
          itemImgEl.classList.add('icon');
          itemImgEl.src = item.icon;
          // div.details
          const itemDetailsEl = document.createElement('div');
          itemDetailsEl.classList.add('details');
          // div.title
          const itemTitleEl = document.createElement('div');
          itemTitleEl.classList.add('title');
          itemTitleEl.innerText = item.title;
          // div.desc
          const itemDescEl = document.createElement('div');
          itemDescEl.classList.add('desc');
          itemDescEl.innerText = item.desc;

          itemDetailsEl.appendChild(itemTitleEl);
          itemDetailsEl.appendChild(itemDescEl);
          itemEl.appendChild(itemImgEl);
          itemEl.appendChild(itemDetailsEl);

          itemsEl.appendChild(itemEl);
        });

        blockEl.appendChild(itemsEl);
        document.body.appendChild(blockEl);
      })
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
