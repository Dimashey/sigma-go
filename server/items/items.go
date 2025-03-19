package items

import (
	"errors"
	"sync"
)

type ItemsRepositry struct {
	items  map[int]Item
	nextId int
	mutex  *sync.Mutex
}

func NewRepository() *ItemsRepositry {
	items := make(map[int]Item)
	nextId := 1
  return &ItemsRepositry{items: items, nextId: nextId, mutex: &sync.Mutex{}}
}

var (
	ErrNotFound     = errors.New("ITEM NOT FOUND")
	ErrNameRequired = errors.New("NAME IS REQUIRED")
)

func (i *ItemsRepositry) Create(item Item) (Item, error) {
	if item.Name == "" {
		return Item{}, ErrNameRequired
	}

	i.mutex.Lock()
	defer i.mutex.Unlock()

	item.ID = i.nextId
	i.items[i.nextId] = item
	i.nextId++

	return Item{}, nil
}

func (i *ItemsRepositry) GetMany() []Item {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	items := make([]Item, 0, len(i.items))

	for _, item := range i.items {
		items = append(items, item)
	}

	return items
}

func (i *ItemsRepositry) GetOne(id int) (Item, error) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	item, exists := i.items[id]

	if !exists {
		return Item{}, ErrNotFound
	}

	return item, nil
}

func (i *ItemsRepositry) Update(id int, updatedItem Item) (Item, error) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	item, exists := i.items[id]

	if !exists {
		return Item{}, ErrNotFound
	}

	if item.Name == "" {
		return Item{}, ErrNameRequired
	}

	item.Name = updatedItem.Name
	i.items[id] = item

	return item, nil
}

func (i *ItemsRepositry) Delete(id int) error {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	_, exists := i.items[id]

	if !exists {
		return ErrNotFound
	}

	delete(i.items, id)

	return nil
}
