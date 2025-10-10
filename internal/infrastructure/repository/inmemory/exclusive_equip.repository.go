package inmemory

import (
	"itemsim-server/internal/config"
	"itemsim-server/internal/domain/exclusive_equip"
	"itemsim-server/internal/infrastructure/file"
)

type exclusiveEquipRepository struct {
	dataJson any
}

func NewExclusiveEquipRepository(config *config.Config) (exclusive_equip.Repository, error) {
	var dataJson any

	if err := file.ReadJson(config.GetFilePath("exclusive-equip.json"), &dataJson); err != nil {
		return nil, err
	}

	return &exclusiveEquipRepository{
		dataJson: dataJson,
	}, nil
}

func (r *exclusiveEquipRepository) FindAllDataAsJson() any {
	return r.dataJson
}
