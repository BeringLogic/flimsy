<html>
<head>
  <title>Flimsy</title>
  <style>
    body {
      background-color: #0d1117;
      color: #fff;
    }
    div.block {
      margin: 50px; 
    }
      div.block h2 {
      }
      div.block div.items {
        display: grid;
        grid-template-columns: repeat(4, 1fr);
        justify-content: space-evenly;
        gap: 10px;
      }
        div.item {
          border: 3px solid darkgreen;
          border-radius: 5px;
          padding: 5px;
        }
        div.item img.icon {
          display: inline-block;
          width: 64px;
          height: 64px;
          padding: 5px;
          vertical-align: -12px;
        }
          div.item div.details {
          display: inline-block;
        }
        div.item div.title {
          font-size: xx-large;
        }
        div.item div.desc {
          font-size: large;
        }
  </style>
</head>

<body>
  <h1>Flimsy Home Page Dev</h1>

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


</div>

</body>
</html>
