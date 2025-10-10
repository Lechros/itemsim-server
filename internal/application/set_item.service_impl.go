package application

import "itemsim-server/internal/domain/set_item"

type setItemServiceImpl struct {
	setItemRepository set_item.Repository
}

func NewSetItemService(setItemRepository set_item.Repository) SetItemService {
	return &setItemServiceImpl{
		setItemRepository: setItemRepository,
	}
}

func (s *setItemServiceImpl) GetAllDataAsJson() (any, error) {
	data := s.setItemRepository.FindAllDataAsJson()
	return data, nil
}
