package application

import (
	"errors"
	"itemsim-server/internal/common/search"
	"itemsim-server/internal/domain/gear"
)

type gearServiceImpl struct {
	gearRepository gear.Repository
	searcher       search.Searcher[gear.Gear]
}

func NewGearService(gearRepository gear.Repository, searcher search.Searcher[gear.Gear]) GearService {
	service := gearServiceImpl{
		gearRepository: gearRepository,
		searcher:       searcher,
	}
	for _, g := range gearRepository.FindAll() {
		service.searcher.Add(g, g.Name)
	}
	return &service
}

func (s *gearServiceImpl) SearchByName(query string) []GearSearchResult {
	searched := s.searcher.Search(query, 100, func(a gear.Gear, b gear.Gear) int {
		return a.Id - b.Id
	})
	results := make([]GearSearchResult, len(searched))
	for i, item := range searched {
		results[i] = GearSearchResult{
			Id:        item.Item.Id,
			Name:      item.Item.Name,
			Icon:      item.Item.Icon,
			Highlight: item.Highlight,
		}
	}
	return results
}

func (s *gearServiceImpl) GetDataById(id int) (map[string]interface{}, error) {
	data, found := s.gearRepository.FindDataById(id)
	if !found {
		return nil, errors.New("not found")
	}
	return data, nil
}

func (s *gearServiceImpl) GetIconOriginById(id int) ([2]int, error) {
	origin, found := s.gearRepository.FindIconOriginById(id)
	if !found {
		return [2]int{}, errors.New("not found")
	}
	return origin, nil
}
