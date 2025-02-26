package storage

import (
	"errors"

	"github.com/BeringLogic/flimsy/internal/db"
)


type FlimsyStorage struct {
  db *db.FlimsyDB
  Config *db.Config
  Lists map[int]*db.List
  Items map[int]*db.Item
  ListsAndItems map[int]*listAndItems
}

type listAndItems struct {
  List *db.List
  Items map[int]*db.Item
}


var OpenError error = errors.New("Failed to open DB")
var SeedError error = errors.New("Failed to seed DB")
var LoadConfigError error = errors.New("Failed to load config")
var LoadListsError error = errors.New("Failed to load lists")
var LoadItemsError error = errors.New("Failed to load items")


func CreateNew() *FlimsyStorage {
  return &FlimsyStorage {
    db: db.CreateNew(),
  }
}

func (storage *FlimsyStorage) Init() error {
  var err error

  if err = storage.db.Open(); err != nil {
    return errors.Join(OpenError, err)
  }

  storage.Config, err = storage.db.LoadConfig(); if err != nil {
  if err = storage.db.Seed(); err != nil {
      return errors.Join(SeedError, err);
    }
    if storage.Config, err = storage.db.LoadConfig(); err != nil {
      return errors.Join(LoadConfigError, err);
    }
  }

  if storage.Lists, err = storage.db.LoadLists(); err != nil {
      return errors.Join(LoadListsError, err);
  }

  if storage.Items, err = storage.db.LoadItems(); err != nil {
      return errors.Join(LoadItemsError, err);
  }

  storage.ListsAndItems = storage.getAllListsAndItems()

  return nil
}

func (flimsyStorage *FlimsyStorage) getAllListsAndItems() map[int]*listAndItems {
  AllListsAndItems := make(map[int]*listAndItems, 0)
  for _, list := range flimsyStorage.Lists {
    lai := listAndItems {
      List: list,
      Items: make(map[int]*db.Item),
    }

    for _, item := range flimsyStorage.Items {
      if item.List_id == list.Id {
        lai.Items[item.Id] = item
      }
    }

    AllListsAndItems[list.Id] = &lai
  }
  
  return AllListsAndItems
}

func (flimsyStorage *FlimsyStorage) Close() {
  flimsyStorage.db.Close()
}

func (flimsyStorage *FlimsyStorage) SaveConfig() error {
  return flimsyStorage.db.SaveConfig(flimsyStorage.Config)
}

func (flimsyStorage *FlimsyStorage) SaveList(list *db.List) (*listAndItems, error) {
  if err := flimsyStorage.db.SaveList(list); err != nil {
    return nil, err
  }

  return flimsyStorage.ListsAndItems[list.Id], nil
}

func (flimsyStorage *FlimsyStorage) SaveItem(item *db.Item) error {
  return flimsyStorage.db.SaveItem(item)
}
