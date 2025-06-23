package application

import (
	"errors"
	"itemsim-server/internal/domain/item"
)

type itemServiceImpl struct {
	itemRepository item.Repository
}

func NewItemService(itemRepository item.Repository) ItemService {
	return &itemServiceImpl{
		itemRepository: itemRepository,
	}
}

func (s *itemServiceImpl) GetIconRawOriginById(id string) ([2]int, error) {
	origin, found := s.itemRepository.FindIconRawOriginById(id)
	if !found {
		return [2]int{}, errors.New("not found")
	}
	return origin, nil
}

func (s *itemServiceImpl) GetAllIconRawOriginsById(ids []string) ([][2]int, error) {
	data := s.itemRepository.FindAllIconRawOriginsById(ids)
	return data, nil
}
