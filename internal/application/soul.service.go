package application

type SoulService interface {
	SearchByName(query string) ([]SoulSearchResult, error)

	GetDataById(id int) (map[string]interface{}, error)

	GetAllDataAsJson() (any, error)
}

type SoulSearchResult struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Highlight string `json:"highlight"`
}
