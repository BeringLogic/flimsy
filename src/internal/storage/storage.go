package storage

import (
	"errors"

	"github.com/BeringLogic/flimsy/internal/db"
)


type FlimsyStorage struct {
  db *db.FlimsyDB
  Config *db.Config
  Lists map[int64]*db.List
  Items map[int64]*db.Item
  ListsAndItems map[int64]*listAndItems
}

type listAndItems struct {
  List *db.List
  Items map[int64]*db.Item
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

func (flimsyStorage *FlimsyStorage) getAllListsAndItems() map[int64]*listAndItems {
  AllListsAndItems := make(map[int64]*listAndItems, 0)
  for _, list := range flimsyStorage.Lists {
    lai := listAndItems {
      List: list,
      Items: make(map[int64]*db.Item),
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

func (flimsyStorage *FlimsyStorage) AddItem(list_id int64, title string, url string, icon string) (*db.Item, error) {
  item, err := flimsyStorage.db.AddItem(list_id, title, url, icon); if err != nil {
    return nil, err
  }

  flimsyStorage.Items[item.Id] = item
  flimsyStorage.ListsAndItems[list_id].Items[item.Id] = item

  return item, nil
}

func (flimsyStorage *FlimsyStorage) SaveItem(item *db.Item) error {
  return flimsyStorage.db.SaveItem(item)
}

func (flimsyStorage *FlimsyStorage) DeleteItem(id int64) error {
  if err := flimsyStorage.db.DeleteItem(id); err != nil {
    return err
  }

  for _, listAndItems := range flimsyStorage.ListsAndItems {
    if listAndItems.Items[id] != nil {
      delete(listAndItems.Items, id)
    }
  }

  return nil
}
