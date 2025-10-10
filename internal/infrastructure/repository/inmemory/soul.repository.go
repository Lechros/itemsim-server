package inmemory

import (
	"itemsim-server/internal/config"
	"itemsim-server/internal/domain/soul"
	"itemsim-server/internal/infrastructure/file"
)

type soulRepository struct {
	dataJson any
}

func NewSoulRepository(config *config.Config) (soul.Repository, error) {
	var dataJson any

	if err := file.ReadJson(config.GetFilePath("soul.json"), &dataJson); err != nil {
		return nil, err
	}

	return &soulRepository{
		dataJson: dataJson,
	}, nil
}

func (r *soulRepository) FindAllDataAsJson() any {
	return r.dataJson
}
