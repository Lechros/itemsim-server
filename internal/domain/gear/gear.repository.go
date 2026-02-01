package gear

type Repository interface {
	FindAll() []Gear

	Count() int

	FindDataById(id int) (map[string]interface{}, bool)

	FindAllDataById(ids []int) ([]map[string]interface{}, error)

	FindHashById(id int) (string, bool)

	FindAllHashesById(ids []int) ([]string, error)

	FindIconOriginById(id int) ([2]int, bool)

	FindAllIconOriginsById(ids []int) [][2]int
}
