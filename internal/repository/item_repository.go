package repository

import util "itemsim-server/internal"

var itemRawOrigins map[string][2]int

func GetItemRawIconOriginById(id string) ([2]int, bool) {
	origin, ok := itemRawOrigins[id]
	return origin, ok
}

func init() {
	util.ReadJson("resources/item-raw-origin.json", &itemRawOrigins)
}
