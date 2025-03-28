package storage

import (
	"errors"
	"slices"

	"github.com/BeringLogic/flimsy/internal/db"
)


type FlimsyStorage struct {
  db *db.FlimsyDB
  Config *db.Config
  Sessions []*db.Session
  Lists []*db.List
  Items []*db.Item
  AllListsAndItems []*listAndItems
}

type listAndItems struct {
  List *db.List
  Items []*db.Item
}


var OpenError error = errors.New("Failed to open DB")
var SeedError error = errors.New("Failed to seed DB")
var LoadConfigError error = errors.New("Failed to load config")
var DeleteExpiredSessionsError error = errors.New("Failed to delete expired sessions")
var LoadSessionsError error = errors.New("Failed to load sessions")
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

  if storage.Config, err = storage.db.LoadConfig(); err != nil {
    return errors.Join(LoadConfigError, err);
  }

  if err = storage.db.DeleteExpiredSessions(); err != nil {
    return errors.Join(DeleteExpiredSessionsError, err);
  }
  if storage.Sessions, err = storage.db.LoadSessions(); err != nil {
    return errors.Join(LoadSessionsError, err);
  }

  if storage.Lists, err = storage.db.LoadLists(); err != nil {
      return errors.Join(LoadListsError, err);
  }

  if storage.Items, err = storage.db.LoadItems(); err != nil {
      return errors.Join(LoadItemsError, err);
  }

  storage.AllListsAndItems = storage.getAllListsAndItems()

  return nil
}

func (flimsyStorage *FlimsyStorage) getAllListsAndItems() []*listAndItems {
  AllListsAndItems := make([]*listAndItems, 0)
  for _, list := range flimsyStorage.Lists {
    lai := listAndItems {
      List: list,
      Items: make([]*db.Item, 0),
    }

    for _, item := range flimsyStorage.Items {
      if item.List_id == list.Id {
        lai.Items = append(lai.Items, item)
      }
    }

    AllListsAndItems = append(AllListsAndItems, &lai)
  }
  
  return AllListsAndItems
}

func (flimsyStorage *FlimsyStorage) Close() {
  flimsyStorage.db.Close()
}

func (flimsyStorage *FlimsyStorage) CreateNewSession() (*db.Session, error) {
  session, err := flimsyStorage.db.CreateNewSession(); if err != nil {
    return nil, err
  }

  flimsyStorage.Sessions = append(flimsyStorage.Sessions, session)

  return session, nil
}

func (flimsyStorage *FlimsyStorage) CheckSessionToken(tokenToCheck string) bool {
  for _, session := range flimsyStorage.Sessions {
    if session.Token == tokenToCheck && !session.IsExpired() {
      return true
    }
  }

  return false
}

func (flimsyStorage *FlimsyStorage) DeleteSession(sessionToken string) error {
  for i, session := range flimsyStorage.Sessions {
    if session.Token == sessionToken {
      if err := flimsyStorage.db.DeleteSession(sessionToken); err != nil {
        return err
      }

      flimsyStorage.Sessions = slices.Delete(flimsyStorage.Sessions, i, i+1)

      return nil
    }
  }

  return errors.New("Session not found in cache")
}

func (flimsyStorage *FlimsyStorage) SaveConfig() error {
  return flimsyStorage.db.SaveConfig(flimsyStorage.Config)
}

func (flimsyStorage *FlimsyStorage) AddList(title string, number_of_cols int) (*listAndItems, error) {
  list, err := flimsyStorage.db.AddList(title, number_of_cols); if err != nil {
    return nil, err
  }

  flimsyStorage.Lists = append(flimsyStorage.Lists, list)

  lai := listAndItems {
    List: list,
    Items: make([]*db.Item, 0),
  }
  flimsyStorage.AllListsAndItems = append(flimsyStorage.AllListsAndItems, &lai)

  return &lai, nil
}

