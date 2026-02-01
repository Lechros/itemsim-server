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

func (s *gearServiceImpl) SearchByName(query string, prefix *int) ([]GearSearchResult, error) {
	cmp := func(a gear.Gear, b gear.Gear) int {
		return a.Id - b.Id
	}
	var filter search.ItemFilter[gear.Gear]
	if prefix != nil {
		if *prefix <= 0 {
			return []GearSearchResult{}, errors.New("invalid type")
		}
		filter = func(g gear.Gear) bool {
			return isPrefix(g.Id, *prefix)
		}
	}
	searched := s.searcher.Search(query, 100, cmp, filter)
	results := make([]GearSearchResult, len(searched))
	for i, item := range searched {
		results[i] = GearSearchResult{
			Id:        item.Item.Id,
			Name:      item.Item.Name,
			Icon:      item.Item.Icon,
			Highlight: item.Highlight,
		}
	}
	return results, nil
}

func (s *gearServiceImpl) GetDataById(id int) (map[string]interface{}, error) {
	data, found := s.gearRepository.FindDataById(id)
	if !found {
		return nil, errors.New("not found")
	}
	return data, nil
}

func (s *gearServiceImpl) GetAllDataById(ids []int) ([]map[string]interface{}, error) {
	data, err := s.gearRepository.FindAllDataById(ids)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *gearServiceImpl) GetHashById(id int) (string, error) {
	data, found := s.gearRepository.FindHashById(id)
	if !found {
		return "", errors.New("not found")
	}
	return data, nil
}

func (s *gearServiceImpl) GetAllHashesById(ids []int) ([]string, error) {
	data, err := s.gearRepository.FindAllHashesById(ids)
	if err != nil {
		return nil, err
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

func (s *gearServiceImpl) GetAllIconOriginsById(ids []int) ([][2]int, error) {
	data := s.gearRepository.FindAllIconOriginsById(ids)
	return data, nil
}

func isPrefix(n int, prefix int) bool {
	for n >= prefix {
		if n == prefix {
			return true
		}
		n /= 10
	}
	return false
}
