package storage


import (
  "fmt"
  "github.com/BeringLogic/flimsy/internal/db"
)


type FlimsyStorage struct {
  db *db.FlimsyDB
  Config *db.Config
  Lists *[]db.List
  Items *[]db.Item
  ListsAndItems []listAndItems
}

type listAndItems struct {
  List *db.List
  Items []*db.Item
}


func CreateNew() *FlimsyStorage {
  return &FlimsyStorage {
    db: db.CreateNew(),
  }
}

func (storage *FlimsyStorage) Init() error {
  var err error

  if err = storage.db.Open(); err != nil {
    return fmt.Errorf("Failed to open DB: %s", err)
  }

  storage.Config, err = storage.db.LoadConfig(); if err != nil {
    if err = storage.db.Seed(); err != nil {
      return fmt.Errorf("Failed to seed DB: %s", err)
    }
    if storage.Config, err = storage.db.LoadConfig(); err != nil {
      return fmt.Errorf("Failed to load config: %s", err)
    }
  }

  if storage.Lists, err = storage.db.LoadLists(); err != nil {
    return fmt.Errorf("Failed to load lists: %s", err)
  }

  if storage.Items, err = storage.db.LoadItems(); err != nil {
    return fmt.Errorf("Failed to load items: %s", err)
  }

  for _, list := range *storage.Lists {
    var lai listAndItems
    lai.List = &list
    for _, item := range *storage.Items {
      if item.List_id == list.Id {
        lai.Items = append(lai.Items, &item)
      }
    }
    storage.ListsAndItems = append(storage.ListsAndItems, lai)
  }

  return nil
}

func (flimsyStorage *FlimsyStorage) Close() {
  flimsyStorage.db.Close()
}
