package application

type ItemService interface {
	GetIconRawOriginById(id string) ([2]int, error)
}