func (flimsyStorage *FlimsyStorage) GetList(id int64) (*db.List, error) {
  for _, list := range flimsyStorage.Lists {
    if list.Id == id {
      return list, nil
    }
  }

  return nil, errors.New("List not found")
}

func (flimsyStorage *FlimsyStorage) SaveList(list *db.List) (*listAndItems, error) {
  if err := flimsyStorage.db.SaveList(list); err != nil {
    return nil, err
  }

  var lai *listAndItems
  for _, n := range flimsyStorage.AllListsAndItems {
    if n.List.Id == list.Id {
      lai = n
    }
  }

  return lai, nil
}

func (flimsyStorage *FlimsyStorage) DeleteList(id int64) error {
  for i, n := range flimsyStorage.AllListsAndItems {
    if n.List.Id == id {
      for _, item := range n.Items {
        if err := flimsyStorage.DeleteItem(item.Id); err != nil {
          return err
        }
      }
      if err := flimsyStorage.db.DeleteList(id); err != nil {
        return err
      }

      flimsyStorage.AllListsAndItems = slices.Delete(flimsyStorage.AllListsAndItems, i, i+1)
      flimsyStorage.Lists = slices.Delete(flimsyStorage.Lists, i, i+1)
      break
    }
  }

  return nil
}

func (flimsyStorage *FlimsyStorage) ReorderLists(list_ids []int64) error {
  if err := flimsyStorage.db.ReorderLists(list_ids); err != nil {
    return err
  }

  var loadListErr error
  if flimsyStorage.Lists, loadListErr = flimsyStorage.db.LoadLists(); loadListErr != nil {
    return errors.Join(LoadListsError, loadListErr);
  }

  var loadItemsErr error
  if flimsyStorage.Items, loadItemsErr = flimsyStorage.db.LoadItems(); loadItemsErr  != nil {
    return errors.Join(LoadItemsError, loadItemsErr);
  }

  flimsyStorage.AllListsAndItems = flimsyStorage.getAllListsAndItems()

  return nil
}

func (flimsyStorage *FlimsyStorage) AddItem(list_id int64, title string, url string, icon string, skip_certificate_verification int, check_url string, shortcut string) (*db.Item, error) {
  item, err := flimsyStorage.db.AddItem(list_id, title, url, icon, skip_certificate_verification, check_url, shortcut); if err != nil {
    return nil, err
  }

  flimsyStorage.Items = append(flimsyStorage.Items, item)

  for _, lislistAndItems := range flimsyStorage.AllListsAndItems {
    if lislistAndItems.List.Id == list_id {
      lislistAndItems.Items = append(lislistAndItems.Items, item)
    }
  }

  return item, nil
}

func (flimsyStorage *FlimsyStorage) GetItem(id int64) (*db.Item, error) {
  for _, item := range flimsyStorage.Items {
    if item.Id == id {
      return item, nil
    }
  }

  return nil, errors.New("Item not found")
}

func (flimsyStorage *FlimsyStorage) SaveItem(item *db.Item) error {
  return flimsyStorage.db.SaveItem(item)
}

func (flimsyStorage *FlimsyStorage) DeleteItem(id int64) error {
  if err := flimsyStorage.db.DeleteItem(id); err != nil {
    return err
  }

  for _, listAndItems := range flimsyStorage.AllListsAndItems {
    for i, item := range listAndItems.Items {
      if item.Id == id {
        listAndItems.Items = slices.Delete(listAndItems.Items, i, i+1)
        return nil
      }
    }
  }

  return errors.New("Item not found in cache")
}

func (flimsyStorage *FlimsyStorage) ReorderItems(list_id int64, item_ids []int64) error {
  err := flimsyStorage.db.ReorderItems(list_id, item_ids); if err != nil {
    return err
  }

  if flimsyStorage.Items, err = flimsyStorage.db.LoadItems(); err != nil {
    return errors.Join(LoadItemsError, err);
  }

  flimsyStorage.AllListsAndItems = flimsyStorage.getAllListsAndItems()

  return nil
}
