package inmemory

import (
	"itemsim-server/internal/config"
	"itemsim-server/internal/domain/soul"
	"itemsim-server/internal/infrastructure/file"
)

type soulRepository struct {
	dataMap map[int]map[string]interface{}
}

func NewSoulRepository(config *config.Config) (soul.Repository, error) {
	dataMap := map[int]map[string]interface{}{}

	if err := file.ReadJson(config.GetFilePath("soul.json"), &dataMap); err != nil {
		return nil, err
	}

	return &soulRepository{
		dataMap: dataMap,
	}, nil
}

func (r *soulRepository) FindAllDataAsJson() any {
	return r.dataMap
}

func (r *soulRepository) FindAll() []soul.Soul {
	souls := make([]soul.Soul, 0, len(r.dataMap))
	for id, data := range r.dataMap {
		souls = append(souls, soul.Soul{
			Id:   id,
			Name: data["name"].(string),
			Data: data,
		})
	}
	return souls
}

func (r *soulRepository) Count() int {
	return len(r.dataMap)
}

func (r *soulRepository) FindDataById(id int) (map[string]interface{}, bool) {
	data, found := r.dataMap[id]
	return data, found
}
