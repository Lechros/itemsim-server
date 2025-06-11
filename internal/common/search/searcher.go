package search

type SearchResult[T any] struct {
	Item      T
	Text      string
	Highlight string
}

type ItemCmp[T any] func(a T, b T) int

type ItemFilter[T any] func(item T) bool

type Searcher[T any] interface {
	Add(item T, text string)

	Search(query string, size int, cmp ItemCmp[T], filter ItemFilter[T]) []SearchResult[T]
}
