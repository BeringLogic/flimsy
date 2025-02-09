package db

type Item struct {
  ID int
  list_id int
  title string
  href string
  icon string
  position int
}

func LoadItems() (*[]Item, error) {
  Items := []Item{}

  rows, err := sqlDb.Query("SELECT * FROM item"); if err != nil {
    return nil, err
  }

  for rows.Next() {
    var item Item
    err = rows.Scan(&item.ID, &item.list_id, &item.title, &item.href, &item.icon, &item.position); if err != nil {
      return nil, err
    }
    Items = append(Items, item)
  }

  return &Items, nil
}

