package item

type Repository interface {
	FindIconRawOriginById(id string) ([2]int, bool)

	FindAllIconRawOriginsById(ids []string) [][2]int
}
