package inmemory

import (
	"errors"
	"itemsim-server/internal/config"
	"itemsim-server/internal/domain/gear"
	"itemsim-server/internal/infrastructure/file"
)

type gearRepository struct {
	dataMap       map[int]map[string]interface{}
	iconOriginMap map[int][2]int
}

func NewGearRepository(config *config.Config) (gear.Repository, error) {
	dataMap := map[int]map[string]interface{}{}
	iconOriginMap := map[int][2]int{}

	if err := file.ReadJson(config.GetFilePath("gear-data.json"), &dataMap); err != nil {
		return nil, err
	}

	if err := file.ReadJson(config.GetFilePath("gear-origin.json"), &iconOriginMap); err != nil {
		return nil, err
	}

	return &gearRepository{
		dataMap:       dataMap,
		iconOriginMap: iconOriginMap,
	}, nil
}

func (r *gearRepository) FindAll() []gear.Gear {
	var gears = make([]gear.Gear, len(r.dataMap))
	i := 0
	for id, data := range r.dataMap {
		gears[i] = gear.Gear{
			Id:   id,
			Name: data["name"].(string),
			Icon: data["icon"].(string),
			Data: data,
		}
		i++
	}
	return gears
}

func (r *gearRepository) Count() int {
	return len(r.dataMap)
}

func (r *gearRepository) FindDataById(id int) (map[string]interface{}, bool) {
	data, ok := r.dataMap[id]
	return data, ok
}

func (r *gearRepository) FindAllDataById(ids []int) ([]map[string]interface{}, error) {
	var result = make([]map[string]interface{}, len(ids))
	for i, id := range ids {
		data, found := r.dataMap[id]
		if !found {
			return nil, errors.New("not found")
		}
		result[i] = data
	}
	return result, nil
}

func (r *gearRepository) FindIconOriginById(id int) ([2]int, bool) {
	origin, ok := r.iconOriginMap[id]
	return origin, ok
}

func (r *gearRepository) FindAllIconOriginsById(ids []int) [][2]int {
	result := make([][2]int, len(ids))
	for i, id := range ids {
		origin, found := r.iconOriginMap[id]
		if found {
			result[i] = origin
		}
	}
	return result
}
