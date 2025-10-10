package inmemory

import (
	"itemsim-server/internal/config"
	"itemsim-server/internal/domain/set_item"
	"itemsim-server/internal/infrastructure/file"
)

type setItemRepository struct {
	dataJson any
}

func NewSetItemRepository(config *config.Config) (set_item.Repository, error) {
	var dataJson any

	if err := file.ReadJson(config.GetFilePath("set-item.json"), &dataJson); err != nil {
		return nil, err
	}

	return &setItemRepository{
		dataJson: dataJson,
	}, nil
}

func (r *setItemRepository) FindAllDataAsJson() any {
	return r.dataJson
}
