package application

type GearService interface {
	SearchByName(query string) []GearSearchResult

	GetDataById(id int) (map[string]interface{}, error)

	GetIconOriginById(id int) ([2]int, error)
}

type GearSearchResult struct {
	Id        int
	Name      string
	Icon      string
	Highlight string
}
