package db

type Item struct {
  Id int
  List_id int
  Title string
  Href string
  Icon string
  Position int
}

func LoadItems() (*[]Item, error) {
  Items := []Item{}

  rows, err := sqlDb.Query("SELECT * FROM item"); if err != nil {
    return nil, err
  }

  for rows.Next() {
    var item Item
    err = rows.Scan(&item.Id, &item.List_id, &item.Title, &item.Href, &item.Icon, &item.Position); if err != nil {
      return nil, err
    }
    Items = append(Items, item)
  }

  return &Items, nil
}

