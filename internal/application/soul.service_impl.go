package application

import "itemsim-server/internal/domain/soul"

type soulServiceImpl struct {
	soulRepository soul.Repository
}

func NewSoulService(soulRepository soul.Repository) SoulService {
	return &soulServiceImpl{
		soulRepository: soulRepository,
	}
}

func (s *soulServiceImpl) GetAllDataAsJson() (any, error) {
	data := s.soulRepository.FindAllDataAsJson()
	return data, nil
}
