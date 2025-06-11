package inmemory

import (
	"itemsim-server/internal/domain/gear"
	util "itemsim-server/internal/infrastructure/file"
)

type gearRepository struct {
	dataMap       map[int]map[string]interface{}
	iconOriginMap map[int][2]int
}

func NewGearRepository() gear.Repository {
	dataMap := map[int]map[string]interface{}{}
	iconOriginMap := map[int][2]int{}
	util.ReadJson("resources/gear-data.json", &dataMap)
	util.ReadJson("resources/gear-origin.json", &iconOriginMap)
	return &gearRepository{
		dataMap:       dataMap,
		iconOriginMap: iconOriginMap,
	}
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

func (r *gearRepository) FindIconOriginById(id int) ([2]int, bool) {
	origin, ok := r.iconOriginMap[id]
	return origin, ok
}
