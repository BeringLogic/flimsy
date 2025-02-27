package db


type Item struct {
  Id int64
  List_id int64
  Title string
  Url string
  Icon string
  Position int
}


func (flimsyDB *FlimsyDB) LoadItems() (map[int64]*Item, error) {
  Items := make(map[int64]*Item);

  rows, err := flimsyDB.sqlDb.Query("SELECT * FROM item"); if err != nil {
    return nil, err
  }

  for rows.Next() {
    var item Item
    if err = rows.Scan(&item.Id, &item.List_id, &item.Title, &item.Url, &item.Icon, &item.Position); err != nil {
      return nil, err
    }
    Items[item.Id] = &item
  }

  return Items, nil
}

func (flimsyDB *FlimsyDB) AddItem(list_id int64, title string, url string, icon string) (*Item, error) {
  row := flimsyDB.sqlDb.QueryRow("SELECT MAX(position) FROM item")

  var position int
  if err := row.Scan(&position); err != nil {
    return nil, err
  }

  item := &Item {
    List_id: list_id,
    Title: title,
    Url: url,
    Icon: icon,
    Position: position,
  }

  result, err := flimsyDB.sqlDb.Exec("INSERT INTO item (list_id, title, url, icon, position) VALUES (?, ?, ?, ?, ?)", list_id, title, url, icon, position); if err != nil {
    return nil, err
  }

  item.Id, err = result.LastInsertId(); if err != nil {
    return nil, err
  }

  return item, nil
}

func (flimsyDB *FlimsyDB) SaveItem(item *Item) error {
  _, err := flimsyDB.sqlDb.Exec("UPDATE item SET list_id = ?, title = ?, url = ?, icon = ?, position = ? WHERE id = ?", item.List_id, item.Title, item.Url, item.Icon, item.Position, item.Id)
  return err
}

func (flimsyDB *FlimsyDB) DeleteItem(id int64) error {
  _, err := flimsyDB.sqlDb.Exec("DELETE FROM item WHERE id = ?", id)
  return err
}
