package db


type Item struct {
  Id int
  List_id int
  Title string
  Url string
  Icon string
  Position int
}


func (flimsyDB *FlimsyDB) LoadItems() (map[int]*Item, error) {
  Items := make(map[int]*Item);

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

func (flimsyDB *FlimsyDB) SaveItem(item *Item) error {
  _, err := flimsyDB.sqlDb.Exec("UPDATE item SET list_id = ?, title = ?, url = ?, icon = ?, position = ? WHERE id = ?", item.List_id, item.Title, item.Url, item.Icon, item.Position, item.Id)
  return err
}
