package db

type List struct {
  ID int
  title string
  number_of_rows int
  position int
}

func LoadLists() (*[]List, error) {
  Lists := []List{}

  rows, err := sqlDb.Query("SELECT * FROM list"); if err != nil {
    return nil, err
  }

  for rows.Next() {
    var list List
    err = rows.Scan(&list.ID, &list.title, &list.number_of_rows, &list.position); if err != nil {
      return nil, err
    }
    Lists = append(Lists, list)
  }

  return &Lists, nil
}

