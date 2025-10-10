package application

import (
	"itemsim-server/internal/domain/exclusive_equip"
)

type exclusiveEquipServiceImpl struct {
	exclusiveEquipRepository exclusive_equip.Repository
}

func NewExclusiveEquipService(exclusiveEquipRepository exclusive_equip.Repository) ExclusiveEquipService {
	return &exclusiveEquipServiceImpl{
		exclusiveEquipRepository: exclusiveEquipRepository,
	}
}

func (s *exclusiveEquipServiceImpl) GetAllDataAsJson() (any, error) {
	data := s.exclusiveEquipRepository.FindAllDataAsJson()
	return data, nil
}
