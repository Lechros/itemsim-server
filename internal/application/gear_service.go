package application

type GearService interface {
	SearchByName(query string, prefix *int) ([]GearSearchResult, error)

	GetDataById(id int) (map[string]interface{}, error)

	GetIconOriginById(id int) ([2]int, error)
}

type GearSearchResult struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Icon      string `json:"icon"`
	Highlight string `json:"highlight"`
}
