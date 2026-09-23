package soul

type Repository interface {
	FindAll() []Soul

	Count() int

	FindDataById(id int) (map[string]interface{}, bool)

	FindAllDataAsJson() any
}
