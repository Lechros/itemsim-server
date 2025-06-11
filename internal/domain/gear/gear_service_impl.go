package gear

import (
	"errors"
	"itemsim-server/internal/common/search"
)

type serviceImpl struct {
	gearRepository Repository
	searcher       search.Searcher[Gear]
}

func NewGearService(gearRepository Repository, searcher search.Searcher[Gear]) Service {
	service := serviceImpl{
		gearRepository: gearRepository,
		searcher:       searcher,
	}
	for _, g := range gearRepository.FindAll() {
		service.searcher.Add(g, g.Name)
	}
	return &service
}

func (s *serviceImpl) SearchByName(query string) []GearSearchResult {
	searched := s.searcher.Search(query, 100, func(a Gear, b Gear) int {
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

func (s *serviceImpl) GetDataById(id int) (map[string]interface{}, error) {
	data, found := s.gearRepository.FindDataById(id)
	if !found {
		return nil, errors.New("not found")
	}
	return data, nil
}

func (s *serviceImpl) GetIconOriginById(id int) ([2]int, error) {
	origin, found := s.gearRepository.FindIconOriginById(id)
	if !found {
		return [2]int{}, errors.New("not found")
	}
	return origin, nil
}
