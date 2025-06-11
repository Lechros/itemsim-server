package item

type Service interface {
	GetIconRawOriginById(id string) ([2]int, error)
}
