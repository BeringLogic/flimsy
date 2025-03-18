package db


type Item struct {
  Id int64
  List_id int64
  Title string
  Url string
  Icon string
  Position int
  Skip_certificate_verification int
  Check_url string
}


func (flimsyDB *FlimsyDB) LoadItems() ([]*Item, error) {
  Items := make([]*Item, 0);

  rows, err := flimsyDB.sqlDb.Query("SELECT * FROM item ORDER BY list_id, position"); if err != nil {
    return nil, err
  }

  for rows.Next() {
    var item Item

    if err = rows.Scan(&item.Id, &item.List_id, &item.Title, &item.Url, &item.Icon, &item.Position, &item.Skip_certificate_verification, &item.Check_url); err != nil {
      return nil, err
    }

    Items = append(Items, &item)
  }

  return Items, nil
}

func (flimsyDB *FlimsyDB) AddItem(list_id int64, title string, url string, icon string, skip_certificate_verification int, check_url string) (*Item, error) {
  row := flimsyDB.sqlDb.QueryRow("SELECT IFNULL(MAX(position),0) FROM item")

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
    Skip_certificate_verification: skip_certificate_verification,
    Check_url: check_url,
  }

  result, err := flimsyDB.sqlDb.Exec("INSERT INTO item (list_id, title, url, icon, position, skip_certificate_verification, check_url) VALUES (?, ?, ?, ?, ?, ?, ?)", list_id, title, url, icon, position, skip_certificate_verification, check_url); if err != nil {
    return nil, err
  }

  item.Id, err = result.LastInsertId(); if err != nil {
    return nil, err
  }

  return item, nil
}

func (flimsyDB *FlimsyDB) SaveItem(item *Item) error {
  _, err := flimsyDB.sqlDb.Exec("UPDATE item SET list_id = ?, title = ?, url = ?, icon = ?, position = ?, skip_certificate_verification = ?, check_url = ? WHERE id = ?", item.List_id, item.Title, item.Url, item.Icon, item.Position, item.Skip_certificate_verification, item.Check_url, item.Id)
  return err
}

func (flimsyDB *FlimsyDB) DeleteItem(id int64) error {
  _, err := flimsyDB.sqlDb.Exec("DELETE FROM item WHERE id = ?", id)
  return err
}

func (flimsyDB *FlimsyDB) ReorderItems(list_id int64, item_ids []int64) error {
  for position, item_id := range item_ids {
    if _, err := flimsyDB.sqlDb.Exec("UPDATE item SET list_id = ?, position = ? WHERE id = ?", list_id, position, item_id); err != nil {
      return err
    }
  }

  return nil
}
