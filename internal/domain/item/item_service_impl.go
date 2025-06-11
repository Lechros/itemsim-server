package item

import "errors"

type serviceImpl struct {
	itemRepository Repository
}

func NewItemService(itemRepository Repository) Service {
	return &serviceImpl{
		itemRepository: itemRepository,
	}
}

func (s *serviceImpl) GetIconRawOriginById(id string) ([2]int, error) {
	origin, found := s.itemRepository.FindIconRawOriginById(id)
	if !found {
		return [2]int{}, errors.New("not found")
	}
	return origin, nil
}
