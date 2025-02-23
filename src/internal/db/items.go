package db

type Item struct {
  Id int
  List_id int
  Title string
  Url string
  Icon string
  Position int
}

func (flimsyDB *FlimsyDB) LoadItems() (*[]Item, error) {
  Items := []Item{}

  rows, err := flimsyDB.sqlDb.Query("SELECT * FROM item"); if err != nil {
    return nil, err
  }

  for rows.Next() {
    var item Item
    if err = rows.Scan(&item.Id, &item.List_id, &item.Title, &item.Url, &item.Icon, &item.Position); err != nil {
      return nil, err
    }
    Items = append(Items, item)
  }

  return &Items, nil
}

