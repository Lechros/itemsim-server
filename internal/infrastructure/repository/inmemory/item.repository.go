package inmemory

import (
	"itemsim-server/internal/config"
	"itemsim-server/internal/domain/item"
	"itemsim-server/internal/infrastructure/file"
)

type itemRepository struct {
	iconRawOriginMap map[string][2]int
}

func NewItemRepository(config *config.Config) (item.Repository, error) {
	iconRawOriginMap := map[string][2]int{}

	if err := file.ReadJson(config.GetFilePath("item-raw-icon-origin.json"), &iconRawOriginMap); err != nil {
		return nil, err
	}

	return &itemRepository{
		iconRawOriginMap: iconRawOriginMap,
	}, nil
}

func (r *itemRepository) FindIconRawOriginById(id string) ([2]int, bool) {
	origin, ok := r.iconRawOriginMap[id]
	return origin, ok
}

func (r *itemRepository) FindAllIconRawOriginsById(ids []string) [][2]int {
	result := make([][2]int, len(ids))
	for i, id := range ids {
		origin, found := r.iconRawOriginMap[id]
		if found {
			result[i] = origin
		}
	}
	return result
}
