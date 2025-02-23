package db


type List struct {
  Id int
  Title string
  Number_of_rows int
  Position int
}


func (flimsyDB *FlimsyDB) LoadLists() (*[]List, error) {
  Lists := []List{}

  rows, err := flimsyDB.sqlDb.Query("SELECT * FROM list"); if err != nil {
    return nil, err
  }

  for rows.Next() {
    var list List
    if err = rows.Scan(&list.Id, &list.Title, &list.Number_of_rows, &list.Position); err != nil {
      return nil, err
    }
    Lists = append(Lists, list)
  }

  return &Lists, nil
}

