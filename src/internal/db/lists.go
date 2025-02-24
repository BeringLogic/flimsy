package db


type List struct {
  Id int
  Title string
  Number_of_rows int
  Position int
}


func (flimsyDB *FlimsyDB) LoadLists() (map[int]*List, error) {
  Lists := make(map[int]*List);

  rows, err := flimsyDB.sqlDb.Query("SELECT * FROM list"); if err != nil {
    return nil, err
  }

  for rows.Next() {
    var list List
    if err = rows.Scan(&list.Id, &list.Title, &list.Number_of_rows, &list.Position); err != nil {
      return nil, err
    }
    Lists[list.Id] = &list
  }

  return Lists, nil
}

func (flimsyDB *FlimsyDB) SaveList(list *List) error {
  _, err := flimsyDB.sqlDb.Exec("UPDATE list SET title = ?, number_of_rows = ?, position = ? WHERE id = ?", list.Title, list.Number_of_rows, list.Position, list.Id)
  return err
}
