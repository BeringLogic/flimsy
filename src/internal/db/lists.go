package db


type List struct {
  Id int64
  Title string
  Number_of_rows int
  Position int
}


func (flimsyDB *FlimsyDB) LoadLists() (map[int64]*List, error) {
  Lists := make(map[int64]*List);

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

func (flimsyDB *FlimsyDB) AddList(title string, number_of_rows int) (*List, error) {
  row := flimsyDB.sqlDb.QueryRow("SELECT MAX(position) + 1 FROM list")

  var position int
  if err := row.Scan(&position); err != nil {
    return nil, err
  }

  result, err := flimsyDB.sqlDb.Exec("INSERT INTO list (title, number_of_rows, position) VALUES (?, ?, ?)", title, number_of_rows, position); if err != nil {
    return nil, err
  }

  id, err := result.LastInsertId(); if err != nil {
    return nil, err
  }

  list := &List {
    Id: id,
    Title: title,
    Number_of_rows: number_of_rows,
    Position: position,
  }

  return list, nil
}

func (flimsyDB *FlimsyDB) SaveList(list *List) error {
  _, err := flimsyDB.sqlDb.Exec("UPDATE list SET title = ?, number_of_rows = ?, position = ? WHERE id = ?", list.Title, list.Number_of_rows, list.Position, list.Id)
  return err
}

func (flimsyDB *FlimsyDB) DeleteList(id int64) error {
  _, err := flimsyDB.sqlDb.Exec("DELETE FROM list WHERE id = ?", id)
  return err
}
