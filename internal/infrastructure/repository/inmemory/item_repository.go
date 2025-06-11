package inmemory

import (
	"itemsim-server/internal/domain/item"
	"itemsim-server/internal/infrastructure/file"
)

type itemRepository struct {
	iconRawOriginMap map[string][2]int
}

func NewItemRepository() item.Repository {
	iconRawOriginMap := map[string][2]int{}
	file.ReadJson("resources/item-raw-origin.json", &iconRawOriginMap)
	return &itemRepository{
		iconRawOriginMap: iconRawOriginMap,
	}
}

func (r *itemRepository) FindIconRawOriginById(id string) ([2]int, bool) {
	origin, ok := r.iconRawOriginMap[id]
	return origin, ok
}
