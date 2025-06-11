package gear

type Repository interface {
	FindAll() []Gear

	Count() int

	FindDataById(id int) (map[string]interface{}, bool)

	FindIconOriginById(id int) ([2]int, bool)
}
