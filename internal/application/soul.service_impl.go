package application

import (
	"errors"
	"itemsim-server/internal/common/search"
	"itemsim-server/internal/domain/soul"
)

type soulServiceImpl struct {
	soulRepository soul.Repository
	searcher       search.Searcher[soul.Soul]
}

func NewSoulService(soulRepository soul.Repository, searcher search.Searcher[soul.Soul]) SoulService {
	service := &soulServiceImpl{
		soulRepository: soulRepository,
		searcher:       searcher,
	}
	for _, s := range soulRepository.FindAll() {
		service.searcher.Add(s, s.Name)
	}
	return service
}

func (s *soulServiceImpl) SearchByName(query string) ([]SoulSearchResult, error) {
	cmp := func(a soul.Soul, b soul.Soul) int {
		return a.Id - b.Id
	}
	searched := s.searcher.Search(query, 100, cmp, nil)
	results := make([]SoulSearchResult, len(searched))
	for i, item := range searched {
		results[i] = SoulSearchResult{
			Id:        item.Item.Id,
			Name:      item.Item.Name,
			Highlight: item.Highlight,
		}
	}
	return results, nil
}

func (s *soulServiceImpl) GetDataById(id int) (map[string]interface{}, error) {
	data, found := s.soulRepository.FindDataById(id)
	if !found {
		return nil, errors.New("not found")
	}
	return data, nil
}

func (s *soulServiceImpl) GetAllDataAsJson() (any, error) {
	data := s.soulRepository.FindAllDataAsJson()
	return data, nil
}
