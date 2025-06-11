package search

import (
	"github.com/BurntSushi/rure-go"
	"github.com/Lechros/hangul_regexp"
	"slices"
	"strings"
)

type matchResult[T any] struct {
	Item  T
	Text  string
	Start int
}

// searcherImpl 아이템 목록에서 regexp를 사용해 검색합니다.
type searcherImpl[T any] struct {
	items   []T
	builder strings.Builder // 각 아이템의 text를 연결한 문자열, 사이에 '\n'을 추가함.
	texts   []string        // 각 아이템의 text 목록
	offsets []int           // 각 아이템 text의 text 내에서의 인덱스
}

func NewSearcher[T any](cap int) Searcher[T] {
	return &searcherImpl[T]{
		items:   make([]T, 0, cap),
		builder: strings.Builder{},
		texts:   make([]string, 0, cap),
		offsets: make([]int, 0, cap),
	}
}

func (s *searcherImpl[T]) Add(item T, text string) {
	s.items = append(s.items, item)
	s.offsets = append(s.offsets, s.builder.Len())
	s.texts = append(s.texts, text)
	s.builder.WriteString(text)
	s.builder.WriteRune('\n')
}

func (s *searcherImpl[T]) Search(query string, size int, cmp ItemCmp[T], filter ItemFilter[T]) []SearchResult[T] {
	matched := make([]matchResult[T], 0) // 매치된 아이템 인덱스 목록

	regex := getNonCapturingRegex(query)
	matches := regex.FindAll(s.builder.String())
	lastIndex := -1
	for i, match := range matches {
		if i%2 == 0 { // start offset을 기준으로 검사
			index := s.findIndexForPosition(match)
			if index > lastIndex {
				matched = append(matched, matchResult[T]{
					Item:  s.items[index],
					Text:  s.texts[index],
					Start: match - s.offsets[index],
				})
				lastIndex = index
			}
		}
	}

	if filter != nil {
		i := 0
		for _, item := range matched {
			if filter(item.Item) {
				matched[i] = item
				i++
			}
		}
		matched = matched[:i]
	}

	slices.SortStableFunc(matched, func(a, b matchResult[T]) int {
		// 첫 매치 위치가 빠른 결과 우선
		if a.Start != b.Start {
			return a.Start - b.Start
		}
		// 이후 입력 정렬 기준으로 정렬
		return cmp(a.Item, b.Item)
	})

	size = min(len(matched), size)
	result := make([]SearchResult[T], size)

	capturingRegex := getCapturingRegex(query)
	for i, item := range matched {
		if i == size {
			break
		}
		result[i] = SearchResult[T]{
			Item:      item.Item,
			Text:      item.Text,
			Highlight: getHighlight(item.Text, capturingRegex),
		}
	}

	return result
}

func (s *searcherImpl[T]) findIndexForPosition(position int) int {
	index, found := slices.BinarySearch(s.offsets, position)
	if !found {
		return index - 1
	}
	return index
}

func getHighlight(text string, regex *rure.Regex) string {
	captures := regex.NewCaptures()
	regex.Captures(captures, text)
	builder := strings.Builder{}
	builder.Grow(len(text)) // Generous amount of buffer is faster than utf8.RuneCountInString, or reallocation

	gi := 1
	start, end, _ := captures.Group(gi)
	for i := range text {
		for gi < captures.Len() && i >= end {
			gi++
			start, end, _ = captures.Group(gi)
		}
		if i >= start && i < end {
			builder.WriteRune('1')
		} else {
			builder.WriteRune('0')
		}
	}
	return builder.String()
}

func getNonCapturingRegex(query string) *rure.Regex {
	pattern, err := hangul_regexp.GetPattern(query, false, true, false, false)
	if err != nil {
		panic(err)
	}
	return rure.MustCompile("(?i)" + pattern)
}

func getCapturingRegex(query string) *rure.Regex {
	pattern, err := hangul_regexp.GetPattern(query, false, true, true, true)
	if err != nil {
		panic(err)
	}
	return rure.MustCompile("(?i)" + pattern)
}
